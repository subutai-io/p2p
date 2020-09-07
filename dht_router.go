package main

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/golang/protobuf/proto"
	ptp "github.com/subutai-io/p2p/lib"
	"github.com/subutai-io/p2p/protocol"
)

// DHTRouter represents a connection to a router
type DHTRouter struct {
	conn          *net.TCPConn             // TCP connection to a bootsrap node
	addr          *net.TCPAddr             // TCP address of a bootstrap node
	router        string                   // Address of a bootstrap node
	running       bool                     // Whether router is running or not
	handshaked    bool                     // Whether handshake has been completed or not
	stop          bool                     // Whether service should be terminated
	fails         int                      // Number of connection fails
	tx            uint64                   // Transferred bytes
	rx            uint64                   // Received bytes
	data          chan *protocol.DHTPacket // Payload channel
	lastContact   time.Time                // Last communication
	packetVersion string                   // Version of packet on DHT
	version       string                   // Version of DHT
}

func (dht *DHTRouter) run() {
	dht.running = false
	dht.handshaked = false
	dht.version = "Unknown"
	dht.packetVersion = "Unknown"

	for !dht.stop {
		for !dht.running {
			dht.connect()
			if dht.running || dht.stop {
				break
			}
			dht.sleep()
		}
		data := make([]byte, ptp.DHTBufferSize)
		n, err := dht.conn.Read(data)
		if err != nil {
			ptp.Warn("BSN socket closed: %s", err)
			dht.running = false
			dht.handshaked = false
			continue
		}
		dht.lastContact = time.Now()
		go dht.handleData(data[:n])
	}
}

func (dht *DHTRouter) handleData(data []byte) {
	length := len(data)
	dht.rx += uint64(length)
	i := 0
	handled := 0
	ptp.Trace("Handling data: data length is [%d]", len(data))
	for i >= 0 {
		i = bytes.Index(data, []byte{0x0a, 0x0b, 0x0c, 0x0a})
		if i <= 0 {
			break
		}
		handled++
		dht.routeData(data[:i])
		if i+4 < len(data) {
			data = data[i+4:]
		} else {
			break
		}
	}
	if handled == 0 {
		dht.routeData(data[:length])
	}
}

func (dht *DHTRouter) routeData(data []byte) {
	packet := &protocol.DHTPacket{}
	err := proto.Unmarshal(data, packet)
	ptp.Trace("DHTPacket size: [%d]", len(data))
	ptp.Trace("DHTPacket contains: %+v --- %+v", bytes.NewBuffer(data).String(), packet)
	if err != nil {
		ptp.Warn("Corrupted data from DHT: %s [%d]", err, len(data))
		return
	}
	ptp.Trace("Received DHT packet: %+v", packet)
	if packet.Type == protocol.DHTPacketType_Ping && dht.handshaked == false {
		supported := false
		for _, v := range ptp.SupportedVersion {
			if v == packet.Version {
				supported = true
			}
		}
		if !supported {
			ptp.Error("Version mismatch. Server have %d. We have %d", packet.Version, ptp.PacketVersion)
			dht.stop = true
			if dht.conn != nil {
				dht.conn.Close()
			}
		} else {
			dht.handshaked = true
			ptp.Info("Connected to a bootstrap node: %s [%s]", dht.addr.String(), packet.Data)
			dht.packetVersion = fmt.Sprintf("%d", packet.Version)
			if packet.Extra != "" {
				ptp.Info("DHT Version: %s", packet.Extra)
				dht.version = packet.Extra
			}
			packet.Query = "handshaked"
			dht.data <- packet
			return
		}
	}
	if !dht.handshaked {
		ptp.Trace("Skipping packet: not handshaked")
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
		ptp.Error("Failed to establish connection with %s: %s", dht.addr.String(), err)
		return
	}
	dht.lastContact = time.Now()
	dht.fails = 0
	dht.running = true
}

func (dht *DHTRouter) sleep() {
	multiplier := dht.fails * 5
	if multiplier > 30 {
		multiplier = 30
	}
	ptp.Info("Waiting for %d second before reconnecting", multiplier)
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
		if time.Since(lastPing) > time.Duration(time.Millisecond*30000) && time.Since(dht.lastContact) > time.Duration(time.Millisecond*40) {
			lastPing = time.Now()
			if dht.ping() != nil {
				ptp.Error("DHT router ping failed")
			}
		}
		if time.Since(dht.lastContact) > time.Duration(time.Millisecond*60000) && dht.running {
			ptp.Warn("Disconnected from DHT router %s by timeout", dht.addr.String())
			dht.handshaked = false
			dht.running = false
			dht.conn.Close()
			dht.conn = nil
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func (dht *DHTRouter) sendRaw(data []byte) (int, error) {
	if dht.conn == nil {
		return -1, fmt.Errorf("Can't send: connection is nil")
	}
	return dht.conn.Write(data)
}

func (dht *DHTRouter) ping() error {
	ptp.Trace("Sending ping to dht %s", dht.addr.String())
	packet := &protocol.DHTPacket{
		Type:    protocol.DHTPacketType_Ping,
		Query:   "req",
		Version: ptp.PacketVersion,
	}
	data, err := proto.Marshal(packet)
	if err != nil {
		return err
	}
	_, err = dht.sendRaw(data)
	if err != nil {
		return err
	}
	return nil
}
