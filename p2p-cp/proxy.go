package main

import (
	"p2p/dht"
	"time"
)

type Proxy struct {
	DHTClient *dht.DHTClient
}

func (p *Proxy) Initialize() {
	p.DHTClient = new(dht.DHTClient)
	config := p.DHTClient.DHTClientConfig()
	config.NetworkHash = p.GenerateHash()
}

func (p *Proxy) GenerateHash() string {
	var hash string
	t := time.Now()
	infohash += "cp" + t.Year() + t.Month() + t.Day()
	return infohash
}
