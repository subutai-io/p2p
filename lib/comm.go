package ptp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
)

// Cross-peer communication handlers

// commPacketCheck is a common packet checker that checks for
// incoming data length
func commPacketCheck(data []byte) error {
	if len(data) < 36 {
		return fmt.Errorf("data is too small for communication packet")
	}
	return nil
}

// commStatusReportHandler handles status reports from another peer
func commStatusReportHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// commSubnetInfoHandler request/response of network subnet. Data format is as follows:
// id[36] - subnet[?]
// If subnet is empty, that means that this is a request. Hash is a mandatory, but just for a sanity check
func commSubnetInfoHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	if p.Interface == nil {
		return nil, fmt.Errorf("nil interface")
	}
	if p.Dht == nil {
		return nil, fmt.Errorf("nil dht")
	}
	//hash := data[0:36]
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}

	if len(data) == 36 {
		// This is a request
		// We are not allowed to reply, if we use automatic IP
		if p.Interface.IsAuto() {
			return nil, nil
		}
		response := make([]byte, 42)
		binary.BigEndian.PutUint16(response[0:2], CommIPSubnet)
		copy(response[2:38], p.Dht.ID)
		copy(response[38:42], p.Interface.GetIP().Mask(net.CIDRMask(24, 32)).To4())
		return response, nil
	}

	if len(data) != 40 {
		return nil, fmt.Errorf("wrong payload size: %d", len(data))
	}

	// This is a response. We just a subnet on the interface
	p.Interface.SetSubnet(net.IP(data[36:40]))

	return nil, nil
}

// commIPInfoHandler will check if we know this IP or not
// id[36] ip[4] res[1]
// When res is empty - packet is a request
// res can be 0 - IP unknown
// res can be 1 - IP known
func commIPInfoHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, fmt.Errorf("nil ptp")
	}
	if p.Swarm == nil {
		return nil, fmt.Errorf("nil peer list")
	}
	if p.Interface == nil {
		return nil, fmt.Errorf("nil interface")
	}
	if p.Dht == nil {
		return nil, fmt.Errorf("nil dht")
	}
	if len(data) < 40 {
		return nil, fmt.Errorf("payload it soo small: %d", len(data))
	}

	//hash := data[0:36]
	ip := net.IP(data[36:40])
	if len(data) == 42 {
		result := binary.BigEndian.Uint16(data[40:42])
		if result == 0 && p.Interface.GetIP() == nil {
			Log(Info, "IP %s is unknown to this swarm. Setting it", ip.String())
			p.Interface.SetIP(ip)
			p.Interface.Configure(false)
			p.Interface.MarkConfigured()
			go p.notifyIP()
			return nil, nil
		}
		Log(Info, "IP %s is already known to this swarm. Ignoring it", ip.String())
		return nil, nil
	}
	if len(data) != 40 {
		return nil, fmt.Errorf("wrong data length: %d", len(data))
	}

	var result uint16

	for _, peer := range p.Swarm.Get() {
		if bytes.Equal(peer.PeerLocalIP.To4(), ip) {
			result = 1
			break
		}
	}

	if result == 1 {
		Log(Debug, "Peer requested info about IP %s. That IP is known to us", ip.String())
	} else {
		Log(Debug, "Peer requested info about IP %s. We don't know that IP", ip.String())
	}

	response := make([]byte, 44)
	binary.BigEndian.PutUint16(response[0:2], CommIPInfo)
	copy(response[2:38], p.Dht.ID)
	copy(response[38:42], ip)
	binary.BigEndian.PutUint16(response[42:44], result)
	return response, nil
}

// commIPSetHandler will handle notification from another peer that he
// is now uses specified ip
// id[36] ip[4]
func commIPSetHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	if p.Swarm == nil {
		return nil, fmt.Errorf("nil swarm")
	}
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}

	if len(data) < 40 {
		return nil, fmt.Errorf("data is too small")
	}

	id := string(data[0:36])
	ip := net.IP(data[36:40])

	for _, peer := range p.Swarm.Get() {
		if peer.ID == id {
			peer.PeerLocalIP = ip
			return nil, nil
		}
	}

	return nil, fmt.Errorf("Can't update peer IP. Peer %s not found", id)
}

func commIPConflictHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
