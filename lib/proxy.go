package ptp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type proxyStatus uint8

const (
	proxyConnecting   proxyStatus = 0
	proxyActive       proxyStatus = 1
	proxyDisconnected proxyStatus = 2
)

// ProxyManager manages TURN servers
type ProxyManager struct {
	proxies map[string]*proxyServer
	lock    sync.RWMutex
}

func (p *ProxyManager) init() error {
	p.proxies = make(map[string]*proxyServer)
	return nil
}

func (p *ProxyManager) operate(operation ListOperation, addr string, proxy *proxyServer) {
	p.lock.Lock()
	defer p.lock.Unlock()
	if operation == OperateUpdate {
		p.proxies[addr] = proxy
	} else if operation == OperateDelete {
		_, exists := p.proxies[addr]
		if exists {
			delete(p.proxies, addr)
		}
	}
}

func (p *ProxyManager) get() map[string]*proxyServer {
	p.lock.RLock()
	result := make(map[string]*proxyServer)
	for i, v := range p.proxies {
		result[i] = v
	}
	p.lock.RUnlock()
	return result
}

// GetList will return a slice of proxyServers
func (p *ProxyManager) GetList() []*proxyServer {
	list := []*proxyServer{}
	proxies := p.get()
	for _, v := range proxies {
		list = append(list, v)
	}
	return list
}

func (p *ProxyManager) new(endpoint *net.UDPAddr) error {
	proxies := p.get()
	_, exists := proxies[endpoint.String()]
	if exists {
		return fmt.Errorf("Proxy %s already exists", endpoint.String())
	}
	proxy := new(proxyServer)
	proxy.Addr = endpoint
	proxy.Status = proxyConnecting
	proxy.Created = time.Now()
	p.operate(OperateUpdate, endpoint.String(), proxy)
	return nil
}

func (p *ProxyManager) check() {
	proxies := p.get()
	for id, proxy := range proxies {
		if proxy.Status == proxyConnecting && time.Since(proxy.Created) > time.Duration(10*time.Second) {
			err := proxy.Close()
			if err != nil {
				Log(Debug, "Failed to close proxy: %s", err)
			}
			Log(Debug, "Failed to connect to proxy %s", id)
		}
		if proxy.Status == proxyActive && time.Since(proxy.LastUpdate) > time.Duration(30*time.Second) {
			err := proxy.Close()
			if err != nil {
				Log(Debug, "Failed to close proxy: %s", err)
			}
			Log(Debug, "Proxy %s has been disconnected by timeout", id)
		}
		if proxy.Status == proxyDisconnected {
			Log(Debug, "Removing proxy %s", id)
			p.operate(OperateDelete, id, nil)
		}
	}
}

func (p *ProxyManager) touch(id string) bool {
	proxies := p.get()
	for pid, proxy := range proxies {
		if pid == id {
			proxy.LastUpdate = time.Now()
			p.operate(OperateUpdate, id, proxy)
			return true
		}
	}
	return false
}

func (p *ProxyManager) activate(id string, endpoint *net.UDPAddr) bool {
	proxies := p.get()
	for pid, proxy := range proxies {
		if pid == id && proxy.Status == proxyConnecting {
			proxy.Status = proxyActive
			proxy.LastUpdate = time.Now()
			proxy.Endpoint = endpoint
			p.operate(OperateUpdate, id, proxy)
			return true
		}
	}
	return false
}

type proxyServer struct {
	Addr       *net.UDPAddr // Address of the proxy
	Endpoint   *net.UDPAddr // Endpoint provided by proxy
	Status     proxyStatus  // Current proxy status
	LastUpdate time.Time    // Last ping
	Created    time.Time    // Creation timestamp
}

// func (p *PeerToPeer) initProxy(addr string) error {
// 	var err error
// 	p.proxyLock.Lock()
// 	defer p.proxyLock.Unlock()

// 	pAddr, err := net.ResolveUDPAddr("udp4", addr)
// 	if err != nil {
// 		Log(Error, "Failed to resolve proxy address")
// 		return fmt.Errorf("Failed to resolve proxy address")
// 	}
// 	for _, pr := range p.Proxies {
// 		if pr.Addr.String() == pAddr.String() {
// 			Log(Debug, "Proxy %s already exists", addr)
// 			return fmt.Errorf("Proxy %s already exists", addr)
// 		}
// 	}
// 	Log(Info, "Initializing proxy %s", addr)
// 	proxy := new(proxyServer)
// 	proxy.LastUpdate = time.Now()
// 	proxy.Addr = pAddr
// 	p.Proxies = append(p.Proxies, proxy)
// 	initStarted := time.Now()
// 	proxy.Status = proxyConnecting

// 	msg := CreateProxyP2PMessage(0, p.Dht.ID, 1)
// 	p.UDPSocket.SendMessage(msg, proxy.Addr)
// 	p.proxyLock.Unlock()
// 	for proxy.Status == proxyConnecting {
// 		time.Sleep(100 * time.Millisecond)
// 		if time.Duration(3*time.Second) < time.Since(initStarted) {
// 			p.disableProxy(proxy.Addr)
// 			Log(Error, "Failed to connect to proxy")
// 			return fmt.Errorf("Failed to connect to proxy")
// 		}
// 	}
// 	if proxy.Status != proxyActive {
// 		p.disableProxy(proxy.Addr)
// 		Log(Error, "Wrong proxy status")
// 		return fmt.Errorf("Wrong proxy status")
// 	}
// 	return nil
// }

// func (p *PeerToPeer) disableProxy(addr *net.UDPAddr) {
// 	for i, proxy := range p.Proxies {
// 		if proxy.Addr == addr {
// 			p.Proxies[i].Close()
// 			//p.Proxies = append(p.Proxies[:i], p.Proxies[i+1:]...)
// 			return
// 		}
// 	}
// }

// Close will stop proxy
func (p *proxyServer) Close() error {
	Log(Info, "Stopping proxy %s, Endpoint: %s", p.Addr.String(), p.Endpoint.String())
	p.Addr = nil
	p.Endpoint = nil
	p.Status = proxyDisconnected
	return nil
}
