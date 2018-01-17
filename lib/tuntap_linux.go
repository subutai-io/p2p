package ptp

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"syscall"
	"unsafe"
)

// Constants
const (
	ConfigDir  string = "/usr/local/etc"
	DefaultMTU string = "1376"
)

// func openDevice(ifPattern string) (*os.File, error) {
// 	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
// 	return file, err
// }

// func createInterface(file *os.File, ifPattern string, kind DevKind) (string, error) {
// 	var req ifReq
// 	req.Flags = 0
// 	copy(req.Name[:15], ifPattern)
// 	switch kind {
// 	case DevTun:
// 		req.Flags |= iffTun
// 	case DevTap:
// 		req.Flags |= iffTap
// 	default:
// 		panic("Unknown interface type")
// 	}
// 	req.Flags |= iffnopi

// 	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
// 	if err != 0 {
// 		return "", err
// 	}
// 	return string(req.Name[:]), nil
// }

// func closeInterface(file *os.File) {
// 	Log(Info, "Closing network interface")
// 	if file != nil {
// 		err := file.Close()
// 		if err != nil {
// 			Log(Error, "Failed to close network interface: %s", err)
// 			return
// 		}
// 		Log(Info, "Interface closed")
// 		return
// 	}
// 	Log(Warning, "Skipping previously closed interface")
// }

// // ConfigureInterface performs a configuration of an existing interface
// func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
// 	Log(Info, "Configuring %s. IP: %s, Mac: %s", device, ip, mac)
// 	err := LinkUp(device, tool)
// 	if err != nil {
// 		return err
// 	}

// 	err = SetMTU(dev, device, tool, DefaultMTU)
// 	if err != nil {
// 		return err
// 	}

// 	// Configure new device
// 	err = SetIP(ip, device, tool)
// 	if err != nil {
// 		return err
// 	}

// 	err = SetMac(mac, device, tool)
// 	return err
// }

// // SetMTU sets an MTU value
// func SetMTU(dev *Interface, device, tool, mtu string) error {
// 	setmtu := exec.Command(tool, "link", "set", "dev", device, "mtu", mtu)
// 	err := setmtu.Run()
// 	if err != nil {
// 		Log(Error, "Failed to set MTU on device %s: %v", device, err)
// 		return err
// 	}
// 	return nil
// }

// // LinkUp brings interface up
// func LinkUp(device, tool string) error {
// 	linkup := exec.Command(tool, "link", "set", "dev", device, "up")
// 	err := linkup.Run()
// 	if err != nil {
// 		Log(Error, "Failed to up link: %v", err)
// 		return err
// 	}
// 	return nil
// }

// // SetIP sets an IP address to an interface
// func SetIP(ip, device, tool string) error {
// 	Log(Info, "Setting %s IP on device %s", ip, device)
// 	setip := exec.Command(tool, "addr", "add", ip+"/24", "dev", device)
// 	err := setip.Run()
// 	if err != nil {
// 		Log(Error, "Failed to set IP: %v", err)
// 		return err
// 	}
// 	return err
// }

// // SetMac sets a MAC address to a device
// func SetMac(mac, device, tool string) error {
// 	// Set MAC to device
// 	Log(Info, "Setting %s MAC on device %s", mac, device)
// 	setmac := exec.Command(tool, "link", "set", "dev", device, "address", mac)
// 	err := setmac.Run()
// 	if err != nil {
// 		Log(Error, "Failed to set MAC: %v", err)
// 		return err
// 	}
// 	return err
// }

// GetDeviceBase returns a default interface name
func GetDeviceBase() string {
	return "vptp"
}

type TAPLinux struct {
	IP   net.IP           // IP
	Mask net.IPMask       // Mask
	Mac  net.HardwareAddr // Hardware Address
	Name string           // Network interface name
	Tool string           // Path to `ip`
	MTU  int              // MTU value
	file *os.File         // Interface descriptor
}

func (t *TAPLinux) GetName() string {
	return t.Name
}

func (t *TAPLinux) GetHardwareAddress() net.HardwareAddr {
	return t.Mac
}

func (t *TAPLinux) GetIP() net.IP {
	return t.IP
}

func (t *TAPLinux) GetMask() net.IPMask {
	return t.Mask
}

func (t *TAPLinux) GetBasename() string {
	return "vptp"
}

func (t *TAPLinux) Init(name string) error {
	t.Name = name
	return nil
}

func (t *TAPLinux) Open() error {
	var err error
	t.file, err = os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		return err
	}

	if t.createInterface() != nil {
		return err
	}
	return nil
}

func (t *TAPLinux) Close() error {
	if file != nil {
		Log(Info, "Closing network interface %s", t.GetName())
		err := file.Close()
		if err != nil {
			return fmt.Errorf("Failed to close network interface: %s", err)
		}
		Log(Info, "Interface closed")
		return nil
	}
	Log(Warning, "Skipping previously closed interface")
	return nil
}

func (t *TAPLinux) Configure() error {
	Log(Info, "Configuring %s. IP: %s, Mac: %s", t.Name, t.IP.String(), t.Mac())
	err := t.linkUp()
	if err != nil {
		return err
	}

	err = t.setMTU(DefaultMTU)
	if err != nil {
		return err
	}

	// Configure new device
	err = t.setIP(ip)
	if err != nil {
		return err
	}

	return t.setMac(mac)
}

func (t *TAPLinux) ReadPacket() (*Packet, error) {
	buf := make([]byte, 4096)

	n, err := t.file.Read(buf)
	if err != nil {
		return nil, err
	}

	var pkt *Packet
	pkt = &Packet{Packet: buf[0:n]}
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[12:14]))
	pkt.Truncated = false
	return pkt, nil
}

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

func newTAP(tool, ip, mac, mask string, mtu int) (*TAPLinux, error) {
	nip, err := net.ParseIP(ip)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse IP during TAP creation: %s", err)
	}
	nmac, err := net.ParseMAC(mac)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse MAC during TAP creation: %s", err)
	}
	return &TAPLinux{
		Tool: tool,
		IP:   nip,
		Mac:  nmac,
		Mask: "255.255.255.0", // Unused yet
		MTU:  DefaultMTU,
	}, nil
}
