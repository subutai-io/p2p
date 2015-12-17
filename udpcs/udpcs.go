package udpcs

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"net"
	"p2p/commons"
	"time"
)

const MAGIC_COOKIE uint16 = 0xabcd

type P2PMessageHeader struct {
	Magic    uint16
	Type     uint16
	Length   uint16
	NetProto uint16
}

func (v *P2PMessageHeader) Serialize() []byte {
	res_buf := make([]byte, 8)
	binary.BigEndian.PutUint16(res_buf[0:2], v.Magic)
	binary.BigEndian.PutUint16(res_buf[2:4], v.Type)
	binary.BigEndian.PutUint16(res_buf[4:6], v.Length)
	binary.BigEndian.PutUint16(res_buf[6:8], v.NetProto)
	return res_buf
}

func P2PMessageHeaderFromBytes(bytes []byte) (*P2PMessageHeader, error) {
	if len(bytes) < 8 {
		return nil, errors.New("P2PMessageHeaderFromBytes_error : less then 6 bytes")
	}

	result := new(P2PMessageHeader)
	result.Magic = binary.BigEndian.Uint16(bytes[0:2])
	result.Type = binary.BigEndian.Uint16(bytes[2:4])
	result.Length = binary.BigEndian.Uint16(bytes[4:6])
	result.NetProto = binary.BigEndian.Uint16(bytes[6:8])
	return result, nil
}

type P2PMessage struct {
	Header *P2PMessageHeader
	Data   []byte
}

func (v *P2PMessage) Serialize() []byte {
	v.Header.Length = uint16(len(v.Data))
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
	if res.Header.Magic != MAGIC_COOKIE {
		return nil, errors.New("magic cookie not presented")
	}
	res.Data = make([]byte, res.Header.Length)
	copy(res.Data[:], bytes[8:len(bytes)])
	return res, err
}

func CreateStringP2PMessage(data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(commons.MT_STRING)
	msg.Header.NetProto = netProto
	msg.Data = []byte(data)
	return msg
}

func CreateIntroP2PMessage(data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(commons.MT_INTRO)
	msg.Header.NetProto = netProto
	msg.Data = []byte(data)
	return msg
}

func CreateNencP2PMessage(data []byte, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(commons.MT_NENC)
	msg.Header.NetProto = netProto
	msg.Data = data
	return msg
}

///////////////////////////////////////////////////////////////////////////////////////////

type UDPClient struct {
	host         string
	port         int16
	addr         *net.UDPAddr
	conn         *net.UDPConn
	input_buffer [1024]byte
	disposed     bool
}

func (uc *UDPClient) Disposed() bool {
	return uc.disposed
}

func (uc *UDPClient) Addr() *net.UDPAddr {
	return uc.addr
}

func (uc *UDPClient) Init(host string, port int16) error {
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

func (uc *UDPClient) GetPort() int {
	addr, _ := net.ResolveUDPAddr("udp", uc.conn.LocalAddr().String())
	return addr.Port
}

type UDPReceivedCallback func(count int, src_addr *net.UDPAddr, err error, buff []byte)

func (uc *UDPClient) Listen(fn_received_callback UDPReceivedCallback) {
	for !uc.Disposed() {
		n, src, err := uc.conn.ReadFromUDP(uc.input_buffer[:])
		log.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!")
		fn_received_callback(n, src, err, uc.input_buffer[:])
	}
}

func (uc *UDPClient) SendMessage(msg *P2PMessage, dst_addr *net.UDPAddr) (int, error) {
	ser_data := msg.Serialize()
	n, err := uc.conn.WriteToUDP(ser_data, dst_addr)
	if err != nil {
		return 0, err
	}
	return n, nil
}

///////////////////////////////////////////////////////////////////////////////////////////

func UDPCSTest() {
	var udp_client_0 *UDPClient = new(UDPClient)
	var udp_client_1 *UDPClient = new(UDPClient)

	udp_client_0.Init("", 5000)
	udp_client_1.Init("", 5001)

	go udp_client_0.Listen(Process_p2p_msg)
	go udp_client_1.Listen(Process_p2p_msg)

	msg := CreateStringP2PMessage("Hello, world!", 0)
	udp_client_0.SendMessage(msg, udp_client_1.Addr())

	for {
		time.Sleep(100 * time.Millisecond)
	}
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

	fmt.Printf("processed message from %s, msg_data : %s\n", src_addr.String(), msg.Data)
}
