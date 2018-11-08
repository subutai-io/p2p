package ptp

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

func TestPeerToPeer_HandleP2PMessage(t *testing.T) {
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
		count    int
		srcAddr  *net.UDPAddr
		err      error
		rcvBytes []byte
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty fields", fields{}, args{}, true},
		{"entry error", fields{}, args{err: errors.New("test error")}, true},
		{"2 bytes header", fields{}, args{count: 2, rcvBytes: []byte{0x1, 0x2}}, true},
		{"decrypt failed", fields{}, args{}, true},
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
			if err := p.HandleP2PMessage(tt.args.count, tt.args.srcAddr, tt.args.err, tt.args.rcvBytes); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleP2PMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleNotEncryptedMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	msg0 := new(P2PMessage)
	msg0.Header = new(P2PMessageHeader)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil msg", fields{}, args{}, true},
		{"nil header", fields{}, args{msg: new(P2PMessage)}, true},
		{"nil source", fields{}, args{msg: msg0}, true},
		{"passing", fields{}, args{msg: msg0, srcAddr: &net.UDPAddr{}}, false},
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
			if err := p.HandleNotEncryptedMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleNotEncryptedMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandlePingMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	msg0 := new(P2PMessage)
	msg0.Header = new(P2PMessageHeader)
	msg0.Data = []byte("a:b")
	msg1 := new(P2PMessage)
	msg1.Header = new(P2PMessageHeader)
	msg1.Data = []byte("192.168.0.1:1234")
	msg2 := new(P2PMessage)
	msg2.Header = new(P2PMessageHeader)
	msg2.Data = []byte("192.168.0.1:0")

	udp1, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	proxy1 := new(proxyServer)
	proxy1.Init(udp1)

	pm0 := new(ProxyManager)
	pm0.init()
	pm1 := new(ProxyManager)
	pm1.init()
	pm1.proxies[udp1.String()] = proxy1

	socket0 := new(Network)
	socket1 := new(Network)
	socket1.remotePort = 1
	socket2 := new(Network)
	socket2.remotePort = 1234

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil msg", fields{}, args{}, true},
		{"nil source", fields{}, args{msg: msg0}, true},
		{"nil proxy manager", fields{}, args{msg: msg0, srcAddr: &net.UDPAddr{}}, true},
		{"nil udp socket", fields{ProxyManager: pm0}, args{msg: msg0, srcAddr: &net.UDPAddr{}}, true},
		{"bad addr", fields{ProxyManager: pm0, UDPSocket: socket0}, args{msg: msg0, srcAddr: &net.UDPAddr{}}, false},
		{"bad addr>real proxy", fields{ProxyManager: pm1, UDPSocket: socket0}, args{msg: msg0, srcAddr: udp1}, false},
		{"empty port", fields{ProxyManager: pm1, UDPSocket: socket0}, args{msg: msg2, srcAddr: udp1}, false},
		{"port translation", fields{ProxyManager: pm1, UDPSocket: socket1}, args{msg: msg1, srcAddr: udp1}, false},
		{"same remote port", fields{ProxyManager: pm1, UDPSocket: socket2}, args{msg: msg1, srcAddr: udp1}, false},
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
			if err := p.HandlePingMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandlePingMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleXpeerPingMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	msg0 := new(P2PMessage)
	msg1 := new(P2PMessage)
	msg1.Data = []byte("q")
	msg2 := new(P2PMessage)
	msg2.Data = []byte("q123e4567-e89b-12d3-a456-426655440000")
	msg3 := new(P2PMessage)
	msg3.Data = []byte("q123e4567-e89b-12d3-a456-426655440000192.168.0.1:1234")
	msg4 := new(P2PMessage)
	msg4.Data = []byte("r192.168.0.1:1234")
	msg5 := new(P2PMessage)
	msg5.Data = []byte("somerandomdata")

	src1, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	kip0 := []*net.UDPAddr{src1}

	ep0 := new(Endpoint)
	ep0.Addr = src1

	pl0 := new(PeerList)
	pl0.Init()
	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:       "123e4567-e89b-12d3-a456-426655440000",
		KnownIPs: kip0,
	}
	pl2 := new(PeerList)
	pl2.Init()
	pl2.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:            "123e4567-e89b-12d3-a456-426655440000",
		State:         PeerStateConnected,
		RemoteState:   PeerStateConnected,
		EndpointsHeap: []*Endpoint{ep0},
	}
	pl3 := new(PeerList)
	pl3.Init()
	pl3.peers["123e4567-e89b-12d3-a456-426655440000"] = nil

	socket0 := new(Network)

	crypto0 := Crypto{}
	crypto0.Active = true

	proxy0 := new(proxyServer)
	proxy0.Addr = src1
	proxy0.Endpoint = src1

	pm0 := new(ProxyManager)
	pm0.init()
	pm1 := new(ProxyManager)
	pm1.init()
	pm1.proxies[src1.String()] = proxy0

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil msg", fields{}, args{}, true},
		{"nil source", fields{}, args{msg: msg0}, true},
		{"nil peer list", fields{}, args{msg: msg0, srcAddr: src1}, true},
		{"nil socket", fields{Peers: pl0}, args{msg: msg0, srcAddr: src1}, true},
		{"nil proxy manager", fields{Peers: pl0, UDPSocket: socket0}, args{msg: msg0, srcAddr: src1}, true},
		{"empty payload", fields{Peers: pl0, UDPSocket: socket0, ProxyManager: pm0}, args{msg: msg0, srcAddr: src1}, true},
		{"q>small payload", fields{Peers: pl0, UDPSocket: socket0, ProxyManager: pm0}, args{msg: msg1, srcAddr: src1}, true},
		{"q>msg create fail", fields{Peers: pl0, UDPSocket: socket0, Crypter: crypto0, ProxyManager: pm0}, args{msg: msg2, srcAddr: src1}, true},
		{"q>unknown endpoint", fields{Peers: pl0, UDPSocket: socket0, ProxyManager: pm0}, args{msg: msg2, srcAddr: src1}, true},
		{"q>known ip", fields{Peers: pl1, UDPSocket: socket0, ProxyManager: pm0}, args{msg: msg2, srcAddr: src1}, false},
		{"q>over proxy", fields{Peers: pl2, UDPSocket: socket0, ProxyManager: pm1}, args{msg: msg3, srcAddr: src1}, false},
		{"r>nil peer", fields{Peers: pl3, UDPSocket: socket0, ProxyManager: pm1}, args{msg: msg4, srcAddr: src1}, true},
		{"r>passing", fields{Peers: pl2, UDPSocket: socket0, ProxyManager: pm1}, args{msg: msg4, srcAddr: src1}, false},
		{"broken message", fields{Peers: pl2, UDPSocket: socket0, ProxyManager: pm1}, args{msg: msg5, srcAddr: src1}, true},
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
			if err := p.HandleXpeerPingMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleXpeerPingMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleIntroMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	msg0 := &P2PMessage{
		Data: []byte(",,,,,,,,,,,,,"),
	}
	msg1 := &P2PMessage{
		Data: []byte("1,00:11:22:33:44:55,10.10.10.1,192.168.0.1:1234"),
	}
	msg2 := &P2PMessage{
		Data: []byte("123e4567-e89b-12d3-a456-426655440000,00:11:22:33:44:55,10.10.10.1,192.168.0.1:1234"),
	}

	src0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")
	src1, _ := net.ResolveUDPAddr("udp4", "192.168.1.1:2345")
	src2, _ := net.ResolveUDPAddr("udp4", "192.168.1.2:3456")

	mac0, _ := net.ParseMAC("00:11:22:33:44:55")

	pl0 := new(PeerList)
	pl0.Init()
	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID: "123e4567-e89b-12d3-a456-426655440000",
	}
	pl1.peers["a23e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:     "a23e4567-e89b-12d3-a456-426655440000",
		PeerHW: mac0,
	}
	pl1.peers["b23e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID: "b23e4567-e89b-12d3-a456-426655440000",
		EndpointsHeap: []*Endpoint{
			&Endpoint{
				Addr: src0,
			},
		},
	}
	pl1.peers["c23e4567-e89b-12d3-a456-426655440000"] = nil

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil message", fields{}, args{}, true},
		{"nil source addr", fields{}, args{msg: new(P2PMessage)}, true},
		{"nil peer list", fields{}, args{msg: new(P2PMessage), srcAddr: &net.UDPAddr{}}, true},
		{"parse failed", fields{Peers: pl0}, args{msg: msg0, srcAddr: src0}, true},
		{"broken id length", fields{Peers: pl0}, args{msg: msg1, srcAddr: src0}, true},
		{"peer not found", fields{Peers: pl0}, args{msg: msg2, srcAddr: src0}, true},
		{"peer the same", fields{Peers: pl1}, args{msg: msg2, srcAddr: src0}, false},
		{"mac duplicate", fields{Peers: pl1}, args{msg: msg2, srcAddr: src0}, false},
		{"ep duplicate", fields{Peers: pl1}, args{msg: msg2, srcAddr: src1}, false},
		{"passing", fields{Peers: pl1}, args{msg: msg2, srcAddr: src2}, false},
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
			if err := p.HandleIntroMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleIntroMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleIntroRequestMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	src0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")
	src1, _ := net.ResolveUDPAddr("udp4", "192.168.0.2:3456")

	pl0 := new(PeerList)
	pl1 := new(PeerList)
	pl1.Init()
	pl1.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:       "123e4567-e89b-12d3-a456-426655440000",
		KnownIPs: []*net.UDPAddr{src0},
	}
	pl2 := new(PeerList)
	pl2.Init()
	pl2.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:       "123e4567-e89b-12d3-a456-426655440000",
		KnownIPs: []*net.UDPAddr{},
	}

	dht0 := new(DHTClient)

	socket0 := new(Network)
	socket1 := new(Network)
	socket1.Init("127.0.0.1", 21345)

	msg0 := new(P2PMessage)
	msg1 := new(P2PMessage)
	msg1.Data = []byte("123e4567-e89b-12d3-a456-426655440000")

	inf0, _ := newTAP("iptool", "10.10.10.1", "00:11:22:33:44:55", "255.255.255.0", 1500, false)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil msg", fields{}, args{}, true},
		{"nil source", fields{}, args{msg: new(P2PMessage)}, true},
		{"nil peer list", fields{}, args{msg: new(P2PMessage), srcAddr: &net.UDPAddr{}}, true},
		{"nil dht", fields{Peers: new(PeerList)}, args{msg: new(P2PMessage), srcAddr: &net.UDPAddr{}}, true},
		{"nil udp socket", fields{Peers: new(PeerList), Dht: new(DHTClient)}, args{msg: new(P2PMessage), srcAddr: &net.UDPAddr{}}, true},
		{"short payload", fields{Peers: pl0, Dht: dht0, UDPSocket: socket0}, args{msg: msg0, srcAddr: src0}, true},
		{"peer not found", fields{Peers: pl0, Dht: dht0, UDPSocket: socket0}, args{msg: msg1, srcAddr: src0}, true},
		{"failed intro message", fields{Peers: pl1, Dht: dht0, UDPSocket: socket0}, args{msg: msg1, srcAddr: src0}, true},
		{"non-existing ep", fields{Peers: pl1, Dht: dht0, UDPSocket: socket0, Interface: inf0}, args{msg: msg1, srcAddr: src1}, true},
		{"existing ep", fields{Peers: pl1, Dht: dht0, UDPSocket: socket0, Interface: inf0}, args{msg: msg1, srcAddr: src0}, true},
		{"passing", fields{Peers: pl2, Dht: dht0, UDPSocket: socket1, Interface: inf0}, args{msg: msg1, srcAddr: src0}, false},
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
			if err := p.HandleIntroRequestMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleIntroRequestMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleProxyMessage(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	src0, _ := net.ResolveUDPAddr("udp4", "192.168.0.1:1234")

	msg0 := &P2PMessage{
		Data: []byte("b:a"),
	}
	msg1 := &P2PMessage{
		Data: []byte("192.168.0.1:1234"),
	}

	proxy0 := new(proxyServer)
	proxy0.Status = proxyConnecting

	pm1 := new(ProxyManager)
	pm2 := new(ProxyManager)
	pm2.init()
	pm2.proxies["192.168.0.1:1234"] = proxy0

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil message", fields{}, args{}, true},
		{"nil source", fields{}, args{msg: &P2PMessage{}}, true},
		{"nil proxy manager", fields{}, args{msg: &P2PMessage{}, srcAddr: &net.UDPAddr{}}, true},
		{"bad udp address", fields{ProxyManager: pm1}, args{msg: msg0, srcAddr: src0}, true},
		{"failed activation", fields{ProxyManager: pm1}, args{msg: msg1, srcAddr: src0}, true},
		{"passed activation", fields{ProxyManager: pm2}, args{msg: msg1, srcAddr: src0}, false},
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
			if err := p.HandleProxyMessage(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleProxyMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleBadTun(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
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
			if err := p.HandleBadTun(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleBadTun() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_HandleLatency(t *testing.T) {
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
		msg     *P2PMessage
		srcAddr *net.UDPAddr
	}

	src0 := &net.UDPAddr{}
	src1, _ := net.ResolveUDPAddr("udp4", "192.168.1.2:3456")

	ts0, _ := time.Now().MarshalBinary()

	d0 := append(LatencyRequestHeader, []byte{0x1, 0x2, 0x3, 0x4, 0x5, 0x6}...)
	d0 = append(d0, []byte("123e4567-e89b-12d3-a456-426655440000")...)
	d0 = append(d0, ts0...)

	msg0 := &P2PMessage{}
	msg1 := &P2PMessage{
		Data: append(LatencyProxyHeader, []byte("bad time for covertion")...),
	}
	msg2 := &P2PMessage{
		Data: append(LatencyProxyHeader, ts0...),
	}
	msg3 := &P2PMessage{
		Data: append(LatencyRequestHeader, []byte("shortpayload")...),
	}
	fmt.Println(len(d0))
	msg4 := &P2PMessage{
		Data: d0,
	}

	proxy0 := &proxyServer{
		Addr: src0,
	}

	pm0 := &ProxyManager{}
	pm1 := &ProxyManager{}
	pm1.init()
	pm1.proxies[src0.String()] = proxy0

	pl0 := &PeerList{}
	pl1 := &PeerList{}
	pl1.Init()
	pl1.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:       "123e4567-e89b-12d3-a456-426655440000",
		Endpoint: nil,
	}
	pl2 := &PeerList{}
	pl2.Init()
	pl2.peers["123e4567-e89b-12d3-a456-426655440000"] = &NetworkPeer{
		ID:       "123e4567-e89b-12d3-a456-426655440000",
		Endpoint: src1,
	}

	cr0 := Crypto{
		Active: true,
	}

	socket0 := new(Network)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil message", fields{}, args{nil, nil}, true},
		{"nil source", fields{}, args{msg0, nil}, true},
		{"nil proxy manager", fields{}, args{msg0, src0}, true},
		{"nil peer list", fields{ProxyManager: pm0}, args{msg0, src0}, true},
		{"nil socket", fields{ProxyManager: pm0, Peers: pl0}, args{msg0, src0}, true},
		{"short payload", fields{ProxyManager: pm0, Peers: pl0, UDPSocket: socket0}, args{msg0, src0}, true},
		{"proxy>bad time", fields{ProxyManager: pm0, Peers: pl0, UDPSocket: socket0}, args{msg1, src0}, true},
		{"proxy>set failed", fields{ProxyManager: pm0, Peers: pl0, UDPSocket: socket0}, args{msg2, src0}, true},
		{"proxy>set passed", fields{ProxyManager: pm1, Peers: pl0, UDPSocket: socket0}, args{msg2, src0}, false},
		{"request>short", fields{ProxyManager: pm0, Peers: pl0, UDPSocket: socket0}, args{msg3, src0}, true},
		{"request>unknown peer", fields{ProxyManager: pm0, Peers: pl0, UDPSocket: socket0}, args{msg4, src1}, true},
		{"request>nil endpoint", fields{ProxyManager: pm0, Peers: pl1, UDPSocket: socket0}, args{msg4, src1}, true},
		// This passes because CreateMessageStatic never fails
		{"request>failed response", fields{Crypter: cr0, ProxyManager: pm0, Peers: pl2, UDPSocket: socket0}, args{msg4, src1}, false},
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
			if err := p.HandleLatency(tt.args.msg, tt.args.srcAddr); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.HandleLatency() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
