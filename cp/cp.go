package main

import (
	"log"
	"net"
)

func handleConnection(c *net.Conn) int {
	return 1

}

func main() {
	// Start a UDP server listener
	listen, err := net.Listen("udp", ":15223")
	if err != nil {
		log.Fatalf("Failed to start UDP server at specified port")
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			// handle error
		}
		go handleConnection(conn)
	}

}
