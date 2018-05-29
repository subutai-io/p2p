package ptp

import (
	"bytes"
	"net"
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
