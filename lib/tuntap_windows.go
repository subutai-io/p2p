// +build windows
package ptp

import (
	"os"
)

var flagTruncated = 0

func createInterface(f *os.File, ifPattern string, kind DevKind) (string, error) {
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
