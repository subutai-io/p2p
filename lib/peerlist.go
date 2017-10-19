package ptp

import (
	"sync"
)

// PeerList is for handling list of peers with all mappings
type PeerList struct {
	peers      map[string]*NetworkPeer
	tableIPID  map[string]string // Mapping for IP->ID
	tableMacID map[string]string // Mapping for MAC->ID
	lock       sync.RWMutex
}

// Init will initialize PeerList's maps
func (l *PeerList) Init() {
	l.peers = make(map[string]*NetworkPeer)
	l.tableIPID = make(map[string]string)
	l.tableMacID = make(map[string]string)
}

// Update will append/edit peer in list
func (l *PeerList) Update(id string, peer *NetworkPeer) error {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.peers[id] = peer
	return nil
}
