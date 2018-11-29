package ptp

import (
	"net"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestPeerToPeer_AssignInterface(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		interfaceName string
	}

	inf1, _ := newTAP("", "192.168.0.1", "00:00:00:00:00:01", "", 1500, false)
	inf1.MarkConfigured()

	inf2, _ := newTAP("", "192.168.0.2", "00:00:00:00:00:02", "", 1500, false)

	inf3 := newEmptyTAP()

	inf4 := newEmptyTAP()
	inf4.IP = net.ParseIP("192.168.0.4")

	f1 := fields{
		Interface: inf1,
	}

	f2 := fields{
		Interface: inf2,
	}

	f3 := fields{
		Interface: inf3,
	}

	f4 := fields{
		Interface: inf4,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil interface", fields{}, args{}, true},
		{"preconfigured interface", f1, args{}, false},
		{"failed init", f2, args{""}, true},
		{"no ip provided", f3, args{"name"}, true},
		{"no mac provided", f4, args{"name"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.AssignInterface(tt.args.interfaceName); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.AssignInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_ListenInterface(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.ListenInterface(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.ListenInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_GenerateDeviceName(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if got := p.GenerateDeviceName(tt.args.i); got != tt.want {
				t.Errorf("PeerToPeer.GenerateDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_IsIPv4(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		ip string
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
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if got := p.IsIPv4(tt.args.ip); got != tt.want {
				t.Errorf("PeerToPeer.IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		mac        string
		hash       string
		keyfile    string
		key        string
		ttl        string
		target     string
		fwd        bool
		port       int
		outboundIP net.IP
	}
	tests := []struct {
		name string
		args args
		want *PeerToPeer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.mac, tt.args.hash, tt.args.keyfile, tt.args.key, tt.args.ttl, tt.args.target, tt.args.fwd, tt.args.port, tt.args.outboundIP); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_ReadDHT(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.ReadDHT(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.ReadDHT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_waitForRemotePort(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	socket0 := new(Network)
	socket1 := new(Network)
	socket1.remotePort = 1

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil socket", fields{}, true},
		{"no port timeout", fields{UDPSocket: socket0}, true},
		{"has port", fields{UDPSocket: socket1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.waitForRemotePort(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.waitForRemotePort() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_PrepareInterfaces(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		ip            string
		interfaceName string
	}

	inf0, _ := newTAP("ip", "10.10.10.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	n0 := "thisisareallylongnetworkinterfacename"

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil interfeace", fields{}, args{}, true},
		{"broken interface name", fields{Interface: inf0}, args{interfaceName: n0}, true},
		{"dhcp>request ip failed", fields{Interface: inf0}, args{"dhcp", ""}, true},
		{"static>broken", fields{Interface: inf0}, args{"badip", ""}, true},
		{"static>report ip failed", fields{Interface: inf0}, args{"10.10.10.1", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.PrepareInterfaces(tt.args.ip, tt.args.interfaceName); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.PrepareInterfaces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_attemptPortForward(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		port uint16
		name string
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
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.attemptPortForward(tt.args.port, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.attemptPortForward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_Init(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.Init(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_validateMac(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		mac string
	}

	hw0, _ := net.ParseMAC("00:11:22:33:44:55")

	tests := []struct {
		name   string
		fields fields
		args   args
		want   net.HardwareAddr
	}{
		{"empty mac", fields{}, args{"00:11:22:33:44:55"}, hw0},
		{"empty mac>broken", fields{}, args{"xx:cc:xx:cc:xx:cc"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if got := p.validateMac(tt.args.mac); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.validateMac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_validateInterfaceName(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			got, err := p.validateInterfaceName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.validateInterfaceName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeerToPeer.validateInterfaceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_setupHandlers(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.setupHandlers(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.setupHandlers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_RequestIP(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		mac    string
		device string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    net.IP
		want1   net.IPMask
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			got, got1, err := p.RequestIP(tt.args.mac, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.RequestIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.RequestIP() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PeerToPeer.RequestIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPeerToPeer_ReportIP(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		ipAddress string
		mac       string
		device    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    net.IP
		want1   net.IPMask
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			got, got1, err := p.ReportIP(tt.args.ipAddress, tt.args.mac, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.ReportIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.ReportIP() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PeerToPeer.ReportIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPeerToPeer_Run(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	dht0 := new(DHTClient)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil dht", fields{}, true},
		{"nil interface", fields{Dht: dht0}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.Run(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_checkLastDHTUpdate(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	dht0 := new(DHTClient)
	dht0.LastUpdate = time.Now()

	dht1 := new(DHTClient)
	dht1.LastUpdate = time.Unix(0, 0)

	pm0 := new(ProxyManager)
	pm0.init()

	pm1 := new(ProxyManager)
	pm1.init()
	pm1.proxies["p0"] = &proxyServer{}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil dht", fields{}, true},
		{"nil proxy manager", fields{Dht: dht0}, true},
		{"dht update timeout not passed", fields{Dht: dht0, ProxyManager: pm0}, false},
		{"dht update passed>proxies", fields{Dht: dht1, ProxyManager: pm0}, true},
		{"passing", fields{Dht: dht1, ProxyManager: pm1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.checkLastDHTUpdate(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.checkLastDHTUpdate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_removeStoppedPeers(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	peer0 := &NetworkPeer{
		State: PeerStateStop,
	}

	pl0 := new(Swarm)

	pl1 := new(Swarm)
	pl1.Init()
	pl1.peers["p1"] = peer0

	pl2 := new(Swarm)
	pl2.Init()
	pl2.peers["p1"] = peer0
	pl2.peers["p2"] = peer0

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil peers", fields{}, true},
		{"no peers", fields{Peers: pl0}, false},
		{"single peer", fields{Peers: pl1}, false},
		{"two peers", fields{Peers: pl2}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.removeStoppedPeers(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.removeStoppedPeers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_checkProxies(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.1.2:3456")

	dht0 := new(DHTClient)

	pm0 := new(ProxyManager)
	pm1 := new(ProxyManager)
	pm1.init()
	pm1.proxies["192.168.0.1:1234"] = &proxyServer{
		Status:     proxyActive,
		Endpoint:   udp0,
		LastUpdate: time.Now(),
	}
	pm1.hasChanges = true

	socket0 := new(Network)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil dht", fields{}, true},
		{"nil proxy manager", fields{Dht: dht0}, true},
		{"nil udp socket", fields{Dht: dht0, ProxyManager: pm0}, true},
		{"check active", fields{Dht: dht0, ProxyManager: pm1, UDPSocket: socket0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.checkProxies(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.checkProxies() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_checkPeers(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	udp0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	socket0 := new(Network)

	dht0 := new(DHTClient)
	dht1 := &DHTClient{
		ID: "90805338-69d0-4bf1-9817-e8f74fb3ebe8",
	}

	pl0 := new(Swarm)
	pl1 := new(Swarm)
	pl1.Init()
	pl1.peers["peer0"] = &NetworkPeer{
		EndpointsHeap: []*Endpoint{
			nil,
			&Endpoint{Addr: udp0},
		},
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil dht", fields{}, true},
		{"nil peer list", fields{Dht: dht0}, true},
		{"nil socket", fields{Dht: dht0, Peers: pl0}, true},
		{"small id", fields{Dht: dht0, Peers: pl0, UDPSocket: socket0}, true},
		{"nil endpoint", fields{Dht: dht1, Peers: pl1, UDPSocket: socket0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.checkPeers(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.checkPeers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_PrepareIntroductionMessage(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		id       string
		endpoint string
	}

	inf0, _ := newTAP("ip", "10.11.12.13", "00:11:22:33:44:55", "255.255.255.0", 1500, false)

	msg0, _ := CreateMessageStatic(MsgTypeIntro, []byte("id0,00:11:22:33:44:55,10.11.12.13,10.10.10.1:1234"))

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *P2PMessage
		wantErr bool
	}{
		{"nil interface", fields{}, args{}, nil, true},
		{"normal message", fields{Interface: inf0}, args{"id0", "10.10.10.1:1234"}, msg0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			got, err := p.PrepareIntroductionMessage(tt.args.id, tt.args.endpoint)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.PrepareIntroductionMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.PrepareIntroductionMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_WriteToDevice(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		b         []byte
		proto     uint16
		truncated bool
	}

	inf0, _ := newTAP("ip", "192.168.0.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	inf1, _ := newTAP("ip", "192.168.0.2", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	inf1.file, _ = os.OpenFile("/tmp/fake-interface-0", os.O_CREATE|os.O_RDWR, 0700)
	defer inf1.file.Close()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, true},
		{"packet0", fields{Interface: inf0}, args{[]byte{0x01}, 0, false}, true},
		{"real interface", fields{Interface: inf1}, args{[]byte{0x01}, 0, false}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.WriteToDevice(tt.args.b, tt.args.proto, tt.args.truncated); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.WriteToDevice() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_SendTo(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		dst net.HardwareAddr
		msg *P2PMessage
	}

	dst0 := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}

	udp0, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")

	pl0 := new(Swarm)
	pl1 := new(Swarm)
	pl1.Init()
	pl1.peers["p0"] = &NetworkPeer{
		ID:       "p0",
		Endpoint: udp0,
	}
	pl1.tableMacID["01:02:03:04:05:06"] = "p0"

	socket0 := new(Network)
	msg0 := new(P2PMessage)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{"nil peer list", fields{}, args{}, -1, true},
		{"nil socket", fields{Peers: pl0}, args{}, -1, true},
		{"nil msg", fields{Peers: pl0, UDPSocket: socket0}, args{nil, nil}, -1, true},
		{"nil dst", fields{Peers: pl0, UDPSocket: socket0}, args{nil, msg0}, -1, true},
		{"non existing endpoint", fields{Peers: pl0, UDPSocket: socket0}, args{dst0, msg0}, 0, false},
		{"existing endpoint", fields{Peers: pl1, UDPSocket: socket0}, args{dst0, msg0}, -1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			got, err := p.SendTo(tt.args.dst, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.SendTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeerToPeer.SendTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_Close(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, false},
		{"with dht", fields{Dht: new(DHTClient)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{

				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.Close(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_deactivateInterface(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	inf0, _ := newTAP("ip", "192.168.0.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	inf1, _ := newTAP("ip", "192.168.0.2", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	ActiveInterfaces = append(ActiveInterfaces, net.ParseIP("192.168.0.2"))
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, true},
		{"inactive interface", fields{Interface: inf0}, true},
		{"active interface", fields{Interface: inf1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.deactivateInterface(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.deactivateInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_stopInterface(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	inf0, _ := newTAP("ip", "192.168.0.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	inf1, _ := newTAP("ip", "192.168.0.2", "00:11:22:33:44:55", "255.255.255.0", 1500, false)
	inf1.file, _ = os.OpenFile("/tmp/test-p2p-close", os.O_CREATE|os.O_RDWR, 0700)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, true},
		{"existing interface", fields{Interface: inf0}, true},
		{"working interface", fields{Interface: inf1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.stopInterface(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.stopInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_stopPeers(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	pl0 := new(Swarm)
	pl0.Init()
	pl0.peers["peer1"] = &NetworkPeer{
		ID: "peer1",
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, true},
		{"single peer", fields{Peers: pl0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.stopPeers(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.stopPeers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_stopDHT(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, true},
		{"closing dht", fields{Dht: new(DHTClient)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.stopDHT(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.stopDHT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_stopSocket(t *testing.T) {
	type fields struct {
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
		Interface       TAP
		Peers           *Swarm
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty test", fields{}, true},
		{"nil connection error", fields{UDPSocket: new(Network)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,

				Hash:         tt.fields.Hash,
				Interface:    tt.fields.Interface,
				Swarm:        tt.fields.Peers,
				HolePunching: tt.fields.HolePunching,
				ProxyManager: tt.fields.ProxyManager,
				outboundIP:   tt.fields.outboundIP,
				UsePMTU:      tt.fields.UsePMTU,
			}
			if err := p.stopSocket(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.stopSocket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
