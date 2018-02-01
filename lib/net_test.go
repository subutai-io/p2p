package ptp

import (
	"bytes"
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

func TestGetProxyAttributes(t *testing.T) {
	bytes := []byte("12345678910")
	var wait1 uint16 = 14641
	var wait2 uint16 = 13108
	get1, get2 := GetProxyAttributes(bytes)
	if get1 != wait1 && get2 != wait2 {
		t.Error(get1, get2)
	}
}
