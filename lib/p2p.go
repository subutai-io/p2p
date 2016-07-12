package ptp

import (
	//"bytes"
	"crypto/rand"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"runtime"
	"strings"
	"sync"
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
	NetworkPeers    map[string]*NetworkPeer              // Knows peers
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
	DHTPeerChannel  chan []PeerIP
	ProxyChannel    chan Forwarder
	RemovePeer      chan string
	MessageBuffer   map[string]map[uint16]map[uint16][]byte
	MessagePacket   map[string]map[uint16][]byte
	BufferLock      sync.Mutex
	PeersLock       sync.Mutex
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
		Log(WARNING, "Failed to load config: %v", err)
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

	// Windows returns a real mac here. However, other systems should return empty string
	mac = ExtractMacFromInterface(p.Device)
	if mac != "" {
		p.Mac = mac
		p.HardwareAddr, _ = net.ParseMAC(mac)
	}

	err = ConfigureInterface(p.Device, p.IP, p.Mac, p.DeviceName, p.IPTool)
	if err != nil {
		return err
	}
	return nil
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
	/*
		dhtClient := new(DHTClient)
		config := dhtClient.DHTClientConfig()
		config.NetworkHash = argHash
		config.Mode = MODE_CLIENT
	*/

	p := new(PTPCloud)
	p.FindNetworkAddresses()
	p.HardwareAddr = hw
	p.NetworkPeers = make(map[string]*NetworkPeer)
	p.IPIDTable = make(map[string]string)
	p.MACIDTable = make(map[string]string)
	p.MessageBuffer = make(map[string]map[uint16]map[uint16][]byte)
	p.MessagePacket = make(map[string]map[uint16][]byte)

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
	p.PacketHandlers[PT_LLDP] = p.handlePacketLLDP

	p.UDPSocket = new(PTPNet)
	p.UDPSocket.Init("", port)
	port = p.UDPSocket.GetPort()
	Log(INFO, "Started UDP Listener at port %d", port)
	/*
		config.P2PPort = port
		if argDht != "" {
			config.Routers = argDht
		}
	*/
	// TODO: Move channels inside DHT
	p.DHTPeerChannel = make(chan []PeerIP)
	p.ProxyChannel = make(chan Forwarder)
	p.StartDHT(argHash, argDht)
	/*
			p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel, p.ProxyChannel)
		for p.Dht == nil {
			Log(WARNING, "Failed to connect to DHT. Retrying in 5 seconds")
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

func (p *PTPCloud) StartDHT(hash, routers string) {
	dhtClient := new(DHTClient)
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = hash
	config.Mode = MODE_CLIENT
	config.P2PPort = p.UDPSocket.GetPort()
	if routers != "" {
		config.Routers = routers
	}
	p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel, p.ProxyChannel)
	for p.Dht == nil {
		Log(WARNING, "Failed to connect to DHT. Retrying in 5 seconds")
		time.Sleep(5 * time.Second)
		p.LocalIPs = p.LocalIPs[:0]
		p.FindNetworkAddresses()
		p.Dht = dhtClient.Initialize(config, p.LocalIPs, p.DHTPeerChannel, p.ProxyChannel)
	}
	Log(INFO, "ID assigned. Continue")
}

func (p *PTPCloud) Run() {
	go p.ReadDHTPeers()
	go p.ReadProxies()
	go func() {
		for {
			if p.Shutdown {
				break
			}
			rm := <-p.Dht.RemovePeerChan
			if rm == "DUMMY" {
				continue
			}
			p.PeersLock.Lock()
			peer, exists := p.NetworkPeers[rm]
			p.PeersLock.Unlock()
			runtime.Gosched()
			if exists {
				Log(INFO, "Stopping %s after STOP command", rm)
				peer.State = P_DISCONNECT
				p.PeersLock.Lock()
				p.NetworkPeers[rm] = peer
				p.PeersLock.Unlock()
				runtime.Gosched()
			} else {
				Log(INFO, "Can't stop peer. ID not found")
			}
		}
		Log(INFO, "Stopping peer state listener")
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
			if peer.State == P_STOP {
				Log(INFO, "Removing peer %s", i)
				time.Sleep(100 * time.Microsecond)
				delete(p.IPIDTable, peer.PeerLocalIP.String())
				delete(p.MACIDTable, peer.PeerHW.String())

				p.PeersLock.Lock()
				delete(p.NetworkPeers, i)
				p.PeersLock.Unlock()
				runtime.Gosched()
			}
		}
		passed := time.Since(p.Dht.LastDHTPing)
		interval := time.Duration(time.Second * 50)
		if passed > interval {
			Log(ERROR, "Lost connection to DHT")
			p.Dht.Shutdown = true
			p.Dht.ID = ""
			hash := p.Dht.NetworkHash
			routers := p.Dht.Routers
			time.Sleep(time.Second * 5)
			p.StartDHT(hash, routers)
			go p.Dht.UpdatePeers()
		}
	}
	Log(INFO, "Shutting down instance %s completed", p.Dht.NetworkHash)
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
			p.PeersLock.Lock()
			delete(p.NetworkPeers, i)
			p.PeersLock.Unlock()
			runtime.Gosched()
		}
	}
	return
}

func (p *PTPCloud) SyncForwarders() int {
	var count int = 0
	for _, fwd := range p.Dht.Forwarders {
		for key, peer := range p.NetworkPeers {
			if peer.Endpoint == nil && fwd.DestinationID == peer.ID && peer.Forwarder == nil {
				Log(INFO, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr
				peer.Forwarder = fwd.Addr
				peer.State = P_HANDSHAKING_FORWARDER
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

func (p *PTPCloud) ParseIntroString(intro string) (string, net.HardwareAddr, net.IP) {
	parts := strings.Split(intro, ",")
	if len(parts) != 3 {
		Log(ERROR, "Failed to parse introduction string: %s", intro)
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
		go callback(msg, src_addr)
	} else {
		Log(WARNING, "Unknown message received")
	}
}

func (p *PTPCloud) HandleNotEncryptedMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	Log(TRACE, "Data: %s, Proto: %d, From: %s", msg.Data, msg.Header.NetProto, src_addr.String())
	p.BufferLock.Lock()
	// Allocate memory
	if p.MessageBuffer[src_addr.String()] == nil {
		p.MessageBuffer[src_addr.String()] = make(map[uint16]map[uint16][]byte)
	}
	if p.MessageBuffer[src_addr.String()][msg.Header.Id] == nil {
		p.MessageBuffer[src_addr.String()][msg.Header.Id] = make(map[uint16][]byte)
	}
	// Append packet contents into queue
	p.MessageBuffer[src_addr.String()][msg.Header.Id][msg.Header.Seq] = msg.Data
	p.BufferLock.Unlock()
	runtime.Gosched()
	if msg.Header.Complete > 0 {
		wcounter := 0
		p.BufferLock.Lock()
		plen := len(p.MessageBuffer[src_addr.String()][msg.Header.Id])
		p.BufferLock.Unlock()
		runtime.Gosched()
		// Wait for packet to arrive
		for plen != int(msg.Header.Complete) {
			time.Sleep(10 * time.Millisecond)
			p.BufferLock.Lock()
			plen = len(p.MessageBuffer[src_addr.String()][msg.Header.Id])
			p.BufferLock.Unlock()
			runtime.Gosched()
			wcounter++
			if wcounter > 100 {
				Log(WARNING, "Packet incomplete. Received %d from %d [%d]", plen, msg.Header.Complete, msg.Header.Id)
				p.BufferLock.Lock()
				delete(p.MessageBuffer[src_addr.String()], msg.Header.Id)
				p.BufferLock.Unlock()
				runtime.Gosched()
				return
			}
		}
		// Combine packet from parts
		var b []byte
		for i := uint16(1); i <= msg.Header.Complete; i++ {
			p.BufferLock.Lock()
			data, exists := p.MessageBuffer[src_addr.String()][msg.Header.Id][i]
			p.BufferLock.Unlock()
			runtime.Gosched()
			if exists {
				b = append(b, data...)
			} else {
				Log(ERROR, "Missing packet: %d/%d", i, msg.Header.Complete)
				p.BufferLock.Lock()
				delete(p.MessageBuffer[src_addr.String()], msg.Header.Id)
				p.BufferLock.Unlock()
				runtime.Gosched()
				return
			}
		}
		p.WriteToDevice(b, msg.Header.NetProto, false)
		p.BufferLock.Lock()
		delete(p.MessageBuffer[src_addr.String()], msg.Header.Id)
		p.BufferLock.Unlock()
		runtime.Gosched()
	}
}

func (p *PTPCloud) HandlePingMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	p.UDPSocket.SendMessage(msg, src_addr)
}

func (p *PTPCloud) HandleXpeerPingMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	pt := PingType(msg.Header.NetProto)
	if pt == PING_REQ {
		Log(DEBUG, "Ping request received")
		// Send a PING response
		r := CreateXpeerPingMessage(PING_RESP, p.HardwareAddr.String())
		addr, err := net.ParseMAC(string(msg.Data))
		if err != nil {
			Log(ERROR, "Failed to parse MAC address in crosspeer ping message")
		} else {
			p.SendTo(addr, r)
			Log(DEBUG, "Sending to %s", addr.String())
		}
	} else {
		Log(DEBUG, "Ping response received")
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

func (p *PTPCloud) HandleIntroMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	p.PeersLock.Lock()
	peer, exists := p.NetworkPeers[id]
	p.PeersLock.Unlock()
	runtime.Gosched()
	if !exists {
		Log(DEBUG, "Received introduction confirmation from unknown peer: %s", id)
		p.Dht.SendUpdateRequest()
		return
	}
	peer.PeerHW = mac
	peer.PeerLocalIP = ip
	peer.State = P_CONNECTED
	peer.LastContact = time.Now()
	p.IPIDTable[ip.String()] = id
	p.MACIDTable[mac.String()] = id
	p.PeersLock.Lock()
	p.NetworkPeers[id] = peer
	p.PeersLock.Unlock()
	runtime.Gosched()
	Log(INFO, "Connection with peer %s has been established", id)
}

func (p *PTPCloud) HandleIntroRequestMessage(msg *P2PMessage, src_addr *net.UDPAddr) {
	id := string(msg.Data)
	peer, exists := p.NetworkPeers[id]
	if !exists {
		Log(DEBUG, "Introduction request came from unknown peer: %s", id)
		p.Dht.SendUpdateRequest()
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
	if msg.Header.ProxyId < 1 {
		return
	}
	ip := string(msg.Data)
	Log(INFO, "Proxy confirmation received from %s. Tunnel ID %d", ip, int(msg.Header.ProxyId))
	for key, peer := range p.NetworkPeers {
		if peer.PeerAddr.String() == ip {
			peer.ProxyID = int(msg.Header.ProxyId)
			p.PeersLock.Lock()
			p.NetworkPeers[key] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
			return
		}
	}
	Log(WARNING, "Can't set Tunnel#%d for %s: Can't find address", int(msg.Header.ProxyId), ip)
}

func (p *PTPCloud) HandleBadTun(msg *P2PMessage, src_addr *net.UDPAddr) {
	for key, peer := range p.NetworkPeers {
		if peer.ProxyID == int(msg.Header.ProxyId) && peer.Endpoint.String() == src_addr.String() {
			Log(DEBUG, "Cleaning bad tunnel %d from %s", msg.Header.ProxyId, src_addr.String())
			peer.ProxyID = 0
			peer.Endpoint = nil
			peer.Forwarder = nil
			peer.PeerAddr = nil
			peer.State = P_INIT
			p.PeersLock.Lock()
			p.NetworkPeers[key] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
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
		p.PeersLock.Lock()
		peer, exists := p.NetworkPeers[id]
		p.PeersLock.Unlock()
		runtime.Gosched()
		if exists {
			msg.Header.ProxyId = uint16(peer.ProxyID)
			Log(DEBUG, "Sending to %s via proxy id %d", dst.String(), msg.Header.ProxyId)
			size, err := p.UDPSocket.SendMessage(msg, peer.Endpoint)
			return size, err
		}
	}
	return 0, nil
}

func (p *PTPCloud) StopInstance() {
	for i, peer := range p.NetworkPeers {
		peer.State = P_DISCONNECT
		p.PeersLock.Lock()
		p.NetworkPeers[i] = peer
		p.PeersLock.Unlock()
		runtime.Gosched()
	}
	var ip net.IP
	if p.Dht == nil || p.Dht.Network == nil {
		Log(WARNING, "DHT isn't in use")
	} else {
		ip = p.Dht.Network.IP
	}
	p.Dht.Stop()
	p.UDPSocket.Stop()
	p.Shutdown = true
	var peers []PeerIP
	var proxy Forwarder
	p.DHTPeerChannel <- peers
	p.ProxyChannel <- proxy
	Log(INFO, "Stopping P2P Message handler")
	// Tricky part: we need to send a message to ourselves to quit blocking operation
	msg := CreateTestP2PMessage(p.Crypter, "STOP", 1)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", p.Dht.P2PPort))
	p.UDPSocket.SendMessage(msg, addr)
	var ipIt int = 200
	if ip != nil {
		for p.IsDeviceExists(p.DeviceName) {
			time.Sleep(1 * time.Second)
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
	}
	time.Sleep(3 * time.Second)
	p.ReadyToStop = true
}

func (p *PTPCloud) ReadDHTPeers() {
	for {
		if p.Shutdown {
			break
		}
		peers := <-p.DHTPeerChannel
		p.UpdatePeers(peers)
	}
	Log(INFO, "Stopped DHT reader channel")
}

func (p *PTPCloud) ReadProxies() {
	for {
		if p.Shutdown {
			break
		}
		proxy := <-p.ProxyChannel
		exists := false
		for i, peer := range p.NetworkPeers {
			if i == proxy.DestinationID {
				peer.State = P_HANDSHAKING_FORWARDER
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
			Log(INFO, "Received forwarder for unknown peer")
			p.Dht.SendUpdateRequest()
		}
	}
	Log(INFO, "Stopped Proxy reader channel")
}

func (p *PTPCloud) UpdatePeers(peers []PeerIP) {
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
			peer.State = P_INIT
			p.PeersLock.Lock()
			p.NetworkPeers[newPeer.ID] = peer
			p.PeersLock.Unlock()
			runtime.Gosched()
			go p.NetworkPeers[newPeer.ID].Run(p)
		}
	}
}
