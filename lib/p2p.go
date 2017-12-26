package ptp

import (
	"crypto/rand"
	"fmt"
	"io/ioutil"
	"net"
	"strings"
	"sync"
	"time"

	upnp "github.com/NebulousLabs/go-upnp"
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
	UDPSocket       *Network                             // Peer-to-peer interconnection socket
	LocalIPs        []net.IP                             // List of IPs available in the system
	Dht             *DHTClient                           // DHT Client
	Crypter         Crypto                               // Cryptography subsystem
	Shutdown        bool                                 // Set to true when instance in shutdown mode
	ForwardMode     bool                                 // Skip local peer discovery
	ReadyToStop     bool                                 // Set to true when instance is ready to stop
	MessageHandlers map[uint16]MessageHandler            // Callbacks for network packets
	PacketHandlers  map[PacketType]PacketHandlerCallback // Callbacks for packets received by TAP interface
	PeersLock       sync.Mutex                           // Lock for peers map
	Hash            string                               // Infohash for this instance
	Routers         string                               // Comma-separated list of Bootstrap nodes
	Interface       NetworkInterface                     // TAP Interface
	Peers           *PeerList                            // Known peers
	HolePunching    sync.Mutex                           // Mutex for hole punching sync
	Proxies         []*proxyServer                       // List of proxies
	outboundIP      net.IP                               // Outbound IP
}

// AssignInterface - Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PeerToPeer) AssignInterface(interfaceName string) error {
	var err error
	p.Interface.Name = interfaceName

	if p.Interface.IP == nil {
		return fmt.Errorf("No IP provided")
	}
	if p.Interface.Mac == nil {
		return fmt.Errorf("No Hardware address provided")
	}

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
func New(argIP, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool, port int, ignoreIPs []string, outboundIP net.IP) *PeerToPeer {
	//argDht = "mdht.subut.ai:6881"
	p := new(PeerToPeer)
	p.outboundIP = outboundIP
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

	p.Hash = argHash
	p.Routers = argDht

	p.setupHandlers()

	p.UDPSocket = new(Network)
	p.UDPSocket.Init("", port)
	go p.UDPSocket.Listen(p.HandleP2PMessage)
	go p.UDPSocket.KeepAlive(p.retrieveFirstDHTRouter())
	p.waitForRemotePort()

	// Create new DHT Client, configure it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code

	Log(Info, "Started UDP Listener at port %d", p.UDPSocket.GetPort())
	// go func() {
	// 	err = p.attemptPortForward(uint16(p.UDPSocket.GetPort()), interfaceName)
	// 	if err != nil {
	// 		Log(Error, "UPnP Failed: %s", err)
	// 	}
	// }()

	err = p.StartDHT(p.Hash, p.Routers)
	if err != nil {
		Log(Info, "Retrying DHT connection")
		attempts := 0
		for attempts < 3 {
			err = p.StartDHT(p.Hash, p.Routers)
			if err != nil {
				Log(Info, "Another attempt failed")
				attempts++
			} else {
				Log(Info, "Connection established")
				err = nil
				break
			}
		}
	}
	if err != nil {
		Log(Error, "Failed to start instance due to problems with bootstrap node connection")
		return nil
	}

	err = p.prepareInterfaces(argIP, interfaceName)
	if err != nil {
		return nil
	}

	go p.ListenInterface()
	return p
}

// This method will block for seconds or unless we receive remote port
// from echo server
func (p *PeerToPeer) waitForRemotePort() {
	started := time.Now()
	for p.UDPSocket.remotePort == 0 {
		time.Sleep(time.Millisecond * 100)
		if time.Since(started) > time.Duration(time.Second*3) {
			break
		}
	}
	if p.UDPSocket.remotePort == 0 {
		Log(Warning, "Didn't received remote port")
		p.UDPSocket.remotePort = p.UDPSocket.GetPort()
		return
	}
	Log(Warning, "Remote port received: %d", p.UDPSocket.remotePort)
}

func (p *PeerToPeer) retrieveFirstDHTRouter() *net.UDPAddr {
	Log(Info, "Routers: %s", p.Routers)
	routers := strings.Split(p.Routers, ",")
	if len(routers) == 0 {
		return nil
	}
	router := strings.Split(routers[0], ":")
	if len(router) != 2 {
		return nil
	}
	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", router[0], 6882))
	if err != nil {
		return nil
	}
	return addr
}

func (p *PeerToPeer) prepareInterfaces(ip, interfaceName string) error {
	if ip == "dhcp" {
		ipn, maskn, err := p.RequestIP(p.Interface.Mac.String(), interfaceName)
		if err != nil {
			Log(Error, "%v", err)
			return err
		}
		p.Interface.IP = ipn
		p.Interface.Mask = maskn
	} else {
		p.Interface.IP = net.ParseIP(ip)
		ipn, maskn, err := p.ReportIP(ip, p.Interface.Mac.String(), interfaceName)
		if err != nil {
			Log(Error, "%v", err)
			return err
		}
		p.Interface.IP = ipn
		p.Interface.Mask = maskn
	}
	return nil
}

func (p *PeerToPeer) attemptPortForward(port uint16, name string) error {
	Log(Info, "Trying to forward port %d", port)
	d, err := upnp.Discover()
	if err != nil {
		return err
	}
	err = d.Forward(port, "subutai-"+name)
	if err != nil {
		return err
	}
	Log(Info, "Port %d has been forwarded", port)
	return nil
}

// Init will initialize PeerToPeer
func (p *PeerToPeer) Init() {
	p.Peers = new(PeerList)
	p.Peers.Init()
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
	Log(Info, "Requesting IP from Bootstrap node")
	requestedAt := time.Now()
	interval := time.Duration(3 * time.Second)
	p.Dht.sendDHCP(nil, nil)
	for p.Dht.IP == nil && p.Dht.Network == nil {
		if time.Since(requestedAt) > interval {
			p.StopInstance()
			return nil, nil, fmt.Errorf("No IP were received. Swarm is empty")
		}
		time.Sleep(100 * time.Millisecond)
	}
	p.Interface.IP = p.Dht.IP
	p.Interface.Mask = p.Dht.Network.Mask
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
		Log(Warning, "IP was not in CIDR format. Assumming /24")
		ip, ipnet, err = net.ParseCIDR(ipAddress)
		if err != nil {
			return nil, nil, fmt.Errorf("Failed to setup provided IP address for local device")
		}
	}
	if ipnet == nil {
		return nil, nil, fmt.Errorf("Can't report network information. Reason: Unknown")
	}
	p.Dht.IP = ip
	p.Dht.Network = ipnet

	p.Dht.sendDHCP(ip, ipnet)
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
	p.Dht.OutboundIP = p.outboundIP
	p.Dht.LocalPort = p.UDPSocket.GetPort()
	if p.UDPSocket.remotePort == 0 {
		p.Dht.LocalPort = p.Dht.RemotePort
	} else {
		p.Dht.RemotePort = p.UDPSocket.remotePort
	}
	err := p.Dht.TCPInit(hash, routers)
	if err != nil {
		return fmt.Errorf("Failed to initialize DHT: %s", err)
	}
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
	err = p.Dht.WaitID()
	if err != nil {
		Log(Error, "Failed to retrieve ID from bootstrap node: %s", err)
	}
	return nil
}

func (p *PeerToPeer) markPeerForRemoval(id, reason string) error {
	peer := p.Peers.GetPeer(id)
	if peer == nil {
		return fmt.Errorf("Peer was not found")
	}
	Log(Info, "Removing peer %s: Reason %s", id, reason)
	peer.SetState(PeerStateDisconnect, p)
	p.Peers.Update(id, peer)
	return nil
}

// Run is a main loop
func (p *PeerToPeer) Run() {
	// p.Dht.LastUpdate = time.Unix(1, 1)
	p.Dht.LastUpdate = time.Now()
	for {
		if p.Shutdown {
			// TODO: Do it more safely
			if p.ReadyToStop {
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}

		select {
		case peer, pd := <-p.Dht.PeerData:
			if pd {
				// Received peer update
				p.handlePeerData(peer)
			}
		case state, s := <-p.Dht.StateChannel:
			if s {
				peer := p.Peers.GetPeer(state.ID)
				if peer != nil {
					peer.RemoteState = state.State
					p.Peers.Update(state.ID, peer)
				} else {
					Log(Warning, "Received state of unknown peer. Updating peers")
					p.Dht.sendFind()
				}
			}
		case proxy, pr := <-p.Dht.ProxyChannel:
			if pr {
				go p.initProxy(proxy)
			}
		default:
			p.removeStoppedPeers()
			p.checkBootstrapNodes()
			p.checkLastDHTUpdate()
			p.checkProxies()
			time.Sleep(100 * time.Millisecond)
		}
	}
	Log(Info, "Shutting down instance %s completed", p.Dht.NetworkHash)
}

func (p *PeerToPeer) checkLastDHTUpdate() {
	passed := time.Since(p.Dht.LastUpdate)
	if passed > time.Duration(30*time.Second) {
		Log(Debug, "DHT Last Update timeout passed")
		p.Dht.sendProxy()
		err := p.Dht.sendFind()
		if err != nil {
			Log(Error, "Failed to send update: %s", err)
		}
	}
}

func (p *PeerToPeer) checkBootstrapNodes() {
	if !p.Dht.Connected {
		p.StartDHT(p.Hash, p.Routers)
		p.Dht.sendDHCP(p.Dht.IP, p.Dht.Network)
	}
}

func (p *PeerToPeer) removeStoppedPeers() {
	peers := p.Peers.Get()
	for id, peer := range peers {
		if peer.State == PeerStateStop {
			Log(Info, "Removing peer %s", id)
			p.Peers.Delete(id)
			Log(Info, "Peer %s has been removed", id)
			break
		}
	}
}

func (p *PeerToPeer) checkProxies() {
	lifetime := time.Duration(time.Second * 30)
	for i, proxy := range p.Proxies {
		if time.Since(proxy.LastUpdate) > lifetime {
			Log(Info, "Proxy connection with %s [EP: %s] has been died", proxy.Addr.String(), proxy.Endpoint.String())
			p.Proxies = append(p.Proxies[:i], p.Proxies[:i+1]...)
			break
		}
	}
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
// TODO: Remove unused function
func (p *PeerToPeer) PurgePeers() {
	/*peers := p.Peers.Get()
	for i, peer := range peers {
		found := false
		for _, newPeer := range p.Dht.Peers {
			if newPeer.ID == peer.ID {
				found = true
			}
		}
		if !found {
			Log(Info, "Removing outdated peer")
			p.Peers.Delete(i)
		}
	}*/
	return
}

// SyncForwarders extracts proxies from DHT and assign them to target peers
func (p *PeerToPeer) SyncForwarders() int {
	count := 0
	for _, fwd := range p.Dht.Forwarders {
		peers := p.Peers.Get()
		for i, peer := range peers {
			if peer.Endpoint == nil && fwd.DestinationID == peer.ID && peer.Forwarder == nil {
				Log(Info, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr
				peer.Forwarder = fwd.Addr
				peer.SetState(PeerStateHandshakingForwarder, p)
				p.Peers.Update(i, peer)
				count++
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
	if msg == nil {
		Log(Error, "Received broken message")
		return
	}
	//var msgType MSG_TYPE = MSG_TYPE(msg.Header.Type)
	// Decrypt message if crypter is active
	if p.Crypter.Active && (msg.Header.Type == MsgTypeIntro || msg.Header.Type == MsgTypeNenc || msg.Header.Type == MsgTypeIntroReq || msg.Header.Type == MsgTypeTest) {
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
	addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		p.UDPSocket.SendMessage(msg, srcAddr)
		for i, proxy := range p.Proxies {
			if proxy.Addr.String() == srcAddr.String() {
				p.Proxies[i].LastUpdate = time.Now()
			}
		}
		return
	}
	port := addr.Port
	if p.UDPSocket.remotePort == 0 {
		p.UDPSocket.remotePort = port
	} else {
		if port != p.UDPSocket.GetPort() && port != p.UDPSocket.remotePort {
			Log(Debug, "Port translation detected %d -> %d", p.UDPSocket.GetPort(), port)
			p.UDPSocket.remotePort = port
		}
	}
}

// HandleXpeerPingMessage receives a cross-peer ping message
func (p *PeerToPeer) HandleXpeerPingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	pt := PingType(msg.Header.NetProto)
	if pt == PingReq {
		Log(Debug, "Ping request received: %s. Responding with %s", string(msg.Data), p.Interface.Mac.String())
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
		peers := p.Peers.Get()
		for i, peer := range peers {
			if peer.PeerHW != nil && peer.PeerHW.String() == string(msg.Data) {
				peer.PingCount = 0
				peer.LastContact = time.Now()
				p.Peers.Update(i, peer)
				break
			}
		}
	}
}

// HandleIntroMessage receives an introduction string from another peer during handshake
func (p *PeerToPeer) HandleIntroMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Info, "Introduction string from %s[%d]", srcAddr, msg.Header.ProxyID)
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	if len(id) != 36 {
		Log(Debug, "Received wrong ID in introduction message: %s", id)
		return
	}
	peer := p.Peers.GetPeer(id)
	// Do nothing when handshaking already done
	if peer.State != PeerStateHandshaking && peer.State != PeerStateHandshakingForwarder {
		return
	}
	if peer == nil {
		Log(Debug, "Received introduction confirmation from unknown peer: %s", id)
		p.Dht.sendFind()
		return
	}

	if mac == nil {
		Log(Error, "Received empty MAC address. Skipping")
		return
	}
	if ip == nil {
		Log(Error, "No IP received. Skipping")
		return
	}
	peer.PeerHW = mac
	peer.PeerLocalIP = ip
	peer.LastContact = time.Now()
	peer.SetState(PeerStateConnected, p)
	p.Peers.Update(id, peer)
	Log(Info, "Connection with peer %s has been established", id)
}

// HandleIntroRequestMessage is a handshake request from another peer
func (p *PeerToPeer) HandleIntroRequestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	id := string(msg.Data)
	if len(id) != 36 {
		Log(Debug, "Introduction request with malformed ID [%s] from %s", id, srcAddr.String())
		return
	}
	peer := p.Peers.GetPeer(id)
	if peer == nil {
		Log(Debug, "Introduction request came from unknown peer: %s [%s]", id, srcAddr.String())
		p.Dht.sendFind()
		return
	}
	proxy := false
	if msg.Header.ProxyID > 0 {
		proxy = true
		Log(Info, "Received introduction request via proxy")
		if len(peer.Proxies) == 0 {
			Log(Warning, "Peer %s has no proxies attached", id)
			p.Dht.sendRequestProxy(id)
			return
		}
	} else {
		Log(Info, "Received introduction request directly")
	}

	response := p.PrepareIntroductionMessage(p.Dht.ID)
	if proxy {
		response.Header.ProxyID = 1
		for _, peerProxy := range peer.Proxies {
			Log(Info, "Sending handshake response over proxy %s", peerProxy.String())
			_, err := p.UDPSocket.SendMessage(response, peerProxy)
			if err != nil {
				Log(Error, "Failed to respond to introduction request over proxy: %v", err)
			}
		}
		return
	}
	Log(Info, "Sending handshake response")
	_, err := p.UDPSocket.SendMessage(response, srcAddr)
	if err != nil {
		Log(Error, "Failed to respond to introduction request: %v", err)
	}
}

// HandleProxyMessage receives a control packet from proxy
// Proxy packets comes in format of UDP connection address
func (p *PeerToPeer) HandleProxyMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Debug, "New proxy message from %s", srcAddr)
	for i, proxy := range p.Proxies {
		if proxy.Addr.String() == srcAddr.String() && proxy.Status == proxyConnecting {
			p.Proxies[i].Status = proxyActive
			addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
			if err != nil {
				Log(Error, "Failed to resolve proxy addr: %s", err)
				return
			}
			Log(Info, "This peer is now available over %s", addr.String())
			p.Dht.sendReportProxy(addr)
		}
	}
}

// HandleBadTun notified peer about proxy being malfunction
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) {
	peers := p.Peers.Get()
	for i, peer := range peers {
		if peer.ProxyID == msg.Header.ProxyID && peer.Endpoint.String() == srcAddr.String() {
			Log(Debug, "Cleaning bad tunnel %d from %s", msg.Header.ProxyID, srcAddr.String())
			peer.ProxyID = 0
			peer.Endpoint = nil
			peer.Forwarder = nil
			peer.PeerAddr = nil
			peer.SetState(PeerStateInit, p)
			p.Peers.Update(i, peer)
		}
	}
}

// HandleTestMessage responses with a test message when another peer trying to
// establish direct connection
func (p *PeerToPeer) HandleTestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	if len(p.Dht.ID) != 36 {
		return
	}

	if len(msg.Data) < 36 || len(msg.Data) > 40 {
		Log(Error, "Malformed data received during test: %s", string(msg.Data))
		return
	}

	// See if we have peer with this ID
	id := string(msg.Data[0:36])
	if len(id) != 36 {
		Log(Error, "Wrong ID during test message")
		return
	}

	peer := p.Peers.GetPeer(id)
	if peer != nil {
		if peer.State == PeerStateConnectingDirectly || peer.State == PeerStateConnectingInternet {
			peer.TestPacketReceived = true
			p.Peers.Update(id, peer)
			response := CreateTestP2PMessage(p.Crypter, p.Dht.ID, 0)
			_, err := p.UDPSocket.SendMessage(response, srcAddr)
			if err != nil {
				Log(Error, "Failed to respond to test message: %v", err)
			}
		} else if peer.State == PeerStateConnected && peer.IsUsingTURN {
			Log(Info, "Received test message from peer which was previously connected over TURN")
			if len(peer.KnownIPs) == 0 {
				return
			}
			peer.Endpoint = peer.KnownIPs[0]
			peer.IsUsingTURN = false
			p.Peers.Update(peer.ID, peer)
			Log(Info, "Peer %s switched to direct UDP connection %s", peer.ID, peer.Endpoint.String())
		} else {
			Log(Info, "Received test message for peer in unsupported state")
		}
	}
}

// SendTo sends a p2p packet by MAC address
func (p *PeerToPeer) SendTo(dst net.HardwareAddr, msg *P2PMessage) (int, error) {
	// TODO: Speed up this by switching to map
	Log(Trace, "Requested Send to %s", dst.String())
	//id, exists := p.MACIDTable[dst.String()]
	endpoint, proxy, err := p.Peers.GetEndpointAndProxy(dst.String())
	if err == nil && endpoint != nil {
		Log(Debug, "Sending to %s via proxy id %d", dst.String(), proxy)
		msg.Header.ProxyID = uint16(proxy)
		size, err := p.UDPSocket.SendMessage(msg, endpoint)
		return size, err
	}
	Log(Trace, "Skipping packet")
	return 0, nil
}

// StopInstance stops current instance
func (p *PeerToPeer) StopInstance() {
	peers := p.Peers.Get()
	for i, peer := range peers {
		peer.SetState(PeerStateDisconnect, p)
		p.Peers.Update(i, peer)
	}
	stopStarted := time.Now()
	for p.Peers.Length() > 0 {
		if time.Since(stopStarted) > time.Duration(time.Second*5) {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	Log(Info, "All peers under this instance has been removed")

	// var ip net.IP
	// if p.Dht == nil || p.Dht.Network == nil {
	// 	Log(Warning, "DHT isn't in use")
	// } else {
	// 	ip = p.Dht.Network.IP
	// }
	p.Dht.Shutdown()
	p.UDPSocket.Stop()
	if p.Interface.Interface != nil && p.Interface.Interface.file != nil {
		closeInterface(p.Interface.Interface.file)
	}
	/*if runtime.GOOS != "windows" {
		p.Interface.Interface.file.Close()
	}*/
	p.Shutdown = true
	// Log(Info, "Stopping P2P Message handler")
	// // Tricky part: we need to send a message to ourselves to quit blocking operation
	// msg := CreateTestP2PMessage(p.Crypter, "STOP", 1)
	// addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("127.0.0.1:%d", p.Dht.P2PPort))
	// p.UDPSocket.SendMessage(msg, addr)
	// var ipIt = 200
	// if ip != nil {
	// 	for p.IsDeviceExists(p.Interface.Name) {
	// 		time.Sleep(1 * time.Second)
	// 		target := fmt.Sprintf("%d.%d.%d.%d:9922", ip[0], ip[1], ip[2], ipIt)
	// 		Log(Info, "Dialing %s", target)
	// 		_, err := net.DialTimeout("tcp", target, 2*time.Second)
	// 		if err != nil {
	// 			Log(Info, "ERROR: %v", err)
	// 		}
	// 		ipIt++
	// 		if ipIt == 255 {
	// 			break
	// 		}
	// 	}
	// }
	time.Sleep(3 * time.Second)
	p.ReadyToStop = true
}

func (p *PeerToPeer) handlePeerData(peerData NetworkPeer) {
	if peerData.ID == "" {
		return
	}
	// Check if such peer exists
	peer := p.Peers.GetPeer(peerData.ID)
	if peer == nil {
		peer := new(NetworkPeer)
		Log(Info, "Received new peer %s", peerData.ID)
		peer.ID = peerData.ID
		peer.SetState(PeerStateInit, p)
		p.Peers.Update(peerData.ID, peer)
		p.Peers.RunPeer(peerData.ID, p)
		return
	}
	// When peer data contains IPs this means we received
	// list of IP addresses of this peer
	if peerData.KnownIPs != nil && len(peerData.KnownIPs) > 0 {
		Log(Info, "Received peer IPs %s", peerData.ID)
		for _, newip := range peerData.KnownIPs {
			found := false
			for _, knownip := range peer.KnownIPs {
				if knownip == newip {
					found = true
				}
			}
			if !found {
				peer.KnownIPs = append(peer.KnownIPs, newip)
			}
		}
		p.Peers.Update(peer.ID, peer)
		return
	}

	if peer != nil && len(peerData.Proxies) > 0 {
		Log(Info, "Received proxies for peer %s", peerData.ID)
		peers := p.Peers.Get()
		for _, proxy := range peerData.Proxies {
			for _, existingPeer := range peers {
				if existingPeer.Endpoint.String() == proxy.String() && existingPeer.ID != peerData.ID {
					existingPeer.SetState(PeerStateDisconnect, p)
					Log(Info, "Peer %s was associated with address %s. Disconnecting", existingPeer.ID, proxy.String())
				}
			}
		}
		peer.Proxies = peerData.Proxies
		p.Peers.Update(peer.ID, peer)
	}
}
