package ptp

import "testing"

func TestInitErrors(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"init"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			InitErrors()
		})
	}
}
