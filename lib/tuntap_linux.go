package ptp

import (
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

const (
	CONFIG_DIR  string = "/usr/local/etc"
	DEFAULT_MTU string = "1600"
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

func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	err := LinkUp(device, tool)
	if err != nil {
		return err
	}

	err = SetMTU(dev, device, tool, DEFAULT_MTU)
	if err != nil {
		return err
	}

	// Configure new device
	err = SetIp(ip, device, tool)
	if err != nil {
		return err
	}

	err = SetMac(mac, device, tool)
	if err != nil {
		return err
	}
	return nil
}

func SetMTU(dev *Interface, device, tool, mtu string) error {
	setmtu := exec.Command(tool, "link", "set", "dev", device, "mtu", mtu)
	err := setmtu.Run()
	if err != nil {
		Log(ERROR, "Failed to set MTU on device %s: %v", device, err)
		return err
	}
	return nil
}

func LinkUp(device, tool string) error {
	linkup := exec.Command(tool, "link", "set", "dev", device, "up")
	err := linkup.Run()
	if err != nil {
		Log(ERROR, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func SetIp(ip, device, tool string) error {
	Log(INFO, "Setting %s IP on device %s", ip, device)
	setip := exec.Command(tool, "addr", "add", ip+"/24", "dev", device)
	err := setip.Run()
	if err != nil {
		Log(ERROR, "Failed to set IP: %v", err)
		return err
	}
	return err
}

func SetMac(mac, device, tool string) error {
	// Set MAC to device
	Log(INFO, "Setting %s MAC on device %s", mac, device)
	setmac := exec.Command(tool, "link", "set", "dev", device, "address", mac)
	err := setmac.Run()
	if err != nil {
		Log(ERROR, "Failed to set MAC: %v", err)
		return err
	}
	return err
}

func GetDeviceBase() string {
	return "vptp"
}
