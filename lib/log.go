package ptp

import (
	"log"
	"os"
	"strings"
	"fmt"
)

// LogLevel is a level of the log message
type LogLevel int32

// Log Levels
const (
	Trace LogLevel = iota
	Debug
	Info
	Warning
	Error
)

var logPrefixes = [...]string{"[TRACE] ", "[DEBUG] ", "[INFO] ", "[WARNING] ", "[ERROR] "}
var logFlags = [...]int{log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime}

var logLevelMin = Info
var syslogSocket = ""
var stdLoggers = [...]*log.Logger{log.New(os.Stdout, logPrefixes[Trace], logFlags[Trace]),
	log.New(os.Stdout, logPrefixes[Debug], logFlags[Debug]),
	log.New(os.Stdout, logPrefixes[Info], logFlags[Info]),
	log.New(os.Stdout, logPrefixes[Warning], logFlags[Warning]),
	log.New(os.Stdout, logPrefixes[Error], logFlags[Error])}

// SetMinLogLevel sets a minimal logging level. Accepts a LogLevel constant for setting
func SetMinLogLevel(level LogLevel) {
	logLevelMin = level
}

// SetMinLogLevelString sets a minimal logging level. Accepts a string for setting
func SetMinLogLevelString(level string) error {
	level = strings.ToLower(level)
	if level == "trace" {
		SetMinLogLevel(Trace)
	} else if level == "debug" {
		SetMinLogLevel(Debug)
	} else if level == "info" {
		SetMinLogLevel(Info)
	} else if level == "warning" {
		SetMinLogLevel(Warning)
	} else if level == "error" {
		SetMinLogLevel(Error)
	} else {
		Log(Warning, "Unknown log level %s was provided. Supported log levels are:\ntrace\ndebug\ninfo\nwarning\nerror\n", level)
		return fmt.Errorf("Could not set provided log level")
	}
	Log(Info, "Logging level has switched to %s level", level)
	return nil
}

// MinLogLevel returns minimal log level
func MinLogLevel() LogLevel { return logLevelMin }

// Log writes a log message
func Log(level LogLevel, format string, v ...interface{}) {
	if level < logLevelMin {
		return
	}
	stdLoggers[level].Printf(format, v...)
	if level != Trace && len(syslogSocket) != 0 {
		go Syslog(level, format, v...)
	}
}

// SetSyslogSocket sets an adders of the syslog server
func SetSyslogSocket(socket string) {
	syslogSocket = socket
}
