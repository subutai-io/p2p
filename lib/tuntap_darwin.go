package ptp

import (
	"os"
	"os/exec"
)

const (
	ConfigDir string = "/usr/local/etc"
)

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/"+ifPattern, os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
	return "1", nil
}

func closeInterface(file *os.File) {
	Log(Info, "Closing network interface")
	if file != nil {
		err := file.Close()
		if err != nil {
			Log(Error, "Failed to close network interface: %s", err)
			return
		}
		Log(Info, "Interface closed")
		return
	}
	Log(Warning, "Skipping previously closed interface")
}

func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	// First we need to set MAC address, because ifconfig requires interface to go down
	// before changing it
	setmac := exec.Command(tool, device, "ether", mac)
	err := setmac.Run()
	if err != nil {
		Log(Error, "Failed to set MAC: %v", err)
	}
	linkup := exec.Command(tool, device, ip, "netmask", "255.255.255.0", "up")
	err = linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func LinkUp(device, tool string) error {
	linkup := exec.Command(tool, "link", "set", "dev", device, "up")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func SetIp(ip, device, tool string) error {
	Log(Info, "Setting %s IP on device %s", ip, device)
	setip := exec.Command(tool, "addr", "add", ip+"/24", "dev", device)
	err := setip.Run()
	if err != nil {
		Log(Error, "Failed to set IP: %v", err)
		return err
	}
	return err
}

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

func GetDeviceBase() string {
	return "tap"
}
