package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/mdlayher/ethernet"
	"io"
	"log"
	"net"
	"p2p/udpcs"
)

var (
	// ErrInvalidHardwareAddr is returned when one or more invalid hardware
	// addresses are passed to NewPacket.
	ErrInvalidHardwareAddr = errors.New("invalid hardware address")

	// ErrInvalidIP is returned when one or more invalid IPv4 addresses are
	// passed to NewPacket.
	ErrInvalidIP = errors.New("invalid IPv4 address")

	// errInvalidARPPacket is returned when an ethernet frame does not
	// indicate that an ARP packet is contained in its payload.
	errInvalidARPPacket = errors.New("invalid ARP packet")
)

// Handles a IPv4 packet and sends it to it's destination
func (ptp *PTPCloud) handlePacketIPv4(contents []byte, proto int) {
	f := new(ethernet.Frame)
	if err := f.UnmarshalBinary(contents); err != nil {
		log.Printf("[ERROR] Failed to unmarshal IPv4 packet")
	}

	if f.EtherType != ethernet.EtherTypeIPv4 {
		return
	}

	msg := udpcs.CreateNencP2PMessage(contents, uint16(proto))
	size, err := ptp.SendTo(f.Destination, msg)
	if err != nil {
		log.Printf("[ERROR] Failed to send message inside P2P: %v", err)
		return
	}
	if size <= 0 {
		log.Printf("[ERROR] Zero bytes were send inside P2P")
	}
}

func (ptp *PTPCloud) handlePacketIPv6(contents []byte) {

}

func (ptp *PTPCloud) handlePacketARP(contents []byte) {
	// Prepare new ethernet frame and fill it with
	// contents of the packet
	f := new(ethernet.Frame)
	if err := f.UnmarshalBinary(contents); err != nil {
		log.Printf("[ERROR] Failed to Unmarshal ARP Binary")
		return
	}

	if f.EtherType != ethernet.EtherTypeARP {
		return
	}

	p := new(ARPPacket)
	if err := p.UnmarshalARP(f.Payload); err != nil {
		return
	}

	//for _, peer := ptp.NetworkPeers {

	//}
	log.Printf("[DEBUG] Peers: %v, Target IP: %s", ptp.NetworkPeers, p.TargetIP.String())
	var hwAddr net.HardwareAddr = nil
	for _, peer := range ptp.NetworkPeers {
		if peer.PeerLocalIP.String() == p.TargetIP.String() {
			hwAddr = peer.PeerHW
		}
	}
	/*
		hwAddr, err := net.ParseMAC("0c:8b:fd:ab:30:ee")
		if err != nil {
			log.Printf("[ERROR] Failed to parse MAC")
		}
	*/
	// TODO: Put there normal IP from list of ips
	// Send a reply
	if hwAddr == nil {
		log.Printf("[ERROR] Cannot find hardware address for requested IP")
		return
	}
	var reply ARPPacket
	ip := net.ParseIP(p.TargetIP.String())
	response, err := reply.NewPacket(OperationReply, hwAddr, ip, p.SenderHardwareAddr, p.SenderIP)
	if err != nil {
		log.Printf("[ERROR] Failed to create ARP reply")
		return
	}
	rp, err := response.MarshalBinary()
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ARP response packet")
		return
	}

	fr := &ethernet.Frame{
		Destination: response.TargetHardwareAddr,
		Source:      response.SenderHardwareAddr,
		EtherType:   ethernet.EtherTypeARP,
		Payload:     rp,
	}

	fb, err := fr.MarshalBinary()
	if err != nil {
		log.Printf("[ERROR] Failed to marshal ARP Ethernet Frame")
	}
	ptp.WriteToDevice(fb)

	log.Printf("[DEBUG] %v", p.String())
}

type Operation uint16

const (
	OperationRequest Operation = 1
	OperationReply   Operation = 2
)

type ARPPacket struct {
	// HardwareType specifies an IANA-assigned hardware type, as described
	// in RFC 826.
	HardwareType uint16

	// ProtocolType specifies the internetwork protocol for which the ARP
	// request is intended.  Typically, this is the IPv4 EtherType.
	ProtocolType uint16

	// HardwareAddrLength specifies the length of the sender and target
	// hardware addresses included in a Packet.
	HardwareAddrLength uint8

	// IPLength specifies the length of the sender and target IPv4 addresses
	// included in a Packet.
	IPLength uint8

	// Operation specifies the ARP operation being performed, such as request
	// or reply.
	Operation Operation

	// SenderHardwareAddr specifies the hardware address of the sender of this
	// Packet.
	SenderHardwareAddr net.HardwareAddr

	// SenderIP specifies the IPv4 address of the sender of this Packet.
	SenderIP net.IP

	// TargetHardwareAddr specifies the hardware address of the target of this
	// Packet.
	TargetHardwareAddr net.HardwareAddr

	// TargetIP specifies the IPv4 address of the target of this Packet.
	TargetIP net.IP
}

func (p *ARPPacket) String() string {
	return fmt.Sprintf("HWType %d, Proto: %d, HWAddrLength: %d, IPLength: %d, Operation: %d, SHWAddr: %s, SIP: %s, THWAddr: %s, TIP: %s", p.HardwareType, p.ProtocolType, p.HardwareAddrLength, p.IPLength, p.Operation, p.SenderHardwareAddr.String(), p.SenderIP.String(), p.TargetHardwareAddr.String(), p.TargetIP.String())
}

// MarshalBinary allocates a byte slice containing the data from a Packet.
//
// MarshalBinary never returns an error.
func (p *ARPPacket) MarshalBinary() ([]byte, error) {
	// 2 bytes: hardware type
	// 2 bytes: protocol type
	// 1 byte : hardware address length
	// 1 byte : protocol length
	// 2 bytes: operation
	// N bytes: source hardware address
	// N bytes: source protocol address
	// N bytes: target hardware address
	// N bytes: target protocol address

	// Though an IPv4 address should always 4 bytes, go-fuzz
	// very quickly created several crasher scenarios which
	// indicated that these values can lie.
	b := make([]byte, 2+2+1+1+2+(p.IPLength*2)+(p.HardwareAddrLength*2))

	// Marshal fixed length data

	binary.BigEndian.PutUint16(b[0:2], p.HardwareType)
	binary.BigEndian.PutUint16(b[2:4], p.ProtocolType)

	b[4] = p.HardwareAddrLength
	b[5] = p.IPLength

	binary.BigEndian.PutUint16(b[6:8], uint16(p.Operation))

	// Marshal variable length data at correct offset using lengths
	// defined in p

	n := 8
	hal := int(p.HardwareAddrLength)
	pl := int(p.IPLength)

	copy(b[n:n+hal], p.SenderHardwareAddr)
	n += hal

	copy(b[n:n+pl], p.SenderIP)
	n += pl

	copy(b[n:n+hal], p.TargetHardwareAddr)
	n += hal

	copy(b[n:n+pl], p.TargetIP)

	return b, nil
}

func (p *ARPPacket) UnmarshalARP(b []byte) error {
	// Must have enough room to retrieve hardware address and IP lengths
	if len(b) < 8 {
		return io.ErrUnexpectedEOF
	}

	// Retrieve fixed length data

	p.HardwareType = binary.BigEndian.Uint16(b[0:2])
	p.ProtocolType = binary.BigEndian.Uint16(b[2:4])

	p.HardwareAddrLength = b[4]
	p.IPLength = b[5]

	p.Operation = Operation(binary.BigEndian.Uint16(b[6:8]))

	// Unmarshal variable length data at correct offset using lengths
	// defined by ml and il
	//
	// These variables are meant to improve readability of offset calculations
	// for the code below
	n := 8
	ml := int(p.HardwareAddrLength)
	ml2 := ml * 2
	il := int(p.IPLength)
	il2 := il * 2

	// Must have enough room to retrieve both hardware address and IP addresses
	addrl := n + ml2 + il2
	if len(b) < addrl {
		return io.ErrUnexpectedEOF
	}

	// Allocate single byte slice to store address information, which
	// is resliced into fields
	bb := make([]byte, addrl-n)

	// Sender hardware address
	copy(bb[0:ml], b[n:n+ml])
	p.SenderHardwareAddr = bb[0:ml]
	n += ml

	// Sender IP address
	copy(bb[ml:ml+il], b[n:n+il])
	p.SenderIP = bb[ml : ml+il]
	n += il

	// Target hardware address
	copy(bb[ml+il:ml2+il], b[n:n+ml])
	p.TargetHardwareAddr = bb[ml+il : ml2+il]
	n += ml

	// Target IP address
	copy(bb[ml2+il:ml2+il2], b[n:n+il])
	p.TargetIP = bb[ml2+il : ml2+il2]

	return nil
}

func (p *ARPPacket) NewPacket(op Operation, srcHW net.HardwareAddr, srcIP net.IP, dstHW net.HardwareAddr, dstIP net.IP) (*ARPPacket, error) {
	// Validate hardware addresses for minimum length, and matching length
	if len(srcHW) < 6 {
		return nil, ErrInvalidHardwareAddr
	}
	if len(dstHW) < 6 {
		return nil, ErrInvalidHardwareAddr
	}
	if len(srcHW) != len(dstHW) {
		return nil, ErrInvalidHardwareAddr
	}

	// Validate IP addresses to ensure they are IPv4 addresses, and
	// correct length
	srcIP = srcIP.To4()
	if srcIP == nil {
		return nil, ErrInvalidIP
	}
	dstIP = dstIP.To4()
	if dstIP == nil {
		return nil, ErrInvalidIP
	}

	return &ARPPacket{
		// There is no Go-native way to detect hardware type of a network
		// interface, so default to 1 (ethernet 10Mb) for now
		HardwareType: 1,

		// Default to EtherType for IPv4
		ProtocolType: uint16(ethernet.EtherTypeIPv4),

		// Populate other fields using input data
		HardwareAddrLength: uint8(len(srcHW)),
		IPLength:           uint8(len(srcIP)),
		Operation:          op,
		SenderHardwareAddr: srcHW,
		SenderIP:           srcIP,
		TargetHardwareAddr: dstHW,
		TargetIP:           dstIP,
	}, nil
}
