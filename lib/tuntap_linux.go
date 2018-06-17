// +build linux

package ptp

import (
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/mdlayher/ethernet"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

// Constants
const (
	ConfigDir string = "/usr/local/etc"
	//DefaultMTU int    = 1376
	DefaultMTU int = 1500
)

// GetDeviceBase returns a default interface name
func GetDeviceBase() string {
	return "vptp"
}

// GetConfigurationTool function will return path to configuration tool on specific platform
func GetConfigurationTool() string {
	path, err := exec.LookPath("ip")
	if err != nil {
		Log(Error, "Failed to find `ip` in path. Returning default /bin/ip")
		return "/bin/ip"
	}
	Log(Info, "Network configuration tool found: %s", path)
	return path
}

func newTAP(tool, ip, mac, mask string, mtu int, pmtu bool) (*TAPLinux, error) {
	Log(Debug, "Acquiring TAP interface [Linux]")
	nip := net.ParseIP(ip)
	if nip == nil {
		return nil, fmt.Errorf("Failed to parse IP during TAP creation")
	}
	nmac, err := net.ParseMAC(mac)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse MAC during TAP creation: %s", err)
	}
	return &TAPLinux{
		Tool: tool,
		IP:   nip,
		Mac:  nmac,
		Mask: net.IPv4Mask(255, 255, 255, 0), // Unused yet
		MTU:  GlobalMTU,
		PMTU: pmtu,
	}, nil
}

// TAPLinux is an interface for TAP device on Linux platform
type TAPLinux struct {
	IP         net.IP           // IP
	Mask       net.IPMask       // Mask
	Mac        net.HardwareAddr // Hardware Address
	Name       string           // Network interface name
	Tool       string           // Path to `ip`
	MTU        int              // MTU value
	file       *os.File         // Interface descriptor
	Configured bool
	PMTU       bool // Enables/Disbles PMTU
}

// GetName returns a name of interface
func (t *TAPLinux) GetName() string {
	return t.Name
}

// GetHardwareAddress returns a MAC address of the interface
func (t *TAPLinux) GetHardwareAddress() net.HardwareAddr {
	return t.Mac
}

// GetIP returns IP addres of the interface
func (t *TAPLinux) GetIP() net.IP {
	return t.IP
}

// GetMask returns an IP mask of the interface
func (t *TAPLinux) GetMask() net.IPMask {
	return t.Mask
}

// GetBasename returns a prefix for automatically generated interface names
func (t *TAPLinux) GetBasename() string {
	return "vptp"
}

// SetName will set interface name
func (t *TAPLinux) SetName(name string) {
	t.Name = name
}

// SetHardwareAddress will set MAC
func (t *TAPLinux) SetHardwareAddress(mac net.HardwareAddr) {
	t.Mac = mac
}

// SetIP will set IP
func (t *TAPLinux) SetIP(ip net.IP) {
	t.IP = ip
}

// SetMask will set mask
func (t *TAPLinux) SetMask(mask net.IPMask) {
	t.Mask = mask
}

// Init will initialize TAP interface creation process
func (t *TAPLinux) Init(name string) error {
	t.Name = name
	return nil
}

// Open will open a file descriptor for a new interface
func (t *TAPLinux) Open() error {
	var err error
	if t.file != nil {
		return fmt.Errorf("TAP device is already acquired")
	}
	t.file, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	err = t.createInterface()
	if err != nil {
		return err
	}
	return nil
}

// Close will close TAP interface by closing it's file descriptor
func (t *TAPLinux) Close() error {
	if t.file != nil {
		Log(Info, "Closing network interface %s", t.GetName())
		err := t.file.Close()
		if err != nil {
			return fmt.Errorf("Failed to close network interface: %s", err)
		}
		Log(Info, "Interface closed")
		return nil
	}
	Log(Warning, "Skipping previously closed interface")
	return nil
}

// Configure will configure interface using system calls to commands
func (t *TAPLinux) Configure() error {
	Log(Info, "Configuring %s. IP: %s, Mac: %s", t.Name, t.IP.String(), t.Mac.String())
	err := t.linkUp()
	if err != nil {
		return err
	}

	err = t.setMTU()
	if err != nil {
		return err
	}

	// Configure new device
	err = t.setIP()
	if err != nil {
		return err
	}
	err = t.linkDown()
	if err != nil {
		return err
	}
	err = t.setMac()
	if err != nil {
		return err
	}
	return t.linkUp()
}

// ReadPacket will read single packet from network interface
func (t *TAPLinux) ReadPacket() (*Packet, error) {
	buf := make([]byte, 4096)

	n, err := t.file.Read(buf)
	if err != nil {
		Log(Error, "Failed to read packet: %+v", err)
		return nil, err
	}

	return t.handlePacket(buf[0:n])
}

func checksum(bytes []byte) uint16 {
	// Clear checksum bytes
	bytes[10] = 0
	bytes[11] = 0

	// Compute checksum
	var csum uint32
	for i := 0; i < len(bytes); i += 2 {
		csum += uint32(bytes[i]) << 8
		csum += uint32(bytes[i+1])
	}
	for {
		// Break when sum is less or equals to 0xFFFF
		if csum <= 65535 {
			break
		}
		// Add carry to the sum
		csum = (csum >> 16) + uint32(uint16(csum))
	}
	// Flip all the bits
	return ^uint16(csum)
}

func (t *TAPLinux) handlePacket(data []byte) (*Packet, error) {
	length := len(data)
	if length < 14 {
		return nil, errPacketTooSmall
	}
	pkt := &Packet{Packet: data[0:length]}
	pkt.Protocol = int(binary.BigEndian.Uint16(data[12:14]))

	if !t.IsPMTUEnabled() {
		return pkt, nil
	}

	if pkt.Protocol == int(PacketIPv4) && length > GlobalMTU-150 {
		header, err := ipv4.ParseHeader(data[14:])
		if err != nil {
			Log(Error, "Failed to parse IPv4 packet: %s", err.Error())
			return nil, nil
		}

		// Don't fragment flag is set. We need to respond with ICMP Destination Unreachable
		if header.Flags == ipv4.DontFragment {
			// Extract packet contents as an ethernet frame for later re-use
			f := new(ethernet.Frame)
			if err := f.UnmarshalBinary(data); err != nil {
				Log(Error, "Failed to Unmarshal IPv4")
				return nil, nil
			}

			// Build "Fragmentation needed" ICMP message
			packetICMP := &icmp.Message{
				Type: ipv4.ICMPTypeDestinationUnreachable,
				Code: 4,
				Body: &icmp.PacketTooBig{
					MTU:  GlobalMTU - 200,    // Next-hop MTU
					Data: data[14 : 14+20+8], // Original header and 64-bits of datagram
				},
			}
			payloadICMP, err := packetICMP.Marshal(nil)
			if err != nil {
				Log(Error, "Failed to marshal ICMP: %s", err.Error())
				return nil, errICMPMarshalFailed
			}

			// Build IPv4 Header
			iph := &ipv4.Header{
				Version:  4,
				Len:      20, // Precalculated header length
				TOS:      0,
				TotalLen: len(payloadICMP) + 20,
				ID:       25,
				TTL:      64,
				Protocol: 1,
				Dst:      header.Src,
				Src:      header.Dst,
				Checksum: 0,
			}
			ipHeader, err := iph.Marshal()
			if err != nil {
				Log(Error, "Failed to marshal header: %s", err.Error())
				return nil, nil
			}

			// Calculate IPv4 header checksum
			hcsum := checksum(ipHeader)
			binary.BigEndian.PutUint16(ipHeader[10:], hcsum)

			// Build new ethernet frame. Swap dst/src
			pl := append(ipHeader, payloadICMP...)
			nf := new(ethernet.Frame)
			nf.Destination = f.Source
			nf.Source = f.Destination
			nf.EtherType = ethernet.EtherTypeIPv4
			nf.Payload = pl
			rpacket, err := nf.MarshalBinary()
			if err != nil {
				Log(Error, "Failed to marshal ethernet")
				return nil, nil
			}

			// Calculate CRC32 checksum for ethernet frame
			crc := make([]byte, 4)
			binary.LittleEndian.PutUint32(crc, crc32.ChecksumIEEE(rpacket))
			rpacket = append(rpacket, crc...)

			// Send frame to the interface
			// P2P will drop packet afterwards
			t.WritePacket(&Packet{int(PacketIPv4), rpacket})
			return nil, nil
		}
	}

	// Return packet
	return pkt, nil
}

// WritePacket will write a single packet to interface
func (t *TAPLinux) WritePacket(packet *Packet) error {
	n, err := t.file.Write(packet.Packet)
	if err != nil {
		return err
	}
	if n != len(packet.Packet) {
		return io.ErrShortWrite
	}
	return nil
}

// Run will start TAP processes
func (t *TAPLinux) Run() {

}

func (t *TAPLinux) createInterface() error {
	var req ifReq
	req.Flags = 0
	copy(req.Name[:15], t.Name)
	req.Flags |= iffTap
	req.Flags |= iffnopi
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, t.file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != 0 {
		return err
	}
	return nil
}

func (t *TAPLinux) setMTU() error {
	mtu := fmt.Sprintf("%d", t.MTU)
	setmtu := exec.Command(t.Tool, "link", "set", "dev", t.Name, "mtu", mtu)
	err := setmtu.Run()
	if err != nil {
		Log(Error, "Failed to set MTU on device %s: %v", t.Name, err)
		return err
	}
	return nil
}

func (t *TAPLinux) linkUp() error {
	linkup := exec.Command(t.Tool, "link", "set", "dev", t.Name, "up")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func (t *TAPLinux) linkDown() error {
	linkup := exec.Command(t.Tool, "link", "set", "dev", t.Name, "down")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func (t *TAPLinux) setIP() error {
	Log(Info, "Setting %s IP on device %s", t.IP.String(), t.Name)
	setip := exec.Command(t.Tool, "addr", "add", t.IP.String()+"/24", "dev", t.Name)
	err := setip.Run()
	if err != nil {
		Log(Error, "Failed to set IP: %v", err)
		return err
	}
	return err
}

func (t *TAPLinux) setMac() error {
	Log(Info, "Setting %s MAC on device %s", t.Mac.String(), t.Name)
	setmac := exec.Command(t.Tool, "link", "set", "dev", t.Name, "address", t.Mac.String())
	err := setmac.Run()
	if err != nil {
		Log(Error, "Failed to set MAC: %v", err)
		return err
	}
	return err
}

func (t *TAPLinux) IsConfigured() bool {
	return t.Configured
}

func (t *TAPLinux) MarkConfigured() {
	t.Configured = true
}

func (t *TAPLinux) EnablePMTU() {
	t.PMTU = true
}

func (t *TAPLinux) DisablePMTU() {
	t.PMTU = false
}

func (t *TAPLinux) IsPMTUEnabled() bool {
	return t.PMTU
}

func (t *TAPLinux) IsBroken() bool {
	return false
}

// FilterInterface will return true if this interface needs to be filtered out
func FilterInterface(infName, infIP string) bool {
	if len(infIP) > 4 && infIP[0:3] == "172" {
		return true
	}
	Log(Trace, "ping -4 -w 1 -c 1 -I %s ptest.subutai.io", infName)
	ping := exec.Command("ping", "-4", "-w", "1", "-c", "1", "-I", infName, "ptest.subutai.io")
	if ping.Run() != nil {
		Log(Debug, "Filtered %s %s", infName, infIP)
		return true
	}
	return false
}
