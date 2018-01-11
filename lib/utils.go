package ptp

import (
	"crypto/rand"
	"fmt"
	"net"
)

// Different utility functions

// GenerateMAC generates a MAC address for a new interface
func GenerateMAC() (string, net.HardwareAddr) {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		Log(Error, "Failed to generate MAC: %v", err)
		return "", nil
	}
	buf[0] |= 2
	mac := fmt.Sprintf("06:%02x:%02x:%02x:%02x:%02x", buf[1], buf[2], buf[3], buf[4], buf[5])
	hw, err := net.ParseMAC(mac)
	if err != nil {
		Log(Error, "Corrupted MAC address generated: %v", err)
		return "", nil
	}
	return mac, hw
}

// This method compares given IP to known private IP address spaces
// and return true if IP is private, false otherwise
func isPrivateIP(ip net.IP) (bool, error) {
	if ip == nil {
		return false, fmt.Errorf("Missing IP")
	}
	_, private24, _ := net.ParseCIDR("10.0.0.0/8")
	_, private20, _ := net.ParseCIDR("172.16.0.0/12")
	_, private16, _ := net.ParseCIDR("192.168.0.0/16")
	isPrivate := private24.Contains(ip) || private20.Contains(ip) || private16.Contains(ip)
	return isPrivate, nil
}
