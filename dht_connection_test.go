/*
Generated TestDHTConnection_init
Generated TestDHTConnection_registerInstance
Generated TestDHTConnection_send
Generated TestDHTConnection_run
Generated TestDHTConnection_unregisterInstance
*/
package main

import (
	"sync"
	"testing"

	ptp "github.com/subutai-io/p2p/lib"
)

func TestDHTConnection_init(t *testing.T) {
	type fields struct {
		routers     []*DHTRouter
		routersList string
		lock        sync.Mutex
		instances   map[string]*P2PInstance
		registered  []string
		incoming    chan *ptp.DHTPacket
		ip          string
		isActive    bool
	}
	type args struct {
		routersSrc string
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
			dht := &DHTConnection{
				routers:     tt.fields.routers,
				routersList: tt.fields.routersList,
				lock:        tt.fields.lock,
				instances:   tt.fields.instances,
				registered:  tt.fields.registered,
				incoming:    tt.fields.incoming,
				ip:          tt.fields.ip,
				isActive:    tt.fields.isActive,
			}
			if err := dht.init(tt.args.routersSrc); (err != nil) != tt.wantErr {
				t.Errorf("DHTConnection.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTConnection_registerInstance(t *testing.T) {
	type fields struct {
		routers     []*DHTRouter
		routersList string
		lock        sync.Mutex
		instances   map[string]*P2PInstance
		registered  []string
		incoming    chan *ptp.DHTPacket
		ip          string
		isActive    bool
	}
	type args struct {
		hash string
		inst *P2PInstance
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
			dht := &DHTConnection{
				routers:     tt.fields.routers,
				routersList: tt.fields.routersList,
				lock:        tt.fields.lock,
				instances:   tt.fields.instances,
				registered:  tt.fields.registered,
				incoming:    tt.fields.incoming,
				ip:          tt.fields.ip,
				isActive:    tt.fields.isActive,
			}
			if err := dht.registerInstance(tt.args.hash, tt.args.inst); (err != nil) != tt.wantErr {
				t.Errorf("DHTConnection.registerInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDHTConnection_send(t *testing.T) {
	type fields struct {
		routers     []*DHTRouter
		routersList string
		lock        sync.Mutex
		instances   map[string]*P2PInstance
		registered  []string
		incoming    chan *ptp.DHTPacket
		ip          string
		isActive    bool
	}
	type args struct {
		packet *ptp.DHTPacket
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
			dht := &DHTConnection{
				routers:     tt.fields.routers,
				routersList: tt.fields.routersList,
				lock:        tt.fields.lock,
				instances:   tt.fields.instances,
				registered:  tt.fields.registered,
				incoming:    tt.fields.incoming,
				ip:          tt.fields.ip,
				isActive:    tt.fields.isActive,
			}
			dht.send(tt.args.packet)
		})
	}
}

func TestDHTConnection_run(t *testing.T) {
	type fields struct {
		routers     []*DHTRouter
		routersList string
		lock        sync.Mutex
		instances   map[string]*P2PInstance
		registered  []string
		incoming    chan *ptp.DHTPacket
		ip          string
		isActive    bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTConnection{
				routers:     tt.fields.routers,
				routersList: tt.fields.routersList,
				lock:        tt.fields.lock,
				instances:   tt.fields.instances,
				registered:  tt.fields.registered,
				incoming:    tt.fields.incoming,
				ip:          tt.fields.ip,
				isActive:    tt.fields.isActive,
			}
			dht.run()
		})
	}
}

func TestDHTConnection_unregisterInstance(t *testing.T) {
	type fields struct {
		routers     []*DHTRouter
		routersList string
		lock        sync.Mutex
		instances   map[string]*P2PInstance
		registered  []string
		incoming    chan *ptp.DHTPacket
		ip          string
		isActive    bool
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
			dht := &DHTConnection{
				routers:     tt.fields.routers,
				routersList: tt.fields.routersList,
				lock:        tt.fields.lock,
				instances:   tt.fields.instances,
				registered:  tt.fields.registered,
				incoming:    tt.fields.incoming,
				ip:          tt.fields.ip,
				isActive:    tt.fields.isActive,
			}
			if err := dht.unregisterInstance(tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("DHTConnection.unregisterInstance() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
