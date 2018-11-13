package ptp

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
		{"nil ptp", fields{}, args{}, true},
		{"non-nil ptp", fields{}, args{new(PeerToPeer)}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.reportState(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.reportState() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		state PeerState
		ptpc  *PeerToPeer
	}

	ptp0 := new(PeerToPeer)
	ptp0.Dht = new(DHTClient)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"nil dht", fields{}, args{ptpc: new(PeerToPeer)}, true},
		{"new state", fields{State: PeerStateConnected}, args{ptpc: ptp0, state: PeerStateDisconnect}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.SetState(tt.args.state, tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.SetState() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"running peer", fields{Running: true}, args{ptp}, true},
		{"stopper peer", fields{State: PeerStateStop}, args{ptp}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.Run(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp := new(PeerToPeer)

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	kip0 := []*net.UDPAddr{
		udp0,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"no ips", fields{}, args{ptp}, false},
		{"no proxies", fields{KnownIPs: kip0}, args{ptp}, false},
		{"passing", fields{KnownIPs: kip0, Proxies: kip0}, args{ptp}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp0 := new(PeerToPeer)
	ptp1 := new(PeerToPeer)
	ptp1.Dht = new(DHTClient)
	ptp2 := new(PeerToPeer)
	ptp2.Dht = new(DHTClient)
	ptp2.Dht.Init("myhash")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"nil dht", fields{}, args{ptp0}, true},
		{"node request failed", fields{}, args{ptp1}, true},
		{"5 attemps", fields{}, args{ptp2}, true},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.RequestForwarder(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.RequestForwarder() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.punchUDPHole(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.punchUDPHole() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"state connect", fields{RemoteState: PeerStateWaitingToConnect}, args{ptp}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	la1, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")
	la2, _ := net.ResolveUDPAddr("udp4", "10.0.0.1:1234")
	la3, _ := net.ResolveUDPAddr("udp4", "172.16.0.1:1234")

	ra1, _ := net.ResolveUDPAddr("udp4", "1.1.1.1:2345")
	ra2, _ := net.ResolveUDPAddr("udp4", "2.2.2.2:2345")

	ep1 := &Endpoint{
		Addr:        la1,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep2 := &Endpoint{
		Addr:        la2,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep3 := &Endpoint{
		Addr:        la3,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep4 := &Endpoint{
		Addr:        ra1,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep5 := &Endpoint{
		Addr:        ra2,
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	ep6 := &Endpoint{
		LastContact: time.Now(),
		LastPing:    time.Now(),
	}

	r1 := []*Endpoint{ep1, ep2, ep3}

	r2 := []*Endpoint{
		&Endpoint{
			Addr:        la1,
			LastContact: time.Unix(1, 1),
			LastPing:    time.Now(),
		},
		ep2, ep3,
	}

	r2_2 := []*Endpoint{
		ep2, ep3,
	}

	r3 := []*Endpoint{
		ep4, ep5,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []*Endpoint
		want1  []*Endpoint
		want2  []*Endpoint
	}{
		{"t1", fields{}, args{}, []*Endpoint{}, []*Endpoint{}, []*Endpoint{}},
		{"t2", fields{EndpointsHeap: r1}, args{}, r1, []*Endpoint{}, []*Endpoint{}},
		{"t3", fields{EndpointsHeap: r2}, args{}, r2_2, []*Endpoint{}, []*Endpoint{}},
		{"t4", fields{EndpointsHeap: r3}, args{}, []*Endpoint{}, r3, []*Endpoint{}},
		{"t5", fields{EndpointsHeap: r1, Proxies: []*net.UDPAddr{la1, la2, la3}}, args{}, []*Endpoint{}, []*Endpoint{}, r1},
		{"t6", fields{EndpointsHeap: []*Endpoint{ep6}}, args{}, []*Endpoint{}, []*Endpoint{}, []*Endpoint{}},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp0 := new(PeerToPeer)

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	heap0 := []*Endpoint{
		&Endpoint{
			Addr:        udp0,
			LastContact: time.Now(),
		},
	}

	heap1 := []*Endpoint{
		&Endpoint{
			Addr:        nil,
			LastContact: time.Now(),
		},
		&Endpoint{
			Addr:        udp0,
			LastContact: time.Now(),
		},
	}

	heap2 := []*Endpoint{
		&Endpoint{
			Addr:        nil,
			LastContact: time.Now(),
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"empty ep heap", fields{}, args{ptp0}, false},
		{"routing is not required", fields{RoutingRequired: false, EndpointsHeap: heap0}, args{ptp0}, false},
		{"routing is required", fields{RoutingRequired: true, EndpointsHeap: heap2}, args{ptp0}, false},
		{"routing is required>sample peer", fields{RoutingRequired: true, EndpointsHeap: heap0}, args{ptp0}, false},
		{"routing is not required>non nil ep", fields{Endpoint: udp0, EndpointsHeap: heap1, Proxies: []*net.UDPAddr{nil, udp0}}, args{ptp0}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
		{"nil ptp", fields{}, args{}, true},
		{"wrong state", fields{State: PeerStateConnecting}, args{new(PeerToPeer)}, true},
		{"new hole punch", fields{State: PeerStateConnected, LastPunch: time.Unix(0, 0)}, args{new(PeerToPeer)}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
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
		{"nil ptp", fields{}, args{}, true},
		{"not nil ptp", fields{}, args{new(PeerToPeer)}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		addr *net.UDPAddr
	}

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	heap0 := []*Endpoint{
		&Endpoint{
			Addr: nil,
		},
	}

	heap1 := []*Endpoint{
		&Endpoint{
			Addr: udp0,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil addr", fields{}, args{}, true},
		{"empty heap", fields{}, args{udp0}, false},
		{"nil ep addr", fields{EndpointsHeap: heap0}, args{udp0}, false},
		{"existing ep", fields{EndpointsHeap: heap1}, args{udp0}, true},
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
				RoutingRequired:    tt.fields.RoutingRequired,
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	ptp0 := new(PeerToPeer)
	ptp1 := new(PeerToPeer)
	ptp1.Dht = new(DHTClient)
	ptp2 := new(PeerToPeer)
	ptp2.Dht = new(DHTClient)
	ptp2.UDPSocket = new(Network)

	heap0 := []*Endpoint{
		&Endpoint{
			LastPing: time.Now(),
		},
	}

	heap1 := []*Endpoint{
		&Endpoint{
			LastPing: time.Unix(0, 0),
			Addr:     nil,
		},
	}

	heap2 := []*Endpoint{
		&Endpoint{
			LastPing: time.Unix(0, 0),
			Addr:     udp0,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{

		{"nil ptp", fields{}, args{}, true},
		{"nil dht", fields{}, args{ptp0}, true},
		{"nil socket", fields{}, args{ptp1}, true},
		{"interval not passed", fields{EndpointsHeap: heap0}, args{ptp2}, false},
		{"interval passed>nil addr", fields{EndpointsHeap: heap1}, args{ptp2}, false},
		{"interval passed>packet sent", fields{EndpointsHeap: heap2}, args{ptp2}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.pingEndpoints(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.pingEndpoints() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		ptpc *PeerToPeer
	}

	ptp := new(PeerToPeer)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil ptp", fields{}, args{}, true},
		{"disconnect", fields{RemoteState: PeerStateDisconnect}, args{ptp}, false},
		{"stop", fields{RemoteState: PeerStateStop}, args{ptp}, false},
		{"init", fields{RemoteState: PeerStateInit}, args{ptp}, false},
		{"waiting to connect", fields{RemoteState: PeerStateWaitingToConnect}, args{ptp}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.syncWithRemoteState(tt.args.ptpc); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.syncWithRemoteState() error = %v, wantErr %v", err, tt.wantErr)
			}
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	type args struct {
		epAddr string
	}

	ua0 := "192.168.0.1:1234"
	udp0, _ := net.ResolveUDPAddr("udp4", ua0)

	heap0 := []*Endpoint{
		&Endpoint{
			Addr: udp0,
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty", fields{}, args{}, true},
		{"ep found", fields{EndpointsHeap: heap0}, args{ua0}, false},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if err := np.BumpEndpoint(tt.args.epAddr); (err != nil) != tt.wantErr {
				t.Errorf("NetworkPeer.BumpEndpoint() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetworkPeer_IsRunning(t *testing.T) {
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
		EndpointsHeap      []*Endpoint
		Lock               sync.RWMutex
		punchingInProgress bool
		LastFind           time.Time
		LastPunch          time.Time
		Stat               PeerStats
		RoutingRequired    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"empty", fields{}, false},
		{"negative", fields{Running: false}, false},
		{"postitive", fields{Running: true}, true},
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
				RoutingRequired:    tt.fields.RoutingRequired,
			}
			if got := np.IsRunning(); got != tt.want {
				t.Errorf("NetworkPeer.IsRunning() = %v, want %v", got, tt.want)
			}
		})
	}
}
