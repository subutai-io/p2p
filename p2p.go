package main

import (
	"flag"
	"fmt"
	"github.com/danderson/tuntap"
	"gopkg.in/yaml.v2"
	"io/ioutil"
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
	IPTool     string `yaml:"iptool"`
	Interface  *os.File
	Device     *tuntap.Interface
}

// Creates Device
func (ptp *PTPCloud) CreateDevice(ip, mac, mask, device string) *PTPCloud {
	var err error

	ptp.IP = ip
	ptp.Mac = mac
	ptp.Mask = mask
	ptp.DeviceName = device

	// Extract necessary information from config file
	// TODO: Remove hard-coded path
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		log.Printf("[ERROR] Failed to load config: %v", err)
	}
	err = yaml.Unmarshal(yamlFile, ptp)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config: %v", err)
	}

	ptp.Device, err = tuntap.Open(ptp.DeviceName, tuntap.DevTap)
	if ptp.Device == nil {
		log.Fatalf("[FATAL] Failed to open TAP device: %v", err)
	} else {
		log.Printf("[INFO] %v TAP Device created", ptp.DeviceName)
	}

	linkup := exec.Command(ptp.IPTool, "link", "set", "dev", ptp.DeviceName, "up")
	err = linkup.Run()
	if err != nil {
		log.Fatalf("Failed to up link")
	}

	// Configure new device
	log.Printf("[INFO] Setting %s IP on device %s\n", ptp.IP, ptp.DeviceName)
	setip := exec.Command(ptp.IPTool, "addr", "add", ptp.IP, "dev", ptp.DeviceName)
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

	// TODO: Improve this
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
