package ptp

import (
	"net"
	"reflect"
	"testing"
)

func TestUpdateTables(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.updateTables("500", "192.168.1.1", "01:02:03:04:05:06")
	_, exists := l.tableIPID["500"]
	_, exist := l.tableMacID["01:02:03:04:05:06"]
	if !exists && !exist {
		t.Error("Error. Can't update peer")
	}
}

func TestDeleteTables(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.tableIPID["800"] = "192.168.8.8"
	l.tableMacID["800"] = "01:02:03:04:05:06"
	l.deleteTables("800", "800")
	_, exists := l.tableIPID["800"]
	_, exist := l.tableMacID["800"]
	if exist && exists {
		t.Error("Error")
	}
}

func TestGet(t *testing.T) {
	l := new(PeerList)
	np1 := new(NetworkPeer)
	np2 := new(NetworkPeer)
	l.Init()
	l.peers["444"] = np1
	l.peers["445"] = np2
	get := l.Get()
	var wait map[string]*NetworkPeer
	wait = make(map[string]*NetworkPeer)
	wait["444"] = np1
	wait["445"] = np2
	if !reflect.DeepEqual(get, wait) {
		t.Error("wait, get", wait, get)
	}
}

func TestLength(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.peers["77"] = new(NetworkPeer)
	l.peers["78"] = new(NetworkPeer)
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
	l.Init()
	l.peers["9"] = new(NetworkPeer)
	l.peers["99"] = new(NetworkPeer)

	get1 := l.GetPeer("9")
	if get1 != l.peers["9"] {
		t.Error("Error")
	}
	get2 := l.GetPeer("-1")
	if get2 != nil {
		t.Error("Error")
	}
}

func TestGetEndpointAndProxy(t *testing.T) {
	l := new(PeerList)
	get, i, err := l.GetEndpointAndProxy("01:02:03:04:05:06")
	if get != nil && i != 0 {
		t.Error(err)
	}
	l.tableMacID = make(map[string]string)
	l.tableMacID["10:11:12:13:14:15"] = "888"
	l.peers = make(map[string]*NetworkPeer)
	np := new(NetworkPeer)
	addr, _ := net.ResolveUDPAddr("udp4", "192.168.44.1:24")
	np.Endpoint = addr
	l.peers["888"] = np
	get2, i2, err2 := l.GetEndpointAndProxy("10:11:12:13:14:15")
	if get2 != addr && i2 != 0 {
		t.Error(err2)
	}
}
