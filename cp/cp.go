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

func CheckError(err error) {
	if err != nil {
		log.Panic("[ERROR] %v", err)
	}
}

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

func handleConnection(c *net.Conn) int {
	return 1
}

func AllocateNodeList() {
	log.Printf("[INFO] Allocating memory for %d nodes slice", MaximumNodes)
	NodeList = make([]Node, MaximumNodes)
}

type DHTRouter struct {
	NodesNumber int
	Port        int
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

func (dht *DHTRouter) Extract(b []byte) (request commons.DHTRequest, err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Panicf("Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if e2 := bencode.Unmarshal(bytes.NewBuffer(b), &request); e2 == nil {
		err = nil
		return
	} else {
		log.Printf("Received from peer: %v %q", request, e2)
		return request, e2
	}

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
	resp.Dest = "0"
	return resp
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
	log.Printf("%s: %s", addr, string(buf[:512]))

	// Try to bencode
	req, err := dht.Extract(buf[:512])
	var resp commons.DHTResponse
	switch req.Command {
	case "conn":
		resp = dht.ResponseConn(req, addr.String(), n)
	case "find":
		resp = dht.ResponseFind(req, addr.String())
	default:
		log.Printf("[ERROR] Unknown command received: %s", req.Command)
	}

	var b bytes.Buffer
	if err := bencode.Marshal(&b, resp); err != nil {
		log.Printf("[ERROR] Failed to Marshal bencode %v", err)
		return
	}

	msg := b.String()
	_, err = conn.WriteToUDP([]byte(msg), addr)
	if err != nil {
		log.Printf("Failed to write to UDP: %v", err)
	}
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
	log.Printf("Starting...")
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
