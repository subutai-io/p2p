package ptp

import (
	"os"
	"os/exec"
)

const (
	CONFIG_DIR string = "/usr/local/etc"
)

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/"+ifPattern, os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
	return "1", nil
}

func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	// First we need to set MAC address, because ifconfig requires interface to go down
	// before changing it
	setmac := exe.Command(tool, device, "ether", mac)
	if err != nil {
		Log(ERROR, "Failed to set MAC: %v", err)
	}
	linkup := exec.Command(tool, device, ip, "netmask", "255.255.255.0", "up")
	err := linkup.Run()
	if err != nil {
		Log(ERROR, "Failed to up link: %v", err)
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
	return "tap"
}
