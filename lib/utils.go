package ptp

import (
	"crypto/rand"
	"fmt"
	"net"

	uuid "github.com/wayn3h0/go-uuid"
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

// GenerateToken produces UUID string that will be used during handshake
// with DHT server. Since we don't have an ID on start - we will use token
// and wait from DHT server to respond with ID and our Token, so later
// we will replace Token with received ID
func GenerateToken() string {
	result := ""
	id, err := uuid.NewTimeBased()
	if err != nil {
		Log(Error, "Failed to generate token for peer")
		return result
	}
	result = id.String()
	Log(Debug, "Token generated: %s", result)
	return result
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

// StringifyState extracts human-readable word that represents a peer status
func StringifyState(state PeerState) string {
	switch state {
	case PeerStateInit:
		return "Initializing"
	case PeerStateRequestedIP:
		return "Waiting for IP"
	case PeerStateRequestingProxy:
		return "Requesting proxies"
	case PeerStateWaitingForProxy:
		return "Waiting for proxies"
	case PeerStateWaitingToConnect:
		return "Waiting for connection"
	case PeerStateConnecting:
		return "Initializing connection"
	case PeerStateConnected:
		return "Connected"
	case PeerStateDisconnect:
		return "Disconnected"
	case PeerStateStop:
		return "Stopped"
	case PeerStateCooldown:
		return "Cooldown"
	}
	return "Unknown"
}

// IsInterfaceLocal will return true if specified IP is in list of
// local network interfaces
func IsInterfaceLocal(ip net.IP) bool {
	for _, localIP := range ActiveInterfaces {
		if localIP.Equal(ip) {
			return true
		}
	}
	return false
}

// FindNetworkAddresses method lists interfaces available in the system and retrieves their
// IP addresses
func (p *PeerToPeer) FindNetworkAddresses() {
	Log(Debug, "Looking for available network interfaces")
	inf, err := net.Interfaces()
	if err != nil {
		Log(Error, "Failed to retrieve list of network interfaces")
		return
	}
	p.LocalIPs = p.LocalIPs[:0]
	for _, i := range inf {
		addresses, err := i.Addrs()

		if err != nil {
			Log(Error, "Failed to retrieve address for interface. %v", err)
			continue
		}
		for _, addr := range addresses {
			var decision = "Ignoring"
			var ipType = "Unknown"
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				Log(Error, "Failed to parse CIDR notation: %v", err)
			}
			if ip.IsLoopback() {
				ipType = "Loopback"
			} else if ip.IsMulticast() {
				ipType = "Multicast"
			} else if ip.IsGlobalUnicast() {
				decision = "Saving"
				ipType = "Global Unicast"
			} else if ip.IsLinkLocalUnicast() {
				ipType = "Link Local Unicast"
			} else if ip.IsLinkLocalMulticast() {
				ipType = "Link Local Multicast"
			} else if ip.IsInterfaceLocalMulticast() {
				ipType = "Interface Local Multicast"
			}
			if !p.IsIPv4(ip.String()) {
				decision = "No IPv4"
			}
			Log(Trace, "Interface %s: %s. Type: %s. %s", i.Name, addr.String(), ipType, decision)
			if decision == "Saving" && FilterInterface(i.Name, addr.String()) {
				p.LocalIPs = append(p.LocalIPs, ip)
			}
		}
	}
	Log(Trace, "%d interfaces were saved", len(p.LocalIPs))
}
