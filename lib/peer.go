package ptp

import (
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"
)

// StateHandlerCallback is a peer method callback executed by peer state
type StateHandlerCallback func(ptpc *PeerToPeer) error

// PeerEndpoint reprsents a UDP address endpoint that instance
// may use for connection with a peer
type PeerEndpoint struct {
	Addr        *net.UDPAddr
	LastContact time.Time
}

// NetworkPeer represents a peer
type NetworkPeer struct {
	ID                 string                             // ID of a peer
	ProxyID            uint16                             // ID of the proxy
	Forwarder          *net.UDPAddr                       // Forwarder address
	PeerAddr           *net.UDPAddr                       // Address of peer
	Endpoint           *net.UDPAddr                       // Endpoint address of a peer. TODO: Make this net.UDPAddr
	KnownIPs           []*net.UDPAddr                     // List of IP addresses that accepts connection on peer
	Proxies            []*net.UDPAddr                     // List of proxies of this peer
	PeerLocalIP        net.IP                             // IP of peers interface. TODO: Rename to IP
	PeerHW             net.HardwareAddr                   // Hardware address of peer interface. TODO: Rename to Mac
	State              PeerState                          // State of a peer on our end
	RemoteState        PeerState                          // State of remote peer
	LastContact        time.Time                          // Last ping with this peer
	PingCount          uint8                              // Number of pings messages sent without response
	LastError          string                             // Test of last error occured during state execution
	ForceProxy         bool                               // Whether we are forced to use proxy or not
	TestPacketReceived bool                               // Whether or not test packet were received
	ConnectionAttempts uint8                              // How many times we tried to connect
	stateHandlers      map[PeerState]StateHandlerCallback // List of callbacks for different peer states
	IsUsingTURN        bool                               // Whether or not we are currently connected via TURN
	Running            bool                               // Whether peer is running or not
	Endpoints          []PeerEndpoint                     // List of active endpoints
	EndpointsLock      sync.RWMutex                       // Mutex for endpoints operations
	punchingInProgress bool                               // Whether or not UDP hole punching is running
	LastFind           time.Time                          // Moment when we got this peer from DHT
}

func (np *NetworkPeer) reportState(ptpc *PeerToPeer) {
	stateStr := strconv.Itoa(int(np.State))
	if stateStr == "" {
		return
	}
	ptpc.Dht.sendState(np.ID, stateStr)
}

// SetState modify local state of peer
func (np *NetworkPeer) SetState(state PeerState, ptpc *PeerToPeer) {
	np.State = state
	if state != PeerStateRouting && state != PeerStateConnected {
		np.reportState(ptpc)
	}
}

// NetworkPeerState represents a state for remote peers
type NetworkPeerState struct {
	ID    string // Peer's ID
	State string // State of peer
}

// Run is main loop for a peer
func (np *NetworkPeer) Run(ptpc *PeerToPeer) {
	np.Running = true
	np.ConnectionAttempts = 0

	np.stateHandlers = make(map[PeerState]StateHandlerCallback)
	np.stateHandlers[PeerStateInit] = np.stateInit
	np.stateHandlers[PeerStateRequestedIP] = np.stateRequestedIP
	np.stateHandlers[PeerStateConnecting] = np.stateConnecting
	np.stateHandlers[PeerStateConnected] = np.stateConnected
	np.stateHandlers[PeerStateDisconnect] = np.stateDisconnect
	np.stateHandlers[PeerStateStop] = np.stateStop
	np.stateHandlers[PeerStateRequestingProxy] = np.stateRequestingProxy
	np.stateHandlers[PeerStateWaitingForProxy] = np.stateWaitingForProxy
	np.stateHandlers[PeerStateWaitingToConnect] = np.stateWaitingToConnect
	np.stateHandlers[PeerStateRouting] = np.stateRouting
	np.stateHandlers[PeerStateCooldown] = np.stateCooldown

	for {
		if np.State == PeerStateStop {
			Log(Info, "Stopping peer %s", np.ID)
			break
		}
		if ptpc.Dht.ID == "" {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		if np.ConnectionAttempts > 1 && np.ConnectionAttempts%10 == 0 {
			np.SetState(PeerStateCooldown, ptpc)
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

// State: Peer Initialization
// Initialize variables
// Automatically switch to PeerStateRequestedIP or PeerStateDisconnect if
// too many connection attempts were failed
func (np *NetworkPeer) stateInit(ptpc *PeerToPeer) error {
	// Send request about IPs of a peer
	Log(Info, "Initializing new peer: %s", np.ID)
	ptpc.Dht.sendNode(np.ID)
	// np.KnownIPs = np.KnownIPs[:0]
	// Do some variables cleanup
	np.Endpoint = nil
	np.PeerAddr = nil
	np.PeerHW = nil
	np.PeerLocalIP = nil
	np.TestPacketReceived = false
	np.IsUsingTURN = false
	if len(np.KnownIPs) == 0 {
		np.SetState(PeerStateRequestedIP, ptpc)
	} else if len(np.Proxies) == 0 {
		np.SetState(PeerStateRequestingProxy, ptpc)
	} else {
		np.SetState(PeerStateWaitingToConnect, ptpc)
	}
	// np.ConnectionAttempts++
	// if np.ConnectionAttempts > 5 {
	// 	np.SetState(PeerStateDisconnect, ptpc)
	// 	return fmt.Errorf("Too many unsuccessfull connection attempts")
	// }
	return nil
}

// stateRequestedIP will wait for a DHT client to receive an IPs for this peer
// State: Requested peer IP
// Send `node` request and wait for known IP addresses of the peer from DHT
// If peer doesn't receive endpoints in the timely manner method will switch to
// PeerStateDisconnect. On success it will switch to PeerStateConnecting
func (np *NetworkPeer) stateRequestedIP(ptpc *PeerToPeer) error {
	Log(Info, "Waiting network addresses for peer: %s", np.ID)
	requestSentAt := time.Now()
	updateInterval := time.Duration(time.Millisecond * 1000)
	attempts := 0
	for {
		if time.Since(requestSentAt) > updateInterval {
			Log(Warning, "Didn't got network addresses for peer. Requesting again")
			requestSentAt = time.Now()
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
			np.SetState(PeerStateRequestingProxy, ptpc)
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
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
		Log(Info, "Missing system information for this peer")
		np.SetState(PeerStateDisconnect, ptpc)
		return nil
	}

	if time.Since(np.LastContact) > time.Duration(time.Millisecond*3000) {
		np.LastContact = time.Now()
		np.EndpointsLock.RLock()
		for _, ep := range np.Endpoints {
			payload := append([]byte("req"), []byte(ep.Addr.String())...)
			msg, err := ptpc.CreateMessage(MsgTypeXpeerPing, payload, 0, true)
			if err != nil {
				continue
			}
			ptpc.UDPSocket.SendMessage(msg, ep.Addr)
		}
		np.EndpointsLock.RUnlock()
	}

	np.SetState(PeerStateRouting, ptpc)
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
	Log(Info, "Peer %s has been stopped", np.ID)
	return nil
}

// Utilities functions

// RequestForwarder sends a request for a proxy with DHT client
func (np *NetworkPeer) RequestForwarder(ptpc *PeerToPeer) {
	ptpc.Dht.sendRequestProxy(np.ID)
}

// Run hope punching in a separate goroutine and switch to
// Routing/Connected mode
func (np *NetworkPeer) stateConnecting(ptpc *PeerToPeer) error {
	Log(Debug, "Connecting to %s", np.ID)
	go func() {
		if np.punchingInProgress {
			return
		}
		Log(Debug, "Hole punching %s", np.ID)

		np.punchingInProgress = true
		round := 0
		for round < 10 {
			for _, ep := range np.KnownIPs {
				alreadyConnected := false
				for _, nep := range np.Endpoints {
					if nep.Addr.String() == ep.String() {
						alreadyConnected = true
					}
				}
				if alreadyConnected {
					continue
				}
				payload := []byte(ptpc.Dht.ID + ep.String())
				msg, err := ptpc.CreateMessage(MsgTypeIntroReq, payload, 0, true)
				if err != nil {
					Log(Error, "Couldn't create an intro message: %s", err)
					continue
				}
				_, err = ptpc.UDPSocket.SendMessage(msg, ep)
				if err != nil {
					Log(Error, "Failed to send message to %s: %s", ep.String(), err)
					continue
				}
				time.Sleep(time.Millisecond * 50)
			}
			time.Sleep(time.Millisecond * 50)
			round++
		}
		np.punchingInProgress = false
	}()
	np.SetState(PeerStateRouting, ptpc)
	return nil
}

func (np *NetworkPeer) stateRequestingProxy(ptpc *PeerToPeer) error {
	ptpc.Dht.sendRequestProxy(np.ID)
	np.SetState(PeerStateWaitingForProxy, ptpc)
	return nil
}

func (np *NetworkPeer) stateWaitingForProxy(ptpc *PeerToPeer) error {
	started := time.Now()
	for time.Since(started) < time.Duration(time.Millisecond*4000) {
		time.Sleep(time.Millisecond * 100)
	}
	np.SetState(PeerStateConnecting, ptpc)
	return nil
}

func (np *NetworkPeer) stateWaitingToConnect(ptpc *PeerToPeer) error {
	Log(Debug, "Waiting for other peer to join connection state")
	started := time.Now()
	for {
		if np.State != PeerStateWaitingToConnect {
			return nil
		}
		if np.RemoteState == PeerStateWaitingToConnect || np.RemoteState == PeerStateConnecting {
			Log(Debug, "Peer %s joined connection state", np.ID)
			np.SetState(PeerStateConnecting, ptpc)
			break
		}
		time.Sleep(100 * time.Millisecond)
		passed := time.Since(started)
		if passed > time.Duration(60000*time.Millisecond) {
			np.LastError = "Peer state desync"
			np.SetState(PeerStateDisconnect, ptpc)
			return fmt.Errorf("Wait for connection failed: Peer doesn't responded in a timely manner")
		}
	}
	return nil
}

func (np *NetworkPeer) stateRouting(ptpc *PeerToPeer) error {
	for len(np.Endpoints) == 0 && np.punchingInProgress {
		time.Sleep(time.Millisecond * 100)
	}
	locals := []PeerEndpoint{}
	internet := []PeerEndpoint{}
	proxies := []PeerEndpoint{}
	np.EndpointsLock.RLock()
	for _, ep := range np.Endpoints {
		if time.Since(ep.LastContact) > time.Duration(time.Second*10) {
			continue
		}
		// Check if it's proxy
		isProxy := false
		for _, proxy := range np.Proxies {
			if proxy.String() == ep.Addr.String() {
				isProxy = true
				break
			}
		}
		if isProxy {
			proxies = append(proxies, ep)
			continue
		}
		// Check if it's LAN
		rc, err := isPrivateIP(ep.Addr.IP)
		if err != nil {
			continue
		}
		if rc {
			locals = append(locals, ep)
			continue
		}
		// Add as Internet Endpoint
		internet = append(internet, ep)
	}
	np.EndpointsLock.RUnlock()

	np.EndpointsLock.Lock()
	np.Endpoints = np.Endpoints[:0]
	np.Endpoints = append(np.Endpoints, locals...)
	np.Endpoints = append(np.Endpoints, internet...)
	np.Endpoints = append(np.Endpoints, proxies...)
	np.EndpointsLock.Unlock()

	if len(np.Endpoints) > 0 {
		np.Endpoint = np.Endpoints[0].Addr
		np.SetState(PeerStateConnected, ptpc)
		np.ConnectionAttempts = 0
	} else {
		np.ConnectionAttempts++
		np.LastError = "No more endpoints"
		if time.Since(np.LastFind) > time.Duration(time.Second*90) {
			Log(Debug, "No endpoints and no updates from DHT")
			np.SetState(PeerStateDisconnect, ptpc)
			return nil
		}
		if len(np.KnownIPs) > 0 && len(np.Proxies) > 0 {
			Log(Debug, "We have IPs and Proxies. Syncing states")
			np.SetState(PeerStateWaitingToConnect, ptpc)
			return nil
		} else if len(np.KnownIPs) == 0 {
			Log(Debug, "Don't know any endpoints. Requesting")
			np.SetState(PeerStateRequestedIP, ptpc)
			return nil
		} else if len(np.Proxies) == 0 {
			Log(Debug, "Don't know any proxies. Requesting")
			np.SetState(PeerStateRequestingProxy, ptpc)
			return nil
		}
		//np.SetState(PeerStateDisconnect, ptpc)
	}
	return nil
}

func (np *NetworkPeer) addEndpoint(addr *net.UDPAddr) error {
	np.EndpointsLock.Lock()
	defer np.EndpointsLock.Unlock()
	for _, ep := range np.Endpoints {
		if ep.Addr.String() == addr.String() {
			return fmt.Errorf("Endpoint already exists")
		}
	}
	np.Endpoints = append(np.Endpoints, PeerEndpoint{Addr: addr, LastContact: time.Now()})
	return nil
}

func (np *NetworkPeer) stateCooldown(ptpc *PeerToPeer) error {
	Log(Debug, "Peer %s in cooldown", np.ID)
	started := time.Now()
	for time.Since(started) < time.Duration(time.Second*30) {
		time.Sleep(time.Millisecond * 100)
	}
	np.ConnectionAttempts++
	np.SetState(PeerStateRouting, ptpc)
	return nil
}
