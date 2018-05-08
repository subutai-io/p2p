/*
Generated TestCommandStop
Generated TestDaemon_execRESTStop
Generated TestDaemon_Stop
*/
package main

import (
	"net"
	"net/http"
	"testing"
)

func TestCommandStop(t *testing.T) {
	type args struct {
		rpcPort int
		hash    string
		dev     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandStop(tt.args.rpcPort, tt.args.hash, tt.args.dev)
		})
	}
}

func TestDaemon_execRESTStop(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
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
			d := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			d.execRESTStop(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_Stop(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *DaemonArgs
		resp *Response
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
			p := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			if err := p.Stop(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
