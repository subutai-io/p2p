package ptp

import (
	"bytes"
	"net"
	"reflect"
	"runtime"
	"testing"
	"sync"
)

func TestIsDeviceExists(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	ptp := new(PeerToPeer)
	dev1 := "lo"
	get1 := ptp.IsDeviceExists(dev1)
	if !get1 {
		t.Error("Error. Device exists.")
	}
	dev2 := "device"
	get2 := ptp.IsDeviceExists(dev2)
	if get2 {
		t.Errorf("Error. There no network interface such %v", dev2)
	}
}

func TestGenerateDeviceName(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	p := new(PeerToPeer)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}

func TestIsIPv4(t *testing.T) {
	ip1 := "194.152.36.143"
	ip2 := "2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"
	ip3 := ""
	ptp := new(PeerToPeer)
	wait1 := true
	get1 := ptp.IsIPv4(ip1)
	if get1 != wait1 {
		t.Errorf("Error: wait %v, get %v", wait1, get1)
	}
	wait2 := false
	get2 := ptp.IsIPv4(ip2)
	if get2 != wait2 {
		t.Errorf("Error: wait %v, get %v", wait2, get2)
	}
	wait3 := false
	get3 := ptp.IsIPv4(ip3)
	if get3 != wait3 {
		t.Errorf("Error: wait %v, get %v", wait3, get3)
	}
}

func TestExtractBestDHTRouters(t *testing.T) {
	ptp := new(PeerToPeer)
	wait, err := net.ResolveUDPAddr("udp4", "192.168.11.5:6882")
	if err != nil {
		t.Error("error")
	}
	ptp.Routers = ""
	get := ptp.extractBestDHTRouter()
	if get != nil {
		t.Error("Length of ptp routers is nil")
	}
	ptp.Routers = "192.168.11.5:24,192.168.22.1:22"
	get2 := ptp.extractBestDHTRouter()

	if bytes.EqualFold(get2.IP, wait.IP) && get2.Port != wait.Port && get2.Zone != wait.Zone {
		t.Errorf("Error.Wait %v, get %v", wait, get2)
	}
}

func TestValidateMac(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestValidateInterfaceName(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestParseIntroString(t *testing.T) {
	ptp := new(PeerToPeer)
	hs := new(PeerHandshake)
	hs.Endpoint, _ = net.ResolveUDPAddr("udp4", "192.168.1.1:24")
	get1, err1 := ptp.ParseIntroString("id,ip,mac")
	if get1 != nil {
		t.Error(err1)
	}
	get2, err2 := ptp.ParseIntroString("1,-,127.0.0.1,192.168.1.1")
	if get2 != nil {
		t.Error(err2)
	}
	get3, err3 := ptp.ParseIntroString("1,01:02:03:04:05:06,-,192.168.1.1")
	if get3 != nil {
		t.Error(err3)
	}
	get4, err4 := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,-")
	if get4 != nil {
		t.Error(err4)
	}
	get5, _ := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,192.168.1.1:24")
	if !reflect.DeepEqual(get5.Endpoint.IP, hs.Endpoint.IP) && get5.Endpoint.Port != hs.Endpoint.Port && get5.Endpoint.Zone != hs.Endpoint.Zone {
		t.Error("Error")
	}
}
/*
Generated TestPeerToPeer_AssignInterface
Generated TestPeerToPeer_ListenInterface
Generated TestPeerToPeer_IsDeviceExists -> OK
Generated TestPeerToPeer_GenerateDeviceName -> OK
Generated TestPeerToPeer_IsIPv4 -> OK
Generated TestNew
Generated TestPeerToPeer_ReadDHT
Generated TestPeerToPeer_waitForRemotePort
Generated TestPeerToPeer_extractBestDHTRouter -> OK
Generated TestPeerToPeer_PrepareInterfaces
Generated TestPeerToPeer_attemptPortForward
Generated TestPeerToPeer_Init
Generated TestPeerToPeer_validateMac -> OK
Generated TestPeerToPeer_validateInterfaceName -> OK
Generated TestPeerToPeer_setupHandlers
Generated TestPeerToPeer_RequestIP
Generated TestPeerToPeer_ReportIP
Generated TestPeerToPeer_Run
Generated TestPeerToPeer_checkLastDHTUpdate
Generated TestPeerToPeer_removeStoppedPeers
Generated TestPeerToPeer_checkProxies
Generated TestPeerToPeer_PrepareIntroductionMessage
Generated TestPeerToPeer_WriteToDevice
Generated TestPeerToPeer_ParseIntroString -> OK
Generated TestPeerToPeer_SendTo
Generated TestPeerToPeer_Close
package ptp
*/

/*
func TestIsDeviceExists(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	ptp := new(PeerToPeer)
	dev1 := "lo"
	get1 := ptp.IsDeviceExists(dev1)
	if !get1 {
		t.Error("Error. Device exists.")
	}
	dev2 := "device"
	get2 := ptp.IsDeviceExists(dev2)
	if get2 {
		t.Errorf("Error. There no network interface such %v", dev2)
	}
}

func TestGenerateDeviceName(t *testing.T) {
	if runtime.GOOS == "darwin" {
		return
	}
	p := new(PeerToPeer)
	dev := p.GenerateDeviceName(12)
	if dev != "vptp12" {
		t.Errorf("Device name generation failed. Received %s", dev)
	}
}

func TestIsIPv4(t *testing.T) {
	ip1 := "194.152.36.143"
	ip2 := "2001:0db8:11a3:09d7:1f34:8a2e:07a0:765d"
	ip3 := ""
	ptp := new(PeerToPeer)
	wait1 := true
	get1 := ptp.IsIPv4(ip1)
	if get1 != wait1 {
		t.Errorf("Error: wait %v, get %v", wait1, get1)
	}
	wait2 := false
	get2 := ptp.IsIPv4(ip2)
	if get2 != wait2 {
		t.Errorf("Error: wait %v, get %v", wait2, get2)
	}
	wait3 := false
	get3 := ptp.IsIPv4(ip3)
	if get3 != wait3 {
		t.Errorf("Error: wait %v, get %v", wait3, get3)
	}
}

func TestExtractBestDHTRouters(t *testing.T) {
	ptp := new(PeerToPeer)
	wait, err := net.ResolveUDPAddr("udp4", "192.168.11.5:6882")
	if err != nil {
		t.Error("error")
	}
	ptp.Routers = ""
	get := ptp.extractBestDHTRouter()
	if get != nil {
		t.Error("Length of ptp routers is nil")
	}
	ptp.Routers = "192.168.11.5:24,192.168.22.1:22"
	get2 := ptp.extractBestDHTRouter()

	if bytes.EqualFold(get2.IP, wait.IP) && get2.Port != wait.Port && get2.Zone != wait.Zone {
		t.Errorf("Error.Wait %v, get %v", wait, get2)
	}
}

func TestValidateMac(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestValidateInterfaceName(t *testing.T) {
	ptp := new(PeerToPeer)
	get1 := ptp.validateMac("-")
	if get1 != nil {
		t.Error("Error. Invalid MAC")
	}
	hw, _ := GenerateMAC()
	var h net.HardwareAddr
	get2 := ptp.validateMac(hw)
	if reflect.DeepEqual(get2, h) {
		t.Error("Error")
	}
	get := ptp.validateMac("")
	if reflect.DeepEqual(get, h) {
		t.Error("Error")
	}
}

func TestParseIntroString(t *testing.T) {
	ptp := new(PeerToPeer)
	hs := new(PeerHandshake)
	hs.Endpoint, _ = net.ResolveUDPAddr("udp4", "192.168.1.1:24")
	get1, err1 := ptp.ParseIntroString("id,ip,mac")
	if get1 != nil {
		t.Error(err1)
	}
	get2, err2 := ptp.ParseIntroString("1,-,127.0.0.1,192.168.1.1")
	if get2 != nil {
		t.Error(err2)
	}
	get3, err3 := ptp.ParseIntroString("1,01:02:03:04:05:06,-,192.168.1.1")
	if get3 != nil {
		t.Error(err3)
	}
	get4, err4 := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,-")
	if get4 != nil {
		t.Error(err4)
	}
	get5, _ := ptp.ParseIntroString("1,01:02:03:04:05:06,127.0.0.1,192.168.1.1:24")
	if !reflect.DeepEqual(get5.Endpoint.IP, hs.Endpoint.IP) && get5.Endpoint.Port != hs.Endpoint.Port && get5.Endpoint.Zone != hs.Endpoint.Zone {
		t.Error("Error")
	}
}
*/

func TestPeerToPeer_AssignInterface(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		interfaceName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if err := p.AssignInterface(tt.args.interfaceName); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.AssignInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_ListenInterface(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.ListenInterface()
		})
	}
}

func TestPeerToPeer_IsDeviceExists(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.IsDeviceExists(tt.args.name); got != tt.want {
				t.Errorf("PeerToPeer.IsDeviceExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_GenerateDeviceName(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		i int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.GenerateDeviceName(tt.args.i); got != tt.want {
				t.Errorf("PeerToPeer.GenerateDeviceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_IsIPv4(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		ip string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.IsIPv4(tt.args.ip); got != tt.want {
				t.Errorf("PeerToPeer.IsIPv4() = %v, want %v", got, tt.want)
			}
		})
	}
}

/*
func TestNew(t *testing.T) {
	type args struct {
		argIP      string
		argMac     string
		argDev     string
		argDirect  string
		argHash    string
		argKeyfile string
		argKey     string
		argTTL     string
		argLog     string
		fwd        bool
		port       int
		ignoreIPs  []string
		outboundIP net.IP
		routers    string
	}
	tests := []struct {
		name string
		args args
		want *PeerToPeer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.argIP, tt.args.argMac, tt.args.argDev, tt.args.argDirect, tt.args.argHash, tt.args.argKeyfile, tt.args.argKey, tt.args.argTTL, tt.args.argLog, tt.args.fwd, tt.args.port, tt.args.ignoreIPs, tt.args.outboundIP, tt.args.routers); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/

func TestPeerToPeer_ReadDHT(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.ReadDHT()
		})
	}
}

func TestPeerToPeer_waitForRemotePort(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.waitForRemotePort()
		})
	}
}

func TestPeerToPeer_extractBestDHTRouter(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
		want   *net.UDPAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.extractBestDHTRouter(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.extractBestDHTRouter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_PrepareInterfaces(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		ip            string
		interfaceName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if err := p.PrepareInterfaces(tt.args.ip, tt.args.interfaceName); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.PrepareInterfaces() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_attemptPortForward(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		port uint16
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if err := p.attemptPortForward(tt.args.port, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.attemptPortForward() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPeerToPeer_Init(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.Init()
		})
	}
}

func TestPeerToPeer_validateMac(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		mac string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   net.HardwareAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.validateMac(tt.args.mac); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.validateMac() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_validateInterfaceName(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			got, err := p.validateInterfaceName(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.validateInterfaceName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeerToPeer.validateInterfaceName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_setupHandlers(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.setupHandlers()
		})
	}
}

func TestPeerToPeer_RequestIP(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		mac    string
		device string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    net.IP
		want1   net.IPMask
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			got, got1, err := p.RequestIP(tt.args.mac, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.RequestIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.RequestIP() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PeerToPeer.RequestIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPeerToPeer_ReportIP(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		ipAddress string
		mac       string
		device    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    net.IP
		want1   net.IPMask
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			got, got1, err := p.ReportIP(tt.args.ipAddress, tt.args.mac, tt.args.device)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.ReportIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.ReportIP() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PeerToPeer.ReportIP() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPeerToPeer_Run(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.Run()
		})
	}
}

func TestPeerToPeer_checkLastDHTUpdate(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.checkLastDHTUpdate()
		})
	}
}

func TestPeerToPeer_removeStoppedPeers(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.removeStoppedPeers()
		})
	}
}

func TestPeerToPeer_checkProxies(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.checkProxies()
		})
	}
}

func TestPeerToPeer_PrepareIntroductionMessage(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		id       string
		endpoint string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *P2PMessage
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if got := p.PrepareIntroductionMessage(tt.args.id, tt.args.endpoint); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.PrepareIntroductionMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_WriteToDevice(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		b         []byte
		proto     uint16
		truncated bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			p.WriteToDevice(tt.args.b, tt.args.proto, tt.args.truncated)
		})
	}
}

func TestPeerToPeer_ParseIntroString(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		intro string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *PeerHandshake
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			got, err := p.ParseIntroString(tt.args.intro)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.ParseIntroString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PeerToPeer.ParseIntroString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_SendTo(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	type args struct {
		dst net.HardwareAddr
		msg *P2PMessage
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			got, err := p.SendTo(tt.args.dst, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.SendTo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("PeerToPeer.SendTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPeerToPeer_Close(t *testing.T) {
	type fields struct {
		Config          Configuration
		UDPSocket       *Network
		LocalIPs        []net.IP
		Dht             *DHTClient
		Crypter         Crypto
		Shutdown        bool
		ForwardMode     bool
		ReadyToStop     bool
		MessageHandlers map[uint16]MessageHandler
		PacketHandlers  map[PacketType]PacketHandlerCallback
		PeersLock       sync.Mutex
		Hash            string
		Routers         string
		Interface       TAP
		Peers           *PeerList
		HolePunching    sync.Mutex
		ProxyManager    *ProxyManager
		outboundIP      net.IP
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &PeerToPeer{
				Config:          tt.fields.Config,
				UDPSocket:       tt.fields.UDPSocket,
				LocalIPs:        tt.fields.LocalIPs,
				Dht:             tt.fields.Dht,
				Crypter:         tt.fields.Crypter,
				Shutdown:        tt.fields.Shutdown,
				ForwardMode:     tt.fields.ForwardMode,
				ReadyToStop:     tt.fields.ReadyToStop,
				MessageHandlers: tt.fields.MessageHandlers,
				PacketHandlers:  tt.fields.PacketHandlers,
				PeersLock:       tt.fields.PeersLock,
				Hash:            tt.fields.Hash,
				Routers:         tt.fields.Routers,
				Interface:       tt.fields.Interface,
				Peers:           tt.fields.Peers,
				HolePunching:    tt.fields.HolePunching,
				ProxyManager:    tt.fields.ProxyManager,
				outboundIP:      tt.fields.outboundIP,
			}
			if err := p.Close(); (err != nil) != tt.wantErr {
				t.Errorf("PeerToPeer.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
