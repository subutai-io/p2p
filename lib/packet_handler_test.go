package ptp

import (
	"net"
	"sync"
	"testing"
)

func TestPeerToPeer_HandleXpeerPingMessage(t *testing.T) {
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

	msg := new(P2PMessage)
	msg2, _ := CreateMessageStatic(MsgTypeXpeerPing, []byte("r123456789012345678901234567890123456"))
	msg3, _ := CreateMessageStatic(MsgTypeXpeerPing, []byte("q123456789012345678901234567890123456"))
	srcAddr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1234")

	pl := new(PeerList)
	pl.Init()
	pl.Update("123456789012345678901234567890123456", &NetworkPeer{ID: "123456789012345678901234567890123456"})

	proxy := new(ProxyManager)
	proxy.init()

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"t1", fields{}, args{}},
		{"t1", fields{}, args{msg: msg}},
		{"t1", fields{Peers: pl}, args{msg: msg, srcAddr: srcAddr}},
		{"t1", fields{Peers: pl}, args{msg: msg2, srcAddr: srcAddr}},
		{"t1", fields{Peers: pl, ProxyManager: proxy}, args{msg: msg3, srcAddr: srcAddr}},
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
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
				UsePMTU:         tt.fields.UsePMTU,
			}
			p.HandleXpeerPingMessage(tt.args.msg, tt.args.srcAddr)
		})
	}
}
