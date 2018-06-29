package ptp

import (
	"net"
	"reflect"
	"sync"
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
	if err != nil {
		t.Error("Error in initializing peer")
	}
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
	np.RoutingRequired = true
	np.stateConnected(ptp)
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

func TestNetworkPeer_sortEndpoints(t *testing.T) {
	type fields struct {
		ID                 string
		Endpoint           *net.UDPAddr
		KnownIPs           []*net.UDPAddr
		Proxies            []*net.UDPAddr
		PeerLocalIP        net.IP
		PeerHW             net.HardwareAddr
		State              PeerState
		RemoteState        PeerState
		LastContact        time.Time
		PingCount          uint8
		LastError          string
		ConnectionAttempts uint8
		handlers           map[PeerState]StateHandlerCallback
		Running            bool
		EndpointsHeap      []*PeerEndpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
	}
	type args struct {
		ptpc *PeerToPeer
	}

	la1, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")
	la2, _ := net.ResolveUDPAddr("udp4", "10.0.0.1:1234")
	la3, _ := net.ResolveUDPAddr("udp4", "172.16.0.1:1234")

	ra1, _ := net.ResolveUDPAddr("udp4", "1.1.1.1:2345")
	ra2, _ := net.ResolveUDPAddr("udp4", "2.2.2.2:2345")

	ep1 := &PeerEndpoint{
		Addr:        la1,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep2 := &PeerEndpoint{
		Addr:        la2,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep3 := &PeerEndpoint{
		Addr:        la3,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep4 := &PeerEndpoint{
		Addr:        ra1,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep5 := &PeerEndpoint{
		Addr:        ra2,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	r1 := []*PeerEndpoint{ep1, ep2, ep3}

	r2 := []*PeerEndpoint{
		&PeerEndpoint{
			Addr:        la1,
			LastContact: time.Unix(1, 1),
			LastPing:    time.Now(),
		},
		ep2, ep3,
	}

	r2_2 := []*PeerEndpoint{
		ep2, ep3,
	}

	r3 := []*PeerEndpoint{
		ep4, ep5,
	}

	// all := r1

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*PeerEndpoint
		want1  []*PeerEndpoint
		want2  []*PeerEndpoint
	}{
		{"t1", fields{}, args{}, []*PeerEndpoint{}, []*PeerEndpoint{}, []*PeerEndpoint{}},
		{"t2", fields{EndpointsHeap: r1}, args{}, r1, []*PeerEndpoint{}, []*PeerEndpoint{}},
		{"t3", fields{EndpointsHeap: r2}, args{}, r2_2, []*PeerEndpoint{}, []*PeerEndpoint{}},
		{"t4", fields{EndpointsHeap: r3}, args{}, []*PeerEndpoint{}, r3, []*PeerEndpoint{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			np := &NetworkPeer{
				ID:                 tt.fields.ID,
				Endpoint:           tt.fields.Endpoint,
				KnownIPs:           tt.fields.KnownIPs,
				Proxies:            tt.fields.Proxies,
				PeerLocalIP:        tt.fields.PeerLocalIP,
				PeerHW:             tt.fields.PeerHW,
				State:              tt.fields.State,
				RemoteState:        tt.fields.RemoteState,
				LastContact:        tt.fields.LastContact,
				PingCount:          tt.fields.PingCount,
				LastError:          tt.fields.LastError,
				ConnectionAttempts: tt.fields.ConnectionAttempts,
				handlers:           tt.fields.handlers,
				Running:            tt.fields.Running,
				EndpointsHeap:      tt.fields.EndpointsHeap,
				Lock:               tt.fields.Lock,
				punchingInProgress: tt.fields.punchingInProgress,
				LastFind:           tt.fields.LastFind,
				LastPunch:          tt.fields.LastPunch,
				Stat:               tt.fields.Stat,
			}
			got, got1, got2 := np.sortEndpoints(tt.args.ptpc)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NetworkPeer.sortEndpoints() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NetworkPeer.sortEndpoints() got1 = %v, want %v", got1, tt.want1)
			}
			if !reflect.DeepEqual(got2, tt.want2) {
				t.Errorf("NetworkPeer.sortEndpoints() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
