package ptp

import (
	"net"
	"os"
	"reflect"
	"testing"
)

func TestGetDeviceBase(t *testing.T) {
	get := GetDeviceBase()
	if get != "vptp" {
		t.Error("Error. Return wrong value")
	}
}

func TestGetConfigurationTool(t *testing.T) {
	get := GetConfigurationTool()
	wait := []string{"/sbin/ip", "/bin/ip", "/usr/bin/ip"}
	for _, w := range wait {
		if get == w {
			return
		}
	}
	t.Error("Error: ", get)
}

func TestNewTAP(t *testing.T) {
	get1, err := newTAP("tool", "", "01:02:03:04:05:06", "255.255.255.0", 1)
	if get1 != nil {
		t.Error(err)
	}
	get2, err2 := newTAP("tool", "192.168.1.1", "-", "255.255.255.0", 1)
	if get2 != nil {
		t.Error(err2)
	}
}

func TestTAPLinux_handlePacket(t *testing.T) {
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
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Packet
		wantErr bool
	}{
		{"Marshal ICMP", fields{}, args{}, nil, true},
		{"Marshal ICMP", fields{}, args{[]byte("This is not a real packet")}, &Packet{24864, []byte("This is not a real packet")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tuntap := &TAPLinux{
				IP:         tt.fields.IP,
				Mask:       tt.fields.Mask,
				Mac:        tt.fields.Mac,
				Name:       tt.fields.Name,
				Tool:       tt.fields.Tool,
				MTU:        tt.fields.MTU,
				file:       tt.fields.file,
				Configured: tt.fields.Configured,
			}
			got, err := tuntap.handlePacket(tt.args.data)
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
