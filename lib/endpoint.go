package ptp

import (
	"net"
	"time"
)

// Endpoint reprsents a UDP address endpoint that instance
// may use for connection with a peer
type Endpoint struct {
	Addr              *net.UDPAddr
	LastContact       time.Time
	LastPing          time.Time
	broken            bool
	Latency           time.Duration
	LastLatencyQuery  time.Time
	MeasureInProgress bool
}

// Measure will prepare and send latency packet to the endpoint
func (e *Endpoint) Measure(n *Network) {
	if e.broken {
		return
	}

	if time.Since(e.LastLatencyQuery) < time.Duration(time.Second*15) {
		return
	}

	if e.MeasureInProgress {
		return
	}

	e.MeasureInProgress = true
	ts, _ := time.Now().MarshalBinary()
	msg, err := CreateMessageStatic(MsgTypeLatency, append(LatencyRequestHeader, ts...))
	if err != nil {
		Log(Error, "Failed to create latency measurement packet for endpoint: %s", err.Error())
		e.LastLatencyQuery = time.Now()
		e.MeasureInProgress = false
		return
	}
	Log(Trace, "Measuring latency with endpoint %s", e.Addr.String())
	n.SendMessage(msg, e.Addr)
}
