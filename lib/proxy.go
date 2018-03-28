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
	proxies    map[string]*proxyServer
	lock       sync.RWMutex
	hasChanges bool
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
		if proxy.Status == proxyActive && time.Since(proxy.LastUpdate) > time.Duration(90*time.Second) {
			err := proxy.Close()
			if err != nil {
				Log(Debug, "Failed to close proxy: %s", err)
			}
			Log(Debug, "Proxy %s has been disconnected by timeout", id)
		}
		if proxy.Status == proxyDisconnected {
			Log(Debug, "Removing proxy %s", id)
			p.operate(OperateDelete, id, nil)
			p.hasChanges = true
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
			p.hasChanges = true
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

// Close will stop proxy
func (p *proxyServer) Close() error {
	Log(Info, "Stopping proxy %s, Endpoint: %s", p.Addr.String(), p.Endpoint.String())
	p.Addr = nil
	p.Endpoint = nil
	p.Status = proxyDisconnected
	return nil
}
