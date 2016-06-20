package ptp

import (
	"errors"
	"fmt"
	"net"
	"time"
)

type StateHandlerCallback func(ptpc *PTPCloud) error

type NetworkPeer struct {
	ID             string                             // ID of a peer
	ProxyID        int                                // ID of the proxy
	Forwarder      *net.UDPAddr                       // Forwarder address
	PeerAddr       *net.UDPAddr                       // Address of peer
	PeerLocalIP    net.IP                             // IP of peers interface. TODO: Rename to IP
	PeerHW         net.HardwareAddr                   // Hardware addres of peer interface. TODO: Rename to Mac
	Endpoint       *net.UDPAddr                       // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs       []*net.UDPAddr                     // List of IP addresses that accepts connection on peer
	Retries        int                                // Number of introduction retries
	State          PeerState                          // State of a peer
	LastContact    time.Time                          // Last ping with this peer
	PingCount      int                                // Number of pings messages sent without response
	StateHandlers  map[PeerState]StateHandlerCallback // List of callbacks for different peer states
	ProxyBlacklist []*net.UDPAddr                     // Blacklist of proxies
	ProxyRequests  int                                // Number of requests sent
}

func (np *NetworkPeer) Run(ptpc *PTPCloud) {
	var initialize bool = false
	for {
		if np.State == P_STOP {
			Log(INFO, "Stopping peer %s", np.ID)
			break
		}
		if ptpc.Dht.ID == "" {
			time.Sleep(time.Millisecond * 500)
			continue
		}
		if !initialize {
			np.StateHandlers = make(map[PeerState]StateHandlerCallback)
			np.StateHandlers[P_INIT] = np.StateInit
			np.StateHandlers[P_REQUESTED_IP] = np.StateRequestedIp
			np.StateHandlers[P_CONNECTING_DIRECTLY] = np.StateConnectingDirectly
			np.StateHandlers[P_CONNECTED] = np.StateConnected
			np.StateHandlers[P_HANDSHAKING] = np.StateHandshaking
			np.StateHandlers[P_WAITING_FORWARDER] = np.StateWaitingForwarder
			np.StateHandlers[P_HANDSHAKING_FORWARDER] = np.StateHandshakingForwarder
			np.StateHandlers[P_HANDSHAKING_FAILED] = np.StateHandshakingFailed
			np.StateHandlers[P_DISCONNECT] = np.StateDisconnect
			np.StateHandlers[P_STOP] = np.StateStop
		}
		callback, exists := np.StateHandlers[np.State]
		if !exists {
			Log(ERROR, "Peer %s is in unknown state: %d", np.ID, int(np.State))
			time.Sleep(1 * time.Second)
			continue
		}
		err := callback(ptpc)
		if err != nil {
			Log(ERROR, "Peer %s: %v", np.ID, err)
		}
		time.Sleep(time.Millisecond * 500)
	}
}

func (np *NetworkPeer) StateInit(ptpc *PTPCloud) error {
	// Send request about IPs of a peer
	Log(INFO, "Initializing new peer: %s", np.ID)
	ptpc.Dht.RequestPeerIPs(np.ID)
	np.State = P_REQUESTED_IP
	return nil
}

func (np *NetworkPeer) StateRequestedIp(ptpc *PTPCloud) error {
	// Waiting for IPs from DHT
	Log(INFO, "Waiting network addresses for peer: %s", np.ID)
	for {
		for _, PeerInfo := range ptpc.Dht.Peers {
			if PeerInfo.ID == np.ID {
				if len(PeerInfo.Ips) >= 1 {
					np.KnownIPs = PeerInfo.Ips
					np.State = P_CONNECTING_DIRECTLY
					return nil
				}
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// In this state we're trying to establish direct connection.
// First we're getting list of local interfaces and see if one of
// received IPs are in the same network. If so, we will try to establish
// local connection across LAN.
// Otherwise, we will try to establish connection over WAN. If every attempt
// will fail we will switch to Proxy mode.
func (np *NetworkPeer) StateConnectingDirectly(ptpc *PTPCloud) error {
	Log(INFO, "Trying direct conection with peer: %s", np.ID)
	if len(np.KnownIPs) == 0 {
		np.State = P_INIT
		return errors.New("Joined connection state without knowing any IPs")
	}
	// If forward mode was activated - skip direction connection attemps
	if ptpc.ForwardMode {
		np.PeerAddr = np.KnownIPs[0]
		np.State = P_WAITING_FORWARDER
		return nil
	}
	// Try to connect locally
	isLocal := np.ProbeLocalConnection(ptpc)
	if isLocal {
		np.PeerAddr = np.Endpoint
		Log(INFO, "Connected with %s over LAN", np.ID)
		np.State = P_HANDSHAKING
		return nil
	}
	// Try direct connection over the internet. If target host is not
	// behind NAT we should connect to it successfully
	// Otherwise we will failback to proxy
	addr := np.KnownIPs[0]
	conn := np.TestConnection(ptpc, addr)
	if conn {
		np.PeerAddr = np.Endpoint
		Log(INFO, "Connected with %s over Internet", np.ID)
		np.State = P_HANDSHAKING
		return nil
	} else {
		Log(INFO, "Direct connection with %s failed", np.ID)
		np.PeerAddr = np.KnownIPs[0]
		np.State = P_WAITING_FORWARDER
	}
	return nil
}

func (np *NetworkPeer) StateConnected(ptpc *PTPCloud) error {
	if np.PingCount > 3 {
		np.State = P_DISCONNECT
		return errors.New(fmt.Sprintf("Peer %s has been timed out", np.ID))
	}
	passed := time.Since(np.LastContact)
	if passed > PEER_PING_TIMEOUT {
		msg := CreateXpeerPingMessage(PING_REQ, ptpc.HardwareAddr.String())
		ptpc.SendTo(np.PeerHW, msg)
		np.PingCount++
	}
	time.Sleep(1 * time.Second)
	return nil
}

func (np *NetworkPeer) StateHandshaking(ptpc *PTPCloud) error {
	Log(INFO, "Sending handshake to %s", np.ID)
	np.SendHandshake(ptpc)
	handshakeSentAt := time.Now()
	interval := time.Duration(time.Second * 3)
	retries := 0
	for np.State == P_HANDSHAKING {
		passed := time.Since(handshakeSentAt)
		if passed > interval {
			if retries >= 3 {
				Log(ERROR, "Failed to handshake with %s", np.ID)
				np.State = P_HANDSHAKING_FAILED
				return errors.New(fmt.Sprintf("Failed to handshake with %s", np.ID))
			} else {
				handshakeSentAt = time.Now()
				np.SendHandshake(ptpc)
				retries++
			}
		}
		time.Sleep(time.Millisecond * 200)
	}
	return nil
}

// Proxy was requested from DHT. This state waits for proxy
// address
func (np *NetworkPeer) StateWaitingForwarder(ptpc *PTPCloud) error {
	Log(INFO, "Looking in a list of cached proxies")
	for _, fwd := range ptpc.Dht.Forwarders {
		if fwd.DestinationID == np.ID {
			np.Forwarder = fwd.Addr
			np.Endpoint = fwd.Addr
			np.State = P_HANDSHAKING_FORWARDER
			Log(INFO, "Found cached forwarder")
			return nil
		}
	}
	if np.ProxyRequests >= 3 {
		Log(INFO, "We've failed to receive any proxies within this period")
		np.State = P_DISCONNECT
		ptpc.Dht.CleanForwarderBlacklist()
		return nil
	}
	Log(INFO, "Requesting proxy for %s", np.ID)
	np.RequestForwarder(ptpc)
	waitStart := time.Now()
	for np.Forwarder == nil {
		time.Sleep(time.Millisecond * 100)
		passed := time.Since(waitStart)
		if passed > WAIT_PROXY_TIMEOUT {
			np.ProxyRequests++
			return errors.New(fmt.Sprintf("No proxy were received for %s", np.ID))
		}
	}
	np.State = P_HANDSHAKING_FORWARDER
	return nil
}

func (np *NetworkPeer) StateHandshakingForwarder(ptpc *PTPCloud) error {
	if np.Forwarder == nil {
		np.State = P_WAITING_FORWARDER
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
		if passed > HANDSHAKE_PROXY_TIMEOUT {
			if attempts >= 3 {
				np.BlacklistCurrentProxy(ptpc)
				a := np.Forwarder
				np.Forwarder = nil
				np.State = P_WAITING_FORWARDER
				return errors.New(fmt.Sprintf("Failed to handshake with proxy %s [%s]", np.ID, a.String()))
			} else {
				err := np.SendProxyHandshake(ptpc)
				if err != nil {
					return err
				}
				handshakeSentAt = time.Now()
				attempts++
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	Log(INFO, "%s handshaked with proxy %s", np.ID, np.Forwarder.String())
	np.State = P_HANDSHAKING
	return nil
}

func (np *NetworkPeer) StateHandshakingFailed(ptpc *PTPCloud) error {
	if np.Forwarder != nil {
		Log(ERROR, "Failed to handshake with %s via proxy %s", np.ID, np.Forwarder.String())
		np.BlacklistCurrentProxy(ptpc)
		np.Forwarder = nil
	} else {
		Log(ERROR, "Failed to handshake directly. Switching to proxy")
	}
	np.State = P_WAITING_FORWARDER
	return nil
}

func (np *NetworkPeer) StateDisconnect(ptpc *PTPCloud) error {
	Log(INFO, "Disconnecting %s", np.ID)
	np.State = P_STOP
	// TODO: Send stop to DHT
	return nil
}

func (np *NetworkPeer) StateStop(ptpc *PTPCloud) error {
	return nil
}

// Utilities functions

func (np *NetworkPeer) BlacklistCurrentProxy(ptpc *PTPCloud) {
	Log(INFO, "%s Adding forwarder %s to a blacklist", np.ID, np.Forwarder.String())
	ptpc.Dht.BlacklistForwarder(np.Forwarder)
	exists := false
	for _, proxy := range np.ProxyBlacklist {
		if proxy.String() == np.Forwarder.String() {
			exists = true
		}
	}
	if exists {
		Log(INFO, "%s already has %s in a blacklist of proxies", np.ID, np.Forwarder.String())
	} else {
		np.ProxyBlacklist = append(np.ProxyBlacklist, np.Forwarder)
	}
}

// This method tests connection with specified endpoint
func (np *NetworkPeer) TestConnection(ptpc *PTPCloud, endpoint *net.UDPAddr) bool {
	msg := CreateTestP2PMessage(ptpc.Crypter, "TEST", 0)
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
	// TODO: Check if it was real TEST message
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

func (np *NetworkPeer) RequestForwarder(ptpc *PTPCloud) {
	ptpc.Dht.RequestControlPeer(np.ID, np.ProxyBlacklist)
}

// ProbeLocalConnection will try to connect to every known IP addr
// over local network interface
func (np *NetworkPeer) ProbeLocalConnection(ptpc *PTPCloud) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		Log(ERROR, "Failed to retrieve list of network interfaces in the system")
		return false
	}

	for _, inf := range interfaces {
		if np.Endpoint != nil {
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
				Log(DEBUG, "Probing new IP %s against network %s", kip.IP.String(), network.String())

				if network.Contains(kip.IP) {
					if np.TestConnection(ptpc, kip) {
						np.Endpoint = kip
						Log(INFO, "Setting endpoint for %s to %s", np.ID, kip.String())
						return true
					}
				}
			}
		}
	}
	return false
}

/*
 * Handshakes remote peer
 */
func (np *NetworkPeer) SendHandshake(ptpc *PTPCloud) {
	Log(DEBUG, "Preparing introduction message for %s", np.ID)
	msg := CreateIntroRequest(ptpc.Crypter, ptpc.Dht.ID)
	msg.Header.ProxyId = uint16(np.ProxyID)
	_, err := ptpc.UDPSocket.SendMessage(msg, np.Endpoint)
	if err != nil {
		Log(ERROR, "Failed to send introduction to %s", np.Endpoint.String())
	} else {
		Log(DEBUG, "Sent introduction handshake to %s [%s %d]", np.ID, np.Endpoint.String(), np.ProxyID)
	}
}

/*
 * Handshakes traffic forwarder
 */
func (np *NetworkPeer) SendProxyHandshake(ptpc *PTPCloud) error {
	Log(INFO, "Handshaking with proxy %s for %s", np.Forwarder.String(), np.ID)
	msg := CreateProxyP2PMessage(-1, np.PeerAddr.String(), uint16(ptpc.UDPSocket.GetPort()))
	_, err := ptpc.UDPSocket.SendMessage(msg, np.Forwarder)
	if err != nil {
		np.BlacklistCurrentProxy(ptpc)
		a := np.Forwarder
		np.Forwarder = nil
		np.State = P_WAITING_FORWARDER
		return errors.New(fmt.Sprintf("%s failed to send handshake to a proxy %s: %v", np.ID, a.String(), err))
	}
	return nil
}
