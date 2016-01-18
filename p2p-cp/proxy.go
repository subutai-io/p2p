package main

// Connect to a DHT
// Register with DHT
// Wait for incoming connections
// Validate incoming connections with DHT

import (
	"fmt"
	"net"
	"p2p/commons"
	"p2p/dht"
	log "p2p/p2p_log"
	"p2p/udpcs"
	"time"
)

type Proxy struct {
	DHTClient *dht.DHTClient
	Tunnels   map[int]Tunnel
	UDPServer *udpcs.UDPClient
	Shutdown  bool
}

// Tunnel established between two peers. Tunnels doesn't
// provide two-way connectivity.
type Tunnel struct {
	Peer1      *net.UDPAddr
	Peer2      *net.UDPAddr
	UniqueHash string
}

func (p *Proxy) Initialize() {
	p.UDPServer = new(udpcs.UDPClient)
	p.UDPServer.Init("", 0)
	p.DHTClient = new(dht.DHTClient)
	config := p.DHTClient.DHTClientConfig()
	config.NetworkHash = p.GenerateHash()
	config.P2PPort = p.UDPServer.GetPort()
	var ips []net.IP
	ips = append(ips, net.ParseIP("127.0.0.1"))
	p.DHTClient = p.DHTClient.Initialize(config, ips)
	p.UDPServer.Listen(p.HandleMessage)
}

func (p *Proxy) GenerateHash() string {
	var infohash string
	t := time.Now()
	infohash = "cp" + fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
	return infohash
}

func (p *Proxy) HandleMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		log.Log(log.ERROR, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := udpcs.P2PMessageFromBytes(buf)
	if des_err != nil {
		log.Log(log.ERROR, "P2PMessageFromBytes error: %v", des_err)
		return
	}
	var msgType commons.MSG_TYPE = commons.MSG_TYPE(msg.Header.Type)
	if msgType == commons.MT_PROXY {
		// Register forwarding
		data = string(msg.Data)
		for key, tunnel := range p.Tunnels {
			if data == tunnel.UniqueHash {
			}
		}
	} else {
		// Forward message
		tunnel, exists := p.Tunnels[int(msg.Header.ProxyId)]
		if !exists {
			log.Log(log.WARNING, "Proxy %d is not registered", msg.Header.ProxyId)
			return
		}
		if tunnel.Peer1.String() == src_addr.String() {
			p.UDPServer.SendMessage(msg, tunnel.Peer2)
		} else if tunnel.Peer2.String() == src_addr.String() {
			p.UDPServer.SendMessage(msg, tunnel.Peer1)
		} else {
			log.Log(log.WARNING, "Connected peer doesn't belong to requested proxy")
		}
	}
}
