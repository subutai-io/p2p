package ptp

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	upnp "github.com/NebulousLabs/go-upnp"
)

// GlobalMTU value specified on daemon start
var GlobalMTU = DefaultMTU

var UsePMTU = false

// PeerToPeer - Main structure
type PeerToPeer struct {
	Config          Configuration                        // Network interface configuration tool
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
	//Routers         map[int]string                       // Comma-separated list of Bootstrap nodes
	Interface    TAP           // TAP Interface
	Peers        *PeerList     // Known peers
	HolePunching sync.Mutex    // Mutex for hole punching sync
	ProxyManager *ProxyManager // Proxy manager
	outboundIP   net.IP        // Outbound IP
	UsePMTU      bool          // Whether PMTU capabilities are enabled or not
}

// PeerHandshake holds handshake information received from peer
type PeerHandshake struct {
	ID           string
	IP           net.IP
	HardwareAddr net.HardwareAddr
	Endpoint     *net.UDPAddr
}

// ActiveInterfaces is a global (daemon-wise) list of reserved IP addresses
var ActiveInterfaces []net.IP

// AssignInterface - Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PeerToPeer) AssignInterface(interfaceName string) error {
	var err error
	if p.Interface == nil {
		return fmt.Errorf("Failed to initialize TAP")
	}
	err = p.Interface.Init(interfaceName)
	if p.Interface.IsConfigured() {
		return nil
	}
	if err != nil {
		return fmt.Errorf("Failed to initialize TAP: %s", err)
	}

	if p.Interface.GetIP() == nil {
		return fmt.Errorf("No IP provided")
	}
	if p.Interface.GetHardwareAddress() == nil {
		return fmt.Errorf("No Hardware address provided")
	}
	if p.Interface.GetName() == "" {
		return fmt.Errorf("Wrong interface name provided: %s", p.Interface.GetName())
	}

	// Extract necessary information from config file
	err = p.Config.Read()
	if err != nil {
		Log(Error, "Failed to extract information from config file: %v", err)
		return err
	}

	err = p.Interface.Open()
	if err != nil {
		Log(Error, "Failed to open TAP device %s: %v", p.Interface.GetName(), err)
		return err
	}
	Log(Debug, "%v TAP Device created", p.Interface.GetName())

	err = p.Interface.Configure()
	if err != nil {
		return err
	}
	ActiveInterfaces = append(ActiveInterfaces, p.Interface.GetIP())
	Log(Debug, "Interface has been configured")
	p.Interface.MarkConfigured()
	return err
}

// ListenInterface - Listens TAP interface for incoming packets
// Read packets received by TAP interface and send them to a handlePacket goroutine
// This goroutine will execute a callback method based on packet type
func (p *PeerToPeer) ListenInterface() {
	if p.Interface == nil {
		Log(Error, "Failed to start TAP listener: nil object")
		return
	}
	p.Interface.Run()
	for {
		if p.Shutdown {
			break
		}
		packet, err := p.Interface.ReadPacket()
		if err != nil && err != errPacketTooBig {
			Log(Error, "Reading packet: %s", err)
			p.Close()
			break
		}
		if packet != nil {
			go p.handlePacket(packet.Packet, packet.Protocol)
		}
	}
	Log(Debug, "Shutting down interface listener")

	if p.Interface != nil {
		p.Interface.Close()
	}
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
	tap, _ := newTAP("", "127.0.0.1", "00:00:00:00:00:00", "", 0, p.UsePMTU)
	var devName = tap.GetBasename() + fmt.Sprintf("%d", i)
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

// New is an entry point of a P2P library.
// This function will return new PeerToPeer object which later
// should be configured and started using Run() method
func New(mac, hash, keyfile, key, ttl, target string, fwd bool, port int, outboundIP net.IP) *PeerToPeer {
	Log(Debug, "Starting new P2P Instance: %s", hash)
	Log(Debug, "Mac: %s", mac)
	p := new(PeerToPeer)
	p.outboundIP = outboundIP
	p.Init()
	var err error
	p.Interface, err = newTAP(GetConfigurationTool(), "127.0.0.1", "00:00:00:00:00:00", "", DefaultMTU, UsePMTU)
	if err != nil {
		Log(Error, "Failed to create TAP object: %s", err)
		return nil
	}
	p.Interface.SetHardwareAddress(p.validateMac(mac))
	p.FindNetworkAddresses()

	if fwd {
		p.ForwardMode = true
	}

	if keyfile != "" {
		p.Crypter.ReadKeysFromFile(keyfile)
	}
	if key != "" {
		// Override key from file
		if ttl == "" {
			ttl = "default"
		}
		var newKey CryptoKey
		newKey = p.Crypter.EnrichKeyValues(newKey, key, ttl)
		p.Crypter.Keys = append(p.Crypter.Keys, newKey)
		p.Crypter.ActiveKey = p.Crypter.Keys[0]
		p.Crypter.Active = true
	}

	if p.Crypter.Active {
		Log(Debug, "Traffic encryption is enabled. Key valid until %s", p.Crypter.ActiveKey.Until.String())
	} else {
		Log(Debug, "No AES key were provided. Traffic encryption is disabled")
	}

	p.Hash = hash

	p.setupHandlers()

	p.UDPSocket = new(Network)
	p.UDPSocket.Init("", port)
	go p.UDPSocket.Listen(p.HandleP2PMessage)
	go p.UDPSocket.KeepAlive(target)
	p.waitForRemotePort()

	// Create new DHT Client, configure it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code

	Log(Debug, "Started UDP Listener at port %d", p.UDPSocket.GetPort())

	p.Dht = new(DHTClient)
	err = p.Dht.Init(p.Hash)
	if err != nil {
		Log(Error, "Failed to initialize DHT: %s", err)
		return nil
	}

	p.setupTCPCallbacks()
	p.ProxyManager = new(ProxyManager)
	p.ProxyManager.init()
	return p
}

// ReadDHT will read packets from bootstrap node
func (p *PeerToPeer) ReadDHT() {
	for !p.Shutdown {
		packet, err := p.Dht.read()
		if err != nil {
			break
		}
		go func() {
			cb, e := p.Dht.TCPCallbacks[packet.Type]
			if !e {
				Log(Error, "Unsupported packet from DHT")
				return
			}
			err = cb(packet)
			if err != nil {
				Log(Error, "DHT: %s", err)
			}
		}()
	}
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
	if p.UDPSocket != nil && p.UDPSocket.remotePort == 0 {
		Log(Warning, "Didn't received remote port")
		p.UDPSocket.remotePort = p.UDPSocket.GetPort()
		return
	}
	Log(Warning, "Remote port received: %d", p.UDPSocket.remotePort)
}

// PrepareInterfaces will assign IPs to interfaces
func (p *PeerToPeer) PrepareInterfaces(ip, interfaceName string) error {

	iface, err := p.validateInterfaceName(interfaceName)
	if err != nil {
		Log(Error, "Interface name validation failed: %s", err)
		return nil
	}
	if p.IsDeviceExists(iface) {
		Log(Error, "Interface is already in use. Can't create duplicate")
		return nil
	}

	if ip == "dhcp" {
		ipn, maskn, err := p.RequestIP(p.Interface.GetHardwareAddress().String(), iface)
		if err != nil {
			return err
		}
		p.Interface.SetIP(ipn)
		p.Interface.SetMask(maskn)
	} else {
		p.Interface.SetIP(net.ParseIP(ip))
		ipn, maskn, err := p.ReportIP(ip, p.Interface.GetHardwareAddress().String(), iface)
		if err != nil {
			return err
		}
		p.Interface.SetIP(ipn)
		p.Interface.SetMask(maskn)
	}
	return nil
}

func (p *PeerToPeer) attemptPortForward(port uint16, name string) error {
	Log(Debug, "Trying to forward port %d", port)
	d, err := upnp.Discover()
	if err != nil {
		return err
	}
	err = d.Forward(port, "subutai-"+name)
	if err != nil {
		return err
	}
	Log(Debug, "Port %d has been forwarded", port)
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
		Log(Debug, "Generate MAC for TAP device: %s", mac)
	}
	return hw
}

func (p *PeerToPeer) validateInterfaceName(name string) (string, error) {
	if name == "" {
		name = p.GenerateDeviceName(1)
	} else {
		if len(name) > MaximumInterfaceNameLength {
			Log(Debug, "Interface name length should be %d symbols max", MaximumInterfaceNameLength)
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
	Log(Debug, "Requesting IP from Bootstrap node")
	requestedAt := time.Now()
	interval := time.Duration(3 * time.Second)
	p.Dht.sendDHCP(nil, nil)
	for p.Dht.IP == nil && p.Dht.Network == nil {
		if time.Since(requestedAt) > interval {
			//p.StopInstance()
			return nil, nil, fmt.Errorf("No IP were received. Swarm is empty")
		}
		time.Sleep(100 * time.Millisecond)
	}
	p.Interface.SetIP(p.Dht.IP)
	p.Interface.SetMask(p.Dht.Network.Mask)
	err := p.AssignInterface(device)
	if err != nil {
		return nil, nil, fmt.Errorf("Failed to configure interface: %s", err)
	}
	return p.Dht.IP, p.Dht.Network.Mask, nil
}

// ReportIP will send IP specified at service start to DHCP-like service
func (p *PeerToPeer) ReportIP(ipAddress, mac, device string) (net.IP, net.IPMask, error) {
	Log(Debug, "Reporting IP to bootstranp node: %s", ipAddress)
	ip, ipnet, err := net.ParseCIDR(ipAddress)
	if err != nil {
		nip := net.ParseIP(ipAddress)
		if nip == nil {
			return nil, nil, fmt.Errorf("Invalid address were provided for network interface. Use -ip \"dhcp\" or specify correct IP address")
		}
		ipAddress += `/24`
		Log(Debug, "IP was not in CIDR format. Assumming /24")
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

// Run is a main loop
func (p *PeerToPeer) Run() {
	// Request proxies from DHT
	p.Dht.sendProxy()

	initialRequestSent := false
	started := time.Now()
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
		p.removeStoppedPeers()
		p.checkLastDHTUpdate()
		p.checkProxies()
		time.Sleep(100 * time.Millisecond)
		if !initialRequestSent && time.Since(started) > time.Duration(time.Millisecond*5000) {
			initialRequestSent = true
			p.Dht.sendFind()
		}
		if p.Interface.IsBroken() {
			Log(Info, "TAP interface is broken. Shutting down instance %s", p.Hash)
			p.Close()
		}
	}
	Log(Info, "Shutting down instance %s completed", p.Dht.NetworkHash)
}

func (p *PeerToPeer) checkLastDHTUpdate() {
	passed := time.Since(p.Dht.LastUpdate)
	if passed > time.Duration(30*time.Second) {
		Log(Debug, "DHT Last Update timeout passed")
		// Request new proxies if we don't have any more
		if len(p.ProxyManager.get()) == 0 {
			p.Dht.sendProxy()
		}
		err := p.Dht.sendFind()
		if err != nil {
			Log(Error, "Failed to send update: %s", err)
		}
	}
}

// TODO: Check if this method is still actual
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
	p.ProxyManager.check()
	// Unlink dead proxies
	proxies := p.ProxyManager.get()
	list := []*net.UDPAddr{}
	for _, proxy := range proxies {
		if proxy.Endpoint != nil && proxy.Status == proxyActive {
			list = append(list, proxy.Endpoint)
		}
	}
	if p.ProxyManager.hasChanges && len(list) > 0 {
		p.ProxyManager.hasChanges = false
		p.Dht.sendReportProxy(list)
	}
}

// PrepareIntroductionMessage collects client ID, mac and IP address
// and create a comma-separated line
// endpoint is an address that received this introduction message
func (p *PeerToPeer) PrepareIntroductionMessage(id, endpoint string) *P2PMessage {
	var intro = id + "," + p.Interface.GetHardwareAddress().String() + "," + p.Interface.GetIP().String() + "," + endpoint
	msg, err := p.CreateMessage(MsgTypeIntro, []byte(intro), 0, true)
	if err != nil {
		return nil
	}
	return msg
}

// WriteToDevice writes data to created TAP interface
func (p *PeerToPeer) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var packet Packet
	packet.Protocol = int(proto)
	packet.Packet = b
	if p.Interface == nil {
		Log(Error, "TAP Interface not initialized")
		return
	}
	err := p.Interface.WritePacket(&packet)
	if err != nil {
		Log(Error, "Failed to write to TAP Interface: %v", err)
	}
}

// ParseIntroString receives a comma-separated string with ID, MAC and IP of a peer
// and returns this data
func (p *PeerToPeer) ParseIntroString(intro string) (*PeerHandshake, error) {
	hs := &PeerHandshake{}
	parts := strings.Split(intro, ",")
	if len(parts) != 4 {
		return nil, fmt.Errorf("Failed to parse introduction string: %s", intro)
	}
	hs.ID = parts[0]
	// Extract MAC
	var err error
	hs.HardwareAddr, err = net.ParseMAC(parts[1])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse MAC address from introduction packet: %v", err)
	}
	// Extract IP
	hs.IP = net.ParseIP(parts[2])
	if hs.IP == nil {
		return nil, fmt.Errorf("Failed to parse IP address from introduction packet")
	}
	hs.Endpoint, err = net.ResolveUDPAddr("udp4", parts[3])
	if err != nil {
		return nil, fmt.Errorf("Failed to parse handshake endpoint: %s", parts[3])
	}

	return hs, nil
}

// SendTo sends a p2p packet by MAC address
func (p *PeerToPeer) SendTo(dst net.HardwareAddr, msg *P2PMessage) (int, error) {
	endpoint, err := p.Peers.GetEndpoint(dst.String())
	if err == nil && endpoint != nil {
		size, err := p.UDPSocket.SendMessage(msg, endpoint)
		return size, err
	}
	return 0, nil
}

// Close stops current instance
func (p *PeerToPeer) Close() error {
	for i, ip := range ActiveInterfaces {
		if ip.Equal(p.Interface.GetIP()) {
			ActiveInterfaces = append(ActiveInterfaces[:i], ActiveInterfaces[i+1:]...)
			break
		}
	}
	hash := p.Dht.NetworkHash
	Log(Info, "Stopping instance %s", hash)
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
	Log(Debug, "All peers under this instance has been removed")

	p.Shutdown = true
	err := p.Dht.Close()
	if err != nil {
		Log(Error, "Failed to stop DHT: %s", err)
	}
	p.UDPSocket.Stop()

	if p.Interface != nil {
		err := p.Interface.Close()
		if err != nil {
			Log(Error, "Failed to close TAP interface: %s", err)
		}
	}
	p.ReadyToStop = true
	Log(Info, "Instance %s stopped", hash)
	return nil
}
