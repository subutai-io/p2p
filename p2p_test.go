package main

import (
	"net"
	"testing"
)

func TestGenerateDeviceName(t *testing.T) {
	p := new(PTPCloud)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}

func TestParseIntroString(t *testing.T) {
	p := new(PTPCloud)
	id, mac, ip := p.ParseIntroString("id,01:02:03:04:05:06,127.0.0.1")
	if id != "id" || mac.String() != "01:02:03:04:05:06" || ip.String() != "127.0.0.1" {
		t.Errorf("Failed to parse intro string")
	}
}

func TestGenerateMac(t *testing.T) {
	macs := make(map[string]net.HardwareAddr)

	for i := 0; i < 10000; i++ {
		smac, mac := GenerateMAC()
		if smac == "" {
			t.Errorf("Failed to generate mac")
			return
		}
		_, e := macs[smac]
		if e {
			t.Errorf("Same MAC was generated")
			return
		}
		macs[smac] = mac
	}
}

func TestPrepareIntroductionMessage(t *testing.T) {
	p := new(PTPCloud)
	p.Mac = "01:02:03:04:05:06"
	p.IP = "127.0.0.1"
	msg := p.PrepareIntroductionMessage("test-id")
	if string(msg.Data) != "test-id,01:02:03:04:05:06,127.0.0.1" {
		t.Errorf("Failed to create introduction message")
	}
}
