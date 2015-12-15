package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const MAGIC_COOKIE uint16 = 0xabcd

type MSG_TYPE uint16

const (
	MT_STRING MSG_TYPE = 0
	//todo add types
)

type p2p_msg_header_t struct {
	Magic    uint16
	Type     uint16
	Length   uint16
	NetProto uint16
}

func (v *p2p_msg_header_t) serialize() []byte {
	res_buf := make([]byte, 8)
	binary.BigEndian.PutUint16(res_buf[0:2], v.Magic)
	binary.BigEndian.PutUint16(res_buf[2:4], v.Type)
	binary.BigEndian.PutUint16(res_buf[4:6], v.Length)
	binary.BigEndian.PutUint16(res_buf[6:8], v.NetProto)
	return res_buf
}

func p2p_msg_header_from_bytes(bytes []byte) (*p2p_msg_header_t, error) {
	if len(bytes) < 8 {
		return nil, errors.New("p2p_msg_header_from_bytes_error : less then 6 bytes")
	}

	result := new(p2p_msg_header_t)
	result.Magic = binary.BigEndian.Uint16(bytes[0:2])
	result.Type = binary.BigEndian.Uint16(bytes[2:4])
	result.Length = binary.BigEndian.Uint16(bytes[4:6])
	result.NetProto = binary.BigEndian.Uint16(bytes[6:8])
	return result, nil
}

type p2p_msg_t struct {
	Header *p2p_msg_header_t
	Data   []byte
}

func (v *p2p_msg_t) serialize() []byte {
	v.Header.Length = uint16(len(v.Data))
	res_buf := v.Header.serialize()
	res_buf = append(res_buf, v.Data...)
	return res_buf
}

func p2p_msg_from_bytes(bytes []byte) (*p2p_msg_t, error) {
	res := new(p2p_msg_t)
	var err error = nil
	res.Header, err = p2p_msg_header_from_bytes(bytes)
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

func create_string_p2p_msg(data string, netProto uint16) *p2p_msg_t {
	msg := new(p2p_msg_t)
	msg.Header = new(p2p_msg_header_t)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(MT_STRING)
	msg.Header.NetProto = netProto
	msg.Data = []byte(data)
	return msg
}

///////////////////////////////////////////////////////////////////////////////////////////

type udp_client_t struct {
	host         string
	port         int16
	addr         *net.UDPAddr
	conn         *net.UDPConn
	input_buffer [1024]byte
	disposed     bool
}

func (uc *udp_client_t) Disposed() bool {
	return uc.disposed
}

func (uc *udp_client_t) Addr() *net.UDPAddr {
	return uc.addr
}

func (uc *udp_client_t) init(host string, port int16) error {
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

type udp_received_callback func(count int, src_addr *net.UDPAddr, err error, buff []byte)

func (uc *udp_client_t) listen(fn_received_callback udp_received_callback) {
	for !uc.Disposed() {
		n, src, err := uc.conn.ReadFromUDP(uc.input_buffer[:])
		fn_received_callback(n, src, err, uc.input_buffer[:])
	}
}

func (uc *udp_client_t) send_msg(msg *p2p_msg_t, dst_addr *net.UDPAddr) {
	ser_data := msg.serialize()
	n, err := uc.conn.WriteToUDP(ser_data, dst_addr)
	if err != nil {
		fmt.Printf("error sending msg : %v\n", err)
		return
	}
	fmt.Printf("sent %d bytes\n", n)
}

///////////////////////////////////////////////////////////////////////////////////////////

func main() {
	var udp_client_0 *udp_client_t = new(udp_client_t)
	var udp_client_1 *udp_client_t = new(udp_client_t)

	udp_client_0.init("", 5000)
	udp_client_1.init("", 5001)

	go udp_client_0.listen(process_p2p_msg)
	go udp_client_1.listen(process_p2p_msg)

	msg := create_string_p2p_msg("Hello, world!", 0)
	udp_client_0.send_msg(msg, udp_client_1.Addr())

	for {
		time.Sleep(100 * time.Millisecond)
	}
}

func process_p2p_msg(count int, src_addr *net.UDPAddr, err error, rcv_bytes []byte) {
	if err != nil {
		fmt.Printf("process_p2p_msg error : %v\n", err)
		return
	}

	buf := make([]byte, count)
	copy(buf[:], rcv_bytes[:])

	msg, des_err := p2p_msg_from_bytes(buf)
	if des_err != nil {
		fmt.Printf("p2p_msg_from_bytes err : %v\n", des_err)
		return
	}

	fmt.Printf("processed message from %s, msg_data : %s\n", src_addr.String(), msg.Data)
}
