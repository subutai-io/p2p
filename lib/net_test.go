package ptp

import (
	"bytes"
	"net"
	"testing"
)

func TestSerialize(t *testing.T) {
	p := new(P2PMessageHeader)
	var wait = make([]byte, 18)
	for i := 0; i < 18; i++ {
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
	wait.ProxyID = 14641
	wait.SerializedLen = 12337
	wait.Complete = 12593
	wait.ID = 12849
	wait.Seq = 13105
	get2, _ := P2PMessageHeaderFromBytes(bytes2)
	if get2.Magic != wait.Magic && get2.Type != wait.Type && get2.Length != wait.Length && get2.NetProto != wait.NetProto && get2.ProxyID != wait.SerializedLen && get2.Complete != wait.Complete && get2.ID != wait.ID && get2.Seq != wait.Seq {
		t.Error("Error. get: ", get2)
	}
	bytes := []byte("12345")
	get, err := P2PMessageHeaderFromBytes(bytes)
	if get != nil {
		t.Error(err)
	}
}

func TestGetProxyAttributes(t *testing.T) {
	bytes := []byte("12345678910")
	var wait1 uint16 = 14641
	var wait2 uint16 = 13108
	get1, get2 := GetProxyAttributes(bytes)
	if get1 != wait1 && get2 != wait2 {
		t.Error(get1, get2)
	}
}

func TestCreatePingP2PMessage(t *testing.T) {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MagicCookie
	msg.Header.Type = uint16(MsgTypePing)
	msg.Header.NetProto = 0
	msg.Header.Length = uint16(len("1"))
	msg.Header.Complete = 1
	msg.Header.ID = 0
	msg.Data = []byte("12345")

	get := CreatePingP2PMessage("12345")

	if get.Header.Magic != msg.Header.Magic && get.Header.Type != msg.Header.Type && get.Header.NetProto != msg.Header.NetProto && get.Header.Length != msg.Header.Length && get.Header.Complete != msg.Header.Complete && !bytes.EqualFold(get.Data, msg.Data) {
		t.Error("Error in func CreatePingP2PMessage")
	}
}

func TestCreateConfP2PMessage(t *testing.T) {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MagicCookie
	msg.Header.Type = uint16(MsgTypeConf)
	msg.Header.NetProto = 0
	msg.Header.Length = uint16(len("1"))
	msg.Header.Complete = 1
	msg.Header.ID = 1
	msg.Header.Seq = 2
	msg.Data = []byte("1")

	get := CreateConfP2PMessage(1, 2)
	if get.Header.Magic != msg.Header.Magic && get.Header.Type != msg.Header.Type && get.Header.NetProto != msg.Header.NetProto && get.Header.Length != msg.Header.Length && get.Header.Complete != msg.Header.Complete && !bytes.EqualFold(get.Data, msg.Data) {
		t.Error("Error in func CreateConfP2PMessage")
	}
}

func TestCreateProxyP2PMessage(t *testing.T) {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MagicCookie
	msg.Header.Type = uint16(MsgTypeProxy)
	msg.Header.NetProto = 2
	msg.Header.Length = uint16(len("12345"))
	msg.Header.Complete = 1
	msg.Header.ProxyID = uint16(4)
	msg.Header.ID = 0
	msg.Data = []byte("12345")

	get := CreateProxyP2PMessage(4, "12345", 2)
	if get.Header.Magic != msg.Header.Magic && get.Header.Type != msg.Header.Type && get.Header.NetProto != msg.Header.NetProto && get.Header.Length != msg.Header.Length && get.Header.Complete != msg.Header.Complete && !bytes.EqualFold(get.Data, msg.Data) {
		t.Error("Error in func CreateProxyP2PMessage")
	}
}

func TestCreateBadTunnelP2PMessage(t *testing.T) {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MagicCookie
	msg.Header.Type = uint16(MsgTypeBadTun)
	msg.Header.NetProto = 1
	msg.Header.Length = uint16(len("rem"))
	msg.Header.ProxyID = uint16(2)
	msg.Header.Complete = 1
	msg.Header.ID = 0
	msg.Data = []byte("rem")

	get := CreateBadTunnelP2PMessage(2, 1)
	if get.Header.Magic != msg.Header.Magic && get.Header.Type != msg.Header.Type && get.Header.NetProto != msg.Header.NetProto && get.Header.Length != msg.Header.Length && get.Header.Complete != msg.Header.Complete && !bytes.EqualFold(get.Data, msg.Data) {
		t.Error("Error in func CreateBadTunnelP2PMessage")
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
	nt.addr, _ = net.ResolveUDPAddr("network", "127.0.0.1")
	get := nt.Addr()
	if get != nil {
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
