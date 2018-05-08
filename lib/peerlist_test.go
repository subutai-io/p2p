package ptp

import (
	"reflect"
	"testing"
	"sync"
	"net"
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
/*
Generated TestPeerList_Init
Generated TestPeerList_operate
Generated TestPeerList_updateTables
Generated TestPeerList_deleteTables
Generated TestPeerList_Delete
Generated TestPeerList_Update
Generated TestPeerList_Get
Generated TestPeerList_GetPeer
Generated TestPeerList_GetEndpointAndProxy
Generated TestPeerList_GetID
Generated TestPeerList_Length
Generated TestPeerList_RunPeer
package ptp

import (
	"net"
	"reflect"
	"sync"
	"testing"
)
*/

/*
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
*/

func TestPeerList_Init(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.Init()
		})
	}
}

func TestPeerList_operate(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		action ListOperation
		id     string
		peer   *NetworkPeer
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.operate(tt.args.action, tt.args.id, tt.args.peer)
		})
	}
}

func TestPeerList_updateTables(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		id  string
		ip  string
		mac string
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.updateTables(tt.args.id, tt.args.ip, tt.args.mac)
		})
	}
}

func TestPeerList_deleteTables(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		ip  string
		mac string
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.deleteTables(tt.args.ip, tt.args.mac)
		})
	}
}

func TestPeerList_Delete(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.Delete(tt.args.id)
		})
	}
}

func TestPeerList_Update(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		id   string
		peer *NetworkPeer
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.Update(tt.args.id, tt.args.peer)
		})
	}
}

func TestPeerList_Get(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*NetworkPeer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			if got := l.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerList.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerList_GetPeer(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *NetworkPeer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			if got := l.GetPeer(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerList.GetPeer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerList_GetEndpointAndProxy(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		mac string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *net.UDPAddr
		want1   uint16
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			got, got1, err := l.GetEndpointAndProxy(tt.args.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerList.GetEndpointAndProxy() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerList.GetEndpointAndProxy() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PeerList.GetEndpointAndProxy() got1 = %v, want %v", got1, tt.want1)
			}
		})
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
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

func TestPeerList_Length(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			if got := l.Length(); got != tt.want {
				t.Errorf("PeerList.Length() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerList_RunPeer(t *testing.T) {
	type fields struct {
		peers      map[string]*NetworkPeer
		tableIPID  map[string]string
		tableMacID map[string]string
		lock       sync.RWMutex
	}
	type args struct {
		id string
		p  *PeerToPeer
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
			l := &PeerList{
				peers:      tt.fields.peers,
				tableIPID:  tt.fields.tableIPID,
				tableMacID: tt.fields.tableMacID,
				lock:       tt.fields.lock,
			}
			l.RunPeer(tt.args.id, tt.args.p)
		})
	}
}
