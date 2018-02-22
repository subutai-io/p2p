package ptp

import (
	"net"
	"testing"
)

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
		Log(Info, "Mac generated: %v", mac)
	}
}

func TestGenerateToken(t *testing.T) {
	token := GenerateToken()
	if token == "" {
		t.Errorf("Failed to generate token")
		return
	}
	if len(token) > 0 {
		Log(Info, "Token generated: %s", token)
	}
}

func TestIsPrivateIP(t *testing.T) {
	var ip1 net.IP
	get1, err1 := isPrivateIP(ip1)
	if get1 != false {
		t.Error(err1)
	}
	ip2 := net.IP([]byte{10, 0, 1, 1})
	get2, _ := isPrivateIP(ip2)
	if get2 != true {
		t.Error("IP is private")
	}
	ip3 := net.IP([]byte{172, 16, 1, 1})
	get3, _ := isPrivateIP(ip3)
	if get3 != true {
		t.Error("IP is private")
	}
	ip4 := net.IP([]byte{192, 168, 1, 1})
	get4, _ := isPrivateIP(ip4)
	if get4 != true {
		t.Error("IP is private")
	}
}

func TestStringifyState(t *testing.T) {
	states := []PeerState{
		PeerStateInit,
		PeerStateRequestedIP,
		PeerStateRequestingProxy,
		PeerStateWaitingForProxy,
		PeerStateWaitingToConnect,
		PeerStateConnecting,
		PeerStateConnected,
		PeerStateDisconnect,
		PeerStateStop,
		PeerStateCooldown,
	}

	ret := []string{
		"Unknown",
		"Initializing",
		"Waiting for IP",
		"Requesting proxies",
		"Waiting for proxies",
		"Waiting for connection",
		"Initializing connection",
		"Connected",
		"Disconnected",
		"Stopped",
		"Cooldown",
	}

	for i := 0; i < len(states)+1; i++ {
		get := StringifyState(PeerState(i))
		if get != ret[i] {
			t.Errorf("Error. Get: %v, wait %v", get, ret[i])
		}
	}
}
