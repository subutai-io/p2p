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
// hash[36] - subnet[?]
// If subnet is empty, that means that this is a request. Hash is a mandatory, but just for a sanity check
func commSubnetInfoHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	//hash := data[0:36]
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// commIPInfoHandler will check if we know this IP or not
// hash[36] ip[4] res[1]
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
	if p.Peers == nil {
		return nil, fmt.Errorf("nil peer list")
	}
	if len(data) < 40 {
		return nil, fmt.Errorf("payload it soo small: %d", len(data))
	}

	hash := data[0:36]
	ip := data[36:40]
	fmt.Println(string(hash))
	fmt.Println(string(ip))
	fmt.Println(len(data))
	if len(data) == 42 {
		fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~0")
		result := binary.BigEndian.Uint16(data[40:42])
		if result == 0 {
			p.Interface.SetIP(net.IP(ip))
			p.Interface.Configure()
		}
		return nil, nil
	}
	if len(data) != 40 {
		return nil, fmt.Errorf("wrong data length: %d", len(data))
	}

	var result uint16

	for _, peer := range p.Peers.Get() {
		if bytes.Equal(peer.PeerLocalIP, ip) {
			result = 1
			break
		}
	}
	response := append(hash, ip...)
	binary.BigEndian.PutUint16(response[len(response):len(response)+2], result)
	return response, nil
}

func commIPSetHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func commIPConflictHandler(data []byte, p *PeerToPeer) ([]byte, error) {
	err := commPacketCheck(data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
