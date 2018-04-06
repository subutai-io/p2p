package ptp

import (
	"testing"
)

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
	if err != nil {
		t.Errorf("Error in TCPInit")
	}
	err1 := dht.Init("hash")
	if err1 != nil {
		t.Errorf("Error. Wait %v, get %v", "dht.cdn.subut.ai:6881", dht.Routers)
	}
}

// func TestProducePacket(t *testing.T) {
// 	dht := new(DHTClient)

// 	dataStr := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
// 	smallArr := []string{dataStr, dataStr}
// 	medArr := []string{dataStr, dataStr, dataStr, dataStr}
// 	bigArr := []string{dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr, dataStr}

// 	empty, e1 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping})
// 	if e1 != nil {
// 		t.Errorf("Failed to produce empty packet: %s", e1)
// 	}
// 	if len(empty) != 1 {
// 		t.Errorf("Wrong length for empty packet: %d", len(empty))
// 	}

// 	s1, e2 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: smallArr})
// 	if e2 != nil {
// 		t.Errorf("Failed to produce s1 packet: %s", e2)
// 	}
// 	if len(s1) != 1 {
// 		t.Errorf("Wrong length for small packet: %d", len(s1))
// 	}

// 	s2, e3 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: smallArr, Arguments: smallArr})
// 	if e3 != nil {
// 		t.Errorf("Failed to produce s2 packet: %s", e3)
// 	}
// 	if len(s2) != 1 {
// 		t.Errorf("Wrong length for small packet: %d", len(s2))
// 	}

// 	s3, e4 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Arguments: smallArr})
// 	if e4 != nil {
// 		t.Errorf("Failed to produce s3 packet: %s", e4)
// 	}
// 	if len(s3) != 1 {
// 		t.Errorf("Wrong length for small packet: %d", len(s3))
// 	}

// 	m1, e5 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: medArr})
// 	if e5 != nil {
// 		t.Errorf("Failed to produce m1 packet: %s", e5)
// 	}
// 	if len(m1) != 1 {
// 		t.Errorf("Wrong length for medium packet: %d", len(m1))
// 	}

// 	m2, e6 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: medArr, Arguments: medArr})
// 	if e6 != nil {
// 		t.Errorf("Failed to produce m2 packet: %s", e6)
// 	}
// 	if len(m2) != 1 {
// 		t.Errorf("Wrong length for medium packet: %d", len(m2))
// 	}

// 	m3, e7 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Arguments: medArr})
// 	if e7 != nil {
// 		t.Errorf("Failed to produce m3 packet: %s", e7)
// 	}
// 	if len(m3) != 1 {
// 		t.Errorf("Wrong length for medium packet: %d", len(m3))
// 	}

// 	b1, e8 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: bigArr})
// 	if e8 != nil {
// 		t.Errorf("Failed to produce b1 packet: %s", e8)
// 	}
// 	if len(b1) != 4 {
// 		t.Errorf("Wrong length for big packet: %d", len(b1))
// 	}

// 	b2, e9 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Proxies: bigArr, Arguments: bigArr})
// 	if e9 != nil {
// 		t.Errorf("Failed to produce b2 packet: %s", e9)
// 	}
// 	if len(b2) != 4 {
// 		t.Errorf("Wrong length for small packet: %d", len(b2))
// 	}

// 	b3, e10 := dht.ProducePacket(&DHTPacket{Type: DHTPacketType_Ping, Arguments: bigArr})
// 	if e10 != nil {
// 		t.Errorf("Failed to produce b3 packet: %s", e10)
// 	}
// 	if len(b3) != 4 {
// 		t.Errorf("Wrong length for small packet: %d", len(b3))
// 	}
// }
