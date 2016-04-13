package main

// Control Peer and DHT Bootstrap Node

import (
	"bytes"
	"flag"
	"fmt"
	bencode "github.com/jackpal/bencode-go"
	ptp "github.com/subutai-io/p2p/lib"
	"github.com/wayn3h0/go-uuid"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	// List of all nodes registered to current DHT bootstrap node
	// This list should always be checked if item is unique by IP and hash
	//PeerList map[string]Peer

	// Ping timeout for variables
	PingTimeout time.Duration = 3 * time.Second
)

type DHTState int
type DHTType int

const (
	ST_RUN      DHTState = 1
	ST_SHUTDOWN DHTState = 2
	T_BOOTSTRAP DHTType  = 1
	T_NORMAL    DHTType  = 2
)

type DHTPeer struct {
	Address  string
	Socket   *net.UDPConn
	PeersNum int
	State    DHTState
	Type     DHTType
	IP       net.IP
	Network  *net.IPNet
}

type DHTCallback func(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error)

// Representation of a DHT Node that was connected to current DHT Bootstrap node
type Peer struct {
	// Unique identifier in a form of UUID generated randomly upoc connection of a node
	ID string

	// IP Address of a node that is listening for incoming connections
	// from future network participants
	Endpoint string

	// Address that was received during connection
	ConnectionAddress string

	// Last time we pinged it.
	LastPing time.Time

	// Infohash that was associated with this node
	AssociatedHash string

	Addr *net.UDPAddr

	MissedPing int

	// When disabled - node will not be interracted.
	Disabled bool

	IPList  []*net.UDPAddr
	IP      net.IP
	Network *net.IPNet
}

// Control Peer represents a connected control peer that can be used by
// normal peers to forward their traffic
type ControlPeer struct {
	ID        string
	Addr      *net.UDPAddr
	TunelsNum int
}

// Infohash is a 20-bytes string and associated IP Address
// There must be multiple infohashes, but each infohash should
// have unique IP address, because we don't want to response
// multiple times with same IP for same infohash
type Infohash struct {
	// 20 bytes infohash string
	Hash string

	// List of Proxies for this hash
	Proxies []string
}

type DHCPSet struct {
	Hash    string
	Network *net.IPNet
}

// Router class
type DHTRouter struct {
	// Number of nodes participating in DHT
	NodesNumber int

	// Port which DHT router listens
	Port int

	// List of infohashes
	Hashes map[string]Infohash

	Connection *net.UDPConn

	ControlPeers []ControlPeer

	PeerList map[string]Peer

	Callbacks map[string]DHTCallback

	DHCPLock bool

	DHCPTable map[string]DHCPSet
	Lock      sync.Mutex
}

// Method ValidateConnection() tries to establish connection with control
// peer to check is it's accessible from outside.
// Return true if CP is able to received connection, false otherwise
func (cp *ControlPeer) ValidateConnection() bool {
	conn, err := net.DialUDP("udp", nil, cp.Addr)
	if err != nil {
		ptp.Log(ptp.ERROR, "Validation failed")
		return false
	}
	// TODO: Send something to CP
	err = conn.Close()
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to close connection with control peer: %v", err)
	}
	return true
}

// Generate UUID, assigns it to a node and returns UUID as a string
// This methods always checks if generated ID is unique
func (node *Peer) GenerateID(hashes map[string]Infohash) string {
	var err error
	var id uuid.UUID
	id, err = uuid.NewTimeBased()

	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to generate UUID: %v", err)
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
func (node *Peer) isPingRequired(n *Peer) bool {
	return false
}

// Currently unused
func handleConnection(c *net.Conn) int {
	return 1
}

// SetupServers prepares a DHT router listening socket that DHT clients
// will send UDP packets to
func (dht *DHTRouter) SetupServer() *net.UDPConn {
	ptp.Log(ptp.INFO, "Setting UDP server at %d port", dht.Port)
	udp, err := net.ListenUDP("udp4", &net.UDPAddr{Port: dht.Port})
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to start UDP Listener: %v", err)
		return nil
	}
	return udp
}

// IsNewPeer returns true if connected peer was not connected yes, false otherwise
func (dht *DHTRouter) IsNewPeer(addr string) bool {
	// TODO: Rewrite with use of ranges
	for _, peer := range dht.PeerList {
		if peer.ConnectionAddress == addr {
			return false
		}
	}
	return true
}

// Extracts DHTMessage from received packet
// This method tries to unmarshal bencode into DHTMessage structure
func (dht *DHTRouter) Extract(b []byte) (request ptp.DHTMessage, err error) {
	defer func() {
		if x := recover(); x != nil {
			ptp.Log(ptp.ERROR, "Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	if err2 := bencode.Unmarshal(bytes.NewBuffer(b), &request); err2 == nil {
		err = nil
		return
	} else {
		ptp.Log(ptp.DEBUG, "Received from peer: %v %q", request, err2)
		return request, err2
	}
}

// Returns a bencoded representation of a DHTMessage
func (dht *DHTRouter) Compose(command, id, dest string) string {
	var resp ptp.DHTMessage
	// Command is mandatory
	resp.Command = command
	// Defaults
	resp.Id = "0"
	resp.Arguments = "0"
	if id != "" {
		resp.Id = id
	}
	if dest != "" {
		resp.Arguments = dest
	}
	return dht.EncodeResponse(resp)
}

// EncodeResponse takes DHTMessage structure and turns it into bencode by
// Marshaling
func (dht *DHTRouter) EncodeResponse(resp ptp.DHTMessage) string {
	if resp.Command == "" {
		return ""
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, resp); err != nil {
		ptp.Log(ptp.ERROR, "Failed to Marshal bencode %v", err)
		return ""
	}
	return b.String()
}

func (dht *DHTRouter) HandleConn(req ptp.DHTMessage, addr *net.UDPAddr, p *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	resp.Command = req.Command
	resp.Id = "0"
	resp.Arguments = "0"
	var supported bool = false

	// Check that current version is supported
	for _, ver := range ptp.SUPPORTED_VERSIONS {
		if ver == req.Query {
			supported = true
		}
	}
	if !supported {
		ptp.Log(ptp.DEBUG, "Unsupported packet version received during connection from %s", addr.String())
		for i, peer := range dht.PeerList {
			if peer.Addr.String() == addr.String() {
				peer.Disabled = true
				dht.PeerList[i] = peer
			}
		}
		return resp, ptp.ErrorList[ptp.ERR_INCOPATIBLE_VERSION]
	}

	// First element should always be a port number
	data := strings.Split(req.Arguments, "|")
	if len(data) <= 1 {
		// We should receive information about at least one network interface
		ptp.Log(ptp.ERROR, "DHT Received malformed handshake")
		return resp, ptp.ErrorList[ptp.ERR_MALFORMED_HANDSHAKE]
	}

	port, err := strconv.Atoi(data[0])
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to parse port from handshake packet")
		return resp, ptp.ErrorList[ptp.ERR_PORT_PARSE_FAILED]
	}

	var ipList []*net.UDPAddr

	for i, d := range data {
		if i == 0 {
			// Put global IP address first
			dIp, _, _ := net.SplitHostPort(addr.String())
			a, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dIp, port))
			if err != nil {
				ptp.Log(ptp.ERROR, "Failed to resolve UDP address during handshake: %v", err)
				return resp, ptp.ErrorList[ptp.ERR_BAD_UDP_ADDR]
			}
			ipList = append(ipList, a)
			continue
		}
		if d == "" {
			continue
		}
		udpAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", d, port))
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to resolve address during handshake: %v", err)
			continue
		}
		var found bool = false
		for _, ip := range ipList {
			if ip.String() == udpAddr.String() {
				// Sometimes when interface IP address is equal to global IP address they will duplicate
				found = true
			}
		}
		if !found {
			ipList = append(ipList, udpAddr)
		}
	}

	for i, peer := range dht.PeerList {
		if peer.ConnectionAddress == addr.String() {
			peer.IPList = ipList
			dht.PeerList[i] = peer
		}
	}

	resp.Id = p.ID
	ptp.Log(ptp.INFO, "Sending greeting with ID %s to %s", p.ID, addr)
	return resp, nil
}

// SendFind sends FIND packet to every peer under specified hash
// to notify about other network participants
func (dht *DHTRouter) SendFind(hash string) {
	ptp.Log(ptp.DEBUG, "Changes in peer list. Notifying everyone")
	var ids []string
	for _, peer := range dht.PeerList {
		if peer.AssociatedHash == hash {
			ids = append(ids, peer.ID)
		}
	}
	for _, id := range ids {
		var list string
		for _, k := range ids {
			if k == id {
				continue
			}
			list = list + k + ","
		}
		var resp ptp.DHTMessage
		resp.Id = id
		resp.Command = ptp.CMD_FIND
		resp.Arguments = list
		dht.Send(dht.Connection, dht.PeerList[id].Addr, dht.EncodeResponse(resp))
	}
}

func (dht *DHTRouter) SendStop(hash, id string) {
	for _, peer := range dht.PeerList {
		if peer.AssociatedHash == hash {
			var resp ptp.DHTMessage
			resp.Command = ptp.CMD_STOP
			resp.Id = peer.ID
			resp.Arguments = id
			dht.Send(dht.Connection, peer.Addr, dht.EncodeResponse(resp))
		}
	}
}

func (dht *DHTRouter) FindFreeProxies() []string {
	var maxProxyNum int = 1
	var proxyNum int = 0
	var result []string
	if len(dht.ControlPeers) > 3 {
		maxProxyNum = 3
	}
	for _, proxy := range dht.ControlPeers {
		if proxyNum >= maxProxyNum {
			break
		}
		result = append(result, proxy.Addr.String())
		proxyNum = proxyNum + 1
	}
	return result
}

func (dht *DHTRouter) RegisterHash(addr string, hash string) {
	for i, peer := range dht.PeerList {
		if peer.ConnectionAddress == addr {
			peer.AssociatedHash = hash
			dht.PeerList[i] = peer
			ptp.Log(ptp.DEBUG, "Registering hash '%s' for %s", hash, addr)
			_, exists := dht.Hashes[hash]
			if !exists {
				var newHash Infohash
				newHash.Hash = hash
				newHash.Proxies = dht.FindFreeProxies()
				dht.Hashes[hash] = newHash
			}
			go dht.SendFind(hash)
		}
	}
}

func (dht *DHTRouter) PeerExists(id string) bool {
	_, exists := dht.PeerList[id]
	return exists
}

func (dht *DHTRouter) HandleFind(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	if len(req.Id) != 36 {
		ptp.Log(ptp.DEBUG, "Malformed ID received. Ignoring")
		return resp, ptp.ErrorList[ptp.ERR_BAD_ID_RECEIVED]
	}
	var foundDest string
	var hashExists bool = false
	for _, node := range dht.PeerList {
		if node.AssociatedHash == req.Query {
			if node.ConnectionAddress == addr.String() {
				hashExists = true
				// Skip if we are the node who requested hash
				continue
			}
			ptp.Log(ptp.TRACE, "Found match in hash '%s' with peer %s", req.Query, node.AssociatedHash)
			foundDest += node.ID + ","
		}
	}
	if !hashExists {
		// Hash was not found for current node. Add it
		dht.RegisterHash(addr.String(), req.Query)
	}
	resp.Command = req.Command
	resp.Id = "0"
	resp.Arguments = foundDest
	return resp, nil
}

func (dht *DHTRouter) HandlePing(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	resp.Command = ""
	peer.MissedPing = 0
	dht.PeerList[req.Id] = *peer
	return resp, nil
}

// ResponseRegCP will check newly connected CP if it was not connected before. Also,
// this method will call a function that will try to connect to CP to see if it's
// accessible from outside it's network and not blocked by NAT, so normal peers
// can connect to it
//func (dht *DHTRouter) ResponseRegCP(req ptp.DHTMessage, addr string) ptp.DHTMessage {
func (dht *DHTRouter) HandleRegCp(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	resp.Command = req.Command
	resp.Id = "0"
	resp.Arguments = "0"
	laddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to extract CP address: %v", err)
		resp.Command = ""
	} else {
		var isNew bool = true
		for _, cp := range dht.ControlPeers {
			if cp.ID == req.Id {
				isNew = false
			}
		}
		if !isNew {
			// At this point we will send an empty response, so CP will try
			// to reconnect later, when it's previous instance will be wiped
			// from list after PING timeout
			ptp.Log(ptp.ERROR, "Connected control peer is already in list")
			resp.Command = ""
		} else {
			var newCP ControlPeer
			newCP.ID = req.Id
			addrStr := laddr.IP.String() + ":" + req.Arguments
			newCP.Addr, _ = net.ResolveUDPAddr("udp", addrStr)
			if !newCP.ValidateConnection() {
				ptp.Log(ptp.ERROR, "Failed to connect to Control Peer. Ignoring")
				resp.Command = ""
			} else {
				// TODO: Consider assigning ID to Control Peers, but currently we
				// don't need such functionality
				ptp.Log(ptp.INFO, "Control peer has been validated. Saving")
				dht.ControlPeers = append(dht.ControlPeers, newCP)
			}
		}
	}
	return resp, nil
}

//func (dht *DHTRouter) ResponseNode(req ptp.DHTMessage, addr string) ptp.DHTMessage {
func (dht *DHTRouter) HandleNode(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	ptp.Log(ptp.DEBUG, "List of peers has been requested from %s", addr.String())

	var resp ptp.DHTMessage
	resp.Command = req.Command
	resp.Id = req.Query
	resp.Arguments = "0"
	p, exists := dht.PeerList[req.Query]
	if exists {
		for _, ip := range p.IPList {
			if resp.Arguments == "0" {
				resp.Arguments = ""
			}
			resp.Arguments += ip.String() + "|"
		}
	}

	return resp, nil
}

func (dht *DHTRouter) HandleNotify(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	resp.Command = req.Command
	resp.Arguments = req.Arguments
	resp.Id = "0"

	return resp, nil
}

func (dht *DHTRouter) HandleStop(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	var resp ptp.DHTMessage
	resp.Command = req.Command
	resp.Arguments = req.Id
	resp.Id = "0"
	return resp, nil
}

// ResponseCP responses to a CP request
// Request Packet contents:
// req.Query - list of CPs that should be excluded
// req.Arguments - ID of target peer
// Response Packet contents:
// resp.Arguments - control peer endpoint
//
func (dht *DHTRouter) HandleCp(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	ptp.Log(ptp.DEBUG, "Received request of control peer from %s", addr.String())
	var resp ptp.DHTMessage
	resp.Command = req.Command

	var candidate string = ""
	var minimal int = 99999

	omitList := strings.Split(req.Query, "|")
	for _, cp := range dht.ControlPeers {
		var omit bool = false
		for _, skip := range omitList {
			if skip == cp.Addr.String() {
				omit = true
			}
		}
		if omit {
			continue
		}
		if cp.ValidateConnection() {
			if cp.TunelsNum < minimal {
				candidate = cp.Addr.String()
				minimal = cp.TunelsNum
			}
		}
	}
	resp.Query = candidate
	//resp.Arguments = candidate
	resp.Arguments = req.Arguments
	// At the same moment we should send this message to a requested address too

	var sr ptp.DHTMessage
	sr.Command = req.Command
	sr.Id = req.Arguments
	sr.Arguments = req.Id
	sr.Query = candidate

	p, exists := dht.PeerList[sr.Id]
	if exists {
		dht.Send(dht.Connection, p.Addr, dht.EncodeResponse(sr))
	}

	return resp, nil
}

func (dht *DHTRouter) HandleBadCp(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	for i, cp := range dht.ControlPeers {
		if cp.Addr.String() == req.Query {
			if !cp.ValidateConnection() {
				// Remove bad control peer
				dht.ControlPeers = append(dht.ControlPeers[:i], dht.ControlPeers[i+1:]...)
				break
			}
		}
	}
	res, err := dht.HandleCp(req, addr, peer)
	return res, err
}

func (dht *DHTRouter) FindNetworkForHash(hash string) *net.IPNet {
	netinfo, exists := dht.DHCPTable[hash]
	if !exists {
		return nil
	}
	return netinfo.Network
}

func (dht *DHTRouter) PickFreeIP(ipnet *net.IPNet, used []net.IP) net.IP {
	ptp.Log(ptp.DEBUG, "Picking free IP for network: %s", ipnet.String())
	ptp.Log(ptp.DEBUG, "Used IPs: %v", used)
	ptp.Log(ptp.DEBUG, "Mask: %s, len: %d", ipnet.Mask.String(), len(ipnet.Mask))
	iplen := len(ipnet.IP)
	ipbase := fmt.Sprintf("%d.%d.%d.", ipnet.IP[iplen-4], ipnet.IP[iplen-3], ipnet.IP[iplen-2])
	for i := 3; i >= 0; i-- {
		k := int(ipnet.Mask[i])
		for j := 1; j < 255-k; j++ {
			nextIp := net.ParseIP(fmt.Sprintf("%s%d", ipbase, j))
			ptp.Log(ptp.DEBUG, "Next IP: %s", nextIp.String())
			var inUse bool = false
			for _, ip := range used {
				if nextIp.String() == ip.String() {
					inUse = true
				}
			}
			if !inUse {
				return nextIp
			}
		}
	}
	return nil
}

func (dht *DHTRouter) HandleDHCP(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	ptp.Log(ptp.INFO, "DHCP Request from %s", addr.String())
	var resp ptp.DHTMessage
	if req.Query == "0" {
		for dht.DHCPLock {
			time.Sleep(100 * time.Microsecond)
		}
		dht.DHCPLock = true
		// This is DHCP request
		for id, peer := range dht.PeerList {
			if peer.ID == req.Id {
				ipnet := dht.FindNetworkForHash(peer.AssociatedHash)
				if ipnet == nil {
					break
				}
				// Collect IPs in use
				var ips []net.IP
				for _, tp := range dht.PeerList {
					if tp.AssociatedHash == peer.AssociatedHash && tp.IP != nil {
						ips = append(ips, tp.IP)
					}
				}
				peer.IP = dht.PickFreeIP(ipnet, ips)
				peer.Network = ipnet
				dht.PeerList[id] = peer
				resp.Command = "dhcp"
				_, bits := peer.Network.Mask.Size()
				resp.Arguments = fmt.Sprintf("%s/%d", peer.IP.String(), bits)
				ptp.Log(ptp.INFO, "DHCP Response: %s", resp.Arguments)
			}
		}
		dht.DHCPLock = false
	} else {
		for dht.DHCPLock {
			time.Sleep(100 * time.Microsecond)
		}
		dht.DHCPLock = true
		// This is DHCP registration
		// We're expecting data in CIDR format
		for id, peer := range dht.PeerList {
			if peer.ID == req.Id {
				ip, ipnet, err := net.ParseCIDR(req.Query)
				if err != nil {
					ptp.Log(ptp.ERROR, "Failed to parse received DHCP information: %v", err)
					dht.DHCPLock = false
					return resp, ptp.ErrorList[ptp.ERR_BAD_DHCP_DATA]
				}
				_, exists := dht.DHCPTable[peer.AssociatedHash]
				if !exists {
					var newnet DHCPSet
					newnet.Hash = peer.AssociatedHash
					newnet.Network = ipnet
					dht.DHCPTable[peer.AssociatedHash] = newnet
				} else {
					if dht.CountParticipants(peer.AssociatedHash) == 0 {
						var newnet DHCPSet
						newnet.Hash = peer.AssociatedHash
						newnet.Network = ipnet
						dht.DHCPTable[peer.AssociatedHash] = newnet
					}
				}
				peer.IP = ip
				peer.Network = ipnet
				resp.Command = "dhcp"
				resp.Arguments = "ok"
				dht.PeerList[id] = peer
			}
		}
		dht.DHCPLock = false
	}
	return resp, nil
}

// Returns number of participants of specified hash
func (dht *DHTRouter) CountParticipants(hash string) int {
	var count int = 0
	for _, peer := range dht.PeerList {
		if peer.AssociatedHash == hash {
			count += 1
		}
	}
	return count
}

func (dht *DHTRouter) HandleLoad(req ptp.DHTMessage, addr *net.UDPAddr, peer *Peer) (ptp.DHTMessage, error) {
	for _, cp := range dht.ControlPeers {
		if cp.ID == req.Id {
			var err error
			cp.TunelsNum, err = strconv.Atoi(req.Arguments)
			if err != nil {
				cp.TunelsNum = 0
			}
		}
	}
	var resp ptp.DHTMessage
	resp.Command = ""
	return resp, nil
}

// Send method send a packet to a connected client over network to a specific UDP address
func (dht *DHTRouter) Send(conn *net.UDPConn, addr *net.UDPAddr, msg string) {
	if msg != "" {
		_, err := conn.WriteToUDP([]byte(msg), addr)
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to write to UDP: %v", err)
		}
	}
}

// This method listens to a UDP connections for incoming packets and
// sends generated responses back to DHT nodes
func (dht *DHTRouter) Listen(conn *net.UDPConn) {
	var buf [512]byte
	_, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to read from UDP socket: %v", err)
		return
	}
	req, err := dht.Extract(buf[:512])
	var peer Peer
	var exists bool
	if req.Command == ptp.CMD_CONN && dht.IsNewPeer(addr.String()) {
		ptp.Log(ptp.INFO, "New Peer connected: %s. Registering", addr)
		peer.ID = peer.GenerateID(dht.Hashes)
		peer.Endpoint = ""
		peer.ConnectionAddress = addr.String()
		peer.Addr = addr
		peer.AssociatedHash = ""
		dht.Lock.Lock()
		dht.PeerList[peer.ID] = peer
		dht.Lock.Unlock()
		runtime.Gosched()
	} else {
		dht.Lock.Lock()
		peer, exists = dht.PeerList[req.Id]
		dht.Lock.Unlock()
		runtime.Gosched()
		if !exists {
			// Send CMD_UNKNOWN for unknown peer
			var resp ptp.DHTMessage
			resp.Command = ptp.CMD_UNKNOWN
			resp.Id = req.Id
			resp.Arguments = ""
			dht.Send(conn, addr, dht.EncodeResponse(resp))
			ptp.Log(ptp.DEBUG, "Received data from unknown node. ID: %s, Command: %s", req.Id, req.Command)
			return
		}
	}
	ptp.Log(ptp.TRACE, "%s: %s", addr, string(buf[:512]))

	if peer.Disabled {
		return
	}

	// Try to bencode
	callback, exists := dht.Callbacks[req.Command]
	if exists {
		resp, err := callback(req, addr, &peer)
		if err != nil {
			dht.SendError(conn, addr, err)
		}
		if resp.Command != "" {
			ptp.Log(ptp.TRACE, "Sending %v", resp)
			dht.Send(conn, addr, dht.EncodeResponse(resp))
		}
	} else {
		ptp.Log(ptp.ERROR, "Unknown command received %s from %s", req.Command, addr.String())
	}
}

func (dht *DHTRouter) UpdateControlPeerLoad(id, amount string) {
	for key, peer := range dht.ControlPeers {
		if peer.ID == id {
			newVal, err := strconv.Atoi(amount)
			if err == nil {
				dht.ControlPeers[key].TunelsNum = newVal
			}
		}
	}
}

// Ping method is running as a goroutine. Ininity loop will
// ping every client after a timeout.
func (dht *DHTRouter) Ping(conn *net.UDPConn) {
	req := new(ptp.DHTMessage)
	req.Command = "ping"
	var removeKeys []string
	for {
		for _, key := range removeKeys {
			hash := dht.PeerList[key].AssociatedHash
			ptp.Log(ptp.WARNING, "%s timeout reached. Disconnecting", dht.PeerList[key].ConnectionAddress)
			delete(dht.PeerList, key)
			if hash != "" {
				go dht.SendStop(hash, key)
				go dht.SendFind(hash)
			}
		}
		dht.SyncControlPeers()
		removeKeys = removeKeys[:0]
		time.Sleep(PingTimeout)
		var resp ptp.DHTMessage
		resp.Command = ptp.CMD_PING
		for i, peer := range dht.PeerList {
			peer.MissedPing = peer.MissedPing + 1
			dht.Send(conn, peer.Addr, dht.EncodeResponse(resp))
			if peer.MissedPing >= 4 {
				removeKeys = append(removeKeys, i)
				peer.Disabled = true
			}
			dht.PeerList[i] = peer
		}
	}
}

func (dht *DHTRouter) SyncControlPeers() {
	for key, cp := range dht.ControlPeers {
		var found bool = false
		for _, p := range dht.PeerList {
			if p.ID == cp.ID {
				found = true
			}
		}
		if !found {
			ptp.Log(ptp.WARNING, "Removing outdated control peer: %s %s", cp.ID, cp.Addr)
			dht.ControlPeers = append(dht.ControlPeers[:key], dht.ControlPeers[key+1:]...)
		}
	}
}

func (dht *DHTRouter) SendError(conn *net.UDPConn, addr *net.UDPAddr, err error) {
	var et ptp.ErrorType = ptp.ERR_UNKNOWN_ERROR
	for t, e := range ptp.ErrorList {
		if e.Error() == err.Error() {
			et = t
		}
	}
	msg := dht.Compose(ptp.CMD_ERROR, "", string(et))
	dht.Send(conn, addr, msg)
}

func main() {
	ptp.InitErrors()
	var (
		argDht    int
		argTarget string
		argListen int
		argLog    string
	)
	flag.IntVar(&argDht, "dht", -1, "Port that DHT Bootstrap will listening to")
	flag.StringVar(&argTarget, "t", "", "Host:Port of DHT Bootstrap node")
	flag.IntVar(&argListen, "listen", 0, "Port for traffic forwarder")
	flag.StringVar(&argLog, "log", "INFO", "Log level: TRACE, DEBUG, INFO, WARNING, ERROR")
	flag.Parse()
	switch argLog {
	case "TRACE":
		ptp.SetMinLogLevel(ptp.TRACE)
	case "DEBUG":
		ptp.SetMinLogLevel(ptp.DEBUG)
	case "WARNING":
		ptp.SetMinLogLevel(ptp.WARNING)
	case "ERROR":
		ptp.SetMinLogLevel(ptp.ERROR)
	default:
		ptp.SetMinLogLevel(ptp.INFO)
	}
	ptp.Log(ptp.DEBUG, "Initialization complete")

	if argDht > 0 {
		var dht DHTRouter
		dht.Port = argDht
		dht.Connection = dht.SetupServer()
		dht.Hashes = make(map[string]Infohash)
		dht.PeerList = make(map[string]Peer)
		dht.DHCPTable = make(map[string]DHCPSet)

		dht.Callbacks = make(map[string]DHTCallback)
		dht.Callbacks[ptp.CMD_CONN] = dht.HandleConn
		dht.Callbacks[ptp.CMD_FIND] = dht.HandleFind
		dht.Callbacks[ptp.CMD_NODE] = dht.HandleNode
		dht.Callbacks[ptp.CMD_PING] = dht.HandlePing
		dht.Callbacks[ptp.CMD_REGCP] = dht.HandleRegCp
		dht.Callbacks[ptp.CMD_BADCP] = dht.HandleBadCp
		dht.Callbacks[ptp.CMD_CP] = dht.HandleCp
		dht.Callbacks[ptp.CMD_NOTIFY] = dht.HandleNotify
		dht.Callbacks[ptp.CMD_LOAD] = dht.HandleLoad
		dht.Callbacks[ptp.CMD_DHCP] = dht.HandleDHCP
		dht.Callbacks[ptp.CMD_STOP] = dht.HandleStop

		go dht.Ping(dht.Connection)

		for {
			dht.Listen(dht.Connection)
		}
	} else {
		// Act as a normal (proxy) control peer
		var proxy Proxy
		proxy.StartUDPServer(argListen)
		proxy.StartDHT(argTarget)
		proxy.Initialize()
		alivePeriod := time.Duration(time.Second * 15)
		for {
			proxy.SendPing()
			time.Sleep(3 * time.Second)
			proxy.CleanTunnels()
			passed := time.Since(proxy.DHTClient.LastDHTPing)
			if passed > alivePeriod {
				proxy.DHTClient.ID = ""
				ptp.Log(ptp.WARNING, "Lost DHT connecton. Restoring")
				proxy.StartDHT(argTarget)
			}
		}
	}
}
