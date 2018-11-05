package ptp

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/subutai-io/p2p/protocol"
)

type dhtCallback func(*protocol.DHTPacket) error

func (p *PeerToPeer) setupTCPCallbacks() {
	if p.Dht == nil {
		Log(Error, "Can't setup TCP callbacks: DHT is nil")
		return
	}
	p.Dht.TCPCallbacks = make(map[protocol.DHTPacketType]dhtCallback)
	p.Dht.TCPCallbacks[protocol.DHTPacketType_BadProxy] = p.packetBadProxy
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Connect] = p.packetConnect
	p.Dht.TCPCallbacks[protocol.DHTPacketType_DHCP] = p.packetDHCP
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Error] = p.packetError
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Find] = p.packetFind
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Forward] = p.packetForward
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Node] = p.packetNode
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Notify] = p.packetNotify
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Ping] = p.packetPing
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Proxy] = p.packetProxy
	p.Dht.TCPCallbacks[protocol.DHTPacketType_RequestProxy] = p.packetRequestProxy
	p.Dht.TCPCallbacks[protocol.DHTPacketType_ReportProxy] = p.packetReportProxy
	p.Dht.TCPCallbacks[protocol.DHTPacketType_RegisterProxy] = p.packetRegisterProxy
	p.Dht.TCPCallbacks[protocol.DHTPacketType_ReportLoad] = p.packetReportLoad
	p.Dht.TCPCallbacks[protocol.DHTPacketType_State] = p.packetState
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Stop] = p.packetStop
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Unknown] = p.packetUnknown
	p.Dht.TCPCallbacks[protocol.DHTPacketType_Unsupported] = p.packetUnsupported
}

func (p *PeerToPeer) packetBadProxy(packet *protocol.DHTPacket) error {
	return nil
}

// Handshake response should be handled here.
func (p *PeerToPeer) packetConnect(packet *protocol.DHTPacket) error {
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if len(packet.Id) != 36 {
		return fmt.Errorf("Received malformed ID")
	}
	p.Dht.ID = packet.Id
	Log(Info, "Received personal ID for this session: %s", p.Dht.ID)
	p.Dht.Connected = true
	return nil
}

func (p *PeerToPeer) packetDHCP(packet *protocol.DHTPacket) error {
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if packet.Data != "" && packet.Extra != "" {
		ip, network, err := net.ParseCIDR(fmt.Sprintf("%s/%s", packet.Data, packet.Extra))
		if err != nil {
			Log(Error, "Failed to parse DHCP packet: %s", err)
			return err
		}
		p.Dht.IP = ip
		p.Dht.Network = network
		Log(Info, "Received network information: %s", network.String())
	}
	return nil
}

func (p *PeerToPeer) packetError(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	lvl := LogLevel(Trace)
	if packet.Data == "" {
		lvl = Error
	} else if packet.Data == "Warning" {
		lvl = Warning
	} else if packet.Data == "Error" {
		lvl = Error
	}
	Log(lvl, "Bootstrap node returns: %s", packet.Extra)
	return nil
}

func (p *PeerToPeer) packetFind(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	if len(packet.Arguments) == 0 {
		Log(Warning, "Received empty peer list")
		return nil
	}
	if packet.Data == p.Dht.ID {
		Log(Debug, "Skipping self [%s = %s]", packet.Data, p.Dht.ID)
		return nil
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}

	Log(Debug, "Received `find`: %+v", packet)
	peer := p.Peers.GetPeer(packet.Data)

	if peer == nil {
		peer := new(NetworkPeer)
		Log(Debug, "Received new peer %s", packet.Data)
		peer.ID = packet.Data
		for _, ip := range packet.Arguments {
			addr, err := net.ResolveUDPAddr("udp4", ip)
			if err != nil {
				continue
			}
			isNew := true
			for _, eip := range peer.KnownIPs {
				if eip.String() == addr.String() {
					isNew = false
				}
			}

			// Check if this endpoint is not local (own) ep
			for _, ip := range p.LocalIPs {
				if ip.Equal(addr.IP) {
					isNew = false
				}
			}

			if isNew {
				peer.KnownIPs = append(peer.KnownIPs, addr)
				Log(Debug, "Adding endpoint: %s", addr.String())
			}
		}
		for _, proxy := range packet.Proxies {
			addr, err := net.ResolveUDPAddr("udp4", proxy)
			if err != nil {
				continue
			}
			isNew := true
			for _, eproxy := range peer.Proxies {
				if eproxy.String() == addr.String() {
					isNew = false
				}
			}

			// Check if this proxy is not ours
			for _, epr := range p.ProxyManager.get() {
				if epr.Endpoint.String() == addr.String() {
					isNew = false
				}
			}

			if isNew {
				peer.Proxies = append(peer.Proxies, addr)
				Log(Debug, "Adding proxy: %s", addr.String())
			}
		}
		if packet.GetExtra() != "skip" {
			peer.SetState(PeerStateInit, p)
			peer.LastFind = time.Now()
			p.Peers.Update(peer.ID, peer)
			p.Peers.RunPeer(peer.ID, p)
		}
	} else {
		// This is an existing peer
		peer.LastFind = time.Now()

		ips := []*net.UDPAddr{}
		proxies := []*net.UDPAddr{}

		for _, ip := range packet.Arguments {
			if ip == "" {
				continue
			}
			addr, err := net.ResolveUDPAddr("udp4", ip)
			if err != nil {
				continue
			}
			isNew := true
			for _, eip := range ips {
				if eip.String() == ip {
					isNew = false
				}
			}
			if isNew {
				ips = append(ips, addr)
				Log(Debug, "Updating endpoint: %s", addr.String())
			}
		}
		peer.KnownIPs = ips
		for _, proxy := range packet.Proxies {
			if proxy == "" {
				continue
			}
			addr, err := net.ResolveUDPAddr("udp4", proxy)
			if err != nil {
				continue
			}
			isNew := true
			for _, eproxy := range proxies {
				if eproxy.String() == proxy {
					isNew = false
				}
			}
			if isNew {
				proxies = append(proxies, addr)
				Log(Debug, "Updating proxy: %s", addr.String())
			}
		}
		peer.Proxies = proxies
		p.Peers.Update(peer.ID, peer)
	}
	return nil
}

func (p *PeerToPeer) packetForward(packet *protocol.DHTPacket) error {
	return nil
}

func (p *PeerToPeer) packetNode(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}

	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}

	if len(packet.Arguments) == 0 {
		return fmt.Errorf("Empty IP's list")
	}

	peer := p.Peers.GetPeer(packet.Data)
	if peer == nil {
		return fmt.Errorf("Peer %s not found", packet.Data)
	}

	Log(Debug, "Received peer %s IPs", packet.Data)
	list := []*net.UDPAddr{}
	for _, addr := range packet.Arguments {
		if addr == "" {
			continue
		}
		ip, err := net.ResolveUDPAddr("udp4", addr)
		if err != nil {
			Log(Error, "Failed to resolve one of peer addresses: %s", err)
			continue
		}
		list = append(list, ip)
	}
	if len(list) > 0 {
		peer.KnownIPs = list
	}
	return nil
}

func (p *PeerToPeer) packetNotify(packet *protocol.DHTPacket) error {
	return nil
}

func (p *PeerToPeer) packetPing(packet *protocol.DHTPacket) error {
	return nil
}

func (p *PeerToPeer) packetProxy(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if p.UDPSocket == nil {
		return fmt.Errorf("nil socket")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	Log(Debug, "Received list of proxies")
	for _, proxy := range packet.Proxies {
		proxyAddr, err := net.ResolveUDPAddr("udp4", proxy)
		if err != nil {
			continue
		}
		if p.ProxyManager.new(proxyAddr) == nil {
			go func() {
				msg, err := p.CreateMessage(MsgTypeProxy, []byte(p.Dht.ID), 0, false)
				if err == nil {
					p.UDPSocket.SendMessage(msg, proxyAddr)
				}
			}()
		}
	}
	return nil
}

// packetRequestProxy received when we was requesting proxy to connect to some peer
func (p *PeerToPeer) packetRequestProxy(packet *protocol.DHTPacket) error {
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	list := []*net.UDPAddr{}
	for _, proxy := range packet.Proxies {
		addr, err := net.ResolveUDPAddr("udp4", proxy)
		if err != nil {
			Log(Error, "Can't parse proxy %s for peer %s", proxy, packet.Data)
			continue
		}
		list = append(list, addr)
	}

	peer := p.Peers.GetPeer(packet.Data)
	if peer != nil {
		peer.Proxies = list
	}
	return nil
}

func (p *PeerToPeer) packetReportProxy(packet *protocol.DHTPacket) error {
	Log(Info, "DHT confirmed proxy registration")
	return nil
}

func (p *PeerToPeer) packetRegisterProxy(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if packet.Data == "OK" {
		Log(Info, "Proxy registration confirmed")
	}
	return nil
}

func (p *PeerToPeer) packetReportLoad(packet *protocol.DHTPacket) error {
	return nil
}

func (p *PeerToPeer) packetState(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if p.Peers == nil {
		return fmt.Errorf("nil peer list")
	}
	if len(packet.Data) != 36 {
		return fmt.Errorf("Receied state packet for unknown/broken ID")
	}
	if len(packet.Extra) == 0 {
		return fmt.Errorf("Received wrong/malformed state")
	}
	numericState, err := strconv.Atoi(packet.Extra)
	if err != nil {
		return fmt.Errorf("Failed to parse state: %s", err)
	}

	peer := p.Peers.GetPeer(packet.Data)
	if peer != nil {
		peer.RemoteState = PeerState(numericState)
		p.Peers.Update(packet.Data, peer)
		Log(Debug, "Peer %s reported state '%s'", peer.ID, StringifyState(peer.RemoteState))
	} else {
		Log(Trace, "Received state of unknown peer. Updating peers")
		//p.Dht.sendFind()
	}
	return nil
}

func (p *PeerToPeer) packetStop(packet *protocol.DHTPacket) error {
	return nil
}

func (p *PeerToPeer) packetUnknown(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	if p.ProxyManager == nil {
		return fmt.Errorf("nil proxy manager")
	}
	if p.Interface == nil {
		return fmt.Errorf("nil interface")
	}
	Log(Debug, "Received unknown packet")
	p.FindNetworkAddresses()
	if len(packet.Data) > 0 && packet.Data == "DHCP" {
		Log(Warning, "Network information was requested")
		p.ReportIP(p.Interface.GetIP().String(), p.Interface.GetHardwareAddress().String(), p.Interface.GetName())
		return nil
	}
	Log(Warning, "Bootstap node refuses our identity. Reconnecting")
	return p.Dht.Connect(p.LocalIPs, p.ProxyManager.GetList())
}

func (p *PeerToPeer) packetUnsupported(packet *protocol.DHTPacket) error {
	if packet == nil {
		return fmt.Errorf("nil packet")
	}
	if p.Dht == nil {
		return fmt.Errorf("nil dht")
	}
	Log(Error, "Bootstap node doesn't support our version. Shutting down")
	return p.Dht.Close()
}
