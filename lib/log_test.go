package ptp

import "testing"

// import "testing"

// func TestSetMinLogLevel(t *testing.T) {
// 	var i LogLevel
// 	for i = 0; i < 10; i++ {
// 		SetMinLogLevel(i)
// 		if logLevelMin != i {
// 			t.Errorf("Error. Wait %v, get %v", i, logLevelMin)
// 		}
// 	}
// }

// func TestMinLogLevel(t *testing.T) {
// 	var level LogLevel
// 	for level = 0; level < 10; level++ {
// 		SetMinLogLevel(level)
// 		get := MinLogLevel()
// 		if get != level {
// 			t.Errorf("Error. Wait %v, get %v", level, get)
// 		}
// 	}
// }

// func TestSetSyslogSocket(t *testing.T) {
// 	syslogs := [...]string{
// 		"socket",
// 		"12345",
// 		"",
// 	}
// 	for i := 0; i < len(syslogs); i++ {
// 		SetSyslogSocket(syslogs[i])
// 		if syslogSocket != syslogs[i] {
// 			t.Errorf("Error. Wait %v, get %v", syslogs[i], syslogSocket)
// 		}
// 	}
// }

func TestSetMinLogLevel(t *testing.T) {
	type args struct {
		level LogLevel
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty test", args{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetMinLogLevel(tt.args.level)
		})
	}
}

func TestSetMinLogLevelString(t *testing.T) {
	type args struct {
		level string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"trace", args{"trace"}, false},
		{"debug", args{"debug"}, false},
		{"info", args{"info"}, false},
		{"warning", args{"warning"}, false},
		{"error", args{"error"}, false},
		{"unknown", args{"unknown"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetMinLogLevelString(tt.args.level); (err != nil) != tt.wantErr {
				t.Errorf("SetMinLogLevelString() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMinLogLevel(t *testing.T) {
	SetMinLogLevel(Info)
	tests := []struct {
		name string
		want LogLevel
	}{
		{"empty test", 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MinLogLevel(); got != tt.want {
				t.Errorf("MinLogLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLog(t *testing.T) {
	type args struct {
		level  LogLevel
		format string
		v      []interface{}
	}
	syslogSocket = "test"
	tests := []struct {
		name string
		args args
	}{
		{"lower level", args{Trace, "%s", nil}},
		{"normal level", args{Info, "%s", nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Log(tt.args.level, tt.args.format, tt.args.v...)
		})
	}
}

func TestSetSyslogSocket(t *testing.T) {
	type args struct {
		socket string
	}
	tests := []struct {
		name string
		args args
	}{
		{"empty test", args{""}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetSyslogSocket(tt.args.socket)
		})
	}
}
