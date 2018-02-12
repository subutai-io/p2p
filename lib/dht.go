package ptp

import (
	fmt "fmt"
	"net"
	"time"

	proto "github.com/golang/protobuf/proto"
	uuid "github.com/wayn3h0/go-uuid"
)

// OperatingMode - Mode in which DHT client is operating
type OperatingMode int

// Possible operating modes
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
	Routers           string                        // Comma-separated list of bootstrap nodes
	NetworkHash       string                        // Saved network hash
	ID                string                        // Current instance ID
	FailedRouters     []string                      // List of routes that we failed to connect to
	Connections       []*net.TCPConn                // TCP connections to bootstrap nodes
	LocalPort         int                           // UDP port number used by this instance
	RemotePort        int                           // UDP port number reported by echo server
	Forwarders        []Forwarder                   // List of worwarders
	TCPCallbacks      map[DHTPacketType]dhtCallback // Callbacks for incoming packets
	Mode              OperatingMode                 // DHT Client mode ???
	IPList            []net.IP                      // List of network active interfaces
	IP                net.IP                        // IP of local interface received from DHCP or specified manually
	Network           *net.IPNet                    // Network information about current network. Used to inform p2p about mask for interface
	StateChannel      chan RemotePeerState          // Channel to pass states to instance
	ProxyChannel      chan string                   // Channel to pass proxies to instance
	PeerData          chan NetworkPeer              // Channel to pass data about changes in peers
	Connected         bool                          // Whether connection with bootstrap nodes established or not
	isShutdown        bool                          // Whether DHT shutting down or not
	LastUpdate        time.Time                     // When last `find` packet was sent
	OutboundIP        net.IP                        // Outbound IP
	ListenerIsRunning bool                          // True if listener is runnning
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

// Init bootstrap for this instance
func (dht *DHTClient) Init(hash string) error {
	dht.LastUpdate = time.Now()
	dht.StateChannel = make(chan RemotePeerState)
	dht.ProxyChannel = make(chan string)
	dht.PeerData = make(chan NetworkPeer)
	dht.NetworkHash = hash
	dht.setupTCPCallbacks()
	dht.ID = GenerateToken()
	if len(dht.ID) != 36 {
		return fmt.Errorf("Failed to produce a token")
	}
	return nil
}

// TCPInit initializes connection to DHT/bootstrap nodes over TCP
// func (dht *DHTClient) TCPInit(hash, routers string) error {
// 	// dht.LastUpdate = time.Unix(1, 1)
// 	dht.LastUpdate = time.Now()
// 	dht.StateChannel = make(chan RemotePeerState)
// 	dht.ProxyChannel = make(chan string)
// 	dht.PeerData = make(chan NetworkPeer)
// 	dht.NetworkHash = hash
// 	dht.Routers = routers
// 	if dht.Routers == "" {
// 		dht.Routers = "dht.cdn.subut.ai:6881"
// 	}
// 	dht.setupTCPCallbacks()
// 	return nil
// }

// Connect will establish TCP connection to bootstrap nodes and
// populate dht.Connetions slice with net.Conn objects
// This method will close all previous connections
// func (dht *DHTClient) Connect() error {
// 	// Close every open connection
// 	dht.Connected = false
// 	for _, con := range dht.Connections {
// 		con.Close()
// 	}

// 	dht.Connections = dht.Connections[:0]
// 	dht.FailedRouters = dht.FailedRouters[:0]
// 	routers := strings.Split(dht.Routers, ",")
// 	for _, router := range routers {
// 		conn, err := dht.ConnectAndHandshake(router, dht.IPList)
// 		if err != nil || conn == nil {
// 			Log(Error, "Failed to handshake with a DHT Server: %v", err)
// 			dht.FailedRouters = append(dht.FailedRouters, router)
// 		} else {
// 			Log(Info, "Handshaked. Starting listener")
// 			dht.Connections = append(dht.Connections, conn)
// 			go dht.Listen(conn)
// 		}
// 	}
// 	if len(dht.Connections) == 0 {
// 		return fmt.Errorf("Failed to establish connection with bootstrap node(s)")
// 	}
// 	return nil
// }

// ConnectAndHandshake will establish TCP connection to DHT Bootstrap node
// and execute dht.Handshake method
// func (dht *DHTClient) ConnectAndHandshake(router string, ipList []net.IP) (*net.TCPConn, error) {
// 	Log(Info, "Connecting to a bootstrap node (BSN) at %s", router)
// 	addr, err := net.ResolveTCPAddr("tcp", router)
// 	if err != nil {
// 		Log(Error, "Wrong address provided: %s router. Error: %s", router, err)
// 		return nil, err
// 	}
// 	conn, err := net.DialTCP("tcp", nil, addr)
// 	if err != nil {
// 		Log(Error, "Failed to establish connectiong with router %s", router)
// 		return nil, err
// 	}
// 	Log(Info, "Connected to BSN %s", router)

// 	err = dht.Handshake(conn)
// 	return conn, err
// }

// Handshake will prepare a new packet with type of DHTPacketType_Connect
// and add list of locally discovered IP addresses, UDP port and
// packet version.
// This packet will be sent immediately to a bootstrap node
// func (dht *DHTClient) Handshake(conn *net.TCPConn) error {
// 	ips := []string{}
// 	if dht.OutboundIP != nil {
// 		ips = append(ips, dht.OutboundIP.String())
// 	}
// 	for _, ip := range dht.IPList {
// 		ips = append(ips, ip.String())
// 	}

// 	packet := DHTPacket{
// 		Arguments: ips,
// 		Type:      DHTPacketType_Connect,
// 		Infohash:  dht.NetworkHash,
// 		Data:      fmt.Sprintf("%d", dht.LocalPort),
// 		Query:     fmt.Sprintf("%d", dht.RemotePort),
// 		Version:   PacketVersion,
// 	}
// 	data, err := proto.Marshal(&packet)
// 	if err != nil {
// 		return fmt.Errorf("Failed to marshal handshake packet: %s", err)
// 	}
// 	conn.Write(data)

// 	return nil
// }

// Listen will wait for incoming data to a TCP connection,
// unmarshal incoming data into DHTPacket and execute callbacks based
// on DHTPacket.Type field's value
// Callback will be executed inside a goroutine
// func (dht *DHTClient) Listen(conn *net.TCPConn) {
// 	Log(Info, "Listening to bootstrap node")
// 	dht.Connected = true
// 	data := make([]byte, 2048)
// 	dht.ListenerIsRunning = true
// 	for dht.Connected {
// 		n, err := conn.Read(data)
// 		if err != nil {
// 			Log(Warning, "BSN socket closed: %s", err)
// 			dht.Connected = false
// 			break
// 		}
// 		go func() {
// 			packet := &DHTPacket{}
// 			err = proto.Unmarshal(data[:n], packet)
// 			if err != nil {
// 				Log(Warning, "Corrupted data: %s", err)
// 				return
// 			}
// 			callback, exists := dht.TCPCallbacks[packet.Type]
// 			if !exists {
// 				Log(Error, "Unknown packet type from BSN")
// 				return
// 			}
// 			Log(Debug, "Received: %+v", packet)
// 			err = callback(packet)
// 			if err != nil {
// 				Log(Error, "%s", err)
// 			}
// 		}()
// 	}
// 	Log(Info, "DHT Listener stopped")
// 	dht.ListenerIsRunning = false
// }

// Sends bytes to all connected bootstrap nodes
func (dht *DHTClient) send(data []byte) error {
	go func() {
		for _, conn := range dht.Connections {
			_, err := conn.Write(data)
			if err != nil {
				continue
			}
		}
	}()
	return nil
}

// This method will send request for network peers known to BSN
// As a response BSN will send array of IDs of peers in this swarm
func (dht *DHTClient) sendFind() error {
	dht.LastUpdate = time.Now()
	if dht.NetworkHash == "" {
		return fmt.Errorf("Failed to find peers: Infohash is not set")
	}
	Log(Debug, "Requesting swarm updates")
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
		Type:     DHTPacketType_Node,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Data:     id,
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
		Infohash:  dht.NetworkHash,
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
		Type:     DHTPacketType_DHCP,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Data:     ip.String(),
		Extra:    subnet,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal DHCP packet: %s", err)
	}
	Log(Debug, "Sending DHCP: %+v", packet)
	return dht.send(data)
}

func (dht *DHTClient) sendProxy() error {
	Log(Debug, "Requesting proxies")
	packet := &DHTPacket{
		Type:     DHTPacketType_Proxy,
		Infohash: dht.NetworkHash,
		Id:       dht.ID,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal DHCP packet: %s", err)
	}
	return dht.send(data)
}

func (dht *DHTClient) sendRequestProxy(id string) error {
	packet := &DHTPacket{
		Type:     DHTPacketType_RequestProxy,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Data:     id,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal DHCP packet: %s", err)
	}
	return dht.send(data)
}

func (dht *DHTClient) sendReportProxy(addr []*net.UDPAddr) error {
	list := []string{}
	for _, proxy := range addr {
		list = append(list, proxy.String())
	}
	packet := &DHTPacket{
		Type:     DHTPacketType_ReportProxy,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Proxies:  list,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return fmt.Errorf("Failed to marshal DHCP packet: %s", err)
	}
	return dht.send(data)
}

// Close will close all connections and switch DHT object to
// shutdown mode, which will terminate every loop/goroutine
func (dht *DHTClient) Close() error {
	dht.Connected = false
	for _, c := range dht.Connections {
		c.Close()
	}
	Log(Info, "Entering shutdown mode. Shutting down connections with bootstrap nodes")
	if dht.ListenerIsRunning {
		Log(Info, "Waiting for DHT listener to stop")
	}
	started := time.Now()
	for dht.ListenerIsRunning {
		time.Sleep(time.Millisecond * 100)
		if time.Since(started) > time.Duration(time.Second*30) {
			Log(Error, "DHT Listener failed to stop within 30 seconds")
			return fmt.Errorf("DHT Listener failed to stop withing 30 seconds")
		}
	}
	dht.isShutdown = true
	return nil
}

// WaitID will block DHT until valid instance ID is received from Bootstrap node
// or specified timeout passes.
func (dht *DHTClient) WaitID() error {
	started := time.Now()
	period := time.Duration(time.Second * 10)
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
		Type:     DHTPacketType_RegisterProxy,
		Id:       id.String(),
		Infohash: dht.NetworkHash,
		Data:     fmt.Sprintf("%s:%d", ip.String(), port),
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
