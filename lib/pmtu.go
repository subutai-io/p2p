package ptp

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"

	"github.com/mdlayher/ethernet"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func checksum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}
	for {
		// Break when sum is less or equals to 0xFFFF
		if csum <= 65535 {
			break
		}
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}

func pmtu(data []byte, tap TAP) (bool, error) {
	protocol := int(binary.BigEndian.Uint16(data[12:14]))
	length := len(data)

	if protocol == int(PacketIPv4) && length > GlobalMTU-150 {
		header, err := ipv4.ParseHeader(data[14:])
		if err != nil {
			Error("Failed to parse IPv4 packet: %s", err.Error())
			return false, nil
		}

		// Don't fragment flag is set. We need to respond with ICMP Destination Unreachable
		if header.Flags == ipv4.DontFragment {
			// Extract packet contents as an ethernet frame for later re-use
			f := new(ethernet.Frame)
			if err := f.UnmarshalBinary(data); err != nil {
				Error("Failed to Unmarshal IPv4")
				return false, nil
			}

			// Build "Fragmentation needed" ICMP message
			packetICMP := &icmp.Message{
				Type: ipv4.ICMPTypeDestinationUnreachable,
				Code: 4,
				Body: &icmp.PacketTooBig{
					MTU:  GlobalMTU - 200,    // Next-hop MTU
					Data: data[14 : 14+20+8], // Original header and 64-bits of datagram
				},
			}
			payloadICMP, err := packetICMP.Marshal(nil)
			if err != nil {
				Error("Failed to marshal ICMP: %s", err.Error())
				return false, errICMPMarshalFailed
			}

			// Build IPv4 Header
			iph := &ipv4.Header{
				Version:  4,
				Len:      20, // Precalculated header length
				TOS:      0,
				TotalLen: len(payloadICMP) + 20,
				ID:       25,
				TTL:      64,
				Protocol: 1,
				Dst:      header.Src,
				Src:      header.Dst,
				Checksum: 0,
			}
			ipHeader, err := iph.Marshal()
			if err != nil {
				Error("Failed to marshal header: %s", err.Error())
				return false, nil
			}

			// Calculate IPv4 header checksum
			hcsum := checksum(ipHeader)
			binary.BigEndian.PutUint16(ipHeader[10:], hcsum)

			// Build new ethernet frame. Swap dst/src
			pl := append(ipHeader, payloadICMP...)
			nf := new(ethernet.Frame)
			nf.Destination = f.Source
			nf.Source = f.Destination
			nf.EtherType = ethernet.EtherTypeIPv4
			nf.Payload = pl
			rpacket, err := nf.MarshalBinary()
			if err != nil {
				Error("Failed to marshal ethernet")
				return false, nil
			}

			// Calculate CRC32 checksum for ethernet frame
			crc := make([]byte, 4)
			binary.LittleEndian.PutUint32(crc, crc32.ChecksumIEEE(rpacket))
			rpacket = append(rpacket, crc...)

			// Send frame to the interface
			// P2P will drop packet afterwards
			tap.WritePacket(&Packet{int(PacketIPv4), rpacket})
			return true, nil
		}
	}
	return false, fmt.Errorf("Unsupported protocol for PMTU")
}
