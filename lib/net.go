package ptp

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
)

const (
	MAGIC_COOKIE uint16 = 0xabcd
	HEADER_SIZE  int    = 18
)

type P2PMessageHeader struct {
	Magic         uint16
	Type          uint16
	Length        uint16
	NetProto      uint16
	ProxyId       uint16
	SerializedLen uint16
	Complete      uint16
	Id            uint16
	Seq           uint16
}

type P2PMessage struct {
	Header *P2PMessageHeader
	Data   []byte
}

func (v *P2PMessageHeader) Serialize() []byte {
	res_buf := make([]byte, HEADER_SIZE)
	binary.BigEndian.PutUint16(res_buf[0:2], v.Magic)
	binary.BigEndian.PutUint16(res_buf[2:4], v.Type)
	binary.BigEndian.PutUint16(res_buf[4:6], v.Length)
	binary.BigEndian.PutUint16(res_buf[6:8], v.NetProto)
	binary.BigEndian.PutUint16(res_buf[8:10], v.ProxyId)
	binary.BigEndian.PutUint16(res_buf[10:12], v.SerializedLen)
	binary.BigEndian.PutUint16(res_buf[12:14], v.Complete)
	binary.BigEndian.PutUint16(res_buf[14:16], v.Id)
	binary.BigEndian.PutUint16(res_buf[16:18], v.Seq)
	return res_buf
}

func P2PMessageHeaderFromBytes(bytes []byte) (*P2PMessageHeader, error) {
	if len(bytes) < HEADER_SIZE {
		return nil, errors.New("P2PMessageHeaderFromBytes_error : less then 14 bytes")
	}

	result := new(P2PMessageHeader)
	result.Magic = binary.BigEndian.Uint16(bytes[0:2])
	result.Type = binary.BigEndian.Uint16(bytes[2:4])
	result.Length = binary.BigEndian.Uint16(bytes[4:6])
	result.NetProto = binary.BigEndian.Uint16(bytes[6:8])
	result.ProxyId = binary.BigEndian.Uint16(bytes[8:10])
	result.SerializedLen = binary.BigEndian.Uint16(bytes[10:12])
	result.Complete = binary.BigEndian.Uint16(bytes[12:14])
	result.Id = binary.BigEndian.Uint16(bytes[14:16])
	result.Seq = binary.BigEndian.Uint16(bytes[16:18])
	return result, nil
}

func (v *P2PMessage) Serialize() []byte {
	v.Header.SerializedLen = uint16(len(v.Data))
	Log(TRACE, "--- Serialize P2PMessage header.SerializedLen : %d", v.Header.SerializedLen)
	res_buf := v.Header.Serialize()
	res_buf = append(res_buf, v.Data...)
	return res_buf
}

func P2PMessageFromBytes(bytes []byte) (*P2PMessage, error) {
	res := new(P2PMessage)
	var err error = nil
	res.Header, err = P2PMessageHeaderFromBytes(bytes)
	if err != nil {
		return nil, err
	}
	Log(TRACE, "--- P2PMessageHeaderFromBytes Length : %d, SerLen : %d", res.Header.Length, res.Header.SerializedLen)
	if res.Header.Magic != MAGIC_COOKIE {
		return nil, errors.New("magic cookie not presented")
	}
	res.Data = make([]byte, res.Header.SerializedLen)
	Log(TRACE, "BYTES : %s", bytes)
	copy(res.Data[:], bytes[HEADER_SIZE:len(bytes)])
	Log(TRACE, "res.Data : %s", res.Data)
	return res, err
}

func CreateStringP2PMessage(c Crypto, data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_STRING)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.Complete = 1
	msg.Header.Id = 1
	msg.Header.Seq = 1
	if c.Active {
		var err error
		msg.Data, err = c.Encrypt(c.ActiveKey.Key, []byte(data))
		if err != nil {
			Log(ERROR, "Failed to encrypt data")
		}
	} else {
		msg.Data = []byte(data)
	}
	return msg
}

func CreatePingP2PMessage() *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_PING)
	msg.Header.NetProto = 0
	msg.Header.Length = uint16(len("1"))
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	msg.Data = []byte("1")
	return msg
}

func CreateXpeerPingMessage(pt PingType, hw string) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_XPEER_PING)
	msg.Header.NetProto = uint16(pt)
	msg.Header.Length = uint16(len(hw))
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	msg.Data = []byte(hw)
	return msg
}

func CreateIntroP2PMessage(c Crypto, data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_INTRO)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	if c.Active {
		var err error
		msg.Data, err = c.Encrypt(c.ActiveKey.Key, []byte(data))
		if err != nil {
			Log(ERROR, "Failed to encrypt data")
		}
	} else {
		msg.Data = []byte(data)
	}
	return msg
}

func CreateIntroRequest(c Crypto, id string) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_INTRO_REQ)
	msg.Header.NetProto = 0
	msg.Header.Length = uint16(len(id))
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	if c.Active {
		var err error
		msg.Data, err = c.Encrypt(c.ActiveKey.Key, []byte(id))
		if err != nil {
			Log(ERROR, "Failed to encrypt data")
		}
	} else {
		msg.Data = []byte(id)
	}
	return msg
}

func CreateNencP2PMessage(c Crypto, data []byte, netProto, complete, id, seq uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_NENC)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.Complete = complete
	msg.Header.Id = id
	msg.Header.Seq = seq
	if c.Active {
		var err error
		msg.Data, err = c.Encrypt(c.ActiveKey.Key, data)
		if err != nil {
			Log(ERROR, "Failed to encrypt data")
		}
	} else {
		msg.Data = data
	}
	return msg
}

func CreateTestP2PMessage(c Crypto, data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_TEST)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	if c.Active {
		var err error
		msg.Data, err = c.Encrypt(c.ActiveKey.Key, []byte(data))
		if err != nil {
			Log(ERROR, "Failed to encrypt data")
		}
	} else {
		msg.Data = []byte(data)
	}
	return msg
}

func CreateProxyP2PMessage(id int, data string, netProto uint16) *P2PMessage {
	// We don't need to encrypt this message
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_PROXY)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.Complete = 1
	msg.Header.ProxyId = uint16(id)
	msg.Header.Id = 0
	msg.Header.Seq = 0
	msg.Data = []byte(data)
	return msg
}

func CreateBadTunnelP2PMessage(id int, netProto uint16) *P2PMessage {
	data := "rem"
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_BAD_TUN)
	msg.Header.NetProto = netProto
	msg.Header.Length = uint16(len(data))
	msg.Header.ProxyId = uint16(id)
	msg.Header.Complete = 1
	msg.Header.Id = 0
	msg.Header.Seq = 0
	msg.Data = []byte(data)
	return msg
}

///////////////////////////////////////////////////////////////////////////////////////////

type PTPNet struct {
	host         string
	port         int
	addr         *net.UDPAddr
	conn         *net.UDPConn
	input_buffer [4096]byte
	disposed     bool
}

func (uc *PTPNet) Stop() {
	uc.disposed = true
}

func (uc *PTPNet) Disposed() bool {
	return uc.disposed
}

func (uc *PTPNet) Addr() *net.UDPAddr {
	return uc.addr
}

func (uc *PTPNet) Init(host string, port int) error {
	var err error = nil
	uc.host = host
	uc.port = port
	uc.disposed = true

	//todo check if we need Host and Port
	uc.addr, err = net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	uc.conn, err = net.ListenUDP("udp", uc.addr)
	if err != nil {
		return err
	}
	uc.disposed = false
	return nil
}

func (uc *PTPNet) GetPort() int {
	addr, _ := net.ResolveUDPAddr("udp", uc.conn.LocalAddr().String())
	return addr.Port
}

type UDPReceivedCallback func(count int, src_addr *net.UDPAddr, err error, buff []byte)

func (uc *PTPNet) Listen(fn_received_callback UDPReceivedCallback) {
	for !uc.Disposed() {
		n, src, err := uc.conn.ReadFromUDP(uc.input_buffer[:])
		fn_received_callback(n, src, err, uc.input_buffer[:])
	}
	Log(INFO, "Stopping UDP Listener")
}

func (uc *PTPNet) Bind(addr *net.UDPAddr, local_addr *net.UDPAddr) {

}

func (uc *PTPNet) SendMessage(msg *P2PMessage, dst_addr *net.UDPAddr) (int, error) {
	ser_data := msg.Serialize()
	n, err := uc.conn.WriteToUDP(ser_data, dst_addr)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func Process_p2p_msg(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		fmt.Printf("process_p2p_msg error : %v\n", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := P2PMessageFromBytes(buf)
	if des_err != nil {
		fmt.Printf("P2PMessageFromBytes err : %v\n", des_err)
		return
	}

	fmt.Printf("processed message from %s, msg_count %d, msg_data : %s\n",
		src_addr.String(),
		count,
		msg.Data)
}
