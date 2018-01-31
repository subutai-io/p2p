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
