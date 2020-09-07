// +build windows

package ptp

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

// Windows Platform specific constants
const (
	TapTool                    = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"
	DriverInf                  = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf"
	TapSuffix                  = ".tap"
	TapID                      = "tap0901"
	MaximumInterfaceNameLength = 128
)

var (
	errorTAPIsNotInstalled          = errors.New("TAP-Windows 9.2x is not installed")
	errorFailedToRemoveInterfaces   = errors.New("Failed to remove TAP interfaces")
	errorFailedToCreateInterface    = errors.New("Failed to create interface")
	errorObjectCreationFailed       = errors.New("Failed to create TAP object")
	errorFailedToRetrieveNetworkKey = errors.New("Failed to retrieve network key from registry")
	errorFailedToQueryInterface     = errors.New("Failed to query network interface")
	errorPreconfigurationFailed     = errors.New("Interface pre-configuration failed")
)

// InitPlatform initializes Windows platform-specific parameters
func InitPlatform() error {
	Info("Initializing Windows Platform")
	if _, err := os.Stat(TapTool); os.IsNotExist(err) {
		Error("TAP-Windows 9.2x is not installed. Go to https://openvpn.net/index.php/open-source/downloads.html and download the latest version. Close P2P now as it will not run properly")
		return errorTAPIsNotInstalled
	}
	// Remove interfaces
	remove := exec.Command(TapTool, "remove", TapID)
	err := remove.Run()
	if err != nil {
		return errorFailedToRemoveInterfaces
	}
	for i := 0; i < 10; i++ {
		adddev := exec.Command(TapTool, "install", DriverInf, TapID)
		err := adddev.Run()
		if err != nil {
			Error("Failed to add TUN/TAP Device: %v", err)
			return errorFailedToCreateInterface
		}
	}

	tap, err := newTAP(GetConfigurationTool(), "127.0.0.1", "00:00:00:00:00:00", "255.255.255.0", DefaultMTU, UsePMTU)
	if err != nil {
		Error("Failed to create TAP object: %s", err)
		return errorObjectCreationFailed
	}

	for i := 0; i < 10; i++ {
		key, err := tap.queryNetworkKey()
		if err != nil {
			Error("Couldn't open Registry Key %s: %s", NetworkKey, err)
			continue
			//return errorFailedToRetrieveNetworkKey
		}
		err = tap.queryAdapters(key)
		if err != nil {
			Error("Failed to query adapters: %s", err)
			syscall.CloseHandle(tap.file)
			continue
			//return errorFailedToQueryInterface
		}
		// Dummy IP address for the interface
		ip := "172." + strconv.Itoa(i) + ".4.100"
		setip := exec.Command("netsh")
		setip.SysProcAttr = &syscall.SysProcAttr{}

		cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, tap.Interface, ip, "255.255.255.0")
		Debug("Executing: %s", cmd)

		setip.SysProcAttr.CmdLine = cmd
		err = setip.Run()
		err2 := syscall.CloseHandle(tap.file)
		if err != nil {
			continue
			//return errorPreconfigurationFailed
		}
		if err2 != nil {
			Error("Failed to close handle: %s", err)
		}
	}
	UsedInterfaces = UsedInterfaces[:0]
	return nil
}

// CheckPermissions return true if started as root/administrator
func HavePrivileges(level int) bool {
	if level != 0 {
		Error("P2P cannot run in daemon mode without Administrator privileges")
		return false
	}
	return true
}

func GetPrivilegesLevel() int {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return 1
	}
	return 0
}

// Syslog provides additional logging to the syslog server
func Syslog(level LogLevel, format string, v ...interface{}) {
	Info("Syslog is not supported on this platform. Please do not use syslog flag.")
}

// func closeInterface(file syscall.Handle) {
// 	err := syscall.CloseHandle(file)
// 	if err != nil {

// 	}
// }

// SetupPlatform will install Windows Service and exit immediatelly
func SetupPlatform(remove bool) {
	// Opening log
	// f, err := os.OpenFile("C:\\ProgramData\\subutai\\log\\service-setup.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	// if err != nil {
	// 	Error( "Failed to open log file")
	// 	f = nil
	// }
	// if f != nil {
	// 	defer f.Close()
	// 	log.SetOutput(f)
	// }

	name := "Subutai P2P"
	desc := "Subutai networking service"
	Info("Setting up Windows Service")

	p2pApp, err := exePath()
	if err != nil {
		Error("Failed to determine path to executable")
		p2pApp = os.Args[0]
	}
	Info("Application: %s", p2pApp)

	manager, err := mgr.Connect()
	if err != nil {
		Error("Failed to open service manager: %s", err)
		os.Exit(1)
	}
	defer manager.Disconnect()

	Info("Opening service manager")
	service, err := manager.OpenService("Subutai P2P")
	if err == nil {
		// Service exists
		if remove {
			restartWindowsService(service, name, true)
			removeWindowsService(service, name)
		} else {
			restartWindowsService(service, name, false)
		}
	} else {
		if !remove {
			installWindowsService(manager, name, p2pApp, desc)
		}
	}
	os.Exit(0)
}

func removeWindowsService(service *mgr.Service, name string) {
	Info("Removing service")
	err := service.Delete()
	if err != nil {
		Error("Failed to remove service: %s", err)
		service.Close()
		os.Exit(15)
	}
	err = eventlog.Remove(name)
	if err != nil {
		Error("Failed to unregister eventlog: %s", err)
		service.Close()
		os.Exit(16)
	}
	Info("Service removed")
	os.Exit(0)
}

func installWindowsService(manager *mgr.Mgr, name, app, desc string) {
	Info("Creating service")
	service, err := manager.CreateService(name, app, mgr.Config{DisplayName: name, Description: desc, StartType: mgr.StartAutomatic}, "service")
	if err != nil {
		Error("Failed to create P2P service: %s", err)
		os.Exit(6)
	}
	defer service.Close()
	Info("Installing service")
	err = eventlog.InstallAsEventCreate(name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		service.Delete()
		Error("SetupEventLogSource() failed: %s", err)
		os.Exit(7)
	}
	Info("Installation complete")
	err = service.Start("service")
	if err != nil {
		Error("Failed to start service: %s", err)
		return
	}
	Info("Service started")
}

func restartWindowsService(service *mgr.Service, name string, noStart bool) {
	Info("Service exists. Stopping")
	status, err := service.Control(svc.Stop)
	if err != nil {
		Error("Failed to get service status on stop: %s", err)

	} else {
		timeout := time.Now().Add(30 * time.Second)
		for status.State != svc.Stopped {
			if timeout.Before(time.Now()) {
				Error("Failed to stop p2p service after timeout")
				service.Close()
				os.Exit(3)
			}
			time.Sleep(time.Millisecond * 300)
			status, err = service.Query()
			if err != nil {
				Error("Couldn't retrieve service status: %s", err)
				service.Close()
				os.Exit(4)
			}
		}
	}
	if !noStart {
		Info("Starting service")
		// Service stopped. Now start it.
		err = service.Start("service")
		if err != nil {
			Error("Failed to start service on restart: %s", err)
			service.Close()
			// TODO Make this non-zero when fix problems with service start
			os.Exit(0)
		}
		service.Close()
		os.Exit(0)
	}
}

func exePath() (string, error) {
	prog := os.Args[0]
	p, err := filepath.Abs(prog)
	if err != nil {
		return "", err
	}
	fi, err := os.Stat(p)
	if err == nil {
		if !fi.Mode().IsDir() {
			return p, nil
		}
		err = fmt.Errorf("%s is directory", p)
	}
	if filepath.Ext(p) == "" {
		p += ".exe"
		fi, err := os.Stat(p)
		if err == nil {
			if !fi.Mode().IsDir() {
				return p, nil
			}
			err = fmt.Errorf("%s is directory", p)
		}
	}
	return "", err
}
