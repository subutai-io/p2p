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
	holePunchNum     int       // Number of hole punch attempts made during peer lifetime
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

func (p *PeerStats) holePunchAttempt() {
	p.holePunchNum++
}

// GetStartedAt returns the time when peer was started
func (p *PeerStats) GetStartedAt() time.Time {
	return p.startedAt
}

// GetConnectedAt returns time when peer was connected for the first time
func (p *PeerStats) GetConnectedAt() time.Time {
	return p.connectedAt
}

// GetConnectionTimeDelta returns difference between `connectedAt` and `startedAt` in seconds
func (p *PeerStats) GetConnectionTimeDelta() int {
	return int(p.connectedAt.Sub(p.startedAt).Seconds()) % 60
}

// GetConnectionLostAt returns time when connection with the peer was lost for the first time
func (p *PeerStats) GetConnectionLostAt() time.Time {
	return p.connectionLostAt
}

// GetReconnectedAt returns the time when connection with the peer was reestablished for the last time
func (p *PeerStats) GetReconnectedAt() time.Time {
	return p.reconnectedAt
}

// GetReconnectionTimeDelta returns difference between `connectionsLostAt` and `reconnectedAt` in seconds
func (p *PeerStats) GetReconnectionTimeDelta() int {
	return int(p.reconnectedAt.Sub(p.connectionLostAt).Seconds()) % 60
}

// GetHolePunchNum returns number of hole punch attempts during peer lifetime
func (p *PeerStats) GetHolePunchNum() int {
	return p.holePunchNum
}

// GetConnectionsNum returns number of connection attempts during first connection cycle
func (p *PeerStats) GetConnectionsNum() int {
	return p.connectionsNum
}

// GetReconnectsNum returns the number of reconnection cycles
func (p *PeerStats) GetReconnectsNum() int {
	return p.reconnectsNum
}
