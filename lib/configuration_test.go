/*
Generated TestConfiguration_GetIPTool
Generated TestConfiguration_GetAddTap
Generated TestConfiguration_GetInfFile
Generated TestConfiguration_Read
*/

package ptp

import "testing"

func TestConfiguration_GetIPTool(t *testing.T) {
	type fields struct {
		IPTool  string
		AddTap  string
		InfFile string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &Configuration{
				IPTool:  tt.fields.IPTool,
				AddTap:  tt.fields.AddTap,
				InfFile: tt.fields.InfFile,
			}
			if got := y.GetIPTool(); got != tt.want {
				t.Errorf("Configuration.GetIPTool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_GetAddTap(t *testing.T) {
	type fields struct {
		IPTool  string
		AddTap  string
		InfFile string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &Configuration{
				IPTool:  tt.fields.IPTool,
				AddTap:  tt.fields.AddTap,
				InfFile: tt.fields.InfFile,
			}
			if got := y.GetAddTap(); got != tt.want {
				t.Errorf("Configuration.GetAddTap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_GetInfFile(t *testing.T) {
	type fields struct {
		IPTool  string
		AddTap  string
		InfFile string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			y := &Configuration{
				IPTool:  tt.fields.IPTool,
				AddTap:  tt.fields.AddTap,
				InfFile: tt.fields.InfFile,
			}
			if got := y.GetInfFile(); got != tt.want {
				t.Errorf("Configuration.GetInfFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfiguration_Read(t *testing.T) {
	type fields struct {
		IPTool  string
		AddTap  string
		InfFile string
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
			y := &Configuration{
				IPTool:  tt.fields.IPTool,
				AddTap:  tt.fields.AddTap,
				InfFile: tt.fields.InfFile,
			}
			if err := y.Read(); (err != nil) != tt.wantErr {
				t.Errorf("Configuration.Read() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
