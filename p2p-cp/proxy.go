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
	UDPServer   *ptp.PTPNet
	Shutdown    bool
	TunnelQueue []WaitingTunnel
}

type WaitingTunnel struct {
	Target     string
	Source     string
	Port       uint16
	Registered bool
}

// Tunnel established between two peers. Tunnels doesn't
// provide two-way connectivity.
type Tunnel struct {
	PingFails   int
	Endpoint    *net.UDPAddr
	RealAddress *net.UDPAddr
	Ready       bool
}

func (p *Proxy) StartUDPServer(port int) {
	p.UDPServer = new(ptp.PTPNet)
	p.UDPServer.Init("", port)
	go p.UDPServer.Listen(p.HandleMessage)
}

func (p *Proxy) StartDHT(target string) {
	p.DHTClient = new(ptp.DHTClient)
	config := p.DHTClient.DHTClientConfig()
	config.Mode = ptp.MODE_CP
	config.NetworkHash = p.GenerateHash()
	config.P2PPort = p.UDPServer.GetPort()
	if target != "" {
		config.Routers = target
	}
	var ips []net.IP
	ips = append(ips, net.ParseIP("127.0.0.1"))
	ch := make(chan []ptp.PeerIP)
	proxych := make(chan ptp.Forwarder)
	p.DHTClient = p.DHTClient.Initialize(config, ips, ch, proxych)
	for len(p.DHTClient.ID) < 32 {
		p.DHTClient.ID = ""
		ptp.Log(ptp.WARNING, "Failed to connect to DHT. Retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		p.DHTClient = p.DHTClient.Initialize(config, ips, ch, proxych)
	}
	p.DHTClient.RegisterControlPeer()
}

func (p *Proxy) Initialize() {
	p.Tunnels = make(map[int]Tunnel)
	go p.RegisterQueue()
	ptp.Log(ptp.INFO, "Control peer initialization process is complete")
}

func (p *Proxy) Stop() {
	p.UDPServer.Stop()
	p.Shutdown = true
}

func (p *Proxy) GenerateHash() string {
	var infohash string
	t := time.Now()
	infohash = "cp" + fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
	return infohash
}

func (p *Proxy) CreateTunnel(addr, origin *net.UDPAddr, ready bool) int {
	var newId int = 0
	var t Tunnel
	t.Endpoint = addr
	t.RealAddress = origin
	t.PingFails = 0
	t.Ready = ready
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
	target_addr, _ := net.ResolveUDPAddr("udp", target)
	source := p.TunnelQueue[0].Source
	src_addr, _ := net.ResolveUDPAddr("udp", source)
	port := p.TunnelQueue[0].Port
	s_ip := src_addr.IP.String()
	realSource, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", s_ip, port))
	ptp.Log(ptp.DEBUG, "Requested proxy for %s from %s [RS: %s]", target, source, realSource)
	// Check if we are in the list
	available := false
	var foundId int
	for id, tun := range p.Tunnels {
		if tun.Endpoint.String() == realSource.String() {
			available = true
			foundId = id
		}
	}
	if !available {
		nId := p.CreateTunnel(realSource, src_addr, true)
		if nId > 0 {
			ptp.Log(ptp.DEBUG, "Requester peer %s was not found in tunnels list. Created new tunnel with ID %d", source, nId)
		}
	} else {
		t, exists := p.Tunnels[foundId]
		if exists && foundId > 0 {
			t.Ready = true
			t.RealAddress = src_addr
			p.Tunnels[foundId] = t
		}
	}
	var responseId int = -1
	for id, tun := range p.Tunnels {
		if tun.Endpoint.String() == target {
			ptp.Log(ptp.DEBUG, "Proxy %d found for peer %s", id, target)
			responseId = int(id)
		}
	}
	if responseId == -1 {
		ptp.Log(ptp.DEBUG, "Tunnel for %s was not found", target)
		responseId = p.CreateTunnel(target_addr, nil, false)
	}
	if responseId < 0 {
		ptp.Log(ptp.ERROR, "Failed to create tunnel from %s to %s", source, target)
	}
	response := ptp.CreateProxyP2PMessage(responseId, target, 0)
	p.UDPServer.SendMessage(response, src_addr)
	p.TunnelQueue[0].Registered = true
	p.Lock = false

	p.DHTClient.ReportControlPeerLoad(len(p.Tunnels))
}

func (p *Proxy) HandleMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	ptp.Log(ptp.TRACE, "Received")
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
		w.Port = msg.Header.NetProto
		p.TunnelQueue = append(p.TunnelQueue, w)
	} else if msgType == ptp.MT_PING {
		for key, tun := range p.Tunnels {
			if tun.RealAddress.String() == src_addr.String() {
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
			if !tunnel.Ready {
				ptp.Log(ptp.DEBUG, "Proxy %d is not ready", msg.Header.ProxyId)
				return
			}
			ptp.Log(ptp.DEBUG, "Forwarding from %s to %s. Proxy ID: %d", src_addr.String(), tunnel.Endpoint.String(), msg.Header.ProxyId)
			p.UDPServer.SendMessage(msg, tunnel.RealAddress)
		}
	}
}

func (p *Proxy) SendPing() {
	for key, tunnel := range p.Tunnels {
		tunnel.PingFails += tunnel.PingFails + 1
		msg := ptp.CreatePingP2PMessage()
		p.UDPServer.SendMessage(msg, tunnel.RealAddress)
		p.Tunnels[key] = tunnel
	}
}

func (p *Proxy) CleanTunnels() {
	for key, tunnel := range p.Tunnels {
		if (tunnel.Ready && tunnel.PingFails > 3) || (!tunnel.Ready && tunnel.PingFails > 20) {
			ptp.Log(ptp.DEBUG, "Removing outdated proxy: %d", key)
			delete(p.Tunnels, key)
			badId := key
			p.NotifyBadTunnel(badId)
			p.DHTClient.ReportControlPeerLoad(len(p.Tunnels))
		}
	}
}

func (p *Proxy) RegisterQueue() {
	for {
		if p.Shutdown {
			break
		}
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

func (p *Proxy) NotifyBadTunnel(id int) {
	msg := ptp.CreateBadTunnelP2PMessage(id, 1)
	for _, t := range p.Tunnels {
		if !t.Ready {
			continue
		}
		p.UDPServer.SendMessage(msg, t.Endpoint)
	}
}
