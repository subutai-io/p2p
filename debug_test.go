/*
Generated TestCommandDebug
Generated TestDaemon_execRESTDebug
Generated TestDaemon_Debug
*/
package main

import (
	"net"
	"net/http"
	"testing"
)

func TestCommandDebug(t *testing.T) {
	type args struct {
		restPort int
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandDebug(tt.args.restPort)
		})
	}
}

func TestDaemon_execRESTDebug(t *testing.T) {
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
			d.execRESTDebug(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_Debug(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *Args
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
			if err := p.Debug(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.Debug() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
