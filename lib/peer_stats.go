package ptp

import "time"

// PeerStats represents different peer statistics
// localNum, internetNum and proxyNum are the number of endpoints in local network, internet and over proxy
// connectionsNum and reconnctsNum represents number of connection attempts made during the lifetime of the peer
// PeerStats also keeps different timestamps related to connections
type PeerStats struct {
	localNum         int       // Number of local network connections
	internetNum      int       // Number of internet connections
	proxyNum         int       // Number of proxy connections
	connectionsNum   int       // Number of connections attempts in a single connection cyclce (not reconnect after connection was established)
	reconnectsNum    int       // Number of reconnects
	startedAt        time.Time // Time when peer was started
	connectedAt      time.Time // Time when peer was connected for the first time
	connectionLostAt time.Time // Time when connection was lost and reconnection cycle was initialized
	reconnectedAt    time.Time // Time when connection to the peer was reeestablished
}

// updateConnectionTime will update timestamp of the `connectedAt`
// if this is the first time when connection was established or
// `reconnectedAt` if connection was established after reconnect
func (p *PeerStats) updateConnectionTime() {
	if p.reconnectsNum == 0 {
		p.connectedAt = time.Now()
		return
	}
	p.reconnectedAt = time.Now()
}

// reconnect must be called when peer join a new connection cycle
// after connection was already established
func (p *PeerStats) reconnect() {
	p.reconnectsNum++
	p.connectionLostAt = time.Now()
}

// connectionAttempt must be called when peer initiates new connection attempt in the
// connection cycle
func (p *PeerStats) connectionAttempt() {
	p.connectionsNum++
}
