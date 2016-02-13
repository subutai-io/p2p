package ptp

import (
	"os"
)

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/"+ifPattern, os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind, meta bool) (string, error) {
	return "1", nil
}
