package main

import (
	"crypto/rand"
	"fmt"
	ptp "github.com/subutai-io/p2p/lib"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

type MSG_TYPE uint16

type MessageHandler func(message *ptp.P2PMessage, src_addr *net.UDPAddr)

// Main structure
type PTPCloud struct {
	IP              string                               // Interface IP address
	Mac             string                               // String representation of a MAC address
	HardwareAddr    net.HardwareAddr                     // MAC address of network interface
	Mask            string                               // Network mask in the dot-decimal notation
	DeviceName      string                               // Name of the network interface
	IPTool          string                               `yaml:"iptool"` // Network interface configuration tool
	Device          *ptp.Interface                       // Network interface
	NetworkPeers    map[string]NetworkPeer               // Knows peers
	UDPSocket       *ptp.PTPNet                          // Peer-to-peer interconnection socket
	LocalIPs        []net.IP                             // List of IPs available in the system
	dht             *ptp.DHTClient                       // DHT Client
	Crypter         ptp.Crypto                           // Instance of crypto
	Shutdown        bool                                 // Set to true when instance in shutdown mode
	IPIDTable       map[string]string                    // Mapping for IP->ID
	MACIDTable      map[string]string                    // Mapping for MAC->ID
	ForwardMode     bool                                 // Skip local peer discovery
	MessageHandlers map[uint16]MessageHandler            // Callbacks
	ReadyToStop     bool                                 // Set to true when instance is ready to stop
	PacketHandlers  map[PacketType]PacketHandlerCallback // Callbacks for network packet handlers
	// Interface       *os.File
}

type NetworkPeer struct {
	ID           string           // ID of a peer
	Unknown      bool             // TODO: Remove after moving to states
	Handshaked   bool             // TODO: Remove after moving to states
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
	State        ptp.PeerState    // State of a peer
}

// Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PTPCloud) CreateDevice(ip, mac, mask, device string) error {
	var err error

	p.IP = ip
	p.Mac = mac
	p.Mask = mask
	p.DeviceName = device

	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to load config: %v", err)
		p.IPTool = "/sbin/ip"
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to parse config: %v", err)
		return err
	}

	p.Device, err = ptp.Open(p.DeviceName, ptp.DevTap, true)
	if p.Device == nil {
		ptp.Log(ptp.ERROR, "Failed to open TAP device: %v", err)
		return err
	} else {
		ptp.Log(ptp.INFO, "%v TAP Device created", p.DeviceName)
	}

	err = ptp.ConfigureInterface(p.IP, mac, p.DeviceName, p.IPTool)
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
		ptp.Log(ptp.WARNING, "Captured undefined packet")
	}
}

// Listen TAP interface for incoming packets
func (p *PTPCloud) ListenInterface() {
	// Read packets received by TUN/TAP device and send them to a handlePacket goroutine
	// This goroutine will decide what to do with this packet
	for {
		if p.Shutdown {
			break
		}
		packet, err := p.Device.ReadPacket()
		if err != nil {
			ptp.Log(ptp.ERROR, "Reading packet %s", err)
		}
		if packet.Truncated {
			ptp.Log(ptp.DEBUG, "Truncated packet")
		}
		// TODO: Make handlePacket as a part of PTPCloud
		go p.handlePacket(packet.Packet, packet.Protocol)
	}
	p.Device.Close()
	ptp.Log(ptp.INFO, "Shutting down interface listener")
}

// This method will generate device name if none were specified at startup
func (p *PTPCloud) GenerateDeviceName(i int) string {
	var devName string = "vptp" + fmt.Sprintf("%d", i)
	inf, err := net.Interfaces()
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to retrieve list of network interfaces")
		return ""
	}
	var exist bool = false
	for _, i := range inf {
		if i.Name == devName {
			exist = true
		}
	}
	if exist {
		return p.GenerateDeviceName(i + 1)
	} else {
		return devName
	}
}

// This method lists interfaces available in the system and retrieves their
// IP addresses
func (p *PTPCloud) FindNetworkAddresses() {
	ptp.Log(ptp.INFO, "Looking for available network interfaces")
	inf, err := net.Interfaces()
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to retrieve list of network interfaces")
		return
	}
	for _, i := range inf {
		addresses, err := i.Addrs()

		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to retrieve address for interface. %v", err)
			continue
		}
		for _, addr := range addresses {
			var decision string = "Ignoring"
			var ipType string = "Unknown"
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				ptp.Log(ptp.ERROR, "Failed to parse CIDR notation: %v", err)
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
			ptp.Log(ptp.INFO, "Interface %s: %s. Type: %s. %s", i.Name, addr.String(), ipType, decision)
			if decision == "Saving" {
				p.LocalIPs = append(p.LocalIPs, ip)
			}
		}
	}
	ptp.Log(ptp.INFO, "%d interfaces were saved", len(p.LocalIPs))
}

func p2pmain(argIp, argMask, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool) *PTPCloud {

	var hw net.HardwareAddr

	if argMac != "" {
		var err2 error
		hw, err2 = net.ParseMAC(argMac)
		if err2 != nil {
			ptp.Log(ptp.ERROR, "Invalid MAC address provided: %v", err2)
			return nil
		}
	} else {
		argMac, hw = GenerateMAC()
		ptp.Log(ptp.INFO, "Generate MAC for TAP device: %s", argMac)
	}

	// Create new DHT Client, configured it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code
	dhtClient := new(ptp.DHTClient)
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = argHash
	config.Mode = ptp.MODE_CLIENT

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
	}

	if argKeyfile != "" {
		p.Crypter.ReadKeysFromFile(argKeyfile)
	}
	if argKey != "" {
		// Override key from file
		if argTTL == "" {
			argTTL = "default"
		}
		var newKey ptp.CryptoKey
		newKey = p.Crypter.EnrichKeyValues(newKey, argKey, argTTL)
		p.Crypter.Keys = append(p.Crypter.Keys, newKey)
		p.Crypter.ActiveKey = p.Crypter.Keys[0]
		p.Crypter.Active = true
	}

	if p.Crypter.Active {
		ptp.Log(ptp.INFO, "Traffic encryption is enabled. Key valid until %s", p.Crypter.ActiveKey.Until.String())
	} else {
		ptp.Log(ptp.INFO, "No AES key were provided. Traffic encryption is disabled")
	}

	// Register network message handlers
	p.MessageHandlers = make(map[uint16]MessageHandler)
	p.MessageHandlers[ptp.MT_NENC] = p.HandleNotEncryptedMessage
	p.MessageHandlers[ptp.MT_PING] = p.HandlePingMessage
	p.MessageHandlers[ptp.MT_ENC] = p.HandleMessage
	p.MessageHandlers[ptp.MT_INTRO] = p.HandleIntroMessage
	p.MessageHandlers[ptp.MT_INTRO_REQ] = p.HandleIntroRequestMessage
	p.MessageHandlers[ptp.MT_PROXY] = p.HandleProxyMessage
	p.MessageHandlers[ptp.MT_TEST] = p.HandleTestMessage

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

	p.CreateDevice(argIp, argMac, argMask, argDev)
	p.UDPSocket = new(ptp.PTPNet)
	p.UDPSocket.Init("", 0)
	port := p.UDPSocket.GetPort()
	ptp.Log(ptp.INFO, "Started UDP Listener at port %d", port)
	config.P2PPort = port
	if argDht != "" {
		config.Routers = argDht
	}
	p.dht = dhtClient.Initialize(config, p.LocalIPs)

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
		p.dht.UpdatePeers()
		// Wait two seconds before synchronizing with catched peers
		time.Sleep(2 * time.Second)
		p.PurgePeers()
		newPeersNum := p.SyncPeers()
		newPeersNum = newPeersNum + p.SyncForwarders()
		//if newPeersNum > 0 {
		p.IntroducePeers()
		//}
	}
	ptp.Log(ptp.INFO, "Shutting down instance %s completed", p.dht.NetworkHash)
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
				ptp.Log(ptp.WARNING, "Failed to introduce to %s via proxy %d", peer.ID, peer.ProxyID)
				p.dht.MakeForwarderFailed(peer.Forwarder)
				peer.Endpoint = ""
				peer.Forwarder = nil
				peer.ProxyID = 0
				peer.State = ptp.P_INIT
				p.NetworkPeers[i] = peer
			} else {
				ptp.Log(ptp.WARNING, "Failed to introduce to %s", peer.ID)
				peer.State = ptp.P_HANDSHAKING_FAILED
				p.NetworkPeers[i] = peer
			}
			continue
		}
		peer.Retries = peer.Retries + 1
		ptp.Log(ptp.DEBUG, "Intoducing to %s", peer.Endpoint)
		addr, err := net.ResolveUDPAddr("udp", peer.Endpoint)
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to resolve UDP address during Introduction: %v", err)
			continue
		}
		//peer.PeerAddr = addr
		p.NetworkPeers[i] = peer
		// Send introduction packet
		//msg := p.PrepareIntroductionMessage(p.dht.ID)
		msg := ptp.CreateIntroRequest(p.Crypter, p.dht.ID)
		msg.Header.ProxyId = uint16(peer.ProxyID)
		_, err = p.UDPSocket.SendMessage(msg, addr)
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to send introduction to %s", addr.String())
		} else {
			ptp.Log(ptp.DEBUG, "Introduction sent to %s", peer.Endpoint)
		}
	}
}

func (p *PTPCloud) PrepareIntroductionMessage(id string) *ptp.P2PMessage {
	var intro string = id + "," + p.Mac + "," + p.IP
	msg := ptp.CreateIntroP2PMessage(p.Crypter, intro, 0)
	return msg
}

// This method goes over peers and removes obsolete ones
// Peer becomes obsolete when it goes out of DHT
func (p *PTPCloud) PurgePeers() {
	for i, peer := range p.NetworkPeers {
		var f bool = false
		for _, newPeer := range p.dht.Peers {
			if newPeer.ID == peer.ID {
				f = true
			}
		}
		if !f {
			ptp.Log(ptp.INFO, ("Removing outdated peer"))
			delete(p.IPIDTable, peer.PeerLocalIP.String())
			delete(p.MACIDTable, peer.PeerHW.String())
			delete(p.NetworkPeers, i)
		}
	}
	return
}

// This method tests connection with specified endpoint
func (p *PTPCloud) TestConnection(endpoint *net.UDPAddr) bool {
	msg := ptp.CreateTestP2PMessage(p.Crypter, "TEST", 0)
	conn, err := net.DialUDP("udp4", nil, endpoint)
	if err != nil {
		ptp.Log(ptp.ERROR, "%v", err)
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
			ptp.Log(ptp.ERROR, "%v", err)
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
	for _, fwd := range p.dht.Forwarders {
		for key, peer := range p.NetworkPeers {
			if peer.Endpoint == "" && fwd.DestinationID == peer.ID && peer.Forwarder == nil {
				ptp.Log(ptp.INFO, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr.String()
				peer.Forwarder = fwd.Addr
				peer.State = ptp.P_HANDSHAKING_FORWARDER
				p.NetworkPeers[key] = peer
				count = count + 1
			}
		}
	}
	p.dht.Forwarders = p.dht.Forwarders[:0]
	return count
}

func (p *PTPCloud) AssignEndpoint(peer NetworkPeer) (string, bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to retrieve list of network interfaces")
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
				ptp.Log(ptp.TRACE, "Probing new IP %s against network %s", kip.IP.String(), network.String())

				if network.Contains(kip.IP) {
					if p.TestConnection(kip) {
						return kip.String(), true
						ptp.Log(ptp.INFO, "Setting endpoint for %s to %s", peer.ID, kip.String())
					}
					// TODO: Test connection
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

	for _, id := range p.dht.Peers {
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
						ptp.Log(ptp.INFO, "Adding new IP (%s) address to %s", ip, peer.ID)
						// TODO: Check IP parsing
						newIp, _ := net.ResolveUDPAddr("udp", ip)
						peer.KnownIPs = append(peer.KnownIPs, newIp)
						p.NetworkPeers[i] = peer
					}
				}

				// Set and Endpoint from peers if no endpoint were set previously
				if peer.State == ptp.P_INIT {
					var added bool
					if !p.ForwardMode {
						peer.Endpoint, added = p.AssignEndpoint(peer)
						peer.State = ptp.P_CONNECTED
					}
					if added {
						p.NetworkPeers[i] = peer
						count += 1
					} else {
						if len(peer.KnownIPs) > 0 {
							if !p.ForwardMode {
								ptp.Log(ptp.INFO, "No peers are available in local network. Switching to global IP if any")
							}
							// If endpoint wasn't set let's test connection from outside of the LAN
							// First one should be the global IP (if DHT works correctly)
							if p.ForwardMode || !p.TestConnection(p.NetworkPeers[i].KnownIPs[0]) {
								if peer.State == ptp.P_WAITING_FORWARDER {
									continue
								}
								// We've failed to establish connection again. Now let's ask for a proxy
								ptp.Log(ptp.INFO, "Requesting Control Peer from Service Discovery Peer")
								p.dht.RequestControlPeer(peer.ID)
								peer.PeerAddr = peer.KnownIPs[0]
								peer.State = ptp.P_WAITING_FORWARDER
								p.NetworkPeers[i] = peer
							} else {
								ptp.Log(ptp.INFO, "Successfully connected to a host over Internet")
								peer.Endpoint = peer.KnownIPs[0].String()
								peer.State = ptp.P_CONNECTED
								p.NetworkPeers[i] = peer
							}
						}
					}

				} else if peer.State == ptp.P_HANDSHAKING_FAILED {
					ptp.Log(ptp.INFO, "Requesting backup control peer")
					p.dht.RequestControlPeer(peer.ID)
					peer.PeerAddr = peer.KnownIPs[0]
					peer.State = ptp.P_WAITING_FORWARDER
					p.NetworkPeers[i] = peer
				} else if peer.State == ptp.P_HANDSHAKING_FORWARDER {
					if peer.ProxyRetries > 3 {
						ptp.Log(ptp.WARNING, "Failed to handshake with control peer %s. Adding to black list", peer.Forwarder.String())
						p.dht.MakeForwarderFailed(peer.Forwarder)
						peer.Forwarder = nil
						peer.Endpoint = ""
						peer.ProxyID = 0
						peer.State = ptp.P_INIT
						continue
					}
					peer.ProxyRetries = peer.ProxyRetries + 1
					p.NetworkPeers[i] = peer
					ptp.Log(ptp.INFO, "Sending proxy request to a forwarder %s", peer.Forwarder.String())
					msg := ptp.CreateProxyP2PMessage(-1, peer.PeerAddr.String(), 0)
					_, err := p.UDPSocket.SendMessage(msg, peer.Forwarder)
					if err != nil {
						ptp.Log(ptp.ERROR, "Failed to send a message to a proxy %v", err)
					}
				}
			}
		}
		if !found {
			ptp.Log(ptp.INFO, "Adding new peer. Requesting peer address")
			var newPeer NetworkPeer
			newPeer.ID = id.ID
			newPeer.Unknown = true
			newPeer.Ready = true
			newPeer.State = ptp.P_INIT
			if p.ForwardMode {
				newPeer.Ready = false
			}
			p.NetworkPeers[newPeer.ID] = newPeer
			p.dht.RequestPeerIPs(id.ID)
		}
	}
	return count
}

// WriteToDevice writes data to created TUN/TAP device
func (p *PTPCloud) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var packet ptp.Packet
	packet.Protocol = int(proto)
	packet.Truncated = truncated
	packet.Packet = b
	if p.Device == nil {
		ptp.Log(ptp.ERROR, "TUN/TAP Device not initialized")
		return
	}
	err := p.Device.WritePacket(&packet)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to write to TUN/TAP device: %v", err)
	}
}

func GenerateMAC() (string, net.HardwareAddr) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to generate MAC: %v", err)
		return "", nil
	}
	buf[0] |= 2
	mac := fmt.Sprintf("06:%02x:%02x:%02x:%02x:%02x", buf[1], buf[2], buf[3], buf[4], buf[5])
	hw, err := net.ParseMAC(mac)
	if err != nil {
		ptp.Log(ptp.ERROR, "Corrupted MAC address generated: %v", err)
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
		ptp.Log(ptp.ERROR, "Failed to parse introduction string")
		return "", nil, nil
	}
	var id string
	id = parts[0]
	// Extract MAC
	mac, err := net.ParseMAC(parts[1])
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to parse MAC address from introduction packet: %v", err)
		return "", nil, nil
	}
	// Extract IP
	ip := net.ParseIP(parts[2])
	if ip == nil {
		ptp.Log(ptp.ERROR, "Failed to parse IP address from introduction packet")
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
		ptp.Log(ptp.ERROR, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := ptp.P2PMessageFromBytes(buf)
	if des_err != nil {
		ptp.Log(ptp.ERROR, "P2PMessageFromBytes error: %v", des_err)
		return
	}
	//var msgType ptp.MSG_TYPE = ptp.MSG_TYPE(msg.Header.Type)
	// Decrypt message if crypter is active
	if p.Crypter.Active {
		var dec_err error
		msg.Data, dec_err = p.Crypter.Decrypt(p.Crypter.ActiveKey.Key, msg.Data)
		if dec_err != nil {
			ptp.Log(ptp.ERROR, "Failed to decrypt message")
		}
		msg.Data = msg.Data[:msg.Header.Length]
	}
	callback, exists := p.MessageHandlers[msg.Header.Type]
	if exists {
		callback(msg, src_addr)
	} else {
		ptp.Log(ptp.WARNING, "Unknown message received")
	}
}

func (p *PTPCloud) HandleMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {

}

func (p *PTPCloud) HandleNotEncryptedMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	ptp.Log(ptp.DEBUG, "Received P2P Message")
	ptp.Log(ptp.TRACE, "Data: %s, Proto: %d, From: %s", msg.Data, msg.Header.NetProto, src_addr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
}

func (p *PTPCloud) HandlePingMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	p.UDPSocket.SendMessage(msg, src_addr)
}

func (p *PTPCloud) IsPeerReady(id string) bool {
	for _, peer := range p.NetworkPeers {
		if peer.ID == id {
			return peer.Ready
		}
	}
	return false
}

func (p *PTPCloud) HandleIntroMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	ptp.Log(ptp.DEBUG, "Introduction message received: %s", string(msg.Data))
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	// Don't do anything if we already know everything about this peer
	if !p.IsPeerReady(id) {
		ptp.Log(ptp.DEBUG, "Introduction will be skipped - peer is not ready")
		return
	}
	if !p.IsPeerUnknown(id) {
		ptp.Log(ptp.DEBUG, "Skipping known peer")
		return
	}
	addr := src_addr
	// TODO: Change PeerAddr with DST addr of real peer
	p.AddPeer(addr, id, ip, mac)
	ptp.Log(ptp.INFO, "Introduced new peer. IP: %s. ID: %s, HW: %s", ip, id, mac)
}

func (p *PTPCloud) HandleIntroRequestMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	id := string(msg.Data)
	peer, exists := p.NetworkPeers[id]
	if !exists {
		ptp.Log(ptp.DEBUG, "Introduction request came from unknown peer")
		return
	}
	response := p.PrepareIntroductionMessage(p.dht.ID)
	response.Header.ProxyId = uint16(peer.ProxyID)
	_, err := p.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to respond to introduction request: %v", err)
	}
}

func (p *PTPCloud) HandleProxyMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	// Proxy registration data
	ptp.Log(ptp.DEBUG, "Proxy confirmation received")
	if msg.Header.ProxyId < 1 {
		return
	}
	id := string(msg.Data)
	for key, peer := range p.NetworkPeers {
		if peer.PeerAddr.String() == id {
			peer.ProxyID = int(msg.Header.ProxyId)
			ptp.Log(ptp.DEBUG, "Setting proxy ID %d for %s", msg.Header.ProxyId, peer.ID)
			peer.Ready = true
			peer.ProxyRetries = 0
			peer.State = ptp.P_HANDSHAKING
			p.NetworkPeers[key] = peer
		}
	}
}

func (p *PTPCloud) HandleTestMessage(msg *ptp.P2PMessage, src_addr *net.UDPAddr) {
	response := ptp.CreateTestP2PMessage(p.Crypter, "TEST", 0)
	_, err := p.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to respond to test message: %v", err)
	}

}

func (p *PTPCloud) SendTo(dst net.HardwareAddr, msg *ptp.P2PMessage) (int, error) {
	// TODO: Speed up this by switching to map
	ptp.Log(ptp.TRACE, "Requested Send to %s", dst.String())
	id, exists := p.MACIDTable[dst.String()]
	if exists {
		peer, exists := p.NetworkPeers[id]
		if exists {
			msg.Header.ProxyId = uint16(peer.ProxyID)
			ptp.Log(ptp.TRACE, "Sending to %s via proxy id %d", dst.String(), msg.Header.ProxyId)
			size, err := p.UDPSocket.SendMessage(msg, peer.PeerAddr)
			return size, err
		}
	}
	return 0, nil
}

func (p *PTPCloud) StopInstance() {
	// Send a packet
	p.dht.Stop()
	p.UDPSocket.Stop()
	p.Shutdown = true
	// Tricky part: we need to send a message to ourselves to quit blocking operation
	msg := ptp.CreateTestP2PMessage(p.Crypter, "STOP", 1)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", p.dht.P2PPort))
	p.UDPSocket.SendMessage(msg, addr)
	time.Sleep(3 * time.Second)
	p.ReadyToStop = true
}
