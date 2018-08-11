package ptp

import (
	"crypto/rand"
	"fmt"
	"net"

	"github.com/google/uuid"
)

// Different utility functions

// GenerateMAC will generate a new MAC address
// Function returns both string representation of MAC and it golang equalient.
// First octet is always 06
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
	id, err := uuid.NewUUID()
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
		return "INITIALIZING"
	case PeerStateRequestedIP:
		return "WAITING_IP"
	case PeerStateRequestingProxy:
		return "REQUESTING_PROXIES"
	case PeerStateWaitingForProxy:
		return "WAITING_PROXIES"
	case PeerStateWaitingToConnect:
		return "WAITING_CONNECTION"
	case PeerStateConnecting:
		return "INITIALIZING_CONNECTION"
	case PeerStateConnected:
		return "CONNECTED"
	case PeerStateDisconnect:
		return "DISCONNECTED"
	case PeerStateStop:
		return "STOPPED"
	case PeerStateCooldown:
		return "COOLDOWN"
	}
	return "UNKNOWN"
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
func (p *PeerToPeer) FindNetworkAddresses() error {
	Log(Debug, "Looking for available network interfaces")
	interfaces, err := net.Interfaces()
	if err != nil {
		Log(Error, "Failed to retrieve list of network interfaces: %s", err.Error())
		return fmt.Errorf("Failed to retrieve list of network interfaces: %s", err.Error())
	}
	p.LocalIPs = p.LocalIPs[:0]
	p.LocalIPs = p.ParseInterfaces(interfaces)
	Log(Trace, "%d interfaces were saved", len(p.LocalIPs))
	return nil
}

// ParseInterfaces accepts list of network interfaces (net.Interface),
// parse their addresses, check they CIDRs and cast type.
// Returns list of IPs
func (p *PeerToPeer) ParseInterfaces(interfaces []net.Interface) []net.IP {
	ips := []net.IP{}
	// We use reserve to collect all multicast interfaces and use them as a fallback
	// in a case when we don't find any interfaces with outbound traffic enabled
	reserve := []net.IP{}
	for _, i := range interfaces {
		addresses, err := i.Addrs()
		if err != nil {
			Log(Error, "Failed to retrieve address for interface: %s", err.Error())
			continue
		}
		if len(addresses) == 0 {
			Log(Warning, "No IPs assigned to interface %s", i.Name)
			continue
		}
		for _, addr := range addresses {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				Log(Error, "Failed to parse CIDR notation: %v", err)
				continue
			}

			if ip.IsGlobalUnicast() && p.IsIPv4(ip.String()) {
				if !FilterInterface(i.Name, ip.String()) {
					ips = append(ips, ip)
				} else {
					reserve = append(reserve, ip)
				}
			}
		}
	}
	if len(ips) == 0 {
		return reserve
	}
	return ips
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

// SrvLookup will search for specified service under provided domain
// and return a map of net.Addr sorted by priority
func SrvLookup(name, proto, domain string) (map[int]string, error) {
	cname, addrs, err := net.LookupSRV(name, proto, domain)
	if err != nil {
		return nil, err
	}
	Log(Debug, "SRV lookup for name cname: %s addrs: %+v", cname, addrs)
	result := make(map[int]string)
	i := 0
	for _, addr := range addrs {
		Log(Trace, "Lookup result: %s:%d", addr.Target, addr.Port)
		result[i] = fmt.Sprintf("%s:%d", addr.Target, addr.Port)
		i++
	}

	return result, nil
}
