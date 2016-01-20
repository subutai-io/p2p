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

func (p *Proxy) Initialize(target string) {
	p.UDPServer = new(udpcs.UDPClient)
	p.UDPServer.Init("", 0)
	p.DHTClient = new(dht.DHTClient)
	config := p.DHTClient.DHTClientConfig()
	if target != "" {
		config.Routers = target
	}
	config.Mode = dht.MODE_CP
	config.NetworkHash = p.GenerateHash()
	config.P2PPort = p.UDPServer.GetPort()
	log.Log(log.INFO, "Listening on a %d port", config.P2PPort)
	var ips []net.IP
	ips = append(ips, net.ParseIP("127.0.0.1"))
	p.DHTClient = p.DHTClient.Initialize(config, ips)
	p.DHTClient.RegisterControlPeer()
	p.UDPServer.Listen(p.HandleMessage)
}

func (p *Proxy) GenerateHash() string {
	var infohash string
	t := time.Now()
	infohash = "cp" + fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
	return infohash
}

func (p *Proxy) HandleMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	log.Log(log.DEBUG, "MSG RECEIVED")
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
		log.Log(log.DEBUG, "Proxy message received")
		// Register forwarding
		// Go over list of proxies in use and find particular IP in it
		// If it's found - return Proxy ID. Create new entry otherwise
		data := string(msg.Data)
		var responseId int = -1
		var respAddr string = data
		targetIp, _ := net.ResolveUDPAddr("udp", data)
		for id, tunnel := range p.Tunnels {
			if tunnel.Peer1 == src_addr {
				if tunnel.Peer2 == targetIp {
					responseId = id
				}
			} else if tunnel.Peer2 == src_addr {
				if tunnel.Peer1 == targetIp {
					responseId = id
				}
			}
		}
		if responseId == -1 {
			// We didn't found any matches. Let's create new entry
			var t Tunnel
			t.Peer1 = src_addr
			t.Peer2, _ = net.ResolveUDPAddr("udp", data)
			for i := 0; i < len(p.Tunnels)+1; i++ {
				_, exists := p.Tunnels[i]
				if !exists {
					log.Log(log.DEBUG, "New tunnel has been created")
					p.Tunnels[i] = t
					break
				}
			}
		}
		msg := udpcs.CreateProxyP2PMessage(responseId, respAddr, 0)
		p.UDPServer.SendMessage(msg, src_addr)
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
