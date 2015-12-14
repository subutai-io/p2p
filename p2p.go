package main

import (
	"flag"
	"fmt"
	"github.com/danderson/tuntap"
	//"golang.org/x/net/icmp"
	//"golang.org/x/net/ipv4"
	//"encoding/binary"
	//"encoding/hex"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"p2p/dht"
	//"strings"
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
func (ptp *PTPCloud) CreateDevice(ip, mac, mask, device string) (*PTPCloud, error) {
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
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, ptp)
	if err != nil {
		log.Printf("[ERROR] Failed to parse config: %v", err)
		return nil, err
	}

	ptp.Device, err = tuntap.Open(ptp.DeviceName, tuntap.DevTap)
	if ptp.Device == nil {
		log.Fatalf("[FATAL] Failed to open TAP device: %v", err)
		return nil, err
	} else {
		log.Printf("[INFO] %v TAP Device created", ptp.DeviceName)
	}

	linkup := exec.Command(ptp.IPTool, "link", "set", "dev", ptp.DeviceName, "up")
	err = linkup.Run()
	if err != nil {
		log.Fatalf("[ERROR] Failed to up link: %v", err)
		return nil, err
	}

	// Configure new device
	log.Printf("[INFO] Setting %s IP on device %s\n", ptp.IP, ptp.DeviceName)
	setip := exec.Command(ptp.IPTool, "addr", "add", ptp.IP+"/24", "dev", ptp.DeviceName)
	err = setip.Run()
	if err != nil {
		log.Fatalf("[FATAL] Failed to set IP: %v", err)
		return nil, err
	}
	return ptp, nil
}

func (ptp *PTPCloud) handlePacket(contents []byte, proto int) {
	/*
		512   (PUP)
		2048  (IP)
		2054  (ARP)
		32821 (RARP)
		33024 (802.1q)
		34525 (IPv6)
		34915 (PPPOE discovery)
		34916 (PPPOE session)
	*/
	switch proto {
	case 512:
		log.Printf("[DEBUG] Received PARC Universal Packet")
	case 2048:
		log.Printf("[DEBUG] Received IPv4 Packet")
		ptp.handlePacketIPv4(contents)
	case 2054:
		log.Printf("[DEBUG] Received ARP Packet")
		ptp.handlePacketARP(contents)
	case 32821:
		log.Printf("[DEBUG] Received RARP Packet")
	case 33024:
		log.Printf("[DEBUG] Received 802.1q Packet")
	case 34525:
		log.Printf("[DEBUG] Received IPv6 Packet")
	case 34915:
		log.Printf("[DEBUG] Received PPPoE Discovery Packet")
	case 34916:
		log.Printf("[DEBUG] Received PPPoE Session Packet")
	default:
		log.Printf("[DEBUG] Received Undefined Packet")
	}
	return

	// Here our TUN device received a packet. Let's parse it

	//log.Printf("Packet received: %s", string(packet.Packet))
	/*
		log.Printf("Packet received: %d", string(packet.Protocol))
		header, err := ipv4.ParseHeader(packet.Packet)
		if packet.Truncated {
			log.Printf("[DEBUG] Truncated packet")
		}
	*/
	// ICMP
	/*
		ipacket, err := icmp.ParseIPv4Header(contents)

		if err != nil {
			log.Printf("[ERROR] Failed to parse IPv4 packet: %v", err)
		} else {
			dstParts := strings.Split(ipacket.Dst.String(), ".")
			if dstParts[0] == "0" {
				continue
			}
			log.Printf("[DEBUG] ICMP: %v", ipacket.String())
		}
	*/

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
	// TODO: Direct connection is not implemented yet
	flag.StringVar(&argDirect, "direct", "none", "IP to connect to directly")
	flag.StringVar(&argHash, "hash", "none", "Infohash for environment")

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
		if packet.Truncated {
			log.Printf("[DEBUG] Truncated packet")
		}
		go ptp.handlePacket(packet.Packet, packet.Protocol)
		//else {
		//	log.Printf("[DEBUG] Destination: %s", header.Dst.String())
		//}
	}
}

// WriteToDevice writes data to created TUN/TAP device
func (ptp *PTPCloud) WriteToDevice(b []byte) {
	var p *tuntap.Packet
	p.Protocol = 2054
	p.Truncated = false
	p.Packet = b
	ptp.Device.WritePacket(p)
}
