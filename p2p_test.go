package main

import (
	"testing"
)

func TestGenerateDeviceName(t *testing.T) {
	p := new(PTPCloud)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}
