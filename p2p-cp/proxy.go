package main

// Connect to a DHT
// Register with DHT
// Wait for incoming connections
// Validate incoming connections with DHT

import (
	"fmt"
	"github.com/subutai-io/p2p/commons"
	"github.com/subutai-io/p2p/dht"
	log "github.com/subutai-io/p2p/p2p_log"
	"github.com/subutai-io/p2p/udpcs"
	"net"
	"time"
)

type Proxy struct {
	DHTClient *dht.DHTClient
	Tunnels   map[uint16]Tunnel
	UDPServer *udpcs.UDPClient
	Shutdown  bool
}

// Tunnel established between two peers. Tunnels doesn't
// provide two-way connectivity.
type Tunnel struct {
	Peer1      *net.UDPAddr
	Peer2      *net.UDPAddr
	UniqueHash string
	PingFails  int
	PingFails1 int
	PingFails2 int
	Endpoint   *net.UDPAddr
}

func (p *Proxy) Initialize(target string) {
	p.UDPServer = new(udpcs.UDPClient)
	p.UDPServer.Init("", 0)
	p.DHTClient = new(dht.DHTClient)
	p.Tunnels = make(map[uint16]Tunnel)
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
	go p.UDPServer.Listen(p.HandleMessage)
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
		// MT_PROXY indicates that peer (src_addr) can't connect to another peer (msg.data)
		data := string(msg.Data)
		var responseId int = -1
		for id, tun := range p.Tunnels {
			if tun.Endpoint.String() == data {
				responseId = int(id)
			}
		}
		if responseId == -1 {
			log.Log(log.DEBUG, "Making new tunnel")
			var t Tunnel
			t.Endpoint, _ = net.ResolveUDPAddr("udp", data)
			t.PingFails = 0
			for i := 1; i < len(p.Tunnels)+2; i++ {
				_, exists := p.Tunnels[uint16(i)]
				if !exists {
					p.Tunnels[uint16(i)] = t
					responseId = i
					break
				}
			}
		}
		msg := udpcs.CreateProxyP2PMessage(responseId, data, 0)
		p.UDPServer.SendMessage(msg, src_addr)
		p.DHTClient.ReportControlPeerLoad(len(p.Tunnels))
	} else if msgType == commons.MT_PING {
		for key, tun := range p.Tunnels {
			if tun.Peer1.String() == src_addr.String() {
				tun.PingFails1 = 0
			} else if tun.Peer2.String() == src_addr.String() {
				tun.PingFails2 = 0
			}
			p.Tunnels[key] = tun
		}
	} else {
		tunnel, exists := p.Tunnels[msg.Header.ProxyId]
		if !exists {
			log.Log(log.DEBUG, "Proxy %d is not registered", msg.Header.ProxyId)
			return
		}
		p.UDPServer.SendMessage(msg, tunnel.Endpoint)
	}
}

func (p *Proxy) HandleMessageOld(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
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
		targetIp, _ := net.ResolveUDPAddr("udp", data)
		for id, tunnel := range p.Tunnels {
			if tunnel.Peer1 == src_addr {
				if tunnel.Peer2 == targetIp {
					responseId = int(id)
				}
			} else if tunnel.Peer2 == src_addr {
				if tunnel.Peer1 == targetIp {
					responseId = int(id)
				}
			}
		}
		if responseId == -1 {
			log.Log(log.DEBUG, "Making new tunnel")
			// We didn't found any matches. Let's create new entry
			var t Tunnel
			t.Peer1 = src_addr
			t.Peer2, _ = net.ResolveUDPAddr("udp", data)
			t.PingFails1 = 0
			t.PingFails2 = 0
			for i := 1; i < len(p.Tunnels)+2; i++ {
				_, exists := p.Tunnels[uint16(i)]
				if !exists {
					log.Log(log.DEBUG, "New tunnel has been created with ID %d", i)
					p.Tunnels[uint16(i)] = t
					responseId = i
					break
				}
			}
		}
		msg := udpcs.CreateProxyP2PMessage(responseId, data, 0)
		p.UDPServer.SendMessage(msg, src_addr)
		// Notify about new tunnel
		p.DHTClient.ReportControlPeerLoad(len(p.Tunnels))
	} else if msgType == commons.MT_PING {
		for key, tun := range p.Tunnels {
			if tun.Peer1.String() == src_addr.String() {
				tun.PingFails1 = 0
			} else if tun.Peer2.String() == src_addr.String() {
				tun.PingFails2 = 0
			}
			p.Tunnels[key] = tun
		}
	} else {
		//log.Log(log.DEBUG, "PROXY: %v", p.Tunnels)
		// Forward message
		tunnel, exists := p.Tunnels[msg.Header.ProxyId]
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

func (p *Proxy) SendPing() {
	for key, tunnel := range p.Tunnels {
		tunnel.PingFails1 += tunnel.PingFails1 + 1
		tunnel.PingFails2 += tunnel.PingFails2 + 1
		msg := udpcs.CreatePingP2PMessage()
		p.UDPServer.SendMessage(msg, tunnel.Peer1)
		p.UDPServer.SendMessage(msg, tunnel.Peer2)
		p.Tunnels[key] = tunnel
	}
}

func (p *Proxy) CleanTunnels() {
	for key, tunnel := range p.Tunnels {
		if tunnel.PingFails1 > 3 || tunnel.PingFails2 > 3 {
			delete(p.Tunnels, key)
		}
	}
}
