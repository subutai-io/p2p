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
	//var msgType MSG_TYPE = MSG_TYPE(msg.Header.Type)
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
	// p.proxyLock.Lock()
	// defer p.proxyLock.Unlock()
	if err != nil {
		if p.ProxyManager.touch(srcAddr.String()) {
			p.UDPSocket.SendMessage(msg, srcAddr)
		}
		// for i, proxy := range p.Proxies {
		// 	if proxy == nil {
		// 		continue
		// 	}
		// 	if p.Proxies[i] != nil && proxy.Addr != nil && srcAddr != nil && proxy.Addr.String() == srcAddr.String() {
		// 		p.Proxies[i].LastUpdate = time.Now()
		// 		p.UDPSocket.SendMessage(msg, srcAddr)
		// 		break
		// 	}
		// }
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

	query := string(msg.Data)[:1]
	if query == "q" {
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
				for i, ep := range peer.Endpoints {
					if ep.Addr.String() == srcAddr.String() {
						peer.Endpoints[i].LastContact = time.Now()
						p.UDPSocket.SendMessage(msg, ep.Addr)
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
			for i, ep := range peer.Endpoints {
				if ep.Addr.String() == string(endpoint) {
					peer.Endpoints[i].LastContact = time.Now()
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
	//peer.Endpoints = append(peer.Endpoints, PeerEndpoint{Addr: hs.Endpoint, LastContact: time.Now()})
	// peer.SetState(PeerStateConnected, p)
	p.Peers.Update(hs.ID, peer)
	Log(Info, "Connection with peer %s has been established over %s", hs.ID, hs.Endpoint.String())
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

	// p.proxyLock.Lock()
	// defer p.proxyLock.Unlock()
	// for i, proxy := range p.Proxies {
	// 	if proxy.Addr.String() == srcAddr.String() && proxy.Status == proxyConnecting {
	// 		p.Proxies[i].Status = proxyActive
	// 		addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	// 		if err != nil {
	// 			Log(Error, "Failed to resolve proxy address: %s", err)
	// 			return
	// 		}
	// 		Log(Debug, "This peer is now available over %s", addr.String())
	// 		p.Dht.sendReportProxy(addr)
	// 		break
	// 	}
	// }
}

// HandleBadTun notified peer about proxy being malfunction
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) {
	// peers := p.Peers.Get()
	// for i, peer := range peers {
	// 	if peer.ProxyID == msg.Header.ProxyID && peer.Endpoint.String() == srcAddr.String() {
	// 		Log(Debug, "Cleaning bad tunnel %d from %s", msg.Header.ProxyID, srcAddr.String())
	// 		peer.ProxyID = 0
	// 		peer.Endpoint = nil
	// 		peer.Forwarder = nil
	// 		peer.PeerAddr = nil
	// 		peer.SetState(PeerStateInit, p)
	// 		p.Peers.Update(i, peer)
	// 	}
	// }
}

// HandleTestMessage responses with a test message when another peer trying to
// establish direct connection
// func (p *PeerToPeer) HandleTestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
// 	if len(p.Dht.ID) != 36 {
// 		return
// 	}

// 	if len(msg.Data) != 36 {
// 		Log(Error, "Malformed data received during test: %s [L: %d]", string(msg.Data), msg.Header.Length)
// 		return
// 	}

// 	// See if we have peer with this ID
// 	id := string(msg.Data[0:36])
// 	if len(id) != 36 {
// 		Log(Error, "Wrong ID during test message")
// 		return
// 	}

// 	peer := p.Peers.GetPeer(id)
// 	if peer != nil {
// 		if peer.State == PeerStateConnectingDirectly || peer.State == PeerStateConnectingInternet {
// 			peer.TestPacketReceived = true
// 			p.Peers.Update(id, peer)
// 			response := CreateTestP2PMessage(p.Crypter, p.Dht.ID, 0)
// 			_, err := p.UDPSocket.SendMessage(response, srcAddr)
// 			if err != nil {
// 				Log(Error, "Failed to respond to test message: %v", err)
// 			}
// 		} else if peer.State == PeerStateConnected && peer.IsUsingTURN {
// 			Log(Info, "Received test message from peer which was previously connected over TURN")
// 			if len(peer.KnownIPs) == 0 {
// 				return
// 			}
// 			peer.Endpoint = peer.KnownIPs[0]
// 			peer.IsUsingTURN = false
// 			p.Peers.Update(peer.ID, peer)
// 			Log(Info, "Peer %s switched to direct UDP connection %s", peer.ID, peer.Endpoint.String())
// 		} else {
// 			Log(Info, "Received test message for peer in unsupported state")
// 		}
// 	}
// }
