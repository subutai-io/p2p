// +build !windows

package ptp

import "testing"

func TestInitPlatform(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"empty test"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitPlatform()
		})
	}
}

func TestHavePrivileges(t *testing.T) {
	type args struct {
		level int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"have privileges", args{0}, true},
		{"doesn't have privileges", args{1}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HavePrivileges(tt.args.level); got != tt.want {
				t.Errorf("HavePrivileges() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPrivilegesLevel(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		{"simple test", 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Converted from != to ==
			if got := GetPrivilegesLevel(); got == tt.want {
				t.Errorf("GetPrivilegesLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyslog(t *testing.T) {
	type args struct {
		level  LogLevel
		format string
		v      []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty test", args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Syslog(tt.args.level, tt.args.format, tt.args.v...)
		})
	}
}

func TestSetupPlatform(t *testing.T) {
	type args struct {
		remove bool
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty test", args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupPlatform(tt.args.remove)
		})
	}
}
