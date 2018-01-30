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
	Addr       *net.UDPAddr // Address of the proxy
	Endpoint   *net.UDPAddr // Endpoint provided by proxy
	Status     proxyStatus  // Current proxy status
	LastUpdate time.Time
}

func (p *PeerToPeer) initProxy(addr string) error {
	var err error
	p.proxyLock.Lock()
	defer p.proxyLock.Unlock()

	pAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		Log(Error, "Failed to resolve proxy address")
		return fmt.Errorf("Failed to resolve proxy address")
	}
	for _, pr := range p.Proxies {
		if pr.Addr.String() == pAddr.String() {
			Log(Debug, "Proxy %s already exists", addr)
			return fmt.Errorf("Proxy %s already exists", addr)
		}
	}
	Log(Info, "Initializing proxy %s", addr)
	proxy := new(proxyServer)
	proxy.LastUpdate = time.Now()
	proxy.Addr = pAddr
	p.Proxies = append(p.Proxies, proxy)
	initStarted := time.Now()
	proxy.Status = proxyConnecting

	msg := CreateProxyP2PMessage(0, p.Dht.ID, 1)
	p.UDPSocket.SendMessage(msg, proxy.Addr)
	p.proxyLock.Unlock()
	for proxy.Status == proxyConnecting {
		time.Sleep(100 * time.Millisecond)
		if time.Duration(3*time.Second) < time.Since(initStarted) {
			p.disableProxy(proxy.Addr)
			Log(Error, "Failed to connect to proxy")
			return fmt.Errorf("Failed to connect to proxy")
		}
	}
	if proxy.Status != proxyActive {
		p.disableProxy(proxy.Addr)
		Log(Error, "Wrong proxy status")
		return fmt.Errorf("Wrong proxy status")
	}
	return nil
}

func (p *PeerToPeer) disableProxy(addr *net.UDPAddr) {
	for i, proxy := range p.Proxies {
		if proxy.Addr == addr {
			p.Proxies[i].Close()
			//p.Proxies = append(p.Proxies[:i], p.Proxies[i+1:]...)
			return
		}
	}
}

// Close will stop proxy
func (p *proxyServer) Close() error {
	Log(Info, "Stopping proxy %s, Endpoint: %s", p.Addr.String(), p.Endpoint.String())
	p.Addr = nil
	p.Endpoint = nil
	p.Status = proxyDisconnected
	return nil
}
