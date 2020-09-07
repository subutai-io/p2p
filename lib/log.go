package ptp

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
)

// LogLevel is a level of the log message
// type LogLevel int32
type LogLevel logrus.Level

// Log Levels
const (
	LTrace LogLevel = iota
	LDebug
	LInfo
	LWarning
	LError
)

var logPrefixes = [...]string{"[TRACE] ", "[DEBUG] ", "[INFO] ", "[WARNING] ", "[ERROR] "}
var logFlags = [...]int{log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime}

var logLevelMin = LInfo
var syslogSocket = ""

/*
var stdLoggers = [...]*log.Logger{log.New(os.Stdout, logPrefixes[Trace], logFlags[Trace]),
	log.New(os.Stdout, logPrefixes[Debug], logFlags[Debug]),
	log.New(os.Stdout, logPrefixes[Info], logFlags[Info]),
	log.New(os.Stdout, logPrefixes[Warning], logFlags[Warning]),
	log.New(os.Stdout, logPrefixes[Error], logFlags[Error])}
*/

// SetMinLogLevel sets a minimal logging level. Accepts a LogLevel constant for setting
func SetMinLogLevel(level LogLevel) {
	logLevelMin = level
}

// SetMinLogLevel sets a minimal logging level. Accepts a string for setting
func SetMinLogLevelString(level string) error {
	level = strings.ToLower(level)
	if level == "trace" {
		SetMinLogLevel(LTrace)
	} else if level == "debug" {
		SetMinLogLevel(LDebug)
	} else if level == "info" {
		SetMinLogLevel(LInfo)
	} else if level == "warning" {
		SetMinLogLevel(LWarning)
	} else if level == "error" {
		SetMinLogLevel(LError)
	} else {
		Warn("Unknown log level %s was provided. Supported log levels are:\ntrace\ndebug\ninfo\nwarning\nerror\n", level)
		return fmt.Errorf("Could not set provided log level")
	}
	Info("Logging level has switched to %s level", level)
	return nil
}

// MinLogLevel returns minimal log level
func MinLogLevel() LogLevel { return logLevelMin }

// Log writes a log message
/*
func Log(level LogLevel, format string, v ...interface{}) {
	if level < logLevelMin {
		return
	}
	stdLoggers[level].Printf(format, v...)
	if level != Trace && len(syslogSocket) != 0 {
		go Syslog(level, format, v...)
	}
}
*/

// SetSyslogSocket sets an adders of the syslog server
/*
func SetSyslogSocket(socket string) {
	syslogSocket = socket
}
*/

func Trace(msg ...interface{}) {
	logrus.Trace(msg)
}

func Debug(msg ...interface{}) {
	logrus.Debug(msg)
}

func Info(msg ...interface{}) {
	logrus.Info(msg...)
}

func Warn(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Warning(msg ...interface{}) {
	logrus.Warn(msg...)
}

func Error(msg ...interface{}) {
	logrus.Error(msg...)
}
