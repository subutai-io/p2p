package main

import (
	"net"
	"testing"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
	"github.com/golang/protobuf/proto"
)

func TestRouteData(t *testing.T) {
	dataStr := "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	dataArr := []string{}
	for i := 0; i < 99; i++ {
		dataArr = append(dataArr, dataStr)
	}
	data := []ptp.DHTPacket{
		ptp.DHTPacket{},
		ptp.DHTPacket{
			Type: ptp.DHTPacketType_Ping,
		},
		ptp.DHTPacket{
			Type:      ptp.DHTPacketType_Ping,
			Arguments: dataArr,
		},
	}
	router := new(DHTRouter)
	for _, d := range data {
		b, _ := proto.Marshal(&d)
		if len(b) > ptp.DHTBufferSize {
			b = b[:ptp.DHTBufferSize]
		}
		router.routeData(b)
	}
}

/*
Generated TestDHTRouter_run
Generated TestDHTRouter_handleData
Generated TestDHTRouter_routeData
Generated TestDHTRouter_connect
Generated TestDHTRouter_sleep
Generated TestDHTRouter_keepAlive
Generated TestDHTRouter_sendRaw
Generated TestDHTRouter_ping
*/

func TestDHTRouter_run(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.run()
		})
	}
}

func TestDHTRouter_handleData(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	type args struct {
		data []byte
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
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.handleData(tt.args.data)
		})
	}
}

func TestDHTRouter_routeData(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	type args struct {
		data []byte
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
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.routeData(tt.args.data)
		})
	}
}

func TestDHTRouter_connect(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.connect()
		})
	}
}

func TestDHTRouter_sleep(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.sleep()
		})
	}
}

func TestDHTRouter_keepAlive(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			dht.keepAlive()
		})
	}
}

func TestDHTRouter_sendRaw(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			got, err := dht.sendRaw(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DHTRouter.sendRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("DHTRouter.sendRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDHTRouter_ping(t *testing.T) {
	type fields struct {
		conn        *net.TCPConn
		addr        *net.TCPAddr
		router      string
		running     bool
		handshaked  bool
		stop        bool
		fails       int
		tx          uint64
		rx          uint64
		data        chan *ptp.DHTPacket
		lastContact time.Time
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
			dht := &DHTRouter{
				conn:        tt.fields.conn,
				addr:        tt.fields.addr,
				router:      tt.fields.router,
				running:     tt.fields.running,
				handshaked:  tt.fields.handshaked,
				stop:        tt.fields.stop,
				fails:       tt.fields.fails,
				tx:          tt.fields.tx,
				rx:          tt.fields.rx,
				data:        tt.fields.data,
				lastContact: tt.fields.lastContact,
			}
			if err := dht.ping(); (err != nil) != tt.wantErr {
				t.Errorf("DHTRouter.ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
