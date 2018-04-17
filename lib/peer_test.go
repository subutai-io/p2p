package ptp

import (
	"net"
	"testing"
	"time"
)

func TestSetState(t *testing.T) {
	np := new(NetworkPeer)
	ptpc := new(PeerToPeer)
	states := [...]int{
		int(PeerStateInit),
		int(PeerStateRequestedIP),
		int(PeerStateRequestingProxy),
		int(PeerStateWaitingForProxy),
		int(PeerStateWaitingToConnect),
		int(PeerStateConnecting),
		int(PeerStateConnected),
		int(PeerStateDisconnect),
		int(PeerStateStop),
		int(PeerStateCooldown),
	}
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
	if err != nil && np.State != PeerStateRequestedIP {
		t.Error("Error in initializing peer")
	}
	np.KnownIPs = make([]*net.UDPAddr, 1)
	addr, _ := net.ResolveUDPAddr("udp4", "192.168.1.1:24")
	np.KnownIPs = append(np.KnownIPs, addr)
	err2 := np.stateInit(ptp)
	if err2 != nil && np.State != PeerStateRequestingProxy {
		t.Error(err2)
	}
	np.Proxies = make([]*net.UDPAddr, 1)
	np.Proxies = append(np.Proxies, addr)
	err3 := np.stateInit(ptp)
	if err3 != nil && np.State != PeerStateWaitingToConnect {
		t.Error(err3)
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

func TestStateConnected(t *testing.T) {
	np := new(NetworkPeer)
	ptp := new(PeerToPeer)
	np.RemoteState = PeerStateDisconnect
	err := np.stateDisconnect(ptp)
	if err != nil && np.State != PeerStateDisconnect {
		t.Error(err)
	}
	np.RemoteState = PeerStateStop
	err2 := np.stateConnected(ptp)
	if err2 != nil && np.State != PeerStateDisconnect {
		t.Error(err2)
	}
	np.RemoteState = PeerStateInit
	err3 := np.stateConnected(ptp)
	if err3 != nil && np.State != PeerStateInit {
		t.Error(err3)
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

func TestAddEndpoint(t *testing.T) {
	np := new(NetworkPeer)
	addr, _ := net.ResolveUDPAddr("udp4", "192.168.1.1:24")
	err := np.addEndpoint(addr)
	if err != nil {
		t.Error(err)
	}
	for _, ep := range np.Endpoints {
		if ep.Addr != addr {
			t.Error("Error.Can't add address")
		}
	}
	addr2, _ := net.ResolveUDPAddr("udp4", "192.168.1.2:24")
	np.Endpoints = make([]PeerEndpoint, 2)
	var peerEp1 PeerEndpoint
	peerEp1.Addr = addr2
	peerEp1.LastContact = time.Now()
	np.Endpoints = append(np.Endpoints, peerEp1)
	err2 := np.addEndpoint(addr2)
	if err2 == nil {
		t.Error(err2)
	}
}
