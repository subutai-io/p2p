package ptp

import (
	"errors"
	"fmt"
	"net"
	"strings"
	"time"
)

// StateHandlerCallback is a peer method callback executed by peer state
type StateHandlerCallback func(ptpc *PeerToPeer) error

// NetworkPeer represents a peer
type NetworkPeer struct {
	ID             string                             // ID of a peer
	ProxyID        int                                // ID of the proxy
	Forwarder      *net.UDPAddr                       // Forwarder address
	PeerAddr       *net.UDPAddr                       // Address of peer
	PeerLocalIP    net.IP                             // IP of peers interface. TODO: Rename to IP
	PeerHW         net.HardwareAddr                   // Hardware address of peer interface. TODO: Rename to Mac
	Endpoint       *net.UDPAddr                       // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs       []*net.UDPAddr                     // List of IP addresses that accepts connection on peer
	Retries        int                                // Number of introduction retries
	State          PeerState                          // State of a peer
	LastContact    time.Time                          // Last ping with this peer
	PingCount      int                                // Number of pings messages sent without response
	StateHandlers  map[PeerState]StateHandlerCallback // List of callbacks for different peer states
	ProxyBlacklist []*net.UDPAddr                     // Blacklist of proxies
	ProxyRequests  int                                // Number of requests sent
	LastError      string
	ForceProxy     bool
}

// Run is main loop for a peer
func (np *NetworkPeer) Run(ptpc *PeerToPeer) {
	var initialize = false
	for {
		if np.State == PeerStateStop {
			Log(Info, "Stopping peer %s", np.ID)
			break
		}
		if ptpc.Dht.ID == "" {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if !initialize {
			np.StateHandlers = make(map[PeerState]StateHandlerCallback)
			np.StateHandlers[PeerStateInit] = np.StateInit
			np.StateHandlers[PeerStateRequestedIP] = np.StateRequestedIP
			np.StateHandlers[PeerStateConnectingDirectly] = np.StateConnectingDirectly
			np.StateHandlers[PeerStateConnected] = np.StateConnected
			np.StateHandlers[PeerStateHandshaking] = np.StateHandshaking
			np.StateHandlers[PeerStateWaitingForwarder] = np.StateWaitingForwarder
			np.StateHandlers[PeerStateHandshakingForwarder] = np.StateHandshakingForwarder
			np.StateHandlers[PeerStateHandshakingFailed] = np.StateHandshakingFailed
			np.StateHandlers[PeerStateDisconnect] = np.StateDisconnect
			np.StateHandlers[PeerStateStop] = np.StateStop
		}
		callback, exists := np.StateHandlers[np.State]
		if !exists {
			Log(Error, "Peer %s is in unknown state: %d", np.ID, int(np.State))
			time.Sleep(1 * time.Second)
			continue
		}
		err := callback(ptpc)
		if err != nil {
			Log(Warning, "Peer %s: %v", np.ID, err)
		}
		time.Sleep(time.Millisecond * 500)
	}
	Log(Info, "Peer %s has been stopped", np.ID)
}

// StateInit executed during peer initialization
func (np *NetworkPeer) StateInit(ptpc *PeerToPeer) error {
	// Send request about IPs of a peer
	Log(Info, "Initializing new peer: %s", np.ID)
	ptpc.Dht.RequestPeerIPs(np.ID)
	np.State = PeerStateRequestedIP
	return nil
}

// StateRequestedIP will wait for a DHT client to receive an IPs for this peer
func (np *NetworkPeer) StateRequestedIP(ptpc *PeerToPeer) error {
	// Waiting for IPs from DHT
	Log(Info, "Waiting network addresses for peer: %s", np.ID)
	requestSentAt := time.Now()
	updateInterval := time.Duration(time.Second * 5)
	attempts := 0
	for {
		if time.Since(requestSendAt) > updateInterval {
			Log(Warning, "Didn't got network addresses for peer. Requesting again")
			ptpc.Dht.RequestPeerIPs(np.ID)
			attempts++
		}
		if attempts > 5 {
			np.PeerState = PeerStateDisconnect
			break
		}
		for _, PeerInfo := range ptpc.Dht.Peers {
			if PeerInfo.ID == np.ID {
				if len(PeerInfo.Ips) >= 1 {
					np.KnownIPs = PeerInfo.Ips
					np.State = PeerStateConnectingDirectly
					return nil
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// SetPeerAddr will update peer address
func (np *NetworkPeer) SetPeerAddr() bool {
	if len(np.KnownIPs) == 0 {
		return false
	}
	Log(Info, "Setting peer address as %s for %s", np.KnownIPs[0].String(), np.ID)
	np.PeerAddr = np.KnownIPs[0]
	return true
}

// StateConnectingDirectly will try to establish direct connection
// First we're getting list of local interfaces and see if one of
// received IPs are in the same network. If so, we will try to establish
// local connection across LAN.
// Otherwise, we will try to establish connection over WAN. If every attempt
// will fail we will switch to Proxy mode.
func (np *NetworkPeer) StateConnectingDirectly(ptpc *PeerToPeer) error {
	Log(Info, "Trying direct connection with peer: %s", np.ID)
	if len(np.KnownIPs) == 0 {
		np.State = PeerStateInit
		np.LastError = fmt.Sprintf("Didn't received any IP addresses")
		return errors.New("Joined connection state without knowing any IPs")
	}
	// If forward mode was activated - skip direct connection attempts
	if ptpc.ForwardMode || np.ForceProxy {
		np.SetPeerAddr()
		np.State = PeerStateWaitingForwarder
		return nil
	}
	// Try to connect locally
	isLocal := np.ProbeLocalConnection(ptpc)
	if isLocal {
		np.PeerAddr = np.Endpoint
		Log(Info, "Connected with %s over LAN", np.ID)
		np.State = PeerStateHandshaking
		return nil
	} else {
		Log(Info, "Can't connect with %s over LAN", np.ID)
	}
	// Try direct connection over the internet. If target host is not
	// behind NAT we should connect to it successfully
	// Otherwise we will failback to proxy
	addr := np.KnownIPs[0]
	conn := np.TestConnection(ptpc, addr)
	if conn {
		np.PeerAddr = np.Endpoint
		Log(Info, "Connected with %s over Internet", np.ID)
		np.State = PeerStateHandshaking
		return nil
	}
	Log(Info, "Direct connection with %s failed", np.ID)
	np.SetPeerAddr()
	np.State = PeerStateWaitingForwarder
	return nil
}

// StateConnected is executed when connection was established and peer is operating normally
func (np *NetworkPeer) StateConnected(ptpc *PeerToPeer) error {
	if np.PingCount > 3 {
		np.LastError = "Disconnected by timeout"
		np.State = PeerStateInit
		np.PeerAddr = nil
		np.Endpoint = nil
		np.PingCount = 0
		return fmt.Errorf("Peer %s has been timed out", np.ID)
	}
	if np.Endpoint == nil {
		np.State = PeerStateInit
		np.PeerAddr = nil
		np.PingCount = 0
		return fmt.Errorf("Peer %s has lost endpoint", np.ID)
	}
	passed := time.Since(np.LastContact)
	if passed > PeerPingTimeout {
		np.LastError = ""
		Log(Trace, "Sending ping")
		msg := CreateXpeerPingMessage(PingReq, ptpc.HardwareAddr.String())
		ptpc.SendTo(np.PeerHW, msg)
		np.PingCount++
	}
	time.Sleep(1 * time.Second)
	return nil
}

// StateHandshaking is executed when we're waiting for handshake to complete
func (np *NetworkPeer) StateHandshaking(ptpc *PeerToPeer) error {
	Log(Info, "Sending handshake to %s", np.ID)
	np.SendHandshake(ptpc)
	handshakeSentAt := time.Now()
	interval := time.Duration(time.Second * 3)
	retries := 0
	for np.State == PeerStateHandshaking {
		passed := time.Since(handshakeSentAt)
		if passed > interval {
			if retries >= 3 {
				np.LastError = "Failed to handshake"
				Log(Error, "Failed to handshake with %s", np.ID)
				np.State = PeerStateHandshakingFailed
				return fmt.Errorf("Failed to handshake with %s", np.ID)
			}
			handshakeSentAt = time.Now()
			np.SendHandshake(ptpc)
			retries++
		}
		time.Sleep(time.Millisecond * 200)
	}
	return nil
}

// StateWaitingForwarder will wait for a proxy address
// Proxy was requested from DHT. This state waits for proxy
// address
func (np *NetworkPeer) StateWaitingForwarder(ptpc *PeerToPeer) error {
	Log(Info, "Looking in a list of cached proxies")
	for _, fwd := range ptpc.Dht.Forwarders {
		if fwd.DestinationID == np.ID {
			np.Forwarder = fwd.Addr
			np.Endpoint = fwd.Addr
			np.State = PeerStateHandshakingForwarder
			Log(Info, "Found cached forwarder")
			return nil
		}
	}
	if np.ProxyRequests >= 3 {
		np.LastError = "No more proxies for this peer"
		Log(Info, "We've failed to receive any proxies within this period")
		np.State = PeerStateInit
		ptpc.Dht.CleanForwarderBlacklist()
		np.ProxyBlacklist = np.ProxyBlacklist[:0]
		np.ProxyRequests = 0
		return nil
	}
	Log(Info, "Requesting proxy for %s", np.ID)
	np.RequestForwarder(ptpc)
	waitStart := time.Now()
	for np.Forwarder == nil {
		time.Sleep(time.Millisecond * 100)
		passed := time.Since(waitStart)
		if passed > WaitProxyTimeout {
			np.ProxyRequests++
			np.LastError = "No forwarders received"
			return fmt.Errorf("No proxy were received for %s", np.ID)
		}
	}
	np.State = PeerStateHandshakingForwarder
	return nil
}

// StateHandshakingForwarder waits for handshake with a proxy to be completed
func (np *NetworkPeer) StateHandshakingForwarder(ptpc *PeerToPeer) error {
	if np.Forwarder == nil {
		np.State = PeerStateWaitingForwarder
		return nil
	}
	np.ProxyRequests = 0
	err := np.SendProxyHandshake(ptpc)
	if err != nil {
		return err
	}
	handshakeSentAt := time.Now()
	attempts := 0
	for np.ProxyID == 0 {
		passed := time.Since(handshakeSentAt)
		if passed > HandshakeProxyTimeout {
			if attempts >= 3 {
				np.BlacklistCurrentProxy(ptpc)
				a := np.Forwarder
				np.Forwarder = nil
				np.State = PeerStateWaitingForwarder
				np.LastError = "Failed to handshake with a forwarder"
				return fmt.Errorf("Failed to handshake with proxy %s [%s]", np.ID, a.String())
			}

			err := np.SendProxyHandshake(ptpc)
			if err != nil {
				return err
			}
			handshakeSentAt = time.Now()
			attempts++
		}
		time.Sleep(time.Millisecond * 100)
	}
	Log(Info, "%s handshaked with proxy %s", np.ID, np.Forwarder.String())
	np.State = PeerStateHandshaking
	return nil
}

// StateHandshakingFailed is executed when we've failed to handshake a peer
func (np *NetworkPeer) StateHandshakingFailed(ptpc *PeerToPeer) error {
	if np.Forwarder != nil {
		np.LastError = "Failed to handshake with this peer over forwarder"
		Log(Error, "Failed to handshake with %s via proxy %s", np.ID, np.Forwarder.String())
		np.BlacklistCurrentProxy(ptpc)
		np.Forwarder = nil
		np.State = PeerStateDisconnect
	} else {
		np.LastError = "Failed to handshake with this peer"
		Log(Error, "Failed to handshake directly. Switching to proxy")
	}
	np.State = PeerStateWaitingForwarder
	return nil
}

// StateDisconnect is executed when we've lost or terminated connection with a peer
func (np *NetworkPeer) StateDisconnect(ptpc *PeerToPeer) error {
	Log(Info, "Disconnecting %s", np.ID)
	np.State = PeerStateStop
	// TODO: Send stop to DHT
	return nil
}

// StateStop is executed when we've terminated connection with a peer
func (np *NetworkPeer) StateStop(ptpc *PeerToPeer) error {
	return nil
}

// Utilities functions

// BlacklistCurrentProxy will add proxy used by this peer to a blacklist
func (np *NetworkPeer) BlacklistCurrentProxy(ptpc *PeerToPeer) {
	Log(Info, "%s Adding forwarder %s to a blacklist", np.ID, np.Forwarder.String())
	ptpc.Dht.BlacklistForwarder(np.Forwarder)
	exists := false
	for _, proxy := range np.ProxyBlacklist {
		if proxy.String() == np.Forwarder.String() {
			exists = true
		}
	}
	if exists {
		Log(Info, "%s already has %s in a blacklist of proxies", np.ID, np.Forwarder.String())
	} else {
		np.ProxyBlacklist = append(np.ProxyBlacklist, np.Forwarder)
	}
}

// TestConnection method tests connection with specified endpoint
func (np *NetworkPeer) TestConnection(ptpc *PeerToPeer, endpoint *net.UDPAddr) bool {
	msg := CreateTestP2PMessage(ptpc.Crypter, "TEST", 0)
	conn, err := net.DialUDP("udp4", nil, endpoint)
	if err != nil {
		Log(Debug, "%v", err)
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
	// TODO: Check if it was real TEST message
	for {
		var buf [4096]byte
		s, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			Log(Debug, "%v", err)
			conn.Close()
			break
		}
		if s > 0 {
			conn.Close()
			return true
			break
		}
	}
	conn.Close()
	return false
}

// RequestForwarder sends a request for a proxy with DHT client
func (np *NetworkPeer) RequestForwarder(ptpc *PeerToPeer) {
	ptpc.Dht.RequestControlPeer(np.ID, np.ProxyBlacklist)
}

// ProbeLocalConnection will try to connect to every known IP addr
// over local network interface
func (np *NetworkPeer) ProbeLocalConnection(ptpc *PeerToPeer) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		Log(Error, "Failed to retrieve list of network interfaces in the system")
		return false
	}

	for _, inf := range interfaces {
		if np.Endpoint != nil {
			Log(Info, "Endpoint already set")
			break
		}
		if inf.Name == ptpc.DeviceName {
			continue
		}
		addrs, _ := inf.Addrs()
		for _, addr := range addrs {
			netip, network, _ := net.ParseCIDR(addr.String())
			if !netip.IsGlobalUnicast() {
				continue
			}
			for _, kip := range np.KnownIPs {
				Log(Debug, "Probing new IP %s against network %s", kip.IP.String(), network.String())

				if network.Contains(kip.IP) {

					for _, i := range GlobalIPBlacklist {
						str := kip.String()
						parts := strings.Split(str, ":")
						if len(parts) > 1 && i == parts[0] {
							continue
						}
					}
					if np.TestConnection(ptpc, kip) {
						np.Endpoint = kip
						Log(Info, "Setting endpoint for %s to %s", np.ID, kip.String())
						return true
					}
				}
			}
		}
	}
	return false
}

// SendHandshake sends a handshakes to a remote peer
func (np *NetworkPeer) SendHandshake(ptpc *PeerToPeer) {
	Log(Debug, "Preparing introduction message for %s", np.ID)
	if ptpc.Dht.ID == "" {
		np.LastError = "DHT Disconnected"
		return
	}
	msg := CreateIntroRequest(ptpc.Crypter, ptpc.Dht.ID)
	msg.Header.ProxyID = uint16(np.ProxyID)
	_, err := ptpc.UDPSocket.SendMessage(msg, np.Endpoint)
	if err != nil {
		np.LastError = "Failed to send intoduction message"
		Log(Error, "Failed to send introduction to %s", np.Endpoint.String())
	} else {
		Log(Debug, "Sent introduction handshake to %s [%s %d]", np.ID, np.Endpoint.String(), np.ProxyID)
	}
}

// SendProxyHandshake sends a handshake packet to a proxy
func (np *NetworkPeer) SendProxyHandshake(ptpc *PeerToPeer) error {
	if np.PeerAddr == nil {
		for !np.SetPeerAddr() {
			time.Sleep(time.Millisecond * 100)
		}
	}
	Log(Info, "Handshaking with proxy %s for %s", np.Forwarder.String(), np.ID)
	msg := CreateProxyP2PMessage(-1, np.PeerAddr.String(), uint16(ptpc.UDPSocket.GetPort()))
	_, err := ptpc.UDPSocket.SendMessage(msg, np.Forwarder)
	if err != nil {
		np.BlacklistCurrentProxy(ptpc)
		a := np.Forwarder
		np.Forwarder = nil
		np.State = PeerStateWaitingForwarder
		np.LastError = "Failed to send handshake to a forwarder"
		return fmt.Errorf("%s failed to send handshake to a proxy %s: %v", np.ID, a.String(), err)
	}
	return nil
}
