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

var GlobalIPBlacklist []string

// MessageHandler is a messages callback
type MessageHandler func(message *P2PMessage, srcAddr *net.UDPAddr)

// PeerToPeer - Main structure
type PeerToPeer struct {
	IP              string                               // Interface IP address
	Mac             string                               // String representation of a MAC address
	HardwareAddr    net.HardwareAddr                     // MAC address of network interface
	Mask            string                               // Network mask in the dot-decimal notation
	DeviceName      string                               // Name of the network interface
	IPTool          string                               `yaml:"iptool"` // Network interface configuration tool
	Device          *Interface                           // Network interface
	NetworkPeers    map[string]*NetworkPeer              // Knows peers
	UDPSocket       *Network                             // Peer-to-peer interconnection socket
	LocalIPs        []net.IP                             // List of IPs available in the system
	Dht             *DHTClient                           // DHT Client
	Crypter         Crypto                               // Instance of crypto
	Shutdown        bool                                 // Set to true when instance in shutdown mode
	Restart         bool                                 // Instance will be restarted
	ForwardMode     bool                                 // Skip local peer discovery
	ReadyToStop     bool                                 // Set to true when instance is ready to stop
	IPIDTable       map[string]string                    // Mapping for IP->ID
	MACIDTable      map[string]string                    // Mapping for MAC->ID
	MessageHandlers map[uint16]MessageHandler            // Callbacks
	PacketHandlers  map[PacketType]PacketHandlerCallback // Callbacks for network packet handlers
	//DHTPeerChannel  chan []PeerIP
	//ProxyChannel    chan Forwarder
	RemovePeer      chan string
	MessageBuffer   map[string]map[uint16]map[uint16][]byte
	MessageLifetime map[string]map[uint16]time.Time
	MessagePacket   map[string][]byte
	BufferLock      sync.Mutex
	PeersLock       sync.Mutex
}

// AssignInterface - Creates TUN/TAP Interface and configures it with provided IP tool
func (p *PeerToPeer) AssignInterface(ip, mac, mask, device string) error {
	var err error

	for _, i := range GlobalIPBlacklist {
		if i == ip {
			Log(Error, "Can't assign IP Address: IP %s is already in use", ip)
			return fmt.Errorf("Can't assign IP Address: IP %s is already in use", ip)
		}
	}
	GlobalIPBlacklist = append(GlobalIPBlacklist, ip)

	p.IP = ip
	p.Mac = mac
	p.Mask = mask
	p.DeviceName = device

	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile(ConfigDir + "/p2p/config.yaml")
	if err != nil {
		Log(Warning, "Failed to load config: %v", err)
		p.IPTool = "/sbin/ip"
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		Log(Error, "Failed to parse config: %v", err)
		return err
	}

	p.Device, err = Open(p.DeviceName, DevTap)
	if p.Device == nil {
		Log(Error, "Failed to open TAP device %s: %v", device, err)
		return err
	}
	Log(Info, "%v TAP Device created", p.DeviceName)

	// Windows returns a real mac here. However, other systems should return empty string
	mac = ExtractMacFromInterface(p.Device)
	if mac != "" {
		p.Mac = mac
		p.HardwareAddr, _ = net.ParseMAC(mac)
	}

	err = ConfigureInterface(p.Device, p.IP, p.Mac, p.DeviceName, p.IPTool)
	Log(Info, "Interface has been configured")
	return err
}

// ListenInterface - Listens TAP interface for incoming packets
func (p *PeerToPeer) ListenInterface() {
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
			Log(Error, "Reading packet %s", err)
		}
		if packet.Truncated {
			Log(Debug, "Truncated packet")
		}
		// TODO: Make handlePacket as a part of PeerToPeer
		go p.handlePacket(packet.Packet, packet.Protocol)
	}
	p.Device.Close()
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
			for _, i := range GlobalIPBlacklist {
				if i == ip.String() {
					decision = "Ignoring"
				}
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
func StartP2PInstance(argIP, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool, port int) *PeerToPeer {

	var hw net.HardwareAddr

	if argMac != "" {
		var err2 error
		hw, err2 = net.ParseMAC(argMac)
		if err2 != nil {
			Log(Error, "Invalid MAC address provided: %v", err2)
			return nil
		}
	} else {
		argMac, hw = GenerateMAC()
		Log(Info, "Generate MAC for TAP device: %s", argMac)
	}

	// Create new DHT Client, configured it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code
	/*
		dhtClient := new(DHTClient)
		config := dhtClient.DHTClientConfig()
		config.NetworkHash = argHash
		config.Mode = DHTModeClient
	*/

	p := new(PeerToPeer)
	p.FindNetworkAddresses()
	p.HardwareAddr = hw
	p.NetworkPeers = make(map[string]*NetworkPeer)
	p.IPIDTable = make(map[string]string)
	p.MACIDTable = make(map[string]string)
	p.MessageBuffer = make(map[string]map[uint16]map[uint16][]byte)
	p.MessageLifetime = make(map[string]map[uint16]time.Time)
	p.MessagePacket = make(map[string][]byte)

	if fwd {
		p.ForwardMode = true
	}

	if argDev == "" {
		argDev = p.GenerateDeviceName(1)
	} else {
		if len(argDev) > 12 {
			Log(Info, "Interface name length should be 12 symbols max")
			return nil
		}
	}
	if p.IsDeviceExists(argDev) {
		Log(Error, "Interface is already in use. Can't create duplicate")
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
		Log(Info, "Traffic encryption is enabled. Key valid until %s", p.Crypter.ActiveKey.Until.String())
	} else {
		Log(Info, "No AES key were provided. Traffic encryption is disabled")
	}

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

	p.UDPSocket = new(Network)
	p.UDPSocket.Init("", port)
	port = p.UDPSocket.GetPort()
	Log(Info, "Started UDP Listener at port %d", port)
	/*
		config.P2PPort = port
		if argDht != "" {
			config.Routers = argDht
		}
	*/
	// TODO: Move channels inside DHT
	//p.DHTPeerChannel = make(chan []PeerIP)
	//p.ProxyChannel = make(chan Forwarder)
	p.StartDHT(argHash, argDht)
	/*
		p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel, p.ProxyChannel)
		for p.Dht == nil {
			Log(Warning, "Failed to connect to DHT. Retrying in 5 seconds")
			time.Sleep(5 * time.Second)
			p.LocalIPs = p.LocalIPs[:0]
			p.FindNetworkAddresses()
			p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel, p.ProxyChannel)
		}
		// Wait for ID
		for len(p.Dht.ID) < 32 {
			time.Sleep(100 * time.Millisecond)
		}
	*/
	var retries = 0
	if argIP == "dhcp" {
		Log(Info, "Requesting IP")
		p.Dht.RequestIP()
		time.Sleep(1 * time.Second)
		for p.Dht.IP == nil && p.Dht.Network == nil {
			Log(Info, "No IP were received. Requesting again")
			p.Dht.RequestIP()
			time.Sleep(3 * time.Second)
			retries++
			if retries >= 10 {
				Log(Error, "Failed to retrieve IP from network after 10 retries")
				return nil
			}
		}
		m := p.Dht.Network.Mask
		mask := fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
		p.AssignInterface(p.Dht.IP.String(), argMac, mask, argDev)
	} else {
		ip, ipnet, err := net.ParseCIDR(argIP)
		if err != nil {
			nip := net.ParseIP(argIP)
			if nip == nil {
				Log(Error, "Invalid address were provided for network interface. Use -ip \"dhcp\" or specify correct IP address")
				return nil
			}
			argIP += `/24`
			Log(Warning, "No CIDR mask was provided. Assumming /24")
			ip, ipnet, err = net.ParseCIDR(argIP)
			if err != nil {
				Log(Error, "Failed to setup provided IP address for local device")
				return nil
			}
		}
		p.Dht.IP = ip
		p.Dht.Network = ipnet
		mask := fmt.Sprintf("%d.%d.%d.%d", ipnet.Mask[0], ipnet.Mask[1], ipnet.Mask[2], ipnet.Mask[3])
		p.Dht.SendIP(argIP, mask)
		err = p.AssignInterface(p.Dht.IP.String(), argMac, mask, argDev)
		if err != nil {
			Log(Error, "Can't configure interface")
			return nil
		}
	}

	go p.UDPSocket.Listen(p.HandleP2PMessage)

	go p.ListenInterface()
	return p
}

// StartDHT starts a DHT client
func (p *PeerToPeer) StartDHT(hash, routers string) {
	dhtClient := new(DHTClient)
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = hash
	config.Mode = DHTModeClient
	config.P2PPort = p.UDPSocket.GetPort()
	if routers != "" {
		config.Routers = routers
	}
	p.Dht = dhtClient.Initialize(config, p.LocalIPs, nil, nil)
	for p.Dht == nil {
		Log(Warning, "Failed to connect to DHT. Retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		p.LocalIPs = p.LocalIPs[:0]
		p.FindNetworkAddresses()
		p.Dht = dhtClient.Initialize(config, p.LocalIPs, nil, nil)
	}
	Log(Info, "ID assigned. Continue")
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
			case rm, r := <-p.Dht.RemovePeerChan:
				if r {
					if rm == "DUMMY" || rm == "" {
						continue
					}
					p.PeersLock.Lock()
					peer, exists := p.NetworkPeers[rm]
					p.PeersLock.Unlock()
					runtime.Gosched()
					if exists {
						Log(Info, "Stopping %s after STOP command", rm)
						peer.State = PeerStateDisconnect
						p.PeersLock.Lock()
						p.NetworkPeers[rm] = peer
						p.PeersLock.Unlock()
						runtime.Gosched()
					} else {
						Log(Info, "Can't stop peer. ID not found")
					}
				} else {
					Log(Trace, "Channel was closed")
				}
			default:
				time.Sleep(100 * time.Microsecond)
			}
			//rm := <-p.Dht.RemovePeerChan
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
				time.Sleep(100 * time.Microsecond)
				lip := peer.PeerLocalIP.String()
				delete(p.IPIDTable, lip)
				delete(p.MACIDTable, peer.PeerHW.String())

				k := 0
				for k, ipb := range GlobalIPBlacklist {
					if ipb == lip {
						//GlobalIPBlacklist = append(GlobalIPBlacklist[:k], GlobalIPBlacklist[k+1:]...)
						GlobalIPBlacklist[k] = ipb
						k++
					}
				}
				GlobalIPBlacklist = GlobalIPBlacklist[:k]

				delete(p.NetworkPeers, i)
				runtime.Gosched()
				Log(Info, "Remove complete")
			}
		}
		passed := time.Since(p.Dht.LastDHTPing)
		interval := time.Duration(time.Second * 45)
		if passed > interval {
			Log(Error, "Lost connection to DHT")
			p.Dht.Shutdown = true
			p.Dht.ID = ""
			hash := p.Dht.NetworkHash
			routers := p.Dht.Routers
			time.Sleep(time.Second * 5)
			p.StartDHT(hash, routers)
			go p.Dht.UpdatePeers()
		}
	}
	Log(Info, "Shutting down instance %s completed", p.Dht.NetworkHash)
}

// PrepareIntroductionMessage collects client ID, mac and IP address
// and create a comma-separated line
func (p *PeerToPeer) PrepareIntroductionMessage(id string) *P2PMessage {
	var intro = id + "," + p.Mac + "," + p.IP
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

// WriteToDevice writes data to created TUN/TAP device
func (p *PeerToPeer) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var packet Packet
	packet.Protocol = int(proto)
	packet.Truncated = truncated
	packet.Packet = b
	if p.Device == nil {
		Log(Error, "TUN/TAP Device not initialized")
		return
	}
	err := p.Device.WritePacket(&packet)
	if err != nil {
		Log(Error, "Failed to write to TUN/TAP device: %v", err)
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
		r := CreateXpeerPingMessage(PingResp, p.HardwareAddr.String())
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
	GlobalIPBlacklist = append(GlobalIPBlacklist, ip.String())
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
	for len(p.NetworkPeers) > 0 {
		time.Sleep(100 * time.Microsecond)
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
		for p.IsDeviceExists(p.DeviceName) {
			time.Sleep(1 * time.Second)
			target := fmt.Sprintf("%d.%d.%d.%d:99", ip[0], ip[1], ip[2], ipIt)
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
			time.Sleep(100 * time.Microsecond)
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
				Log(Trace, "Clossed channel")
			}
		default:
			time.Sleep(100 * time.Microsecond)
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
