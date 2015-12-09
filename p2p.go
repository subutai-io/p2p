package main

import (
	"fmt"
	"github.com/danderson/tuntap"
	"os"
	"os/signal"
)

func main() {
	dev, err := tuntap.Open("tuntap-t0", tuntap.DevTap)
	if err != nil {
		fmt.Errorf("Failed to open tuntap device: %v", err)
	}

	// Capture SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Println("Received signal: ", sig)
			dev.Close()
			os.Exit(0)
		}
	}()

	for {

	}

	dev.Close()

	fmt.Println("Hello, World")
}
