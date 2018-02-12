package main

import (
	"errors"
	"net"
	"strings"
	"time"

	"github.com/gogo/protobuf/proto"
	ptp "github.com/subutai-io/p2p/lib"
)

// DHT Errors
var (
	ErrorNoRouters        = errors.New("Routers wasn't specified")
	ErrorBadRouterAddress = errors.New("Bad router address")
)

// DHTConnection to a DHT bootstrap node
type DHTConnection struct {
	routers []*DHTRouter
}

type DHTRouter struct {
	conn       *net.TCPConn
	addr       *net.TCPAddr
	router     string
	running    bool
	handshaked bool
	fails      int
}

func (dht *DHTConnection) init(routersSrc string) error {
	ptp.Log(ptp.Info, "Initializing connection to a bootstrap nodes")
	routers := strings.Split(routersSrc, ",")
	if len(routers) == 0 {
		return ErrorNoRouters
	}
	for _, r := range routers {
		if r == "" {
			continue
		}
		addr, err := net.ResolveTCPAddr("tcp4", r)
		if err != nil {
			ptp.Log(ptp.Error, "Bad router address provided [%s]: %s", r, err)
			return ErrorBadRouterAddress
		}
		router := new(DHTRouter)
		router.addr = addr
		router.router = r
		dht.routers = append(dht.routers, router)
	}
	return nil
}

func (dht *DHTRouter) run() {
	dht.running = false
	dht.handshaked = false
	data := make([]byte, 4096)
	for {
		for !dht.running {
			dht.connect()
			if dht.running {
				break
			}
			dht.sleep()
		}
		n, err := dht.conn.Read(data)
		if err != nil {
			ptp.Log(ptp.Warning, "BSN socket closed: %s", err)
			dht.running = false
			dht.handshaked = false
			continue
		}
		go dht.routeData(data[:n])
	}
}

func (dht *DHTRouter) routeData(data []byte) {
	packet := &ptp.DHTPacket{}
	err := proto.Unmarshal(data, packet)
	if err != nil {
		ptp.Log(ptp.Warning, "Corrupted data from DHT: %s", err)
		return
	}
	if packet.Type == ptp.DHTPacketType_Ping {
		dht.handshaked = true
		ptp.Log(ptp.Info, "Connected to a bootstrap node: %s [%s]", dht.addr.String(), packet.Data)
		return
	}
	if !dht.handshaked {
		return
	}

}

func (dht *DHTRouter) connect() {
	dht.handshaked = false
	dht.running = false
	var err error
	dht.conn, err = net.DialTCP("tcp4", nil, dht.addr)
	if err != nil {
		dht.fails++
		ptp.Log(ptp.Error, "Failed to establish connection with %s", dht.addr.String())
		return
	}
	dht.fails = 0
	dht.running = true
}

func (dht *DHTRouter) sleep() {
	multiplier := dht.fails * 5
	if multiplier > 30 {
		multiplier = 30
	}
	ptp.Log(ptp.Info, "Waiting for %d second before reconnecting", multiplier)
	started := time.Now()
	timeout := time.Duration(time.Second * time.Duration(multiplier))
	for time.Since(started) < timeout {
		time.Sleep(time.Millisecond * 100)
	}
}
