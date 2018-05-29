package ptp

import (
	"net"
	"testing"
)

func TestClose(t *testing.T) {
	d := new(proxyServer)
	d.Addr = new(net.UDPAddr)
	d.Addr.IP = []byte("192.168.34.2")
	d.Addr.Port = 8787
	d.Addr.Zone = "Zone"
	ips := "192.168.11.5"
	d.Endpoint, _ = net.ResolveUDPAddr("network", ips)
	d.Status = proxyActive

	d.Close()

	if d.Addr != nil && d.Status != 2 && d.Endpoint != nil {
		t.Error("Close Error")
	}
}

func TestOperate(t *testing.T) {
	p := new(ProxyManager)
	prsrv := new(proxyServer)
	p.init()
	p.proxies["1"] = prsrv
	oper := OperateUpdate
	p.operate(oper, "2", prsrv)
	for i := 0; i < len(p.proxies); i++ {
		_, exists := p.proxies["2"]
		if !exists {
			t.Error("Error in update operation")
		}
	}
	oper2 := OperateDelete
	p.operate(oper2, "1", prsrv)
	for i := 0; i < len(p.proxies); i++ {
		_, exist := p.proxies["1"]
		if exist {
			t.Error("Error in delete operation")
		}
	}
}

func TestNew(t *testing.T) {
	p := new(ProxyManager)
	prsrv1 := new(proxyServer)
	prsrv2 := new(proxyServer)
	prsrv1.Endpoint, _ = net.ResolveUDPAddr("24", "192.168.1.1")
	prsrv2.Endpoint, _ = net.ResolveUDPAddr("24", "192.168.1.2")
	p.init()
	p.proxies["1"] = prsrv1
	p.proxies["2"] = prsrv2

	endpoint, _ := net.ResolveUDPAddr("udp", "192.168.1.1")
	err := p.new(endpoint)
	if err != nil {
		t.Error("Error")
	}
}

func TestCheck(t *testing.T) {
	p := new(ProxyManager)
	prsrv1 := new(proxyServer)
	prsrv2 := new(proxyServer)
	prsrv1.Status = proxyConnecting
	prsrv2.Status = proxyActive
	p.init()
	p.proxies["10"] = prsrv1
	p.proxies["11"] = prsrv2
	p.check()
	_, exists := p.proxies["10"]
	if prsrv1.Addr != nil && prsrv1.Endpoint != nil && prsrv1.Status != proxyDisconnected && !exists {
		t.Error("Error")
	}
	if prsrv2.Addr != nil && prsrv1.Endpoint != nil && prsrv1.Status != proxyDisconnected && !exists {
		t.Error("Error")
	}
}

func TestTouch(t *testing.T) {
	p := new(ProxyManager)
	prsvr1 := new(proxyServer)
	prsvr2 := new(proxyServer)
	p.init()
	p.proxies["100"] = prsvr1
	p.proxies["101"] = prsvr2

	id1 := "100"
	get := p.touch(id1)
	if !get {
		t.Error("Error")
	}
	id2 := "0"
	get2 := p.touch(id2)
	if get2 {
		t.Error("Error. ProxyId is not exists")
	}
}

func TestActivate(t *testing.T) {
	p := new(ProxyManager)
	endpoint, _ := net.ResolveUDPAddr("24", "192.168.1.1")
	prsrv1 := new(proxyServer)
	prsrv1.Status = proxyConnecting
	prsrv2 := new(proxyServer)
	p.init()
	p.proxies["5"] = prsrv1
	p.proxies["6"] = prsrv2

	get := p.activate("6", endpoint)
	if !get {
		t.Error("Error")
	}
	get2 := p.activate("0", endpoint)
	if get2 {
		t.Error("Error")
	}
}
