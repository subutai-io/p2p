package ptp

import (
	"reflect"
	"testing"
	"time"
)

func TestPeerStats_updateConnectionTime(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"0 rec", fields{}},
		{">0 rec", fields{reconnectsNum: 1}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			p.updateConnectionTime()
		})
	}
}

func TestPeerStats_reconnect(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"empty test", fields{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			p.reconnect()
		})
	}
}

func TestPeerStats_connectionAttempt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"empty test", fields{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			p.connectionAttempt()
		})
	}
}

func TestPeerStats_holePunchAttempt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{"empty test", fields{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			p.holePunchAttempt()
		})
	}
}

func TestPeerStats_GetStartedAt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}

	tr := time.Now()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"empty test", fields{startedAt: tr}, tr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetStartedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerStats.GetStartedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetConnectedAt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}

	tr := time.Now()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"simple test", fields{connectedAt: tr}, tr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetConnectedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerStats.GetConnectedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetConnectionTimeDelta(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}

	ca := time.Unix(1, 1)
	sa := time.Unix(0, 0)

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple test", fields{startedAt: sa, connectedAt: ca}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetConnectionTimeDelta(); got != tt.want {
				t.Errorf("PeerStats.GetConnectionTimeDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetConnectionLostAt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}

	cla := time.Now()

	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"simple test", fields{connectionLostAt: cla}, cla},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetConnectionLostAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerStats.GetConnectionLostAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetReconnectedAt(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	ra := time.Now()
	tests := []struct {
		name   string
		fields fields
		want   time.Time
	}{
		{"simple test", fields{reconnectedAt: ra}, ra},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetReconnectedAt(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerStats.GetReconnectedAt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetReconnectionTimeDelta(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}

	ra := time.Unix(1, 1)
	cla := time.Unix(0, 0)

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple test", fields{reconnectedAt: ra, connectionLostAt: cla}, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetReconnectionTimeDelta(); got != tt.want {
				t.Errorf("PeerStats.GetReconnectionTimeDelta() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetHolePunchNum(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple test", fields{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetHolePunchNum(); got != tt.want {
				t.Errorf("PeerStats.GetHolePunchNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetConnectionsNum(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple test", fields{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetConnectionsNum(); got != tt.want {
				t.Errorf("PeerStats.GetConnectionsNum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerStats_GetReconnectsNum(t *testing.T) {
	type fields struct {
		localNum         int
		internetNum      int
		proxyNum         int
		connectionsNum   int
		reconnectsNum    int
		startedAt        time.Time
		connectedAt      time.Time
		connectionLostAt time.Time
		reconnectedAt    time.Time
		holePunchNum     int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"simple test", fields{}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerStats{
				localNum:         tt.fields.localNum,
				internetNum:      tt.fields.internetNum,
				proxyNum:         tt.fields.proxyNum,
				connectionsNum:   tt.fields.connectionsNum,
				reconnectsNum:    tt.fields.reconnectsNum,
				startedAt:        tt.fields.startedAt,
				connectedAt:      tt.fields.connectedAt,
				connectionLostAt: tt.fields.connectionLostAt,
				reconnectedAt:    tt.fields.reconnectedAt,
				holePunchNum:     tt.fields.holePunchNum,
			}
			if got := p.GetReconnectsNum(); got != tt.want {
				t.Errorf("PeerStats.GetReconnectsNum() = %v, want %v", got, tt.want)
			}
		})
	}
}
