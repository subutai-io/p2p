package ptp

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

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
		EndpointsHeap      []*Endpoint
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

	// all := r1

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
		name   string
		fields fields
		args   args
	}{
		{"nil ptp", fields{}, args{}},
		{"running peer", fields{Running: true}, args{ptp}},
		{"stopper peer", fields{State: PeerStateStop}, args{ptp}},
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
