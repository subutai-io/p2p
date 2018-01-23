package ptp

import (
	"testing"
)

func TestLength(t *testing.T) {
	l := new(PeerList)
	count := 0
	for i := 0; i < len(l.peers); i++ {
		count++
	}
	get := l.Length()
	if get != count {
		t.Errorf("Error. Wait: %v, get: %v", &count, &get)
	}
}
