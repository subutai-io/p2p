// +build windows

package ptp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os/exec"
	"syscall"
	"unicode/utf16"

	"golang.org/x/sys/windows"
)

// Windows platform specific constants
const (
	ConfigDir        string         = "C:\\ProgramData\\Subutai\\etc"
	DefaultMTU       int            = 1376
	NetworkKey       string         = "SYSTEM\\CurrentControlSet\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	AdapterKey       string         = "SYSTEM\\CurrentControlSet\\Control\\Class\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	NoMoreItems      syscall.Errno  = 259
	UsermodDeviceDir string         = "\\\\.\\Global\\"
	SysDeviceDir     string         = "\\Device\\"
	UserDeviceDir    string         = "\\DosDevices\\Global\\"
	InvalidHandle    syscall.Handle = 0
)

var UsedInterfaces []string // List of interfaces currently in use by p2p daemon

var (
	getMacIOCTL             = tapControlCode(1, 0)
	getVersionIOCTL         = tapControlCode(2, 0)
	getMTUValueIOCTL        = tapControlCode(3, 0)
	getInfoIOCTL            = tapControlCode(4, 0)
	configPointToPointIOCTL = tapControlCode(5, 0)
	setMediaStatusIOCTL     = tapControlCode(6, 0)
	configDHCPMasqIOCTL     = tapControlCode(7, 0)
	configGetLogLineIOCTL   = tapControlCode(8, 0)
	configDHCPSetOptIOCTL   = tapControlCode(9, 0)
	configTUNIOCTL          = tapControlCode(10, 0)
)

// GetDeviceBase returns a default interface name
func GetDeviceBase() string {
	return "Local Area Network "
}

// GetConfigurationTool function will return path to configuration tool on specific platform
func GetConfigurationTool() string {
	path, err := exec.LookPath("netsh")
	if err != nil {
		Log(Error, "Failed to find `netsh` in path. Returning default netsh")
		return "netsh"
	}
	Log(Info, "Network configuration tool found: %s", path)
	return path
}

func newTAP(tool, ip, mac, mask string, mtu int) (*TAPWindows, error) {
	Log(Debug, "Acquiring TAP interface [Windows]")
	nip := net.ParseIP(ip)
	if nip == nil {
		return nil, fmt.Errorf("Failed to parse IP during TAP creation")
	}
	nmac, err := net.ParseMAC(mac)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse MAC during TAP creation: %s", err)
	}
	return &TAPWindows{
		Tool:      tool,
		IP:        nip,
		Mac:       nmac,
		Mask:      net.IPv4Mask(255, 255, 255, 0), // Unused yet
		MTU:       DefaultMTU,
		MacNotSet: true,
	}, nil
}

// TAPLinux is an interface for TAP device on Linux platform
type TAPWindows struct {
	IP        net.IP           // IP
	Mask      net.IPMask       // Mask
	Mac       net.HardwareAddr // Hardware Address
	MacNotSet bool
	Name      string // Network interface name
	Interface string // ?????????????????
	Tool      string // Path to `ip`
	MTU       int    // MTU value
	file      syscall.Handle
	Handle    syscall.Handle
	Rx        chan []byte
	Tx        chan []byte
}

// GetName returns a name of interface
func (t *TAPWindows) GetName() string {
	return t.Name
}

// GetHardwareAddress returns a MAC address of the interface
func (t *TAPWindows) GetHardwareAddress() net.HardwareAddr {
	if t.MacNotSet {
		mac := make([]byte, 6)
		var length uint32
		err := syscall.DeviceIoControl(t.file, getMacIOCTL, &mac[0], uint32(len(mac)), &mac[0], uint32(len(mac)), &length, nil)
		if err != nil {
			Log(Error, "Failed to retrieve Mac")
			return t.Mac
		}
		var macAddr bytes.Buffer

		i := 0
		for _, a := range mac {
			if a == 0 {
				macAddr.WriteString("00")
			} else if a < 16 {
				macAddr.WriteString(fmt.Sprintf("0%x", a))
			} else {
				macAddr.WriteString(fmt.Sprintf("%x", a))
			}
			if i < 5 {
				macAddr.WriteString(":")
			}
			i++
		}
		Log(Debug, "MAC: %s", macAddr.String())
		deviceMac, err := net.ParseMAC(macAddr.String())
		if err != nil {
			Log(Error, "Failed to extract mac: %s", err)
		}
		t.Mac = deviceMac
		t.MacNotSet = false
	}
	return t.Mac
}

// GetIP returns IP addres of the interface
func (t *TAPWindows) GetIP() net.IP {
	return t.IP
}

// GetMask returns an IP mask of the interface
func (t *TAPWindows) GetMask() net.IPMask {
	return t.Mask
}

// GetBasename returns a prefix for automatically generated interface names
func (t *TAPWindows) GetBasename() string {
	return "vptp"
}

// SetName will set interface name
func (t *TAPWindows) SetName(name string) {
	t.Name = name
}

// SetHardwareAddress will set MAC
func (t *TAPWindows) SetHardwareAddress(mac net.HardwareAddr) {
	t.Mac = mac
}

// SetIP will set IP
func (t *TAPWindows) SetIP(ip net.IP) {
	t.IP = ip
}

// SetMask will set mask
func (t *TAPWindows) SetMask(mask net.IPMask) {
	t.Mask = mask
}

// Init will initialize TAP interface creation process
func (t *TAPWindows) Init(name string) error {
	t.Name = name
	return nil
}

func (t *TAPWindows) Open() error {
	handle, err := t.queryNetworkKey()
	if err != nil {
		Log(Error, "Failed to query Windows registry: %v", err)
		return err
	}
	err = t.queryAdapters(handle)
	if err != nil {
		Log(Error, "Failed to query network adapters: %v", err)
		return err
	}
	if t.Name == "" {
		Log(Error, "Failed to query network adapters: %v", err)
		return errors.New("Empty network adapter")
	}
	err = syscall.CloseHandle(handle)
	if err != nil {
		Log(Error, "Failed to close retrieved handle: %v", err)
	}
	return nil
}

// Close will close handle for TAP interface
func (t *TAPWindows) Close() error {
	for i, iface := range UsedInterfaces {
		if iface == t.Interface {
			UsedInterfaces = append(UsedInterfaces[:i], UsedInterfaces[i+1:]...)
			break
		}
	}
	return syscall.CloseHandle(t.file)
}

// Configure will configure TAP interface and set it's IP, Mask and other
// parameters
func (t *TAPWindows) Configure() error {
	Log(Debug, "Configuring %s. IP: %s Mask: %s", t.Interface, t.IP.String(), t.Mask.String())
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}
	// TODO: Unhardcode mask
	cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, t.Interface, t.IP.String(), "255.255.255.0")
	Log(Debug, "Executing: %s", cmd)
	setip.SysProcAttr.CmdLine = cmd
	err := setip.Run()
	if err != nil {
		return fmt.Errorf("Failed to properly configure TAP device with netsh: %v", err)
	}

	in := []byte("\x01\x00\x00\x00")
	var length uint32
	err = syscall.DeviceIoControl(t.file, setMediaStatusIOCTL,
		&in[0],
		uint32(len(in)),
		&in[0],
		uint32(len(in)),
		&length,
		nil)
	if err != nil {
		return fmt.Errorf("Failed to change device status to 'connected': %v", err)
	}
	return nil
}

// Run will start read/write goroutines
func (t *TAPWindows) Run() {
	Log(Info, "Listening for TAP interface")
	t.Rx = make(chan []byte, 1500)
	t.Tx = make(chan []byte, 1500)
	go func() {
		if err := t.read(t.Rx); err != nil {
			Log(Error, "Failed to read packet: %v", err)
		}
	}()
	go func() {
		if err := t.write(t.Tx); err != nil {
			Log(Error, "Failed ro write packet: %v", err)
		}
	}()
}

func (t *TAPWindows) ReadPacket() (*Packet, error) {
	buf := <-t.Rx
	n := len(buf)
	if n <= 4 {
		return nil, nil
	}
	p := 12
	var pkt *Packet
	pkt = &Packet{Packet: buf[0:n]}
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[p : p+2]))
	// flags := int(*(*uint16)(unsafe.Pointer(&buf[0])))
	// if flags&flagTruncated != 0 {
	// 	pkt.Truncated = true
	// }
	// pkt.Truncated = false
	return pkt, nil
}

func (t *TAPWindows) WritePacket(pkt *Packet) error {
	n := len(pkt.Packet)
	buf := make([]byte, n)
	copy(buf, pkt.Packet)
	t.Tx <- buf[:n]
	return nil
}

func (t *TAPWindows) read(ch chan []byte) (err error) {
	rx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	rx.HEvent = syscall.Handle(hevent)
	buf := make([]byte, 1500)
	var l uint32
	for {
		if err := syscall.ReadFile(t.file, buf, &l, &rx); err != nil {
		}
		if _, err := syscall.WaitForSingleObject(rx.HEvent, syscall.INFINITE); err != nil {
			Log(Error, "Failed to read from TUN/TAP: %v", err)
		}
		rx.Offset += l
		ch <- buf
	}
}

func (t *TAPWindows) write(ch chan []byte) (err error) {
	tx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	tx.HEvent = syscall.Handle(hevent)
	for {
		select {
		case data := <-ch:
			var l uint32
			syscall.WriteFile(t.file, data, &l, &tx)
			syscall.WaitForSingleObject(tx.HEvent, syscall.INFINITE)
			tx.Offset += uint32(len(data))
		}
	}
}

func (t *TAPWindows) queryNetworkKey() (syscall.Handle, error) {
	var handle syscall.Handle
	err := syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(NetworkKey), 0, syscall.KEY_READ, &handle)
	if err != nil {
		return 0, err
	}
	return handle, nil
}

func (t *TAPWindows) queryAdapters(handle syscall.Handle) error {
	var index uint32
	for {
		var length uint32 = 72
		adapter := make([]uint16, length)
		err := syscall.RegEnumKeyEx(handle, index, &adapter[0], &length, nil, nil, nil, nil)
		if err == NoMoreItems {
			Log(Debug, "No more items in Windows Registry")
			return nil
		}
		index++
		adapterID := string(utf16.Decode(adapter[0:length]))
		adapterID = t.removeZeroes(adapterID)
		path := fmt.Sprintf("%s\\%s\\Connection", NetworkKey, adapterID)
		var iHandle syscall.Handle
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &iHandle)
		if err != nil {
			syscall.RegCloseKey(iHandle)
			continue
		}
		length = 1024
		aName := make([]byte, length)
		err = syscall.RegQueryValueEx(iHandle, syscall.StringToUTF16Ptr("Name"), nil, nil, &aName[0], &length)

		if err != nil {
			syscall.RegCloseKey(iHandle)
			continue
		}

		aNameUtf16 := make([]uint16, length/2)
		for i := 0; i < int(length)-2; i += 2 {
			aNameUtf16[i/2] = binary.LittleEndian.Uint16(aName[i:])
		}
		aNameUtf16[length/2-1] = 0

		adapterName := string(utf16.Decode(aNameUtf16))
		adapterName = t.removeZeroes(adapterName)
		Log(Debug, "AdapterName : %s, len : %d", adapterName, len(adapterName))

		var isInUse = false
		for _, i := range UsedInterfaces {
			if i == adapterName {
				isInUse = true
			}
		}
		if isInUse {
			Log(Debug, "Adapter already in use. Skipping.")
			continue
		}
		UsedInterfaces = append(UsedInterfaces, adapterName)

		tapname := fmt.Sprintf("%s%s%s", UsermodDeviceDir, adapterID, TapSuffix)

		t.file, err = syscall.CreateFile(syscall.StringToUTF16Ptr(tapname),
			syscall.GENERIC_WRITE|syscall.GENERIC_READ,
			0,
			nil,
			syscall.OPEN_EXISTING,
			syscall.FILE_ATTRIBUTE_SYSTEM|syscall.FILE_FLAG_OVERLAPPED,
			0)
		if err != nil {
			syscall.CloseHandle(t.Handle)
			continue
		}
		Log(Info, "Acquired control over TAP interface: %s", adapterName)
		t.Name = adapterID
		t.Interface = adapterName
		return nil
	}
	return nil
}

func (t *TAPWindows) removeZeroes(s string) string {
	bytes := []byte(s)
	var (
		res  []byte
		prev bool
		size int
	)
	for _, b := range bytes {
		if b == 0 && prev {
			break
		} else if b == 0 && !prev {
			prev = true
		} else {
			prev = false
			res = append(res, b)
			size++
		}
	}
	return string(res[:size])
}

func tapControlCode(request, method uint32) uint32 {
	return controlCode(34, request, method, 0)
}

func controlCode(device_type, function, method, access uint32) uint32 {
	return (device_type << 16) | (access << 14) | (function << 2) | method
}
