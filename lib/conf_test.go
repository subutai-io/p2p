package ptp

import (
	"io/ioutil"
	"testing"
)

func Test_Conf_Load(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		filepath string
	}

	d1 := []byte("-")
	ioutil.WriteFile("/tmp/test-yaml-Config-p2p-bad", d1, 0777)

	d2 := []byte("iptool: /sbin/ip")
	ioutil.WriteFile("/tmp/test-yaml-Config-p2p-ok", d2, 0777)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"empty filepath", fields{}, args{""}, false},
		{"wrong filepath", fields{}, args{"/"}, true},
		{"bad yaml", fields{}, args{"/tmp/test-yaml-Config-p2p-bad"}, true},
		{"normal yaml", fields{}, args{"/tmp/test-yaml-Config-p2p-ok"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if err := c.Load(tt.args.filepath); (err != nil) != tt.wantErr {
				t.Errorf("Conf.Load() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_Conf_getIPTool(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		preset string
	}

	c1 := new(Conf)
	c1.Load("/")

	f1 := fields{
		IPTool:  c1.IPTool,
		TAPTool: c1.TAPTool,
		INFFile: c1.INFFile,
		MTU:     c1.MTU,
		PMTU:    c1.PMTU,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"empty string test", fields{}, args{}, ""},
		{"default value", f1, args{""}, c1.IPTool},
		{"preset value", f1, args{"preset-val"}, "preset-val"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if got := c.GetIPTool(tt.args.preset); got != tt.want {
				t.Errorf("Conf.getIPTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Conf_getTAPTool(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		preset string
	}

	c1 := new(Conf)
	c1.Load("/")

	f1 := fields{
		IPTool:  c1.IPTool,
		TAPTool: c1.TAPTool,
		INFFile: c1.INFFile,
		MTU:     c1.MTU,
		PMTU:    c1.PMTU,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"empty string test", fields{}, args{}, ""},
		{"default value", f1, args{""}, c1.TAPTool},
		{"preset value", f1, args{"preset-val"}, "preset-val"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if got := c.GetTAPTool(tt.args.preset); got != tt.want {
				t.Errorf("Conf.getTAPTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Conf_getINFFile(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		preset string
	}

	c1 := new(Conf)
	c1.Load("/")

	f1 := fields{
		IPTool:  c1.IPTool,
		TAPTool: c1.TAPTool,
		INFFile: c1.INFFile,
		MTU:     c1.MTU,
		PMTU:    c1.PMTU,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"empty string test", fields{}, args{}, ""},
		{"default value", f1, args{""}, c1.INFFile},
		{"preset value", f1, args{"preset-val"}, "preset-val"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if got := c.GetINFFile(tt.args.preset); got != tt.want {
				t.Errorf("Conf.getINFFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Conf_getMTU(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		preset int
	}

	c1 := new(Conf)
	c1.Load("/")

	f1 := fields{
		IPTool:  c1.IPTool,
		TAPTool: c1.TAPTool,
		INFFile: c1.INFFile,
		MTU:     c1.MTU,
		PMTU:    c1.PMTU,
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{"empty string test", fields{}, args{}, 0},
		{"default value", f1, args{0}, c1.MTU},
		{"preset value", f1, args{256}, 256},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if got := c.GetMTU(tt.args.preset); got != tt.want {
				t.Errorf("Conf.getMTU() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Conf_getPMTU(t *testing.T) {
	type fields struct {
		IPTool  string
		TAPTool string
		INFFile string
		MTU     int
		PMTU    bool
	}
	type args struct {
		preset bool
	}

	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"empty pmtu val", fields{}, args{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Conf{
				IPTool:  tt.fields.IPTool,
				TAPTool: tt.fields.TAPTool,
				INFFile: tt.fields.INFFile,
				MTU:     tt.fields.MTU,
				PMTU:    tt.fields.PMTU,
			}
			if got := c.GetPMTU(); got != tt.want {
				t.Errorf("Conf.getPMTU() = %v, want %v", got, tt.want)
			}
		})
	}
}
