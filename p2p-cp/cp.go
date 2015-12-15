package main

// Control Peer

import (
	"bytes"
	bencode "github.com/jackpal/bencode-go"
	"github.com/wayn3h0/go-uuid"
	"log"
	"net"
	"p2p/commons"
	"time"
)

var (
	// List of all nodes registered to current DHT bootstrap node
	// This list should always be checked if item is unique by IP and hash
	NodeList []Node

	// Maximum number of nodes that can connect to DHT bootstrap node
	// When maximum number is exceeded, the most oldest entries should be wiped
	// from list. However, bootstrap node can deal with millions of entries
	// so possibility of MaximumNodes exceeding is fairly low
	MaximumNodes int
)

// Representation of a DHT Node that was connected to current DHT Bootstrap node
type Node struct {
	// Unique identifier in a form of UUID generated randomly upoc connection of a node
	ID string

	// IP Address of a node
	Endpoint string

	// Last time we pinged it.
	LastPing time.Time

	// Infohash that was associated with this node
	AssociatedHash string
}

// Infohash is a 20-bytes string and associated IP Address
// There must be multiple infohashes, but each infohash should
// have unique IP address, because we don't want to response
// multiple times with same IP for same infohash
type Infohash struct {
	// 20 bytes infohash string
	Hash string

	// Address associated with this hash
	NodeAddr string
}

// Router class
type DHTRouter struct {
	// Number of nodes participating in DHT
	NodesNumber int

	// Port which DHT router listens
	Port int

	// List of infohashes
	Hashes []Infohash
}

// Generate UUID, assigns it to a node and returns UUID as a string
// This methods always checks if generated ID is unique
func (node *Node) GenerateID(hashes []Infohash) string {
	var err error
	var id uuid.UUID
	id, err = uuid.NewTimeBased()
	if err != nil {
		log.Panic("[ERROR] Failed to generate UUID: %v", err)
		node.ID = ""
	} else {
		// Check if UUID is unique here
		var unique bool
		unique = true
		for _, hash := range hashes {
			if hash.Hash == id.String() {
				unique = false
			}
		}
		if unique {
			node.ID = id.String()
		} else {
			node.ID = node.GenerateID(hashes)
		}
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

// Allocated NodeList slice with maximum nodes
func AllocateNodeList() {
	log.Printf("[INFO] Allocating memory for %d nodes slice", MaximumNodes)
	NodeList = make([]Node, MaximumNodes)
}

// SetupServers prepares a DHT router listening socket that DHT clients
// will send UDP packets to
func (dht *DHTRouter) SetupServer() *net.UDPConn {
	log.Printf("[INFO] Setting UDP server at %d port", dht.Port)
	udp, err := net.ListenUDP("udp4", &net.UDPAddr{Port: dht.Port})
	if err != nil {
		log.Printf("[ERROR] Failed to start UDP Listener: %v", err)
		return nil
	}
	return udp
}

// IsNewPeer returns true if connected peer was not connected yes, false otherwise
func (dht *DHTRouter) IsNewPeer(addr string) bool {
	// TODO: Rewrite with use of ranges
	for i := 0; i < MaximumNodes; i++ {
		if NodeList[i].Endpoint == addr {
			return false
		}
	}
	return true
}

// Adds newly connected DHT node to a list of DHT participants
// New nodes not always added to the end of list as due to timeout
// some nodes may be wiped from the middle of the list. Therefore
// we go through full slice unless we find wiped node with empty ID
func (dht *DHTRouter) RegisterNode(n Node) {
	for i := 0; i < MaximumNodes; i++ {
		if NodeList[i].ID == "" {
			NodeList[i] = n
		}
	}
}

// Extracts DHTRequest from received packet
// This method tries to unmarshal bencode into DHTRequest structure
func (dht *DHTRouter) Extract(b []byte) (request commons.DHTRequest, err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("[ERROR] Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if err2 := bencode.Unmarshal(bytes.NewBuffer(b), &request); err2 == nil {
		err = nil
		return
	} else {
		log.Printf("[DEBUG] Received from peer: %v %q", request, err2)
		return request, err2
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

// EncodeResponse takes DHTResponse structure and turns it into bencode by
// Marshaling
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

// ResponseConn method generates a response to a "conn" network message received as a first packet
// from a newly connected node. Response writes an ID of the node
func (dht *DHTRouter) ResponseConn(req commons.DHTRequest, addr string, n Node) commons.DHTResponse {
	var resp commons.DHTResponse
	resp.Command = req.Command
	resp.Id = n.ID
	resp.Dest = "0"
	return resp
}

// ResponseFind method generates a response to a "find" network message which sent by DHT client
// when they want to build a p2p network based on infohash string.
// This method goes over list of hashes and collects information about all nodes with the
// same hash separated by comma
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

// ResponsePing responses to a received "ping" message
func (dht *DHTRouter) ResponsePing(req commons.DHTRequest, addr string) commons.DHTResponse {
	var resp commons.DHTResponse
	resp.Command = req.Command
	resp.Id = "0"
	resp.Dest = "0"
	return resp
}

// Send method send a packet to a connected client over network to a specific UDP address
func (dht *DHTRouter) Send(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	if msg != "" {
		_, err := conn.WriteToUDP([]byte(msg), addr)
		if err != nil {
			log.Printf("[ERROR] Failed to write to UDP: %v", err)
		}
	}
}

// This method listens to a UDP connections for incoming packets and
// sends generated responses back to DHT nodes
func (dht *DHTRouter) Listen(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		log.Printf("[ERROR] Failed to read from UDP socket: %v", err)
		return
	}
	var n Node
	if dht.IsNewPeer(addr.String()) {
		log.Printf("[INFO] New Peer connected: %s. Registering", addr)
		n.ID = n.GenerateID(dht.Hashes)
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
		newNode.GenerateID(dht.Hashes)
		newNode.LastPing = time.Now()
		NodeList[i] = newNode
	}
}
