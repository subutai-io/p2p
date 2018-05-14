/*
Generated TestCommandShow
Generated TestDaemon_execRESTShow
Generated TestDaemon_Show
Generated TestDaemon_showOutput
Generated TestDaemon_showIP
Generated TestDaemon_showHash
Generated TestDaemon_showInterfaces
Generated TestDaemon_showAllInterfaces
Generated TestDaemon_showBindInterfaces
Generated TestDaemon_showInstances
*/
package main

import (
	"net"
	"net/http"
	"reflect"
	"testing"
)

func TestCommandShow(t *testing.T) {
	type args struct {
		queryPort  int
		hash       string
		ip         string
		interfaces bool
		all        bool
		bind       bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CommandShow(tt.args.queryPort, tt.args.hash, tt.args.ip, tt.args.interfaces, tt.args.all, tt.args.bind)
		})
	}
}

func TestDaemon_execRESTShow(t *testing.T) {
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
			d.execRESTShow(tt.args.w, tt.args.r)
		})
	}
}

func TestDaemon_Show(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *request
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
			got, err := d.Show(tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.Show() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.Show() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showOutput(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		data []ShowOutput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
			got, err := d.showOutput(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showOutput() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showIP(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		ip       string
		instance *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
			got, err := d.showIP(tt.args.ip, tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showHash(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		instance *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
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
			got, err := d.showHash(tt.args.instance)
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showHash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showInterfaces(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
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
			got, err := d.showInterfaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showInterfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showInterfaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showAllInterfaces(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
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
			got, err := d.showAllInterfaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showAllInterfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showAllInterfaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showBindInterfaces(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
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
			got, err := d.showBindInterfaces()
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showBindInterfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showBindInterfaces() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_showInstances(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
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
			got, err := d.showInstances()
			if (err != nil) != tt.wantErr {
				t.Errorf("Daemon.showInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Daemon.showInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}
