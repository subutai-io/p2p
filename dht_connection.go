package main

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/golang/protobuf/proto"
	ptp "github.com/subutai-io/p2p/lib"
	"github.com/subutai-io/p2p/protocol"
)

// DHT Errors
var (
	ErrorNoRouters        = errors.New("Routers wasn't specified")
	ErrorBadRouterAddress = errors.New("Bad router address")
)

// DHTConnection to a DHT bootstrap node
type DHTConnection struct {
	routers     []*DHTRouter             // Bootstrap nodes
	routersList map[int]string           // List of bootstrap nodes received from SRV lookup
	lock        sync.Mutex               // Mutex for register/unregister
	instances   map[string]*P2PInstance  // Instances
	registered  []string                 // List of registered swarm IDs
	incoming    chan *protocol.DHTPacket // Packets received by routers
	ip          string                   // Our outbound IP
	isActive    bool                     // Whether DHT connection is active or not
}

func (dht *DHTConnection) init(target string) error {
	ptp.Log(ptp.Info, "Initializing connection to a bootstrap nodes")
	dht.incoming = make(chan *protocol.DHTPacket)
	var err error
	dht.routersList, err = ptp.SrvLookup(target, "tcp", "subutai.io")
	if err != nil {
		ptp.Log(ptp.Error, "Failed to get bootstrap nodes: %s", err.Error)
		dht.routersList = make(map[int]string)
	}
	if len(dht.routersList) == 0 {
		return ErrorNoRouters
	}
	for _, r := range dht.routersList {
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
		router.data = dht.incoming
		dht.routers = append(dht.routers, router)
	}
	dht.instances = make(map[string]*P2PInstance)
	return nil
}

func (dht *DHTConnection) registerInstance(hash string, inst *P2PInstance) error {
	dht.lock.Lock()
	defer dht.lock.Unlock()
	ptp.Log(ptp.Debug, "Registering instance %s on bootstrap", hash)

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
	inst.PTP.Dht.IncomingData = make(chan *protocol.DHTPacket)
	inst.PTP.Dht.OutgoingData = make(chan *protocol.DHTPacket)
	go func() {
		for {
			packet := <-inst.PTP.Dht.OutgoingData
			if packet == nil {
				break
			}
			dht.send(packet)
		}
	}()
	ptp.Log(ptp.Debug, "Instance was registered with bootstrap client")
	return nil
}

func (dht *DHTConnection) send(packet *protocol.DHTPacket) {
	if packet == nil {
		return
	}
	ptp.Log(ptp.Trace, "Sending DHT packet %+v", packet)
	data, err := proto.Marshal(packet)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to marshal DHT Packet: %s", err)
	}
	ptp.Log(ptp.Trace, "Sending marshaled DHT Packet of size [%d]", len(data))
	for i, router := range dht.routers {
		if router.running && router.handshaked {
			n, err := router.sendRaw(data)
			if err != nil {
				ptp.Log(ptp.Error, "Failed to send data to %s", router.addr.String())
				continue
			}
			if n >= 0 {
				dht.routers[i].tx += uint64(n)
			}
		}
	}
}

func (dht *DHTConnection) run() {
	for {
		packet := <-dht.incoming
		if packet == nil {
			continue
		}
		ptp.Log(ptp.Trace, "Routing DHT Packet %+v", packet)
		// Ping should always provide us with outbound IP value
		if packet.Type == protocol.DHTPacketType_Ping && packet.Data != "" {
			dht.ip = packet.Data
			continue
		}
		if packet.Infohash == "" {
			continue
		}
		i, e := dht.instances[packet.Infohash]
		if e && i != nil && i.PTP != nil && !i.PTP.Shutdown && i.PTP.Dht != nil && i.PTP.Dht.IncomingData != nil {
			i.PTP.Dht.IncomingData <- packet
		} else {
			ptp.Log(ptp.Debug, "DHT received data for unknown instance %s: %+v", packet.Infohash, packet)
		}
	}
}

func (dht *DHTConnection) unregisterInstance(hash string) error {
	dht.lock.Lock()
	defer dht.lock.Unlock()
	ptp.Log(ptp.Debug, "Unregistering instance %s from bootstrap", hash)
	inst, e := dht.instances[hash]
	if !e {
		return fmt.Errorf("Can't unregister hash %s: Instance doesn't exists", hash)
	}
	if inst != nil && inst.PTP != nil && inst.PTP.Dht != nil {
		err := inst.PTP.Dht.Close()
		if err != nil {
			ptp.Log(ptp.Error, "Failed to stop DHT on instance %s", hash)
		}
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
