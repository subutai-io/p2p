package ptp

/*

import (
	"net"
	"runtime"
	"sync"
)

type Proxy struct {
	Addr *net.UDPAddr
}

var Proxies map[string]Proxy
var ProxyLock sync.Mutex

func RequestProxies() {
	var c DHTClient
	config := dhtClient.DHTClientConfig()

}

func AddProxy(addr *net.UDPAddr) {
	ProxyLock.Lock()
	_, e := Proxies[addr.String()]
	ProxyLock.Unlock()
	runtime.Gosched()
	if e {
		ProxyLock.Unlock()
		runtime.Gosched()
		return
	}
	Proxies[addr.String()] = addr
	ProxyLock.Unlock()
	runtime.Gosched()
}

func RemoveProxy(addr *net.UDPAddr) {
	ProxyLock.Lock()
	key, e := Proxies[addr.String()]
	if !e {
		ProxyLock.Unlock()
		runtime.Gosched()
		return
	}
	delete(Proxies, e)

	ProxyLock.Unlock()
	runtime.Gosched()
}
*/
