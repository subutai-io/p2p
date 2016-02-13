package ptp

import (
	"os"
	"syscall"
	"unsafe"
)

func openDevice(ifPattern string) (*os.File, error) {
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind, meta bool) (string, error) {
	var req ifReq
	//req.Flags = iffOneQueue
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
	if !meta {
		req.Flags |= iffnopi
	}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != 0 {
		return "", err
	}
	return string(req.Name[:]), nil
}
