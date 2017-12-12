package ptp

import (
	"fmt"
	"net"
	"time"
)

type proxyStatus uint8

const (
	proxyConnecting   proxyStatus = 0
	proxyActive       proxyStatus = 1
	proxyDisconnected proxyStatus = 2
)

type proxyServer struct {
	Addr       *net.UDPAddr
	Endpoint   *net.UDPAddr // Endpoint provided by proxy
	Status     proxyStatus
	LastUpdate time.Time
}

func (p *PeerToPeer) initProxy(addr string) error {
	Log(Info, "Initializing proxy %s", addr)
	var err error
	proxy := new(proxyServer)
	proxy.LastUpdate = time.Now()
	proxy.Addr, err = net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return fmt.Errorf("Failed to resolve proxy address")
	}
	for _, pr := range p.Proxies {
		if pr.Addr.String() == proxy.Addr.String() {
			return fmt.Errorf("Proxy %s already exists", addr)
		}
	}
	p.Proxies = append(p.Proxies, proxy)
	initStarted := time.Now()
	proxy.Status = proxyConnecting

	msg := CreateProxyP2PMessage(0, p.Dht.ID, 1)
	p.UDPSocket.SendMessage(msg, proxy.Addr)
	for proxy.Status == proxyConnecting {
		time.Sleep(100 * time.Millisecond)
		if time.Duration(3*time.Second) < time.Since(initStarted) {
			p.removeProxy(proxy.Addr)
			return fmt.Errorf("Failed to connect to proxy")
		}
	}
	if proxy.Status != proxyActive {
		p.removeProxy(proxy.Addr)
		return fmt.Errorf("Wrong proxy status")
	}
	Log(Info, "Proxy %s initialization complete", addr)
	return nil
}

func (p *PeerToPeer) removeProxy(addr *net.UDPAddr) {
	for i, proxy := range p.Proxies {
		if proxy.Addr == addr {
			p.Proxies = append(p.Proxies[:i], p.Proxies[i+1:]...)
			return
		}
	}
}
