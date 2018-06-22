package ptp

import (
	"net"
	"reflect"
	"runtime"
	"testing"
)

func TestIsDeviceExists(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	ptp := new(PeerToPeer)
	dev1 := "lo"
	get1 := ptp.IsDeviceExists(dev1)
	if !get1 {
		t.Error("Error. Device exists.")
	}
	dev2 := "device"
	get2 := ptp.IsDeviceExists(dev2)
	if get2 {
		t.Errorf("Error. There no network interface such %v", dev2)
	}
}

func TestGenerateDeviceName(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	p := new(PeerToPeer)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}

func TestIsIPv4(t *testing.T) {
	ip1 := "194.152.36.143"
	ip2 := "2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"
	ip3 := ""
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
	wait3 := false
	get3 := ptp.IsIPv4(ip3)
	if get3 != wait3 {
		t.Errorf("Error: wait %v, get %v", wait3, get3)
	}
}

func TestValidateMac(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestValidateInterfaceName(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestParseIntroString(t *testing.T) {
	ptp := new(PeerToPeer)
	hs := new(PeerHandshake)
	hs.Endpoint, _ = net.ResolveUDPAddr("udp4", "192.168.1.1:24")
	get1, err1 := ptp.ParseIntroString("id,ip,mac")
	if get1 != nil {
		t.Error(err1)
	}
	get2, err2 := ptp.ParseIntroString("1,-,127.0.0.1,192.168.1.1")
	if get2 != nil {
		t.Error(err2)
	}
	get3, err3 := ptp.ParseIntroString("1,01:02:03:04:05:06,-,192.168.1.1")
	if get3 != nil {
		t.Error(err3)
	}
	get4, err4 := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,-")
	if get4 != nil {
		t.Error(err4)
	}
	get5, _ := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,192.168.1.1:24")
	if !reflect.DeepEqual(get5.Endpoint.IP, hs.Endpoint.IP) && get5.Endpoint.Port != hs.Endpoint.Port && get5.Endpoint.Zone != hs.Endpoint.Zone {
		t.Error("Error")
	}
}
