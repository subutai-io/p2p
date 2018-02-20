package ptp

import "testing"

// import (
// 	"testing"
// )

// func TestCompose(t *testing.T) {
// 	var dht DHTClient
// 	t1 := dht.Compose("ping", "00000000-1111-2222-3333-444444444444", "QUERY", "ARGUMENT")
// 	if t1 != "d1:a8:ARGUMENT1:c4:ping1:i36:00000000-1111-2222-3333-4444444444441:p0:1:q5:QUERYe" {
// 		t.Errorf("dht.Compose failed (1)")
// 	}
// 	t2 := dht.Compose("", "", "", "")
// 	if t2 != "" {
// 		t.Errorf("dht.Compose failed (2)")
// 	}
// }

// func TestExtract(t *testing.T) {
// 	var m = "d1:a0:1:c4:ping1:i36:00000000-1111-2222-3333-4444444444441:p0:1:q1:0e"
// 	var dht DHTClient
// 	result, err := dht.Extract([]byte(m))
// 	if result.ID != "00000000-1111-2222-3333-444444444444" {
// 		t.Errorf("Failed to extract DHT message")
// 	}
// 	if err != nil {
// 		t.Errorf("Error during DHT message extraction")
// 	}
// }

// func TestEncodeRequest(t *testing.T) {
// 	var dht DHTClient
// 	t1 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test1", Query: "Query1", Arguments: "Argument1"})
// 	if t1 != "d1:a9:Argument11:c5:Test11:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query1e" {
// 		t.Errorf("EncodeRequest failed (1)")
// 	}
// 	t2 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test2", Query: "Query2", Arguments: "Argument2"})
// 	if t2 != "d1:a9:Argument21:c5:Test21:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query2e" {
// 		t.Errorf("EncodeRequest failed (2)")
// 	}
// 	t3 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test3", Query: "Query3", Arguments: "Argument3"})
// 	if t3 != "d1:a9:Argument31:c5:Test31:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query3e" {
// 		t.Errorf("EncodeRequest failed (3)")
// 	}
// }

func TestInit(t *testing.T) {
	dht := new(DHTClient)
	err := dht.Init("hash")
	if err == nil && dht.NetworkHash != "hash" {
		t.Error(err)
	}
}

func TestRead(t *testing.T) {
	dht := new(DHTClient)
	pct := new(DHTPacket)
	packetChan := make(chan *DHTPacket, 1)
	packetChan <- pct
	dht.IncomingData = packetChan
	get, _ := dht.read()
	if get != pct {
		t.Error("Error")
	}
	pct2 := new(DHTPacket)
	pct2 = nil
	packetChan <- pct2
	dht.IncomingData = packetChan
	get2, err := dht.read()
	if get2 != nil {
		t.Error(err)
	}
}

func TestWaitID(t *testing.T) {
	dht := new(DHTClient)
	dht.ID = "12345"
	err := dht.WaitID()
	if err == nil {
		t.Error(err)
	}
	dht.ID = "123456789101112131415161718192212223"
	err2 := dht.WaitID()
	if err2 != nil {
		t.Error(err2)
	}
}
