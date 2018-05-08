/*
Generated TestP2PService_Execute
Generated TestExecService
*/
package main

import (
	"testing"

	"golang.org/x/sys/windows/svc"
)

func TestP2PService_Execute(t *testing.T) {
	type args struct {
		args    []string
		r       <-chan svc.ChangeRequest
		changes chan<- svc.Status
	}
	tests := []struct {
		name      string
		m         *P2PService
		args      args
		wantSsec  bool
		wantErrno uint32
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &P2PService{}
			gotSsec, gotErrno := m.Execute(tt.args.args, tt.args.r, tt.args.changes)
			if gotSsec != tt.wantSsec {
				t.Errorf("P2PService.Execute() gotSsec = %v, want %v", gotSsec, tt.wantSsec)
			}
			if gotErrno != tt.wantErrno {
				t.Errorf("P2PService.Execute() gotErrno = %v, want %v", gotErrno, tt.wantErrno)
			}
		})
	}
}

func TestExecService(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecService(); (err != nil) != tt.wantErr {
				t.Errorf("ExecService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
