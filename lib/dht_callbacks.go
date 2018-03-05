package ptp

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

func (dht *DHTClient) setupTCPCallbacks() {
	dht.TCPCallbacks = make(map[DHTPacketType]dhtCallback)
	dht.TCPCallbacks[DHTPacketType_BadProxy] = dht.packetBadProxy
	dht.TCPCallbacks[DHTPacketType_Connect] = dht.packetConnect
	dht.TCPCallbacks[DHTPacketType_DHCP] = dht.packetDHCP
	dht.TCPCallbacks[DHTPacketType_Error] = dht.packetError
	dht.TCPCallbacks[DHTPacketType_Find] = dht.packetFind
	dht.TCPCallbacks[DHTPacketType_Forward] = dht.packetForward
	dht.TCPCallbacks[DHTPacketType_Node] = dht.packetNode
	dht.TCPCallbacks[DHTPacketType_Notify] = dht.packetNotify
	dht.TCPCallbacks[DHTPacketType_Ping] = dht.packetPing
	dht.TCPCallbacks[DHTPacketType_Proxy] = dht.packetProxy
	dht.TCPCallbacks[DHTPacketType_RequestProxy] = dht.packetRequestProxy
	dht.TCPCallbacks[DHTPacketType_ReportProxy] = dht.packetReportProxy
	dht.TCPCallbacks[DHTPacketType_RegisterProxy] = dht.packetRegisterProxy
	dht.TCPCallbacks[DHTPacketType_ReportLoad] = dht.packetReportLoad
	dht.TCPCallbacks[DHTPacketType_State] = dht.packetState
	dht.TCPCallbacks[DHTPacketType_Stop] = dht.packetStop
	dht.TCPCallbacks[DHTPacketType_Unknown] = dht.packetUnknown
	dht.TCPCallbacks[DHTPacketType_Unsupported] = dht.packetUnsupported
}

func (dht *DHTClient) packetBadProxy(packet *DHTPacket) error {
	return nil
}

// Handshake response should be handled here.
func (dht *DHTClient) packetConnect(packet *DHTPacket) error {
	if len(packet.Id) != 36 {
		return fmt.Errorf("Received malformed ID")
	}
	dht.ID = packet.Id
	Log(Info, "Received personal ID for this session: %s", dht.ID)
	return nil
}

func (dht *DHTClient) packetDHCP(packet *DHTPacket) error {
	Log(Info, "Received DHCP packet")
	if packet.Data != "" && packet.Extra != "" {
		ip, network, err := net.ParseCIDR(fmt.Sprintf("%s/%s", packet.Data, packet.Extra))
		if err != nil {
			Log(Error, "Failed to parse DHCP packet: %s", err)
			return err
		}
		dht.IP = ip
		dht.Network = network
		Log(Info, "Received network information: %s", network.String())
	}
	return nil
}

func (dht *DHTClient) packetError(packet *DHTPacket) error {
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

func (dht *DHTClient) packetFind(packet *DHTPacket) error {
	if len(packet.Arguments) == 0 {
		Log(Warning, "Received empty peer list")
		return nil
	}
	Log(Debug, "Received peer list")
	for _, id := range packet.Arguments {
		if id == dht.ID {
			continue
		}
		peer := NetworkPeer{ID: id}
		dht.PeerData <- peer
	}
	return nil
}

func (dht *DHTClient) packetForward(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetNode(packet *DHTPacket) error {
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

	if len(list) == 0 {
		return fmt.Errorf("Received empty IP list for peer %s", packet.Data)
	}
	peer := NetworkPeer{ID: packet.Data, KnownIPs: list}
	dht.PeerData <- peer
	return nil
}

func (dht *DHTClient) packetNotify(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetPing(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetProxy(packet *DHTPacket) error {
	Log(Debug, "Received list of proxies")
	for _, proxy := range packet.Proxies {
		dht.ProxyChannel <- proxy
	}
	return nil
}

// packetRequestProxy received when we was requesting proxy to connect to some peer
func (dht *DHTClient) packetRequestProxy(packet *DHTPacket) error {
	list := []*net.UDPAddr{}
	for _, proxy := range packet.Proxies {
		addr, err := net.ResolveUDPAddr("udp4", proxy)
		if err != nil {
			Log(Error, "Can't parse proxy %s for peer %s", proxy, packet.Data)
			continue
		}
		list = append(list, addr)
	}
	peer := NetworkPeer{ID: packet.Data, Proxies: list}
	dht.PeerData <- peer
	return nil
}

func (dht *DHTClient) packetReportProxy(packet *DHTPacket) error {
	Log(Info, "DHT confirmed proxy registration")
	return nil
}

func (dht *DHTClient) packetRegisterProxy(packet *DHTPacket) error {
	if packet.Data == "OK" {
		Log(Info, "Proxy registration confirmed")
	}
	return nil
}

func (dht *DHTClient) packetReportLoad(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetState(packet *DHTPacket) error {
	if len(packet.Data) != 36 {
		return fmt.Errorf("Receied state packet for unknown/broken ID")
	}
	if len(packet.Arguments) != 1 {
		return fmt.Errorf("Received wrong/malformed state")
	}
	numericState, err := strconv.Atoi(packet.Arguments[0])
	if err != nil {
		Log(Error, "Failed to parse state: %s", err)
	}
	state := RemotePeerState{}
	state.ID = packet.Data
	state.State = PeerState(numericState)
	dht.StateChannel <- state
	return nil
}

func (dht *DHTClient) packetStop(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetUnknown(packet *DHTPacket) error {
	Log(Warning, "Bootstap node refuses our identity. Shutting down")
	return dht.Close()
}

func (dht *DHTClient) packetUnsupported(packet *DHTPacket) error {
	Log(Error, "Bootstap node doesn't support our version. Shutting down")
	os.Exit(0)
	return dht.Close()
}
