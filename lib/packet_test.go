package ptp

import (
	"net"
	"os"
	"testing"

	"github.com/google/gofuzz"
	"sync"
	"reflect"
)

func TestUnmarshalARP(t *testing.T) {
	arp := new(ARPPacket)

	b1 := make([]byte, 7)
	err := arp.UnmarshalARP(b1)
	if err == nil {
		t.Error(err)
	}

	b2 := make([]byte, 23)
	err1 := arp.UnmarshalARP(b2)
	if err1 == nil {
		t.Error(err1)
	}

	f := fuzz.New().NilChance(0.5)
	var a struct {
		Ht   uint16
		Pt   uint16
		Hal  uint8
		Ipl  uint8
		O    Operation
		Shwa net.HardwareAddr
		Sip  net.IP
		Thwa net.HardwareAddr
		Tip  net.IP
	}

	f.Fuzz(&a)

	arp.HardwareType = 2
	arp.ProtocolType = 0x0800
	arp.HardwareAddrLength = 6
	arp.IPLength = 4
	arp.Operation = 2
	arp.SenderHardwareAddr = a.Shwa
	arp.SenderIP = a.Sip
	arp.TargetHardwareAddr = a.Thwa
	arp.TargetIP = a.Tip

	b, _ := arp.MarshalBinary()

	file, e := os.Create("MarshalBinary")
	if e != nil {
		t.Error("Unable to create file:", e)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(b)

	err3 := arp.UnmarshalARP(b)
	if err3 != nil {
		t.Error(err3)
	}
}
/*
Generated TestPeerToPeer_handlePacket
Generated TestPeerToPeer_handlePacketIPv4
Generated TestPeerToPeer_handlePacketIPv6
Generated TestPeerToPeer_handlePARCUniversalPacket
Generated TestPeerToPeer_handleRARPPacket
Generated TestPeerToPeer_handle8021qPacket
Generated TestPeerToPeer_handlePPPoEDiscoveryPacket
Generated TestPeerToPeer_handlePPPoESessionPacket
Generated TestPeerToPeer_handlePacketARP
Generated TestPeerToPeer_handlePacketLLDP
Generated TestARPPacket_String
Generated TestARPPacket_MarshalBinary
Generated TestARPPacket_UnmarshalARP
Generated TestARPPacket_NewPacket
package ptp

import (
	"net"
	"os"
	"reflect"
	"sync"
	"testing"

	"github.com/google/gofuzz"
)
*/

/*
func TestUnmarshalARP(t *testing.T) {
	arp := new(ARPPacket)

	b1 := make([]byte, 7)
	err := arp.UnmarshalARP(b1)
	if err == nil {
		t.Error(err)
	}

	b2 := make([]byte, 23)
	err1 := arp.UnmarshalARP(b2)
	if err1 == nil {
		t.Error(err1)
	}

	f := fuzz.New().NilChance(0.5)
	var a struct {
		Ht   uint16
		Pt   uint16
		Hal  uint8
		Ipl  uint8
		O    Operation
		Shwa net.HardwareAddr
		Sip  net.IP
		Thwa net.HardwareAddr
		Tip  net.IP
	}

	f.Fuzz(&a)

	arp.HardwareType = 2
	arp.ProtocolType = 0x0800
	arp.HardwareAddrLength = 6
	arp.IPLength = 4
	arp.Operation = 2
	arp.SenderHardwareAddr = a.Shwa
	arp.SenderIP = a.Sip
	arp.TargetHardwareAddr = a.Thwa
	arp.TargetIP = a.Tip

	b, _ := arp.MarshalBinary()

	file, e := os.Create("MarshalBinary")
	if e != nil {
		t.Error("Unable to create file:", e)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(b)

	err3 := arp.UnmarshalARP(b)
	if err3 != nil {
		t.Error(err3)
	}
}
*/

func TestPeerToPeer_handlePacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePacketIPv4(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePacketIPv4(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePacketIPv6(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePacketIPv6(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePARCUniversalPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePARCUniversalPacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handleRARPPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handleRARPPacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handle8021qPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handle8021qPacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePPPoEDiscoveryPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePPPoEDiscoveryPacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePPPoESessionPacket(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePPPoESessionPacket(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePacketARP(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePacketARP(tt.args.contents, tt.args.proto)
		})
	}
}

func TestPeerToPeer_handlePacketLLDP(t *testing.T) {
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
		contents []byte
		proto    int
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
			p.handlePacketLLDP(tt.args.contents, tt.args.proto)
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *ARPPacket
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
			got, err := p.NewPacket(tt.args.op, tt.args.srcHW, tt.args.srcIP, tt.args.dstHW, tt.args.dstIP)
			if (err != nil) != tt.wantErr {
				t.Errorf("ARPPacket.NewPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ARPPacket.NewPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}
