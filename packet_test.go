package ptp

import (
	"net"
	"os"
	"testing"

	"github.com/google/gofuzz"
)

func TestUnmarshalARP(t *testing.T) {
	arp := new(ARPPacket)

	b1 := make([]byte, 7)
	err := arp.UnmarshalARP(b1)
	if err == nil {
		t.Error(err)
	}

	b2 := make([]byte, 23)
	err1 := arp.UnmarshalARP(b2)
	if err1 == nil {
		t.Error(err1)
	}

	f := fuzz.New().NilChance(0.5)
	var a struct {
		Ht   uint16
		Pt   uint16
		Hal  uint8
		Ipl  uint8
		O    Operation
		Shwa net.HardwareAddr
		Sip  net.IP
		Thwa net.HardwareAddr
		Tip  net.IP
	}

	f.Fuzz(&a)

	arp.HardwareType = 2
	arp.ProtocolType = 0x0800
	arp.HardwareAddrLength = 6
	arp.IPLength = 4
	arp.Operation = 2
	arp.SenderHardwareAddr = a.Shwa
	arp.SenderIP = a.Sip
	arp.TargetHardwareAddr = a.Thwa
	arp.TargetIP = a.Tip

	b, _ := arp.MarshalBinary()

	file, e := os.Create("MarshalBinary")
	if e != nil {
		t.Error("Unable to create file:", e)
		os.Exit(1)
	}
	defer file.Close()

	file.Write(b)

	err3 := arp.UnmarshalARP(b)
	if err3 != nil {
		t.Error(err3)
	}
}
