package ptp

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// MessageHandler is a messages callback
type MessageHandler func(message *P2PMessage, srcAddr *net.UDPAddr)

// NetworkInterface keeps information about P2P network interface
type NetworkInterface struct {
	IP        net.IP           // IP
	Mask      net.IPMask       // Mask
	Mac       net.HardwareAddr // Hardware Address
	Name      string           // Network interface name
	Interface *Interface       // TAP Interface
}

// PeerToPeer - Main structure
type PeerToPeer struct {
	IPTool          string                               `yaml:"iptool"`  // Network interface configuration tool
	AddTap          string                               `yaml:"addtap"`  // Path to addtap.bat
	InfFile         string                               `yaml:"inffile"` // Path to deltap.bat
	NetworkPeers    map[string]*NetworkPeer              // Knows peers
	UDPSocket       *Network                             // Peer-to-peer interconnection socket
	LocalIPs        []net.IP                             // List of IPs available in the system
	Dht             *DHTClient                           // DHT Client
	Crypter         Crypto                               // Cryptography subsystem
	Shutdown        bool                                 // Set to true when instance in shutdown mode
	ForwardMode     bool                                 // Skip local peer discovery
	ReadyToStop     bool                                 // Set to true when instance is ready to stop
	IPIDTable       map[string]string                    // Mapping for IP->ID
	MACIDTable      map[string]string                    // Mapping for MAC->ID
	MessageHandlers map[uint16]MessageHandler            // Callbacks for network packets
	PacketHandlers  map[PacketType]PacketHandlerCallback // Callbacks for packets received by TAP interface
	PeersLock       sync.Mutex                           // Lock for peers map
	Hash            string                               // Infohash for this instance
	Routers         string                               // Comma-separated list of Bootstrap nodes
	Interface       NetworkInterface                     // TAP Interface
}

// AssignInterface - Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PeerToPeer) AssignInterface(interfaceName string) error {
	var err error
	p.Interface.Name = interfaceName
	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile(ConfigDir + "/p2p/config.yaml")
	if err != nil {
		Log(Warning, "Failed to load config: %v", err)
		p.IPTool = "/sbin/ip"
		p.AddTap = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"
		p.InfFile = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf"
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		Log(Error, "Failed to parse config: %v", err)
		return err
	}

	p.Interface.Interface, err = Open(p.Interface.Name, DevTap)
	if p.Interface.Interface == nil {
		Log(Error, "Failed to open TAP device %s: %v", p.Interface.Name, err)
		return err
	}
	Log(Info, "%v TAP Device created", p.Interface.Name)

	// Windows returns a real mac here. However, other systems should return empty string
	hwaddr := ExtractMacFromInterface(p.Interface.Interface)
	if hwaddr != "" {
		p.Interface.Mac, _ = net.ParseMAC(hwaddr)
	}

	err = ConfigureInterface(p.Interface.Interface, p.Interface.IP.String(), p.Interface.Mac.String(), p.Interface.Name, p.IPTool)
	Log(Info, "Interface has been configured")
	return err
}

// ListenInterface - Listens TAP interface for incoming packets
func (p *PeerToPeer) ListenInterface() {
	// Read packets received by TUN/TAP device and send them to a handlePacket goroutine
	// This goroutine will decide what to do with this packet

	// Run is for windows only
	p.Interface.Interface.Run()
	for {
		if p.Shutdown {
			break
		}
		packet, err := p.Interface.Interface.ReadPacket()
		if err != nil {
			Log(Error, "Reading packet %s", err)
		}
		if packet.Truncated {
			Log(Debug, "Truncated packet")
		}
		go p.handlePacket(packet.Packet, packet.Protocol)
	}
	p.Interface.Interface.Close()
	Log(Info, "Shutting down interface listener")
}

// IsDeviceExists - checks whether interface with the given name exists in the system or not
func (p *PeerToPeer) IsDeviceExists(name string) bool {
	inf, err := net.Interfaces()
	if err != nil {
		Log(Error, "Failed to retrieve list of network interfaces")
		return true
	}
	for _, i := range inf {
		if i.Name == name {
			return true
		}
	}
	return false
}

// GenerateDeviceName method will generate device name if none were specified at startup
func (p *PeerToPeer) GenerateDeviceName(i int) string {
	var devName = GetDeviceBase() + fmt.Sprintf("%d", i)
	if p.IsDeviceExists(devName) {
		return p.GenerateDeviceName(i + 1)
	}
	return devName
}

// IsIPv4 checks whether interface is IPv4 or IPv6
func (p *PeerToPeer) IsIPv4(ip string) bool {
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

// FindNetworkAddresses method lists interfaces available in the system and retrieves their
// IP addresses
func (p *PeerToPeer) FindNetworkAddresses() {
	Log(Info, "Looking for available network interfaces")
	inf, err := net.Interfaces()
	if err != nil {
		Log(Error, "Failed to retrieve list of network interfaces")
		return
	}
	for _, i := range inf {
		addresses, err := i.Addrs()

		if err != nil {
			Log(Error, "Failed to retrieve address for interface. %v", err)
			continue
		}
		for _, addr := range addresses {
			var decision = "Ignoring"
			var ipType = "Unknown"
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				Log(Error, "Failed to parse CIDR notation: %v", err)
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
			Log(Info, "Interface %s: %s. Type: %s. %s", i.Name, addr.String(), ipType, decision)
			if decision == "Saving" {
				p.LocalIPs = append(p.LocalIPs, ip)
			}
		}
	}
	Log(Info, "%d interfaces were saved", len(p.LocalIPs))
}

// StartP2PInstance is an entry point of a P2P library.
func StartP2PInstance(argIP, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool, port int, ignoreIPs []string) *PeerToPeer {
	p := new(PeerToPeer)
	p.Init()
	p.Interface.Mac = p.validateMac(argMac)
	p.FindNetworkAddresses()
	interfaceName, err := p.validateInterfaceName(argDev)
	if err != nil {
		Log(Error, "Interface name validation failed: %s", err)
		return nil
	}
	if p.IsDeviceExists(interfaceName) {
		Log(Error, "Interface is already in use. Can't create duplicate")
		return nil
	}

	if fwd {
		p.ForwardMode = true
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
		Log(Info, "Traffic encryption is enabled. Key valid until %s", p.Crypter.ActiveKey.Until.String())
	} else {
		Log(Info, "No AES key were provided. Traffic encryption is disabled")
	}

	p.setupHandlers()

	p.UDPSocket = new(Network)
	p.UDPSocket.Init("", port)

	// Create new DHT Client, configure it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code
	Log(Info, "Started UDP Listener at port %d", p.UDPSocket.GetPort())
	p.Hash = argHash
	p.StartDHT(p.Hash, argDht)
	p.Routers = p.Dht.Routers
	if argIP == "dhcp" {
		ipn, maskn, err := p.RequestIP(p.Interface.Mac.String(), interfaceName)
		if err != nil {
			Log(Error, "%v", err)
			return nil
		}
		p.Interface.IP = ipn
		p.Interface.Mask = maskn
	} else {
		p.Interface.IP = net.ParseIP(argIP)
		ipn, maskn, err := p.ReportIP(argIP, p.Interface.Mac.String(), interfaceName)
		if err != nil {
			Log(Error, "%v", err)
			return nil
		}
		p.Interface.IP = ipn
		p.Interface.Mask = maskn
	}

	go p.UDPSocket.Listen(p.HandleP2PMessage)
	go p.ListenInterface()
	return p
}

// Init will initialize PeerToPeer
func (p *PeerToPeer) Init() {
	p.NetworkPeers = make(map[string]*NetworkPeer)
	p.IPIDTable = make(map[string]string)
	p.MACIDTable = make(map[string]string)
}

func (p *PeerToPeer) validateMac(mac string) net.HardwareAddr {
	var hw net.HardwareAddr
	var err error
	if mac != "" {
		hw, err = net.ParseMAC(mac)
		if err != nil {
			Log(Error, "Invalid MAC address provided: %v", err)
			return nil
		}
	} else {
		mac, hw = GenerateMAC()
		Log(Info, "Generate MAC for TAP device: %s", mac)
	}
	return hw
}

func (p *PeerToPeer) validateInterfaceName(name string) (string, error) {
	if name == "" {
		name = p.GenerateDeviceName(1)
	} else {
		if len(name) > 12 {
			Log(Info, "Interface name length should be 12 symbols max")
			return "", fmt.Errorf("Interface name is too big")
		}
	}
	return name, nil
}

func (p *PeerToPeer) setupHandlers() {
	// Register network message handlers
	p.MessageHandlers = make(map[uint16]MessageHandler)
	p.MessageHandlers[MsgTypeNenc] = p.HandleNotEncryptedMessage
	p.MessageHandlers[MsgTypePing] = p.HandlePingMessage
	p.MessageHandlers[MsgTypeXpeerPing] = p.HandleXpeerPingMessage
	p.MessageHandlers[MsgTypeIntro] = p.HandleIntroMessage
	p.MessageHandlers[MsgTypeIntroReq] = p.HandleIntroRequestMessage
	p.MessageHandlers[MsgTypeProxy] = p.HandleProxyMessage
	p.MessageHandlers[MsgTypeTest] = p.HandleTestMessage
	p.MessageHandlers[MsgTypeBadTun] = p.HandleBadTun

	// Register packet handlers
	p.PacketHandlers = make(map[PacketType]PacketHandlerCallback)
	p.PacketHandlers[PacketPARCUniversal] = p.handlePARCUniversalPacket
	p.PacketHandlers[PacketIPv4] = p.handlePacketIPv4
	p.PacketHandlers[PacketARP] = p.handlePacketARP
	p.PacketHandlers[PacketRARP] = p.handleRARPPacket
	p.PacketHandlers[Packet8021Q] = p.handle8021qPacket
	p.PacketHandlers[PacketIPv6] = p.handlePacketIPv6
	p.PacketHandlers[PacketPPPoEDiscovery] = p.handlePPPoEDiscoveryPacket
	p.PacketHandlers[PacketPPPoESession] = p.handlePPPoESessionPacket
	p.PacketHandlers[PacketLLDP] = p.handlePacketLLDP
}

// RequestIP asks DHT to get IP from DHCP-like service
func (p *PeerToPeer) RequestIP(mac, device string) (net.IP, net.IPMask, error) {
	Log(Info, "Requesting IP")
	p.Dht.RequestIP()
	time.Sleep(1 * time.Second)
	retries := 0
	for p.Dht.IP == nil && p.Dht.Network == nil {
		Log(Info, "No IP were received. Requesting again")
		p.Dht.RequestIP()
		time.Sleep(3 * time.Second)
		retries++
		if retries >= 10 {
			return nil, nil, fmt.Errorf("Failed to retrieve IP from network after 10 retries")
		}
	}
	err := p.AssignInterface(device)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to configure interface: %s", err)
	}
	return p.Dht.IP, p.Dht.Network.Mask, nil
}

// ReportIP will send IP specified at service start to DHCP-like service
func (p *PeerToPeer) ReportIP(ipAddress, mac, device string) (net.IP, net.IPMask, error) {
	ip, ipnet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		nip := net.ParseIP(ipAddress)
		if nip == nil {
			return nil, nil, fmt.Errorf("Invalid address were provided for network interface. Use -ip \"dhcp\" or specify correct IP address")
		}
		ipAddress += `/24`
		Log(Warning, "No CIDR mask was provided. Assumming /24")
		ip, ipnet, err = net.ParseCIDR(ipAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to setup provided IP address for local device")
		}
	}
	p.Dht.IP = ip
	p.Dht.Network = ipnet
	mask := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
	p.Dht.SendIP(ipAddress, mask)
	err = p.AssignInterface(device)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to configure interface: %s", err)
	}
	return ip, ipnet.Mask, nil
}

// StartDHT starts a DHT client
func (p *PeerToPeer) StartDHT(hash, routers string) error {
	if p.Dht != nil {
		Log(Info, "Stopping previous DHT instance")
		p.Dht.Shutdown()
		p.Dht = nil
	}
	p.Dht = new(DHTClient)
	err := p.Dht.Init(hash, routers)
	if err != nil {
		return fmt.Errorf("Failed to initialize DHT: %s", err)
	}
	p.Dht.setupCallbacks()
	p.Dht.IPList = p.LocalIPs
	err = p.Dht.Connect()
	if err != nil {
		Log(Error, "Failed to establish connection with Bootstrap node: %s")
		for err != nil {
			Log(Warning, "Retrying connection")
			err = p.Dht.Connect()
			time.Sleep(3 * time.Second)
		}
	}
	err = p.Dht.WaitForID()
	if err != nil {
		Log(Error, "Failed to retrieve ID from bootstrap node: %s", err)
	}
	return nil
}

func (p *PeerToPeer) markPeerForRemoval(id, reason string) error {
	p.PeersLock.Lock()
	peer, exists := p.NetworkPeers[id]
	p.PeersLock.Unlock()
	runtime.Gosched()
	if exists {
		Log(Info, "Removing peer %s: Reason %s", id, reason)
		peer.State = PeerStateDisconnect
		p.PeersLock.Lock()
		p.NetworkPeers[id] = peer
		p.PeersLock.Unlock()
		runtime.Gosched()
	} else {
		return fmt.Errorf("Peer not found")
	}
	return nil
}

// Run is a main loop
func (p *PeerToPeer) Run() {
	go p.ReadDHTPeers()
	go p.ReadProxies()
	go func() {
		for {
			if p.Shutdown {
				break
			}
			select {
			// Handle STOP Command from DHT
			case rm, r := <-p.Dht.RemovePeerChan:
				if r {
					if rm == "DUMMY" || rm == "" {
						continue
					}
					err := p.markPeerForRemoval(rm, "Stop")
					if err != nil {
						Log(Error, "Failed to mark peer for removal: %s", err)
					}
				}
			default:
				time.Sleep(100 * time.Millisecond)
			}
		}
		Log(Info, "Stopping peer state listener")
	}()
	go p.Dht.UpdatePeers()
	for {
		if p.Shutdown {
			// TODO: Do it more safely
			if p.ReadyToStop {
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(time.Second * 1)
		for i, peer := range p.NetworkPeers {
			if peer.State == PeerStateStop {
				Log(Info, "Removing peer %s", i)
				time.Sleep(100 * time.Millisecond)
				lip := peer.PeerLocalIP.String()
				if peer.ID == p.IPIDTable[lip] {
					delete(p.IPIDTable, lip)
				}
				delete(p.MACIDTable, peer.PeerHW.String())
				delete(p.NetworkPeers, i)
				err := p.Dht.CleanPeer(i)
				if err != nil {
					Log(Error, "Failed to remove peer from DHT: %s", err)
				}
				//runtime.Gosched()
				Log(Info, "Remove complete")
				break
			}
		}
		passed := time.Since(p.Dht.LastDHTPing)
		interval := time.Duration(time.Second * 45)
		if passed > interval {
			Log(Error, "Lost connection to DHT")
			time.Sleep(time.Second * 3)
			p.StartDHT(p.Hash, p.Routers)
			p.Dht.SendIP(p.Interface.IP.To4().String(), p.Interface.Mask.String())
			go p.Dht.UpdatePeers()
		}
	}
	Log(Info, "Shutting down instance %s completed", p.Dht.NetworkHash)
}

// PrepareIntroductionMessage collects client ID, mac and IP address
// and create a comma-separated line
func (p *PeerToPeer) PrepareIntroductionMessage(id string) *P2PMessage {
	var intro = id + "," + p.Interface.Mac.String() + "," + p.Interface.IP.String()
	msg := CreateIntroP2PMessage(p.Crypter, intro, 0)
	return msg
}

// PurgePeers method goes over peers and removes obsolete ones
// Peer becomes obsolete when it goes out of DHT
func (p *PeerToPeer) PurgePeers() {
	for i, peer := range p.NetworkPeers {
		var f = false
		for _, newPeer := range p.Dht.Peers {
			if newPeer.ID == peer.ID {
				f = true
			}
		}
		if !f {
			Log(Info, ("Removing outdated peer"))
			delete(p.IPIDTable, peer.PeerLocalIP.String())
			delete(p.MACIDTable, peer.PeerHW.String())
			p.PeersLock.Lock()
			delete(p.NetworkPeers, i)
			p.PeersLock.Unlock()
			runtime.Gosched()
		}
	}
	return
}

// SyncForwarders extracts proxies from DHT and assign them to target peers
func (p *PeerToPeer) SyncForwarders() int {
	var count = 0
	for _, fwd := range p.Dht.Forwarders {
		for key, peer := range p.NetworkPeers {
			if peer.Endpoint == nil && fwd.DestinationID == peer.ID && peer.Forwarder == nil {
				Log(Info, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr
				peer.Forwarder = fwd.Addr
				peer.State = PeerStateHandshakingForwarder
				p.PeersLock.Lock()
				p.NetworkPeers[key] = peer
				p.PeersLock.Unlock()
				runtime.Gosched()
				count = count + 1
			}
		}
	}
	p.Dht.Forwarders = p.Dht.Forwarders[:0]
	return count
}

// WriteToDevice writes data to created TAP interface
func (p *PeerToPeer) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var packet Packet
	packet.Protocol = int(proto)
	packet.Truncated = truncated
	packet.Packet = b
	if p.Interface.Interface == nil {
		Log(Error, "TAP Interface not initialized")
		return
	}
	err := p.Interface.Interface.WritePacket(&packet)
	if err != nil {
		Log(Error, "Failed to write to TAP Interface: %v", err)
	}
}

// GenerateMAC generates a MAC address for a new interface
func GenerateMAC() (string, net.HardwareAddr) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		Log(Error, "Failed to generate MAC: %v", err)
		return "", nil
	}
	buf[0] |= 2
	mac := fmt.Sprintf("06:%02x:%02x:%02x:%02x:%02x", buf[1], buf[2], buf[3], buf[4], buf[5])
	hw, err := net.ParseMAC(mac)
	if err != nil {
		Log(Error, "Corrupted MAC address generated: %v", err)
		return "", nil
	}
	return mac, hw
}

// ParseIntroString receives a comma-separated string with ID, MAC and IP of a peer
// and returns this data
func (p *PeerToPeer) ParseIntroString(intro string) (string, net.HardwareAddr, net.IP) {
	parts := strings.Split(intro, ",")
	if len(parts) != 3 {
		Log(Error, "Failed to parse introduction string: %s", intro)
		return "", nil, nil
	}
	var id string
	id = parts[0]
	// Extract MAC
	mac, err := net.ParseMAC(parts[1])
	if err != nil {
		Log(Error, "Failed to parse MAC address from introduction packet: %v", err)
		return "", nil, nil
	}
	// Extract IP
	ip := net.ParseIP(parts[2])
	if ip == nil {
		Log(Error, "Failed to parse IP address from introduction packet")
		return "", nil, nil
	}

	return id, mac, ip
}

// HandleP2PMessage is a handler for new messages received from P2P network
func (p *PeerToPeer) HandleP2PMessage(count int, srcAddr *net.UDPAddr, err error, rcvBytes []byte) {
	if err != nil {
		Log(Error, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcvBytes[:])

	msg, desErr := P2PMessageFromBytes(buf)
	if desErr != nil {
		Log(Error, "P2PMessageFromBytes error: %v", desErr)
		return
	}
	//var msgType MSG_TYPE = MSG_TYPE(msg.Header.Type)
	// Decrypt message if crypter is active
	if p.Crypter.Active && (msg.Header.Type == MsgTypeIntro || msg.Header.Type == MsgTypeNenc || msg.Header.Type == MsgTypeIntroReq) {
		var decErr error
		msg.Data, decErr = p.Crypter.Decrypt(p.Crypter.ActiveKey.Key, msg.Data)
		if decErr != nil {
			Log(Error, "Failed to decrypt message")
		}
		msg.Data = msg.Data[:msg.Header.Length]
	}
	callback, exists := p.MessageHandlers[msg.Header.Type]
	if exists {
		callback(msg, srcAddr)
	} else {
		Log(Warning, "Unknown message received")
	}
}

// HandleNotEncryptedMessage is a normal message sent over p2p network
func (p *PeerToPeer) HandleNotEncryptedMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Trace, "Data: %s, Proto: %d, From: %s", msg.Data, msg.Header.NetProto, srcAddr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
}

// HandlePingMessage is a PING message from a proxy handler
func (p *PeerToPeer) HandlePingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	p.UDPSocket.SendMessage(msg, srcAddr)
}

// HandleXpeerPingMessage receives a cross-peer ping message
func (p *PeerToPeer) HandleXpeerPingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	pt := PingType(msg.Header.NetProto)
	if pt == PingReq {
		Log(Debug, "Ping request received")
		// Send a PING response
		r := CreateXpeerPingMessage(PingResp, p.Interface.Mac.String())
		addr, err := net.ParseMAC(string(msg.Data))
		if err != nil {
			Log(Error, "Failed to parse MAC address in crosspeer ping message")
		} else {
			p.SendTo(addr, r)
			Log(Debug, "Sending to %s", addr.String())
		}
	} else {
		Log(Debug, "Ping response received")
		// Handle PING response
		for i, peer := range p.NetworkPeers {
			if peer.PeerHW.String() == string(msg.Data) {
				peer.PingCount = 0
				peer.LastContact = time.Now()
				p.PeersLock.Lock()
				p.NetworkPeers[i] = peer
				p.PeersLock.Unlock()
				runtime.Gosched()
			}
		}
	}
}

// HandleIntroMessage receives an introduction string from another peer during handshake
func (p *PeerToPeer) HandleIntroMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Info, "Introduction string from %s[%d]", srcAddr, msg.Header.ProxyID)
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	p.PeersLock.Lock()
	peer, exists := p.NetworkPeers[id]
	p.PeersLock.Unlock()
	runtime.Gosched()
	if !exists {
		Log(Debug, "Received introduction confirmation from unknown peer: %s", id)
		p.Dht.SendUpdateRequest()
		return
	}
	if msg.Header.ProxyID > 0 && peer.ProxyID == 0 {
		peer.ForceProxy = true
		peer.PeerAddr = nil
		peer.Endpoint = nil
		peer.State = PeerStateInit
		peer.KnownIPs = peer.KnownIPs[:0]
		p.PeersLock.Lock()
		p.NetworkPeers[id] = peer
		p.PeersLock.Unlock()
		runtime.Gosched()
		return
	}
	peer.PeerHW = mac
	peer.PeerLocalIP = ip
	peer.State = PeerStateConnected
	peer.LastContact = time.Now()
	p.PeersLock.Lock()
	p.IPIDTable[ip.String()] = id
	p.MACIDTable[mac.String()] = id
	p.NetworkPeers[id] = peer
	p.PeersLock.Unlock()
	runtime.Gosched()
	Log(Info, "Connection with peer %s has been established", id)
}

// HandleIntroRequestMessage is a handshake request from another peer
func (p *PeerToPeer) HandleIntroRequestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	id := string(msg.Data)
	p.PeersLock.Lock()
	peer, exists := p.NetworkPeers[id]
	p.PeersLock.Unlock()
	runtime.Gosched()
	if !exists {
		Log(Debug, "Introduction request came from unknown peer: %s", id)
		p.Dht.SendUpdateRequest()
		return
	}
	response := p.PrepareIntroductionMessage(p.Dht.ID)
	response.Header.ProxyID = uint16(peer.ProxyID)
	_, err := p.UDPSocket.SendMessage(response, srcAddr)
	if err != nil {
		Log(Error, "Failed to respond to introduction request: %v", err)
	}
}

// HandleProxyMessage receives a control packet from proxy
func (p *PeerToPeer) HandleProxyMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	// Proxy registration data
	if msg.Header.ProxyID < 1 {
		return
	}
	ip := string(msg.Data)
	Log(Info, "Proxy confirmation received from %s. Tunnel ID %d", ip, int(msg.Header.ProxyID))
	for key, peer := range p.NetworkPeers {
		if peer.PeerAddr.String() == ip {
			peer.ProxyID = int(msg.Header.ProxyID)
			p.PeersLock.Lock()
			p.NetworkPeers[key] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
			return
		}
	}
	Log(Warning, "Can't set Tunnel#%d for %s: Can't find address", int(msg.Header.ProxyID), ip)
}

// HandleBadTun notified peer about proxy being malfunction
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) {
	for key, peer := range p.NetworkPeers {
		if peer.ProxyID == int(msg.Header.ProxyID) && peer.Endpoint.String() == srcAddr.String() {
			Log(Debug, "Cleaning bad tunnel %d from %s", msg.Header.ProxyID, srcAddr.String())
			peer.ProxyID = 0
			peer.Endpoint = nil
			peer.Forwarder = nil
			peer.PeerAddr = nil
			peer.State = PeerStateInit
			p.PeersLock.Lock()
			p.NetworkPeers[key] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
		}
	}
}

// HandleTestMessage responses with a test message when another peer trying to
// establish direct connection
func (p *PeerToPeer) HandleTestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	response := CreateTestP2PMessage(p.Crypter, "TEST", 0)
	_, err := p.UDPSocket.SendMessage(response, srcAddr)
	if err != nil {
		Log(Error, "Failed to respond to test message: %v", err)
	}

}

// SendTo sends a p2p packet by MAC address
func (p *PeerToPeer) SendTo(dst net.HardwareAddr, msg *P2PMessage) (int, error) {
	// TODO: Speed up this by switching to map
	Log(Trace, "Requested Send to %s", dst.String())
	id, exists := p.MACIDTable[dst.String()]
	if exists {
		p.PeersLock.Lock()
		peer, exists := p.NetworkPeers[id]
		p.PeersLock.Unlock()
		runtime.Gosched()
		if exists {
			msg.Header.ProxyID = uint16(peer.ProxyID)
			Log(Debug, "Sending to %s via proxy id %d", dst.String(), msg.Header.ProxyID)
			size, err := p.UDPSocket.SendMessage(msg, peer.Endpoint)
			return size, err
		}
	}
	return 0, nil
}

// StopInstance stops current instance
func (p *PeerToPeer) StopInstance() {
	p.PeersLock.Lock()
	for i, peer := range p.NetworkPeers {
		peer.State = PeerStateDisconnect
		p.NetworkPeers[i] = peer
	}
	p.PeersLock.Unlock()
	runtime.Gosched()
	stopStarted := time.Now()
	for len(p.NetworkPeers) > 0 {
		if time.Since(stopStarted) > time.Duration(time.Second*5) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	Log(Info, "All peers under this instance has been removed")

	var ip net.IP
	if p.Dht == nil || p.Dht.Network == nil {
		Log(Warning, "DHT isn't in use")
	} else {
		ip = p.Dht.Network.IP
	}
	p.Dht.Stop()
	p.UDPSocket.Stop()
	p.Shutdown = true
	Log(Info, "Stopping P2P Message handler")
	// Tricky part: we need to send a message to ourselves to quit blocking operation
	msg := CreateTestP2PMessage(p.Crypter, "STOP", 1)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", p.Dht.P2PPort))
	p.UDPSocket.SendMessage(msg, addr)
	var ipIt = 200
	if ip != nil {
		for p.IsDeviceExists(p.Interface.Name) {
			time.Sleep(1 * time.Second)
			target := fmt.Sprintf("%d.%d.%d.%d:9922", ip[0], ip[1], ip[2], ipIt)
			Log(Info, "Dialing %s", target)
			_, err := net.DialTimeout("tcp", target, 2*time.Second)
			if err != nil {
				Log(Info, "ERROR: %v", err)
			}
			ipIt++
			if ipIt == 255 {
				break
			}
		}
	}
	time.Sleep(3 * time.Second)
	p.ReadyToStop = true
}

// ReadDHTPeers - reads a list of peers received by DHT client
func (p *PeerToPeer) ReadDHTPeers() {
	for {
		if p.Shutdown {
			break
		}
		select {
		case peers, hasData := <-p.Dht.PeerChannel:
			if hasData {
				p.UpdatePeers(peers)
			} else {
				Log(Trace, "Clossed channel")
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	Log(Info, "Stopped DHT reader channel")
}

// ReadProxies - reads a list of proxies received by DHT client
func (p *PeerToPeer) ReadProxies() {
	for {
		if p.Shutdown {
			break
		}
		if p.Dht == nil {
			time.Sleep(10 * time.Millisecond)
			continue
		}
		select {
		case proxy, hasData := <-p.Dht.ProxyChannel:
			if hasData {
				exists := false
				for i, peer := range p.NetworkPeers {
					if i == proxy.DestinationID {
						peer.State = PeerStateHandshakingForwarder
						peer.Forwarder = proxy.Addr
						peer.Endpoint = proxy.Addr
						p.PeersLock.Lock()
						p.NetworkPeers[i] = peer
						p.PeersLock.Unlock()
						runtime.Gosched()
						exists = true
					}
				}
				if !exists {
					Log(Info, "Received forwarder for unknown peer")
					p.Dht.SendUpdateRequest()
				}

			} else {
				Log(Trace, "Closed channel")
			}
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
	Log(Info, "Stopped Proxy reader channel")
}

// UpdatePeers updates information about known peers
func (p *PeerToPeer) UpdatePeers(peers []PeerIP) {
	for _, newPeer := range peers {
		if newPeer.ID == "" {
			continue
		}
		found := false
		for _, peer := range p.NetworkPeers {
			if peer.ID == newPeer.ID {
				found = true
			}
		}
		if !found && newPeer.ID != p.Dht.ID {
			peer := new(NetworkPeer)
			peer.ID = newPeer.ID
			peer.KnownIPs = newPeer.Ips
			peer.State = PeerStateInit
			p.PeersLock.Lock()
			p.NetworkPeers[newPeer.ID] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
			go p.NetworkPeers[newPeer.ID].Run(p)
		}
	}
}
