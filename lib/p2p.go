package ptp

import (
	"crypto/rand"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type MessageHandler func(message *P2PMessage, src_addr *net.UDPAddr)

// Main structure
type PTPCloud struct {
	IP              string                               // Interface IP address
	Mac             string                               // String representation of a MAC address
	HardwareAddr    net.HardwareAddr                     // MAC address of network interface
	Mask            string                               // Network mask in the dot-decimal notation
	DeviceName      string                               // Name of the network interface
	IPTool          string                               `yaml:"iptool"` // Network interface configuration tool
	Device          *Interface                           // Network interface
	NetworkPeers    map[string]NetworkPeer               // Knows peers
	UDPSocket       *PTPNet                              // Peer-to-peer interconnection socket
	LocalIPs        []net.IP                             // List of IPs available in the system
	Dht             *DHTClient                           // DHT Client
	Crypter         Crypto                               // Instance of crypto
	Shutdown        bool                                 // Set to true when instance in shutdown mode
	Restart         bool                                 // Instance will be restarted
	IPIDTable       map[string]string                    // Mapping for IP->ID
	MACIDTable      map[string]string                    // Mapping for MAC->ID
	ForwardMode     bool                                 // Skip local peer discovery
	MessageHandlers map[uint16]MessageHandler            // Callbacks
	ReadyToStop     bool                                 // Set to true when instance is ready to stop
	PacketHandlers  map[PacketType]PacketHandlerCallback // Callbacks for network packet handlers
	DHTPeerChannel  chan []string
}

type NetworkPeer struct {
	ID           string           // ID of a peer
	Unknown      bool             // TODO: Remove after moving to states
	Handshaked   bool             // TODO: Remove after moving to states
	WaitingPing  bool             // True if ping request was sent
	ProxyID      int              // ID of the proxy
	ProxyRetries int              // Number of retries to reach proxy
	Forwarder    *net.UDPAddr     // Forwarder address
	PeerAddr     *net.UDPAddr     // Address of peer
	PeerLocalIP  net.IP           // IP of peers interface. TODO: Rename to IP
	PeerHW       net.HardwareAddr // Hardware addres of peer interface. TODO: Rename to Mac
	Endpoint     string           // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs     []*net.UDPAddr   // List of IP addresses that accepts connection on peer
	Retries      int              // Number of introduction retries
	Ready        bool             // Set to true when peer is ready to communicate with p2p network
	State        PeerState        // State of a peer
	LastContact  time.Time        // Last ping with this peer
}

// Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PTPCloud) AssignInterface(ip, mac, mask, device string) error {
	var err error

	p.IP = ip
	p.Mac = mac
	p.Mask = mask
	p.DeviceName = device

	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile(CONFIG_DIR + "/p2p/config.yaml")
	if err != nil {
		Log(ERROR, "Failed to load config: %v", err)
		p.IPTool = "/sbin/ip"
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		Log(ERROR, "Failed to parse config: %v", err)
		return err
	}

	p.Device, err = Open(p.DeviceName, DevTap)
	if p.Device == nil {
		Log(ERROR, "Failed to open TAP device %s: %v", device, err)
		return err
	} else {
		Log(INFO, "%v TAP Device created", p.DeviceName)
	}

	err = ConfigureInterface(p.Device, p.IP, mac, p.DeviceName, p.IPTool)
	if err != nil {
		return err
	}
	return nil
}

// Handles a packet that was received by TUN/TAP device
// Receiving a packet by device means that some application sent a network
// packet within a subnet in which our application works.
// This method calls appropriate gorouting for extracted packet protocol
func (p *PTPCloud) handlePacket(contents []byte, proto int) {
	callback, exists := p.PacketHandlers[PacketType(proto)]
	if exists {
		callback(contents, proto)
	} else {
		Log(WARNING, "Captured undefined packet: %d", PacketType(proto))
	}
}

// Listen TAP interface for incoming packets
func (p *PTPCloud) ListenInterface() {
	// Read packets received by TUN/TAP device and send them to a handlePacket goroutine
	// This goroutine will decide what to do with this packet

	// Run is for windows only
	p.Device.Run()
	for {
		if p.Shutdown {
			break
		}
		packet, err := p.Device.ReadPacket()
		if err != nil {
			Log(ERROR, "Reading packet %s", err)
		}
		if packet.Truncated {
			Log(DEBUG, "Truncated packet")
		}
		// TODO: Make handlePacket as a part of PTPCloud
		go p.handlePacket(packet.Packet, packet.Protocol)
	}
	p.Device.Close()
	Log(INFO, "Shutting down interface listener")
}

func (p *PTPCloud) IsDeviceExists(name string) bool {
	inf, err := net.Interfaces()
	if err != nil {
		Log(ERROR, "Failed to retrieve list of network interfaces")
		return true
	}
	for _, i := range inf {
		if i.Name == name {
			return true
		}
	}
	return false
}

// This method will generate device name if none were specified at startup
func (p *PTPCloud) GenerateDeviceName(i int) string {
	var devName string = GetDeviceBase() + fmt.Sprintf("%d", i)
	if p.IsDeviceExists(devName) {
		return p.GenerateDeviceName(i + 1)
	} else {
		return devName
	}
}

func (p *PTPCloud) IsIPv4(ip string) bool {
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case ':':
			return false
		case '.':
			return true
		}
	}
	return false
}

// This method lists interfaces available in the system and retrieves their
// IP addresses
func (p *PTPCloud) FindNetworkAddresses() {
	Log(INFO, "Looking for available network interfaces")
	inf, err := net.Interfaces()
	if err != nil {
		Log(ERROR, "Failed to retrieve list of network interfaces")
		return
	}
	for _, i := range inf {
		addresses, err := i.Addrs()

		if err != nil {
			Log(ERROR, "Failed to retrieve address for interface. %v", err)
			continue
		}
		for _, addr := range addresses {
			var decision string = "Ignoring"
			var ipType string = "Unknown"
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				Log(ERROR, "Failed to parse CIDR notation: %v", err)
			}
			if ip.IsLoopback() {
				ipType = "Loopback"
			} else if ip.IsMulticast() {
				ipType = "Multicast"
			} else if ip.IsGlobalUnicast() {
				decision = "Saving"
				ipType = "Global Unicast"
			} else if ip.IsLinkLocalUnicast() {
				ipType = "Link Local Unicast"
			} else if ip.IsLinkLocalMulticast() {
				ipType = "Link Local Multicast"
			} else if ip.IsInterfaceLocalMulticast() {
				ipType = "Interface Local Multicast"
			}
			if !p.IsIPv4(ip.String()) {
				decision = "No IPv4"
			}
			Log(INFO, "Interface %s: %s. Type: %s. %s", i.Name, addr.String(), ipType, decision)
			if decision == "Saving" {
				p.LocalIPs = append(p.LocalIPs, ip)
			}
		}
	}
	Log(INFO, "%d interfaces were saved", len(p.LocalIPs))
}

func StartP2PInstance(argIp, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool, port int) *PTPCloud {

	var hw net.HardwareAddr

	if argMac != "" {
		var err2 error
		hw, err2 = net.ParseMAC(argMac)
		if err2 != nil {
			Log(ERROR, "Invalid MAC address provided: %v", err2)
			return nil
		}
	} else {
		argMac, hw = GenerateMAC()
		Log(INFO, "Generate MAC for TAP device: %s", argMac)
	}

	// Create new DHT Client, configured it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code
	dhtClient := new(DHTClient)
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = argHash
	config.Mode = MODE_CLIENT

	p := new(PTPCloud)
	p.FindNetworkAddresses()
	p.HardwareAddr = hw
	p.NetworkPeers = make(map[string]NetworkPeer)
	p.IPIDTable = make(map[string]string)
	p.MACIDTable = make(map[string]string)

	if fwd {
		p.ForwardMode = true
	}

	if argDev == "" {
		argDev = p.GenerateDeviceName(1)
	} else {
		if len(argDev) > 12 {
			Log(INFO, "Interface name lenght should be 12 symbols max")
			return nil
		}
	}
	if p.IsDeviceExists(argDev) {
		Log(ERROR, "Interface is already in use. Can't create duplicate")
		return nil
	}

	if argKeyfile != "" {
		p.Crypter.ReadKeysFromFile(argKeyfile)
	}
	if argKey != "" {
		// Override key from file
		if argTTL == "" {
			argTTL = "default"
		}
		var newKey CryptoKey
		newKey = p.Crypter.EnrichKeyValues(newKey, argKey, argTTL)
		p.Crypter.Keys = append(p.Crypter.Keys, newKey)
		p.Crypter.ActiveKey = p.Crypter.Keys[0]
		p.Crypter.Active = true
	}

	if p.Crypter.Active {
		Log(INFO, "Traffic encryption is enabled. Key valid until %s", p.Crypter.ActiveKey.Until.String())
	} else {
		Log(INFO, "No AES key were provided. Traffic encryption is disabled")
	}

	// Register network message handlers
	p.MessageHandlers = make(map[uint16]MessageHandler)
	p.MessageHandlers[MT_NENC] = p.HandleNotEncryptedMessage
	p.MessageHandlers[MT_PING] = p.HandlePingMessage
	p.MessageHandlers[MT_XPEER_PING] = p.HandleXpeerPingMessage
	p.MessageHandlers[MT_ENC] = p.HandleMessage
	p.MessageHandlers[MT_INTRO] = p.HandleIntroMessage
	p.MessageHandlers[MT_INTRO_REQ] = p.HandleIntroRequestMessage
	p.MessageHandlers[MT_PROXY] = p.HandleProxyMessage
	p.MessageHandlers[MT_TEST] = p.HandleTestMessage
	p.MessageHandlers[MT_BAD_TUN] = p.HandleBadTun

	// Register packet handlers
	p.PacketHandlers = make(map[PacketType]PacketHandlerCallback)
	p.PacketHandlers[PT_PARC_UNIVERSAL] = p.handlePARCUniversalPacket
	p.PacketHandlers[PT_IPV4] = p.handlePacketIPv4
	p.PacketHandlers[PT_ARP] = p.handlePacketARP
	p.PacketHandlers[PT_RARP] = p.handleRARPPacket
	p.PacketHandlers[PT_8021Q] = p.handle8021qPacket
	p.PacketHandlers[PT_IPV6] = p.handlePacketIPv6
	p.PacketHandlers[PT_PPPOE_DISCOVERY] = p.handlePPPoEDiscoveryPacket
	p.PacketHandlers[PT_PPPOE_SESSION] = p.handlePPPoESessionPacket

	p.UDPSocket = new(PTPNet)
	p.UDPSocket.Init("", port)
	port = p.UDPSocket.GetPort()
	Log(INFO, "Started UDP Listener at port %d", port)
	config.P2PPort = port
	if argDht != "" {
		config.Routers = argDht
	}
	p.DHTPeerChannel = make(chan []string)
	p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel)
	for p.Dht == nil {
		Log(WARNING, "Failed to connect to DHT. Retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		p.LocalIPs = p.LocalIPs[:0]
		p.FindNetworkAddresses()
		p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel)
	}
	// Wait for ID
	for len(p.Dht.ID) < 32 {
		time.Sleep(100 * time.Millisecond)
	}
	Log(INFO, "ID assigned. Continue")
	var retries int = 0
	if argIp == "dhcp" {
		Log(INFO, "Requesting IP")
		p.Dht.RequestIP()
		time.Sleep(1 * time.Second)
		for p.Dht.IP == nil && p.Dht.Network == nil {
			Log(INFO, "No IP were received. Requesting again")
			p.Dht.RequestIP()
			time.Sleep(3 * time.Second)
			retries++
			if retries >= 10 {
				Log(ERROR, "Failed to retrieve IP from network after 10 retries")
				return nil
			}
		}
		m := p.Dht.Network.Mask
		mask := fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
		p.AssignInterface(p.Dht.IP.String(), argMac, mask, argDev)
	} else {
		ip, ipnet, err := net.ParseCIDR(argIp)
		if err != nil {
			nip := net.ParseIP(argIp)
			if nip == nil {
				Log(ERROR, "Invalid address were provided for network interface. Use -ip \"dhcp\" or specify correct IP address")
				return nil
			}
			argIp += `/24`
			Log(WARNING, "No CIDR mask was provided. Assumming /24")
			ip, ipnet, err = net.ParseCIDR(argIp)
			if err != nil {
				Log(ERROR, "Failed to setup provided IP address for local device")
				return nil
			}
		}
		p.Dht.IP = ip
		p.Dht.Network = ipnet
		mask := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
		p.Dht.SendIP(argIp, mask)
		p.AssignInterface(p.Dht.IP.String(), argMac, mask, argDev)
	}

	go p.UDPSocket.Listen(p.HandleP2PMessage)

	go p.ListenInterface()
	return p
}

func (p *PTPCloud) Run() {
	for {
		if p.Shutdown {
			// TODO: Do it more safely
			if p.ReadyToStop {
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(3 * time.Second)
		p.Dht.UpdatePeers()
		// Wait two seconds before synchronizing with catched peers
		time.Sleep(2 * time.Second)
		p.PurgePeers()
		newPeersNum := p.SyncPeers()
		newPeersNum = newPeersNum + p.SyncForwarders()
		p.IntroducePeers()
		if p.Dht.Listeners == 0 {
			p.Shutdown = true
			p.Restart = true
		}
	}
	Log(INFO, "Shutting down instance %s completed", p.Dht.NetworkHash)
}

func (p *PTPCloud) TouchPeers() {
	for i, peer := range p.NetworkPeers {
		if peer.State != P_CONNECTED {
			continue
		}
		passed := time.Since(peer.LastContact)
		if passed > PEER_PING_TIMEOUT {
			if peer.WaitingPing {
				Log(INFO, "Removing timeout peer: %s", peer.ID)
				delete(p.NetworkPeers, i)
				break
			}
			peer.LastContact = time.Now()
			peer.WaitingPing = true
			p.NetworkPeers[i] = peer
			p.Ping(peer.PeerAddr)
		}
	}
	time.Sleep(100 * time.Microsecond)
}

func (p *PTPCloud) Ping(addr *net.UDPAddr) {
	msg := CreateXpeerPingMessage(PING_REQ)
	p.UDPSocket.SendMessage(msg, addr)
}

// This method sends information about himself to empty peers
// Empty peers is a peer that was not sent us information
// about his device
func (p *PTPCloud) IntroducePeers() {
	for i, peer := range p.NetworkPeers {
		// Skip if know this peer
		if !peer.Unknown {
			continue
		}
		// Skip if we don't have an endpoint address for this peer
		if peer.Endpoint == "" {
			continue
		}

		if !peer.Ready {
			continue
		}

		if peer.Retries >= 10 {
			if peer.ProxyID != 0 {
				Log(WARNING, "Failed to introduce to %s via proxy %d", peer.ID, peer.ProxyID)
				p.Dht.MakeForwarderFailed(peer.Forwarder)
				peer.Endpoint = ""
				peer.Forwarder = nil
				peer.ProxyID = 0
				peer.State = P_INIT
				p.NetworkPeers[i] = peer
			} else {
				Log(WARNING, "Failed to introduce to %s", peer.ID)
				peer.State = P_HANDSHAKING_FAILED
				p.NetworkPeers[i] = peer
			}
			continue
		}
		peer.Retries = peer.Retries + 1
		Log(DEBUG, "Intoducing to %s", peer.Endpoint)
		addr, err := net.ResolveUDPAddr("udp", peer.Endpoint)
		if err != nil {
			Log(ERROR, "Failed to resolve UDP address during Introduction: %v", err)
			continue
		}
		//peer.PeerAddr = addr
		p.NetworkPeers[i] = peer
		// Send introduction packet
		//msg := p.PrepareIntroductionMessage(p.Dht.ID)
		msg := CreateIntroRequest(p.Crypter, p.Dht.ID)
		msg.Header.ProxyId = uint16(peer.ProxyID)
		_, err = p.UDPSocket.SendMessage(msg, addr)
		if err != nil {
			Log(ERROR, "Failed to send introduction to %s", addr.String())
		} else {
			Log(DEBUG, "Introduction sent to %s", peer.Endpoint)
		}
	}
}

func (p *PTPCloud) PrepareIntroductionMessage(id string) *P2PMessage {
	var intro string = id + "," + p.Mac + "," + p.IP
	msg := CreateIntroP2PMessage(p.Crypter, intro, 0)
	return msg
}

// This method goes over peers and removes obsolete ones
// Peer becomes obsolete when it goes out of DHT
func (p *PTPCloud) PurgePeers() {
	for i, peer := range p.NetworkPeers {
		var f bool = false
		for _, newPeer := range p.Dht.Peers {
			if newPeer.ID == peer.ID {
				f = true
			}
		}
		if !f {
			Log(INFO, ("Removing outdated peer"))
			delete(p.IPIDTable, peer.PeerLocalIP.String())
			delete(p.MACIDTable, peer.PeerHW.String())
			delete(p.NetworkPeers, i)
		}
	}
	return
}

// This method tests connection with specified endpoint
func (p *PTPCloud) TestConnection(endpoint *net.UDPAddr) bool {
	msg := CreateTestP2PMessage(p.Crypter, "TEST", 0)
	conn, err := net.DialUDP("udp4", nil, endpoint)
	if err != nil {
		Log(ERROR, "%v", err)
		return false
	}
	ser := msg.Serialize()
	_, err = conn.Write(ser)
	if err != nil {
		conn.Close()
		return false
	}
	t := time.Now()
	t = t.Add(3 * time.Second)
	conn.SetReadDeadline(t)
	for {
		var buf [4096]byte
		s, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			Log(ERROR, "%v", err)
			conn.Close()
			return false
		}
		if s > 0 {
			conn.Close()
			return true
		}
	}
	conn.Close()
	return false
}

func (p *PTPCloud) SyncForwarders() int {
	var count int = 0
	for _, fwd := range p.Dht.Forwarders {
		for key, peer := range p.NetworkPeers {
			if peer.Endpoint == "" && fwd.DestinationID == peer.ID && peer.Forwarder == nil {
				Log(INFO, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr.String()
				peer.Forwarder = fwd.Addr
				peer.State = P_HANDSHAKING_FORWARDER
				p.NetworkPeers[key] = peer
				count = count + 1
			}
		}
	}
	p.Dht.Forwarders = p.Dht.Forwarders[:0]
	return count
}

func (p *PTPCloud) AssignEndpoint(peer NetworkPeer) (string, bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		Log(ERROR, "Failed to retrieve list of network interfaces")
		//failback = true
	}

	for _, inf := range interfaces {
		if peer.Endpoint != "" {
			break
		}
		if inf.Name == p.DeviceName {
			continue
		}
		addrs, _ := inf.Addrs()
		for _, addr := range addrs {
			netip, network, _ := net.ParseCIDR(addr.String())
			if !netip.IsGlobalUnicast() {
				continue
			}
			for _, kip := range peer.KnownIPs {
				Log(TRACE, "Probing new IP %s against network %s", kip.IP.String(), network.String())

				if network.Contains(kip.IP) {
					if p.TestConnection(kip) {
						return kip.String(), true
						Log(INFO, "Setting endpoint for %s to %s", peer.ID, kip.String())
					}
				}
			}
		}
	}
	return "", false
}

// This method takes a list of catched peers from DHT and
// adds every new peer into list of peers
// Returns amount of peers that has been added
func (p *PTPCloud) SyncPeers() int {
	var count int = 0

	for _, id := range p.Dht.Peers {
		if id.ID == "" {
			continue
		}
		var found bool = false
		for i, peer := range p.NetworkPeers {
			if peer.ID == id.ID {
				found = true
				// Check if we know something new about this peer, e.g. new addresses were
				// assigned to it
				for _, ip := range id.Ips {
					if ip == "" || ip == "0" {
						continue
					}
					var ipFound bool = false
					for _, kip := range peer.KnownIPs {
						if kip.String() == ip {
							ipFound = true
						}
					}
					if !ipFound {
						Log(INFO, "Adding new IP (%s) address to %s", ip, peer.ID)
						// TODO: Check IP parsing
						newIp, _ := net.ResolveUDPAddr("udp", ip)
						peer.KnownIPs = append(peer.KnownIPs, newIp)
						p.NetworkPeers[i] = peer
					}
				}

				// Set and Endpoint from peers if no endpoint were set previously
				if peer.State == P_INIT {
					var added bool
					if !p.ForwardMode {
						peer.Endpoint, added = p.AssignEndpoint(peer)
						peer.State = P_CONNECTED
						peer.LastContact = time.Now()
					}
					if added {
						p.NetworkPeers[i] = peer
						count += 1
					} else {
						if len(peer.KnownIPs) > 0 {
							if !p.ForwardMode {
								Log(INFO, "No peers are available in local network. Switching to global IP if any")
							}
							// If endpoint wasn't set let's test connection from outside of the LAN
							// First one should be the global IP (if DHT works correctly)
							if p.ForwardMode || !p.TestConnection(p.NetworkPeers[i].KnownIPs[0]) {
								if peer.State == P_WAITING_FORWARDER {
									continue
								}
								// We've failed to establish connection again. Now let's ask for a proxy
								Log(INFO, "Requesting Control Peer from Service Discovery Peer")
								p.Dht.RequestControlPeer(peer.ID)
								peer.PeerAddr = peer.KnownIPs[0]
								peer.State = P_WAITING_FORWARDER
								p.NetworkPeers[i] = peer
							} else {
								Log(INFO, "Successfully connected to a host over Internet")
								peer.Endpoint = peer.KnownIPs[0].String()
								peer.State = P_CONNECTED
								peer.LastContact = time.Now()
								p.NetworkPeers[i] = peer
							}
						}
					}

				} else if peer.State == P_HANDSHAKING_FAILED {
					Log(INFO, "Requesting backup control peer")
					p.Dht.RequestControlPeer(peer.ID)
					peer.PeerAddr = peer.KnownIPs[0]
					peer.State = P_WAITING_FORWARDER
					p.NetworkPeers[i] = peer
				} else if peer.State == P_HANDSHAKING_FORWARDER {
					if peer.ProxyRetries > 3 {
						Log(WARNING, "Failed to handshake with control peer %s. Adding to black list", peer.Forwarder.String())
						p.Dht.MakeForwarderFailed(peer.Forwarder)
						peer.Forwarder = nil
						peer.Endpoint = ""
						peer.ProxyID = 0
						peer.State = P_INIT
						continue
					}
					peer.ProxyRetries = peer.ProxyRetries + 1
					p.NetworkPeers[i] = peer
					Log(INFO, "Sending proxy request to a forwarder %s", peer.Forwarder.String())
					msg := CreateProxyP2PMessage(-1, peer.PeerAddr.String(), 0)
					_, err := p.UDPSocket.SendMessage(msg, peer.Forwarder)
					if err != nil {
						Log(ERROR, "Failed to send a message to a proxy %v", err)
					}
				}
			}
		}
		if !found {
			Log(INFO, "Adding new peer. Requesting peer address")
			var newPeer NetworkPeer
			newPeer.ID = id.ID
			newPeer.Unknown = true
			newPeer.Ready = true
			newPeer.State = P_INIT
			if p.ForwardMode {
				newPeer.Ready = false
			}
			p.NetworkPeers[newPeer.ID] = newPeer
			p.Dht.RequestPeerIPs(id.ID)
		}
	}
	return count
}

// WriteToDevice writes data to created TUN/TAP device
func (p *PTPCloud) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var packet Packet
	packet.Protocol = int(proto)
	packet.Truncated = truncated
	packet.Packet = b
	if p.Device == nil {
		Log(ERROR, "TUN/TAP Device not initialized")
		return
	}
	err := p.Device.WritePacket(&packet)
	if err != nil {
		Log(ERROR, "Failed to write to TUN/TAP device: %v", err)
	}
}

func GenerateMAC() (string, net.HardwareAddr) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		Log(ERROR, "Failed to generate MAC: %v", err)
		return "", nil
	}
	buf[0] |= 2
	mac := fmt.Sprintf("06:%02x:%02x:%02x:%02x:%02x", buf[1], buf[2], buf[3], buf[4], buf[5])
	hw, err := net.ParseMAC(mac)
	if err != nil {
		Log(ERROR, "Corrupted MAC address generated: %v", err)
		return "", nil
	}
	return mac, hw
}

// AddPeer adds new peer into list of network participants. If peer was added previously
// information about him will be updated. If not, new entry will be added
func (p *PTPCloud) AddPeer(addr *net.UDPAddr, id string, ip net.IP, mac net.HardwareAddr) {
	var found bool = false
	for i, peer := range p.NetworkPeers {
		if peer.ID == id {
			found = true
			peer.ID = id
			peer.PeerAddr = addr
			peer.PeerLocalIP = ip
			peer.PeerHW = mac
			peer.Unknown = false
			peer.Handshaked = true
			p.NetworkPeers[i] = peer
			p.IPIDTable[ip.String()] = id
			p.MACIDTable[mac.String()] = id
		}
	}
	if !found {
		var newPeer NetworkPeer
		newPeer.ID = id
		newPeer.PeerAddr = addr
		newPeer.PeerLocalIP = ip
		newPeer.PeerHW = mac
		newPeer.Unknown = false
		newPeer.Handshaked = true
		p.NetworkPeers[newPeer.ID] = newPeer
		p.IPIDTable[ip.String()] = id
		p.MACIDTable[mac.String()] = id
	}
}

func (p *NetworkPeer) ProbeConnection() bool {
	return false
}

func (p *PTPCloud) ParseIntroString(intro string) (string, net.HardwareAddr, net.IP) {
	parts := strings.Split(intro, ",")
	if len(parts) != 3 {
		Log(ERROR, "Failed to parse introduction string")
		return "", nil, nil
	}
	var id string
	id = parts[0]
	// Extract MAC
	mac, err := net.ParseMAC(parts[1])
	if err != nil {
		Log(ERROR, "Failed to parse MAC address from introduction packet: %v", err)
		return "", nil, nil
	}
	// Extract IP
	ip := net.ParseIP(parts[2])
	if ip == nil {
		Log(ERROR, "Failed to parse IP address from introduction packet")
		return "", nil, nil
	}

	return id, mac, ip
}

func (p *PTPCloud) IsPeerUnknown(id string) bool {
	for _, peer := range p.NetworkPeers {
		if peer.ID == id {
			return peer.Unknown
		}
	}
	return true
}

// Handler for new messages received from P2P network
func (p *PTPCloud) HandleP2PMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		Log(ERROR, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := P2PMessageFromBytes(buf)
	if des_err != nil {
		Log(ERROR, "P2PMessageFromBytes error: %v", des_err)
		return
	}
	//var msgType MSG_TYPE = MSG_TYPE(msg.Header.Type)
	// Decrypt message if crypter is active
	if p.Crypter.Active && (msg.Header.Type == MT_INTRO || msg.Header.Type == MT_NENC || msg.Header.Type == MT_INTRO_REQ) {
		var dec_err error
		msg.Data, dec_err = p.Crypter.Decrypt(p.Crypter.ActiveKey.Key, msg.Data)
		if dec_err != nil {
			Log(ERROR, "Failed to decrypt message")
		}
		msg.Data = msg.Data[:msg.Header.Length]
	}
	callback, exists := p.MessageHandlers[msg.Header.Type]
	if exists {
		callback(msg, src_addr)
	} else {
		Log(WARNING, "Unknown message received")
	}
}

func (p *PTPCloud) HandleMessage(msg *P2PMessage, src_addr *net.UDPAddr) {

}

func (p *PTPCloud) HandleNotEncryptedMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	Log(DEBUG, "Received P2P Message")
	Log(TRACE, "Data: %s, Proto: %d, From: %s", msg.Data, msg.Header.NetProto, src_addr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
}

func (p *PTPCloud) HandlePingMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	p.UDPSocket.SendMessage(msg, src_addr)
}

func (p *PTPCloud) HandleXpeerPingMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	pt := PingType(msg.Data)
	if pt == PING_REQ {
		// Send a PING response
		msg := CreateXpeerPingMessage(PING_RESP)
		p.UDPSocket.SendMessage(msg, src_addr)
	} else {
		// Handle PING response
		for i, peer := range p.NetworkPeers {
			if peer.PeerAddr == src_addr {
				peer.WaitingPing = false
				p.NetworkPeers[i] = peer
			}
		}
	}
}

func (p *PTPCloud) IsPeerReady(id string) bool {
	for _, peer := range p.NetworkPeers {
		if peer.ID == id {
			return peer.Ready
		}
	}
	return false
}

func (p *PTPCloud) HandleIntroMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	Log(DEBUG, "Introduction message received: %s", string(msg.Data))
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	// Don't do anything if we already know everything about this peer
	if !p.IsPeerReady(id) {
		Log(DEBUG, "Introduction will be skipped - peer is not ready")
		return
	}
	if !p.IsPeerUnknown(id) {
		Log(DEBUG, "Skipping known peer")
		return
	}
	addr := src_addr
	// TODO: Change PeerAddr with DST addr of real peer
	p.AddPeer(addr, id, ip, mac)
	Log(INFO, "Introduced new peer. IP: %s. ID: %s, HW: %s", ip, id, mac)
}

func (p *PTPCloud) HandleIntroRequestMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	id := string(msg.Data)
	peer, exists := p.NetworkPeers[id]
	if !exists {
		Log(DEBUG, "Introduction request came from unknown peer: %s", id)
		return
	}
	response := p.PrepareIntroductionMessage(p.Dht.ID)
	response.Header.ProxyId = uint16(peer.ProxyID)
	_, err := p.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		Log(ERROR, "Failed to respond to introduction request: %v", err)
	}
}

func (p *PTPCloud) HandleProxyMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	// Proxy registration data
	Log(DEBUG, "Proxy confirmation received")
	if msg.Header.ProxyId < 1 {
		return
	}
	id := string(msg.Data)
	for key, peer := range p.NetworkPeers {
		if peer.PeerAddr.String() == id {
			peer.ProxyID = int(msg.Header.ProxyId)
			Log(DEBUG, "Setting proxy ID %d for %s", msg.Header.ProxyId, peer.ID)
			peer.Ready = true
			peer.ProxyRetries = 0
			peer.State = P_HANDSHAKING
			p.NetworkPeers[key] = peer
		}
	}
}

func (p *PTPCloud) HandleBadTun(msg *P2PMessage, src_addr *net.UDPAddr) {
	Log(INFO, "Cleaning bad tunnel with ID: %d", msg.Header.ProxyId)
	for key, peer := range p.NetworkPeers {
		if peer.ProxyID == int(msg.Header.ProxyId) {
			peer.ProxyID = 0
			peer.Endpoint = ""
			peer.Forwarder = nil
			peer.State = P_INIT
			p.NetworkPeers[key] = peer
		}
	}
}

func (p *PTPCloud) HandleTestMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	response := CreateTestP2PMessage(p.Crypter, "TEST", 0)
	_, err := p.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		Log(ERROR, "Failed to respond to test message: %v", err)
	}

}

func (p *PTPCloud) SendTo(dst net.HardwareAddr, msg *P2PMessage) (int, error) {
	// TODO: Speed up this by switching to map
	Log(TRACE, "Requested Send to %s", dst.String())
	id, exists := p.MACIDTable[dst.String()]
	if exists {
		peer, exists := p.NetworkPeers[id]
		if exists {
			msg.Header.ProxyId = uint16(peer.ProxyID)
			Log(TRACE, "Sending to %s via proxy id %d", dst.String(), msg.Header.ProxyId)
			size, err := p.UDPSocket.SendMessage(msg, peer.PeerAddr)
			return size, err
		}
	}
	return 0, nil
}

func (p *PTPCloud) StopInstance() {
	// Send a packet
	p.Dht.Stop()
	p.UDPSocket.Stop()
	p.Shutdown = true
	// Tricky part: we need to send a message to ourselves to quit blocking operation
	msg := CreateTestP2PMessage(p.Crypter, "STOP", 1)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", p.Dht.P2PPort))
	p.UDPSocket.SendMessage(msg, addr)
	var ipIt int = 200
	for p.IsDeviceExists(p.DeviceName) {
		time.Sleep(1 * time.Second)
		ip := p.Dht.Network.IP
		target := fmt.Sprintf("%d.%d.%d.%d:99", ip[0], ip[1], ip[2], ipIt)
		Log(INFO, "Dialing %s", target)
		_, err := net.DialTimeout("tcp", target, 2*time.Second)
		if err != nil {
			Log(INFO, "ERROR: %v", err)
		}
		ipIt++
		if ipIt == 255 {
			break
		}
	}
	time.Sleep(3 * time.Second)
	p.ReadyToStop = true
}
