package ptp

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/mdlayher/ethernet"
)

func TestPeerToPeer_handlePacket(t *testing.T) {
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
		contents []byte
		proto    int
	}

	ptp := new(PeerToPeer)
	ptp.setupHandlers()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty packet test", fields{}, args{}, true},
		{"existing packet type", fields{PacketHandlers: ptp.PacketHandlers}, args{proto: int(PacketLLDP)}, false},
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
			if err := p.handlePacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePacketIPv4(t *testing.T) {
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
		contents []byte
		proto    int
	}

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, uint16(2048))

	p0 := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xcd, 0xce, 0xcf}
	p1 := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb}
	p1 = append(p1, buf.Bytes()...)
	p1 = append(p1, []byte{0x01, 0x02}...)

	pl0 := new(PeerList)
	pl0.Init()

	socket0 := new(Network)
	socket0.Init("127.0.0.1", 1234)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, true},
		{"bad ether frame", fields{}, args{p0, 0}, true},
		{"good ether frame", fields{Peers: pl0, UDPSocket: socket0}, args{p1, 0}, false},
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
			if err := p.handlePacketIPv4(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePacketIPv4() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePacketIPv6(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handlePacketIPv6(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePacketIPv6() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePARCUniversalPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handlePARCUniversalPacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePARCUniversalPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handleRARPPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handleRARPPacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handleRARPPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handle8021qPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handle8021qPacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handle8021qPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePPPoEDiscoveryPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handlePPPoEDiscoveryPacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePPPoEDiscoveryPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePPPoESessionPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handlePPPoESessionPacket(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePPPoESessionPacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePacketARP(t *testing.T) {
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
		contents []byte
		proto    int
	}

	ip := "170.187.204.221"

	mac0, _ := net.ParseMAC("00:00:00:00:00:00")

	inf0, _ := newTAP("ip", "10.10.10.1", "00:00:00:00:00:00", "255.255.255.255", 1500, false)
	inf0.file, _ = os.OpenFile("/tmp/p2p-test-interface", os.O_CREATE|os.O_RDWR, 0700)
	defer inf0.file.Close()

	p0 := []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xcd, 0xce, 0xcf}
	payload0 := []byte{0x01, 0x02, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff,
		0x01, 0x02, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0xaa, 0xbb, 0xcc, 0xdd, 0xee, 0xff}
	p1 := append(p0, payload0...)

	pl0 := new(PeerList)
	pl0.Init()

	pl1 := new(PeerList)
	pl1.Init()
	pl1.tableIPID[ip] = "broken-id"

	pl2 := new(PeerList)
	pl2.Init()
	pl2.tableIPID[ip] = "peer-id0"
	pl2.peers["peer-id0"] = &NetworkPeer{
		ID:          "peer-id0",
		PeerLocalIP: net.ParseIP(ip),
	}

	pl3 := new(PeerList)
	pl3.Init()
	pl3.tableIPID[ip] = "peer-id0"
	pl3.peers["peer-id0"] = &NetworkPeer{
		ID:          "peer-id0",
		PeerLocalIP: net.ParseIP(ip),
		PeerHW:      mac0,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty test", fields{}, args{}, true},
		{"arp unmarshal failed", fields{}, args{p0, 0}, true},
		{"nil peer list", fields{}, args{p1, 0}, true},
		{"peer not found", fields{Peers: pl0}, args{p1, 0}, true},
		{"broken id", fields{Peers: pl1}, args{p1, 0}, true},
		{"no hw address", fields{Peers: pl2}, args{p1, 0}, true},
		{"empty hw address", fields{Peers: pl3}, args{p1, 0}, true},
		{"interface exists", fields{Peers: pl3, Interface: inf0}, args{p1, 0}, false},
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
			if err := p.handlePacketARP(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePacketARP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_handlePacketLLDP(t *testing.T) {
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
		contents []byte
		proto    int
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
			if err := p.handlePacketLLDP(tt.args.contents, tt.args.proto); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.handlePacketLLDP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestARPPacket_String(t *testing.T) {
	type fields struct {
		HardwareType       uint16
		ProtocolType       uint16
		HardwareAddrLength uint8
		IPLength           uint8
		Operation          Operation
		SenderHardwareAddr net.HardwareAddr
		SenderIP           net.IP
		TargetHardwareAddr net.HardwareAddr
		TargetIP           net.IP
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ARPPacket{
				HardwareType:       tt.fields.HardwareType,
				ProtocolType:       tt.fields.ProtocolType,
				HardwareAddrLength: tt.fields.HardwareAddrLength,
				IPLength:           tt.fields.IPLength,
				Operation:          tt.fields.Operation,
				SenderHardwareAddr: tt.fields.SenderHardwareAddr,
				SenderIP:           tt.fields.SenderIP,
				TargetHardwareAddr: tt.fields.TargetHardwareAddr,
				TargetIP:           tt.fields.TargetIP,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("ARPPacket.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARPPacket_MarshalBinary(t *testing.T) {
	type fields struct {
		HardwareType       uint16
		ProtocolType       uint16
		HardwareAddrLength uint8
		IPLength           uint8
		Operation          Operation
		SenderHardwareAddr net.HardwareAddr
		SenderIP           net.IP
		TargetHardwareAddr net.HardwareAddr
		TargetIP           net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ARPPacket{
				HardwareType:       tt.fields.HardwareType,
				ProtocolType:       tt.fields.ProtocolType,
				HardwareAddrLength: tt.fields.HardwareAddrLength,
				IPLength:           tt.fields.IPLength,
				Operation:          tt.fields.Operation,
				SenderHardwareAddr: tt.fields.SenderHardwareAddr,
				SenderIP:           tt.fields.SenderIP,
				TargetHardwareAddr: tt.fields.TargetHardwareAddr,
				TargetIP:           tt.fields.TargetIP,
			}
			got, err := p.MarshalBinary()
			if (err != nil) != tt.wantErr {
				t.Errorf("ARPPacket.MarshalBinary() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ARPPacket.MarshalBinary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARPPacket_UnmarshalARP(t *testing.T) {
	type fields struct {
		HardwareType       uint16
		ProtocolType       uint16
		HardwareAddrLength uint8
		IPLength           uint8
		Operation          Operation
		SenderHardwareAddr net.HardwareAddr
		SenderIP           net.IP
		TargetHardwareAddr net.HardwareAddr
		TargetIP           net.IP
	}
	type args struct {
		b []byte
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
			p := &ARPPacket{
				HardwareType:       tt.fields.HardwareType,
				ProtocolType:       tt.fields.ProtocolType,
				HardwareAddrLength: tt.fields.HardwareAddrLength,
				IPLength:           tt.fields.IPLength,
				Operation:          tt.fields.Operation,
				SenderHardwareAddr: tt.fields.SenderHardwareAddr,
				SenderIP:           tt.fields.SenderIP,
				TargetHardwareAddr: tt.fields.TargetHardwareAddr,
				TargetIP:           tt.fields.TargetIP,
			}
			if err := p.UnmarshalARP(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("ARPPacket.UnmarshalARP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestARPPacket_NewPacket(t *testing.T) {
	type fields struct {
		HardwareType       uint16
		ProtocolType       uint16
		HardwareAddrLength uint8
		IPLength           uint8
		Operation          Operation
		SenderHardwareAddr net.HardwareAddr
		SenderIP           net.IP
		TargetHardwareAddr net.HardwareAddr
		TargetIP           net.IP
	}
	type args struct {
		op    Operation
		srcHW net.HardwareAddr
		srcIP net.IP
		dstHW net.HardwareAddr
		dstIP net.IP
	}

	hw0, _ := net.ParseMAC("00:00:00:00:00:01")
	hw1 := net.HardwareAddr{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
	ip := net.ParseIP("10.10.10.1")

	res := &ARPPacket{
		HardwareType:       1,
		ProtocolType:       uint16(ethernet.EtherTypeIPv4),
		HardwareAddrLength: uint8(6),
		IPLength:           uint8(4),
		Operation:          0,
		SenderHardwareAddr: hw0,
		SenderIP:           ip,
		TargetHardwareAddr: hw0,
		TargetIP:           ip,
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ARPPacket
		wantErr bool
	}{
		{"small src", fields{}, args{}, nil, true},
		{"small dst", fields{}, args{srcHW: hw0}, nil, true},
		{"different size", fields{}, args{srcHW: hw0, dstHW: hw1}, nil, true},
		{"invalid src ip", fields{}, args{srcHW: hw0, dstHW: hw0}, nil, true},
		{"invalid dst ip", fields{}, args{srcHW: hw0, dstHW: hw0, srcIP: ip}, nil, true},
		{"passing", fields{}, args{srcHW: hw0, dstHW: hw0, srcIP: ip, dstIP: ip}, res, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ARPPacket{
				HardwareType:       tt.fields.HardwareType,
				ProtocolType:       tt.fields.ProtocolType,
				HardwareAddrLength: tt.fields.HardwareAddrLength,
				IPLength:           tt.fields.IPLength,
				Operation:          tt.fields.Operation,
				SenderHardwareAddr: tt.fields.SenderHardwareAddr,
				SenderIP:           tt.fields.SenderIP,
				TargetHardwareAddr: tt.fields.TargetHardwareAddr,
				TargetIP:           tt.fields.TargetIP,
			}
			got, err := p.NewPacket(tt.args.op, tt.args.srcHW, tt.args.srcIP, tt.args.dstHW, tt.args.dstIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("ARPPacket.NewPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if !reflect.DeepEqual(got, tt.want) {
			if got != nil && got.SenderHardwareAddr != nil && tt.want.SenderHardwareAddr != nil && !bytes.Equal(got.SenderHardwareAddr, tt.want.SenderHardwareAddr) {
				t.Errorf("ARPPacket.NewPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}
