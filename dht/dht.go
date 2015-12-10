package dht

import (
	"net"
	"strings"
	//"p2p/commons"
	"log"
)

type DHTClient struct {
	Routers    string ""
	Connection []*net.UDPConn
}

func DHTClientConfig() *DHTClient {
	return &DHTClient{
		Routers: "localhost:6881,dht1.subut.ai:6881,dht2.subut.ai:6881,dht3.subut.ai:6881",
	}
}

func (dht *DHTClient) ConnectAndHandshake(router string) (*net.UDPConn, error) {
	log.Printf("Connecting to a router %s", router)
	return nil, nil
}

func (dht *DHTClient) Initialize(config *DHTClient) {
	dht = config
	routers := strings.Split(dht.Routers, ",")
	for _, router := range routers {
		dht.ConnectAndHandshake(router)
	}
}
