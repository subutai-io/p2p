package ptp

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// Constants
const (
	ConfigDir  string = "/usr/local/etc"
	DefaultMTU string = "1376"
)

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
	var req ifReq
	req.Flags = 0
	copy(req.Name[:15], ifPattern)
	switch kind {
	case DevTun:
		req.Flags |= iffTun
	case DevTap:
		req.Flags |= iffTap
	default:
		panic("Unknown interface type")
	}
	req.Flags |= iffnopi

	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != 0 {
		return "", err
	}
	return string(req.Name[:]), nil
}

// ConfigureInterface performs a configuration of an existing interface
func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	err := LinkUp(device, tool)
	if err != nil {
		return err
	}

	err = SetMTU(dev, device, tool, DefaultMTU)
	if err != nil {
		return err
	}

	// Configure new device
	err = SetIP(ip, device, tool)
	if err != nil {
		return err
	}

	err = SetMac(mac, device, tool)
	return err
}

// SetMTU sets an MTU value
func SetMTU(dev *Interface, device, tool, mtu string) error {
	setmtu := exec.Command(tool, "link", "set", "dev", device, "mtu", mtu)
	err := setmtu.Run()
	if err != nil {
		Log(Error, "Failed to set MTU on device %s: %v", device, err)
		return err
	}
	return nil
}

// LinkUp brings interface up
func LinkUp(device, tool string) error {
	linkup := exec.Command(tool, "link", "set", "dev", device, "up")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

// SetIP sets an IP address to an interface
func SetIP(ip, device, tool string) error {
	Log(Info, "Setting %s IP on device %s", ip, device)
	setip := exec.Command(tool, "addr", "add", ip+"/24", "dev", device)
	err := setip.Run()
	if err != nil {
		Log(Error, "Failed to set IP: %v", err)
		return err
	}
	return err
}

// SetMac sets a MAC address to a device
func SetMac(mac, device, tool string) error {
	// Set MAC to device
	Log(Info, "Setting %s MAC on device %s", mac, device)
	setmac := exec.Command(tool, "link", "set", "dev", device, "address", mac)
	err := setmac.Run()
	if err != nil {
		Log(Error, "Failed to set MAC: %v", err)
		return err
	}
	return err
}

// GetDeviceBase returns a default interface name
func GetDeviceBase() string {
	return "vptp"
}
