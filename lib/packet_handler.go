package ptp

import (
	"net"
	"time"
)

// Handlers for P2P packets received from other network peers or TURN servers

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
			Log(Error, "Failed to decrypt message")
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
	Log(Trace, "Data: %s, Proto: %d, From: %s", msg.Data, msg.Header.NetProto, srcAddr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
}

// HandlePingMessage is a PING message from a proxy handler
func (p *PeerToPeer) HandlePingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		p.UDPSocket.SendMessage(msg, srcAddr)
		for i, proxy := range p.Proxies {
			if proxy == nil {
				continue
			}
			if p.Proxies[i] != nil && proxy.Addr != nil && srcAddr != nil && proxy.Addr.String() == srcAddr.String() {
				p.Proxies[i].LastUpdate = time.Now()
				break
			}
		}
		return
	}
	port := addr.Port
	if p.UDPSocket.remotePort == 0 {
		p.UDPSocket.remotePort = port
	} else {
		if port != p.UDPSocket.GetPort() && port != p.UDPSocket.remotePort {
			Log(Debug, "Port translation detected %d -> %d", p.UDPSocket.GetPort(), port)
			p.UDPSocket.remotePort = port
		}
	}
}

// HandleXpeerPingMessage receives a cross-peer ping message
func (p *PeerToPeer) HandleXpeerPingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	pt := PingType(msg.Header.NetProto)
	if pt == PingReq {
		Log(Debug, "Ping request received: %s. Responding with %s", string(msg.Data), p.Interface.GetHardwareAddress().String())
		// Send a PING response
		r := CreateXpeerPingMessage(p.Crypter, PingResp, p.Interface.GetHardwareAddress().String())
		addr, err := net.ParseMAC(string(msg.Data))
		if err != nil {
			Log(Error, "Failed to parse MAC address in crosspeer ping message")
		} else {
			p.SendTo(addr, r)
			Log(Debug, "Sending crosspeer PING response to %s", addr.String())
		}
	} else {
		Log(Debug, "Ping response received")
		// Handle PING response
		peers := p.Peers.Get()
		for i, peer := range peers {
			if peer.PeerHW != nil && peer.PeerHW.String() == string(msg.Data) {
				peer.PingCount = 0
				peer.LastContact = time.Now()
				p.Peers.Update(i, peer)
				break
			}
		}
	}
}

// HandleIntroMessage receives an introduction string from another peer during handshake
func (p *PeerToPeer) HandleIntroMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Info, "Introduction string from %s[%d]", srcAddr, msg.Header.ProxyID)
	id, mac, ip := p.ParseIntroString(string(msg.Data))
	if len(id) != 36 {
		Log(Debug, "Received wrong ID in introduction message: %s", id)
		return
	}
	peer := p.Peers.GetPeer(id)
	// Do nothing when handshaking already done
	if peer.State != PeerStateHandshaking && peer.State != PeerStateHandshakingForwarder {
		return
	}
	if peer == nil {
		Log(Debug, "Received introduction confirmation from unknown peer: %s", id)
		//p.Dht.sendFind()
		return
	}

	if mac == nil {
		Log(Error, "Received empty MAC address. Skipping")
		return
	}
	if ip == nil {
		Log(Error, "No IP received. Skipping")
		return
	}
	peer.PeerHW = mac
	peer.PeerLocalIP = ip
	peer.LastContact = time.Now()
	peer.SetState(PeerStateConnected, p)
	p.Peers.Update(id, peer)
	Log(Info, "Connection with peer %s has been established", id)
}

// HandleIntroRequestMessage is a handshake request from another peer
func (p *PeerToPeer) HandleIntroRequestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	id := string(msg.Data)
	if len(id) != 36 {
		Log(Debug, "Introduction request with malformed ID [%s] from %s", id, srcAddr.String())
		return
	}
	peer := p.Peers.GetPeer(id)
	if peer == nil {
		Log(Debug, "Introduction request came from unknown peer: %s [%s]", id, srcAddr.String())
		//p.Dht.sendFind()
		return
	}
	proxy := false
	if msg.Header.ProxyID > 0 {
		proxy = true
		Log(Info, "Received introduction request via proxy")
		if len(peer.Proxies) == 0 {
			Log(Warning, "Peer %s has no proxies attached", id)
			p.Dht.sendRequestProxy(id)
			return
		}
	} else {
		Log(Info, "Received introduction request directly")
	}

	response := p.PrepareIntroductionMessage(p.Dht.ID)
	if proxy {
		response.Header.ProxyID = 1
		for _, peerProxy := range peer.Proxies {
			Log(Info, "Sending handshake response over proxy %s", peerProxy.String())
			_, err := p.UDPSocket.SendMessage(response, peerProxy)
			if err != nil {
				Log(Error, "Failed to respond to introduction request over proxy: %v", err)
			}
		}
		return
	}
	Log(Info, "Sending handshake response")
	_, err := p.UDPSocket.SendMessage(response, srcAddr)
	if err != nil {
		Log(Error, "Failed to respond to introduction request: %v", err)
	}
}

// HandleProxyMessage receives a control packet from proxy
// Proxy packets comes in format of UDP connection address
func (p *PeerToPeer) HandleProxyMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	Log(Debug, "New proxy message from %s", srcAddr)
	for i, proxy := range p.Proxies {
		if proxy.Addr.String() == srcAddr.String() && proxy.Status == proxyConnecting {
			p.Proxies[i].Status = proxyActive
			addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
			if err != nil {
				Log(Error, "Failed to resolve proxy addr: %s", err)
				return
			}
			Log(Info, "This peer is now available over %s", addr.String())
			p.Dht.sendReportProxy(addr)
		}
	}
}

// HandleBadTun notified peer about proxy being malfunction
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) {
	peers := p.Peers.Get()
	for i, peer := range peers {
		if peer.ProxyID == msg.Header.ProxyID && peer.Endpoint.String() == srcAddr.String() {
			Log(Debug, "Cleaning bad tunnel %d from %s", msg.Header.ProxyID, srcAddr.String())
			peer.ProxyID = 0
			peer.Endpoint = nil
			peer.Forwarder = nil
			peer.PeerAddr = nil
			peer.SetState(PeerStateInit, p)
			p.Peers.Update(i, peer)
		}
	}
}

// HandleTestMessage responses with a test message when another peer trying to
// establish direct connection
func (p *PeerToPeer) HandleTestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) {
	if len(p.Dht.ID) != 36 {
		return
	}

	if len(msg.Data) != 36 {
		Log(Error, "Malformed data received during test: %s [L: %d]", string(msg.Data), msg.Header.Length)
		return
	}

	// See if we have peer with this ID
	id := string(msg.Data[0:36])
	if len(id) != 36 {
		Log(Error, "Wrong ID during test message")
		return
	}

	peer := p.Peers.GetPeer(id)
	if peer != nil {
		if peer.State == PeerStateConnectingDirectly || peer.State == PeerStateConnectingInternet {
			peer.TestPacketReceived = true
			p.Peers.Update(id, peer)
			response := CreateTestP2PMessage(p.Crypter, p.Dht.ID, 0)
			_, err := p.UDPSocket.SendMessage(response, srcAddr)
			if err != nil {
				Log(Error, "Failed to respond to test message: %v", err)
			}
		} else if peer.State == PeerStateConnected && peer.IsUsingTURN {
			Log(Info, "Received test message from peer which was previously connected over TURN")
			if len(peer.KnownIPs) == 0 {
				return
			}
			peer.Endpoint = peer.KnownIPs[0]
			peer.IsUsingTURN = false
			p.Peers.Update(peer.ID, peer)
			Log(Info, "Peer %s switched to direct UDP connection %s", peer.ID, peer.Endpoint.String())
		} else {
			Log(Info, "Received test message for peer in unsupported state")
		}
	}
}
