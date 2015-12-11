package dht

import (
	"bytes"
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
}

func (dht *DHTClient) DHTClientConfig() *DHTClient {
	return &DHTClient{
		Routers:     "localhost:6881,dht1.subut.ai:6881,dht2.subut.ai:6881,dht3.subut.ai:6881",
		NetworkHash: "",
	}
}

func (dht *DHTClient) AddConnection(connections []*net.UDPConn, conn *net.UDPConn) []*net.UDPConn {
	n := len(connections)
	if n == cap(connections) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]*net.UDPConn, len(connections), 2*len(connections)+1)
		copy(newSlice, connections)
		connections = newSlice
	}
	connections = connections[0 : n+1]
	connections[n] = conn
	return connections
}

func (dht *DHTClient) ConnectAndHandshake(router string) (*net.UDPConn, error) {
	log.Printf("[DHT-INFO] Connecting to a router %s", router)
	addr, err := net.ResolveUDPAddr("udp", router)
	if err != nil {
		log.Printf("[DHT-ERROR]: Failed to resolve router address: %v", err)
		return nil, err
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Printf("[DHT-ERROR]: Failed to establish connection: %v", err)
		return nil, err
	}
	//defer conn.Close()

	dht.Connection = dht.AddConnection(dht.Connection, conn)

	// Handshake
	var req commons.DHTRequest
	req.Id = "0"
	req.Hash = "0"
	req.Command = "conn"
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

func (dht *DHTClient) ListenDHT(conn *net.UDPConn) {
	for {
		var buf [512]byte
		_, addr, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Printf("[DHT-ERROR] Failed to read from DHT bootstrap node: %v", err)
		} else {
			log.Printf("[DHT-DEBUG] %s", addr.String())
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
						log.Printf("[DHT-INFO] Received connection confirmation from tracker")
					}
				} else if data.Command == "ping" {
					msg := dht.Compose("ping", "", "")
					_, err = conn.Write([]byte(msg))
					if err != nil {
						log.Printf("[DHT-ERROR] Failed to send PING packet: %v", err)
					}
				} else if data.Command == "find" {
					log.Printf("[DHT-INFO] Found peers: %s", data.Dest)
				}
			}
		}
	}

	/*
		for i := 0; i < len(dht.Connection); i++ {
			var buf [512]byte
			_, addr, err := dht.Connection[i].ReadFromUDP(buf[0:])
			if err != nil {
				log.Printf("[DHT-ERROR] Failed to read from DHT bootstrap node: %v", err)
			} else {
				log.Printf("[DHT-DEBUG] %s", addr.String())
			}
		}
	*/
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
