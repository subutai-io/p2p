package ptp

import (
	"testing"
)

func TestExtract(t *testing.T) {
	var m string = "d1:a0:1:c4:ping1:i36:00000000-1111-2222-3333-4444444444441:p0:1:q1:0e"
	var dht DHTClient
	result, err := dht.Extract([]byte(m))
	if result.Id != "00000000-1111-2222-3333-444444444444" {
		t.Errorf("Failed to extract DHT message")
	}
	if err != nil {
		t.Errorf("Error during DHT message extraction")
	}
}
