// +build linux

package ptp

import (
	"encoding/binary"
	"fmt"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
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

func newEmptyTAP() *TAPLinux {
	return &TAPLinux{}
}

// TAPLinux is an interface for TAP device on Linux platform
type TAPLinux struct {
	IP         net.IP           // IP
	Subnet     net.IP           // Subnet
	Mask       net.IPMask       // Mask
	Mac        net.HardwareAddr // Hardware Address
	Name       string           // Network interface name
	Tool       string           // Path to `ip`
	MTU        int              // MTU value
	fd         int              // File descriptor
	Configured bool             // Whether interface was configured
	PMTU       bool             // Enables/Disbles PMTU
	Auto       bool
	Status     InterfaceStatus
	file       *os.File // Interface descriptor
	//file       unix.FileHandle  // TAP Interface File Handle
}

// GetName returns a name of interface
func (tap *TAPLinux) GetName() string {
	return tap.Name
}

// GetHardwareAddress returns a MAC address of the interface
func (tap *TAPLinux) GetHardwareAddress() net.HardwareAddr {
	return tap.Mac
}

// GetIP returns IP addres of the interface
func (tap *TAPLinux) GetIP() net.IP {
	return tap.IP
}

func (tap *TAPLinux) GetSubnet() net.IP {
	return tap.Subnet
}

// GetMask returns an IP mask of the interface
func (tap *TAPLinux) GetMask() net.IPMask {
	return tap.Mask
}

// GetBasename returns a prefix for automatically generated interface names
func (tap *TAPLinux) GetBasename() string {
	return "vptp"
}

// SetName will set interface name
func (tap *TAPLinux) SetName(name string) {
	tap.Name = name
}

// SetHardwareAddress will set MAC
func (tap *TAPLinux) SetHardwareAddress(mac net.HardwareAddr) {
	tap.Mac = mac
}

// SetIP will set IP
func (tap *TAPLinux) SetIP(ip net.IP) {
	tap.IP = ip
}

func (tap *TAPLinux) SetSubnet(subnet net.IP) {
	tap.Subnet = subnet
}

// SetMask will set mask
func (tap *TAPLinux) SetMask(mask net.IPMask) {
	tap.Mask = mask
}

// Init will initialize TAP interface creation process
func (tap *TAPLinux) Init(name string) error {
	if name == "" {
		return fmt.Errorf("Failed to configure interface: empty name")
	}
	tap.Name = name
	return nil
}

// Open will open a file descriptor for a new interface
func (tap *TAPLinux) Open() error {
	var err error
	if tap.file != nil {
		return fmt.Errorf("TAP device is already acquired")
	}
	tap.fd, err = unix.Open("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return err
	}
	tap.file = os.NewFile(uintptr(tap.fd), "/dev/net/tun")
	err = tap.createInterface()
	if err != nil {
		return err
	}
	return nil
}

// Close will close TAP interface by closing it's file descriptor
func (tap *TAPLinux) Close() error {
	if tap.file == nil {
		return fmt.Errorf("nil interface file descriptor")
	}
	Log(Info, "Closing network interface %s", tap.GetName())
	err := tap.file.Close()
	if err != nil {
		return fmt.Errorf("Failed to close network interface: %s", err)
	}
	Log(Info, "Interface closed")
	return nil
}

// Configure will configure interface using system calls to commands
func (tap *TAPLinux) Configure(lazy bool) error {
	tap.Status = InterfaceConfiguring
	if lazy {
		return nil
	}
	Log(Info, "Configuring %s. IP: %s, Mac: %s", tap.Name, tap.IP.String(), tap.Mac.String())
	err := tap.linkUp()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}

	err = tap.setMTU()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}

	// Configure new device
	err = tap.setIP()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}
	err = tap.linkDown()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}
	err = tap.setMac()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}
	err = tap.linkUp()
	if err != nil {
		tap.Status = InterfaceBroken
		return err
	}
	tap.Status = InterfaceConfigured
	return nil
}

func (tap *TAPLinux) Deconfigure() error {
	tap.Status = InterfaceDeconfigured
	return nil
}

// ReadPacket will read single packet from network interface
func (tap *TAPLinux) ReadPacket() (*Packet, error) {
	buf := make([]byte, 4096)

	n, err := tap.file.Read(buf)
	if err != nil {
		Log(Error, "Failed to read packet: %+v", err)
		return nil, err
	}

	return tap.handlePacket(buf[0:n])
}

func (tap *TAPLinux) handlePacket(data []byte) (*Packet, error) {
	length := len(data)
	if length < 14 {
		return nil, errPacketTooSmall
	}
	pkt := &Packet{Packet: data[0:length]}
	pkt.Protocol = int(binary.BigEndian.Uint16(data[12:14]))

	if !tap.IsPMTUEnabled() {
		return pkt, nil
	}

	if pkt.Protocol == int(PacketIPv4) {
		// Return packet
		skip, err := pmtu(data, tap)
		if skip {
			return nil, err
		}
	}
	return pkt, nil
}

// WritePacket will write a single packet to interface
func (tap *TAPLinux) WritePacket(packet *Packet) error {
	n, err := tap.file.Write(packet.Packet)
	if err != nil {
		return err
	}
	if n != len(packet.Packet) {
		return io.ErrShortWrite
	}
	return nil
}

// Run will start TAP processes
func (tap *TAPLinux) Run() {

}

func (tap *TAPLinux) createInterface() error {
	var req ifReq
	req.Flags = 0
	copy(req.Name[:15], tap.Name)
	req.Flags |= iffTap
	req.Flags |= iffnopi
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, uintptr(tap.fd), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if err != 0 {
		return err
	}
	return nil
}

func (tap *TAPLinux) setMTU() error {
	mtu := fmt.Sprintf("%d", tap.MTU)
	setmtu := exec.Command(tap.Tool, "link", "set", "dev", tap.Name, "mtu", mtu)
	err := setmtu.Run()
	if err != nil {
		Log(Error, "Failed to set MTU on device %s: %v", tap.Name, err)
		return err
	}
	return nil
}

func (tap *TAPLinux) linkUp() error {
	linkup := exec.Command(tap.Tool, "link", "set", "dev", tap.Name, "up")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func (tap *TAPLinux) linkDown() error {
	linkup := exec.Command(tap.Tool, "link", "set", "dev", tap.Name, "down")
	err := linkup.Run()
	if err != nil {
		Log(Error, "Failed to up link: %v", err)
		return err
	}
	return nil
}

func (tap *TAPLinux) setIP() error {
	Log(Info, "Setting %s IP on device %s", tap.IP.String(), tap.Name)
	setip := exec.Command(tap.Tool, "addr", "add", tap.IP.String()+"/24", "dev", tap.Name)
	err := setip.Run()
	if err != nil {
		Log(Error, "Failed to set IP: %v", err)
		return err
	}
	return err
}

func (tap *TAPLinux) setMac() error {
	Log(Info, "Setting %s MAC on device %s", tap.Mac.String(), tap.Name)
	setmac := exec.Command(tap.Tool, "link", "set", "dev", tap.Name, "address", tap.Mac.String())
	err := setmac.Run()
	if err != nil {
		Log(Error, "Failed to set MAC: %v", err)
		return err
	}
	return err
}

func (tap *TAPLinux) IsConfigured() bool {
	return tap.Configured
}

func (tap *TAPLinux) MarkConfigured() {
	tap.Configured = true
}

func (tap *TAPLinux) EnablePMTU() {
	tap.PMTU = true
}

func (tap *TAPLinux) DisablePMTU() {
	tap.PMTU = false
}

func (tap *TAPLinux) IsPMTUEnabled() bool {
	return tap.PMTU
}

func (tap *TAPLinux) IsBroken() bool {
	return false
}

func (tap *TAPLinux) SetAuto(auto bool) {
	tap.Auto = auto
}

func (tap *TAPLinux) IsAuto() bool {
	return tap.Auto
}

func (tap *TAPLinux) GetStatus() InterfaceStatus {
	return tap.Status
}

// FilterInterface will return true if this interface needs to be filtered out
func FilterInterface(infName, infIP string) bool {
	if len(infIP) > 4 && infIP[0:3] == "172" {
		return true
	}
	for _, ip := range ActiveInterfaces {
		if ip.String() == infIP {
			return true
		}
	}
	Log(Trace, "ping -4 -w 1 -c 1 -I %s ptest.subutai.io", infName)
	ping := exec.Command("ping", "-4", "-w", "1", "-c", "1", "-I", infName, "ptest.subutai.io")
	if ping.Run() != nil {
		Log(Debug, "Filtered %s %s", infName, infIP)
		return true
	}
	return false
}
