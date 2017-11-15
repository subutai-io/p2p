package ptp

import (
	fmt "fmt"
	"net"
	"strings"
	"time"

	"github.com/ccding/go-stun/stun"
	proto "github.com/golang/protobuf/proto"
	uuid "github.com/wayn3h0/go-uuid"
)

// OperatingMode - Mode in which DHT client is operating
type OperatingMode int

// DHTModeClient - Indicates DHT works as a P2P client
// DHTModeProxy - Inidicates DHT works as a P2P proxy
const (
	DHTModeClient OperatingMode = 1
	DHTModeProxy  OperatingMode = 2
)

// RemotePeerState is a state information of another peer received from DHT
type RemotePeerState struct {
	ID    string
	State PeerState
}

type dhtCallback func(*DHTPacket) error

// DHTClient is a main structure of a DHT client
type DHTClient struct {
	Routers       string                        // Comma-separated list of bootstrap nodes
	FailedRouters []string                      // List of routes that we failed to connect to
	TCPConnection []*net.TCPConn                // TCP connections to bootstrap nodes
	NetworkHash   string                        // Saved network hash
	P2PPort       int                           // UDP port number used by this instance
	ID            string                        // Current instance ID
	Forwarders    []Forwarder                   // List of worwarders
	TCPCallbacks  map[DHTPacketType]dhtCallback // Callbacks for incoming packets
	Mode          OperatingMode                 // DHT Client mode ???
	IPList        []net.IP                      // List of network active interfaces
	IP            net.IP                        // IP of local interface received from DHCP or specified manually
	Network       *net.IPNet                    // Network information about current network. Used to inform p2p about mask for interface
	StateChannel  chan RemotePeerState          // Channel to pass states to instance
	ProxyChannel  chan Forwarder                // Channel to pass proxies to instance
	isShutdown    bool                          // Whether DHT shutting down or not
	PeerData      chan NetworkPeer              // Channel to pass data about changes in peers
	Connected     bool                          // Whether connection with bootstrap nodes established or not
	// LastDHTPing    time.Time
	// RemovePeerChan chan string
	// ForwardersLock sync.Mutex // To avoid multiple read-write
	// Listeners      int
	// DataChannel    chan []byte
	// Connection       []*net.UDPConn // List of connection objects
	// CommandChannel chan []byte
	// LastCatch        []string
	// ProxyBlacklist   []*net.UDPAddr
	// State            DHTState
	// ResponseHandlers map[string]DHTResponseCallback
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

// Init initialized DHT
func (dht *DHTClient) TCPInit(hash, routers string) error {
	// dht.State = DHTStateInitializing
	// dht.RemovePeerChan = make(chan string)
	//dht.PeerChannel = make(chan []PeerIP)
	dht.StateChannel = make(chan RemotePeerState)
	dht.ProxyChannel = make(chan Forwarder)
	dht.PeerData = make(chan NetworkPeer)
	dht.NetworkHash = hash
	dht.Routers = routers
	if dht.Routers == "" {
		dht.Routers = "dht1.subut.ai:6881"
	}
	dht.setupTCPCallbacks()
	return nil
}

func (dht *DHTClient) setupTCPCallbacks() {
	dht.TCPCallbacks = make(map[DHTPacketType]dhtCallback)
	dht.TCPCallbacks[DHTPacketType_BadProxy] = dht.packetBadProxy
	dht.TCPCallbacks[DHTPacketType_Connect] = dht.packetConnect
	dht.TCPCallbacks[DHTPacketType_DHCP] = dht.packetDHCP
	dht.TCPCallbacks[DHTPacketType_Error] = dht.packetError
	dht.TCPCallbacks[DHTPacketType_Find] = dht.packetFind
	dht.TCPCallbacks[DHTPacketType_Forward] = dht.packetForward
	dht.TCPCallbacks[DHTPacketType_Node] = dht.packetNode
	dht.TCPCallbacks[DHTPacketType_Notify] = dht.packetNotify
	dht.TCPCallbacks[DHTPacketType_Ping] = dht.packetPing
	dht.TCPCallbacks[DHTPacketType_Proxy] = dht.packetProxy
	dht.TCPCallbacks[DHTPacketType_RegisterProxy] = dht.packetRegisterProxy
	dht.TCPCallbacks[DHTPacketType_ReportLoad] = dht.packetReportLoad
	dht.TCPCallbacks[DHTPacketType_State] = dht.packetState
	dht.TCPCallbacks[DHTPacketType_Stop] = dht.packetStop
	dht.TCPCallbacks[DHTPacketType_Unknown] = dht.packetUnknown
	dht.TCPCallbacks[DHTPacketType_Unsupported] = dht.packetUnsupported
}

func (dht *DHTClient) TCPConnect() error {
	// Close every open connection
	for _, con := range dht.TCPConnection {
		con.Close()
	}
	dht.TCPConnection = dht.TCPConnection[:0]
	dht.FailedRouters = dht.FailedRouters[:0]
	routers := strings.Split(dht.Routers, ",")
	for _, router := range routers {
		conn, err := dht.TCPConnectAndHandshake(router, dht.IPList)
		if err != nil || conn == nil {
			Log(Error, "Failed to handshake with a DHT Server: %v", err)
			dht.FailedRouters = append(dht.FailedRouters, router)
		} else {
			Log(Info, "Handshaked. Starting listener")
			dht.TCPConnection = append(dht.TCPConnection, conn)
			go dht.TCPListen(conn)
		}
	}
	if len(dht.TCPConnection) == 0 {
		return fmt.Errorf("Failed to establish connection with bootstrap node(s)")
	}
	// dht.LastDHTPing = time.Now()
	return nil
}

func (dht *DHTClient) TCPConnectAndHandshake(router string, ipList []net.IP) (*net.TCPConn, error) {
	// TODO: Determine if we tsill need this
	// dht.State = DHTStateConnecting
	Log(Info, "Connecting to a bootstrap node (BSN) at %s", router)
	addr, err := net.ResolveTCPAddr("tcp", router)
	if err != nil {
		Log(Error, "Wrong address provided: %s router. Error: %s", router, err)
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		Log(Error, "Failed to establish connectiong with router %s", router)
		return nil, err
	}
	Log(Info, "Connected to BSN %s", router)

	err = dht.TCPHandshake(conn)
	return conn, err
}

func (dht *DHTClient) TCPHandshake(conn *net.TCPConn) error {
	Log(Info, "Requesting outbound IP")
	_, host, err := stun.NewClient().Discover()
	if err != nil {
		return fmt.Errorf("Failed to discover outbound IP: %s", err)
	}
	Log(Info, "Our IP: %s", host.IP())

	ips := []string{host.IP()}
	for _, ip := range dht.IPList {
		ips = append(ips, ip.String())
	}

	packet := DHTPacket{
		Type:      DHTPacketType_Connect,
		Arguments: ips,
		Data:      fmt.Sprintf("%d", dht.P2PPort),
		Extra:     PacketVersion,
	}
	data, err := proto.Marshal(&packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal handshake packet: %s", err)
	}
	conn.Write(data)

	return nil
}

func (dht *DHTClient) TCPListen(conn *net.TCPConn) {
	dht.Connected = true
	data := make([]byte, 2048)
	for dht.Connected {
		n, err := conn.Read(data)
		if err != nil {
			Log(Warning, "BSN socket closed: %s", err)
			dht.Connected = false
			break
		}
		packet := &DHTPacket{}
		err = proto.Unmarshal(data[:n], packet)
		if err != nil {
			Log(Warning, "Corrupted data: %s", err)
			continue
		}
		go func() {
			callback, exists := dht.TCPCallbacks[packet.Type]
			if !exists {
				Log(Error, "Unknown packet type from BSN")
				return
			}
			Log(Debug, "Received: %+v", packet)
			err = callback(packet)
			if err != nil {
				Log(Error, "%s", err)
			}
		}()
	}
}

// Sends bytes to all connected bootstrap nodes
func (dht *DHTClient) send(data []byte) error {
	for _, conn := range dht.TCPConnection {
		_, err := conn.Write(data)
		if err != nil {
			return err
		}
	}
	return nil
}

// This method will send request for network peers known to BSN
// As a response BSN will send array of IDs of peers in this swarm
func (dht *DHTClient) sendFind() error {
	if dht.NetworkHash == "" {
		return fmt.Errorf("Failed to find peers: Infohash is not set")
	}
	packet := &DHTPacket{
		Type:     DHTPacketType_Find,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal find: %s", err)
	}
	return dht.send(data)
}

// This method will send request of IPs of particular peer known to BSN
func (dht *DHTClient) sendNode(id string) error {
	if len(id) != 36 {
		return fmt.Errorf("Failed to send node: Malformed ID")
	}
	packet := &DHTPacket{
		Type: DHTPacketType_Node,
		Id:   dht.ID,
		Data: id,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal node: %s", err)
	}
	return dht.send(data)
}

func (dht *DHTClient) sendState(id, state string) error {
	if len(id) != 36 {
		return fmt.Errorf("Failed to send state: Malformed ID")
	}
	packet := &DHTPacket{
		Type:      DHTPacketType_State,
		Id:        dht.ID,
		Data:      id,
		Arguments: []string{state},
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal state: %s", err)
	}
	return dht.send(data)
}

func (dht *DHTClient) sendDHCP(ip net.IP, network *net.IPNet) error {
	subnet := "0"
	if ip == nil {
		ip = net.ParseIP("127.0.0.1")
	}
	if network != nil {
		ones, _ := network.Mask.Size()
		subnet = fmt.Sprintf("%d", ones)
	}
	packet := &DHTPacket{
		Type:  DHTPacketType_DHCP,
		Id:    dht.ID,
		Data:  ip.String(),
		Extra: subnet,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal DHCP packet: %s", err)
	}
	return dht.send(data)
}

func (dht *DHTClient) sendProxy(id string) error {

	return nil
}

func (dht *DHTClient) Shutdown() {
	Log(Info, "Entering shutdown mode. Shutting down connections with bootstrap nodes")
	dht.isShutdown = true
}

func (dht *DHTClient) WaitID() error {
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
	// dht.LastDHTPing = time.Now()
	return nil
}

// RegisterProxy will register current node as a proxy on
// bootstrap node
func (dht *DHTClient) RegisterProxy(ip net.IP, port int) error {
	id, err := uuid.NewTimeBased()
	if err != nil {
		return fmt.Errorf("Failed to generate ID: %s", err)
	}

	packet := &DHTPacket{
		Type: DHTPacketType_RegisterProxy,
		Id:   id.String(),
		Data: fmt.Sprintf("%s:%d", ip.String(), port),
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal RegProxy: %s", err)
	}
	dht.send(data)
	return nil
}

// ReportLoad will send amount of tunnels created on particular proxy
func (dht *DHTClient) ReportLoad(clientsNum int) error {
	return nil
}
