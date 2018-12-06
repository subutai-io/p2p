// +build windows

package ptp

const DefaultConfigLocation = "C:\\ProgramData\\subutai\\bin\\p2p.yaml"

// Platform specific defaults
const (
	DefaultIPTool  = "netsh.exe"                                            // Default network interface configuration tool for Darwin OS
	DefaultTAPTool = "C:\\Program Files\\TAP-Windows\\bin\\tapinstall.exe"  // Default path to TAP configuration tool on Windows OS
	DefaultINFFile = "C:\\Program Files\\TAP-Windows\\driver\\OemVista.inf" // Default path to TAP INF file used by Windows OS
	DefaultMTU     = 1500                                                   // Default MTU value
	DefaultPMTU    = false                                                  // Default PMTU switch
)
