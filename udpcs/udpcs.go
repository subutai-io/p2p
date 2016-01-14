package udpcs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net"
	"p2p/commons"
	"time"
)

const (
	MAGIC_COOKIE uint16 = 0xabcd
	BLOCK_SIZE   int    = 32
	IV_SIZE      int    = aes.BlockSize
)

type CryptoKey struct {
	TTLConfig string `yaml:"ttl"`
	KeyConfig string `yaml:"key"`
	Until     time.Time
	Key       []byte
}

type Crypto struct {
	Keys   []CryptoKey
	Active bool
}

type P2PMessageHeader struct {
	Magic    uint16
	Type     uint16
	Length   uint16
	NetProto uint16
	ProxyId  uint16
}

func (v *P2PMessageHeader) Serialize() []byte {
	res_buf := make([]byte, 10)
	binary.BigEndian.PutUint16(res_buf[0:2], v.Magic)
	binary.BigEndian.PutUint16(res_buf[2:4], v.Type)
	binary.BigEndian.PutUint16(res_buf[4:6], v.Length)
	binary.BigEndian.PutUint16(res_buf[6:8], v.NetProto)
	binary.BigEndian.PutUint16(res_buf[8:10], v.ProxyId)
	return res_buf
}

func P2PMessageHeaderFromBytes(bytes []byte) (*P2PMessageHeader, error) {
	if len(bytes) < 10 {
		return nil, errors.New("P2PMessageHeaderFromBytes_error : less then 6 bytes")
	}

	result := new(P2PMessageHeader)
	result.Magic = binary.BigEndian.Uint16(bytes[0:2])
	result.Type = binary.BigEndian.Uint16(bytes[2:4])
	result.Length = binary.BigEndian.Uint16(bytes[4:6])
	result.NetProto = binary.BigEndian.Uint16(bytes[6:8])
	result.ProxyId = binary.BigEndian.Uint16(bytes[8:10])
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
	copy(res.Data[:], bytes[10:len(bytes)])
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

func CreateTestP2PMessage(c *Crypto, data string, netProto uint16) *P2PMessage {
	msg := new(P2PMessage)
	msg.Header = new(P2PMessageHeader)
	msg.Header.Magic = MAGIC_COOKIE
	msg.Header.Type = uint16(commons.MT_TEST)
	msg.Header.NetProto = netProto
	msg.Data = []byte(data)
	return msg
}

///////////////////////////////////////////////////////////////////////////////////////////

type UDPClient struct {
	host         string
	port         int16
	addr         *net.UDPAddr
	conn         *net.UDPConn
	input_buffer [4096]byte
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

	fmt.Printf("processed message from %s, msg_count %d, msg_data : %s\n",
		src_addr.String(),
		count,
		msg.Data)
}

// Cryptography

func (c *Crypto) EncrichKeyValues(ckey CryptoKey, key, datetime string) CryptoKey {
	var err error
	ckey.Until, err = time.Parse("2016-01-14 01:18:18.032507415 +0600 KGT", datetime)
	ckey.Key = []byte(key)
	if err != nil {
		log.Printf("[ERROR] Failed to parse provided TTL value: %v", err)
		return ckey
	}
	return ckey
}

func (c *Crypto) ReadKeysFromFile(filepath string) {
	yamlFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Printf("[ERROR] Failed to read key file yaml: %v", err)
		c.Active = false
		return
	}
	var ckey CryptoKey
	err = yaml.Unmarshal(yamlFile, ckey)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config: %v", err)
		c.Active = false
		return
	}
	ckey = c.EncrichKeyValues(ckey, ckey.KeyConfig, ckey.TTLConfig)
	c.Active = true
	c.Keys = append(c.Keys, ckey)
}

func (c *Crypto) Encrypt(key []byte, data []byte) ([]byte, error) {
	cb, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	data_len := len(data)
	result_data_len := (data_len + BLOCK_SIZE - 1) & (^(BLOCK_SIZE - 1))
	encrypted_data := make([]byte, IV_SIZE+result_data_len)
	// The IV needs to be unique, but not secured.
	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	copy(encrypted_data[:IV_SIZE], iv)
	count := result_data_len / BLOCK_SIZE
	for i := 0; i < count-1; i++ {
		mode := cipher.NewCBCEncrypter(cb, iv)
		mode.CryptBlocks(encrypted_data[i*BLOCK_SIZE+IV_SIZE:], data[i*BLOCK_SIZE:(i+1)*BLOCK_SIZE])
	}

	tmp_arr := make([]byte, BLOCK_SIZE)
	copy(tmp_arr, data[(count-1)*BLOCK_SIZE:])
	mode := cipher.NewCBCEncrypter(cb, iv)
	mode.CryptBlocks(encrypted_data[(count-1)*BLOCK_SIZE+IV_SIZE:], tmp_arr)
	return encrypted_data, nil
}

/////////////////////////////////////////////////////

func (c *Crypto) Decrypt(key []byte, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	iv := data[:IV_SIZE]
	data_len := len(data) - IV_SIZE
	decrypted_data := make([]byte, data_len)
	count := data_len / BLOCK_SIZE
	for i := 0; i < count; i++ {
		mode := cipher.NewCBCDecrypter(block, iv)
		mode.CryptBlocks(decrypted_data[i*BLOCK_SIZE:], data[i*BLOCK_SIZE+IV_SIZE:])
	}
	return decrypted_data, nil
}
