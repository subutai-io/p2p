
// +build !windows

/*
Generated TestInitPlatform
Generated TestCheckPermissions
Generated TestSyslog
Generated TestSetupPlatform
*/

package ptp

import "testing"

func TestInitPlatform(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitPlatform()
		})
	}
}

func TestCheckPermissions(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPermissions(); got != tt.want {
				t.Errorf("CheckPermissions() = %v, want %v", got, tt.want)
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetupPlatform(tt.args.remove)
		})
	}
}
