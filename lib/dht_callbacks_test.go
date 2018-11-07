package ptp

import (
	"net"
	"sync"
	"testing"

	"github.com/subutai-io/p2p/protocol"
)

func TestPeerToPeer_setupTCPCallbacks(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"nil dht", fields{}},
		{"setup", fields{Dht: new(DHTClient)}},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			p.setupTCPCallbacks()
		})
	}
}

func TestPeerToPeer_packetBadProxy(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"bad proxy", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetBadProxy(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetBadProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetConnect(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil dht", fields{}, args{}, true},
		{"nil packet", fields{Dht: new(DHTClient)}, args{}, true},
		{"small id", fields{Dht: new(DHTClient)}, args{&protocol.DHTPacket{Id: "123"}}, true},
		{"normal id", fields{Dht: new(DHTClient)}, args{&protocol.DHTPacket{Id: "123e4567-e89b-12d3-a456-426655440000"}}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetConnect(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetConnect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetDHCP(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil dht", fields{}, args{}, true},
		{"nil packet", fields{Dht: new(DHTClient)}, args{}, true},
		{"broken cidr", fields{Dht: new(DHTClient)}, args{&protocol.DHTPacket{Data: "d", Extra: "e"}}, true},
		{"normal packet", fields{Dht: new(DHTClient)}, args{&protocol.DHTPacket{Data: "192.168.0.1", Extra: "32"}}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetDHCP(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetDHCP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetError(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"empty level", fields{}, args{&protocol.DHTPacket{}}, false},
		{"warning level", fields{}, args{&protocol.DHTPacket{Data: "Warning"}}, false},
		{"error level", fields{}, args{&protocol.DHTPacket{Data: "Error"}}, false},
		{"unknown level", fields{}, args{&protocol.DHTPacket{Data: "unknown"}}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetError(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetFind(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}

	SetMinLogLevel(Debug)

	tdht := new(DHTClient)
	tdht.ID = "testid"

	pl0 := new(PeerList)
	pl0.Init()

	knownIP, _ := net.ResolveUDPAddr("udp4", "192.168.1.2:3456")

	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["testid"] = &NetworkPeer{
		Proxies: []*net.UDPAddr{knownIP, knownIP},
	}
	pl1.peers["testid2"] = &NetworkPeer{
		KnownIPs: []*net.UDPAddr{knownIP},
		Proxies:  []*net.UDPAddr{knownIP},
	}

	lip0 := []net.IP{}
	lip0 = append(lip0, net.ParseIP("192.168.1.2"))

	proxy1 := []string{"b:p"}
	proxy2 := []string{"192.168.0.1:1234", "192.168.0.1:1234"}
	proxy3 := []string{"192.168.1.2:3456"}

	proxyServer := new(proxyServer)
	proxyServer.Endpoint = knownIP
	proxyServer.Addr = knownIP

	pm1 := new(ProxyManager)
	pm2 := new(ProxyManager)
	pm2.init()
	pm2.proxies[knownIP.String()] = proxyServer

	f1 := fields{}
	f2 := fields{Dht: tdht}
	f3 := fields{Dht: tdht, Peers: pl0}
	f4 := fields{Dht: tdht, Peers: pl0, ProxyManager: pm1}
	f5 := fields{Dht: tdht, Peers: pl1, LocalIPs: lip0, ProxyManager: pm1}
	f6 := fields{Dht: tdht, Peers: pl1, LocalIPs: lip0, ProxyManager: pm2}

	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", f1, args{}, true},
		{"nil dht", f1, args{new(protocol.DHTPacket)}, true},
		{"empty peer list", f2, args{new(protocol.DHTPacket)}, false},
		{"skip self", f2, args{&protocol.DHTPacket{Data: tdht.ID, Arguments: []string{"arg1"}, Extra: "skip"}}, false},
		{"nil peer list", f2, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{"arg1"}, Extra: "skip"}}, true},
		{"nil proxy manager", f3, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{"arg1"}, Extra: "skip"}}, true},
		{"new peer>bad ip", f4, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{"arg1"}, Extra: "skip"}}, false},
		{"new peer>known ip", f4, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{knownIP.String(), knownIP.String()}, Extra: "skip"}}, false},
		{"new peer>local ip", f5, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{knownIP.String()}, Extra: "skip"}}, false},
		{"new peer>bad proxy", f5, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy1}}, false},
		{"new peer>existing proxy", f5, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy2}}, false},
		{"new peer>own proxy", f6, args{&protocol.DHTPacket{Data: "ttt", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy3}}, false},
		{"existing peer>empty ip", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{""}, Extra: "skip"}}, false},
		{"existing peer>bad ip", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{"arg1"}, Extra: "skip"}}, false},
		{"existing peer>known ip", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String(), knownIP.String()}, Extra: "skip"}}, false},
		{"existing peer>local ip", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String()}, Extra: "skip"}}, false},
		{"existing peer>empty proxy", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: []string{""}}}, false},
		{"existing peer>bad proxy", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy1}}, false},
		{"existing peer>existing proxy", f5, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy2}}, false},
		{"existing peer>own proxy", f6, args{&protocol.DHTPacket{Data: "testid2", Arguments: []string{knownIP.String()}, Extra: "skip", Proxies: proxy3}}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetFind(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetFind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetForward(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetForward(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetForward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetNode(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}

	p1 := &protocol.DHTPacket{}
	p2 := &protocol.DHTPacket{
		Data:      "unknown-id",
		Arguments: []string{""},
	}
	p3 := &protocol.DHTPacket{
		Data:      "test-id-1",
		Arguments: []string{""},
	}
	p4 := &protocol.DHTPacket{
		Data:      "test-id-1",
		Arguments: []string{"b/p"},
	}
	p5 := &protocol.DHTPacket{
		Data:      "test-id-1",
		Arguments: []string{"192.168.0.1:1234"},
	}

	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["test-id-1"] = new(NetworkPeer)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"nil peer list", fields{}, args{p1}, true},
		{"empty arguments", fields{Peers: new(PeerList)}, args{p1}, true},
		{"unknown peer", fields{Peers: pl1}, args{p2}, true},
		{"empty addr", fields{Peers: pl1}, args{p3}, false},
		{"bad addr", fields{Peers: pl1}, args{p4}, false},
		{"passing test", fields{Peers: pl1}, args{p5}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetNode(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetNotify(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetNotify(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetNotify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetPing(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetPing(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetPing() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetProxy(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}

	n1 := new(Network)
	n1.Init("127.0.0.1", 1234)

	pm1 := new(ProxyManager)
	pm1.init()

	f1 := fields{
		UDPSocket:    n1,
		ProxyManager: pm1,
	}

	f2 := fields{
		UDPSocket:    n1,
		ProxyManager: pm1,
		Dht:          &DHTClient{},
	}

	p1 := &protocol.DHTPacket{
		Proxies: []string{"s:b"},
	}

	p2 := &protocol.DHTPacket{
		Proxies: []string{"127.0.0.1:1111", "127.0.0.1:1111"},
	}

	p3 := &protocol.DHTPacket{
		Proxies: []string{"127.0.0.1:1111"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"nil socket", fields{}, args{new(protocol.DHTPacket)}, true},
		{"nil proxy manager", fields{UDPSocket: new(Network)}, args{new(protocol.DHTPacket)}, true},
		{"nil dht", f1, args{new(protocol.DHTPacket)}, true},
		{"bad proxy addr", f2, args{p1}, false},
		{"two same proxies", f2, args{p2}, false},
		{"new proxy", f2, args{p3}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetProxy(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetRequestProxy(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}

	p1 := &protocol.DHTPacket{
		Proxies: []string{"b:a"},
	}

	p2 := &protocol.DHTPacket{
		Proxies: []string{"127.0.0.1:1234"},
	}

	p3 := &protocol.DHTPacket{
		Proxies: []string{"127.0.0.1:1234"},
		Data:    "test-peer",
	}

	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["test-peer"] = &NetworkPeer{}

	f1 := fields{
		Peers: pl1,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil peer list", fields{}, args{}, true},
		{"bad proxy addr", fields{Peers: new(PeerList)}, args{p1}, false},
		{"non existing peer", f1, args{p2}, false},
		{"existing peer", f1, args{p3}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetRequestProxy(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetRequestProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetReportProxy(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetReportProxy(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetReportProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetRegisterProxy(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"ok", fields{}, args{&protocol.DHTPacket{Data: "OK"}}, false},
		{"not ok", fields{}, args{&protocol.DHTPacket{Data: "!ok"}}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetRegisterProxy(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetRegisterProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetReportLoad(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetReportLoad(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetReportLoad() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetState(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}

	pl1 := new(PeerList)
	pl1.Init()

	pl2 := new(PeerList)
	pl2.Init()
	pl2.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{}

	f1 := fields{
		Peers: pl1,
	}

	f2 := fields{
		Peers: pl2,
	}

	p1 := &protocol.DHTPacket{
		Data: "short",
	}

	p2 := &protocol.DHTPacket{
		Data:  "123e4567-e89b-12d3-a456-426655440000",
		Extra: "",
	}

	p3 := &protocol.DHTPacket{
		Data:  "123e4567-e89b-12d3-a456-426655440000",
		Extra: "error",
	}

	p4 := &protocol.DHTPacket{
		Data:  "123e4567-e89b-12d3-a456-426655440000",
		Extra: "4",
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"nil peer list", fields{}, args{new(protocol.DHTPacket)}, true},
		{"short data", f1, args{p1}, true},
		{"empty extra", f1, args{p2}, true},
		{"broken extra", f1, args{p3}, true},
		{"peer not exists", f1, args{p4}, false},
		{"peer exists", f2, args{p4}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetState(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetStop(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetStop(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetUnknown(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}

	p1 := &protocol.DHTPacket{
		Data: "DHCP",
	}

	p2 := &protocol.DHTPacket{
		Data: "Anything, but DHCP",
	}

	pm1 := new(ProxyManager)
	pm1.init()

	dht := new(DHTClient)
	dht.Init("")

	inf, _ := newTAP("", "192.168.0.1", "00:11:22:33:44:55", "", 1500, false)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"nil dht", fields{}, args{p1}, true},
		{"nil proxy manager", fields{Dht: dht}, args{p1}, true},
		{"nil interface", fields{Dht: dht, ProxyManager: pm1}, args{p1}, true},
		{"normal data", fields{Dht: dht, ProxyManager: pm1, Interface: inf}, args{p1}, false},
		{"refuse data", fields{Dht: dht, ProxyManager: pm1, Interface: inf}, args{p2}, true},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetUnknown(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetUnknown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_packetUnsupported(t *testing.T) {
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
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
		UsePMTU         bool
	}
	type args struct {
		packet *protocol.DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil packet", fields{}, args{}, true},
		{"nil dht", fields{}, args{new(protocol.DHTPacket)}, true},
		{"passing", fields{Dht: new(DHTClient)}, args{new(protocol.DHTPacket)}, false},
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
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			if err := p.packetUnsupported(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.packetUnsupported() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
