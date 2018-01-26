package ptp

import (
	"bytes"
	"reflect"
	"testing"
)

func TestString(t *testing.T) {
	var dht DHTPacketType
	types := [...]DHTPacketType{
		DHTPacketType_Undefined,
		DHTPacketType_Connect,
		DHTPacketType_Forward,
		DHTPacketType_Find,
		DHTPacketType_Node,
		DHTPacketType_Ping,
		DHTPacketType_RegisterProxy,
		DHTPacketType_RequestProxy,
		DHTPacketType_ReportProxy,
		DHTPacketType_BadProxy,
		DHTPacketType_Proxy,
		DHTPacketType_Notify,
		DHTPacketType_ReportLoad,
		DHTPacketType_Stop,
		DHTPacketType_Unknown,
		DHTPacketType_DHCP,
		DHTPacketType_Error,
		DHTPacketType_Unsupported,
		DHTPacketType_State,
	}

	names := [...]string{
		"Undefined",
		"Connect",
		"Forward",
		"Find",
		"Node",
		"Ping",
		"RegisterProxy",
		"RequestProxy",
		"ReportProxy",
		"BadProxy",
		"Proxy",
		"Notify",
		"ReportLoad",
		"Stop",
		"Unknown",
		"DHCP",
		"Error",
		"Unsupported",
		"State",
	}

	for i := 0; i < len(types); i++ {
		dht = types[i]
		get := dht.String()
		if get != names[i] {
			t.Errorf("Error. Wait %v, get %v", names[i], get)
		}
	}
}

func TestEnumDescriptor(t *testing.T) {
	var dht DHTPacketType
	ints := make([]int, 1)
	get1, get2 := dht.EnumDescriptor()
	if !bytes.EqualFold(get1, fileDescriptor0) && !reflect.DeepEqual(get2, ints) {
		t.Errorf("get1: %v, get2: %v", get1, get2)
	}
}

func TestDescriptor(t *testing.T) {
	dht := new(DHTPacket)
	get1, get2 := dht.Descriptor()
	i := make([]int, 1)
	if !bytes.EqualFold(get1, fileDescriptor0) && !reflect.DeepEqual(get2, i) {
		t.Error("Error", get1, get2)
	}
}

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

	dht2 := new(DHTPacket)
	dht2 = nil
	wait := DHTPacketType_Undefined
	get := dht2.GetType()
	if get != wait {
		t.Errorf("Error. Wait %v, get %v", wait, get)
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

	dht2 := new(DHTPacket)
	dht2 = nil
	wait := ""
	get2 := dht2.GetId()
	if get2 != wait {
		t.Errorf("Error. Wait %v, get %v", wait, get2)
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
			t.Errorf("Error.Wait %v, get %v", dht.Data, get)
		}
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	wait := ""
	get2 := dht2.GetInfohash()
	if get2 != wait {
		t.Errorf("Error. Wait %v, get %v", wait, get2)
	}
}

func TestGetData(t *testing.T) {
	dht := new(DHTPacket)
	data := [...]string{
		"data1",
		"string",
		"12345",
		"",
	}
	for i := 0; i < len(data); i++ {
		dht.Data = data[i]
		get := dht.GetData()
		if get != data[i] {
			t.Errorf("Error. Wait %v, get %v", data[i], get)
		}
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	wait := ""
	get2 := dht2.GetData()
	if get2 != wait {
		t.Errorf("Error. Wait %v, get %v", wait, get2)
	}
}

func TestGetQuery(t *testing.T) {
	dht := new(DHTPacket)
	queries := [...]string{
		"query",
		"",
		"12345",
	}

	for i := 0; i < len(queries); i++ {
		dht.Query = queries[i]
		get := dht.GetQuery()
		if get != queries[i] {
			t.Errorf("Error. Wait %v, get %v", queries[i], get)
		}
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	get2 := dht2.GetQuery()
	if get2 != "" {
		t.Errorf("Error. Wait %v, get %v", "", get2)
	}
}

func TestGetArguments(t *testing.T) {
	dht := new(DHTPacket)
	var argums []string
	argums = append(argums, "string")
	argums = append(argums, "12345")
	argums = append(argums, "")
	dht.Arguments = argums
	get := dht.GetArguments()
	if !reflect.DeepEqual(get, argums) {
		t.Errorf("Error. Wait %v, get %v", dht.Arguments, get)
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	get2 := dht2.GetArguments()
	if get2 != nil {
		t.Errorf("Error. Wait %v, get %v", nil, get)
	}
}

func TestGetProxies(t *testing.T) {
	dht := new(DHTPacket)
	var proxies []string
	proxies = append(proxies, "proxy1")
	proxies = append(proxies, "proxy2")
	proxies = append(proxies, "12345")
	proxies = append(proxies, "")

	dht.Proxies = proxies
	get := dht.GetProxies()
	if !reflect.DeepEqual(get, proxies) {
		t.Errorf("Error. Wait %v, get %v", dht.Proxies, get)
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	get2 := dht2.GetProxies()
	if get2 != nil {
		t.Errorf("Error. Wait %v, get %v", nil, get)
	}
}

func TestGetExtra(t *testing.T) {
	dht := new(DHTPacket)
	var extras = [...]string{
		"extra",
		"12345",
		"",
	}
	for i := 0; i < len(extras); i++ {
		dht.Extra = extras[i]
		get := dht.GetExtra()
		if get != dht.Extra {
			t.Errorf("Error. Wait %v, get %v", dht.Extra, get)
		}
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	get2 := dht2.GetExtra()
	if get2 != "" {
		t.Errorf("Error. Wait %v, get %v", "", get2)
	}
}

func TestGetPayload(t *testing.T) {
	dht := new(DHTPacket)
	payloads := make([]byte, 5)
	dht.Payload = payloads
	get := dht.GetPayload()
	if !bytes.EqualFold(get, dht.Payload) {
		t.Errorf("Error. Wait %v, get %v", dht.Payload, get)
	}

	dht2 := new(DHTPacket)
	dht2 = nil
	get2 := dht2.GetPayload()
	if get2 != nil {
		t.Errorf("Error. Wait %v, get %v", nil, get)
	}
}
