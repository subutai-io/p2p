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
	"p2p/dht"
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
	setip := exec.Command(ptp.IPTool, "addr", "add", ptp.IP+"/24", "dev", ptp.DeviceName)
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
	var argDirect string
	var argHash string

	// TODO: Improve this
	flag.StringVar(&argIp, "ip", "none", "IP Address to be used")
	// TODO: Parse this properly
	flag.StringVar(&argMask, "mask", "none", "Network mask")
	// TODO: Implement this
	flag.StringVar(&argMac, "mac", "none", "MAC Address for a TUN/TAP interface")
	flag.StringVar(&argDev, "dev", "none", "TUN/TAP interface name")
	flag.StringVar(&argDirect, "direct", "none", "IP to connect to directly")
	flag.StringVar(&argHash, "hash", "none", "Infohash")

	flag.Parse()
	if argIp == "none" || argMask == "none" || argDev == "none" {
		fmt.Println("USAGE: p2p [OPTIONS]")
		fmt.Printf("\nOPTIONS:\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var dhtClient dht.DHTClient
	config := dhtClient.DHTClientConfig()
	config.NetworkHash = argHash
	dhtClient.Initialize(config)

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
		packet, err := ptp.Device.ReadPacket()
		if err != nil {
			log.Printf("Error reading packet: %s", err)
		}
		//log.Printf("Packet received: %s", string(packet.Packet))
		log.Printf("Packet received: %d", string(packet.Protocol))

	}
}
