package main

import (
	"bytes"
	"net"
	"time"

	"github.com/gogo/protobuf/proto"
	ptp "github.com/subutai-io/p2p/lib"
)

// DHTRouter represents a connection to a router
type DHTRouter struct {
	conn        *net.TCPConn // TCP connection to a bootsrap node
	addr        *net.TCPAddr // TCP address of a bootstrap node
	router      string       // Address of a bootstrap node
	running     bool         // Whether router is running or not
	handshaked  bool         // Whether handshake has been completed or not
	stop        bool         // Whether service should be terminated
	fails       int          // Number of connection fails
	tx          uint64
	rx          uint64
	data        chan *ptp.DHTPacket
	lastContact time.Time
}

func (dht *DHTRouter) run() {
	dht.running = false
	dht.handshaked = false

	for !dht.stop {
		for !dht.running {
			dht.connect()
			if dht.running {
				break
			}
			if dht.stop {
				break
			}
			dht.sleep()
		}
		data := make([]byte, 1024)
		n, err := dht.conn.Read(data)
		if err != nil {
			ptp.Log(ptp.Warning, "BSN socket closed: %s", err)
			dht.running = false
			dht.handshaked = false
			continue
		}
		dht.lastContact = time.Now()
		go dht.handleData(data, n)
	}
}

func (dht *DHTRouter) handleData(data []byte, length int) {
	dht.rx += uint64(length)
	i := 0
	handled := 0
	for i >= 0 {
		i = bytes.Index(data, []byte{0x0a, 0x0b, 0x0c, 0x0a})
		if i <= 0 {
			break
		}
		handled++
		dht.routeData(data[:i])
		data = data[i:]
	}
	if handled == 0 {
		dht.routeData(data[:length])
	}
}

func (dht *DHTRouter) routeData(data []byte) {
	packet := &ptp.DHTPacket{}
	err := proto.Unmarshal(data, packet)
	if err != nil {
		ptp.Log(ptp.Warning, "Corrupted data from DHT: %s [%d]", err, len(data))
		return
	}
	ptp.Log(ptp.Trace, "Received DHT packet: %+v", packet)
	if packet.Type == ptp.DHTPacketType_Ping && dht.handshaked == false {
		supported := false
		for _, v := range ptp.SupportedVersion {
			if v == packet.Version {
				supported = true
			}
		}
		if !supported {
			ptp.Log(ptp.Error, "Version mismatch. Server have %d. We have %d", packet.Version, ptp.PacketVersion)
			dht.stop = true
			if dht.conn != nil {
				dht.conn.Close()
			}
		} else {
			dht.handshaked = true
			ptp.Log(ptp.Info, "Connected to a bootstrap node: %s [%s]", dht.addr.String(), packet.Data)
			dht.data <- packet
			return
		}
	}
	if !dht.handshaked {
		ptp.Log(ptp.Trace, "Skipping packet: not handshaked")
		return
	}
	dht.data <- packet
}

func (dht *DHTRouter) connect() {
	dht.handshaked = false
	dht.running = false

	if dht.conn != nil {
		dht.conn.Close()
		dht.conn = nil
	}

	var err error
	dht.conn, err = net.DialTCP("tcp4", nil, dht.addr)
	if err != nil {
		dht.fails++
		ptp.Log(ptp.Error, "Failed to establish connection with %s: %s", dht.addr.String(), err)
		return
	}
	dht.fails = 0
	dht.running = true
}

func (dht *DHTRouter) sleep() {
	multiplier := dht.fails * 5
	if multiplier > 30 {
		multiplier = 30
	}
	ptp.Log(ptp.Info, "Waiting for %d second before reconnecting", multiplier)
	started := time.Now()
	timeout := time.Duration(time.Second * time.Duration(multiplier))
	for time.Since(started) < timeout {
		time.Sleep(time.Millisecond * 200)
	}
}

func (dht *DHTRouter) keepAlive() {
	lastPing := time.Now()
	dht.lastContact = time.Now()
	for {
		if time.Since(lastPing) > time.Duration(time.Millisecond*30000) {
			lastPing = time.Now()
		}
		if time.Since(dht.lastContact) > time.Duration(time.Millisecond*120000) && dht.running {
			ptp.Log(ptp.Warning, "Disconnected from DHT router %s by timeout", dht.addr.String())
			dht.handshaked = false
			dht.running = false
			dht.conn.Close()
			dht.conn = nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (dht *DHTRouter) sendRaw(data []byte) (int, error) {
	return dht.conn.Write(data)
}
