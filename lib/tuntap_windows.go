// +build windows
package ptp

import (
	"fmt"
	"github.com/lxn/win"
	"os"
	"syscall"
)

const (
	NETWORK_CONNECTIONS_KEY string = "SYSTEM\\CurrentControlSet\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	ADAPTER_KEY             string = "SYSTEM\\CurrentControlSet\\Control\\Class\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
)

func InitTuntap() int {
	var root win.HKEY
	rootpath, _ := syscall.UTF16PtrFromString(NETWORK_CONNECTIONS_KEY)
	ret := win.RegOpenKeyEx(win.HKEY_LOCAL_MACHINE, rootpath, 0, win.KEY_READ, &root)
	if ret != 0 {
		return ret
	}
	var (
		name_length  uint32 = 72
		key_type     uint32
		lpDataLength uint32 = 72
		zero_unit    uint32 = 0
	)
	name := make([]uint16, 72)
	lpData := make([]byte, 72)

	ret = win.RegEnumValue(root, zero_unit, &name[0], &name_length, nil, &key_type, &lpData[0], &lpDataLength)
	fmt.Printf("Execution result is: %d", ret)

	fmt.Printf("lpDataLength: %d\n", lpDataLength)
	fmt.Printf("name: %d\n", name)
	fmt.Printf("lpData: %s\n", string(lpData))

	win.RegCloseKey(root)
	return 0
}

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind, meta bool) (string, error) {
	panic("TUN/TAP functionality is not supported on this platform")
}

func ConfigureInterface(ip, mac, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
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
