package ptp

import (
	"net"
	"testing"
	"time"
)

func Test_proxyServer_Measure(t *testing.T) {
	type fields struct {
		Addr              *net.UDPAddr
		Endpoint          *net.UDPAddr
		Status            proxyStatus
		LastUpdate        time.Time
		Created           time.Time
		Latency           time.Duration
		LastLatencyQuery  time.Time
		MeasureInProgress bool
	}
	type args struct {
		n *Network
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"t1", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &proxyServer{
				Addr:              tt.fields.Addr,
				Endpoint:          tt.fields.Endpoint,
				Status:            tt.fields.Status,
				LastUpdate:        tt.fields.LastUpdate,
				Created:           tt.fields.Created,
				Latency:           tt.fields.Latency,
				LastLatencyQuery:  tt.fields.LastLatencyQuery,
				MeasureInProgress: tt.fields.MeasureInProgress,
			}
			p.Measure(tt.args.n)
		})
	}
}
