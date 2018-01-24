package ptp

import "testing"

func TestGetType(t *testing.T) {
	dht := new(DHTPacket)
	types := [...]int{
		int(DHTPacketType_Undefined),
		int(DHTPacketType_Connect),
		int(DHTPacketType_Forward),
		int(DHTPacketType_Find),
		int(DHTPacketType_Node),
		int(DHTPacketType_Ping),
		int(DHTPacketType_RegisterProxy),
		int(DHTPacketType_RequestProxy),
		int(DHTPacketType_ReportProxy),
		int(DHTPacketType_BadProxy),
		int(DHTPacketType_Proxy),
		int(DHTPacketType_Notify),
		int(DHTPacketType_ReportLoad),
		int(DHTPacketType_Stop),
		int(DHTPacketType_Unknown),
		int(DHTPacketType_DHCP),
		int(DHTPacketType_Error),
		int(DHTPacketType_Unsupported),
		int(DHTPacketType_State),
	}

	for i := 0; i < len(types); i++ {
		dht.Type = DHTPacketType(i)
		get := dht.GetType()
		if get != DHTPacketType(types[i]) {
			t.Errorf("Error. Wait %v, get %v ", dht.Type, get)
		}
	}
}

func TestGetId(t *testing.T) {
	dht := new(DHTPacket)
	for i := 0; i < 10; i++ {
		dht.Id = string(i)
		get := dht.GetId()
		if get != dht.Id {
			t.Errorf("Error. Wait %v, get %v", dht.Id, get)
		}
	}
}

func TestGetInfohash(t *testing.T) {
	dht := new(DHTPacket)

	Infohashs := [...]string{
		"infohash",
		"",
		"123456",
		"",
	}

	for i := 0; i < len(Infohashs); i++ {
		dht.Infohash = Infohashs[i]
		get := dht.GetInfohash()
		if get != Infohashs[i] {
			t.Error("Error.")
		}
	}
}
