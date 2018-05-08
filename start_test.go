/*
Generated TestCommandStart
Generated TestDaemon_execRESTStart
Generated TestDaemon_run
*/
package main

import (
	"net"
	"net/http"
	"testing"
)

func TestCommandStart(t *testing.T) {
	type args struct {
		restPort int
		ip       string
		hash     string
		mac      string
		dev      string
		keyfile  string
		key      string
		ttl      string
		fwd      bool
		port     int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandStart(tt.args.restPort, tt.args.ip, tt.args.hash, tt.args.mac, tt.args.dev, tt.args.keyfile, tt.args.key, tt.args.ttl, tt.args.fwd, tt.args.port)
		})
	}
}

func TestDaemon_execRESTStart(t *testing.T) {
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
			d.execRESTStart(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_run(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *RunArgs
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
			d := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			if err := d.run(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
