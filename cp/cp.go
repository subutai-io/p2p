package main

import (
	"bytes"
	bencode "github.com/jackpal/bencode-go"
	"github.com/wayn3h0/go-uuid"
	"log"
	"net"
	"p2p/commons"
	"time"
)

var NodeList []Node
var MaximumNodes int

type Node struct {
	ID       string
	Endpoint string
	LastPing time.Time
}

type Infohash struct {
	Hash     string
	NodeAddr string
}

type DHTRouter struct {
	NodesNumber int
	Port        int
	Hashes      []Infohash
}

func CheckError(err error) {
	if err != nil {
		log.Panic("[ERROR] %v", err)
	}
}

// Generate UUID, assigns it to a node and returns UUID as a string
func (node *Node) GenerateID() string {
	var err error
	var id uuid.UUID
	id, err = uuid.NewTimeBased()
	if err != nil {
		log.Panic("[ERROR] Failed to generate UUID: %v", err)
		node.ID = ""
	} else {
		node.ID = id.String()
	}
	return node.ID
}

// Functions returns true if timeout period has passed since last ping
func (node *Node) isPingRequired(n *Node) bool {
	return false
}

// Currently unused
func handleConnection(c *net.Conn) int {
	return 1
}

func AllocateNodeList() {
	log.Printf("[INFO] Allocating memory for %d nodes slice", MaximumNodes)
	NodeList = make([]Node, MaximumNodes)
}

func (dht *DHTRouter) SetupServer() *net.UDPConn {
	log.Printf("[INFO] Setting UDP server at %d port", dht.Port)
	udp, err := net.ListenUDP("udp4", &net.UDPAddr{Port: dht.Port})
	CheckError(err)
	return udp
}

func (dht *DHTRouter) IsNewPeer(addr string) bool {
	for i := 0; i < MaximumNodes; i++ {
		if NodeList[i].Endpoint == addr {
			return false
		}
	}
	return true
}

func (dht *DHTRouter) RegisterNode(n Node) {
	for i := 0; i < MaximumNodes; i++ {
		if NodeList[i].ID == "" {
			NodeList[i] = n
		}
	}
}

// Extracts DHTRequest from received packet
func (dht *DHTRouter) Extract(b []byte) (request commons.DHTRequest, err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("[ERROR] Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if e2 := bencode.Unmarshal(bytes.NewBuffer(b), &request); e2 == nil {
		err = nil
		return
	} else {
		log.Printf("[DEBUG] Received from peer: %v %q", request, e2)
		return request, e2
	}
}

// Returns a bencoded representation of a DHTResponse
func (dht *DHTRouter) Compose(command, id, dest string) string {
	var resp commons.DHTResponse
	// Command is mandatory
	resp.Command = command
	// Defaults
	resp.Id = "0"
	resp.Dest = "0"
	if id != "" {
		resp.Id = id
	}
	if dest != "" {
		resp.Dest = dest
	}
	return dht.EncodeResponse(resp)
}

func (dht *DHTRouter) EncodeResponse(resp commons.DHTResponse) string {
	if resp.Command == "" {
		return ""
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, resp); err != nil {
		log.Printf("[ERROR] Failed to Marshal bencode %v", err)
		return ""
	}
	return b.String()
}

func (dht *DHTRouter) ResponseConn(req commons.DHTRequest, addr string, n Node) commons.DHTResponse {
	var resp commons.DHTResponse
	resp.Command = req.Command
	resp.Id = n.ID
	resp.Dest = "0"
	return resp
}

func (dht *DHTRouter) ResponseFind(req commons.DHTRequest, addr string) commons.DHTResponse {
	var resp commons.DHTResponse
	resp.Command = req.Command
	resp.Id = "0"
	var foundDest string
	// Comma separated list of found destinations
	for i := 0; i < len(dht.Hashes); i++ {
		if dht.Hashes[i].Hash == req.Hash {
			// We found required hash. Check if we're not the node who requested it
			if dht.Hashes[i].NodeAddr == addr {
				continue
			} else {
				foundDest += dht.Hashes[i].NodeAddr + ","
			}
		}
	}
	if foundDest == "" {
		// Save new hash
		var h Infohash
		h.Hash = req.Hash
		h.NodeAddr = addr
		dht.Hashes = append(dht.Hashes, h)
	}
	/*
		for hash := range dht.Hashes {
			if hash.Hash == req.Hash {
				// We found required hash. Check if we're not the node who requested it
				if hash.NodeAddr == addr {
					continue
				} else {
					foundDest += "," + hash.NodeAddr
				}
			}
		}
	*/
	resp.Dest = foundDest
	return resp
}

func (dht *DHTRouter) ResponsePing(req commons.DHTRequest, addr string) commons.DHTResponse {
	var resp commons.DHTResponse
	resp.Command = req.Command
	resp.Id = "0"
	resp.Dest = "0"
	return resp
}

func (dht *DHTRouter) Send(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	if msg != "" {
		_, err := conn.WriteToUDP([]byte(msg), addr)
		if err != nil {
			log.Printf("[ERROR] Failed to write to UDP: %v", err)
		}
	}

}

func (dht *DHTRouter) Listen(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	CheckError(err)
	var n Node
	if dht.IsNewPeer(addr.String()) {
		log.Printf("[INFO] New Peer connected: %s. Registering", addr)
		n.ID = n.GenerateID()
		n.Endpoint = addr.String()
		dht.RegisterNode(n)
	}
	log.Printf("[DEBUG] %s: %s", addr, string(buf[:512]))

	// Try to bencode
	req, err := dht.Extract(buf[:512])
	var resp commons.DHTResponse
	switch req.Command {
	case "conn":
		resp = dht.ResponseConn(req, addr.String(), n)
	case "find":
		resp = dht.ResponseFind(req, addr.String())
	case "ping":
		resp = dht.ResponsePing(req, addr.String())
	default:
		log.Printf("[ERROR] Unknown command received: %s", req.Command)
		resp.Command = ""
	}

	dht.Send(conn, addr, dht.EncodeResponse(resp))
}

func init() {
}

func main() {
	/*f, err := os.OpenFile("cp.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("[ERROR] Error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	*/
	MaximumNodes = 100
	AllocateNodeList()
	log.Printf("[INFO] Initialization complete")
	log.Printf("[INFO] Starting bootstrap node")
	var dht DHTRouter
	dht.Port = 6881
	listener := dht.SetupServer()

	for {
		dht.Listen(listener)
	}

	log.Printf("[INFO] Starting Control Peer")
	for i := 0; i < MaximumNodes; i++ {
		var newNode Node
		newNode.Endpoint = "IP"
		newNode.GenerateID()
		newNode.LastPing = time.Now()
		NodeList[i] = newNode
	}
}
