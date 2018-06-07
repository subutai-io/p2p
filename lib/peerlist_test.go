package ptp

import (
	"reflect"
	"sync"
	"testing"
)

func TestUpdateTables(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.updateTables("500", "192.168.1.1", "01:02:03:04:05:06")
	_, exists := l.tableIPID["500"]
	_, exist := l.tableMacID["01:02:03:04:05:06"]
	if !exists && !exist {
		t.Error("Error. Can't update peer")
	}
}

func TestDeleteTables(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.tableIPID["800"] = "192.168.8.8"
	l.tableMacID["800"] = "01:02:03:04:05:06"
	l.deleteTables("800", "800")
	_, exists := l.tableIPID["800"]
	_, exist := l.tableMacID["800"]
	if exist && exists {
		t.Error("Error")
	}
}

func TestGet(t *testing.T) {
	l := new(PeerList)
	np1 := new(NetworkPeer)
	np2 := new(NetworkPeer)
	l.Init()
	l.peers["444"] = np1
	l.peers["445"] = np2
	get := l.Get()
	var wait map[string]*NetworkPeer
	wait = make(map[string]*NetworkPeer)
	wait["444"] = np1
	wait["445"] = np2
	if !reflect.DeepEqual(get, wait) {
		t.Error("wait, get", wait, get)
	}
}

func TestLength(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.peers["77"] = new(NetworkPeer)
	l.peers["78"] = new(NetworkPeer)
	count := 0
	for i := 0; i < len(l.peers); i++ {
		count++
	}
	get := l.Length()
	if get != count {
		t.Errorf("Error. Wait: %v, get: %v", count, get)
	}
}

func TestGetPeer(t *testing.T) {
	l := new(PeerList)
	l.Init()
	l.peers["9"] = new(NetworkPeer)
	l.peers["99"] = new(NetworkPeer)

	get1 := l.GetPeer("9")
	if get1 != l.peers["9"] {
		t.Error("Error")
	}
	get2 := l.GetPeer("-1")
	if get2 != nil {
		t.Error("Error")
	}
}

func TestPeerList_GetID(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		ip string
	}

	data := make(map[string]string)
	data["127.0.0.1"] = "test_id"

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"t1", fields{tableIPID: data}, args{"127.0.0.1"}, "test_id", false},
		{"t1", fields{tableIPID: data}, args{"127.0.0.2"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			got, err := l.GetID(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerList.GetID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeerList.GetID() = %v, want %v", got, tt.want)
			}
		})
	}
}
