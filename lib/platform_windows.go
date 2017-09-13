// +build windows

package ptp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
	"unicode/utf16"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// UsedInterfaces - List of interfaces currently in use by p2p daemon
	UsedInterfaces []string
)

// Interface - represents a TAP interface in the system
type Interface struct {
	Name      string
	file      syscall.Handle
	Handle    syscall.Handle
	Interface string
	IP        string
	Mask      string
	Mac       string
	Rx        chan []byte
	Tx        chan []byte
}

const (
	CONFIG_DIR          string         = "C:\\"
	NETWORK_KEY         string         = "SYSTEM\\CurrentControlSet\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	ADAPTER_KEY         string         = "SYSTEM\\CurrentControlSet\\Control\\Class\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	NO_MORE_ITEMS       syscall.Errno  = 259
	USERMODE_DEVICE_DIR string         = "\\\\.\\Global\\"
	SYS_DEVICE_DIR      string         = "\\Device\\"
	USER_DEVICE_DIR     string         = "\\DosDevices\\Global\\"
	TAP_SUFFIX          string         = ".tap"
	TAP_ID              string         = "tap0901"
	INVALID_HANDLE      syscall.Handle = 0
	TAP_TOOL                           = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"
	DRIVER_INF                         = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf"
	ConfigDir                          = "C:\\"
)

var (
	TAP_IOCTL_GET_MAC               = TAP_CONTROL_CODE(1, 0)
	TAP_IOCTL_GET_VERSION           = TAP_CONTROL_CODE(2, 0)
	TAP_IOCTL_GET_MTU               = TAP_CONTROL_CODE(3, 0)
	TAP_IOCTL_GET_INFO              = TAP_CONTROL_CODE(4, 0)
	TAP_IOCTL_CONFIG_POINT_TO_POINT = TAP_CONTROL_CODE(5, 0)
	TAP_IOCTL_SET_MEDIA_STATUS      = TAP_CONTROL_CODE(6, 0)
	TAP_IOCTL_CONFIG_DHCP_MASQ      = TAP_CONTROL_CODE(7, 0)
	TAP_IOCTL_GET_LOG_LINE          = TAP_CONTROL_CODE(8, 0)
	TAP_IOCTL_CONFIG_DHCP_SET_OPT   = TAP_CONTROL_CODE(9, 0)
	TAP_IOCTL_CONFIG_TUN            = TAP_CONTROL_CODE(10, 0)
)

func InitPlatform() {
	// Remove interfaces
	remove := exec.Command(TAP_TOOL, "remove", TAP_ID)
	err := remove.Run()
	if err != nil {
		Log(Error, "Failed to remove TAP interfaces")
	}

	for i := 0; i < 10; i++ {
		adddev := exec.Command(TAP_TOOL, "install", DRIVER_INF, TAP_ID)
		err := adddev.Run()
		if err != nil {
			Log(Error, "Failed to add TUN/TAP Device: %v", err)
		}
	}

	/*remdev := exec.Command(TAP_TOOL, "")
	err := remdev.Run()
	if err != nil {
		Log(Error, "Failed to remove TUN/TAP Devices: %v", err)
	}*/
	/*
		for i := 0; i < 10; i++ {
			adddev := exec.Command(ADD_DEV)
			err := adddev.Run()
			if err != nil {
				Log(Error, "Failed to add TUN/TAP Device: %v", err)
			}
		}
	*/
	for i := 0; i < 10; i++ {
		key, err := queryNetworkKey()
		if err != nil {
			Log(Error, "Failed to retrieve network key: %v", err)
			continue
		}
		inf, err := queryAdapters(key)
		if err != nil || inf == nil {
			Log(Error, "Failed to query interface: %v", err)
			continue
		}
		ip := "172." + strconv.Itoa(i) + ".4.100"
		setip := exec.Command("netsh")
		setip.SysProcAttr = &syscall.SysProcAttr{}

		cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, inf.Interface, ip, "255.255.255.0")
		Log(Info, "Executing: %s", cmd)

		setip.SysProcAttr.CmdLine = cmd
		err = setip.Run()
		syscall.CloseHandle(inf.file)
		if err != nil {
			Log(Error, "Failed to properly preconfigure TAP device with netsh: %v", err)
			continue
		}
	}
	UsedInterfaces = nil
}

func TAP_CONTROL_CODE(request, method uint32) uint32 {
	return CTL_CODE(34, request, method, 0)
}

func CTL_CODE(device_type, function, method, access uint32) uint32 {
	return (device_type << 16) | (access << 14) | (function << 2) | method
}

func removeZeroes(s string) string {
	bytes := []byte(s)
	var (
		res  []byte
		prev bool
		size int
	)
	for _, b := range bytes {
		if b == 0 && prev {
			break
		} else if b == 0 && !prev {
			prev = true
		} else {
			prev = false
			res = append(res, b)
			size++
		}
	}
	return string(res[:size])
}

func queryNetworkKey() (syscall.Handle, error) {
	var handle syscall.Handle
	err := syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(NETWORK_KEY), 0, syscall.KEY_READ, &handle)
	if err != nil {
		return 0, err
	}
	return handle, nil
}

func queryAdapters(handle syscall.Handle) (*Interface, error) {
	var dev Interface
	var index uint32
	for {
		var length uint32 = 72
		adapter := make([]uint16, length)
		err := syscall.RegEnumKeyEx(handle, index, &adapter[0], &length, nil, nil, nil, nil)
		if err == NO_MORE_ITEMS {
			Log(Warning, "No more items in Windows Registry")
			break
		}
		index++
		adapterID := string(utf16.Decode(adapter[0:length]))
		adapterID = removeZeroes(adapterID)
		path := fmt.Sprintf("%s\\%s\\Connection", NETWORK_KEY, adapterID)
		var iHandle syscall.Handle
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &iHandle)
		if err != nil {
			syscall.RegCloseKey(iHandle)
			continue
		}
		length = 1024
		aName := make([]byte, length)
		err = syscall.RegQueryValueEx(iHandle, syscall.StringToUTF16Ptr("Name"), nil, nil, &aName[0], &length)

		if err != nil {
			syscall.RegCloseKey(iHandle)
			continue
		}

		aNameUtf16 := make([]uint16, length/2)
		for i := 0; i < int(length)-2; i += 2 {
			aNameUtf16[i/2] = binary.LittleEndian.Uint16(aName[i:])
		}
		aNameUtf16[length/2-1] = 0

		adapterName := string(utf16.Decode(aNameUtf16))
		adapterName = removeZeroes(adapterName)
		Log(Warning, "AdapterName : %s, len : %d", adapterName, len(adapterName))

		var isInUse = false
		for _, i := range UsedInterfaces {
			if i == adapterName {
				isInUse = true
			}
		}
		if isInUse {
			Log(Warning, "Adapter already in use. Skipping.")
			continue
		}
		UsedInterfaces = append(UsedInterfaces, adapterName)

		tapname := fmt.Sprintf("%s%s%s", USERMODE_DEVICE_DIR, adapterID, TAP_SUFFIX)

		dev.file, err = syscall.CreateFile(syscall.StringToUTF16Ptr(tapname),
			syscall.GENERIC_WRITE|syscall.GENERIC_READ,
			0,
			nil,
			syscall.OPEN_EXISTING,
			syscall.FILE_ATTRIBUTE_SYSTEM|syscall.FILE_FLAG_OVERLAPPED,
			0)
		if err != nil {
			syscall.CloseHandle(dev.Handle)
			continue
		}
		Log(Info, "Acquired control over TAP interface: %s", adapterName)
		dev.Name = adapterID
		dev.Interface = adapterName
		return &dev, nil
	}
	return nil, nil
}

func createNewTAPDevice() {
	// Check if we already have devices
	/*
		if len(UsedInterfaces) == 0 {
			// If not, remove interfaces from previous instances and/or created by other software
			// Yes, this will active OpenVPN Connections
			Log(WARNING, "Removing TUN/TAP Devices created by other applications or previous instances")
			remdev := exec.Command(REMOVE_DEV)
			err := remdev.Run()
			if err != nil {
				Log(ERROR, "Failed to remove TUN/TAP Devices: %v", err)
			}
		}

		// Now add a new device
		Log(INFO, "Creating new TUN/TAP Device")
		adddev := exec.Command(ADD_DEV)
		err := adddev.Run()
		if err != nil {
			Log(ERROR, "Failed to add TUN/TAP Device: %v", err)
		}*/
}

func openDevice(ifPattern string) (*Interface, error) {
	createNewTAPDevice()
	handle, err := queryNetworkKey()
	if err != nil {
		Log(Error, "Failed to query Windows registry: %v", err)
		return nil, err
	}
	dev, err := queryAdapters(handle)
	if err != nil {
		Log(Error, "Failed to query network adapters: %v", err)
		return nil, err
	}
	if dev == nil {
		Log(Error, "All devices are in use")
		return nil, errors.New("No free TAP devices")
	}
	if dev.Name == "" {
		Log(Error, "Failed to query network adapters: %v", err)
		return nil, errors.New("Empty network adapter")
	}
	err = syscall.CloseHandle(handle)
	if err != nil {
		Log(Error, "Failed to close retrieved handle: %v", err)
	}

	return dev, nil
}

func createInterface(file syscall.Handle, ifPattern string, kind DevKind) (string, error) {
	return "1", nil
}

// ConfigureInterface configures TAP interface by assigning it's IP and netmask using netsh Command
func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	dev.IP = ip
	dev.Mask = "255.255.255.0"
	Log(Info, "Configuring %s. IP: %s Mask: %s", dev.Interface, dev.IP, dev.Mask)
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}
	cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, dev.Interface, dev.IP, dev.Mask)
	Log(Info, "Executing: %s", cmd)
	setip.SysProcAttr.CmdLine = cmd
	err := setip.Run()
	if err != nil {
		Log(Error, "Failed to properly configure TAP device with netsh: %v", err)
		return err
	}

	in := []byte("\x01\x00\x00\x00")
	var length uint32
	err = syscall.DeviceIoControl(dev.file, TAP_IOCTL_SET_MEDIA_STATUS,
		&in[0],
		uint32(len(in)),
		&in[0],
		uint32(len(in)),
		&length,
		nil)
	if err != nil {
		Log(Error, "Failed to change device status to 'connected': %v", err)
		return err
	}
	return nil
}

func ExtractMacFromInterface(dev *Interface) string {
	mac := make([]byte, 6)
	var length uint32
	err := syscall.DeviceIoControl(dev.file, TAP_IOCTL_GET_MAC, &mac[0], uint32(len(mac)), &mac[0], uint32(len(mac)), &length, nil)
	if err != nil {
		Log(Error, "Failed to get MAC from device")
	}
	var macAddr bytes.Buffer

	i := 0
	for _, a := range mac {
		if a == 0 {
			macAddr.WriteString("00")
		} else if a < 16 {
			macAddr.WriteString(fmt.Sprintf("0%x", a))
		} else {
			macAddr.WriteString(fmt.Sprintf("%x", a))
		}
		if i < 5 {
			macAddr.WriteString(":")
		}
		i++
	}
	Log(Info, "MAC: %s", macAddr.String())
	return macAddr.String()
}

func (t *Interface) Run() {
	t.Rx = make(chan []byte, 1500)
	t.Tx = make(chan []byte, 1500)
	go func() {
		if err := t.Read(t.Rx); err != nil {
			Log(Error, "Failed to read packet: %v", err)
		}
	}()
	go func() {
		if err := t.Write(t.Tx); err != nil {
			Log(Error, "Failed ro write packet: %v", err)
		}
	}()
}

func LinkUp(device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func SetIp(ip, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func SetMac(mac, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func (t *Interface) ReadPacket() (*Packet, error) {
	buf := <-t.Rx
	n := len(buf)
	if n <= 4 {
		return nil, nil
	}
	p := 12
	var pkt *Packet
	pkt = &Packet{Packet: buf[0:n]}
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[p : p+2]))
	flags := int(*(*uint16)(unsafe.Pointer(&buf[0])))
	if flags&flagTruncated != 0 {
		pkt.Truncated = true
	}
	pkt.Truncated = false
	return pkt, nil
}

func (t *Interface) WritePacket(pkt *Packet) error {
	n := len(pkt.Packet)
	buf := make([]byte, n)
	copy(buf, pkt.Packet)
	t.Tx <- buf[:n]
	return nil
}

func (t *Interface) Close() error {
	for i, iface := range UsedInterfaces {
		if iface == t.Interface {
			UsedInterfaces = append(UsedInterfaces[:i], UsedInterfaces[i+1:]...)
		}
	}
	syscall.Close(t.Handle)
	return nil
}

func CheckPermissions() bool {
	return true
}

func Open(ifPattern string, kind DevKind) (*Interface, error) {
	inf, err := openDevice(ifPattern)
	if err != nil {
		return nil, err
	}
	return inf, err
}

func (t *Interface) Read(ch chan []byte) (err error) {
	rx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	rx.HEvent = syscall.Handle(hevent)
	buf := make([]byte, 1500)
	var l uint32
	for {
		if err := syscall.ReadFile(t.file, buf, &l, &rx); err != nil {
		}
		if _, err := syscall.WaitForSingleObject(rx.HEvent, syscall.INFINITE); err != nil {
			Log(Error, "Failed to read from TUN/TAP: %v", err)
		}
		rx.Offset += l
		ch <- buf
	}
}

func (t *Interface) Write(ch chan []byte) (err error) {
	tx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	tx.HEvent = syscall.Handle(hevent)
	for {
		select {
		case data := <-ch:
			var l uint32
			syscall.WriteFile(t.file, data, &l, &tx)
			syscall.WaitForSingleObject(tx.HEvent, syscall.INFINITE)
			tx.Offset += uint32(len(data))
		}
	}
}

func GetDeviceBase() string {
	return "Local Area Network "
}

// Syslog provides additional logging to the syslog server
func Syslog(level LogLevel, format string, v ...interface{}) {
	Log(Info, "Syslog is not supported on this platform. Please do not use syslog flag.")
}
