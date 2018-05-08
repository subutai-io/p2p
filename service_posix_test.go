/*
Generated TestExecService
*/
// +build !windows

package main

import "testing"

func TestExecService(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ExecService(); (err != nil) != tt.wantErr {
				t.Errorf("ExecService() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
