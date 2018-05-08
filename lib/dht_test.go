package ptp

/*
import (
 	"testing"
)

func TestCompose(t *testing.T) {
 	var dht DHTClient
 	t1 := dht.Compose("ping", "00000000-1111-2222-3333-444444444444", "QUERY", "ARGUMENT")
 	if t1 != "d1:a8:ARGUMENT1:c4:ping1:i36:00000000-1111-2222-3333-4444444444441:p0:1:q5:QUERYe" {
 		t.Fatalf("dht.Compose failed (1)")
 	}
 	t2 := dht.Compose("", "", "", "")
 	if t2 != "" {
 		t.Fatalf("dht.Compose failed (2)")
 	}
}

func TestExtract(t *testing.T) {
 	var m = "d1:a0:1:c4:ping1:i36:00000000-1111-2222-3333-4444444444441:p0:1:q1:0e"
 	var dht DHTClient
 	result, err := dht.Extract([]byte(m))
 	if result.ID != "00000000-1111-2222-3333-444444444444" {
 		t.Fatalf("Failed to extract DHT message")
 	}
 	if err != nil {
 		t.Fatalf("Error during DHT message extraction")
 	}
}

func TestEncodeRequest(t *testing.T) {
	var dht DHTClient
 	t1 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test1", Query: "Query1", Arguments: "Argument1"})
 	if t1 != "d1:a9:Argument11:c5:Test11:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query1e" {
 		t.Fatalf("EncodeRequest failed (1)")
 	}
 	t2 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test2", Query: "Query2", Arguments: "Argument2"})
 	if t2 != "d1:a9:Argument21:c5:Test21:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query2e" {
 		t.Fatalf("EncodeRequest failed (2)")
 	}
 	t3 := dht.EncodeRequest(DHTMessage{ID: "00000000-1111-2222-3333-444444444444", Command: "Test3", Query: "Query3", Arguments: "Argument3"})
 	if t3 != "d1:a9:Argument31:c5:Test31:i36:00000000-1111-2222-3333-4444444444441:p0:1:q6:Query3e" {
 		t.Fatalf("EncodeRequest failed (3)")
 	}
}
*/


import (
	"testing"
	"net"
	"github.com/golang/protobuf/proto"
	"time"
	"reflect"
)


func TestInit(t *testing.T) {
	dht := new(DHTClient)
	err := dht.Init("hash")
	if err != nil {
		t.Fatalf("Failed to init (1): %v", err)
	}
	err = dht.Init("hash")
	if err != nil {
		t.Fatalf("Failed to init (2): %v. Expected %v, got %v", err, "dht.cdn.subut.ai:6881", dht.Routers)
	}
}

func TestConnect(t *testing.T) {
	finish := make(chan bool)
	defer close(finish)
	dht := new(DHTClient)
	ActiveInterfaces = []net.IP{net.IP("127.0.0.1")}
	go func() {
		errChan := make(chan error)
		go func() {
			errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080,}}})
		}()
		err := <-errChan
		if err == nil {
			t.Fatalf("Failed to connect (1): must have returned non-nil but returned nil")
		}
		finish <- true
	}()
	breakFirstFor:
	for {
		select {
		case <-finish:
			go func() {
				dht.OutgoingData = make(chan *DHTPacket, 1)
				defer close(dht.OutgoingData)
				errChan := make(chan error)
				go func() {
					errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080,}}})
				}()
				time.Sleep(2 * time.Second)
				dht.Connected = true
				err := <-errChan
				if err != nil {
					t.Fatalf("Failed to connect (2): %v", err)
				}
				finish <- true
			}()
			break breakFirstFor
		}
	}
	breaKSecondFor:
	for {
		select {
		case <-finish:
			go func() {
				dht.OutgoingData = make(chan *DHTPacket, 1)
				defer close(dht.OutgoingData)
				errChan := make(chan error)
				go func() {
					errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080,}}})
				}()
				err := <-errChan
				if err == nil {
					t.Fatalf("Failed to connect (3): must have returned non-nil but returned nil")
				}
				finish <- true
			}()
			break breaKSecondFor
		}
	}
	breaKThirdFor:
	for {
		select {
		case <-finish:
			break breaKThirdFor
		}
	}
}

func TestRead(t *testing.T) {
	dht := new(DHTClient)
	dht.IncomingData = nil
	_, err := dht.read()
	if err == nil {
		t.Fatalf("Failed to read (1): must have returned non-nil but returned nil")
	}
	dht.IncomingData = make(chan *DHTPacket)
	go func() {
		dht.IncomingData <- new(DHTPacket)
	}()
	packet, err := dht.read()
	close(dht.IncomingData)
	if err != nil {
		t.Fatalf("Failed to read (2): %v", err)
	}
	if packet == nil {
		t.Fatalf("Failed to read (3): must have returned non-nil packet but returned nil packet")
	}
	packet = nil
	dht.IncomingData = make(chan *DHTPacket)
	go func() {
		dht.IncomingData <- packet
	}()
	packet, err = dht.read()
	close(dht.IncomingData)
	if err == nil {
		t.Fatalf("Failed to read (4): must have returned non-nil but returned nil")
	}
	if packet != nil {
		t.Fatalf("Failed to read (5): must have returned nil packet but returned non-nil packet")
	}
}

func TestSend(t *testing.T) {
	{
		dht := new(DHTClient)
		dht.IncomingData = make(chan *DHTPacket)
		dht.OutgoingData = make(chan *DHTPacket)
		dht.Close()
		err := dht.send(&DHTPacket{})
		if err == nil {
			t.Fatalf("Failed to send (1): must have returned non-nil but returned nil: %v", err)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type: DHTPacketType_Connect,
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (2): Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (3): failed to marshal data (1): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (4): failed to unmarshal data (1): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (5): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (6): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (7): Proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (8): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (9): failed to marshal data (2): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (10): failed to unmarshal data (2): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (11): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (12): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (13): proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:    DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (14): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (15): failed to marshal data (3): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (16): failed to unmarshal data (3): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (17): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (18): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (19): proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies:   []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (20): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (21): failed to marshal data (4): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (22): failed to unmarshal data (4): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (23): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (24): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (25): Proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
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
			p1.Arguments = append(p1.Arguments, "Argument")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (26): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (27): failed to marshal data (5): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (28): failed to unmarshal data (5): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (29): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (30): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (31): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:    DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Proxy")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (32): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (33): failed to marshal data (6): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (34): failed to unmarshal data (6): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (35): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (36): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (37): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies:   []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Arguments = append(p1.Arguments, "Argument")
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Proxy")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (38): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000 + 1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (39): failed to marshal data (7): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (40): failed to unmarshal data (7): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (41): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (42): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (43): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
}

func TestSendFind(t *testing.T) {
	dht := new(DHTClient)
	dht.NetworkHash = ""
	err := dht.sendFind()
	if err == nil {
		t.Fatalf("Failed to sendFind (1): must have returned non-nil but returned nil")
	}
	dht.NetworkHash = "NetworkHash"
	err = dht.sendFind()
	if err == nil {
		t.Fatalf("Failed to sendFind (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendFind()
	if err != nil {
		t.Fatalf("Failed to sendFind (3): %v", err)
	}
}

func TestSendNode(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendNode("ID", []net.IP{net.IP("192.168.0.1")})
	if err == nil {
		t.Fatalf("Failed to sendNode (1): must have returned non-nil but returned nil")
	}
	err = dht.sendNode("123456789012345678901234567890123456", []net.IP{net.IP("192.168.0.1"), net.IP(nil), net.IP("192.168.0.1")})
	if err == nil {
		t.Fatalf("Failed to sendNode (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendNode("123456789012345678901234567890123456", []net.IP{net.IP("192.168.0.1"), net.IP(nil), net.IP("192.168.0.1")})
	if err != nil {
		t.Fatalf("Failed to sendNode (3): %v", err)
	}
}

func TestSendState(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendState("ID", "1")
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendState("123456789012345678901234567890123456", "1")
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendState("123456789012345678901234567890123456", "1")
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}

func TestSendDHCP(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendDHCP(nil, nil)
	if err == nil {
		t.Fatalf("Failed to sendDHCP (1): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	err = dht.sendDHCP(nil, nil)
	close(dht.OutgoingData)
	if err != nil {
		t.Fatalf("Failed to sendDHCP (2): %v", err)
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendDHCP(nil, new(net.IPNet))
	if err != nil {
		t.Fatalf("Failed to sendDHCP (3): %v", err)
	}
}

func TestSendProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendProxy()
	if err == nil {
		t.Fatalf("Failed to sendProxy (1): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendProxy()
	if err != nil {
		t.Fatalf("Failed to sendDHCP (2): %v", err)
	}
}

func TestSendRequestProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendRequestProxy("ID")
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendRequestProxy("123456789012345678901234567890123456")
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendRequestProxy("123456789012345678901234567890123456")
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}

func TestSendReportProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendReportProxy([]*net.UDPAddr{})
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendReportProxy([]*net.UDPAddr{{IP: net.IP("127.0.0.1"), Port: 8080,}})
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendReportProxy([]*net.UDPAddr{{IP: net.IP("127.0.0.1"), Port: 8080,}})
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}
/*
Generated TestDHTClient_Init
Generated TestDHTClient_Connect
Generated TestDHTClient_read
Generated TestDHTClient_send
Generated TestDHTClient_sendFind
Generated TestDHTClient_sendNode
Generated TestDHTClient_sendState
Generated TestDHTClient_sendDHCP
Generated TestDHTClient_sendProxy
Generated TestDHTClient_sendRequestProxy
Generated TestDHTClient_sendReportProxy
Generated TestDHTClient_Close
Generated TestDHTClient_WaitID
Generated TestDHTClient_RegisterProxy
Generated TestDHTClient_ReportLoad
package ptp

import (
	"net"
	"reflect"
	"testing"
	"time"

	"github.com/golang/protobuf/proto"
)
*/

/*
func TestInit(t *testing.T) {
	dht := new(DHTClient)
	err := dht.Init("hash")
	if err != nil {
		t.Fatalf("Failed to init (1): %v", err)
	}
	err = dht.Init("hash")
	if err != nil {
		t.Fatalf("Failed to init (2): %v. Expected %v, got %v", err, "dht.cdn.subut.ai:6881", dht.Routers)
	}
}

func TestConnect(t *testing.T) {
	finish := make(chan bool)
	defer close(finish)
	dht := new(DHTClient)
	ActiveInterfaces = []net.IP{net.IP("127.0.0.1")}
	go func() {
		errChan := make(chan error)
		go func() {
			errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080}}})
		}()
		err := <-errChan
		if err == nil {
			t.Fatalf("Failed to connect (1): must have returned non-nil but returned nil")
		}
		finish <- true
	}()
breakFirstFor:
	for {
		select {
		case <-finish:
			go func() {
				dht.OutgoingData = make(chan *DHTPacket, 1)
				defer close(dht.OutgoingData)
				errChan := make(chan error)
				go func() {
					errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080}}})
				}()
				time.Sleep(2 * time.Second)
				dht.Connected = true
				err := <-errChan
				if err != nil {
					t.Fatalf("Failed to connect (2): %v", err)
				}
				finish <- true
			}()
			break breakFirstFor
		}
	}
breaKSecondFor:
	for {
		select {
		case <-finish:
			go func() {
				dht.OutgoingData = make(chan *DHTPacket, 1)
				defer close(dht.OutgoingData)
				errChan := make(chan error)
				go func() {
					errChan <- dht.Connect([]net.IP{net.IP("127.0.0.1"), net.IP(nil), net.IP("127.0.0.2")}, []*proxyServer{{Endpoint: &net.UDPAddr{IP: net.IP("192.168.0.1"), Port: 8080}}})
				}()
				err := <-errChan
				if err == nil {
					t.Fatalf("Failed to connect (3): must have returned non-nil but returned nil")
				}
				finish <- true
			}()
			break breaKSecondFor
		}
	}
breaKThirdFor:
	for {
		select {
		case <-finish:
			break breaKThirdFor
		}
	}
}

func TestRead(t *testing.T) {
	dht := new(DHTClient)
	dht.IncomingData = nil
	_, err := dht.read()
	if err == nil {
		t.Fatalf("Failed to read (1): must have returned non-nil but returned nil")
	}
	dht.IncomingData = make(chan *DHTPacket)
	go func() {
		dht.IncomingData <- new(DHTPacket)
	}()
	packet, err := dht.read()
	close(dht.IncomingData)
	if err != nil {
		t.Fatalf("Failed to read (2): %v", err)
	}
	if packet == nil {
		t.Fatalf("Failed to read (3): must have returned non-nil packet but returned nil packet")
	}
	packet = nil
	dht.IncomingData = make(chan *DHTPacket)
	go func() {
		dht.IncomingData <- packet
	}()
	packet, err = dht.read()
	close(dht.IncomingData)
	if err == nil {
		t.Fatalf("Failed to read (4): must have returned non-nil but returned nil")
	}
	if packet != nil {
		t.Fatalf("Failed to read (5): must have returned nil packet but returned non-nil packet")
	}
}

func TestSend(t *testing.T) {
	{
		dht := new(DHTClient)
		dht.IncomingData = make(chan *DHTPacket)
		dht.OutgoingData = make(chan *DHTPacket)
		dht.Close()
		err := dht.send(&DHTPacket{})
		if err == nil {
			t.Fatalf("Failed to send (1): must have returned non-nil but returned nil: %v", err)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type: DHTPacketType_Connect,
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (2): Could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (3): failed to marshal data (1): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (4): failed to unmarshal data (1): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (5): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (6): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (7): Proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (8): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (9): failed to marshal data (2): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (10): failed to unmarshal data (2): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (11): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (12): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (13): proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:    DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (14): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (15): failed to marshal data (3): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (16): failed to unmarshal data (3): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (17): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (18): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (19): proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies:   []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (20): could not send packet")
			}
		}()
		data := <-dht.OutgoingData
		packetBytes, err := proto.Marshal(data)
		if err != nil {
			t.Fatalf("Failed to send (21): failed to marshal data (4): %v", err)
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(packetBytes, packet)
		if err != nil {
			t.Fatalf("Failed to send (22): failed to unmarshal data (4): %v", err)
		}
		close(dht.OutgoingData)
		if packet.Type != p1.Type {
			t.Fatalf("Failed to send (23): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
		}
		if len(packet.Arguments) != lenArguments {
			t.Fatalf("Failed to send (24): arguments length mismatch: %d -> %d", len(packet.Arguments), lenArguments)
		}
		if len(packet.Proxies) != lenProxies {
			t.Fatalf("Failed to send (25): Proxies length mismatch: %d -> %d", len(packet.Proxies), lenProxies)
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
			p1.Arguments = append(p1.Arguments, "Argument")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (26): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000+1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (27): failed to marshal data (5): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (28): failed to unmarshal data (5): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (29): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (30): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (31): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:    DHTPacketType_Connect,
			Proxies: []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Proxy")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (32): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000+1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (33): failed to marshal data (6): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (34): failed to unmarshal data (6): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (35): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (36): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (37): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
	{
		dht := new(DHTClient)
		dht.OutgoingData = make(chan *DHTPacket)
		p1 := &DHTPacket{
			Type:      DHTPacketType_Connect,
			Arguments: []string{"ARGUMENT_1", "ARGUMENT_2", "ARGUMENT_3", "ARGUMENT_4", "ARGUMENT_5", "ARGUMENT_6"},
			Proxies:   []string{"PROXY_1", "PROXY_2", "PROXY_3", "PROXY_4", "PROXY_5", "PROXY_6"},
		}
		for i := 0; i < 100000; i++ {
			p1.Arguments = append(p1.Arguments, "Argument")
		}
		for i := 0; i < 100000; i++ {
			p1.Proxies = append(p1.Proxies, "Proxy")
		}
		lenArguments := len(p1.Arguments)
		lenProxies := len(p1.Proxies)
		go func() {
			err := dht.send(p1)
			if err != nil {
				t.Fatalf("Failed to send (38): could not send packet")
			}
		}()
		data := []*DHTPacket{}
		for i := 0; i < 10000+1; i++ {
			item := <-dht.OutgoingData
			packetBytes, err := proto.Marshal(item)
			if err != nil {
				t.Fatalf("Failed to send (39): failed to marshal data (7): %v", err)
			}
			packet := &DHTPacket{}
			err = proto.Unmarshal(packetBytes, packet)
			if err != nil {
				t.Fatalf("Failed to send (40): failed to unmarshal data (7): %v", err)
			}
			data = append(data, packet)
		}
		close(dht.OutgoingData)
		allArguments := []string{}
		allProxies := []string{}
		for _, packet := range data {
			if packet.Type != p1.Type {
				t.Fatalf("Failed to send (41): data mismatch on type: %d -> %d", int(packet.Type), int(p1.Type))
			}
			allArguments = append(allArguments, packet.Arguments[:]...)
			allProxies = append(allProxies, packet.Proxies[:]...)
		}
		if len(allArguments) != lenArguments {
			t.Fatalf("Failed to send (42): arguments length mismatch: %d -> %d", len(allArguments), lenArguments)
		}
		if len(allProxies) != lenProxies {
			t.Fatalf("Failed to send (43): proxies length mismatch: %d -> %d", len(allProxies), lenProxies)
		}
	}
}

func TestSendFind(t *testing.T) {
	dht := new(DHTClient)
	dht.NetworkHash = ""
	err := dht.sendFind()
	if err == nil {
		t.Fatalf("Failed to sendFind (1): must have returned non-nil but returned nil")
	}
	dht.NetworkHash = "NetworkHash"
	err = dht.sendFind()
	if err == nil {
		t.Fatalf("Failed to sendFind (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendFind()
	if err != nil {
		t.Fatalf("Failed to sendFind (3): %v", err)
	}
}

func TestSendNode(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendNode("ID", []net.IP{net.IP("192.168.0.1")})
	if err == nil {
		t.Fatalf("Failed to sendNode (1): must have returned non-nil but returned nil")
	}
	err = dht.sendNode("123456789012345678901234567890123456", []net.IP{net.IP("192.168.0.1"), net.IP(nil), net.IP("192.168.0.1")})
	if err == nil {
		t.Fatalf("Failed to sendNode (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendNode("123456789012345678901234567890123456", []net.IP{net.IP("192.168.0.1"), net.IP(nil), net.IP("192.168.0.1")})
	if err != nil {
		t.Fatalf("Failed to sendNode (3): %v", err)
	}
}

func TestSendState(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendState("ID", "1")
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendState("123456789012345678901234567890123456", "1")
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendState("123456789012345678901234567890123456", "1")
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}

func TestSendDHCP(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendDHCP(nil, nil)
	if err == nil {
		t.Fatalf("Failed to sendDHCP (1): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	err = dht.sendDHCP(nil, nil)
	close(dht.OutgoingData)
	if err != nil {
		t.Fatalf("Failed to sendDHCP (2): %v", err)
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendDHCP(nil, new(net.IPNet))
	if err != nil {
		t.Fatalf("Failed to sendDHCP (3): %v", err)
	}
}

func TestSendProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendProxy()
	if err == nil {
		t.Fatalf("Failed to sendProxy (1): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendProxy()
	if err != nil {
		t.Fatalf("Failed to sendDHCP (2): %v", err)
	}
}

func TestSendRequestProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendRequestProxy("ID")
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendRequestProxy("123456789012345678901234567890123456")
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendRequestProxy("123456789012345678901234567890123456")
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}

func TestSendReportProxy(t *testing.T) {
	dht := new(DHTClient)
	err := dht.sendReportProxy([]*net.UDPAddr{})
	if err == nil {
		t.Fatalf("Failed to sendState (1): must have returned non-nil but returned nil")
	}
	err = dht.sendReportProxy([]*net.UDPAddr{{IP: net.IP("127.0.0.1"), Port: 8080}})
	if err == nil {
		t.Fatalf("Failed to sendState (2): must have returned non-nil but returned nil")
	}
	dht.OutgoingData = make(chan *DHTPacket, 1)
	defer close(dht.OutgoingData)
	err = dht.sendReportProxy([]*net.UDPAddr{{IP: net.IP("127.0.0.1"), Port: 8080}})
	if err != nil {
		t.Fatalf("Failed to sendState (3): %v", err)
	}
}
*/

func TestDHTClient_Init(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.Init(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_Connect(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		ipList    []net.IP
		proxyList []*proxyServer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.Connect(tt.args.ipList, tt.args.proxyList); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.Connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_read(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		want    *DHTPacket
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			got, err := dht.read()
			if (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.read() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DHTClient.read() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDHTClient_send(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		packet *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.send(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendFind(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendFind(); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendFind() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendNode(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		id     string
		ipList []net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendNode(tt.args.id, tt.args.ipList); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendNode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendState(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		id    string
		state string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendState(tt.args.id, tt.args.state); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendState() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendDHCP(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		ip      net.IP
		network *net.IPNet
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendDHCP(tt.args.ip, tt.args.network); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendDHCP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendProxy(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendProxy(); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendRequestProxy(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendRequestProxy(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendRequestProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_sendReportProxy(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		addr []*net.UDPAddr
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.sendReportProxy(tt.args.addr); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.sendReportProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_Close(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.Close(); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_WaitID(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.WaitID(); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.WaitID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_RegisterProxy(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		ip   net.IP
		port int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.RegisterProxy(tt.args.ip, tt.args.port); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.RegisterProxy() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTClient_ReportLoad(t *testing.T) {
	type fields struct {
		Routers           string
		NetworkHash       string
		ID                string
		FailedRouters     []string
		Connections       []*net.TCPConn
		LocalPort         int
		RemotePort        int
		Forwarders        []Forwarder
		TCPCallbacks      map[DHTPacketType]dhtCallback
		Mode              OperatingMode
		IPList            []net.IP
		IP                net.IP
		Network           *net.IPNet
		Connected         bool
		LastUpdate        time.Time
		OutboundIP        net.IP
		ListenerIsRunning bool
		IncomingData      chan *DHTPacket
		OutgoingData      chan *DHTPacket
	}
	type args struct {
		clientsNum int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTClient{
				Routers:           tt.fields.Routers,
				NetworkHash:       tt.fields.NetworkHash,
				ID:                tt.fields.ID,
				FailedRouters:     tt.fields.FailedRouters,
				Connections:       tt.fields.Connections,
				LocalPort:         tt.fields.LocalPort,
				RemotePort:        tt.fields.RemotePort,
				Forwarders:        tt.fields.Forwarders,
				TCPCallbacks:      tt.fields.TCPCallbacks,
				Mode:              tt.fields.Mode,
				IPList:            tt.fields.IPList,
				IP:                tt.fields.IP,
				Network:           tt.fields.Network,
				Connected:         tt.fields.Connected,
				LastUpdate:        tt.fields.LastUpdate,
				OutboundIP:        tt.fields.OutboundIP,
				ListenerIsRunning: tt.fields.ListenerIsRunning,
				IncomingData:      tt.fields.IncomingData,
				OutgoingData:      tt.fields.OutgoingData,
			}
			if err := dht.ReportLoad(tt.args.clientsNum); (err != nil) != tt.wantErr {
				t.Errorf("DHTClient.ReportLoad() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
