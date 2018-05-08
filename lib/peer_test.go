package ptp

import (
	"testing"
	"net"
	"time"
	"sync"
	"reflect"
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
/*
Generated TestNetworkPeer_reportState
Generated TestNetworkPeer_SetState
Generated TestNetworkPeer_Run
Generated TestNetworkPeer_stateInit
Generated TestNetworkPeer_stateRequestedIP
Generated TestNetworkPeer_stateDisconnect
Generated TestNetworkPeer_stateStop
Generated TestNetworkPeer_RequestForwarder
Generated TestNetworkPeer_stateConnecting
Generated TestNetworkPeer_punchUDPHole
Generated TestNetworkPeer_isEndpointActive
Generated TestNetworkPeer_stateRequestingProxy
Generated TestNetworkPeer_stateWaitingForProxy
Generated TestNetworkPeer_stateWaitingToConnect
Generated TestNetworkPeer_sortEndpoints
Generated TestNetworkPeer_route
Generated TestNetworkPeer_stateConnected
Generated TestNetworkPeer_stateCooldown
Generated TestNetworkPeer_addEndpoint
Generated TestNetworkPeer_pingEndpoints
Generated TestNetworkPeer_syncWithRemoteState
Generated TestNetworkPeer_BumpEndpoint
package ptp

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)
*/

/*
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
*/

func TestNetworkPeer_reportState(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.reportState(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_SetState(t *testing.T) {
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
		state PeerState
		ptpc  *PeerToPeer
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.SetState(tt.args.state, tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_Run(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.Run(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_stateInit(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateInit(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateInit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateRequestedIP(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateRequestedIP(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateRequestedIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateDisconnect(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateDisconnect(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateDisconnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateStop(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateStop(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_RequestForwarder(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.RequestForwarder(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_stateConnecting(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateConnecting(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateConnecting() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_punchUDPHole(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.punchUDPHole(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_isEndpointActive(t *testing.T) {
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
		ep *net.UDPAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
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
			if got := np.isEndpointActive(tt.args.ep); got != tt.want {
				t.Errorf("NetworkPeer.isEndpointActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetworkPeer_stateRequestingProxy(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateRequestingProxy(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateRequestingProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateWaitingForProxy(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateWaitingForProxy(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateWaitingForProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateWaitingToConnect(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateWaitingToConnect(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateWaitingToConnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
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
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*PeerEndpoint
		want1  []*PeerEndpoint
		want2  []*PeerEndpoint
	}{
		// TODO: Add test cases.
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

func TestNetworkPeer_route(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.route(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.route() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateConnected(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateConnected(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateConnected() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_stateCooldown(t *testing.T) {
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.stateCooldown(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.stateCooldown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_addEndpoint(t *testing.T) {
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
		addr *net.UDPAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
			if err := np.addEndpoint(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.addEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_pingEndpoints(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.pingEndpoints(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_syncWithRemoteState(t *testing.T) {
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.syncWithRemoteState(tt.args.ptpc)
		})
	}
}

func TestNetworkPeer_BumpEndpoint(t *testing.T) {
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
		epAddr string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
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
			np.BumpEndpoint(tt.args.epAddr)
		})
	}
}
