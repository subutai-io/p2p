// +build windows
package ptp

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"unicode/utf16"
)

type TAPDevice struct {
	Handle    syscall.Handle
	Name      string
	Interface string
	IP        string
	Mask      string
}

const (
	NETWORK_KEY         string        = "SYSTEM\\CurrentControlSet\\Control\\Network\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	ADAPTER_KEY         string        = "SYSTEM\\CurrentControlSet\\Control\\Class\\{4D36E972-E325-11CE-BFC1-08002BE10318}"
	NO_MORE_ITEMS       syscall.Errno = 259
	USERMODE_DEVICE_DIR string        = "\\\\.\\Global\\"
	SYS_DEVICE_DIR      string        = "\\Device\\"
	USER_DEVICE_DIR     string        = "\\DosDevices\\Global\\"
	TAP_SUFFIX          string        = ".tap"
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

func queryAdapters(handle syscall.Handle) (TAPDevice, error) {
	var dev TAPDevice
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
		dev.Handle, err = syscall.CreateFile(syscall.StringToUTF16Ptr(tapname),
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
		dev.Name = adapterId
		dev.Interface = adapterName
		return dev, nil
	}
	return dev, nil
}

func openDevice(ifPattern string) (*os.File, error) {
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
	setip := exec.Command("netsh")
	setip.SysProcAttr = &syscall.SysProcAttr{}
	setip.SysProcAttr.CmdLine = fmt.Sprintf(`netsh interface ip set address "%s" static %s %s`, dev.Interface, dev.IP, dev.Mask)
	err = setip.Run()
	if err != nil {
		Log(ERROR, "Failed to properly configure TAP device with netsh: %v", err)
		return nil, err
	}

	in := []byte("\x01\x00\x00\x00")
	var length uint32
	err = syscall.DeviceIoControl(dev.Handle, TAP_IOCTL_SET_MEDIA_STATUS,
		&in[0],
		uint32(len(in)),
		&in[0],
		uint32(len(in)),
		&length,
		nil)
	if err != nil {
		Log(ERROR, "Failed to change device status to 'connected': %v", err)
		return nil, err
	}

	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	return file, err
}

func createInterface(file *os.File, ifPattern string, kind DevKind, meta bool) (string, error) {
	panic("TUN/TAP functionality is not supported on this platform")
}

func ConfigureInterface(ip, mac, device, tool string) error {
	panic("TUN/TAP functionality is not supported on this platform")
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

func (t *tun) Read(ch chan []byte) (err error) {
	overlappedRx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	overlappedRx.HEvent = syscall.Handle(hevent)
	buf := make([]byte, t.mtu)
	var l uint32
	for {
		if err := syscall.ReadFile(t.fd, buf, &l, &overlappedRx); err != nil {
		}
		if _, err := syscall.WaitForSingleObject(overlappedRx.HEvent, syscall.INFINITE); err != nil {
			fmt.Println(err)
		}
		overlappedRx.Offset += l
		totalLen := 0
		switch buf[0] & 0xf0 {
		case 0x40:
			totalLen = 256*int(buf[2]) + int(buf[3])
		case 0x60:
			continue
			totalLen = 256*int(buf[4]) + int(buf[5]) + IPv6_HEADER_LENGTH
		}
		fmt.Println("read data", buf[:totalLen])
		send := make([]byte, totalLen)
		copy(send, buf)
		ch <- send
	}
}

func (t *tun) Write(ch chan []byte) (err error) {
	overlappedRx := syscall.Overlapped{}
	var hevent windows.Handle
	hevent, err = windows.CreateEvent(nil, 0, 0, nil)
	if err != nil {
		return
	}
	overlappedRx.HEvent = syscall.Handle(hevent)
	for {
		select {
		case data := <-ch:
			var l uint32
			syscall.WriteFile(t.fd, data, &l, &overlappedRx)
			syscall.WaitForSingleObject(overlappedRx.HEvent, syscall.INFINITE)
			overlappedRx.Offset += uint32(len(data))
		}
	}
}
