// +build linux

package ptp

const DefaultConfigLocation = "/etc/p2p.yaml"

// Platform specific defaults
const (
	DefaultIPTool  = "/sbin/ip" // Default network interface configuration tool for Darwin OS
	DefaultTAPTool = ""         // Default path to TAP configuration tool on Windows OS
	DefaultINFFile = ""         // Default path to TAP INF file used by Windows OS
	DefaultMTU     = 1500       // Default MTU value
	DefaultPMTU    = false      // Default PMTU switch
)
