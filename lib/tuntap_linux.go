// +build linux

package ptp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"syscall"
	"unsafe"
)

// Constants
const (
	ConfigDir  string = "/usr/local/etc"
	DefaultMTU int    = 1376
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

func newTAP(tool, ip, mac, mask string, mtu int) (*TAPLinux, error) {
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
		MTU:  DefaultMTU,
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

	var pkt *Packet
	pkt = &Packet{Packet: buf[0:n]}
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[12:14]))
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
