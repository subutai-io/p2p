package main

import (
	"flag"
	"fmt"
	"github.com/danderson/tuntap"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

type PTPCloud struct {
	IP         string
	Mac        string
	Mask       string
	DeviceName string
	Device     *tuntap.Interface
}

func (ptp *PTPCloud) CreateDevice(ip, mac, mask, device string) *PTPCloud {
	ptp.IP = ip
	ptp.Mac = mac
	ptp.Mask = mask
	ptp.DeviceName = device

	var err error

	ptp.Device, err = tuntap.Open(ptp.DeviceName, tuntap.DevTap)
	if ptp.Device == nil {
		log.Fatalf("[FATAL] Failed to open TAP device: %v", err)
	} else {
		log.Printf("[INFO] %v TAP Device created", ptp.DeviceName)
	}

	// Configure new device
	log.Printf("[INFO] Setting %s IP on device %s\n", ptp.IP, ptp.DeviceName)
	setip := exec.Command("/usr/bin/ip", "addr", "add", ptp.IP, "dev", ptp.DeviceName)
	err = setip.Run()
	if err != nil {
		log.Fatalf("[FATAL] Failed to set IP: %v", err)
	}
	return ptp
}

func main() {
	var argIp string
	var argMask string
	var argMac string
	var argDev string

	flag.StringVar(&argIp, "ip", "none", "IP Address to be used")
	flag.StringVar(&argMask, "mask", "none", "IP Address to be used")
	flag.StringVar(&argMac, "mac", "none", "IP Address to be used")
	flag.StringVar(&argDev, "dev", "none", "IP Address to be used")

	flag.Parse()
	if argIp == "none" || argMask == "none" || argDev == "none" {
		fmt.Println("USAGE: p2p [OPTIONS]")
		fmt.Printf("\nOPTIONS:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var ptp PTPCloud
	ptp.CreateDevice(argIp, argMac, argMask, argDev)

	// Capture SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Println("Received signal: ", sig)
			os.Exit(0)
		}
	}()

	for {

	}
}
