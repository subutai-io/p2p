// +build !linux,!darwin

package ptp

import (
	"os"
)

var flagTruncated = 0

func createInterface(f *os.File, ifPattern string, kind DevKind) (string, error) {
	panic("Not implemented on this platform")
}
