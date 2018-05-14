/*
Generated TestCommandSet
Generated TestDaemon_execRESTSet
Generated TestDaemon_SetLog
Generated TestDaemon_AddKey
*/
package main

import (
	"net"
	"net/http"
	"testing"
)

func TestCommandSet(t *testing.T) {
	type args struct {
		rpcPort int
		log     string
		hash    string
		keyfile string
		key     string
		ttl     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandSet(tt.args.rpcPort, tt.args.log, tt.args.hash, tt.args.keyfile, tt.args.key, tt.args.ttl)
		})
	}
}

func TestDaemon_execRESTSet(t *testing.T) {
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
			d.execRESTSet(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_SetLog(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *NameValueArg
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
			if err := d.SetLog(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.SetLog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDaemon_AddKey(t *testing.T) {
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
			p := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			if err := p.AddKey(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.AddKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
