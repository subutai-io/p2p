/*
Generated TestGetDeviceBase
Generated TestGetConfigurationTool
Generated Test_newTAP
Generated TestTAPWindows_GetName
Generated TestTAPWindows_GetHardwareAddress
Generated TestTAPWindows_GetIP
Generated TestTAPWindows_GetMask
Generated TestTAPWindows_GetBasename
Generated TestTAPWindows_SetName
Generated TestTAPWindows_SetHardwareAddress
Generated TestTAPWindows_SetIP
Generated TestTAPWindows_SetMask
Generated TestTAPWindows_Init
Generated TestTAPWindows_Open
Generated TestTAPWindows_Close
Generated TestTAPWindows_Configure
Generated TestTAPWindows_Run
Generated TestTAPWindows_ReadPacket
Generated TestTAPWindows_WritePacket
Generated TestTAPWindows_read
Generated TestTAPWindows_write
Generated TestTAPWindows_queryNetworkKey
Generated TestTAPWindows_queryAdapters
Generated TestTAPWindows_removeZeroes
Generated TestTAPWindows_IsConfigured
Generated TestTAPWindows_MarkConfigured
Generated Test_tapControlCode
Generated Test_controlCode
Generated TestFilterInterface
*/
// +build windows

package ptp

import (
	"net"
	"reflect"
	"syscall"
	"testing"
)

func TestGetDeviceBase(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDeviceBase(); got != tt.want {
				t.Errorf("GetDeviceBase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetConfigurationTool(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetConfigurationTool(); got != tt.want {
				t.Errorf("GetConfigurationTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_newTAP(t *testing.T) {
	type args struct {
		tool string
		ip   string
		mac  string
		mask string
		mtu  int
	}
	tests := []struct {
		name    string
		args    args
		want    *TAPWindows
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTAP(tt.args.tool, tt.args.ip, tt.args.mac, tt.args.mask, tt.args.mtu)
			if (err != nil) != tt.wantErr {
				t.Errorf("newTAP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newTAP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_GetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.GetName(); got != tt.want {
				t.Errorf("TAPWindows.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_GetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   net.HardwareAddr
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.GetHardwareAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPWindows.GetHardwareAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_GetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IP
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.GetIP(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPWindows.GetIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_GetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IPMask
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.GetMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPWindows.GetMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_GetBasename(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.GetBasename(); got != tt.want {
				t.Errorf("TAPWindows.GetBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_SetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		name string
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.SetName(tt.args.name)
		})
	}
}

func TestTAPWindows_SetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		mac net.HardwareAddr
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.SetHardwareAddress(tt.args.mac)
		})
	}
}

func TestTAPWindows_SetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		ip net.IP
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.SetIP(tt.args.ip)
		})
	}
}

func TestTAPWindows_SetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		mask net.IPMask
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.SetMask(tt.args.mask)
		})
	}
}

func TestTAPWindows_Init(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.Init(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_Open(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_Close(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_Configure(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.Configure(); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_Run(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.Run()
		})
	}
}

func TestTAPWindows_ReadPacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Packet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			got, err := t.ReadPacket()
			if (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.ReadPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPWindows.ReadPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_WritePacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		pkt *Packet
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.WritePacket(tt.args.pkt); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.WritePacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_read(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		ch chan []byte
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.read(tt.args.ch); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_write(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		ch chan []byte
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.write(tt.args.ch); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_queryNetworkKey(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name    string
		fields  fields
		want    syscall.Handle
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			got, err := t.queryNetworkKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.queryNetworkKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPWindows.queryNetworkKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_queryAdapters(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		handle syscall.Handle
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if err := t.queryAdapters(tt.args.handle); (err != nil) != tt.wantErr {
				t.Errorf("TAPWindows.queryAdapters() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPWindows_removeZeroes(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	type args struct {
		s string
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
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.removeZeroes(tt.args.s); got != tt.want {
				t.Errorf("TAPWindows.removeZeroes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_IsConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			if got := t.IsConfigured(); got != tt.want {
				t.Errorf("TAPWindows.IsConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPWindows_MarkConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		MacNotSet  bool
		Name       string
		Interface  string
		Tool       string
		MTU        int
		file       syscall.Handle
		Handle     syscall.Handle
		Rx         chan []byte
		Tx         chan []byte
		Configured bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t := &TAPWindows{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				MacNotSet:  tt.fields.MacNotSet,
				Name:       tt.fields.Name,
				Interface:  tt.fields.Interface,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Handle:     tt.fields.Handle,
				Rx:         tt.fields.Rx,
				Tx:         tt.fields.Tx,
				Configured: tt.fields.Configured,
			}
			t.MarkConfigured()
		})
	}
}

func Test_tapControlCode(t *testing.T) {
	type args struct {
		request uint32
		method  uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tapControlCode(tt.args.request, tt.args.method); got != tt.want {
				t.Errorf("tapControlCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_controlCode(t *testing.T) {
	type args struct {
		device_type uint32
		function    uint32
		method      uint32
		access      uint32
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := controlCode(tt.args.device_type, tt.args.function, tt.args.method, tt.args.access); got != tt.want {
				t.Errorf("controlCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilterInterface(t *testing.T) {
	type args struct {
		infName string
		infIP   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FilterInterface(tt.args.infName, tt.args.infIP); got != tt.want {
				t.Errorf("FilterInterface() = %v, want %v", got, tt.want)
			}
		})
	}
}
