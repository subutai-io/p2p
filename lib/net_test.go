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
