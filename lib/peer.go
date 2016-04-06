package ptp

import (
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
	var added bool
	if ptpc.ForwardMode {
		np.Endpoint, added = ptpc.AssignEndpoint(np)
		np.State = P_CONNECTED
		np.LastContact = time.Now()
	}
	if added {
		ptpc.NetworkPeers[np.ID] = np
	}
}
