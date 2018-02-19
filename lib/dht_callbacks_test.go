package ptp

import (
	"net"
	"reflect"
	"testing"
)

func TestPacketConnect(t *testing.T) {
	dht := new(DHTClient)
	pct := new(DHTPacket)
	id := "123456789101112131415161415161415161"
	pct.Id = id
	err := dht.packetConnect(pct)
	if err != nil && dht.ID != id {
		t.Error("Error")
	}
	id2 := "12345"
	pct.Id = id2
	err2 := dht.packetConnect(pct)
	if err2 == nil {
		t.Error("Wrong value of identificator")
	}
}

func TestPacketDHCP(t *testing.T) {
	dht := new(DHTClient)
	pct := new(DHTPacket)
	ipstr := "192.168.21.1"
	netstr := "24"
	ip := net.IP(ipstr)
	mask := net.IPMask(netstr)
	net := new(net.IPNet)
	net.IP = ip
	net.Mask = mask
	pct.Data = ipstr
	pct.Extra = netstr
	err := dht.packetDHCP(pct)
	if err != nil && reflect.DeepEqual(dht.IP, ip) && dht.Network != net {
		t.Error("Error")
	}

	pct2 := new(DHTPacket)
	pct2.Data = "-"
	pct2.Extra = "-"
	err2 := dht.packetDHCP(pct2)
	if err2 == nil {
		t.Error("Error")
	}
}

func TestPacketError(t *testing.T) {
	dht := new(DHTClient)
	pct := new(DHTPacket)
	data1 := ""
	pct.Data = data1
	err := dht.packetError(pct)
	if err != nil {
		t.Error("err")
	}
	data2 := "Warning"
	pct.Data = data2
	err2 := dht.packetError(pct)
	if err2 != nil {
		t.Error("Error")
	}
	data := "Error"
	pct.Data = data
	err3 := dht.packetError(pct)
	if err3 != nil {
		t.Error("Error")
	}
}
