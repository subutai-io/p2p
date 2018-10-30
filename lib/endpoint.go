package ptp

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// Endpoint reprsents a UDP address endpoint that instance
// may use for connection with a peer
type Endpoint struct {
	Addr             *net.UDPAddr
	LastContact      time.Time
	LastPing         time.Time
	broken           bool
	Latency          time.Duration
	LastLatencyQuery time.Time
}

// Measure will prepare and send latency packet to the endpoint
// id is an ID of this peer
func (e *Endpoint) Measure(n *Network, id string) {
	if e.broken {
		return
	}

	if time.Since(e.LastLatencyQuery) < EndpointLatencyRequestInterval {
		return
	}

	e.LastLatencyQuery = time.Now()

	ts, _ := time.Now().MarshalBinary()

	payload := []byte{}
	payload = append(payload, LatencyRequestHeader...)
	payload = append(payload, e.addrToBytes()...)
	payload = append(payload, []byte(id)...)
	payload = append(payload, ts...)

	msg, err := CreateMessageStatic(MsgTypeLatency, payload)
	if err != nil {
		Log(Error, "Failed to create latency measurement packet for endpoint: %s", err.Error())
		e.LastLatencyQuery = time.Now()
		return
	}
	Log(Trace, "Measuring latency with endpoint %s", e.Addr.String())
	n.SendMessage(msg, e.Addr)
}

func (e *Endpoint) addrToBytes() []byte {
	if e.Addr == nil {
		return nil
	}

	fmt.Println(e.Addr.String())

	// 4 bytes of IP and 2 bytes of port
	ipfield := make([]byte, 6)

	ip4 := e.Addr.IP.To4()
	port := e.Addr.Port

	copy(ipfield[0:4], ip4[:4])
	binary.BigEndian.PutUint16(ipfield[4:6], uint16(port))

	// Data extract
	// net.IP{ipfield[0], ipfield[1], ipfield[2], ipfield[3]}
	// binary.BigEndian.Uint16(ipfield[4:6])

	fmt.Println(ipfield)
	return ipfield
}
