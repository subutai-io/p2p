package ptp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

// Constants
const (
	MagicCookie uint16 = 0xabcd
	HeaderSize  int    = 10
)

// P2PMessageHeader is header used in cross-peer packets
type P2PMessageHeader struct {
	Magic         uint16
	Type          uint16
	Length        uint16
	SerializedLen uint16
	NetProto      uint16
}

// P2PMessage is a cross-peer message packet
type P2PMessage struct {
	Header *P2PMessageHeader
	Data   []byte
}

// Serialize does a header serialization
func (v *P2PMessageHeader) Serialize() []byte {
	resBuf := make([]byte, HeaderSize)
	binary.BigEndian.PutUint16(resBuf[0:2], v.Magic)
	binary.BigEndian.PutUint16(resBuf[2:4], v.Type)
	binary.BigEndian.PutUint16(resBuf[4:6], v.Length)
	binary.BigEndian.PutUint16(resBuf[6:8], v.NetProto)
	binary.BigEndian.PutUint16(resBuf[8:10], v.SerializedLen)
	return resBuf
}

// P2PMessageHeaderFromBytes extracts message header from received packet
func P2PMessageHeaderFromBytes(bytes []byte) (*P2PMessageHeader, error) {
	if len(bytes) < HeaderSize {
		if len(bytes) == 2 {
			return nil, nil
		}
		return nil, errors.New("P2PMessageHeaderFromBytes_error : less then 14 bytes")
	}

	result := new(P2PMessageHeader)
	result.Magic = binary.BigEndian.Uint16(bytes[0:2])
	result.Type = binary.BigEndian.Uint16(bytes[2:4])
	result.Length = binary.BigEndian.Uint16(bytes[4:6])
	result.NetProto = binary.BigEndian.Uint16(bytes[6:8])
	result.SerializedLen = binary.BigEndian.Uint16(bytes[8:10])
	return result, nil
}

// Serialize constructs a P2P message
func (v *P2PMessage) Serialize() []byte {
	v.Header.SerializedLen = uint16(len(v.Data))
	resBuf := v.Header.Serialize()
	resBuf = append(resBuf, v.Data...)
	return resBuf
}

// P2PMessageFromBytes extract a payload from received packet
func P2PMessageFromBytes(bytes []byte) (*P2PMessage, error) {
	res := new(P2PMessage)
	var err error
	res.Header, err = P2PMessageHeaderFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	if res.Header == nil {
		return nil, nil
	}
	// Log(Trace, "--- P2PMessageHeaderFromBytes Length : %d, SerLen : %d", res.Header.Length, res.Header.SerializedLen)
	if res.Header.Magic != MagicCookie {
		return nil, errors.New("magic cookie not presented")
	}
	res.Data = make([]byte, res.Header.SerializedLen)
	copy(res.Data[:], bytes[HeaderSize:])
	return res, err
}

// CreateMessage create internal P2P Message
func (p *PeerToPeer) CreateMessage(msgType MsgType, payload []byte, proto uint16, encrypt bool) (*P2PMessage, error) {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MagicCookie
	msg.Header.Type = uint16(msgType)
	msg.Header.NetProto = proto
	msg.Header.Length = uint16(len(payload))
	if p.Crypter.Active && encrypt {
		var err error
		msg.Data, err = p.Crypter.encrypt(p.Crypter.ActiveKey.Key, payload)
		if err != nil {
			return nil, err
		}
	} else {
		msg.Data = payload
	}
	return msg, nil
}

// CreateMessageStatic is a static method for a P2P Message
func CreateMessageStatic(msgType MsgType, payload []byte) (*P2PMessage, error) {
	p := PeerToPeer{}
	return p.CreateMessage(msgType, payload, 0, false)
}

///////////////////////////////////////////////////////////////////////////////////////////

// Network is a network subsystem
type Network struct {
	host       string
	port       int
	remotePort int
	addr       *net.UDPAddr
	conn       *net.UDPConn
	inBuffer   [4096]byte
	disposed   bool
}

// Stop will terminate packet reader
func (uc *Network) Stop() error {
	uc.disposed = true
	if uc.conn != nil {
		return uc.conn.Close()
	}
	return fmt.Errorf("Nil Connection")
}

// Disposed returns whether service is willing to stop or not
func (uc *Network) Disposed() bool {
	return uc.disposed
}

// Addr returns assigned address
func (uc *Network) Addr() *net.UDPAddr {
	if uc.addr != nil {
		return uc.addr
	}
	return nil
}

// Init creates a UDP connection
func (uc *Network) Init(host string, port int) error {
	var err error
	uc.host = host
	uc.port = port
	uc.disposed = true

	//todo check if we need Host and Port
	uc.addr, err = net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	uc.conn, err = net.ListenUDP("udp4", uc.addr)
	if err != nil {
		return err
	}
	uc.disposed = false
	return nil
}

// KeepAlive will send keep alive packet periodically to keep
// UDP port bind
func (uc *Network) KeepAlive(target string) error {
	if uc.conn == nil {
		return fmt.Errorf("Nil Connection")
	}

	addresses, err := SrvLookup(target, "udp", "subutai.io")
	if err != nil {
		return fmt.Errorf("Failed to lookup address for keep alive session: %s", err.Error())
	}
	if len(addresses) == 0 {
		return fmt.Errorf("No suitable address for keep alive")
	}

	firstAddr, e := addresses[0]
	if !e {
		return fmt.Errorf("Failed to retrieve keep alive address at index 0")
	}

	addr, err := net.ResolveUDPAddr("udp4", firstAddr)
	if err != nil {
		return fmt.Errorf("Failed to resolve UDP addr for keep alive session: %s", err.Error())
	}

	data := []byte{0x0D, 0x0A}
	keepAlive := time.Now()
	Log(Debug, "Started keep alive session with %s", addr)
	i := 0
	for i < 20 {
		uc.SendRawBytes(data, addr)
		i++
		time.Sleep(time.Millisecond * 500)
	}
	for !uc.disposed {
		if time.Duration(time.Second*3) < time.Since(keepAlive) {
			keepAlive = time.Now()
			uc.SendRawBytes(data, addr)
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

// GetPort return a port assigned
func (uc *Network) GetPort() int {
	if uc.conn == nil {
		return -1
	}
	addr, _ := net.ResolveUDPAddr("udp4", uc.conn.LocalAddr().String())
	return addr.Port
}

// UDPReceivedCallback is executed when message is received
type UDPReceivedCallback func(count int, src_addr *net.UDPAddr, err error, buff []byte)

// Listen is a main listener of a network traffic
func (uc *Network) Listen(receivedCallback UDPReceivedCallback) error {
	Log(Info, "Started UDP listener")
	if uc.conn == nil {
		return fmt.Errorf("Nil connection")
	}
	for !uc.Disposed() {
		n, src, err := uc.conn.ReadFromUDP(uc.inBuffer[:])
		receivedCallback(n, src, err, uc.inBuffer[:])
	}
	Log(Info, "Stopping UDP Listener")
	return nil
}

// SendMessage sends message over network
func (uc *Network) SendMessage(msg *P2PMessage, dstAddr *net.UDPAddr) (int, error) {
	if uc.conn == nil {
		return -1, fmt.Errorf("Nil connection")
	}
	if msg == nil {
		return 0, fmt.Errorf("Nil message")
	}
	n, err := uc.conn.WriteToUDP(msg.Serialize(), dstAddr)
	if err != nil {
		return 0, err
	}
	return n, nil
}

// SendRawBytes sends bytes over network
func (uc *Network) SendRawBytes(bytes []byte, dstAddr *net.UDPAddr) (int, error) {
	if uc.conn == nil {
		return -1, fmt.Errorf("Nil connection")
	}
	n, err := uc.conn.WriteToUDP(bytes, dstAddr)
	if err != nil {
		return 0, err
	}
	return n, nil
}
