package ptp

import (
	"errors"
	"net"
	"time"
)

type StateHandlerCallback func(ptpc *PTPCloud) error

type NetworkPeer struct {
	ID            string                             // ID of a peer
	Unknown       bool                               // TODO: Remove after moving to states
	Handshaked    bool                               // TODO: Remove after moving to states
	WaitingPing   bool                               // True if ping request was sent
	ProxyID       int                                // ID of the proxy
	ProxyRetries  int                                // Number of retries to reach proxy
	Forwarder     *net.UDPAddr                       // Forwarder address
	PeerAddr      *net.UDPAddr                       // Address of peer
	PeerLocalIP   net.IP                             // IP of peers interface. TODO: Rename to IP
	PeerHW        net.HardwareAddr                   // Hardware addres of peer interface. TODO: Rename to Mac
	Endpoint      string                             // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs      []*net.UDPAddr                     // List of IP addresses that accepts connection on peer
	Retries       int                                // Number of introduction retries
	Ready         bool                               // Set to true when peer is ready to communicate with p2p network
	State         PeerState                          // State of a peer
	LastContact   time.Time                          // Last ping with this peer
	StateHandlers map[PeerState]StateHandlerCallback // List of callbacks for different peer states
}

func (np *NetworkPeer) Run(ptpc *PTPCloud) {
	var initialize bool = false
	for {
		if np.State == P_DISCONNECT {
			Log(INFO, "Stopping peer %s", np.ID)
			break
		}
		if !initialize {
			np.StateHandlers = make(map[PeerState]StateHandlerCallback)
			np.StateHandlers[P_INIT] = np.StateInit
			np.StateHandlers[P_REQUESTED_IP] = np.StateRequestedIp
			np.StateHandlers[P_CONNECTING_DIRECTLY] = np.StateConnectingDirectly
		}
		callback, exists := np.StateHandlers[np.State]
		if !exists {
			Log(ERROR, "Peer %s is in unknown state")
			time.Sleep(1 * time.Second)
			continue
		}
		err := callback(ptpc)
		if err != nil {
			Log(ERROR, "Error with peer %s: %v", np.ID, err)
		}
	}
}

func (np *NetworkPeer) StateInit(ptpc *PTPCloud) error {
	// Send request about IPs of a peer
	ptpc.Dht.RequestPeerIPs(np.ID)
	np.State = P_REQUESTED_IP
	return nil
}

func (np *NetworkPeer) StateRequestedIp(ptpc *PTPCloud) error {
	// Waiting for IPs from DHT
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
		time.Sleep(100 * time.Microsecond)
	}
}

// In this state we're trying to establish direct connection.
// First we're getting list of local interfaces and see if one of
// received IPs are in the same network. If so, we will try to establish
// local connection across LAN.
// Otherwise, we will try to establish connection over WAN. If every attempt
// will fail we will switch to Proxy mode.
func (np *NetworkPeer) StateConnectingDirectly(ptpc *PTPCloud) error {
	if len(np.KnownIPs) == 0 {
		np.State = P_INIT
		return errors.New("Joined connection state without knowing any IPs")
	}
	// If forward mode was activated - skip direction connection attemps
	if ptpc.ForwardMode {
		np.RequestForwarder()
		np.State = P_WAITING_FORWARDER
		return nil
	}
	// Try to connect locally
	isLocal := np.ProbeLocalConnection(ptpc)
	if isLocal {
		np.SendHandshake(ptpc)
		np.State = P_HANDSHAKING
		return nil
	}
	// Try direct connection over the internet. If target host is not
	// behind NAT we should connect to it successfully
	// Otherwise we will failback to proxy
	addr := np.KnownIPs[0]
	conn := np.TestConnection(ptpc, addr)
	if conn {
		np.State = P_HANDSHAKING
		return nil
	} else {
		np.RequestForwarder()
		np.State = P_WAITING_FORWARDER
	}
	return nil
}

// Utilities functions

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

func (np *NetworkPeer) RequestForwarder() {

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
		if np.Endpoint != "" {
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
						return true
						Log(INFO, "Setting endpoint for %s to %s", np.ID, kip.String())
					}
				}
			}
		}
	}
	return false
}

func (np *NetworkPeer) SendHandshake(ptpc *PTPCloud) {

}
