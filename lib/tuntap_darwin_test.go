/*
Generated TestGetDeviceBase
Generated TestGetConfigurationTool
Generated Test_newTAP
Generated TestTAPDarwin_GetName
Generated TestTAPDarwin_GetHardwareAddress
Generated TestTAPDarwin_GetIP
Generated TestTAPDarwin_GetMask
Generated TestTAPDarwin_GetBasename
Generated TestTAPDarwin_SetName
Generated TestTAPDarwin_SetHardwareAddress
Generated TestTAPDarwin_SetIP
Generated TestTAPDarwin_SetMask
Generated TestTAPDarwin_Init
Generated TestTAPDarwin_Open
Generated TestTAPDarwin_Close
Generated TestTAPDarwin_Configure
Generated TestTAPDarwin_ReadPacket
Generated TestTAPDarwin_WritePacket
Generated TestTAPDarwin_Run
Generated TestTAPDarwin_IsConfigured
Generated TestTAPDarwin_MarkConfigured
Generated TestFilterInterface
*/

package ptp

import (
	"net"
	"os"
	"reflect"
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
		want    *TAPDarwin
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

func TestTAPDarwin_GetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.GetName(); got != tt.want {
				t.Errorf("TAPDarwin.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_GetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.GetHardwareAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPDarwin.GetHardwareAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_GetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.GetIP(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPDarwin.GetIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_GetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.GetMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPDarwin.GetMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_GetBasename(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.GetBasename(); got != tt.want {
				t.Errorf("TAPDarwin.GetBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_SetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.SetName(tt.args.name)
		})
	}
}

func TestTAPDarwin_SetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.SetHardwareAddress(tt.args.mac)
		})
	}
}

func TestTAPDarwin_SetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.SetIP(tt.args.ip)
		})
	}
}

func TestTAPDarwin_SetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.SetMask(tt.args.mask)
		})
	}
}

func TestTAPDarwin_Init(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if err := t.Init(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPDarwin_Open(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if err := t.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPDarwin_Close(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if err := t.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPDarwin_Configure(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if err := t.Configure(); (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPDarwin_ReadPacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			got, err := t.ReadPacket()
			if (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.ReadPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPDarwin.ReadPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_WritePacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
	}
	type args struct {
		packet *Packet
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if err := t.WritePacket(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("TAPDarwin.WritePacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPDarwin_Run(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.Run()
		})
	}
}

func TestTAPDarwin_IsConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			if got := t.IsConfigured(); got != tt.want {
				t.Errorf("TAPDarwin.IsConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_MarkConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
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
			t := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			t.MarkConfigured()
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
