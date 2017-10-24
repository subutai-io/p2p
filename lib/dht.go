package ptp

import (
	"bytes"
	"fmt"
	"net"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	bencode "github.com/jackpal/bencode-go"
)

// OperatingMode - Mode in which DHT client is operating
type OperatingMode int

// DHTState - Represents a state of current DHT client
type DHTState int

// DHTModeClient - Indicates DHT works as a P2P client
// DHTModeProxy - Inidicates DHT works as a P2P proxy
const (
	DHTModeClient OperatingMode = 1
	DHTModeProxy  OperatingMode = 2
)

// DHTStateConnecting indicates DHT client is trying to reach DHT server
// DHTStateReconnecting indicates DHT client has lost connection to server and tries to recreate it
// DHTStateOperating indicates DHT is connected and operating normally
const (
	DHTStateConnecting   DHTState = 0 + iota
	DHTStateReconnecting DHTState = 1
	DHTStateOperating    DHTState = 2
	DHTStateInitializing DHTState = 3
)

// RemotePeerState is a state information of another peer received from DHT
type RemotePeerState struct {
	ID    string
	State PeerState
}

// DHTClient is a main structure of a DHT client
type DHTClient struct {
	Routers          string         // Comma separated list of bootstrap nodes
	FailedRouters    []string       // List of routes that we failed to connect to
	Connection       []*net.UDPConn // List of connection objects
	NetworkHash      string         // Saved network hash
	P2PPort          int
	LastCatch        []string
	ID               string
	Peers            []PeerIP
	Forwarders       []Forwarder
	ProxyBlacklist   []*net.UDPAddr
	ResponseHandlers map[string]DHTResponseCallback
	Mode             OperatingMode
	IPList           []net.IP
	State            DHTState
	IP               net.IP     // IP of local interface received from DHCP or specified manually
	Network          *net.IPNet // Network information about current network. Used to inform p2p about mask for interface
	DataChannel      chan []byte
	CommandChannel   chan []byte
	StateChannel     chan RemotePeerState
	Listeners        int
	PeerChannel      chan []PeerIP
	ProxyChannel     chan Forwarder
	LastDHTPing      time.Time
	RemovePeerChan   chan string
	ForwardersLock   sync.Mutex // To avoid multiple read-write
	isShutdown       bool       // Whether DHT shutting down or not
	//NetworkPeers     []string       // List of peers
}

// Forwarder structure represents a Proxy received from DHT server
type Forwarder struct {
	Addr          *net.UDPAddr
	DestinationID string
}

// PeerIP structure represents a pair of peer ID and associated list of IP addresses
type PeerIP struct {
	ID  string
	Ips []*net.UDPAddr
}

// DHTResponseCallback is a callback executed upon receiving packets from DHT Server
type DHTResponseCallback func(data DHTMessage, conn *net.UDPConn)

// DHTClientConfig sets a default configuration for DHTClient structure
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

// Handshake performs data exchange between DHT client and server
func (dht *DHTClient) Handshake(conn *net.UDPConn) error {
	// Handshake
	var req DHTMessage
	req.ID = "0"
	req.Query = PacketVersion
	req.Command = DhtCmdConn
	// TODO: rename Port to something more clear
	req.Arguments = fmt.Sprintf("%d", dht.P2PPort)
	req.Payload = dht.NetworkHash
	for _, ip := range dht.IPList {
		req.Arguments = req.Arguments + "|" + ip.String()
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		conn.Close()
		return err
	}
	// TODO: Optimize types here
	msg := b.String()
	if dht.isShutdown {
		return nil
	}
	_, err := conn.Write([]byte(msg))
	if err != nil {
		Log(Error, "Failed to send packet: %v", err)
		conn.Close()
		return err
	}
	return nil
}

// ConnectAndHandshake sends an initial packet to a DHT bootstrap node
func (dht *DHTClient) ConnectAndHandshake(router string, ips []net.IP) (*net.UDPConn, error) {
	dht.State = DHTStateConnecting
	Log(Info, "Connecting to a router %s", router)
	addr, err := net.ResolveUDPAddr("udp", router)
	if err != nil {
		Log(Error, "Failed to resolve discovery service address: %v", err)
		return nil, err
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		Log(Error, "Failed to establish connection to discovery service: %v", err)
		return nil, err
	}

	Log(Info, "Ready to peer discovery via %s [%s]", router, conn.RemoteAddr().String())

	err = dht.Handshake(conn)

	return conn, err
}

// Extract - Extracts DHTMessage from received packet
func (dht *DHTClient) Extract(b []byte) (response DHTMessage, err error) {
	defer func() {
		if x := recover(); x != nil {
			Log(Error, "Bencode Unmarshal failed %q, %v", string(b), x)
		}
	}()
	var e2 error
	if e2 = bencode.Unmarshal(bytes.NewBuffer(b), &response); e2 == nil {
		err = nil
		return
	}
	Log(Debug, "Received from peer: %v %q", response, e2)
	return response, e2
}

// Compose - creates and returns a bencoded representation of a DHTMessage
func (dht *DHTClient) Compose(command, id, query, arguments string) string {
	var req DHTMessage
	// Command is mandatory
	req.Command = command
	// Defaults
	req.ID = "0"
	req.Query = "0"
	if id != "" {
		req.ID = id
	}

	if (req.ID == "0" || req.ID == "") && command != DhtCmdConn {
		Log(Error, "Failed to compose message, ID is empty")
		return ""
	}

	if query != "" {
		req.Query = query
	}
	req.Arguments = arguments
	return dht.EncodeRequest(req)
}

// EncodeRequest - Marshals message onto Bencode format
func (dht *DHTClient) EncodeRequest(req DHTMessage) string {
	if req.Command == "" {
		return ""
	}
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		return ""
	}
	return b.String()
}

// UpdateLastCatch - After receiving a list of peers from DHT we will parse the list
// and add every new peer into list of peers
func (dht *DHTClient) UpdateLastCatch(catch string) {
	peers := strings.Split(catch, ",")
	for _, p := range peers {
		if p == "" {
			continue
		}
		var found = false
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

// RequestPeerIPs sends a request to DHT bootstrap node with ID of
// target node we want to connect to
func (dht *DHTClient) RequestPeerIPs(id string) {
	msg := dht.Compose(DhtCmdNode, dht.ID, id, "")
	for _, conn := range dht.Connection {
		if dht.isShutdown {
			continue
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			Log(Error, "Failed to send 'node' request to %s: %v", conn.RemoteAddr().String(), err)
		}
	}
}

// UpdatePeers sends "find" request to a DHT Bootstrap node, so it can respond
// with a list of peers that we can connect to
// This method should be called periodically in case any new peers was discovered
func (dht *DHTClient) UpdatePeers() {
	for {
		if dht.isShutdown {
			break
		}
		dht.SendUpdateRequest()
		// Just in case do an update
		time.Sleep(1 * time.Minute)
	}
	Log(Info, "Stopped DHT updater")
}

// SendUpdateRequest requests a new list of peer from DHT server
func (dht *DHTClient) SendUpdateRequest() {
	msg := dht.Compose(DhtCmdFind, dht.ID, dht.NetworkHash, "")
	for _, conn := range dht.Connection {
		if dht.isShutdown {
			continue
		}
		Log(Debug, "Updating peers from %s", conn.RemoteAddr().String())
		_, err := conn.Write([]byte(msg))
		if err != nil {
			Log(Error, "Failed to send 'find' request to %s: %v", conn.RemoteAddr().String(), err)
		}
	}
}

// ListenDHT - listens for packets received from DHT bootstrap node
// Every packet is unmarshaled and turned into Request structure
// which we should analyze and respond
func (dht *DHTClient) ListenDHT(conn *net.UDPConn) {
	defer conn.Close()
	Log(Info, "Bootstraping via %s", conn.RemoteAddr().String())
	dht.Listeners++
	var failCounter = 0
	for {
		if dht.isShutdown {
			Log(Info, "Closing DHT Connection to %s", conn.RemoteAddr().String())
			conn.Close()
			for i, c := range dht.Connection {
				if c.RemoteAddr().String() == conn.RemoteAddr().String() {
					dht.Connection = append(dht.Connection[:i], dht.Connection[i+1:]...)
				}
			}
			break
		}
		var buf [2048]byte
		_, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			Log(Debug, "Failed to read from Discovery Service: %v", err)
			failCounter++
		} else {
			failCounter = 0
			data, err := dht.Extract(buf[:2048])
			if err != nil {
				Log(Error, "Failed to extract a message received from discovery service: %v", err)
			} else {
				callback, exists := dht.ResponseHandlers[data.Command]
				if exists {
					Log(Trace, "DHT Received %v", data)
					callback(data, conn)
				} else {
					Log(Debug, "Unsupported packet type received from DHT: %s", data.Command)
				}
			}
		}
		if failCounter > 1000 {
			Log(Error, "Multiple errors reading from DHT")
			break
		}
	}
	dht.Listeners--
}

// HandleConn analyzes received connecting message and assigns received
// cliend ID if any
func (dht *DHTClient) HandleConn(data DHTMessage, conn *net.UDPConn) {
	if dht.State != DHTStateConnecting && dht.State != DHTStateReconnecting {
		return
	}
	if data.ID == "" {
		Log(Error, "Empty ID was received")
		return
	}
	if data.ID == "0" {
		Log(Error, "Empty ID were received. Stopping")
		return
	}
	dht.State = DHTStateOperating
	dht.ID = data.ID
	Log(Info, "Received connection confirmation from router %s",
		conn.RemoteAddr().String())
	Log(Info, "Received personal ID for this session: %s", data.ID)
}

// HandlePing - Receives a Ping message from server and sends a response
func (dht *DHTClient) HandlePing(data DHTMessage, conn *net.UDPConn) {
	Log(Trace, "Ping message from DHT")
	dht.LastDHTPing = time.Now()
	msg := dht.Compose(DhtCmdPing, dht.ID, "", "")
	_, err := conn.Write([]byte(msg))
	if err != nil {
		Log(Error, "Failed to send 'ping' packet: %v", err)
	}
}

// ForcePing will send ping response to DHT without any requests
func (dht *DHTClient) ForcePing() {
	msg := dht.Compose(DhtCmdPing, dht.ID, "", "")
	dht.Send(msg)
}

// HandleFind - Receives a Find message with a list of peers in this environment
func (dht *DHTClient) HandleFind(data DHTMessage, conn *net.UDPConn) {
	// This means we've received a list of nodes we can connect to
	if data.Arguments != "" {
		ids := strings.Split(data.Arguments, ",")
		if len(ids) == 0 {
			Log(Error, "Malformed list of peers received")
		} else {
			// Go over list of received peer IDs and look if we know
			// anything about them. Add every new peer into list of peers
			for _, id := range ids {
				var found = false
				for _, peer := range dht.Peers {
					if peer.ID == id && len(peer.ID) > 0 {
						found = true
					}
				}
				if !found {
					var p PeerIP
					p.ID = id
					dht.Peers = append(dht.Peers, p)
				}
			}
			k := 0
			for _, peer := range dht.Peers {
				var found = false
				for _, id := range ids {
					if peer.ID == id && len(peer.ID) > 0 {
						found = true
					}
				}
				if found {
					dht.Peers[k] = peer
					k++
				}
			}
			dht.Peers = dht.Peers[:k]
			if dht.PeerChannel == nil {
				dht.PeerChannel = make(chan []PeerIP)
			}
			dht.PeerChannel <- dht.Peers
			Log(Debug, "Received peers from %s: %s", conn.RemoteAddr().String(), data.Arguments)
			dht.UpdateLastCatch(data.Arguments)
		}
	} else {
		dht.Peers = dht.Peers[:0]
	}
}

// HandleRegCp - Confirms our Proxy has been registered in DHT
func (dht *DHTClient) HandleRegCp(data DHTMessage, conn *net.UDPConn) {
	Log(Info, "This proxy has been registered in Service Discovery Peer")
	// We've received a registration confirmation message from DHT bootstrap node
}

// HandleNode - Receives a Node message
func (dht *DHTClient) HandleNode(data DHTMessage, conn *net.UDPConn) {
	// We've received an IPs associated with target node
	Log(Debug, "Received IPs from %s: %v", data.ID, data.Arguments)
	for i, peer := range dht.Peers {
		if peer.ID == data.ID {
			ips := strings.Split(data.Arguments, "|")
			var list []*net.UDPAddr
			for _, addr := range ips {
				if addr == "" {
					continue
				}
				ip, err := net.ResolveUDPAddr("udp", addr)
				if err != nil {
					Log(Error, "Failed to resolve address of peer: %v", err)
					continue
				}
				list = append(list, ip)
			}
			dht.Peers[i].Ips = list
		}
	}
}

// NotifyPeerAboutProxy - sends a notification to another peer about proxy
func (dht *DHTClient) NotifyPeerAboutProxy(id string) {
	Log(Info, "Notifying %s about proxy", id)

}

// HandleCp - receives a message with a proxy address
func (dht *DHTClient) HandleCp(data DHTMessage, conn *net.UDPConn) {
	// We've received information about proxy
	if data.Query == "0" || data.Query == "" {
		return
	}
	Log(Info, "Received forwarder %s", data.Query)
	addr, err := net.ResolveUDPAddr("udp", data.Query)
	if err != nil {
		Log(Error, "Received invalid forwarder: %v", err)
		return
	}
	var fwd Forwarder
	fwd.Addr = addr
	fwd.DestinationID = data.Arguments
	if dht.ProxyChannel == nil {
		dht.ProxyChannel = make(chan Forwarder)
	}
	dht.ProxyChannel <- fwd
	found := false
	for _, f := range dht.Forwarders {
		if f.Addr.String() == fwd.Addr.String() && f.DestinationID == fwd.DestinationID {
			found = true
		}
	}
	if !found {
		dht.Forwarders = append(dht.Forwarders, fwd)
	}
}

// HandleState will accept state message from DHT server sent by other network
// participants that we should be aware of already.
func (dht *DHTClient) HandleState(data DHTMessage, conn *net.UDPConn) {
	// We have received some state from another peer
	if dht.StateChannel == nil {
		dht.StateChannel = make(chan RemotePeerState)
	}
	var state RemotePeerState
	state.ID = data.Arguments
	numericState, err := strconv.Atoi(data.Query)
	if err != nil {
		Log(Error, "Failed to parse remote state: %s", err)
		return
	}
	state.State = PeerState(numericState)
	dht.StateChannel <- state
}

// ReportState will send specified state to DHT
func (dht *DHTClient) ReportState(targetID, state string) {
	msg := dht.Compose(DhtCmdState, dht.ID, state, targetID)
	dht.Send(msg)
}

// HandleNotify - we've received a proxy from another peer that is tries to reach us
func (dht *DHTClient) HandleNotify(data DHTMessage, conn *net.UDPConn) {
	// Notify means we should ask DHT bootstrap node for a control peer
	// in order to connect to a node that can't reach us
	// TODO: Fix this
	var l []*net.UDPAddr
	dht.RequestControlPeer(data.ID, l)
}

// HandleStop - receives a stop command from DHT server. Stop means peer should be removed from environments
func (dht *DHTClient) HandleStop(data DHTMessage, conn *net.UDPConn) {
	if data.Arguments != "" {
		// We need to stop particular peer by changing it's state to
		// P_DISCONNECT
		Log(Info, "Stop command for %s", data.Arguments)
		if dht.RemovePeerChan == nil {
			dht.RemovePeerChan = make(chan string)
		}
		dht.RemovePeerChan <- data.Arguments
	} else {
		conn.Close()
	}
}

// HandleDHCP - Received a DHCP information from server
func (dht *DHTClient) HandleDHCP(data DHTMessage, conn *net.UDPConn) {
	if data.Arguments == "ok" {
		Log(Info, "DHCP Registration confirmed")
		return
	}
	Log(Info, "Received DHCP Information: %v", data.Arguments)
	ip, ipnet, err := net.ParseCIDR(data.Arguments)
	if err != nil {
		Log(Error, "Failed to parse received DHCP packet: %v", err)
		return
	}
	Log(Info, "Saving IP/Net data: %v", ip)
	dht.IP = ip
	dht.Network = ipnet
}

// HandleUnknown - received when we was not handshaked with a DHT server
// but tried to reach some endpoints that is available only for
// handshaked clients
func (dht *DHTClient) HandleUnknown(data DHTMessage, conn *net.UDPConn) {
	Log(Warning, "DHT server refuses our identity")
	dht.ID = ""
	if dht.State == DHTStateConnecting || dht.State == DHTStateReconnecting {
		time.Sleep(3 * time.Second)
	}
	dht.State = DHTStateReconnecting
	Log(Info, "Restoring connection to a DHT bootstrap node")
	err := dht.Handshake(conn)
	if err != nil {
		Log(Error, "Failed to send new handshake packet")
	}
}

// HandleError - received an error from DHT server
func (dht *DHTClient) HandleError(data DHTMessage, conn *net.UDPConn) {
	e, exists := ErrorList[ErrorType(data.Arguments)]
	if !exists {
		Log(Error, "Unknown error were received from DHT: %s", data.Arguments)
	} else {
		Log(Error, "DHT returned error: %s", e.Error())
	}
}

// Init initialized DHT
func (dht *DHTClient) Init(hash, routers string) error {
	dht.State = DHTStateInitializing
	dht.RemovePeerChan = make(chan string)
	dht.PeerChannel = make(chan []PeerIP)
	dht.StateChannel = make(chan RemotePeerState)
	dht.ProxyChannel = make(chan Forwarder)
	dht.NetworkHash = hash
	dht.Routers = routers
	if dht.Routers == "" {
		dht.Routers = "dht1.subut.ai:6881"
	}
	dht.setupCallbacks()
	return nil
}

func (dht *DHTClient) setupCallbacks() {
	// Fallback to default working mode
	if dht.Mode != DHTModeProxy {
		dht.Mode = DHTModeClient
	}
	dht.ResponseHandlers = make(map[string]DHTResponseCallback)
	if dht.Mode != DHTModeProxy && dht.Mode != DHTModeClient {
		dht.Mode = DHTModeClient
	}
	if dht.Mode == DHTModeClient {
		Log(Info, "DHT operating in CLIENT mode")
		dht.ResponseHandlers[DhtCmdNode] = dht.HandleNode
		dht.ResponseHandlers[DhtCmdProxy] = dht.HandleCp
		dht.ResponseHandlers[DhtCmdNotify] = dht.HandleNotify
		dht.ResponseHandlers[DhtCmdStop] = dht.HandleStop
		dht.ResponseHandlers[DhtCmdState] = dht.HandleState
	} else {
		Log(Info, "DHT operating in CONTROL PEER mode")
		dht.ResponseHandlers[DhtCmdRegProxy] = dht.HandleRegCp
	}
	dht.ResponseHandlers[DhtCmdDhcp] = dht.HandleDHCP
	dht.ResponseHandlers[DhtCmdFind] = dht.HandleFind
	dht.ResponseHandlers[DhtCmdConn] = dht.HandleConn
	dht.ResponseHandlers[DhtCmdPing] = dht.HandlePing
	dht.ResponseHandlers[DhtCmdUnknown] = dht.HandleUnknown
	dht.ResponseHandlers[DhtCmdError] = dht.HandleError
}

// Connect will establish connection to bootstrap nodes
func (dht *DHTClient) Connect() error {
	if len(dht.IPList) == 0 {
		return fmt.Errorf("IP List is empty. Can't proceed with connection")
	}
	// Close every open connection
	for _, con := range dht.Connection {
		con.Close()
	}
	dht.Connection = dht.Connection[:0]
	dht.FailedRouters = dht.FailedRouters[:0]
	routers := strings.Split(dht.Routers, ",")
	for _, router := range routers {
		conn, err := dht.ConnectAndHandshake(router, dht.IPList)
		if err != nil || conn == nil {
			Log(Error, "Failed to handshake with a DHT Server: %v", err)
			dht.FailedRouters = append(dht.FailedRouters, router)
		} else {
			Log(Info, "Handshaked. Starting listener")
			dht.Connection = append(dht.Connection, conn)
			go dht.ListenDHT(conn)
		}
	}
	if len(dht.Connection) == 0 {
		return fmt.Errorf("Failed to establish connection with bootstrap node(s)")
	}
	dht.LastDHTPing = time.Now()
	return nil
}

// WaitForID will wait for ID from bootstrap node
func (dht *DHTClient) WaitForID() error {
	started := time.Now()
	period := time.Duration(time.Second * 3)
	for len(dht.ID) != 36 {
		time.Sleep(time.Millisecond * 100)
		passed := time.Since(started)
		if passed > period {
			break
		}
	}
	if len(dht.ID) != 36 {
		return fmt.Errorf("Didn't received ID from bootstrap node")
	}
	dht.LastDHTPing = time.Now()
	return nil
}

// Initialize - This method initializes DHT by splitting list of routers and connect to each one
func (dht *DHTClient) Initialize(config *DHTClient, ips []net.IP, peerChan chan []PeerIP, proxyChan chan Forwarder) *DHTClient {
	dht.RemovePeerChan = make(chan string)
	dht.PeerChannel = make(chan []PeerIP)
	dht.ProxyChannel = make(chan Forwarder)
	dht.StateChannel = make(chan RemotePeerState)
	dht = config
	//dht.PeerChannel = peerChan
	//dht.ProxyChannel = proxyChan
	routers := strings.Split(dht.Routers, ",")
	dht.FailedRouters = make([]string, len(routers))
	dht.ResponseHandlers = make(map[string]DHTResponseCallback)
	if dht.Mode != DHTModeProxy && dht.Mode != DHTModeClient {
		dht.Mode = DHTModeClient
	}
	if dht.Mode == DHTModeClient {
		Log(Info, "DHT operating in CLIENT mode")
		dht.ResponseHandlers[DhtCmdNode] = dht.HandleNode
		dht.ResponseHandlers[DhtCmdProxy] = dht.HandleCp
		dht.ResponseHandlers[DhtCmdNotify] = dht.HandleNotify
		dht.ResponseHandlers[DhtCmdStop] = dht.HandleStop
	} else {
		Log(Info, "DHT operating in CONTROL PEER mode")
		dht.ResponseHandlers[DhtCmdRegProxy] = dht.HandleRegCp
	}
	dht.ResponseHandlers[DhtCmdDhcp] = dht.HandleDHCP
	dht.ResponseHandlers[DhtCmdFind] = dht.HandleFind
	dht.ResponseHandlers[DhtCmdConn] = dht.HandleConn
	dht.ResponseHandlers[DhtCmdPing] = dht.HandlePing
	dht.ResponseHandlers[DhtCmdUnknown] = dht.HandleUnknown
	dht.ResponseHandlers[DhtCmdError] = dht.HandleError
	dht.IPList = ips
	var connected int
	for _, con := range dht.Connection {
		con.Close()
	}
	dht.Connection = dht.Connection[:0]
	for _, router := range routers {
		conn, err := dht.ConnectAndHandshake(router, dht.IPList)
		if err != nil || conn == nil {
			Log(Error, "Failed to handshake with a DHT Server: %v", err)
			dht.FailedRouters[0] = router
		} else {
			Log(Info, "Handshaked. Starting listener")
			dht.Connection = append(dht.Connection, conn)
			connected++
			go dht.ListenDHT(conn)
		}
	}
	started := time.Now()
	period := time.Duration(time.Second * 3)
	for len(dht.ID) != 36 {
		time.Sleep(time.Millisecond * 100)
		passed := time.Since(started)
		if passed > period {
			break
		}
	}
	dht.LastDHTPing = time.Now()
	if connected == 0 {
		return nil
	}
	return dht
}

// RegisterControlPeer - This method register control peer on a Bootstrap node
func (dht *DHTClient) RegisterControlPeer() {
	for len(dht.ID) != 36 {
		time.Sleep(1 * time.Second)
	}
	var req DHTMessage
	var err error
	req.ID = dht.ID
	req.Query = "0"
	req.Command = DhtCmdRegProxy
	req.Arguments = fmt.Sprintf("%d", dht.P2PPort)
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		return
	}
	// TODO: Optimize types here
	msg := b.String()
	for _, conn := range dht.Connection {
		if dht.isShutdown {
			continue
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			Log(Error, "Failed to send packet: %v", err)
			conn.Close()
			return
		}
	}
}

// RequestControlPeer - This method request a new control peer for particular host
func (dht *DHTClient) RequestControlPeer(id string, omit []*net.UDPAddr) {
	var req DHTMessage
	var err error
	req.ID = dht.ID
	req.Query = ""
	// Collect list of failed forwarders
	for _, fwd := range omit {
		req.Query += fwd.String() + "|"
	}
	req.Command = DhtCmdProxy
	req.Arguments = id
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		return
	}
	msg := b.String()
	// TODO: Move sending to a separate method
	for _, conn := range dht.Connection {
		if dht.isShutdown {
			continue
		}
		_, err = conn.Write([]byte(msg))
		if err != nil {
			Log(Error, "Failed to send packet: %v", err)
			conn.Close()
			return
		}
	}
}

// ReportControlPeerLoad - sends current amount of clients on this proxy
func (dht *DHTClient) ReportControlPeerLoad(amount int) {
	var req DHTMessage
	req.ID = dht.ID
	req.Command = DhtCmdLoad
	req.Arguments = fmt.Sprintf("%d", amount)
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		return
	}
	dht.Send(b.String())
}

// Send - sends a DHT message to a DHT server
func (dht *DHTClient) Send(msg string) bool {
	if msg == "" {
		Log(Error, "Failed to send DHT packet: empty msg")
		return false
	}
	for _, conn := range dht.Connection {
		if dht.isShutdown {
			continue
		}
		_, err := conn.Write([]byte(msg))
		if err != nil {
			Log(Error, "Failed to send DHT packet: %v", err)
			return false
		}
	}
	return true
}

// RequestIP - Requests an IP from DHT. DHT Server will understand empty query field
// and send IP in response
func (dht *DHTClient) RequestIP() {
	Log(Info, "Sending DHCP request")
	req := dht.Compose(DhtCmdDhcp, dht.ID, "", "")
	dht.Send(req)
}

// SendIP - Notify DHT about configured IP and netmask
func (dht *DHTClient) SendIP(ip, mask string) {
	Log(Info, "Sending DHCP information. IP: %s, Mask: %s", ip, mask)
	req := dht.Compose(DhtCmdDhcp, dht.ID, ip, mask)
	dht.Send(req)
}

// Stop - sends a STOP message about current peer
func (dht *DHTClient) Stop() {
	dht.Shutdown()
	var req DHTMessage
	req.ID = dht.ID
	req.Command = DhtCmdStop
	req.Arguments = "0"
	var b bytes.Buffer
	if err := bencode.Marshal(&b, req); err != nil {
		Log(Error, "Failed to Marshal bencode %v", err)
		return
	}
	msg := b.String()
	for _, conn := range dht.Connection {
		conn.Write([]byte(msg))
	}
}

// BlacklistForwarder - adds a proxy to a blacklist
func (dht *DHTClient) BlacklistForwarder(addr *net.UDPAddr) {
	dht.ForwardersLock.Lock()
	// Remove it from list of cached forwarders
	for i, fwd := range dht.Forwarders {
		if fwd.Addr.String() == addr.String() {
			dht.Forwarders = append(dht.Forwarders[:i], dht.Forwarders[i+1:]...)
			break
		}
	}
	found := false
	for _, fwd := range dht.ProxyBlacklist {
		if fwd.String() == addr.String() {
			found = true
		}
	}
	if !found {
		dht.ProxyBlacklist = append(dht.ProxyBlacklist, addr)
	}
	dht.ForwardersLock.Unlock()
	runtime.Gosched()
}

// CleanForwarderBlacklist - removes all entries about blacklisted proxies
func (dht *DHTClient) CleanForwarderBlacklist() {
	Log(Debug, "Cleaning forwarders blacklist")
	dht.ProxyBlacklist = dht.ProxyBlacklist[:0]
}

// CleanPeer will remove information about peer with specified ID
func (dht *DHTClient) CleanPeer(id string) error {
	for i, p := range dht.Peers {
		if p.ID == id {
			dht.Peers = append(dht.Peers[:i], dht.Peers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Specified peer was not found")
}

// Shutdown will turn DHT to shutdown state
func (dht *DHTClient) Shutdown() {
	dht.isShutdown = true
}
