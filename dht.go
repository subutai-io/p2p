package main

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
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
	routers    []*DHTRouter            // Routers
	lock       sync.Mutex              // Mutex for register/unregister
	instances  map[string]*P2PInstance // Instances
	registered []string                // List of registered swarm IDs
}

type DHTRouter struct {
	conn       *net.TCPConn // TCP connection to a bootsrap node
	addr       *net.TCPAddr // TCP address of a bootstrap node
	router     string       // Address of a bootstrap node
	running    bool         // Whether router is running or not
	handshaked bool         // Whether handshake has been completed or not
	stop       bool         // Whether service should be terminated
	fails      int          // Number of connection fails
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
	dht.instances = make(map[string]*P2PInstance)
	return nil
}

func (dht *DHTConnection) registerInstance(hash string, inst *P2PInstance) error {
	dht.lock.Lock()
	defer dht.lock.Unlock()
	ptp.Log(ptp.Debug, "Registering instance %s on bootstrap")

	exists := false
	for ihash, _ := range dht.instances {
		if hash == ihash {
			exists = true
			break
		}
	}
	for _, ihash := range dht.registered {
		if ihash == hash {
			exists = true
			break
		}
	}
	if exists {
		return fmt.Errorf("Hash already registered on bootstrap")
	}
	dht.instances[hash] = inst
	dht.registered = append(dht.registered, hash)
	inst.PTP.Dht.IncomingData = make(chan *ptp.DHTPacket)
	return nil
}

func (dht *DHTConnection) unregisterInstance(hash string) error {
	dht.lock.Lock()
	defer dht.lock.Unlock()
	ptp.Log(ptp.Debug, "Unregistering instance %s from bootstrap")
	_, e := dht.instances[hash]
	if !e {
		return fmt.Errorf("Can't unregister hash %s: Instance doesn't exists", hash)
	}
	delete(dht.instances, hash)
	for i, ihash := range dht.registered {
		if ihash == hash {
			dht.registered = append(dht.registered[:i], dht.registered[i+1:]...)
			break
		}
	}
	return nil
}

func (dht *DHTRouter) run() {
	dht.running = false
	dht.handshaked = false
	data := make([]byte, 4096)
	for !dht.stop {
		for !dht.running {
			dht.connect()
			if dht.running {
				break
			}
			if dht.stop {
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

		supported := false
		for _, v := range ptp.SupportedVersion {
			if v == packet.Version {
				supported = true
			}
		}
		if !supported {
			ptp.Log(ptp.Error, "Version mismatch. Server have %d. We have %d", packet.Version, ptp.PacketVersion)
			dht.stop = true
			dht.conn.Close()
		} else {
			dht.handshaked = true
			ptp.Log(ptp.Info, "Connected to a bootstrap node: %s [%s]", dht.addr.String(), packet.Data)
		}
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
		ptp.Log(ptp.Error, "Failed to establish connection with %s: %s", dht.addr.String(), err)
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
		time.Sleep(time.Millisecond * 200)
	}
}
