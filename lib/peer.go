package ptp

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"time"
)

// StateHandlerCallback is a peer method callback executed by peer state
type StateHandlerCallback func(ptpc *PeerToPeer) error

// NetworkPeer represents a peer
type NetworkPeer struct {
	ID                 string                             // ID of a peer
	ProxyID            uint16                             // ID of the proxy
	Forwarder          *net.UDPAddr                       // Forwarder address
	PeerAddr           *net.UDPAddr                       // Address of peer
	PeerLocalIP        net.IP                             // IP of peers interface. TODO: Rename to IP
	PeerHW             net.HardwareAddr                   // Hardware address of peer interface. TODO: Rename to Mac
	Endpoint           *net.UDPAddr                       // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs           []*net.UDPAddr                     // List of IP addresses that accepts connection on peer
	State              PeerState                          // State of a peer
	RemoteState        PeerState                          // State of remote peer
	LastContact        time.Time                          // Last ping with this peer
	PingCount          uint8                              // Number of pings messages sent without response
	ProxyBlacklist     []*net.UDPAddr                     // Blacklist of proxies
	ProxyRequests      uint8                              // Number of requests sent
	LastError          string                             // Test of last error occured during state execution
	ForceProxy         bool                               // Whether we are forced to use proxy or not
	TestPacketReceived bool                               // Whether or not test packet were received
	ConnectionAttempts uint8                              // How many times we tried to connect
	stateHandlers      map[PeerState]StateHandlerCallback // List of callbacks for different peer states
}

func (np *NetworkPeer) reportState(ptpc *PeerToPeer) {
	stateStr := strconv.Itoa(int(np.State))
	if stateStr == "" {
		return
	}
	//ptpc.Dht.ReportState(np.ID, stateStr)
	ptpc.Dht.sendState(np.ID, stateStr)
}

// SetState modify local state of peer
func (np *NetworkPeer) SetState(state PeerState, ptpc *PeerToPeer) {
	np.State = state
	np.reportState(ptpc)
}

// NetworkPeerState represents a state for remote peers
type NetworkPeerState struct {
	ID    string // Peer's ID
	State string // State of peer
}

// Run is main loop for a peer
func (np *NetworkPeer) Run(ptpc *PeerToPeer) {
	var initialize = false
	np.ConnectionAttempts = 0
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
			np.stateHandlers = make(map[PeerState]StateHandlerCallback)
			np.stateHandlers[PeerStateInit] = np.stateInit
			np.stateHandlers[PeerStateRequestedIP] = np.stateRequestedIP
			np.stateHandlers[PeerStateConnectingDirectlyWait] = np.stateConnectingDirectlyWait
			np.stateHandlers[PeerStateConnectingDirectly] = np.stateConnectingDirectly
			np.stateHandlers[PeerStateConnectingInternetWait] = np.stateConnectingInternetWait
			np.stateHandlers[PeerStateConnectingInternet] = np.stateConnectingInternet
			np.stateHandlers[PeerStateConnected] = np.stateConnected
			np.stateHandlers[PeerStateHandshaking] = np.stateHandshaking
			np.stateHandlers[PeerStateWaitingForwarder] = np.stateWaitingForwarder
			np.stateHandlers[PeerStateWaitingForwarderFailed] = np.stateWaitingForwarderFailed
			np.stateHandlers[PeerStateHandshakingForwarder] = np.stateHandshakingForwarder
			np.stateHandlers[PeerStateHandshakingFailed] = np.stateHandshakingFailed
			np.stateHandlers[PeerStateDisconnect] = np.stateDisconnect
			np.stateHandlers[PeerStateStop] = np.stateStop
		}
		callback, exists := np.stateHandlers[np.State]
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

// stateInit executed during peer initialization
func (np *NetworkPeer) stateInit(ptpc *PeerToPeer) error {
	// Send request about IPs of a peer
	Log(Info, "Initializing new peer: %s", np.ID)
	ptpc.Dht.sendNode(np.ID)
	np.KnownIPs = np.KnownIPs[:0]
	// Do some variables cleanup
	np.Endpoint = nil
	np.PeerAddr = nil
	np.PeerHW = nil
	np.PeerLocalIP = nil
	np.TestPacketReceived = false
	np.SetState(PeerStateRequestedIP, ptpc)
	np.ConnectionAttempts++
	if np.ConnectionAttempts > 5 {
		np.SetState(PeerStateDisconnect, ptpc)
		return fmt.Errorf("Too many unsuccessfull connection attempts")
	}
	return nil
}

// stateRequestedIP will wait for a DHT client to receive an IPs for this peer
func (np *NetworkPeer) stateRequestedIP(ptpc *PeerToPeer) error {
	// Waiting for IPs from DHT
	Log(Info, "Waiting network addresses for peer: %s", np.ID)
	requestSentAt := time.Now()
	updateInterval := time.Duration(time.Second * 5)
	attempts := 0
	for {
		if time.Since(requestSentAt) > updateInterval {
			Log(Warning, "Didn't got network addresses for peer. Requesting again")
			requestSentAt = time.Now()
			//ptpc.Dht.RequestPeerIPs(np.ID)
			err := ptpc.Dht.sendNode(np.ID)
			if err != nil {
				np.SetState(PeerStateDisconnect, ptpc)
				return fmt.Errorf("Failed to request IPs: %s", err)
			}
			attempts++
		}
		if attempts > 5 {
			np.SetState(PeerStateDisconnect, ptpc)
			break
		}
		if len(np.KnownIPs) > 0 {
			np.SetState(PeerStateConnectingDirectlyWait, ptpc)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// stateConnectingDirectlyWait - Wait for other peer to synchronize connection with us
func (np *NetworkPeer) stateConnectingDirectlyWait(ptpc *PeerToPeer) error {
	// We don't want to do this for more than 5 minutes
	Log(Info, "Waiting for other peer to start connecting directly")
	started := time.Now()
	for {
		if np.State != PeerStateConnectingDirectlyWait {
			return nil
		}
		if np.RemoteState == PeerStateConnectingDirectlyWait || np.RemoteState == PeerStateConnectingDirectly {
			Log(Info, "Second peer has joined required state")
			np.SetState(PeerStateConnectingDirectly, ptpc)
			break
		}
		time.Sleep(100 * time.Millisecond)
		passed := time.Since(started)
		if passed > time.Duration(4*time.Minute) {
			np.SetState(PeerStateConnectingDirectly, ptpc)
			return fmt.Errorf("Wait for direct connection failed: Peer doesn't responded in a timely manner")
		}
	}
	return nil
}

// stateConnectingDirectly will try to establish direct connection
// First we're getting list of local interfaces and see if one of
// received IPs are in the same network. If so, we will try to establish
// local connection across LAN.
// Otherwise, we will try to establish connection over WAN. If every attempt
// will fail we will switch to Proxy mode.
func (np *NetworkPeer) stateConnectingDirectly(ptpc *PeerToPeer) error {
	Log(Info, "Trying direct connection with peer: %s", np.ID)
	if len(np.KnownIPs) == 0 {
		np.SetState(PeerStateInit, ptpc)
		np.LastError = fmt.Sprintf("Didn't received any IP addresses")
		return errors.New("Joined connection state without knowing any IPs")
	}
	// If forward mode was activated - skip direct connection attempts
	if ptpc.ForwardMode || np.ForceProxy {
		Log(Info, "Forcing switch to proxy usage")
		np.SetPeerAddr()
		np.SetState(PeerStateWaitingForwarder, ptpc)
		return nil
	}
	// Try to connect locally
	isLocal := np.ProbeLocalConnection(ptpc)

	if isLocal {
		np.PeerAddr = np.Endpoint
		Log(Info, "Connected with %s over LAN", np.ID)
		np.SetState(PeerStateHandshaking, ptpc)
		return nil
	}
	Log(Info, "Can't connect with %s over LAN", np.ID)

	np.SetState(PeerStateConnectingInternetWait, ptpc)
	return nil
}

// stateConnectingInternetWait will wait for this peer to join the same state
func (np *NetworkPeer) stateConnectingInternetWait(ptpc *PeerToPeer) error {
	// We don't want to do this for more than 5 minutes
	Log(Info, "Waiting for other peer to start connecting over Internet")
	started := time.Now()
	for {
		if np.State != PeerStateConnectingInternetWait {
			return nil
		}
		if np.RemoteState == PeerStateConnectingInternetWait || np.RemoteState == PeerStateConnectingInternet {
			newState := "Waiting for internet connection"
			if np.RemoteState == PeerStateConnectingInternet {
				newState = "Connecting over internet"
			}
			Log(Info, "Second peer joined required state: %s", newState)
			np.SetState(PeerStateConnectingInternet, ptpc)
			break
		}
		time.Sleep(100 * time.Millisecond)
		passed := time.Since(started)
		if passed > time.Duration(4*time.Minute) {
			np.SetState(PeerStateConnectingInternet, ptpc)
			return fmt.Errorf("Wait for internet connection failed: Peer doesn't responded in a timely manner")
		}
	}
	return nil
}

// stateConnectingInternet will try to establish connection with peer over internet
// and in case if direct connection is not possible (peer is behind NAT) it
// will continue to send requests in a cycle (UDP Hole punching)
func (np *NetworkPeer) stateConnectingInternet(ptpc *PeerToPeer) error {
	// Try direct connection over the internet. If target host is not
	// behind NAT we should connect to it successfully
	// Otherwise we will failback to proxy
	addr := np.KnownIPs[0]
	np.Endpoint = addr
	Log(Info, "Attempting to connect with %s over Internet", np.ID)
	success := np.holePunch(addr, ptpc)
	if success {
		np.PeerAddr = np.Endpoint
		Log(Info, "Connected with %s over Internet", np.ID)
		np.SetState(PeerStateHandshaking, ptpc)
		return nil
	}
	np.SetPeerAddr()
	np.SetState(PeerStateWaitingForwarder, ptpc)
	return fmt.Errorf("Internet connection with %s failed", np.ID)
}

// stateConnected is executed when connection was established and peer is operating normally
func (np *NetworkPeer) stateConnected(ptpc *PeerToPeer) error {
	if np.RemoteState == PeerStateDisconnect {
		Log(Info, "Peer %s started disconnect procedure", np.ID)
		np.SetState(PeerStateDisconnect, ptpc)
		return nil
	}
	if np.RemoteState == PeerStateStop {
		Log(Info, "Peer %s has been stopped", np.ID)
		np.SetState(PeerStateDisconnect, ptpc)
		return nil
	}
	if np.RemoteState == PeerStateInit {
		Log(Info, "Remote peer %s decided to reconnect", np.ID)
		np.SetState(PeerStateInit, ptpc)
		return nil
	}
	if np.PeerHW == nil || np.PeerLocalIP == nil {
		np.SetState(PeerStateHandshaking, ptpc)
		return nil
	}
	if np.PingCount > 15 {
		np.LastError = "Disconnected by timeout"
		np.SetState(PeerStateInit, ptpc)
		np.PeerAddr = nil
		np.Endpoint = nil
		np.PingCount = 0
		return fmt.Errorf("Peer %s has been timed out", np.ID)
	}
	if np.Endpoint == nil {
		np.SetState(PeerStateInit, ptpc)
		np.PeerAddr = nil
		np.PingCount = 0
		return fmt.Errorf("Peer %s has lost endpoint", np.ID)
	}
	passed := time.Since(np.LastContact)
	if passed > PeerPingTimeout {
		np.LastError = ""
		Log(Trace, "Sending ping")
		msg := CreateXpeerPingMessage(PingReq, ptpc.Interface.Mac.String())
		ptpc.SendTo(np.PeerHW, msg)
		np.PingCount++
	}
	return nil
}

// stateHandshaking is executed when we're waiting for handshake to complete
func (np *NetworkPeer) stateHandshaking(ptpc *PeerToPeer) error {
	Log(Info, "Sending handshake to %s", np.ID)
	handshakeSentAt := time.Now()
	for np.State == PeerStateHandshaking {
		passed := time.Since(handshakeSentAt)
		if passed > time.Duration(time.Second*14) {
			np.SetState(PeerStateHandshakingFailed, ptpc)
			return fmt.Errorf("Failed to handshake with peer %s", np.ID)
		}
		np.sendHandshake(ptpc)
		time.Sleep(time.Millisecond * 500)
	}
	return nil
}

// stateHandshakingFailed is executed when we've failed to handshake a peer
func (np *NetworkPeer) stateHandshakingFailed(ptpc *PeerToPeer) error {
	if np.Forwarder != nil {
		np.LastError = "Failed to handshake with this peer over forwarder"
		Log(Error, "Failed to handshake with %s via proxy %s", np.ID, np.Forwarder.String())
		np.BlacklistCurrentProxy(ptpc)
		np.Forwarder = nil
		np.SetState(PeerStateDisconnect, ptpc)
	} else {
		np.LastError = "Failed to handshake with this peer"
		Log(Error, "Failed to handshake directly. Switching to proxy")
	}
	np.SetState(PeerStateWaitingForwarder, ptpc)
	return nil
}

// stateWaitingForwarder will wait for a proxy address
// Proxy was requested from DHT. This state waits for proxy
// address
func (np *NetworkPeer) stateWaitingForwarder(ptpc *PeerToPeer) error {
	Log(Info, "Looking in a list of cached proxies")

	for _, fwd := range ptpc.Dht.Forwarders {
		if fwd.DestinationID == np.ID {
			np.Forwarder = fwd.Addr
			np.Endpoint = fwd.Addr
			np.SetState(PeerStateHandshakingForwarder, ptpc)
			Log(Info, "Found cached forwarder")
			return nil
		}
	}

	Log(Info, "Requesting proxy for %s", np.ID)
	np.RequestForwarder(ptpc)
	waitStart := time.Now()
	for np.Forwarder == nil {
		time.Sleep(time.Millisecond * 100)
		passed := time.Since(waitStart)
		if passed > WaitProxyTimeout {
			np.SetState(PeerStateWaitingForwarderFailed, ptpc)
			np.LastError = "No forwarders received"
			return fmt.Errorf("No proxy were received for %s", np.ID)
		}
	}
	np.SetState(PeerStateHandshakingForwarder, ptpc)
	return nil
}

func (np *NetworkPeer) stateWaitingForwarderFailed(ptpc *PeerToPeer) error {
	Log(Info, "Didn't received any proxies")
	np.SetState(PeerStateInit, ptpc)
	np.ProxyBlacklist = np.ProxyBlacklist[:0]
	return nil
}

// stateHandshakingForwarder waits for handshake with a proxy to be completed
func (np *NetworkPeer) stateHandshakingForwarder(ptpc *PeerToPeer) error {
	if np.Forwarder == nil {
		np.SetState(PeerStateWaitingForwarder, ptpc)
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
				np.SetState(PeerStateWaitingForwarder, ptpc)
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
	np.SetState(PeerStateHandshaking, ptpc)
	return nil
}

// stateDisconnect is executed when we've lost or terminated connection with a peer
func (np *NetworkPeer) stateDisconnect(ptpc *PeerToPeer) error {
	Log(Info, "Disconnecting %s", np.ID)
	np.SetState(PeerStateStop, ptpc)
	// TODO: Send stop to DHT
	return nil
}

// stateStop is executed when we've terminated connection with a peer
func (np *NetworkPeer) stateStop(ptpc *PeerToPeer) error {
	return nil
}

// Utilities functions

// TestConnection method tests connection with specified endpoint
func (np *NetworkPeer) TestConnection(ptpc *PeerToPeer, endpoint *net.UDPAddr) bool {
	if endpoint == nil || ptpc == nil {
		return false
	}
	msg := CreateTestP2PMessage(ptpc.Crypter, ptpc.Dht.ID, 0)
	conn, err := net.DialUDP("udp4", nil, endpoint)
	defer conn.Close()
	if err != nil {
		Log(Debug, "%v", err)
		return false
	}
	ser := msg.Serialize()
	_, err = conn.Write(ser)
	if err != nil {
		return false
	}
	t := time.Now()
	t = t.Add(1500 * time.Millisecond)
	conn.SetReadDeadline(t)
	// TODO: Check if it was real TEST message
	for {
		var buf [4096]byte
		s, _, err := conn.ReadFromUDP(buf[0:])
		if err != nil {
			Log(Debug, "%v", err)
			break
		}
		if s > 0 {
			return true
		}
	}
	return false
}

// RequestForwarder sends a request for a proxy with DHT client
func (np *NetworkPeer) RequestForwarder(ptpc *PeerToPeer) {
	ptpc.Dht.sendProxy(np.ID)
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
		if inf.Name == ptpc.Interface.Name {
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
					result := np.holePunch(kip, ptpc)
					if result {
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

func (np *NetworkPeer) sendHandshake(ptpc *PeerToPeer) {
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
		np.SetState(PeerStateWaitingForwarder, ptpc)
		np.LastError = "Failed to send handshake to a forwarder"
		return fmt.Errorf("%s failed to send handshake to a proxy %s: %v", np.ID, a.String(), err)
	}
	return nil
}

func (np *NetworkPeer) holePunch(endpoint *net.UDPAddr, ptpc *PeerToPeer) bool {
	ptpc.HolePunching.Lock()
	defer ptpc.HolePunching.Unlock()
	Log(Info, "Starting UDP hole punching to %s", endpoint.String())
	if endpoint == nil {
		Log(Error, "Endpoint is not set")
		return false
	}
	msg := CreateTestP2PMessage(ptpc.Crypter, ptpc.Dht.ID, 0)
	packet := msg.Serialize()

	punchStarted := time.Now()
	for {
		if np.TestPacketReceived {
			np.TestPacketReceived = false
			return true
		}
		_, err := ptpc.UDPSocket.SendRawBytes(packet, endpoint)
		if err != nil {
			Log(Error, "Failed to send data: %s", err)
		}
		passed := time.Since(punchStarted)
		if passed > time.Duration(5*time.Second) {
			Log(Warning, "Stopping UDP hole punching to %s after timeout", endpoint.String())
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
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
