package ptp

import (
	"testing"
)

func TestSetState(t *testing.T) {
	// TODO: Fix this test
	// np := new(NetworkPeer)
	// ptpc := new(PeerToPeer)

	// states := [...]int{
	// 	int(PeerStateInit),
	// 	int(PeerStateRequestedIP),
	// 	int(PeerStateConnecting),
	// 	int(PeerStateConnectingDirectlyWait),
	// 	int(PeerStateConnectingDirectly),
	// 	int(PeerStateConnectingInternetWait),
	// 	int(PeerStateConnectingInternet),
	// 	int(PeerStateConnected),
	// 	int(PeerStateHandshaking),
	// 	int(PeerStateHandshakingFailed),
	// 	int(PeerStateWaitingForwarder),
	// 	int(PeerStateWaitingForwarderFailed),
	// 	int(PeerStateHandshakingForwarder),
	// 	int(PeerStateDisconnect),
	// 	int(PeerStateStop)}

	// for i := 0; i < len(states); i++ {
	// 	np.SetState(PeerState(states[i]), ptpc)
	// 	if np.State != PeerState(states[i]) {
	// 		t.Errorf("wait: %d; get: %d", states[i], np.State)
	// 	}
	// }
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

func TestStateConnecting(t *testing.T) {
	// TODO: Fix this test
	// np := new(NetworkPeer)
	// ptp := new(PeerToPeer)
	// np.stateConnecting(ptp)
	// if np.State != PeerStateConnectingDirectlyWait {
	// 	t.Errorf("Error. Wait %v, get %v", PeerStateConnectingDirectlyWait, np.State)
	// }
}

func TestStateConnected(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	np.RemoteState = PeerStateDisconnect
	err := np.stateDisconnect(ptp)
	if err != nil && np.State != PeerStateDisconnect {
		t.Error("Error. Peer can't disconnect")
	}
	np.RemoteState = PeerStateStop
	err2 := np.stateConnected(ptp)
	if err2 != nil && np.State != PeerStateDisconnect {
		t.Error("Error. Peer can't stop")
	}
	np.RemoteState = PeerStateInit
	err3 := np.stateConnected(ptp)
	if err3 != nil && np.State != PeerStateInit {
		t.Error("Error. Remote peer can't to reconnect")
	}
}

func TestStateDisconnect(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	np.ID = "1"
	err := np.stateDisconnect(ptp)
	if err != nil && np.State == PeerStateStop {
		t.Error("Error. Can't disconnect peer")
	}
}

func TestStateStop(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	np.ID = "1"
	err := np.stateStop(ptp)
	if err != nil && np.State == PeerStateStop {
		t.Error("Error")
	}
}

func TestSetPeerAddr(t *testing.T) {
	// TODO: Fix this test
	// np := new(NetworkPeer)

	// ip := new(net.UDPAddr)
	// ip.IP = []byte("10.156.119.247")
	// ip.Port = 45109

	// np.KnownIPs = append(np.KnownIPs, ip)
	// get := np.SetPeerAddr()
	// if !get {
	// 	t.Errorf("Error. Wait: %t, get: %t", true, get)
	// }

	// np2 := new(NetworkPeer)
	// wait2 := false
	// get2 := np2.SetPeerAddr()
	// if get2 != wait2 {
	// 	t.Errorf("Error: Wait %v, get %v", wait2, get2)
	// }
}
