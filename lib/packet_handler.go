package ptp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Handlers for P2P packets received from other network peers or TURN servers

// MessageHandler is a messages callback
type MessageHandler func(message *P2PMessage, srcAddr *net.UDPAddr) error

// HandleP2PMessage is a handler for new messages received from P2P network
func (p *PeerToPeer) HandleP2PMessage(count int, srcAddr *net.UDPAddr, err error, rcvBytes []byte) error {
	if err != nil {
		Log(Error, "P2P Message Handle: %v", err)
		return err
	}
	buf := make([]byte, count)
	copy(buf[:], rcvBytes[:])

	msg, desErr := P2PMessageFromBytes(buf)
	if desErr != nil {
		Log(Error, "P2PMessageFromBytes error: %v", desErr)
		return fmt.Errorf("Failed to unmarshal message: %s", desErr.Error())
	}
	if msg == nil {
		Log(Error, "Received broken message")
		return fmt.Errorf("Broken P2P message")
	}
	// Decrypt message if crypter is active
	if p.Crypter.Active && (msg.Header.Type == MsgTypeIntro || msg.Header.Type == MsgTypeNenc || msg.Header.Type == MsgTypeIntroReq || msg.Header.Type == MsgTypeTest || msg.Header.Type == MsgTypeXpeerPing) {
		var decErr error
		msg.Data, decErr = p.Crypter.decrypt(p.Crypter.ActiveKey.Key, msg.Data)
		if decErr != nil {
			Log(Error, "Failed to decrypt message: %s", decErr)
			return fmt.Errorf("Failed to decrypt message: %s", decErr)
		}
		msg.Data = msg.Data[:msg.Header.Length]
	}

	callback, exists := p.MessageHandlers[msg.Header.Type]
	if exists {
		return callback(msg, srcAddr)
	}
	Log(Warning, "Unknown message received")
	return fmt.Errorf("Unknown message received")
}

// HandleNotEncryptedMessage is a normal message sent over p2p network
func (p *PeerToPeer) HandleNotEncryptedMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if msg.Header == nil {
		return fmt.Errorf("nil header")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	Log(Trace, "Data: %s, From: %s", msg.Data, srcAddr.String())
	p.WriteToDevice(msg.Data, msg.Header.NetProto, false)
	return nil
}

// HandlePingMessage is a PING message from a proxy handler
func (p *PeerToPeer) HandlePingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}
	if p.UDPSocket == nil {
		return fmt.Errorf("nil udp socket")
	}

	addr, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		if p.ProxyManager.touch(srcAddr.String()) {
			p.UDPSocket.SendMessage(msg, srcAddr)
		}
		return nil
	}

	port := addr.Port
	if p.UDPSocket.remotePort == 0 {
		p.UDPSocket.remotePort = port
		return nil
	}
	if port != p.UDPSocket.GetPort() && port != p.UDPSocket.remotePort && port != 0 {
		Log(Debug, "Port translation detected %d -> %d", p.UDPSocket.GetPort(), port)
		p.UDPSocket.remotePort = port
	}
	return nil
}

// HandleXpeerPingMessage receives a cross-peer ping message
func (p *PeerToPeer) HandleXpeerPingMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	if p.UDPSocket == nil {
		return fmt.Errorf("nil socket")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}
	if len(msg.Data) < 1 {
		return fmt.Errorf("message is payload is empty")
	}
	query := string(msg.Data)[:1]
	if query == "q" {
		if len(msg.Data) < 37 {
			return fmt.Errorf("payload length is too small for xpeer ping query")
		}
		id := string(msg.Data)[1:37]
		endpoint := string(msg.Data)[37:]
		response := append([]byte("r"), []byte(endpoint)...)

		msg, err := p.CreateMessage(MsgTypeXpeerPing, response, 0, true)
		if err != nil {
			Log(Debug, "Failed to create ping response: %s", err)
			return fmt.Errorf("failed to create crosspeer ping message")
		}

		// Look if we really know this peer
		for _, peer := range p.Peers.Get() {
			if peer.ID == id {
				for _, ep := range peer.KnownIPs {
					if ep.String() == srcAddr.String() {
						p.UDPSocket.SendMessage(msg, ep)
						peer.BumpEndpoint(ep.String())
						return nil
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
					return nil
				}
			}
		}
		Log(Debug, "Received ping from unknown endpoint: %s [%s ID: %s]", srcAddr.String(), endpoint, id)
		return fmt.Errorf("Received ping from unknown endpoint: %s [%s ID: %s]", srcAddr.String(), endpoint, id)
	} else if query == "r" {
		endpoint := msg.Data[1:]
		for _, peer := range p.Peers.Get() {
			if peer == nil {
				continue
			}
			for i, ep := range peer.EndpointsHeap {
				if ep.Addr.String() == string(endpoint) {
					peer.EndpointsHeap[i].LastContact = time.Now()
					return nil
				}
			}
		}
	}
	return fmt.Errorf("Broken cross peer ping message")
}

// HandleIntroMessage receives an introduction string from another peer during handshake
func (p *PeerToPeer) HandleIntroMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	Log(Debug, "Introduction string from %s", srcAddr)
	hs, err := p.ParseIntroString(string(msg.Data))
	if err != nil {
		Log(Debug, "Failed to parse handshake response: %s", err)
		return err
	}
	if len(hs.ID) != 36 {
		Log(Debug, "Received wrong ID in introduction message: %s", hs.ID)
		return fmt.Errorf("ID length mismatch in introduction message: %d", len(hs.ID))
	}
	peer := p.Peers.GetPeer(hs.ID)
	if peer == nil {
		Log(Trace, "Unknown peer in handshke response")
		return fmt.Errorf("Received unknown peer in handshake response")
	}

	peer.PeerHW = hs.HardwareAddr
	peer.PeerLocalIP = hs.IP
	peer.LastContact = time.Now()
	peer.addEndpoint(hs.Endpoint)
	for _, np := range p.Peers.Get() {
		if np == nil {
			continue
		}
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

	p.Peers.Update(hs.ID, peer)
	Log(Debug, "Connection with peer %s has been established over %s", hs.ID, hs.Endpoint.String())
	return nil
}

// HandleIntroRequestMessage is a handshake request from another peer
// First 36 bytes is an ID of original sender, data after byte 36 is an
// endpoint on which sender was trying to communicate with this peer.
// We need to send this data back to him, so he knows which endpoint
// replied
func (p *PeerToPeer) HandleIntroRequestMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	if p.UDPSocket == nil {
		return fmt.Errorf("nil udp socket")
	}
	if len(msg.Data) < 36 {
		return fmt.Errorf("payload is too short")
	}
	id := string(msg.Data[0:36])
	peer := p.Peers.GetPeer(id)
	if peer == nil {
		Log(Trace, "Introduction request came from unknown peer: %s -> %s [%s]", id, msg.Data[36:], srcAddr.String())
		return fmt.Errorf("Introduction request from unknown peer: %s -> %s [%s]", id, msg.Data[36:], srcAddr.String())
	}
	response, err := p.PrepareIntroductionMessage(p.Dht.ID, string(msg.Data[36:]))
	if err != nil {
		Log(Error, "Failed to prepare intro message: %s", err.Error())
		return fmt.Errorf("Failed to prepare introduction message: %s", err.Error())
	}
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
			Log(Error, "Failed to respond to introduction request: %s", err.Error())
			return fmt.Errorf("Failed to response to introduction reuqest: %s", err.Error())
		}
	}
	return nil
}

// HandleProxyMessage receives a control packet from proxy
// Proxy packets comes in format of UDP connection address
func (p *PeerToPeer) HandleProxyMessage(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}

	Log(Debug, "New proxy message from %s", srcAddr)
	ep, err := net.ResolveUDPAddr("udp4", string(msg.Data))
	if err != nil {
		Log(Error, "Failed to resolve proxy address: %s", err.Error())
		return fmt.Errorf("Failed to resolve proxy address: %s", err.Error())
	}
	rc := p.ProxyManager.activate(srcAddr.String(), ep)
	if rc {
		Log(Debug, "This peer is now available over %s", ep.String())
		return nil
	}
	return fmt.Errorf("Failed to activate proxy %s", ep.String())
}

// HandleBadTun notified peer about proxy being malfunction
// This method is not used in currenct scheme
// TODO: Consider to remove
func (p *PeerToPeer) HandleBadTun(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	return nil
}

// HandleLatency will handle latency respones from another peer/proxy
func (p *PeerToPeer) HandleLatency(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	if len(msg.Data) < 12 {
		return fmt.Errorf("payload is too short")
	}
	Log(Trace, "Latency response from %s", srcAddr.String())

	if bytes.Equal(msg.Data[:4], LatencyProxyHeader) {
		// This is a response from proxy

		ts := time.Time{}
		err := ts.UnmarshalBinary(msg.Data[4:])
		if err != nil {
			Log(Error, "Failed to unmarshal latency packet from %s: %s", srcAddr.String(), err.Error())
			return fmt.Errorf("Failed to unmarshal latency from %s: %s", srcAddr.String(), err.Error())
		}
		latency := time.Since(ts)

		if p.ProxyManager.setLatency(latency, srcAddr) != nil {
			Log(Error, "Couldn't set latency for proxy: %s", srcAddr)
			return fmt.Errorf("Failed to set latency for proxy %s", srcAddr.String())
		}
		return nil
	} else if bytes.Equal(msg.Data[:4], LatencyRequestHeader) {
		// This is a request of latency from endpoint

		if len(msg.Data) < 52 {
			Log(Error, "Broken latency request packet: too small [%d]", len(msg.Data))
			return fmt.Errorf("latency packet request is too small: %d bytes", len(msg.Data))
		}

		// Find this peer
		peerID := string(msg.Data[10:46])
		peer := p.Peers.GetPeer(peerID)
		if peer == nil {
			Log(Trace, "Received latency request from unknown peers: %s [Origin: %s]", peerID, srcAddr.String())
			return fmt.Errorf("latency request from unknown peer: %s [Origin: %s]", peerID, srcAddr.String())
		}
		if peer.Endpoint == nil {
			Log(Trace, "Received latency request from not integrated peer %s [Origin: %s]", peerID, srcAddr.String())
			return fmt.Errorf("Received latency request from not integrated peer %s [Origin: %s]", peerID, srcAddr.String())
		}

		Log(Trace, "Latency request from %s", srcAddr.String())
		response, err := p.CreateMessage(MsgTypeLatency, append(LatencyResponseHeader, msg.Data[4:]...), 0, false)
		if err != nil {
			Log(Error, "Failed to create latency response for %s: %s", srcAddr.String(), err.Error())
			return fmt.Errorf("Failed to create latency response for %s: %s", srcAddr.String(), err.Error())
		}

		p.UDPSocket.SendMessage(response, peer.Endpoint)
		return nil
	} else if bytes.Equal(msg.Data[:4], LatencyResponseHeader) {
		// This is a response of latency from endpoint

		if len(msg.Data) < 52 {
			Log(Error, "Broken latency response packet: too small [%d]", len(msg.Data))
			return fmt.Errorf("latency response packet is too small: %d bytes", len(msg.Data))
		}

		// Extract IP and Port
		ipfield := msg.Data[4:10]

		ip := net.IP{ipfield[0], ipfield[1], ipfield[2], ipfield[3]}
		port := binary.BigEndian.Uint16(ipfield[4:6])
		addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", ip.String(), port))
		if err != nil || addr.String() == "255.255.255.255:65535" || addr.String() == "0.0.0.0:0" {
			Log(Error, "Received malformed latency packet: address is broken")
			return fmt.Errorf("malformed latency packet: broken address")
		}

		ts := time.Time{}
		err = ts.UnmarshalBinary(msg.Data[46:])
		if err != nil {
			Log(Error, "Failed to unmarshal latency packet from %s: %s", srcAddr.String(), err.Error())
			return fmt.Errorf("failed to unmarshal latency packet from %s: %s", srcAddr.String(), err.Error())
		}
		latency := time.Since(ts)

		for _, peer := range p.Peers.Get() {
			for i, ep := range peer.EndpointsHeap {
				if ep.Addr.String() == addr.String() {
					peer.EndpointsHeap[i].Latency = latency
					return nil
				}
			}
		}
		Log(Error, "Can't set latency value for endpoint %s: Peer or endpoint wasn't found", addr.String())
		return fmt.Errorf("couldn't set latency value for endpoint %s: not found", addr.String())
	}
	Log(Error, "Malformed Latency packet from %s", srcAddr.String())
	return fmt.Errorf("malformed latency packet from %s", srcAddr.String())
}

// HandleComm is an internal communication packet for peers
func (p *PeerToPeer) HandleComm(msg *P2PMessage, srcAddr *net.UDPAddr) error {
	if p.UDPSocket == nil {
		return fmt.Errorf("nil udp socket")
	}
	if msg == nil {
		return fmt.Errorf("nil message")
	}
	if srcAddr == nil {
		return fmt.Errorf("nil source addr")
	}
	if msg.Data == nil {
		return fmt.Errorf("nil data")
	}
	if len(msg.Data) < 3 {
		return fmt.Errorf("payload is too small")
	}
	commType := binary.BigEndian.Uint16(msg.Data[0:2])
	data := msg.Data[2:]

	var response []byte
	var err error

	switch commType {
	case CommStatusReport:
		response, err = commStatusReportHandler(data, p)
		if err != nil {
			return err
		}
	case CommIPSubnet:
		response, err = commSubnetInfoHandler(data, p)
		if err != nil {
			return err
		}
	case CommIPInfo:
		response, err = commIPInfoHandler(data, p)
		if err != nil {
			return err
		}
	case CommIPSet:
		response, err = commIPSetHandler(data, p)
		if err != nil {
			return err
		}
	case CommIPConflict:
		response, err = commIPConflictHandler(data, p)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown comm type")
	}

	if response != nil {
		packet, err := p.CreateMessage(MsgTypeComm, response, 0, true)
		if err != nil {
			return err
		}
		_, err = p.UDPSocket.SendMessage(packet, srcAddr)
		return err
	}

	return fmt.Errorf("nil response")
}
