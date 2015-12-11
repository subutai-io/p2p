package main

import (
	"fmt"
	"net"
	//"strconv"
	"bytes"
	bencode "github.com/jackpal/bencode-go"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Println("Error: ", err)
	}
}

type OutgoingData struct {
	command string "c"
	id      string "i"
	hash    string "h"
}

func main() {
	ServerAddr, err := net.ResolveUDPAddr("udp", "dht2.subut.ai:6881")
	CheckError(err)

	LocalAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	CheckError(err)

	Conn, err := net.DialUDP("udp", LocalAddr, ServerAddr)
	CheckError(err)

	defer Conn.Close()
	for {
		// Send
		var b bytes.Buffer
		var query OutgoingData
		query.command = "Hello"
		query.id = "asdf"
		if err := bencode.Marshal(&b, query); err != nil {
			fmt.Printf("Failed to Marshal\n")
		}
		msg := b.String()
		//i++
		//buf := []byte(msg)
		_, err = Conn.Write([]byte(msg))
		if err != nil {
			fmt.Println(msg, err)
		}
		time.Sleep(time.Second * 1)
		// Receive
		var rbuf [512]byte
		_, addr, err := Conn.ReadFromUDP(rbuf[0:])
		CheckError(err)
		fmt.Printf("Receiced from %s: %s\n", addr, string(rbuf[:10]))
	}
}
