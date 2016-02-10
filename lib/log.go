package ptp

import (
	"log"
	"os"
)

type LOG_LEVEL int32

const (
	TRACE LOG_LEVEL = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

var log_prefixes = [...]string{"[TRACE]", "[DEBUG] ", "[INFO] ", "[WARNING] ", "[ERROR] "}
var log_flags = [...]int{log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime,
	log.Ldate | log.Ltime}

var log_level_min LOG_LEVEL = INFO
var std_loggers = [...]*log.Logger{log.New(os.Stdout, log_prefixes[TRACE], log_flags[TRACE]),
	log.New(os.Stdout, log_prefixes[DEBUG], log_flags[DEBUG]),
	log.New(os.Stdout, log_prefixes[INFO], log_flags[INFO]),
	log.New(os.Stdout, log_prefixes[WARNING], log_flags[WARNING]),
	log.New(os.Stdout, log_prefixes[ERROR], log_flags[ERROR])}

func SetMinLogLevel(level LOG_LEVEL) {
	log_level_min = level
}
func MinLogLevel() LOG_LEVEL { return log_level_min }

func Log(level LOG_LEVEL, format string, v ...interface{}) {
	if level < log_level_min {
		return
	}
	std_loggers[level].Printf(format, v...)
}
