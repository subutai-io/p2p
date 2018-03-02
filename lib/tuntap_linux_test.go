package ptp

import "testing"

func TestGetDeviceBase(t *testing.T) {
	get := GetDeviceBase()
	if get != "vptp" {
		t.Error("Error. Return wrong value")
	}
}

func TestGetConfigurationTool(t *testing.T) {
	get := GetConfigurationTool()
	wait := "/sbin/ip"
	if get != wait {
		t.Error("Error", get)
	}
}

func TestNewTAP(t *testing.T) {
	get1, err := newTAP("tool", "", "01:02:03:04:05:06", "255.255.255.0", 1)
	if get1 != nil {
		t.Error(err)
	}
	get2, err2 := newTAP("tool", "192.168.1.1", "-", "255.255.255.0", 1)
	if get2 != nil {
		t.Error(err2)
	}
}
