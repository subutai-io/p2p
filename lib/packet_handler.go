package ptp

import (
	"net"
	"time"
)

// Handlers for P2P packets received from other network peers or TURN servers

// MessageHandler is a messages callback
type MessageHandler func(message *P2PMessage, srcAddr *net.UDPAddr)

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
	if msg == nil {
		Log(Error, "Received broken message")
		return
	}
	// Decrypt message if crypter is active
	if p.Crypter.Active && (msg.Header.Type == MsgTypeIntro || msg.Header.Type == MsgTypeNenc || msg.Header.Type == MsgTypeIntroReq || msg.Header.Type == MsgTypeTest || msg.Header.Type == MsgTypeXpeerPing) {
		var decErr error
		msg.Data, decErr = p.Crypter.decrypt(p.Crypter.ActiveKey.Key, msg.Data)
		if decErr != nil {
			Log(Error, "Failed to decrypt message: %s", decErr)
			return
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
	Log(Trace, "Data: %s, From: %s", msg.Data, srcAddr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
}

// HandlePingMessage is a PING message from a proxy handler
func (p *PeerToPeer) HandlePingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		if p.ProxyManager.touch(srcAddr.String()) {
			p.UDPSocket.SendMessage(msg, srcAddr)
		}
		return
	}
	port := addr.Port
	if p.UDPSocket.remotePort == 0 {
		p.UDPSocket.remotePort = port
	} else {
		if port != p.UDPSocket.GetPort() && port != p.UDPSocket.remotePort && port != 0 {
			Log(Debug, "Port translation detected %d -> %d", p.UDPSocket.GetPort(), port)
			p.UDPSocket.remotePort = port
		}
	}
}

// HandleXpeerPingMessage receives a cross-peer ping message
func (p *PeerToPeer) HandleXpeerPingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {

	if msg == nil {
		return
	}

	if srcAddr == nil {
		return
	}

	if len(msg.Data) < 1 {
		return
	}
	query := string(msg.Data)[:1]
	if query == "q" {
		if len(msg.Data) < 37 {
			return
		}
		id := string(msg.Data)[1:37]
		endpoint := string(msg.Data)[37:]
		response := append([]byte("r"), []byte(endpoint)...)

		msg, err := p.CreateMessage(MsgTypeXpeerPing, response, 0, true)
		if err != nil {
			Log(Debug, "Failed to create ping response: %s", err)
			return
		}
		// Look if we really know this peer

		for _, peer := range p.Peers.Get() {
			if peer.ID == id {
				for _, ep := range peer.KnownIPs {
					if ep.String() == srcAddr.String() {
						p.UDPSocket.SendMessage(msg, ep)
						peer.BumpEndpoint(ep.String())
						return
					}
				}
				// It is possible that we received ping over proxy. In this case
				// origin address will not match any of the endpoints. Therefore
				// we are going to iterate over registered proxies
				overProxy := false
				for _, proxy := range p.ProxyManager.get() {
					if proxy.Endpoint.String() == string(endpoint) {
						overProxy = true
						break
					}
				}
				if overProxy && peer.State == PeerStateConnected && peer.RemoteState == PeerStateConnected {
					p.UDPSocket.SendMessage(msg, peer.Endpoint)
					return
				}
			}
		}
		Log(Debug, "Received ping from unknown endpoint: %s [%s ID: %s]", srcAddr.String(), endpoint, id)
	} else if query == "r" {
		endpoint := msg.Data[1:]
		for _, peer := range p.Peers.Get() {
			if peer == nil {
				continue
			}
			for i, ep := range peer.EndpointsHeap {
				if ep.Addr.String() == string(endpoint) {
					peer.EndpointsHeap[i].LastContact = time.Now()
					return
				}
			}
		}
	} else {
		Log(Trace, "Wrong xpeer ping message")
	}
}

// HandleIntroMessage receives an introduction string from another peer during handshake
func (p *PeerToPeer) HandleIntroMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Debug, "Introduction string from %s", srcAddr)
	hs, err := p.ParseIntroString(string(msg.Data))
	if err != nil {
		Log(Debug, "Failed to parse handshake response: %s", err)
		return
	}
	if len(hs.ID) != 36 {
		Log(Debug, "Received wrong ID in introduction message: %s", hs.ID)
		return
	}
	peer := p.Peers.GetPeer(hs.ID)
	if peer == nil {
		Log(Trace, "Unknown peer in handshke response")
		return
	}

	if hs.HardwareAddr == nil {
		Log(Debug, "Received empty MAC address. Skipping")
		return
	}
	if hs.IP == nil {
		Log(Debug, "No IP received. Skipping")
		return
	}
	peer.PeerHW = hs.HardwareAddr
	peer.PeerLocalIP = hs.IP
	peer.LastContact = time.Now()
	peer.addEndpoint(hs.Endpoint)
	for _, np := range p.Peers.Get() {
		if np.ID == peer.ID {
			continue
		}
		if np.PeerHW.String() == peer.PeerHW.String() {
			Log(Warning, "%s: Duplicate MAC has been detected on peer %s. Disconnecting it", peer.ID, np.ID)
			np.SetState(PeerStateDisconnect, p)
			continue
		}
		for _, ep := range np.EndpointsHeap {
			if ep.Addr.String() == hs.Endpoint.String() {
				Log(Warning, "%s: Endpoint %s was used by another peer %s. Disconnecting it.", peer.ID, ep.Addr.String(), np.ID)
				np.SetState(PeerStateDisconnect, p)
				break
			}
		}
	}
	//peer.Endpoints = append(peer.Endpoints, PeerEndpoint{Addr: hs.Endpoint, LastContact: time.Now()})
	// peer.SetState(PeerStateConnected, p)
	p.Peers.Update(hs.ID, peer)
	Log(Debug, "Connection with peer %s has been established over %s", hs.ID, hs.Endpoint.String())
}

// HandleIntroRequestMessage is a handshake request from another peer
// First 36 bytes is an ID of original sender, data after byte 36 is an
// endpoint on which sender was trying to communicate with this peer.
// We need to send this data back to him, so he knows which endpoint
// replied
func (p *PeerToPeer) HandleIntroRequestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	id := string(msg.Data[0:36])
	peer := p.Peers.GetPeer(id)
	if peer == nil {
		Log(Trace, "Introduction request came from unknown peer: %s -> %s [%s]", id, msg.Data[36:], srcAddr.String())
		//p.Dht.sendFind()
		return
	}
	response := p.PrepareIntroductionMessage(p.Dht.ID, string(msg.Data[36:]))
	eps := []*net.UDPAddr{}
	eps = append(eps, peer.KnownIPs...)
	eps = append(eps, peer.Proxies...)
	Log(Debug, "Sending handshake response")

	srcFound := false
	for _, ep := range eps {
		if ep.String() == srcAddr.String() {
			srcFound = true
		}
	}
	if !srcFound {
		eps = append(eps, srcAddr)
	}

	for _, ep := range eps {
		time.Sleep(time.Millisecond * 10)
		_, err := p.UDPSocket.SendMessage(response, ep)
		if err != nil {
			Log(Error, "Failed to respond to introduction request: %v", err)
		}
	}
}

// HandleProxyMessage receives a control packet from proxy
// Proxy packets comes in format of UDP connection address
func (p *PeerToPeer) HandleProxyMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Debug, "New proxy message from %s", srcAddr)

	ep, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		Log(Error, "Failed to resolve proxy address: %s", err)
		return
	}
	rc := p.ProxyManager.activate(srcAddr.String(), ep)
	if rc {
		Log(Debug, "This peer is now available over %s", ep.String())
	}
}

// HandleBadTun notified peer about proxy being malfunction
// This method is not used in currenct scheme
// TODO: Consider to remove
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) {

}

// HandleLatency will handle latency respones from another peer/proxy
func (p *PeerToPeer) HandleLatency(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Trace, "Latency response from %s", srcAddr.String())

	ts := time.Time{}
	err := ts.UnmarshalBinary(msg.Data)
	if err != nil {
		Log(Error, "Failed to unmarshal latency packet from %s: %s", srcAddr.String(), err.Error())
		return
	}

	latency := time.Since(ts)

	// Lookup where this packet comes from
	if p.ProxyManager.setLatency(latency, srcAddr) == nil {
		return
	}

}
