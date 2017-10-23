package ptp

import (
	"fmt"
	"net"
	"sync"
)

// ProxyStatus represents current status of a proxy
type ProxyStatus uint8

// Types of proxy statuses
const (
	ProxyConnecting ProxyStatus = 0
	ProxyConnected  ProxyStatus = 1
	ProxyFailed     ProxyStatus = 2
)

// ProxyList manages all proxies within daemon, not per p2p instance
type ProxyList struct {
	proxies map[string]*Proxy
	lock    sync.RWMutex
}

func (l *ProxyList) operate(action ListOperation, key string, proxy *Proxy) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if action == OperateUpdate {
		l.proxies[key] = proxy
	} else if action == OperateDelete {
		delete(l.proxies, key)
	}
}

// Exists returns whether proxy with specified already in the list or not
func (l *ProxyList) Exists(endpoint string) bool {
	_, exists := l.proxies[endpoint]
	return exists
}

// Add will add new proxy, connect to it and get our tunnel ID from it
func (l *ProxyList) Add(endpoint string) error {
	if l.Exists(endpoint) {
		return fmt.Errorf("Proxy already exists")
	}
	proxy := new(Proxy)
	err := proxy.Connect(endpoint)
	if err != nil {
		return fmt.Errorf("Failed to connect to proxy: %s", err)
	}
	// Wait for our ID
	return nil
}

// Proxy is user a traffic proxy when peer is behind NAT and can't
// connect to other peers in any way
type Proxy struct {
	Conn   *net.UDPConn
	ID     uint16
	Status ProxyStatus
}

// Connect will send initial handshake packet to the specified proxy
func (p *Proxy) Connect(endpoint string) error {

	return nil
}
