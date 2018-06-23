package ptp

import (
	"bytes"
	"net"
	"reflect"
	"testing"
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
	}{}
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
		{"disposed", fields{disposed: true}, true},
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
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	tests := []struct {
		name   string
		fields fields
		want   *net.UDPAddr
	}{
		{"nil addr", fields{}, nil},
		{"addr", fields{addr: addr}, addr},
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
	addr1, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	conn, _ := net.ListenUDP("udp4", addr1)
	if conn != nil {
		defer conn.Close()
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"bad addr", fields{}, args{"n", 999999}, true},
		{"wrong listen", fields{}, args{"127.0.0.1", 2555}, true},
		{"wrong listen", fields{}, args{"127.0.0.1", 2556}, false},
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
	addr1, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	conn, _ := net.ListenUDP("udp4", addr1)
	if conn != nil {
		defer conn.Close()
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"GetPort()", fields{}, -1},
		{"GetPort()", fields{conn: conn}, 2555},
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
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1")
	addr2, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	conn, err := net.ListenUDP("udp4", addr)
	if err == nil {
		defer conn.Close()
	} else {
		t.Errorf("%s\n", err)
	}
	msg, _ := CreateMessageStatic(MsgTypePing, []byte{})
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{"nil conn", fields{}, args{}, -1, true},
		{"nil message", fields{conn: &net.UDPConn{}}, args{}, 0, true},
		{"empty conn", fields{conn: &net.UDPConn{}}, args{msg, addr2}, 0, true},
		{"empty conn", fields{conn: conn}, args{msg, addr2}, 10, false},
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
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1")
	addr2, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	conn, err := net.ListenUDP("udp4", addr)
	if err == nil {
		defer conn.Close()
	} else {
		t.Errorf("%s\n", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{"nil conn", fields{}, args{}, -1, true},
		{"empty conn", fields{conn: &net.UDPConn{}}, args{}, 0, true},
		{"non-conn", fields{conn: conn}, args{[]byte{0x01}, addr2}, 1, false},
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
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1")
	conn, err := net.ListenUDP("udp4", addr)
	if err == nil {
		defer conn.Close()
	} else {
		t.Errorf("%s\n", err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"nil conn", fields{}, args{}, true},
		{"disposed", fields{disposed: true, conn: conn}, args{}, false},
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
			if err := uc.Listen(tt.args.receivedCallback); (err != nil) != tt.wantErr {
				t.Errorf("Network.Listen() error = %v, wantErr %v", err, tt.wantErr)
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
	addr, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:2555")
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		t.Errorf("Failed to start connection: %v", err)
	}
	// if conn != nil {
	// 	defer conn.Close()
	// }
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil conn", fields{}, true},
		{"non-nil conn", fields{conn: conn}, false},
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
			if err := uc.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Network.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
