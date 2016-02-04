package main

import (
	//"github.com/subutai-io/p2p/udpcs"
	"testing"
)

func TestGenerateDeviceName(t *testing.T) {
	ptp := new(PTPCloud)
	dev := ptp.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed")
	}
}
