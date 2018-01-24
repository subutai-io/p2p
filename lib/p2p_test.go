package ptp

import (
	"testing"
)

func TestGenerateDeviceName(t *testing.T) {
	p := new(PeerToPeer)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}

func TestIsIPv4(t *testing.T) {
	ip1 := "194.152.36.143"
	ip2 := "2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"
	ptp := new(PeerToPeer)
	wait1 := true
	get1 := ptp.IsIPv4(ip1)
	if get1 != wait1 {
		t.Errorf("Error: wait %v, get %v", wait1, get1)
	}
	wait2 := false
	get2 := ptp.IsIPv4(ip2)
	if get2 != wait2 {
		t.Errorf("Error: wait %v, get %v", wait2, get2)
	}
}

func TestParseIntroString(t *testing.T) {
	p := new(PeerToPeer)
	id, mac, ip := p.ParseIntroString("id,01:02:03:04:05:06,127.0.0.1")
	if id != "id" || mac.String() != "01:02:03:04:05:06" || ip.String() != "127.0.0.1" {
		t.Errorf("Failed to parse intro string")
	}
}

func TestPrepareIntroductionMessage(t *testing.T) {
	p := new(PeerToPeer)
	p.Interface, _ = newTAP("", "127.0.0.1", "01:02:03:04:05:06", "", 1)
	msg := p.PrepareIntroductionMessage("test-id")
	if string(msg.Data) != "test-id,01:02:03:04:05:06,127.0.0.1" {
		t.Errorf("Failed to create introduction message")
	}
}
