package dht

import (
	"bytes"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	"github.com/subutai-io/p2p/commons"
	"github.com/subutai-io/p2p/go-stun/stun"
	log "github.com/subutai-io/p2p/p2p_log"
	"net"
	"os"
	"strings"
)

type OperatingMode int

const (
	MODE_CLIENT OperatingMode = 1
	MODE_CP     OperatingMode = 2
)

type DHTClient struct {
	Routers          string
	FailedRouters    []string
	Connection       []*net.UDPConn
	NetworkHash      string
	NetworkPeers     []string
	P2PPort          int
	LastCatch        []string
	ID               string
	Peers            []PeerIP
	Forwarders       []Forwarder
	ResponseHandlers map[string]DHTResponseCallback
	Mode             OperatingMode
	Shutdown         bool
}

type Forwarder struct {
	Addr          *net.UDPAddr
	DestinationID string
}

type PeerIP struct {
	ID  string
	Ips []string
}

type DHTResponseCallback func(data commons.DHTResponse, conn *net.UDPConn)

func (dht *DHTClient) DHTClientConfig() *DHTClient {
	return &DHTClient{
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
func (dht *DHTClient) ConnectAndHandshake(router string, ips []net.IP) (*net.UDPConn, error) {
	log.Log(log.INFO, "Connecting to a router %s", router)
	addr, err := net.ResolveUDPAddr("udp", router)
	if err != nil {
		log.Log(log.ERROR, "Failed to resolve discovery service address: %v", err)
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		log.Log(log.ERROR, "Failed to establish connection to discovery service: %v", err)
		return nil, err
	}

	log.Log(log.INFO, "Ready to peer discovery via %s [%s]", router, conn.RemoteAddr().String())

	// Handshake
	var req commons.DHTRequest
	req.Id = "0"
	req.Hash = "0"
	req.Command = commons.CMD_CONN
	// TODO: rename Port to something more clear
	req.Port = fmt.Sprintf("%d", dht.P2PPort)
	for _, ip := range ips {
		req.Port = req.Port + "|" + ip.String()
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		conn.Close()
		return nil, err
	}
	// TODO: Optimize types here
	msg := b.String()
	if dht.Shutdown {
		return nil, nil
	}
	_, err = conn.Write([]byte(msg))
	if err != nil {
		log.Log(log.ERROR, "Failed to send packet: %v", err)
		conn.Close()
		return nil, err
	}

	return conn, nil
}

// Extracts DHTRequest from received packet
func (dht *DHTClient) Extract(b []byte) (response commons.DHTResponse, err error) {
	defer func() {
		if x := recover(); x != nil {
			log.Log(log.ERROR, "Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if e2 := bencode.Unmarshal(bytes.NewBuffer(b), &response); e2 == nil {
		err = nil
		return
	} else {
		log.Log(log.DEBUG, "Received from peer: %v %q", response, e2)
		return response, e2
	}
}

// Returns a bencoded representation of a DHTRequest
func (dht *DHTClient) Compose(command, id, hash string, port string) string {
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
	req.Port = port
	return dht.EncodeRequest(req)
}

func (dht *DHTClient) EncodeRequest(req commons.DHTRequest) string {
	if req.Command == "" {
		return ""
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		return ""
	}
	return b.String()
}

// After receiving a list of peers from DHT we will parse the list
// and add every new peer into list of peers
func (dht *DHTClient) UpdateLastCatch(catch string) {
	peers := strings.Split(catch, ",")
	for _, p := range peers {
		if p == "" {
			continue
		}
		var found bool = false
		for _, catchedPeer := range dht.LastCatch {
			if p == catchedPeer {
				found = true
			}
		}
		if !found {
			dht.LastCatch = append(dht.LastCatch, p)
		}
	}
}

// This function sends a request to DHT bootstrap node with ID of
// target node we want to connect to
func (dht *DHTClient) RequestPeerIPs(id string) {
	msg := dht.Compose(commons.CMD_NODE, id, "", "")
	for _, conn := range dht.Connection {
		if dht.Shutdown {
			continue
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Log(log.ERROR, "Failed to send 'node' request to %s: %v", conn.RemoteAddr().String(), err)
		}
	}
}

// UpdatePeers sends "find" request to a DHT Bootstrap node, so it can respond
// with a list of peers that we can connect to
// This method should be called periodically in case any new peers was discovered
func (dht *DHTClient) UpdatePeers() {
	msg := dht.Compose(commons.CMD_FIND, "", dht.NetworkHash, "")
	for _, conn := range dht.Connection {
		if dht.Shutdown {
			continue
		}
		log.Log(log.TRACE, "Updating peer %s", conn.RemoteAddr().String())
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Log(log.ERROR, "Failed to send 'find' request to %s: %v", conn.RemoteAddr().String(), err)
		}
	}
}

// Listens for packets received from DHT bootstrap node
// Every packet is unmarshaled and turned into Request structure
// which we should analyze and respond
func (dht *DHTClient) ListenDHT(conn *net.UDPConn) string {
	log.Log(log.INFO, "Bootstraping via %s", conn.RemoteAddr().String())
	for {
		if dht.Shutdown {
			log.Log(log.INFO, "Closing DHT Connection to %s", conn.RemoteAddr().String())
			conn.Close()
			for i, c := range dht.Connection {
				if c.RemoteAddr().String() == conn.RemoteAddr().String() {
					dht.Connection = append(dht.Connection[:i], dht.Connection[i+1:]...)
				}
			}
			break
		}
		var buf [512]byte
		_, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			log.Log(log.ERROR, "Failed to read from Discovery Service: %v", err)
		} else {
			data, err := dht.Extract(buf[:512])
			if err != nil {
				log.Log(log.ERROR, "Failed to extract a message received from discovery service: %v", err)
			} else {
				callback, exists := dht.ResponseHandlers[data.Command]
				if exists {
					callback(data, conn)
				} else {
					log.Log(log.ERROR, "Unknown packet received from DHT: %s", data.Command)
				}
			}
		}
	}
	return ""
}

func (dht *DHTClient) HandleConn(data commons.DHTResponse, conn *net.UDPConn) {
	log.Log(log.DEBUG, "CONN packet receied")
	if dht.ID != "" {
		log.Log(log.ERROR, "Empty ID was received")
		return
	}
	dht.ID = data.Id
	if dht.ID == "0" {
		log.Log(log.ERROR, "Empty ID were received. Stopping")
		os.Exit(1)
	}
	// Send a hash within FIND command
	// Afterwards application should wait for response from DHT
	// with list of clients. This may not happen if this client is the
	// first connected node.
	msg := dht.Compose(commons.CMD_FIND, "", dht.NetworkHash, "")
	if dht.Shutdown {
		return
	}
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Log(log.ERROR, "Failed to send 'find' request: %v", err)
	} else {
		log.Log(log.INFO, "Received connection confirmation from router %s",
			conn.RemoteAddr().String())
		log.Log(log.INFO, "Received personal ID for this session: %s", data.Id)
	}
}

func (dht *DHTClient) HandlePing(data commons.DHTResponse, conn *net.UDPConn) {
	msg := dht.Compose(commons.CMD_PING, "", "", "")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		log.Log(log.ERROR, "Failed to send 'ping' packet: %v", err)
	}
}

func (dht *DHTClient) HandleFind(data commons.DHTResponse, conn *net.UDPConn) {
	// This means we've received a list of nodes we can connect to
	if data.Dest != "" {
		ids := strings.Split(data.Dest, ",")
		if len(ids) == 0 {
			log.Log(log.ERROR, "Malformed list of peers received")
		} else {
			// Go over list of received peer IDs and look if we know
			// anything about them. Add every new peer into list of peers
			for _, id := range ids {
				var found bool = false
				for _, peer := range dht.Peers {
					if peer.ID == id {
						found = true
					}
				}
				if !found {
					var p PeerIP
					p.ID = id
					dht.Peers = append(dht.Peers, p)
				}
			}
			for i, peer := range dht.Peers {
				var found bool = false
				for _, id := range ids {
					if peer.ID == id {
						found = true
					}
				}
				if !found {
					log.Log(log.INFO, "Removing")
					dht.Peers = append(dht.Peers[:i], dht.Peers[i+1:]...)
				}
			}
			log.Log(log.DEBUG, "Received peers from %s: %s", conn.RemoteAddr().String(), data.Dest)
			dht.UpdateLastCatch(data.Dest)
		}
	} else {
		dht.Peers = dht.Peers[:0]
	}
}

func (dht *DHTClient) HandleRegCp(data commons.DHTResponse, conn *net.UDPConn) {
	log.Log(log.INFO, "Control peer has been registered in Service Discovery Peer")
	// We've received a registration confirmation message from DHT bootstrap node
}

func (dht *DHTClient) HandleNode(data commons.DHTResponse, conn *net.UDPConn) {
	// We've received an IPs associated with target node
	for i, peer := range dht.Peers {
		if peer.ID == data.Id {
			ips := strings.Split(data.Dest, "|")
			dht.Peers[i].Ips = ips
		}
	}
}

func (dht *DHTClient) HandleCp(data commons.DHTResponse, conn *net.UDPConn) {
	// We've received information about proxy
	log.Log(log.INFO, "Received control peer %s. Saving", data.Dest)
	var found bool = false
	for _, fwd := range dht.Forwarders {
		if fwd.Addr.String() == data.Dest && fwd.DestinationID == data.Id {
			found = true
		}
	}
	if !found {
		var fwd Forwarder
		a, err := net.ResolveUDPAddr("udp", data.Dest)
		if err != nil {
			log.Log(log.ERROR, "Failed to resolve UDP Address for proxy %s", data.Dest)
		} else {
			fwd.Addr = a
			fwd.DestinationID = data.Id
			dht.Forwarders = append(dht.Forwarders, fwd)
			log.Log(log.DEBUG, "Control peer has been added to the list of forwarders")
			log.Log(log.DEBUG, "Sending notify request back to the DHT")
			msg := dht.Compose(commons.CMD_NOTIFY, "", dht.ID, data.Id)
			for _, conn := range dht.Connection {
				if dht.Shutdown {
					continue
				}
				_, err := conn.Write([]byte(msg))
				if err != nil {
					log.Log(log.ERROR, "Failed to send 'node' request to %s: %v", conn.RemoteAddr().String(), err)
				}
			}
		}
	}
}

func (dht *DHTClient) HandleNotify(data commons.DHTResponse, conn *net.UDPConn) {
	// Notify means we should ask DHT bootstrap node for a control peer
	// in order to connect to a node that can't reach us
	dht.RequestControlPeer(data.Id)
}

func (dht *DHTClient) HandleStop(data commons.DHTResponse, conn *net.UDPConn) {
	conn.Close()
}

// This method initializes DHT by splitting list of routers and connect to each one
func (dht *DHTClient) Initialize(config *DHTClient, ips []net.IP) *DHTClient {
	dht = config
	routers := strings.Split(dht.Routers, ",")
	dht.FailedRouters = make([]string, len(routers))
	dht.ResponseHandlers = make(map[string]DHTResponseCallback)
	if dht.Mode != MODE_CP && dht.Mode != MODE_CLIENT {
		dht.Mode = MODE_CLIENT
	}
	if dht.Mode == MODE_CLIENT {
		log.Log(log.INFO, "DHT operating in CLIENT mode")
		dht.ResponseHandlers[commons.CMD_NODE] = dht.HandleNode
		dht.ResponseHandlers[commons.CMD_CP] = dht.HandleCp
		dht.ResponseHandlers[commons.CMD_NOTIFY] = dht.HandleNotify
	} else {
		log.Log(log.INFO, "DHT operating in CONTROL PEER mode")
		dht.ResponseHandlers[commons.CMD_REGCP] = dht.HandleRegCp
	}
	dht.ResponseHandlers[commons.CMD_FIND] = dht.HandleFind
	dht.ResponseHandlers[commons.CMD_CONN] = dht.HandleConn
	dht.ResponseHandlers[commons.CMD_PING] = dht.HandlePing
	dht.ResponseHandlers[commons.CMD_STOP] = dht.HandleStop
	for _, router := range routers {
		conn, err := dht.ConnectAndHandshake(router, ips)
		if err != nil || conn == nil {
			log.Log(log.ERROR, "Failed to handshake with a DHT Server: %v", err)
			dht.FailedRouters[0] = router
		} else {
			log.Log(log.INFO, "Handshaked. Starting listener")
			dht.Connection = append(dht.Connection, conn)
			go dht.ListenDHT(conn)
		}
	}
	return dht
}

var nat_type_str = [...]string{"NAT_ERROR", "NAT_UNKNOWN", "NAT_NONE", "NAT_BLOCKED",
	"NAT_FULL", "NAT_SYMETRIC", "NAT_RESTRICTED", "NAT_PORT_RESTRICTED", "NAT_SYMETRIC_UDP_FIREWALL"}

func DetectIP() string {
	stun_client := stun.NewClient()
	stun_client.LocalAddr = ""
	stun_client.LocalPort = 15000
	stun_client.SetSoftwareName("subutai")
	stun_client.SetServerHost("stun.iptel.org", 3478)
	_, host, err := stun_client.Discover()
	if err != nil {
		log.Log(log.ERROR, "Stun discover error : %v", err)
		return ""
	}
	if host != nil {
		return host.IP()
	}
	return ""
}

// This method register control peer on a Bootstrap node
func (dht *DHTClient) RegisterControlPeer() {
	var req commons.DHTRequest
	var err error
	req.Id = dht.ID
	req.Hash = "0"
	req.Command = commons.CMD_REGCP
	req.Port = fmt.Sprintf("%d", dht.P2PPort)
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		return
	}
	// TODO: Optimize types here
	msg := b.String()
	for _, conn := range dht.Connection {
		if dht.Shutdown {
			continue
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Log(log.ERROR, "Failed to send packet: %v", err)
			conn.Close()
			return
		}
	}
}

// This method request a new control peer for particular host
func (dht *DHTClient) RequestControlPeer(id string) {
	var req commons.DHTRequest
	var err error
	req.Id = dht.ID
	req.Hash = dht.NetworkHash
	req.Command = commons.CMD_CP
	req.Port = id
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		return
	}
	msg := b.String()
	// TODO: Move sending to a separate method
	for _, conn := range dht.Connection {
		if dht.Shutdown {
			continue
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			log.Log(log.ERROR, "Failed to send packet: %v", err)
			conn.Close()
			return
		}
	}
}

func (dht *DHTClient) ReportControlPeerLoad(amount int) {
	var req commons.DHTRequest
	req.Id = dht.ID
	req.Command = commons.CMD_LOAD
	req.Port = fmt.Sprintf("%d", amount)
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		return
	}
	msg := b.String()
	// TODO: Move sending to a separate method
	for _, conn := range dht.Connection {
		if dht.Shutdown {
			continue
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			log.Log(log.ERROR, "Failed to send packet: %v", err)
			conn.Close()
			return
		}
	}
}

func (dht *DHTClient) Stop() {
	dht.Shutdown = true
	var req commons.DHTRequest
	req.Id = dht.ID
	req.Command = commons.CMD_STOP
	req.Port = "0"
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		log.Log(log.ERROR, "Failed to Marshal bencode %v", err)
		return
	}
	msg := b.String()
	for _, conn := range dht.Connection {
		conn.Write([]byte(msg))
	}
}
