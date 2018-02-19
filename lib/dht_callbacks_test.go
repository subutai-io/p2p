package ptp

import (
	"testing"
)

func TestPacketConnect(t *testing.T) {
	ptp := new(PeerToPeer)
	pct := new(DHTPacket)
	id := "12345"
	pct.Id = id
	err := ptp.packetConnect(pct)
	if err == nil {
		t.Error("Wrong value of identificator")
	}
}

func TestPacketError(t *testing.T) {
	ptp := new(PeerToPeer)
	pct := new(DHTPacket)
	data1 := ""
	pct.Data = data1
	err := ptp.packetError(pct)
	if err != nil {
		t.Error("err")
	}
	data2 := "Warning"
	pct.Data = data2
	err2 := ptp.packetError(pct)
	if err2 != nil {
		t.Error("Error")
	}
	data := "Error"
	pct.Data = data
	err3 := ptp.packetError(pct)
	if err3 != nil {
		t.Error("Error")
	}
}
