package ptp

import (
	"reflect"
	"testing"
)

func TestGet(t *testing.T) {
	l := new(PeerList)
	get := l.Get()
	var wait map[string]*NetworkPeer

	if reflect.DeepEqual(get, wait) {
		t.Error("wait, get", wait, get)
	}
}

func TestLength(t *testing.T) {
	l := new(PeerList)
	count := 0
	for i := 0; i < len(l.peers); i++ {
		count++
	}
	get := l.Length()
	if get != count {
		t.Errorf("Error. Wait: %v, get: %v", count, get)
	}
}

func TestGetPeer(t *testing.T) {
	l := new(PeerList)
	wait := new(NetworkPeer)
	wait = l.peers["1"]
	get := l.GetPeer("1")
	if get != wait {
		t.Error("wait, get: ", wait, get)
	}
}
