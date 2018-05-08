package ptp

import (
	"bytes"
	"net"
	"testing"
	"reflect"
	"sync"
)

func TestSerialize(t *testing.T) {
	p := new(P2PMessageHeader)
	var wait = make([]byte, 10)
	for i := 0; i < 10; i++ {
		wait[i] = 0
	}
	get := p.Serialize()
	if !bytes.EqualFold(wait, get) {
		t.Errorf("Error. Wait: %v, get: %v", wait, get)
	}
}

func TestP2PMessageHeaderFromBytes(t *testing.T) {
	bytes1 := []byte("12")
	get1, _ := P2PMessageHeaderFromBytes(bytes1)
	if get1 != nil {
		t.Error("Error")
	}
	bytes2 := []byte("12345678910111213140")
	wait := new(P2PMessageHeader)
	wait.Magic = 12594
	wait.Type = 13108
	wait.Length = 13622
	wait.NetProto = 14136
	wait.SerializedLen = 12337
	get2, _ := P2PMessageHeaderFromBytes(bytes2)
	if get2.Magic != wait.Magic && get2.Type != wait.Type && get2.Length != wait.Length && get2.NetProto != wait.NetProto {
		t.Error("Error. get: ", get2)
	}
	bytes := []byte("12345")
	get, err := P2PMessageHeaderFromBytes(bytes)
	if get != nil {
		t.Error(err)
	}
}

func TestDisposed(t *testing.T) {
	nt := new(Network)
	nt.disposed = true
	get := nt.Disposed()
	if !get {
		t.Error("Error.Return wrong value.")
	}
	nt.disposed = false
	get2 := nt.Disposed()
	if get2 {
		t.Error("Error.Return wrong value")
	}
}

func TestAddr(t *testing.T) {
	nt := new(Network)
	get := nt.Addr()
	if get != nil {
		t.Error("Error")
	}
	nt.addr, _ = net.ResolveUDPAddr("udp4", "127.0.0.1:24")
	get2 := nt.Addr()
	if get2 != nt.addr {
		t.Error("Error")
	}
}

func TestSendRawBytes(t *testing.T) {
	nt := new(Network)
	nt.conn = nil
	bytes := []byte("12345")
	addr, _ := net.ResolveUDPAddr("network", "127.0.0.1")
	get, _ := nt.SendRawBytes(bytes, addr)
	if get != -1 {
		t.Errorf("Error.Wait: %v, get: %v", -1, get)
	}
}
/*
Generated TestP2PMessageHeader_Serialize
Generated TestP2PMessage_Serialize
Generated TestP2PMessageFromBytes
Generated TestPeerToPeer_CreateMessage
Generated TestCreateMessageStatic
Generated TestNetwork_Stop
Generated TestNetwork_Disposed
Generated TestNetwork_Addr
Generated TestNetwork_Init
Generated TestNetwork_KeepAlive
Generated TestNetwork_GetPort
Generated TestNetwork_Listen
Generated TestNetwork_SendMessage
Generated TestNetwork_SendRawBytes
package ptp

import (
	"bytes"
	"net"
	"reflect"
	"sync"
	"testing"
)
*/

/*
func TestSerialize(t *testing.T) {
	p := new(P2PMessageHeader)
	var wait = make([]byte, 10)
	for i := 0; i < 10; i++ {
		wait[i] = 0
	}
	get := p.Serialize()
	if !bytes.EqualFold(wait, get) {
		t.Errorf("Error. Wait: %v, get: %v", wait, get)
	}
}

func TestP2PMessageHeaderFromBytes(t *testing.T) {
	bytes1 := []byte("12")
	get1, _ := P2PMessageHeaderFromBytes(bytes1)
	if get1 != nil {
		t.Error("Error")
	}
	bytes2 := []byte("12345678910111213140")
	wait := new(P2PMessageHeader)
	wait.Magic = 12594
	wait.Type = 13108
	wait.Length = 13622
	wait.NetProto = 14136
	wait.SerializedLen = 12337
	get2, _ := P2PMessageHeaderFromBytes(bytes2)
	if get2.Magic != wait.Magic && get2.Type != wait.Type && get2.Length != wait.Length && get2.NetProto != wait.NetProto {
		t.Error("Error. get: ", get2)
	}
	bytes := []byte("12345")
	get, err := P2PMessageHeaderFromBytes(bytes)
	if get != nil {
		t.Error(err)
	}
}

func TestDisposed(t *testing.T) {
	nt := new(Network)
	nt.disposed = true
	get := nt.Disposed()
	if !get {
		t.Error("Error.Return wrong value.")
	}
	nt.disposed = false
	get2 := nt.Disposed()
	if get2 {
		t.Error("Error.Return wrong value")
	}
}

func TestAddr(t *testing.T) {
	nt := new(Network)
	get := nt.Addr()
	if get != nil {
		t.Error("Error")
	}
	nt.addr, _ = net.ResolveUDPAddr("udp4", "127.0.0.1:24")
	get2 := nt.Addr()
	if get2 != nt.addr {
		t.Error("Error")
	}
}

func TestSendRawBytes(t *testing.T) {
	nt := new(Network)
	nt.conn = nil
	bytes := []byte("12345")
	addr, _ := net.ResolveUDPAddr("network", "127.0.0.1")
	get, _ := nt.SendRawBytes(bytes, addr)
	if get != -1 {
		t.Errorf("Error.Wait: %v, get: %v", -1, get)
	}
}
*/

func TestP2PMessageHeader_Serialize(t *testing.T) {
	type fields struct {
		Magic         uint16
		Type          uint16
		Length        uint16
		SerializedLen uint16
		NetProto      uint16
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &P2PMessageHeader{
				Magic:         tt.fields.Magic,
				Type:          tt.fields.Type,
				Length:        tt.fields.Length,
				SerializedLen: tt.fields.SerializedLen,
				NetProto:      tt.fields.NetProto,
			}
			if got := v.Serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("P2PMessageHeader.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP2PMessage_Serialize(t *testing.T) {
	type fields struct {
		Header *P2PMessageHeader
		Data   []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &P2PMessage{
				Header: tt.fields.Header,
				Data:   tt.fields.Data,
			}
			if got := v.Serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("P2PMessage.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestP2PMessageFromBytes(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *P2PMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := P2PMessageFromBytes(tt.args.bytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("P2PMessageFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("P2PMessageFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_CreateMessage(t *testing.T) {
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
		msgType MsgType
		payload []byte
		proto   uint16
		encrypt bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *P2PMessage
		wantErr bool
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
			got, err := p.CreateMessage(tt.args.msgType, tt.args.payload, tt.args.proto, tt.args.encrypt)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.CreateMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateMessageStatic(t *testing.T) {
	type args struct {
		msgType MsgType
		payload []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *P2PMessage
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateMessageStatic(tt.args.msgType, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateMessageStatic() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateMessageStatic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_Stop(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			uc.Stop()
		})
	}
}

func TestNetwork_Disposed(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			if got := uc.Disposed(); got != tt.want {
				t.Errorf("Network.Disposed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_Addr(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   *net.UDPAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			if got := uc.Addr(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Network.Addr() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_Init(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	type args struct {
		host string
		port int
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
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			if err := uc.Init(tt.args.host, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("Network.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNetwork_KeepAlive(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	type args struct {
		addr *net.UDPAddr
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
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			uc.KeepAlive(tt.args.addr)
		})
	}
}

func TestNetwork_GetPort(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			if got := uc.GetPort(); got != tt.want {
				t.Errorf("Network.GetPort() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_Listen(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	type args struct {
		receivedCallback UDPReceivedCallback
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
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			uc.Listen(tt.args.receivedCallback)
		})
	}
}

func TestNetwork_SendMessage(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	type args struct {
		msg     *P2PMessage
		dstAddr *net.UDPAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			got, err := uc.SendMessage(tt.args.msg, tt.args.dstAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Network.SendMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNetwork_SendRawBytes(t *testing.T) {
	type fields struct {
		host       string
		port       int
		remotePort int
		addr       *net.UDPAddr
		conn       *net.UDPConn
		inBuffer   [4096]byte
		disposed   bool
	}
	type args struct {
		bytes   []byte
		dstAddr *net.UDPAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &Network{
				host:       tt.fields.host,
				port:       tt.fields.port,
				remotePort: tt.fields.remotePort,
				addr:       tt.fields.addr,
				conn:       tt.fields.conn,
				inBuffer:   tt.fields.inBuffer,
				disposed:   tt.fields.disposed,
			}
			got, err := uc.SendRawBytes(tt.args.bytes, tt.args.dstAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Network.SendRawBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Network.SendRawBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
