package main

import (
	"fmt"
	"net"
	"p2p/dht"
	"time"
)

type Proxy struct {
	DHTClient *dht.DHTClient
	Tunnels   map[int]Tunnel
}

// Tunnel established between two peers. Tunnels doesn't
// provide two-way connectivity.
type Tunnel struct {
	Src *net.UDPAddr
	Dst *net.UDPAddr
}

func (p *Proxy) Initialize() {
	p.DHTClient = new(dht.DHTClient)
	config := p.DHTClient.DHTClientConfig()
	config.NetworkHash = p.GenerateHash()
	//p.DHTClient.Initialize(config)
}

func (p *Proxy) GenerateHash() string {
	var infohash string
	t := time.Now()
	infohash = "cp" + fmt.Sprintf("%d%d%d", t.Year(), t.Month(), t.Day())
	return infohash
}
