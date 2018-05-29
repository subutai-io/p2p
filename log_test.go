package ptp

import "testing"

func TestSetMinLogLevel(t *testing.T) {
	var i LogLevel
	for i = 0; i < 10; i++ {
		SetMinLogLevel(i)
		if logLevelMin != i {
			t.Errorf("Error. Wait %v, get %v", i, logLevelMin)
		}
	}
}

func TestMinLogLevel(t *testing.T) {
	var level LogLevel
	for level = 0; level < 10; level++ {
		SetMinLogLevel(level)
		get := MinLogLevel()
		if get != level {
			t.Errorf("Error. Wait %v, get %v", level, get)
		}
	}
}

func TestSetSyslogSocket(t *testing.T) {
	syslogs := [...]string{
		"socket",
		"12345",
		"",
	}
	for i := 0; i < len(syslogs); i++ {
		SetSyslogSocket(syslogs[i])
		if syslogSocket != syslogs[i] {
			t.Errorf("Error. Wait %v, get %v", syslogs[i], syslogSocket)
		}
	}
}
