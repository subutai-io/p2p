// +build linux

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
		{"simple test", "vptp"},
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
	}{}
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

	res := &TAPLinux{
		IP:   ip,
		Mac:  hwa,
		Mask: net.IPv4Mask(255, 255, 255, 0),
		MTU:  1500,
		PMTU: false,
	}

	tests := []struct {
		name    string
		args    args
		want    *TAPLinux
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
		want *TAPLinux
	}{
		{"simple test", &TAPLinux{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newEmptyTAP(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newEmptyTAP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty test", fields{}, ""},
		{"normal test", fields{Name: "name"}, "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetName(); got != tt.want {
				t.Errorf("TAPLinux.GetName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   net.HardwareAddr
	}{
		{"empty", fieds{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetHardwareAddress(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.GetHardwareAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IP
	}{
		{"empty", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetIP(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.GetIP() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetSubnet(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   net.IP
	}{
		{"empty", fields{}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetSubnet(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.GetSubnet() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.GetMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetBasename(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"empty", fields{}, "vptp"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetBasename(); got != tt.want {
				t.Errorf("TAPLinux.GetBasename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_SetName(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetName(tt.args.name)
		})
	}
}

func TestTAPLinux_SetHardwareAddress(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetHardwareAddress(tt.args.mac)
		})
	}
}

func TestTAPLinux_SetIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetIP(tt.args.ip)
		})
	}
}

func TestTAPLinux_SetSubnet(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	type args struct {
		subnet net.IP
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetSubnet(tt.args.subnet)
		})
	}
}

func TestTAPLinux_SetMask(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetMask(tt.args.mask)
		})
	}
}

func TestTAPLinux_Init(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
		{"passing", fields{}, args{"name"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.Init(tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_Open(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}

	f0, _ := os.OpenFile("/tmp/t", os.O_RWDR, 0)
	defer f0.Close()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"existing file", fields{file: f0}, true},
		{"create failed", fields{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.Open(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_Close(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}

	f0, _ := os.OpenFile("/tmp/t", os.O_RWDR, 0)

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"nil file descriptor", fields{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.Close(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_Configure(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	type args struct {
		lazy bool
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.Configure(tt.args.lazy); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.Configure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_Deconfigure(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.Deconfigure(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.Deconfigure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_ReadPacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			got, err := tap.ReadPacket()
			if (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.ReadPacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.ReadPacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_checksum(t *testing.T) {
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name string
		args args
		want uint16
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checksum(tt.args.bytes); got != tt.want {
				t.Errorf("checksum() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_handlePacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Packet
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			got, err := tap.handlePacket(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.handlePacket() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TAPLinux.handlePacket() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_WritePacket(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.WritePacket(tt.args.packet); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.WritePacket() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_Run(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.Run()
		})
	}
}

func TestTAPLinux_createInterface(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.createInterface(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.createInterface() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_setMTU(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.setMTU(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.setMTU() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_linkUp(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.linkUp(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.linkUp() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_linkDown(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.linkDown(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.linkDown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_setIP(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.setIP(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.setIP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_setMac(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if err := tap.setMac(); (err != nil) != tt.wantErr {
				t.Errorf("TAPLinux.setMac() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTAPLinux_IsConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.IsConfigured(); got != tt.want {
				t.Errorf("TAPLinux.IsConfigured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_MarkConfigured(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.MarkConfigured()
		})
	}
}

func TestTAPLinux_EnablePMTU(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.EnablePMTU()
		})
	}
}

func TestTAPLinux_DisablePMTU(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.DisablePMTU()
		})
	}
}

func TestTAPLinux_IsPMTUEnabled(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.IsPMTUEnabled(); got != tt.want {
				t.Errorf("TAPLinux.IsPMTUEnabled() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_IsBroken(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.IsBroken(); got != tt.want {
				t.Errorf("TAPLinux.IsBroken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_SetAuto(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	type args struct {
		auto bool
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			tap.SetAuto(tt.args.auto)
		})
	}
}

func TestTAPLinux_IsAuto(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
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
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.IsAuto(); got != tt.want {
				t.Errorf("TAPLinux.IsAuto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTAPLinux_GetStatus(t *testing.T) {
	type fields struct {
		IP         net.IP
		Subnet     net.IP
		Mask       net.IPMask
		Mac        net.HardwareAddr
		Name       string
		Tool       string
		MTU        int
		file       *os.File
		Configured bool
		PMTU       bool
		Auto       bool
		Status     InterfaceStatus
	}
	tests := []struct {
		name   string
		fields fields
		want   InterfaceStatus
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tap := &TAPLinux{
				IP:         tt.fields.IP,
				Subnet:     tt.fields.Subnet,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
				PMTU:       tt.fields.PMTU,
				Auto:       tt.fields.Auto,
				Status:     tt.fields.Status,
			}
			if got := tap.GetStatus(); got != tt.want {
				t.Errorf("TAPLinux.GetStatus() = %v, want %v", got, tt.want)
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
