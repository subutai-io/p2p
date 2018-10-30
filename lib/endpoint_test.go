package ptp

import (
	"encoding/binary"
	"net"
	"reflect"
	"testing"
	"time"
)

func TestEndpoint_addrToBytes(t *testing.T) {
	type fields struct {
		Addr              *net.UDPAddr
		LastContact       time.Time
		LastPing          time.Time
		broken            bool
		Latency           time.Duration
		LastLatencyQuery  time.Time
		MeasureInProgress bool
	}

	a1, _ := net.ResolveUDPAddr("udp4", "127.0.0.1:1111")
	r1 := []byte{127, 0, 0, 1, 0, 0}
	binary.BigEndian.PutUint16(r1[4:6], uint16(1111))

	a2, _ := net.ResolveUDPAddr("udp4", "0.0.0.0:0000")
	r2 := []byte{0, 0, 0, 0, 0, 0}
	binary.BigEndian.PutUint16(r2[4:6], uint16(0))

	a3, _ := net.ResolveUDPAddr("udp4", "255.255.255.255:65535")
	r3 := []byte{255, 255, 255, 255, 255, 255}

	a4, _ := net.ResolveUDPAddr("udp4", "254.254.254.254:65534")
	r4 := []byte{254, 254, 254, 254, 0, 0}
	binary.BigEndian.PutUint16(r4[4:6], uint16(65534))

	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{"Testing 127.0.0.1:1111", fields{Addr: a1}, r1},
		{"Testing 0.0.0.0:0000", fields{Addr: a2}, r2},
		{"Testing 255.255.255.255:65535", fields{Addr: a3}, r3},
		{"Testing 254.254.254.254:65534", fields{Addr: a4}, r4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Endpoint{
				Addr:              tt.fields.Addr,
				LastContact:       tt.fields.LastContact,
				LastPing:          tt.fields.LastPing,
				broken:            tt.fields.broken,
				Latency:           tt.fields.Latency,
				LastLatencyQuery:  tt.fields.LastLatencyQuery,
				MeasureInProgress: tt.fields.MeasureInProgress,
			}
			if got := e.addrToBytes(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Endpoint.addrToBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndpoint_Measure(t *testing.T) {
	type fields struct {
		Addr             *net.UDPAddr
		LastContact      time.Time
		LastPing         time.Time
		broken           bool
		Latency          time.Duration
		LastLatencyQuery time.Time
	}
	type args struct {
		n  *Network
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Endpoint{
				Addr:             tt.fields.Addr,
				LastContact:      tt.fields.LastContact,
				LastPing:         tt.fields.LastPing,
				broken:           tt.fields.broken,
				Latency:          tt.fields.Latency,
				LastLatencyQuery: tt.fields.LastLatencyQuery,
			}
			e.Measure(tt.args.n, tt.args.id)
		})
	}
}
