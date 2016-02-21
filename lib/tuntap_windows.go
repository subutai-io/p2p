// +build windows

package ptp

import (
	"fmt"
	//"os"
	"encoding/binary"
	"os/exec"
	"syscall"
	"unicode/utf16"
	"unsafe"
	"golang.org/x/sys/windows"
)

type Interface struct {
	Name string
	//file *os.File
	file      syscall.Handle
	meta      bool
	Handle    syscall.Handle
	Interface string
	IP        string
	Mask      string
	Rx        syscall.Overlapped
	RxE       windows.Handle
	Tx        syscall.Overlapped
	TxE       windows.Handle
	Rl        uint32
	Wl        uint32
}

const (
	NETWORK_KEY         string         = "SYSTEM\\CurrentControlSet\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	ADAPTER_KEY         string         = "SYSTEM\\CurrentControlSet\\Control\\Class\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	NO_MORE_ITEMS       syscall.Errno  = 259
	USERMODE_DEVICE_DIR string         = "\\\\.\\Global\\"
	SYS_DEVICE_DIR      string         = "\\Device\\"
	USER_DEVICE_DIR     string         = "\\DosDevices\\Global\\"
	TAP_SUFFIX          string         = ".tap"
	INVALID_HANDLE      syscall.Handle = 0
)

var (
	TAP_IOCTL_GET_MAC               = TAP_CONTROL_CODE(1, 0)
	TAP_IOCTL_GET_VERSION           = TAP_CONTROL_CODE(2, 0)
	TAP_IOCTL_GET_MTU               = TAP_CONTROL_CODE(3, 0)
	TAP_IOCTL_GET_INFO              = TAP_CONTROL_CODE(4, 0)
	TAP_IOCTL_CONFIG_POINT_TO_POINT = TAP_CONTROL_CODE(5, 0)
	TAP_IOCTL_SET_MEDIA_STATUS      = TAP_CONTROL_CODE(6, 0)
	TAP_IOCTL_CONFIG_DHCP_MASQ      = TAP_CONTROL_CODE(7, 0)
	TAP_IOCTL_GET_LOG_LINE          = TAP_CONTROL_CODE(8, 0)
	TAP_IOCTL_CONFIG_DHCP_SET_OPT   = TAP_CONTROL_CODE(9, 0)
	TAP_IOCTL_CONFIG_TUN            = TAP_CONTROL_CODE(10, 0)
)

func TAP_CONTROL_CODE(request, method uint32) uint32 {
	return CTL_CODE(34, request, method, 0)
}

func CTL_CODE(device_type, function, method, access uint32) uint32 {
	return (device_type << 16) | (access << 14) | (function << 2) | method
}

func removeZeroes(s string) string {
	bytes := []byte(s)
	var (
		res  []byte
		prev bool
		size int = 0
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

func queryNetworkKey() (syscall.Handle, error) {
	var handle syscall.Handle
	err := syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(NETWORK_KEY), 0, syscall.KEY_READ, &handle)
	if err != nil {
		return 0, err
	}
	return handle, nil
}

func queryAdapters(handle syscall.Handle) (*Interface, error) {
	var dev Interface
	var index uint32 = 0
	for {
		var length uint32 = 72
		adapter := make([]uint16, length)
		err := syscall.RegEnumKeyEx(handle, index, &adapter[0], &length, nil, nil, nil, nil)
		if err == NO_MORE_ITEMS {
			break
		}
		index++
		adapterId := string(utf16.Decode(adapter[0:72]))
		adapterId = removeZeroes(adapterId)
		path := fmt.Sprintf("%s\\%s\\Connection", NETWORK_KEY, adapterId)
		var iHandle syscall.Handle
		err = syscall.RegOpenKeyEx(syscall.HKEY_LOCAL_MACHINE, syscall.StringToUTF16Ptr(path), 0, syscall.KEY_READ, &iHandle)
		if err != nil {
			continue
		}
		length = 1024
		aName := make([]byte, length)
		err = syscall.RegQueryValueEx(iHandle, syscall.StringToUTF16Ptr("Name"), nil, nil, &aName[0], &length)
		if err != nil {
			continue
		}
		syscall.RegCloseKey(iHandle)
		adapterName := removeZeroes(string(aName))
		tapname := fmt.Sprintf("%s%s%s", USERMODE_DEVICE_DIR, adapterId, TAP_SUFFIX)

		dev.file, err = syscall.CreateFile(syscall.StringToUTF16Ptr(tapname),
			syscall.GENERIC_WRITE|syscall.GENERIC_READ,
			0,
			nil,
			syscall.OPEN_EXISTING,
			syscall.FILE_ATTRIBUTE_SYSTEM|syscall.FILE_FLAG_OVERLAPPED,
			0)
		if err != nil {
			syscall.CloseHandle(dev.Handle)
			continue
		}
		Log(INFO, "Acquired control over TAP interface: %s", adapterName)
		dev.Name = adapterId
		dev.Interface = adapterName
		return &dev, nil
	}
	return nil, nil
}

func openDevice(ifPattern string) (*Interface, error) {
	handle, err := queryNetworkKey()
	if err != nil {
		Log(ERROR, "Failed to query Windows registry: %v", err)
		return nil, err
	}
	dev, err := queryAdapters(handle)
	if err != nil {
		Log(ERROR, "Failed to query network adapters: %v", err)
		return nil, err
	}
	if dev.Name == "" {
		Log(ERROR, "Failed to query network adapters: %v", err)
		return nil, nil
	}

	return dev, nil
}

func createInterface(file syscall.Handle, ifPattern string, kind DevKind, meta bool) (string, error) {
	return "1", nil
}

func ConfigureInterface(dev *Interface, ip, mac, device, tool string) error {
	
	dev.IP = ip
	dev.Mask = "255.255.255.0"
	Log(INFO, "Configuring %s. IP: %s Mask: %s", dev.Interface, dev.IP, dev.Mask)
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}
	cmd := fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, dev.Interface, dev.IP, dev.Mask)
	Log(INFO, "Executing: %s", cmd)
	setip.SysProcAttr.CmdLine = cmd
	err := setip.Run()
	if err != nil {
		Log(ERROR, "Failed to properly configure TAP device with netsh: %v", err)
		return err
	}

	in := []byte("\x01\x00\x00\x00")
	var length uint32
	err = syscall.DeviceIoControl(dev.file, TAP_IOCTL_SET_MEDIA_STATUS,
		&in[0],
		uint32(len(in)),
		&in[0],
		uint32(len(in)),
		&length,
		nil)
	if err != nil {
		Log(ERROR, "Failed to change device status to 'connected': %v", err)
		return err
	}

	Log(INFO, "Configuring overlapped Rx & Tx for Windows-TAP i/o operations")
	dev.Rx = syscall.Overlapped{}
	dev.RxE, err = windows.CreateEvent(nil, 0, 0, nil)
	dev.Rx.HEvent = syscall.Handle(dev.RxE)
	dev.Tx = syscall.Overlapped{}
	dev.TxE, err = windows.CreateEvent(nil, 0, 0, nil)
	dev.Tx.HEvent = syscall.Handle(dev.TxE)

	return nil
}

func LinkUp(device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func SetIp(ip, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func SetMac(mac, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
}

func (t *Interface) ReadPacket() (*Packet, error) {
	buf := make([]byte, 100000)
	err := syscall.ReadFile(t.file, buf, &t.Rl, &t.Rx)
	if err != nil {
		/*Log(ERROR, "Failed to read from TAP device: %s", err)
		return nil, err*/
	}
	if _, err := syscall.WaitForSingleObject(t.Rx.HEvent, syscall.INFINITE); err != nil {
		Log(ERROR, "Failed to read from TAP device: %s", err)
	}
	t.Rx.Offset += t.Rl
	Log(INFO, "1")
	l := 0
	switch buf[0] & 0xf0 {
	case 0x40:
		Log(INFO, "2")
		l = 256*int(buf[2]) + int(buf[3])
	case 0x60:
		Log(INFO, "3")
		continue
		// 40 is ipv6 packet header length
		l = 256*int(buf[4]) + int(buf[5]) + 40
	}
	Log(INFO, "4")
	var pkt *Packet
	pkt = &Packet{Packet: buf[4:l]}
	Log(INFO, "5")
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[2:4]))
	Log(INFO, "6")
	flags := int(*(*uint16)(unsafe.Pointer(&buf[0])))
	Log(INFO, "7")
	if flags&flagTruncated != 0 {
		pkt.Truncated = true
	}
	return pkt, nil
}

func (t *Interface) WritePacket(pkt *Packet) error {
	buf := make([]byte, len(pkt.Packet)+4)
	binary.BigEndian.PutUint16(buf[2:4], uint16(pkt.Protocol))
	copy(buf[4:], pkt.Packet)
	var l uint32
	syscall.WriteFile(t.file, buf, &l, &t.Tx)
	t.Tx.Offset += uint32(len(buf))
	return nil
}

func (t *Interface) Close() error {
	syscall.Close(t.Handle)
	return nil
}

func CheckPermissions() bool {
	return true
}

func Open(ifPattern string, kind DevKind, meta bool) (*Interface, error) {
	inf, err := openDevice(ifPattern)
	if err != nil {
		return nil, err
	}
	return inf, err
}
