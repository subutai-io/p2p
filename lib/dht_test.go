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

func TestSend(t *testing.T) {
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		close(dht.OutgoingData)
		if data.Type != p1.Type {
			t.Fatalf("Data mismatch on type: %d -> %d", int(data.Type), int(p1.Type))
		}
		if len(data.Arguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(data.Arguments), len(p1.Arguments))
		}
		if len(data.Proxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(data.Proxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		close(dht.OutgoingData)
		if data.Type != p1.Type {
			t.Fatalf("Data mismatch on type: %d -> %d", int(data.Type), int(p1.Type))
		}
		if len(data.Arguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(data.Arguments), len(p1.Arguments))
		}
		if len(data.Proxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(data.Proxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		close(dht.OutgoingData)
		if data.Type != p1.Type {
			t.Fatalf("Data mismatch on type: %d -> %d", int(data.Type), int(p1.Type))
		}
		if len(data.Arguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(data.Arguments), len(p1.Arguments))
		}
		if len(data.Proxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(data.Proxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		close(dht.OutgoingData)
		if data.Type != p1.Type {
			t.Fatalf("Data mismatch on type: %d -> %d", int(data.Type), int(p1.Type))
		}
		if len(data.Arguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(data.Arguments), len(p1.Arguments))
		}
		if len(data.Proxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(data.Proxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Arguments = append(p1.Arguments, "Azret Argument")
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			data = append(data, <-dht.OutgoingData)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(allArguments), len(p1.Arguments))
		}
		if len(allProxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(allProxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Azret Proxy")
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			data = append(data, <-dht.OutgoingData)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(allArguments), len(p1.Arguments))
		}
		if len(allProxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(allProxies), len(p1.Proxies))
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Arguments = append(p1.Arguments, "Azret Argument")
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Azret Proxy")
		}
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			data = append(data, <-dht.OutgoingData)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != len(p1.Arguments) {
			t.Fatalf("Arguments length mismatch: %d -> %d", len(allArguments), len(p1.Arguments))
		}
		if len(allProxies) != len(p1.Proxies) {
			t.Fatalf("Proxies length mismatch: %d -> %d", len(allProxies), len(p1.Proxies))
		}
	}
}
