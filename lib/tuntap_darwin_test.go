// +build darwin

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
		{"passing", "tun"},
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
		{"ifconfig", "/sbin/ifconfig"},
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
		pmtu bool
	}

	ip := net.ParseIP("10.0.0.1")
	hwa, _ := net.ParseMAC("00:11:22:33:44:55")

	res := &TAPDarwin{
		IP:   ip,
		Mac:  hwa,
		Mask: net.IPv4Mask(255, 255, 255, 0),
		MTU:  1500,
		PMTU: false,
	}

	tests := []struct {
		name    string
		args    args
		want    *TAPDarwin
		wantErr bool
	}{
		{"bad ip", args{ip: "badip"}, nil, true},
		{"bad mac", args{mac: "badmac"}, nil, true},
		{"passing", args{ip: "10.0.0.1", mac: "00:11:22:33:44:55"}, res, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newTAP(tt.args.tool, tt.args.ip, tt.args.mac, tt.args.mask, tt.args.mtu, tt.args.pmtu)
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

func Test_newEmptyTAP(t *testing.T) {
	tests := []struct {
		name string
		want *TAPDarwin
	}{
		{"empty test", &TAPDarwin{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEmptyTAP(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEmptyTAP() = %v, want %v", got, tt.want)
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, ""},
		{"predefined", fields{Name: "test-name"}, "test-name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.GetName(); got != tt.want {
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
		PMTU       bool
	}

	tests := []struct {
		name   string
		fields fields
		want   net.HardwareAddr
	}{
		{"empty", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.GetHardwareAddress(); !reflect.DeepEqual(got, tt.want) {
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IP
	}{
		{"epmty", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.GetIP(); !reflect.DeepEqual(got, tt.want) {
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IPMask
	}{
		{"empty", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.GetMask(); !reflect.DeepEqual(got, tt.want) {
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "tap"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.GetBasename(); got != tt.want {
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
		PMTU       bool
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.SetName(tt.args.name)
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
		PMTU       bool
	}
	type args struct {
		mac net.HardwareAddr
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.SetHardwareAddress(tt.args.mac)
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
		PMTU       bool
	}
	type args struct {
		ip net.IP
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.SetIP(tt.args.ip)
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
		PMTU       bool
	}
	type args struct {
		mask net.IPMask
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"empty", fields{}, args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.SetMask(tt.args.mask)
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
		PMTU       bool
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
		{"empty name", fields{}, args{}, true},
		{"passing", fields{}, args{"test-name"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if err := tap.Init(tt.args.name); (err != nil) != tt.wantErr {
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
		PMTU       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"empty", fields{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if err := tap.Open(); (err != nil) != tt.wantErr {
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
		PMTU       bool
	}

	dsc0, _ := os.OpenFile("/tmp/p2p-tap-close-test-0", os.O_CREATE|os.O_RDWR, 0700)
	dsc1, _ := os.OpenFile("/tmp/p2p-tap-close-test-1", os.O_CREATE|os.O_RDWR, 0700)
	dsc1.Close()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil interface", fields{}, true},
		{"closing closed file", fields{file: dsc1}, true},
		{"passing", fields{file: dsc0}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if err := tap.Close(); (err != nil) != tt.wantErr {
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
		PMTU       bool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"no tool", fields{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if err := tap.Configure(false); (err != nil) != tt.wantErr {
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
		PMTU       bool
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
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			got, err := tap.ReadPacket()
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
		PMTU       bool
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
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if err := tap.WritePacket(tt.args.packet); (err != nil) != tt.wantErr {
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.Run()
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
		PMTU       bool
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
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.IsConfigured(); got != tt.want {
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
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.MarkConfigured()
		})
	}
}

func TestTAPDarwin_EnablePMTU(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.EnablePMTU()
		})
	}
}

func TestTAPDarwin_DisablePMTU(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			tap.DisablePMTU()
		})
	}
}

func TestTAPDarwin_IsPMTUEnabled(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
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
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.IsPMTUEnabled(); got != tt.want {
				t.Errorf("TAPDarwin.IsPMTUEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPDarwin_IsBroken(t *testing.T) {
	type fields struct {
		IP         net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
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
			tap := &TAPDarwin{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
			}
			if got := tap.IsBroken(); got != tt.want {
				t.Errorf("TAPDarwin.IsBroken() = %v, want %v", got, tt.want)
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
