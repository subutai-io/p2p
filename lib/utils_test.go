package ptp

import (
	"net"
	"reflect"
	"sync"
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
	}
}

// func TestFindNetworkAddresses(t *testing.T) {
// 	ptp := new(PeerToPeer)
// 	ptp.FindNetworkAddresses()
// 	// fmt.Printf("%+v\n", ptp.LocalIPs)
// }

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{"t1", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateToken(); got == tt.want {
				t.Errorf("GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isPrivateIP(t *testing.T) {
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{"Empty IP", args{}, false, true},
		{"10.x subnet", args{net.ParseIP("10.12.13.14")}, true, false},
		{"10.x subnet", args{net.ParseIP("172.16.17.18")}, true, false},
		{"10.x subnet", args{net.ParseIP("192.168.0.1")}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isPrivateIP(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("isPrivateIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("isPrivateIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifyState(t *testing.T) {
	type args struct {
		state PeerState
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Stringify state: Init", args{PeerStateInit}, "INITIALIZING"},
		{"Stringify state: Waiting IP", args{PeerStateRequestedIP}, "WAITING_IP"},
		{"Stringify state: Requesting Proxies", args{PeerStateRequestingProxy}, "REQUESTING_PROXIES"},
		{"Stringify state: Waiting Proxies", args{PeerStateWaitingForProxy}, "WAITING_PROXIES"},
		{"Stringify state: Waiting Connection", args{PeerStateWaitingToConnect}, "WAITING_CONNECTION"},
		{"Stringify state: Initializing Connection", args{PeerStateConnecting}, "INITIALIZING_CONNECTION"},
		{"Stringify state: Connected", args{PeerStateConnected}, "CONNECTED"},
		{"Stringify state: Disconnected", args{PeerStateDisconnect}, "DISCONNECTED"},
		{"Stringify state: Stopped", args{PeerStateStop}, "STOPPED"},
		{"Stringify state: Cooldown", args{PeerStateCooldown}, "COOLDOWN"},
		{"Stringify state: Unknown", args{PeerState(99)}, "UNKNOWN"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringifyState(tt.args.state); got != tt.want {
				t.Errorf("StringifyState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsInterfaceLocal(t *testing.T) {
	type args struct {
		ip net.IP
	}
	ActiveInterfaces = []net.IP{net.ParseIP("10.10.10.1")}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"Interface in list", args{net.ParseIP("10.10.10.1")}, true},
		{"Interface not in list", args{net.ParseIP("192.168.0.1")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInterfaceLocal(tt.args.ip); got != tt.want {
				t.Errorf("IsInterfaceLocal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_FindNetworkAddresses(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.FindNetworkAddresses()
		})
	}
}

func Test_min(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"min 1", args{0, 0}, 0},
		{"min 2", args{0, 1}, 0},
		{"min 3", args{1, 0}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := min(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("min() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_ParseInterfaces(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		interfaces []net.Interface
	}

	// hw1, _ := net.ParseMAC("06:01:02:03:04:05")

	// inf1 := net.Interface{
	// 	Index:        0,
	// 	MTU:          0,
	// 	Name:         "",
	// 	HardwareAddr: hw1,
	// }

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []net.IP
	}{
	//{"t1", fields{}, args{[]net.Interface{inf1}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.ParseInterfaces(tt.args.interfaces); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.ParseInterfaces() = %v, want %v", got, tt.want)
			}
		})
	}
}
