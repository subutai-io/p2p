package ptp

// DevKind Type of the device
type DevKind int

const (
	// DevTun Receive/send layer routable 3 packets (IP, IPv6...). Notably,
	// you don't receive link-local multicast with this interface
	// type.
	DevTun DevKind = iota
	// DevTap Receive/send Ethernet II frames. You receive all packets that
	// would be visible on an Ethernet link, including broadcast and
	// multicast traffic.
	DevTap
)

// Packet represents a packet received on TUN/TAP interface
type Packet struct {
	// The Ethernet type of the packet. Commonly seen values are
	// 0x8000 for IPv4 and 0x86dd for IPv6.
	Protocol int
	// True if the packet was too large to be read completely.
	Truncated bool
	// The raw bytes of the Ethernet payload (for DevTun) or the full
	// Ethernet frame (for DevTap).
	Packet []byte
}

// InterfaceName - The name of the interface. May be different from the name given to
// Open(), if the latter was a pattern.
func (t *Interface) InterfaceName() string {
	return t.Name
}
