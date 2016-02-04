package main

// Connect to a DHT
// Register with DHT
// Wait for incoming connections
// Validate incoming connections with DHT

import (
	"fmt"
	ptp "github.com/subutai-io/p2p/lib"
	"net"
	"time"
)

type Proxy struct {
	DHTClient   *ptp.DHTClient
	Tunnels     map[int]Tunnel
	Lock        bool
	UDPServer   *ptp.UDPClient
	Shutdown    bool
	TunnelQueue []WaitingTunnel
}

type WaitingTunnel struct {
	Target     string
	Source     string
	Registered bool
}

// Tunnel established between two peers. Tunnels doesn't
// provide two-way connectivity.
type Tunnel struct {
	PingFails int
	Endpoint  *net.UDPAddr
}

func (p *Proxy) Initialize(target string) {
	p.UDPServer = new(ptp.UDPClient)
	p.UDPServer.Init("", 0)
	p.DHTClient = new(ptp.DHTClient)
	p.Tunnels = make(map[int]Tunnel)
	config := p.DHTClient.DHTClientConfig()
	if target != "" {
		config.Routers = target
	}
	config.Mode = ptp.MODE_CP
	config.NetworkHash = p.GenerateHash()
	config.P2PPort = p.UDPServer.GetPort()
	ptp.Log(ptp.INFO, "Listening on a %d port", config.P2PPort)
	var ips []net.IP
	ips = append(ips, net.ParseIP("127.0.0.1"))
	go p.UDPServer.Listen(p.HandleMessage)
	go p.RegisterQueue()
	p.DHTClient = p.DHTClient.Initialize(config, ips)
	p.DHTClient.RegisterControlPeer()
	ptp.Log(ptp.INFO, "Control peer initialization process is complete")
}

func (p *Proxy) GenerateHash() string {
	var infohash string
	t := time.Now()
	infohash = "cp" + fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
	return infohash
}

func (p *Proxy) CreateTunnel(addr string) int {
	var newId int = 0
	var t Tunnel
	t.Endpoint, _ = net.ResolveUDPAddr("udp", addr)
	t.PingFails = 0
	for i := 1; i < len(p.Tunnels)+2; i++ {
		_, exists := p.Tunnels[i]
		if !exists {
			p.Tunnels[i] = t
			newId = i
			break
		}
	}
	ptp.Log(ptp.DEBUG, "Created new tunnel. ID: %d Endpoint: %s", newId, addr)
	return newId
}

func (p *Proxy) RegisterTunnel() {
	if len(p.TunnelQueue) == 0 {
		return
	}
	p.Lock = true
	target := p.TunnelQueue[0].Target
	source := p.TunnelQueue[0].Source
	ptp.Log(ptp.DEBUG, "Size of map is %d", len(p.Tunnels))
	ptp.Log(ptp.DEBUG, "Requested proxy for %s from %s", target, source)
	// Check if we are in the list
	available := false
	for _, tun := range p.Tunnels {
		if tun.Endpoint.String() == source {
			available = true
		}
	}
	if !available {
		nId := p.CreateTunnel(source)
		if nId > 0 {
			ptp.Log(ptp.DEBUG, "Requester peer %s was not found in tunnels list. Creating new one with ID %d", source, nId)
		}
	}
	// MT_PROXY indicates that peer (src_addr) can't connect to another peer (msg.data)
	var responseId int = -1
	for id, tun := range p.Tunnels {
		if tun.Endpoint.String() == target {
			ptp.Log(ptp.DEBUG, "Proxy %d found for peer %s", id, target)
			responseId = int(id)
		}
	}
	if responseId == -1 {
		ptp.Log(ptp.DEBUG, "Tunnel for %s was not found", target)
		responseId = p.CreateTunnel(target)
	}
	if responseId < 0 {
		ptp.Log(ptp.ERROR, "Failed to create tunnel from %s to %s", source, target)
	}
	response := ptp.CreateProxyP2PMessage(responseId, target, 0)
	src_addr, _ := net.ResolveUDPAddr("udp", source)
	p.UDPServer.SendMessage(response, src_addr)
	p.TunnelQueue[0].Registered = true
	p.Lock = false

	p.DHTClient.ReportControlPeerLoad(len(p.Tunnels))
}

func (p *Proxy) HandleMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		ptp.Log(ptp.ERROR, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := ptp.P2PMessageFromBytes(buf)
	if des_err != nil {
		ptp.Log(ptp.ERROR, "P2PMessageFromBytes error: %v", des_err)
		return
	}
	var msgType ptp.MSG_TYPE = ptp.MSG_TYPE(msg.Header.Type)
	if msgType == ptp.MT_PROXY {
		var w WaitingTunnel
		w.Target = string(msg.Data)
		w.Source = src_addr.String()
		p.TunnelQueue = append(p.TunnelQueue, w)
	} else if msgType == ptp.MT_PING {
		for key, tun := range p.Tunnels {
			if tun.Endpoint.String() == src_addr.String() {
				tun.PingFails = 0
			}
			p.Tunnels[key] = tun
		}
	} else {
		if msg.Header.ProxyId > 0 {
			tunnel, exists := p.Tunnels[int(msg.Header.ProxyId)]
			if !exists {
				ptp.Log(ptp.DEBUG, "Proxy %d is not registered", msg.Header.ProxyId)
				return
			}
			ptp.Log(ptp.DEBUG, "Forwarding from %s to %s. Proxy ID: %d", src_addr.String(), tunnel.Endpoint.String(), msg.Header.ProxyId)
			p.UDPServer.SendMessage(msg, tunnel.Endpoint)
		}
	}
}

func (p *Proxy) SendPing() {
	for key, tunnel := range p.Tunnels {
		tunnel.PingFails += tunnel.PingFails + 1
		msg := ptp.CreatePingP2PMessage()
		p.UDPServer.SendMessage(msg, tunnel.Endpoint)
		p.Tunnels[key] = tunnel
	}
}

func (p *Proxy) CleanTunnels() {
	for key, tunnel := range p.Tunnels {
		if tunnel.PingFails > 3 {
			delete(p.Tunnels, key)
		}
	}
}

func (p *Proxy) RegisterQueue() {
	for {
		time.Sleep(1 * time.Second)
		if len(p.TunnelQueue) == 0 {
			continue
		}
		if p.Lock {
			continue
		}
		p.RegisterTunnel()
		for i, t := range p.TunnelQueue {
			if t.Registered {
				p.TunnelQueue = append(p.TunnelQueue[:i], p.TunnelQueue[i+1:]...)
			}
		}
	}
}
