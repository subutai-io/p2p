package ptp

import (
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestProxyManager_init(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Dummy", fields{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if err := p.init(); (err != nil) != tt.wantErr {
				t.Errorf("ProxyManager.init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyManager_operate(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	type args struct {
		operation ListOperation
		addr      string
		proxy     *proxyServer
	}
	proxies := make(map[string]*proxyServer)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"No operation", fields{proxies: proxies}, args{}},
		{"Update", fields{proxies: proxies}, args{operation: OperateUpdate}},
		{"Delete", fields{proxies: proxies}, args{operation: OperateDelete}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			p.operate(tt.args.operation, tt.args.addr, tt.args.proxy)
		})
	}
}

func TestProxyManager_get(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	empty := make(map[string]*proxyServer)
	r1 := make(map[string]*proxyServer)
	addr, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")
	r1[addr.String()] = &proxyServer{
		Addr:       addr,
		Endpoint:   addr,
		Status:     proxyConnecting,
		LastUpdate: time.Now(),
		Created:    time.Now(),
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*proxyServer
	}{
		{"Empty", fields{}, empty},
		{"Map", fields{proxies: r1}, r1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if got := p.get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProxyManager.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyManager_GetList(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}

	addr, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")

	r1 := make(map[string]*proxyServer)
	pr1 := &proxyServer{
		Addr:       addr,
		Endpoint:   addr,
		Status:     proxyConnecting,
		LastUpdate: time.Now(),
		Created:    time.Now(),
	}
	r1[addr.String()] = pr1
	r2 := []*proxyServer{pr1}

	tests := []struct {
		name   string
		fields fields
		want   []*proxyServer
	}{
		{"Empty", fields{}, []*proxyServer{}},
		{"Map", fields{proxies: r1}, r2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if got := p.GetList(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProxyManager.GetList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyManager_new(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	type args struct {
		endpoint *net.UDPAddr
	}

	ep1, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")
	r0 := make(map[string]*proxyServer)
	r1 := make(map[string]*proxyServer)
	r1[ep1.String()] = new(proxyServer)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"New", fields{proxies: r0}, args{ep1}, false},
		{"Existing", fields{proxies: r1}, args{ep1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if err := p.new(tt.args.endpoint); (err != nil) != tt.wantErr {
				t.Errorf("ProxyManager.new() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProxyManager_check(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}

	r1 := make(map[string]*proxyServer)
	r2 := make(map[string]*proxyServer)
	r3 := make(map[string]*proxyServer)
	a1, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")
	p1 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Now(),
		LastUpdate: time.Unix(1, 1),
		Status:     proxyActive,
	}
	p2 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Unix(1, 1),
		LastUpdate: time.Now(),
		Status:     proxyConnecting,
	}
	p3 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Unix(1, 1),
		LastUpdate: time.Now(),
		Status:     proxyDisconnected,
	}

	r1[a1.String()] = p1
	r2[a1.String()] = p2
	r3[a1.String()] = p3

	tests := []struct {
		name   string
		fields fields
	}{
		{"Empty", fields{}},
		{"Active proxy", fields{proxies: r1}},
		{"Connecting proxy", fields{proxies: r2}},
		{"Disconnecting proxy", fields{proxies: r3}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			p.check()
		})
	}
}

func TestProxyManager_touch(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	type args struct {
		id string
	}

	r1 := make(map[string]*proxyServer)
	a1, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")
	p1 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Now(),
		LastUpdate: time.Unix(1, 1),
		Status:     proxyActive,
	}
	r1[a1.String()] = p1

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Empty", fields{}, args{}, false},
		{"Wrong ID", fields{proxies: r1}, args{"wrongid"}, false},
		{"Correct ID", fields{proxies: r1}, args{a1.String()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if got := p.touch(tt.args.id); got != tt.want {
				t.Errorf("ProxyManager.touch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProxyManager_activate(t *testing.T) {
	type fields struct {
		proxies    map[string]*proxyServer
		lock       sync.RWMutex
		hasChanges bool
	}
	type args struct {
		id       string
		endpoint *net.UDPAddr
	}

	r1 := make(map[string]*proxyServer)
	r2 := make(map[string]*proxyServer)
	a1, _ := net.ResolveUDPAddr("udp4", "1.2.3.4:1234")
	p1 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Now(),
		LastUpdate: time.Unix(1, 1),
		Status:     proxyActive,
	}
	p2 := &proxyServer{
		Addr:       a1,
		Endpoint:   a1,
		Created:    time.Now(),
		LastUpdate: time.Unix(1, 1),
		Status:     proxyConnecting,
	}
	r1[a1.String()] = p1
	r2[a1.String()] = p2

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"Empty", fields{}, args{}, false},
		{"Different state", fields{proxies: r1}, args{id: a1.String()}, false},
		{"Correct state", fields{proxies: r2}, args{id: a1.String()}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ProxyManager{
				proxies:    tt.fields.proxies,
				lock:       tt.fields.lock,
				hasChanges: tt.fields.hasChanges,
			}
			if got := p.activate(tt.args.id, tt.args.endpoint); got != tt.want {
				t.Errorf("ProxyManager.activate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_proxyServer_Close(t *testing.T) {
	type fields struct {
		Addr       *net.UDPAddr
		Endpoint   *net.UDPAddr
		Status     proxyStatus
		LastUpdate time.Time
		Created    time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
	// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &proxyServer{
				Addr:       tt.fields.Addr,
				Endpoint:   tt.fields.Endpoint,
				Status:     tt.fields.Status,
				LastUpdate: tt.fields.LastUpdate,
				Created:    tt.fields.Created,
			}
			if err := p.Close(); (err != nil) != tt.wantErr {
				t.Errorf("proxyServer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
