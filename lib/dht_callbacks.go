package ptp

import (
	"fmt"
	"net"
	"strconv"
)

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

	return dht.sendFind()
}

func (dht *DHTClient) packetDHCP(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetError(packet *DHTPacket) error {
	return nil
}

func (dht *DHTClient) packetFind(packet *DHTPacket) error {
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
		ip, err := net.ResolveUDPAddr("udp", addr)
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
	return nil
}

func (dht *DHTClient) packetRegisterProxy(packet *DHTPacket) error {
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
	Log(Warning, "Bootstap node refuses our identity")
	for _, conn := range dht.TCPConnection {
		conn.Close()
	}
	dht.shutdown()
	return nil
}

func (dht *DHTClient) packetUnsupported(packet *DHTPacket) error {
	Log(Error, "Bootstap node doesn't support our version")
	for _, conn := range dht.TCPConnection {
		conn.Close()
	}
	dht.shutdown()
	return nil
}
