package ptp

import (
	"net"
	"time"
)

type proxyServer struct {
	Addr              *net.UDPAddr  // Address of the proxy
	Endpoint          *net.UDPAddr  // Endpoint provided by proxy
	Status            proxyStatus   // Current proxy status
	LastUpdate        time.Time     // Last ping
	Created           time.Time     // Creation timestamp
	Latency           time.Duration // Measured latency
	LastLatencyQuery  time.Time     // Last latency request
	MeasureInProgress bool          // Whether or not this proxy is measuring latency currently
}

// Init will initialize Proxy Server
func (p *proxyServer) Init(addr *net.UDPAddr) error {
	p.Addr = addr
	p.Endpoint = nil
	p.Status = proxyConnecting
	p.Created = time.Now()
	p.LastLatencyQuery = time.Unix(0, 0)
	return nil
}

// Close will stop proxy
func (p *proxyServer) Close() error {
	Log(Info, "Stopping proxy %s, Endpoint: %s", p.Addr.String(), p.Endpoint.String())
	p.Addr = nil
	p.Endpoint = nil
	p.Status = proxyDisconnected
	return nil
}

// Measure will send request to a proxy peer with timestamp in it and
// proxy peer must response with the same message
func (p *proxyServer) Measure(n *Network) {
	if p.Status != proxyActive {
		return
	}

	if time.Since(p.LastLatencyQuery) < time.Duration(time.Second*15) {
		return
	}

	if p.MeasureInProgress {
		return
	}

	p.MeasureInProgress = true
	ts, _ := time.Now().MarshalBinary()
	msg, err := CreateMessageStatic(MsgTypeLatency, append(LatencyProxyHeader, ts...))
	if err != nil {
		Log(Error, "Failed to create latency measurement packet for proxy: %s", err.Error())
		p.LastLatencyQuery = time.Now()
		p.MeasureInProgress = false
		return
	}
	Log(Trace, "Measuring latency with proxy %s", p.Addr.String())
	n.SendMessage(msg, p.Addr)
}
