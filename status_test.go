/*
Generated TestCommandStatus
Generated TestDaemon_execRESTStatus
Generated TestDaemon_Status
*/
package main

import (
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestCommandStatus(t *testing.T) {
	type args struct {
		restPort int
		hash     string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandStatus(tt.args.restPort, tt.args.hash)
		})
	}
}

func TestDaemon_execRESTStatus(t *testing.T) {
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
			d.execRESTStatus(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_Status(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		hash string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *statusResponse
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
			got, err := d.Status(tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
