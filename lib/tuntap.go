package ptp

import (
	"errors"
	"net"
)

const (
	flagMF      = 0x10
	flagDF      = 0x1
	iffTun      = 0x1
	iffTap      = 0x2
	iffOneQueue = 0x2000
	iffnopi     = 0x1000
)

var (
	errPacketTooBig      = errors.New("Packet exceeds MTU")
	errICMPMarshalFailed = errors.New("Failed to marshal ICMP")
	errPacketTooSmall    = errors.New("Packet is too small")
)

type ifReq struct {
	Name  [0x10]byte
	Flags uint16
	pad   [0x28 - 0x10 - 2]byte
}

// Packet represents a packet received on TUN/TAP interface
type Packet struct {
	Protocol int
	Packet   []byte
}

// InterfaceStatus holds Status of the network Interface
type InterfaceStatus uint8

// Interface Statuses
const (
	InterfaceWaiting      InterfaceStatus = 0
	InterfaceConfiguring  InterfaceStatus = 1
	InterfaceConfigured   InterfaceStatus = 2
	InterfaceDeconfigured InterfaceStatus = 3
	InterfaceRunning      InterfaceStatus = 4
	InterfaceBroken       InterfaceStatus = 5
	InterfaceShutdown     InterfaceStatus = 6
)

// TAP interface
type TAP interface {
	GetName() string
	GetHardwareAddress() net.HardwareAddr
	GetIP() net.IP
	GetSubnet() net.IP
	GetMask() net.IPMask
	GetBasename() string
	SetName(string)
	SetHardwareAddress(net.HardwareAddr)
	SetIP(net.IP)
	SetSubnet(net.IP)
	SetMask(net.IPMask)
	Init(string) error
	Open() error
	Close() error
	Configure(bool) error
	Deconfigure() error
	ReadPacket() (*Packet, error)
	WritePacket(*Packet) error
	Run()
	IsConfigured() bool
	MarkConfigured()
	EnablePMTU()
	DisablePMTU()
	IsPMTUEnabled() bool
	IsBroken() bool
	SetAuto(bool)
	IsAuto() bool
	GetStatus() InterfaceStatus
}
