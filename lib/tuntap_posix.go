// +build !windows

package ptp

import (
	"encoding/binary"
	"io"
	"os"
	"os/user"
)

type Interface struct {
	Name string
	file *os.File
}

func InitPlatform() {

}

func (t *Interface) ReadPacket() (*Packet, error) {
	buf := make([]byte, 4096)

	n, err := t.file.Read(buf)
	if err != nil {
		return nil, err
	}

	var pkt *Packet
	pkt = &Packet{Packet: buf[0:n]}
	pkt.Protocol = int(binary.BigEndian.Uint16(buf[12:14]))
	/*flags := int(*(*uint16)(unsafe.Pointer(&buf[0])))
	if flags&flagTruncated != 0 {
		Log(Info, "TRUNCATED")
		//pkt.Truncated = true
	}
	*/
	pkt.Truncated = false
	return pkt, nil
}

func (t *Interface) WritePacket(pkt *Packet) error {
	n, err := t.file.Write(pkt.Packet)
	if err != nil {
		return err
	}
	if n != len(pkt.Packet) {
		return io.ErrShortWrite
	}
	return nil
}

func (t *Interface) Close() error {
	return t.file.Close()
}

func CheckPermissions() bool {
	user, err := user.Current()
	if err != nil {
		Log(Error, "Failed to retrieve information about user: %v", err)
		return false
	}
	if user.Uid != "0" {
		Log(Error, "P2P cannot run in daemon mode without root privileges")
		return false
	}
	return true
}

func Open(ifPattern string, kind DevKind) (*Interface, error) {
	file, err := openDevice(ifPattern)
	if err != nil {
		return nil, err
	}

	ifName, err := createInterface(file, ifPattern, kind)
	if err != nil {
		return nil, err
	}

	inf := new(Interface)
	inf.Name = ifName
	inf.file = file

	return inf, nil
}

func (t *Interface) Run() {

	// Dummy, used for windows only

}

func ExtractMacFromInterface(dev *Interface) string {
	return ""
}
