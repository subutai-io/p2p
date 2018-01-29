package ptp

import (
	"net"
	"testing"
)

func TestSetState(t *testing.T) {
	np := new(NetworkPeer)
	ptpc := new(PeerToPeer)

	states := [...]int{
		int(PeerStateInit),
		int(PeerStateRequestedIP),
		int(PeerStateConnecting),
		int(PeerStateConnectingDirectlyWait),
		int(PeerStateConnectingDirectly),
		int(PeerStateConnectingInternetWait),
		int(PeerStateConnectingInternet),
		int(PeerStateConnected),
		int(PeerStateHandshaking),
		int(PeerStateHandshakingFailed),
		int(PeerStateWaitingForwarder),
		int(PeerStateWaitingForwarderFailed),
		int(PeerStateHandshakingForwarder),
		int(PeerStateDisconnect),
		int(PeerStateStop)}

	for i := 0; i < len(states); i++ {
		np.SetState(PeerState(states[i]), ptpc)
		if np.State != PeerState(states[i]) {
			t.Errorf("wait: %d; get: %d", states[i], np.State)
		}
	}
}

func TestRun(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	if np.Running == false {
		if !true {
			t.Error("Error in Run. np.Running is False")
		}
	}
	np.Running = true
	np.State = PeerStateStop
	np.Run(ptp)
	if !true {
		t.Error("Error. Can't stop peer")
	}
}

func TestStateInit(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	err := np.stateInit(ptp)
	if err != nil {
		t.Error("Error in initializing peer")
	}
}

func TestSetPeerAddr(t *testing.T) {
	np := new(NetworkPeer)

	ip := new(net.UDPAddr)
	ip.IP = []byte("10.156.119.247")
	ip.Port = 45109

	np.KnownIPs = append(np.KnownIPs, ip)
	get := np.SetPeerAddr()
	if !get {
		t.Errorf("Error. Wait: %t, get: %t", true, get)
	}

	np2 := new(NetworkPeer)
	wait2 := false
	get2 := np2.SetPeerAddr()
	if get2 != wait2 {
		t.Errorf("Error: Wait %v, get %v", wait2, get2)
	}
}
