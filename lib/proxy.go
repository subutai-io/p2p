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
	addr     *net.UDPAddr
	endpoint *net.UDPAddr // Endpoint provided by proxy
	status   proxyStatus
}

func (p *PeerToPeer) initProxy(addr string) error {
	var err error
	proxy := new(proxyServer)
	proxy.addr, err = net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return fmt.Errorf("Failed to resolve proxy address")
	}

	for _, pr := range p.Proxies {
		if pr.addr == proxy.addr {
			return fmt.Errorf("Proxy %s already exists", addr)
		}
	}
	p.Proxies = append(p.Proxies, proxy)
	initStarted := time.Now()
	proxy.status = proxyConnecting

	msg := CreateProxyP2PMessage(0, p.Dht.ID, 1)
	p.UDPSocket.SendMessage(msg, proxy.addr)
	for proxy.status == proxyConnecting {
		time.Sleep(100 * time.Millisecond)
		if time.Duration(3*time.Second) < time.Since(initStarted) {
			p.removeProxy(proxy.addr)
			return fmt.Errorf("Failed to connect to proxy")
		}
	}
	if proxy.status != proxyActive {
		p.removeProxy(proxy.addr)
		return fmt.Errorf("Wrong proxy status")
	}
	return nil
}

func (p *PeerToPeer) removeProxy(addr *net.UDPAddr) {
	for i, proxy := range p.Proxies {
		if proxy.addr == addr {
			p.Proxies = append(p.Proxies[:i], p.Proxies[i+1:]...)
			return
		}
	}
}
