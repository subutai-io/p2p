package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/danderson/tuntap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"p2p/commons"
	"p2p/dht"
	log "p2p/p2p_log"
	"p2p/udpcs"
	"strings"
	"time"
)

type MSG_TYPE uint16

type MessageHandler func(message *udpcs.P2PMessage, src_addr *net.UDPAddr)

// Main structure
type PTPCloud struct {
	// IP Address assigned to device at startup
	IP string

	// MAC Address assigned to device or generated by the application (TODO: Implement random generation and MAC assignment)
	Mac string

	HardwareAddr net.HardwareAddr

	// Netmask for device
	Mask string

	// Name of the device
	DeviceName string

	// Path to tool that is used to configure network device (only "ip" tools is supported at this moment)
	IPTool string `yaml:"iptool"`

	// TUN/TAP Interface
	Interface *os.File

	// Representation of TUN/TAP Device
	Device *tuntap.Interface

	//NetworkPeers []NetworkPeer
	NetworkPeers map[string]NetworkPeer

	UDPSocket *udpcs.UDPClient

	LocalIPs []net.IP

	dht *dht.DHTClient

	Crypter udpcs.Crypto

	// If true, instance will shutdown itself on a next iteration
	Shutdown bool

	// IP -> ID Table for faster ARP lookup
	IPIDTable map[string]string

	// If yes, client will not try to establish direct connection over LAN and will always switch to proxy
	ForwardMode bool

	MessageHandlers map[uint16]MessageHandler

	ReadyToStop bool

	PacketHandlers map[PacketType]PacketHandlerCallback
}

type NetworkPeer struct {
	// ID of the node received from DHT Bootstrap node
	ID string
	// Whether informaton about this node is filled or not
	// Normally it should be filled after peer-to-peer handshake procedure
	Unknown bool
	// This variables indicates whether handshake mechanism was started or not
	Handshaked bool
	// ID of the proxy used to communicate with the node
	ProxyID   int
	Forwarder *net.UDPAddr
	PeerAddr  *net.UDPAddr
	// IP of the peer we are connected to.
	PeerLocalIP net.IP
	// Hardware address of node's TUN/TAP device
	PeerHW net.HardwareAddr
	// Endpoint is the same as CleanAddr TODO: Remove CleanAddr
	Endpoint string
	// List of peer IP addresses
	KnownIPs []*net.UDPAddr
	// Number of retries of introduce
	Retries int
}

// Creates TUN/TAP Interface and configures it with provided IP tool
func (ptp *PTPCloud) CreateDevice(ip, mac, mask, device string) error {
	var err error

	ptp.IP = ip
	ptp.Mac = mac
	ptp.Mask = mask
	ptp.DeviceName = device

	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Log(log.ERROR, "Failed to load config: %v", err)
		ptp.IPTool = "/sbin/ip"
	}
	err = yaml.Unmarshal(yamlFile, ptp)
	if err != nil {
		log.Log(log.ERROR, "Failed to parse config: %v", err)
		return err
	}

	ptp.Device, err = tuntap.Open(ptp.DeviceName, tuntap.DevTap)
	if ptp.Device == nil {
		log.Log(log.ERROR, "Failed to open TAP device: %v", err)
		return err
	} else {
		log.Log(log.INFO, "%v TAP Device created", ptp.DeviceName)
	}

	linkup := exec.Command(ptp.IPTool, "link", "set", "dev", ptp.DeviceName, "up")
	err = linkup.Run()
	if err != nil {
		log.Log(log.ERROR, "Failed to up link: %v", err)
		return err
	}

	// Configure new device
	log.Log(log.INFO, "Setting %s IP on device %s", ptp.IP, ptp.DeviceName)
	setip := exec.Command(ptp.IPTool, "addr", "add", ptp.IP+"/24", "dev", ptp.DeviceName)
	err = setip.Run()
	if err != nil {
		log.Log(log.ERROR, "Failed to set IP: %v", err)
		return err
	}

	// Set MAC to device
	log.Log(log.INFO, "Setting %s MAC on device %s", mac, ptp.DeviceName)
	setmac := exec.Command(ptp.IPTool, "link", "set", "dev", ptp.DeviceName, "address", mac)
	err = setmac.Run()
	if err != nil {
		log.Log(log.ERROR, "Failed to set MAC: %v", err)
		return err
	}
	return nil
}

// Handles a packet that was received by TUN/TAP device
// Receiving a packet by device means that some application sent a network
// packet within a subnet in which our application works.
// This method calls appropriate gorouting for extracted packet protocol
func (ptp *PTPCloud) handlePacket(contents []byte, proto int) {
	callback, exists := ptp.PacketHandlers[PacketType(proto)]
	if exists {
		callback(contents, proto)
	} else {
		log.Log(log.WARNING, "Captured undefined packet")
	}
}

// Listen TAP interface for incoming packets
func (ptp *PTPCloud) ListenInterface() {
	// Read packets received by TUN/TAP device and send them to a handlePacket goroutine
	// This goroutine will decide what to do with this packet
	for {
		if ptp.Shutdown {
			break
		}
		packet, err := ptp.Device.ReadPacket()
		if err != nil {
			log.Log(log.ERROR, "Reading packet %s", err)
		}
		if packet.Truncated {
			log.Log(log.DEBUG, "Truncated packet")
		}
		// TODO: Make handlePacket as a part of PTPCloud
		go ptp.handlePacket(packet.Packet, packet.Protocol)
	}
	ptp.Device.Close()
	log.Log(log.INFO, "Shutting down interface listener")
}

// This method will generate device name if none were specified at startup
func (ptp *PTPCloud) GenerateDeviceName(i int) string {
	var devName string = "vptp" + fmt.Sprintf("%d", i)
	inf, err := net.Interfaces()
	if err != nil {
		log.Log(log.ERROR, "Failed to retrieve list of network interfaces")
		return ""
	}
	var exist bool = false
	for _, i := range inf {
		if i.Name == devName {
			exist = true
		}
	}
	if exist {
		return ptp.GenerateDeviceName(i + 1)
	} else {
		return devName
	}
}

// This method lists interfaces available in the system and retrieves their
// IP addresses
func (ptp *PTPCloud) FindNetworkAddresses() {
	log.Log(log.INFO, "Looking for available network interfaces")
	inf, err := net.Interfaces()
	if err != nil {
		log.Log(log.ERROR, "Failed to retrieve list of network interfaces")
		return
	}
	for _, i := range inf {
		addresses, err := i.Addrs()

		if err != nil {
			log.Log(log.ERROR, "Failed to retrieve address for interface. %v", err)
			continue
		}
		for _, addr := range addresses {
			var decision string = "Ignoring"
			var ipType string = "Unknown"
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				log.Log(log.ERROR, "Failed to parse CIDR notation: %v", err)
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
			log.Log(log.INFO, "Interface %s: %s. Type: %s. %s", i.Name, addr.String(), ipType, decision)
			if decision == "Saving" {
				ptp.LocalIPs = append(ptp.LocalIPs, ip)
			}
		}
	}
	log.Log(log.INFO, "%d interfaces were saved", len(ptp.LocalIPs))
}

func p2pmain(argIp, argMask, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog string, fwd bool) *PTPCloud {

	if argIp == "none" {
		fmt.Println("USAGE: p2p [OPTIONS]")
		fmt.Printf("\nOPTIONS:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var hw net.HardwareAddr

	if argMac != "" {
		var err2 error
		hw, err2 = net.ParseMAC(argMac)
		if err2 != nil {
			log.Log(log.ERROR, "Invalid MAC address provided: %v", err2)
			return nil
		}
	} else {
		argMac, hw = GenerateMAC()
		log.Log(log.INFO, "Generate MAC for TAP device: %s", argMac)
	}

	// Create new DHT Client, configured it and initialize
	// During initialization procedure, DHT Client will send
	// a introduction packet along with a hash to a DHT bootstrap
	// nodes that was hardcoded into it's code
	dhtClient := new(dht.DHTClient)
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = argHash
	config.Mode = dht.MODE_CLIENT

	ptp := new(PTPCloud)
	ptp.FindNetworkAddresses()
	ptp.HardwareAddr = hw
	ptp.NetworkPeers = make(map[string]NetworkPeer)
	ptp.IPIDTable = make(map[string]string)

	if fwd {
		ptp.ForwardMode = true
	}

	if argDev == "" {
		argDev = ptp.GenerateDeviceName(1)
	}

	if argKeyfile != "" {
		ptp.Crypter.ReadKeysFromFile(argKeyfile)
	}
	if argKey != "" {
		// Override key from file
		if argTTL == "" {
			argTTL = "default"
		}
		var newKey udpcs.CryptoKey
		newKey = ptp.Crypter.EncrichKeyValues(newKey, argKey, argTTL)
		ptp.Crypter.Keys = append(ptp.Crypter.Keys, newKey)
		ptp.Crypter.ActiveKey = ptp.Crypter.Keys[0]
		ptp.Crypter.Active = true
	}

	if ptp.Crypter.Active {
		log.Log(log.INFO, "Traffic encryption is enabled. Key valid until %s", ptp.Crypter.ActiveKey.Until.String())
	} else {
		log.Log(log.INFO, "No AES key were provided. Traffic encryption is disabled")
	}

	// Register network message handlers
	ptp.MessageHandlers = make(map[uint16]MessageHandler)
	ptp.MessageHandlers[commons.MT_NENC] = ptp.HandleNotEncryptedMessage
	ptp.MessageHandlers[commons.MT_PING] = ptp.HandlePingMessage
	ptp.MessageHandlers[commons.MT_ENC] = ptp.HandleMessage
	ptp.MessageHandlers[commons.MT_INTRO] = ptp.HandleIntroMessage
	ptp.MessageHandlers[commons.MT_PROXY] = ptp.HandleProxyMessage
	ptp.MessageHandlers[commons.MT_TEST] = ptp.HandleTestMessage

	// Register packet handlers
	ptp.PacketHandlers = make(map[PacketType]PacketHandlerCallback)
	ptp.PacketHandlers[PT_PARC_UNIVERSAL] = ptp.handlePARCUniversalPacket
	ptp.PacketHandlers[PT_IPV4] = ptp.handlePacketIPv4
	ptp.PacketHandlers[PT_ARP] = ptp.handlePacketARP
	ptp.PacketHandlers[PT_RARP] = ptp.handleRARPPacket
	ptp.PacketHandlers[PT_8021Q] = ptp.handle8021qPacket
	ptp.PacketHandlers[PT_IPV6] = ptp.handlePacketIPv6
	ptp.PacketHandlers[PT_PPPOE_DISCOVERY] = ptp.handlePPPoEDiscoveryPacket
	ptp.PacketHandlers[PT_PPPOE_SESSION] = ptp.handlePPPoESessionPacket

	ptp.CreateDevice(argIp, argMac, argMask, argDev)
	ptp.UDPSocket = new(udpcs.UDPClient)
	ptp.UDPSocket.Init("", 0)
	port := ptp.UDPSocket.GetPort()
	log.Log(log.INFO, "Started UDP Listener at port %d", port)
	config.P2PPort = port
	if argDht != "" {
		config.Routers = argDht
	}
	ptp.dht = dhtClient.Initialize(config, ptp.LocalIPs)

	go ptp.UDPSocket.Listen(ptp.HandleP2PMessage)

	go ptp.ListenInterface()
	return ptp
}

func (ptp *PTPCloud) Run() {
	for {
		if ptp.Shutdown {
			// TODO: Do it more safely
			if ptp.ReadyToStop {
				break
			}
			time.Sleep(1 * time.Second)
			continue
		}
		time.Sleep(3 * time.Second)
		ptp.dht.UpdatePeers()
		// Wait two seconds before synchronizing with catched peers
		time.Sleep(2 * time.Second)
		ptp.PurgePeers()
		newPeersNum := ptp.SyncPeers()
		newPeersNum = newPeersNum + ptp.SyncForwarders()
		//if newPeersNum > 0 {
		ptp.IntroducePeers()
		//}
	}
	log.Log(log.INFO, "Shutting down instance %s completed", ptp.dht.NetworkHash)
}

// This method sends information about himself to empty peers
// Empty peers is a peer that was not sent us information
// about his device
func (ptp *PTPCloud) IntroducePeers() {
	for i, peer := range ptp.NetworkPeers {
		// Skip if know this peer
		if !peer.Unknown {
			continue
		}
		// Skip if we don't have an endpoint address for this peer
		if peer.Endpoint == "" {
			continue
		}
		if peer.Retries >= 10 {
			log.Log(log.WARNING, "Failed to introduce to %s", peer.ID)
			// TODO: Perform necessary action
		}
		peer.Retries = peer.Retries + 1
		log.Log(log.DEBUG, "Intoducing to %s", peer.Endpoint)
		addr, err := net.ResolveUDPAddr("udp", peer.Endpoint)
		if err != nil {
			log.Log(log.ERROR, "Failed to resolve UDP address during Introduction: %v", err)
			continue
		}
		//peer.PeerAddr = addr
		ptp.NetworkPeers[i] = peer
		// Send introduction packet
		msg := ptp.PrepareIntroductionMessage(ptp.dht.ID)
		msg.Header.ProxyId = uint16(peer.ProxyID)
		_, err = ptp.UDPSocket.SendMessage(msg, addr)
		if err != nil {
			log.Log(log.ERROR, "Failed to send introduction to %s", addr.String())
		} else {
			log.Log(log.DEBUG, "Introduction sent to %s", peer.Endpoint)
		}
	}
}

func (ptp *PTPCloud) PrepareIntroductionMessage(id string) *udpcs.P2PMessage {
	var intro string = id + "," + ptp.Mac + "," + ptp.IP
	msg := udpcs.CreateIntroP2PMessage(ptp.Crypter, intro, 0)
	return msg
}

// This method goes over peers and removes obsolete ones
// Peer becomes obsolete when it goes out of DHT
func (ptp *PTPCloud) PurgePeers() {
	for i, peer := range ptp.NetworkPeers {
		var f bool = false
		for _, newPeer := range ptp.dht.Peers {
			if newPeer.ID == peer.ID {
				f = true
			}
		}
		if !f {
			log.Log(log.DEBUG, ("Peer not found in DHT peer table. Remove it"))
			delete(ptp.IPIDTable, peer.PeerLocalIP.String())
			delete(ptp.NetworkPeers, i)
		}
	}
	return
}

// This method tests connection with specified endpoint
func (ptp *PTPCloud) TestConnection(endpoint *net.UDPAddr) bool {
	msg := udpcs.CreateTestP2PMessage(ptp.Crypter, "TEST", 0)
	conn, err := net.DialUDP("udp4", nil, endpoint)
	if err != nil {
		log.Log(log.ERROR, "%v", err)
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
			log.Log(log.ERROR, "%v", err)
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

func (ptp *PTPCloud) SyncForwarders() int {
	var count int = 0
	for _, fwd := range ptp.dht.Forwarders {
		for key, peer := range ptp.NetworkPeers {
			if peer.Endpoint == "" && fwd.DestinationID == peer.ID {
				log.Log(log.INFO, "Saving control peer as a proxy destination for %s", peer.ID)
				peer.Endpoint = fwd.Addr.String()
				peer.Forwarder = fwd.Addr
				ptp.NetworkPeers[key] = peer
				count = count + 1
			}
		}
	}
	ptp.dht.Forwarders = ptp.dht.Forwarders[:0]
	return count
}

// This method takes a list of catched peers from DHT and
// adds every new peer into list of peers
// Returns amount of peers that has been added
func (ptp *PTPCloud) SyncPeers() int {
	var count int = 0

	for _, id := range ptp.dht.Peers {
		if id.ID == "" {
			continue
		}
		var found bool = false
		for i, peer := range ptp.NetworkPeers {
			if peer.ID == id.ID {
				found = true
				// Check if know something new about this peer, e.g. new addresses were
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
						log.Log(log.INFO, "Adding new IP (%s) address to %s", ip, peer.ID)
						// TODO: Check IP parsing
						newIp, _ := net.ResolveUDPAddr("udp", ip)
						peer.KnownIPs = append(peer.KnownIPs, newIp)
						ptp.NetworkPeers[i] = peer
					}
				}

				// Set and Endpoint from peers if no endpoint were set previously
				if peer.Endpoint == "" {
					// First we need to go over each network and see if some of addresses are inside LAN
					// TODO: Implement
					//var failback bool = false
					interfaces, err := net.Interfaces()
					if err != nil {
						log.Log(log.ERROR, "Failed to retrieve list of network interfaces")
						//failback = true
					}

					for _, inf := range interfaces {
						if ptp.ForwardMode {
							// Don't try to connect over local network
							break
						}
						if ptp.NetworkPeers[i].Endpoint != "" {
							break
						}
						if inf.Name == ptp.DeviceName {
							continue
						}
						addrs, _ := inf.Addrs()
						for _, addr := range addrs {
							netip, network, _ := net.ParseCIDR(addr.String())
							if !netip.IsGlobalUnicast() {
								continue
							}
							for _, kip := range ptp.NetworkPeers[i].KnownIPs {
								log.Log(log.TRACE, "Probing new IP %s against network %s", kip.IP.String(), network.String())

								if network.Contains(kip.IP) {
									if ptp.TestConnection(kip) {
										//ptp.NetworkPeers[i].Endpoint = kip.String()
										peer.Endpoint = kip.String()
										ptp.NetworkPeers[i] = peer
										count = count + 1
										log.Log(log.INFO, "Setting endpoint for %s to %s", peer.ID, kip.String())
									}
									// TODO: Test connection
								}
							}
						}
					}

					// If we still don't have an endpoint we will try to reach peer from outside of network

					if ptp.NetworkPeers[i].Endpoint == "" && len(ptp.NetworkPeers[i].KnownIPs) > 0 {
						log.Log(log.INFO, "No peers are available in local network. Switching to global IP if any")
						// If endpoint wasn't set let's test connection from outside of the LAN
						// First one should be the global IP (if DHT works correctly)
						if !ptp.TestConnection(ptp.NetworkPeers[i].KnownIPs[0]) {
							// We've failed to establish connection again. Now let's ask for a proxy
							log.Log(log.INFO, "Failed to establish connection. Requesting Control Peer from Service Discovery Peer")
							ptp.dht.RequestControlPeer(peer.ID)
							peer.PeerAddr = peer.KnownIPs[0]
							ptp.NetworkPeers[i] = peer
						} else {
							log.Log(log.INFO, "Successfully connected to a host over Internet")
							peer.Endpoint = peer.KnownIPs[0].String()
							ptp.NetworkPeers[i] = peer
						}
					}
				} else if peer.Endpoint != "" && peer.Forwarder != nil && peer.ProxyID == 0 {
					// This peer received a forwarder but it doesn't have a proxy yet
					log.Log(log.INFO, "Sending proxy request to a forwarder %s", peer.Forwarder.String())
					msg := udpcs.CreateProxyP2PMessage(-1, peer.PeerAddr.String(), 0)
					_, err := ptp.UDPSocket.SendMessage(msg, peer.Forwarder)
					if err != nil {
						log.Log(log.ERROR, "Failed to send a message to a proxy %v", err)
					}
				}
			}
		}
		if !found {
			log.Log(log.INFO, "Adding new peer. Requesting peer address")
			var newPeer NetworkPeer
			newPeer.ID = id.ID
			newPeer.Unknown = true
			ptp.NetworkPeers[newPeer.ID] = newPeer
			ptp.dht.RequestPeerIPs(id.ID)
		}
	}
	return count
}

// WriteToDevice writes data to created TUN/TAP device
func (ptp *PTPCloud) WriteToDevice(b []byte, proto uint16, truncated bool) {
	var p tuntap.Packet
	p.Protocol = int(proto)
	p.Truncated = truncated
	p.Packet = b
	if ptp.Device == nil {
		log.Log(log.ERROR, "TUN/TAP Device not initialized")
		return
	}
	err := ptp.Device.WritePacket(&p)
	if err != nil {
		log.Log(log.ERROR, "Failed to write to TUN/TAP device: %v", err)
	}
}

func GenerateMAC() (string, net.HardwareAddr) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		log.Log(log.ERROR, "Failed to generate MAC: %v", err)
		return "", nil
	}
	buf[0] |= 2
	mac := fmt.Sprintf("06:%02x:%02x:%02x:%02x:%02x", buf[1], buf[2], buf[3], buf[4], buf[5])
	hw, err := net.ParseMAC(mac)
	if err != nil {
		log.Log(log.ERROR, "Corrupted MAC address generated: %v", err)
		return "", nil
	}
	return mac, hw
}

// AddPeer adds new peer into list of network participants. If peer was added previously
// information about him will be updated. If not, new entry will be added
func (ptp *PTPCloud) AddPeer(addr *net.UDPAddr, id string, ip net.IP, mac net.HardwareAddr) {
	var found bool = false
	for i, peer := range ptp.NetworkPeers {
		if peer.ID == id {
			found = true
			peer.ID = id
			peer.PeerAddr = addr
			peer.PeerLocalIP = ip
			peer.PeerHW = mac
			peer.Unknown = false
			peer.Handshaked = true
			ptp.NetworkPeers[i] = peer
			ptp.IPIDTable[ip.String()] = id
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
		ptp.NetworkPeers[newPeer.ID] = newPeer
		ptp.IPIDTable[ip.String()] = id
	}
}

func (p *NetworkPeer) ProbeConnection() bool {
	return false
}

func (ptp *PTPCloud) ParseIntroString(intro string) (string, net.HardwareAddr, net.IP) {
	parts := strings.Split(intro, ",")
	if len(parts) != 3 {
		log.Log(log.ERROR, "Failed to parse introduction string")
		return "", nil, nil
	}
	var id string
	id = parts[0]
	// Extract MAC
	mac, err := net.ParseMAC(parts[1])
	if err != nil {
		log.Log(log.ERROR, "Failed to parse MAC address from introduction packet: %v", err)
		return "", nil, nil
	}
	// Extract IP
	ip := net.ParseIP(parts[2])
	if ip == nil {
		log.Log(log.ERROR, "Failed to parse IP address from introduction packet")
		return "", nil, nil
	}

	return id, mac, ip
}

func (ptp *PTPCloud) IsPeerUnknown(id string) bool {
	for _, peer := range ptp.NetworkPeers {
		if peer.ID == id {
			return peer.Unknown
		}
	}
	return true
}

// Handler for new messages received from P2P network
func (ptp *PTPCloud) HandleP2PMessage(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		log.Log(log.ERROR, "P2P Message Handle: %v", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := udpcs.P2PMessageFromBytes(buf)
	if des_err != nil {
		log.Log(log.ERROR, "P2PMessageFromBytes error: %v", des_err)
		return
	}
	//var msgType commons.MSG_TYPE = commons.MSG_TYPE(msg.Header.Type)
	// Decrypt message if crypter is active
	if ptp.Crypter.Active {
		var dec_err error
		msg.Data, dec_err = ptp.Crypter.Decrypt(ptp.Crypter.ActiveKey.Key, msg.Data)
		if dec_err != nil {
			log.Log(log.ERROR, "Failed to decrypt message")
		}
		msg.Data = msg.Data[:msg.Header.Length]
	}
	callback, exists := ptp.MessageHandlers[msg.Header.Type]
	if exists {
		callback(msg, src_addr)
	} else {
		log.Log(log.WARNING, "Unknown message received")
	}
}

func (ptp *PTPCloud) HandleMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {

}

func (ptp *PTPCloud) HandleNotEncryptedMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {
	log.Log(log.DEBUG, "Received P2P Message")
	ptp.WriteToDevice(msg.Data, msg.Header.NetProto, false)

}

func (ptp *PTPCloud) HandlePingMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {
	ptp.UDPSocket.SendMessage(msg, src_addr)
}

func (ptp *PTPCloud) HandleIntroMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {
	log.Log(log.DEBUG, "Introduction message received: %s", string(msg.Data))
	id, mac, ip := ptp.ParseIntroString(string(msg.Data))
	// Don't do anything if we already know everything about this peer
	if !ptp.IsPeerUnknown(id) {
		log.Log(log.DEBUG, "Skipping known peer")
		return
	}
	addr := src_addr
	// TODO: Change PeerAddr with DST addr of real peer
	ptp.AddPeer(addr, id, ip, mac)
	response := ptp.PrepareIntroductionMessage(ptp.dht.ID)
	response.Header.ProxyId = msg.Header.ProxyId
	_, err := ptp.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		log.Log(log.ERROR, "Failed to respond to introduction message: %v", err)
	}
}

func (ptp *PTPCloud) HandleProxyMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {
	// Proxy registration data
	log.Log(log.DEBUG, "Proxy confirmation received")
	if msg.Header.ProxyId < 1 {
		return
	}
	id := string(msg.Data)
	for key, peer := range ptp.NetworkPeers {
		if peer.PeerAddr.String() == id {
			peer.ProxyID = int(msg.Header.ProxyId)
			log.Log(log.DEBUG, "Settings proxy ID %d", msg.Header.ProxyId)
			ptp.NetworkPeers[key] = peer
		}
	}

}

func (ptp *PTPCloud) HandleTestMessage(msg *udpcs.P2PMessage, src_addr *net.UDPAddr) {
	response := udpcs.CreateTestP2PMessage(ptp.Crypter, "TEST", 0)
	_, err := ptp.UDPSocket.SendMessage(response, src_addr)
	if err != nil {
		log.Log(log.ERROR, "Failed to respond to test message: %v", err)
	}

}

func (ptp *PTPCloud) SendTo(dst net.HardwareAddr, msg *udpcs.P2PMessage) (int, error) {
	// TODO: Speed up this by switching to map
	for _, peer := range ptp.NetworkPeers {
		if peer.PeerHW.String() == dst.String() {
			msg.Header.ProxyId = uint16(peer.ProxyID)
			size, err := ptp.UDPSocket.SendMessage(msg, peer.PeerAddr)
			return size, err
		}
	}
	return 0, nil
}

func (ptp *PTPCloud) StopInstance() {
	// Send a packet
	ptp.dht.Stop()
	ptp.UDPSocket.Stop()
	ptp.Shutdown = true
	// Tricky part: we need to send a message to ourselves to quit blocking operation
	msg := udpcs.CreateTestP2PMessage(ptp.Crypter, "STOP", 1)
	addr, _ := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", ptp.dht.P2PPort))
	ptp.UDPSocket.SendMessage(msg, addr)
	time.Sleep(3 * time.Second)
	ptp.ReadyToStop = true
}
