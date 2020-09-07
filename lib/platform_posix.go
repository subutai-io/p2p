// +build !windows

package ptp

import (
	"fmt"
	"log/syslog"
	"os"
)

const (
	MaximumInterfaceNameLength int = 12
)

var syslogLevel = [...]syslog.Priority{syslog.LOG_DEBUG, syslog.LOG_DEBUG, syslog.LOG_INFO, syslog.LOG_WARNING, syslog.LOG_ERR}

// InitPlatform does a platform specific preparation
func InitPlatform() {

}

// CheckPermissions validates platform specific permissions to run TUNTAP utilities
func HavePrivileges(level int) bool {
	if level != 0 {
		Error("P2P cannot run in daemon mode without root privileges")
		return false
	}
	return true
}

func GetPrivilegesLevel() int {
	return os.Getuid()
}

// Syslog provides additional logging to the syslog server
func Syslog(level LogLevel, format string, v ...interface{}) {
	if l3, err := syslog.Dial("udp", syslogSocket, syslogLevel[level], "p2p"); err == nil {
		l3.Write([]byte(fmt.Sprintf(format, v...)))
		l3.Close()
	}
}

// SetupPlatform runs platform specific preparations during p2p daemon creation
func SetupPlatform(remove bool) {
	// Not used on POSIX
}
