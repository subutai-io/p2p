package dht

import (
	"bytes"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"log"
	"net"
	"p2p/commons"
	"strings"
)

type DHTClient struct {
	Routers       string
	FailedRouters []string
	Connection    []*net.UDPConn
	NetworkHash   string
	NetworkPeers  []string
	P2PPort       int
}

func (dht *DHTClient) DHTClientConfig() *DHTClient {
	return &DHTClient{
		//Routers: "localhost:6881",
		Routers: "dht1.subut.ai:6881",
		//Routers:     "dht1.subut.ai:6881,dht2.subut.ai:6881,dht3.subut.ai:6881,dht4.subut.ai:6881,dht5.subut.ai:6881",
		NetworkHash: "",
	}
}

// AddConnection adds new UDP Connection reference onto list of DHT node connections
func (dht *DHTClient) AddConnection(connections []*net.UDPConn, conn *net.UDPConn) []*net.UDPConn {
	n := len(connections)
	if n == cap(connections) {
		newSlice := make([]*net.UDPConn, len(connections), 2*len(connections)+1)
		copy(newSlice, connections)
		connections = newSlice
	}
	connections = connections[0 : n+1]
	connections[n] = conn
	return connections
}

// ConnectAndHandshake sends an initial packet to a DHT bootstrap node
func (dht *DHTClient) ConnectAndHandshake(router string) (*net.UDPConn, error) {
	log.Printf("[DHT-INFO] Connecting to a router %s", router)
	addr, err := net.ResolveUDPAddr("udp", router)
	if err != nil {
		log.Printf("[DHT-ERROR]: Failed to resolve router address: %v", err)
		return nil, err
	}
	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Printf("[DHT-ERROR]: Failed to establish connection: %v", err)
		return nil, err
	}

	log.Printf("[DHT-INFO] Ready to bootstrap with %s [%s]", router, conn.RemoteAddr().String())
	dht.Connection = dht.AddConnection(dht.Connection, conn)

	// Handshake
	var req commons.DHTRequest
	req.Id = "0"
	req.Hash = "0"
	req.Command = "conn"
	req.Port = fmt.Sprintf("%d", dht.P2PPort)
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Printf("[DHT-ERROR] Failed to Marshal bencode %v", err)
		conn.Close()
		return nil, err
	}
	// TODO: Optimize types here
	msg := b.String()
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Printf("[DHT-ERROR] Failed to send packet: %v", err)
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// Extracts DHTRequest from received packet
func (dht *DHTClient) Extract(b []byte) (response commons.DHTResponse, err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("[DHT-ERROR] Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if e2 := bencode.Unmarshal(bytes.NewBuffer(b), &response); e2 == nil {
		err = nil
		return
	} else {
		log.Printf("[DHT-DEBUG] Received from peer: %v %q", response, e2)
		return response, e2
	}
}

// Returns a bencoded representation of a DHTRequest
func (dht *DHTClient) Compose(command, id, hash string) string {
	var req commons.DHTRequest
	// Command is mandatory
	req.Command = command
	// Defaults
	req.Id = "0"
	req.Hash = "0"
	if id != "" {
		req.Id = id
	}
	if hash != "" {
		req.Hash = hash
	}
	return dht.EncodeRequest(req)
}

func (dht *DHTClient) EncodeRequest(req commons.DHTRequest) string {
	if req.Command == "" {
		return ""
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Printf("[ERROR] Failed to Marshal bencode %v", err)
		return ""
	}
	return b.String()
}

func (dht *DHTClient) ListenDHT(conn *net.UDPConn) string {
	log.Printf("[DHT-INFO] Bootstraping via %s", conn.RemoteAddr().String())
	for {
		var buf [512]byte
		//_, addr, err := conn.ReadFromUDP(buf[0:])
		_, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Printf("[DHT-ERROR] Failed to read from DHT bootstrap node: %v", err)
		} else {
			data, err := dht.Extract(buf[:512])
			if err != nil {
				log.Printf("[DHT-ERROR] Failed to extract a message: %v", err)
			} else {
				if data.Command == "conn" {
					// Send a hash
					msg := dht.Compose("find", "", dht.NetworkHash)
					_, err = conn.Write([]byte(msg))
					if err != nil {
						log.Printf("[DHT-ERROR] Failed to send FIND packet: %v", err)
					} else {
						log.Printf("[DHT-INFO] Received connection confirmation from tracker %s", conn.RemoteAddr().String())
					}
				} else if data.Command == "ping" {
					msg := dht.Compose("ping", "", "")
					_, err = conn.Write([]byte(msg))
					if err != nil {
						log.Printf("[DHT-ERROR] Failed to send PING packet: %v", err)
					}
				} else if data.Command == "find" {
					log.Printf("[DHT-INFO] Found peers from %s: %s", conn.RemoteAddr().String(), data.Dest)
					/*
						hosts := strings.Split(data.Dest, ",")
						var hostExists bool
						for _, host := range hosts {
							if host == "" {
								continue
							}
							hostExists = false
							for _, ehost := range dht.NetworkPeers {
								if host != ehost {
									continue
								}
								hostExists = true
							}
							if !hostExists {
								dht.NetworkPeers = append(dht.NetworkPeers, host)
							}
						}
					*/
				}
			}
		}
	}
}

func (dht *DHTClient) Initialize(config *DHTClient) {
	dht = config
	routers := strings.Split(dht.Routers, ",")
	dht.FailedRouters = make([]string, len(routers))
	for _, router := range routers {
		conn, err := dht.ConnectAndHandshake(router)
		if err != nil || conn == nil {
			dht.FailedRouters[0] = router
		} else {
			go dht.ListenDHT(conn)
		}
	}
}
