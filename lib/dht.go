package ptp

import (
	"fmt"
	"net"
	"time"

	"github.com/wayn3h0/go-uuid"
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

// DHTClient is a main structure of a DHT client
type DHTClient struct {
	Routers       string                        // Comma-separated list of bootstrap nodes
	NetworkHash   string                        // Saved network hash
	ID            string                        // Current instance ID
	FailedRouters []string                      // List of routes that we failed to connect to
	Connections   []*net.TCPConn                // TCP connections to bootstrap nodes
	LocalPort     int                           // UDP port number used by this instance
	RemotePort    int                           // UDP port number reported by echo server
	Forwarders    []Forwarder                   // List of worwarders
	TCPCallbacks  map[DHTPacketType]dhtCallback // Callbacks for incoming packets
	Mode          OperatingMode                 // DHT Client mode ???
	IPList        []net.IP                      // List of network active interfaces
	IP            net.IP                        // IP of local interface received from DHCP or specified manually
	Network       *net.IPNet                    // Network information about current network. Used to inform p2p about mask for interface
	Connected     bool                          // Whether connection with bootstrap nodes established or not
	//isShutdown        bool                          // Whether DHT shutting down or not
	LastUpdate        time.Time // When last `find` packet was sent
	OutboundIP        net.IP    // Outbound IP
	ListenerIsRunning bool      // True if listener is runnning
	IncomingData      chan *DHTPacket
	OutgoingData      chan *DHTPacket
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
	dht.NetworkHash = hash
	dht.ID = GenerateToken()
	if len(dht.ID) != 36 {
		return fmt.Errorf("Failed to produce a token")
	}
	return nil
}

// Connect sends `conn` packet to a DHT
func (dht *DHTClient) Connect(ipList []net.IP, proxyList []*proxyServer) error {
	dht.Connected = false
	if dht.RemotePort == 0 {
		dht.RemotePort = dht.LocalPort
	}

	ips := []string{}
	proxies := []string{}
	for _, ip := range ipList {
		skip := false
		for _, a := range ActiveInterfaces {
			if a.Equal(ip) {
				skip = true
			}
		}
		if skip {
			continue
		}
		ips = append(ips, ip.String())
	}
	for _, proxy := range proxyList {
		proxies = append(proxies, proxy.Endpoint.String())
	}

	packet := &DHTPacket{
		Type:      DHTPacketType_Connect,
		Infohash:  dht.NetworkHash,
		Id:        dht.ID,
		Version:   PacketVersion,
		Data:      fmt.Sprintf("%d", dht.LocalPort),
		Query:     fmt.Sprintf("%d", dht.RemotePort),
		Arguments: ips,
		Proxies:   proxies,
	}
	err := dht.send(packet)
	if err != nil {
		return fmt.Errorf("Failed to handshake with bootstrap node: %s", err)
	}
	// Waiting for 3 seconds to get connection confirmation
	sent := time.Now()
	for time.Since(sent) < time.Duration(5000*time.Millisecond) {
		if dht.Connected {
			return nil
		}
		time.Sleep(time.Millisecond * 100)
	}
	Log(Error, "DHT handshake didn't finish")
	return fmt.Errorf("Couldn't handshake with bootstrap node")
}

func (dht *DHTClient) read() (*DHTPacket, error) {
	packet := <-dht.IncomingData
	if packet == nil {
		return nil, fmt.Errorf("Received nil packet: channel is closed")
	}
	return packet, nil
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// Sends bytes to all connected bootstrap nodes
func (dht *DHTClient) send(packet *DHTPacket) error {
	// if dht.OutgoingData != nil && !dht.isShutdown {
	if dht.OutgoingData != nil {
		// dht.OutgoingData <- packet
		if len(packet.Arguments) == 0 && len(packet.Proxies) == 0 {
			dht.OutgoingData <- packet
		} else {
			for len(packet.Arguments) != 0 || len(packet.Proxies) != 0 {
				blockLengthArgs := min(10, len(packet.Arguments))
				blockLengthProxies := min(10, len(packet.Proxies))
				args := packet.Arguments[:blockLengthArgs]
				proxies := packet.Proxies[:blockLengthProxies]
				currentPacket := &DHTPacket{
					Type:      packet.Type,
					Id:        packet.Id,
					Infohash:  packet.Infohash,
					Data:      packet.Data,
					Query:     packet.Query,
					Arguments: args,
					Proxies:   proxies,
					Extra:     packet.Extra,
					Payload:   packet.Payload,
					Version:   packet.Version,
				}
				dht.OutgoingData <- currentPacket
				packet.Arguments = packet.Arguments[blockLengthArgs:]
				packet.Proxies = packet.Proxies[blockLengthProxies:]
			}
		}
	} else {
		// Log(Debug, "%+v ||| %+v", dht.OutgoingData, dht.isShutdown)
		return fmt.Errorf("Trying to send to closed channel")
	}
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
		Version:  PacketVersion,
	}
	return dht.send(packet)
}

// This method will send request of IPs of particular peer known to BSN
func (dht *DHTClient) sendNode(id string, ipList []net.IP) error {
	if len(id) != 36 {
		return fmt.Errorf("Failed to send node: Malformed ID %s", id)
	}

	ips := []string{}
	for _, ip := range ipList {
		if ip == nil {
			continue
		}
		exists := false
		for _, eip := range ips {
			if eip == ip.String() {
				exists = true
			}
		}
		if !exists {
			ips = append(ips, ip.String())
		}
	}

	packet := &DHTPacket{
		Type:      DHTPacketType_Node,
		Id:        dht.ID,
		Infohash:  dht.NetworkHash,
		Data:      id,
		Arguments: ips,
		Version:   PacketVersion,
	}
	return dht.send(packet)
}

func (dht *DHTClient) sendState(id, state string) error {
	if len(id) != 36 {
		return fmt.Errorf("Failed to send state: Malformed ID")
	}
	packet := &DHTPacket{
		Type:     DHTPacketType_State,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Data:     id,
		Extra:    state,
		Version:  PacketVersion,
	}
	return dht.send(packet)
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
		Version:  PacketVersion,
	}
	return dht.send(packet)
}

func (dht *DHTClient) sendProxy() error {
	Log(Debug, "Requesting proxies")
	packet := &DHTPacket{
		Type:     DHTPacketType_Proxy,
		Infohash: dht.NetworkHash,
		Id:       dht.ID,
		Version:  PacketVersion,
	}
	return dht.send(packet)
}

func (dht *DHTClient) sendRequestProxy(id string) error {
	packet := &DHTPacket{
		Type:     DHTPacketType_RequestProxy,
		Id:       dht.ID,
		Infohash: dht.NetworkHash,
		Data:     id,
		Version:  PacketVersion,
	}
	return dht.send(packet)
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
		Version:  PacketVersion,
	}
	return dht.send(packet)
}

// Close will close all connections and switch DHT object to
// shutdown mode, which will terminate every loop/goroutine
func (dht *DHTClient) Close() error {
	if dht.IncomingData != nil {
		close(dht.IncomingData)
		dht.IncomingData = nil
	}
	if dht.OutgoingData != nil {
		close(dht.OutgoingData)
		dht.OutgoingData = nil
	}
	// dht.isShutdown = true
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
		Version:  PacketVersion,
	}
	return dht.send(packet)
}

// ReportLoad will send amount of tunnels created on particular proxy
func (dht *DHTClient) ReportLoad(clientsNum int) error {
	return nil
}
