package ptp

import (
	"bytes"
	"net"
	"reflect"
	"testing"
)

func TestIsDeviceExists(t *testing.T) {
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

func TestFindNetworkAddresses(t *testing.T) {
	ptp := new(PeerToPeer)
	ptp.FindNetworkAddresses()
	if !true {
		t.Error("Error in function")
	}
}

func TestRetrieveFirstDHTRouters(t *testing.T) {
	ptp := new(PeerToPeer)
	wait, err := net.ResolveUDPAddr("udp4", "192.168.11.5:6882")
	if err != nil {
		t.Error("error")
	}
	ptp.Routers = ""
	get := ptp.retrieveFirstDHTRouter()
	if get != nil {
		t.Error("Length of ptp routers is nil")
	}
	ptp.Routers = "192.168.11.5:24,192.168.22.1:22"
	get2 := ptp.retrieveFirstDHTRouter()

	if bytes.EqualFold(get2.IP, wait.IP) && get2.Port != wait.Port && get2.Zone != wait.Zone {
		t.Errorf("Error.Wait %v, get %v", wait, get2)
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
	get1, err := ptp.validateInterfaceName("lo")
	if get1 != "lo" {
		t.Error(err)
	}
	get2, err2 := ptp.validateInterfaceName("")
	if get2 != "vptp1" {
		t.Error(err2)
	}
	get3, err3 := ptp.validateInterfaceName("123456789101112")
	if get3 != "" {
		t.Error(err3)
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

func TestMarkPeerForRemoval(t *testing.T) {
	ptp := new(PeerToPeer)
	np := new(NetworkPeer)
	ptp.Init()
	ptp.Peers.peers["1"] = np
	get := ptp.markPeerForRemoval("1", "Some reasons")
	if get != nil && np.State != PeerStateDisconnect {
		t.Error("Error")
	}
	get2 := ptp.markPeerForRemoval("0", "some reasons")
	if get2 == nil {
		t.Error("Error")
	}
}

func TestParseIntroString(t *testing.T) {
	// TODO: Fix this test
	// p := new(PeerToPeer)
	// id, mac, ip := p.ParseIntroString("id,01:02:03:04:05:06,127.0.0.1")
	// if id != "id" || mac.String() != "01:02:03:04:05:06" || ip.String() != "127.0.0.1" {
	// 	t.Errorf("Failed to parse intro string")
	// }
	// id2, mac2, ip2 := p.ParseIntroString("id,192.168.14.1")
	// if id2 != "" && mac2 != nil && ip2 != nil {
	// 	t.Error("Insufficient number of parameters")
	// }
	// id3, mac3, ip3 := p.ParseIntroString("id,mac,192.168.14.1")
	// if id3 != "" && mac3 != nil && ip3 != nil {
	// 	t.Error("Error in parse MAC")
	// }
	// id4, mac4, ip4 := p.ParseIntroString("id,01:02:03:04:05:06,ip")
	// if id4 != "" && mac4 != nil && ip4 != nil {
	// 	t.Error("Error in parse ip")
	// }
}
