/*
Generated TestInitPlatform
Generated TestCheckPermissions
Generated TestSyslog
Generated TestSetupPlatform
Generated Test_removeWindowsService
Generated Test_installWindowsService
Generated Test_restartWindowsService
Generated Test_exePath
*/
// +build windows

package ptp

import (
	"testing"

	"golang.org/x/sys/windows/svc/mgr"
)

func TestInitPlatform(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InitPlatform(); (err != nil) != tt.wantErr {
				t.Errorf("InitPlatform() error = %v, wantErr %v", err, tt.wantErr)
			}
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

func Test_removeWindowsService(t *testing.T) {
	type args struct {
		service *mgr.Service
		name    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			removeWindowsService(tt.args.service, tt.args.name)
		})
	}
}

func Test_installWindowsService(t *testing.T) {
	type args struct {
		manager *mgr.Mgr
		name    string
		app     string
		desc    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			installWindowsService(tt.args.manager, tt.args.name, tt.args.app, tt.args.desc)
		})
	}
}

func Test_restartWindowsService(t *testing.T) {
	type args struct {
		service *mgr.Service
		name    string
		noStart bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			restartWindowsService(tt.args.service, tt.args.name, tt.args.noStart)
		})
	}
}

func Test_exePath(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exePath()
			if (err != nil) != tt.wantErr {
				t.Errorf("exePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exePath() = %v, want %v", got, tt.want)
			}
		})
	}
}
