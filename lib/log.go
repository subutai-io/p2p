package ptp

import (
	"log"
	"os"
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
var stdLoggers = [...]*log.Logger{log.New(os.Stdout, logPrefixes[Trace], logFlags[Trace]),
	log.New(os.Stdout, logPrefixes[Debug], logFlags[Debug]),
	log.New(os.Stdout, logPrefixes[Info], logFlags[Info]),
	log.New(os.Stdout, logPrefixes[Warning], logFlags[Warning]),
	log.New(os.Stdout, logPrefixes[Error], logFlags[Error])}

// SetMinLogLevel sets a minimal logging level
func SetMinLogLevel(level LogLevel) {
	logLevelMin = level
}

// MinLogLevel returns minimal log level
func MinLogLevel() LogLevel { return logLevelMin }

// Log writes a log message
func Log(level LogLevel, format string, v ...interface{}) {
	if level < logLevelMin {
		return
	}
	stdLoggers[level].Printf(format, v...)
}
