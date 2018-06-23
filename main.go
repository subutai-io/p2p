package main

import (
	"fmt"
	"net"
	"os"
	"runtime/pprof"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
	"github.com/urfave/cli"
)

// These variables must be customized at build time

// AppVersion is a Version of P2P
var AppVersion = "Unknown"

// BuildID usually holds output of `git describe`
var BuildID = "Unknown"

// TargetURL will point p2p to specified service under default domain for SRV lookup
var TargetURL = "dht"

// DefaultLog is used when it was not specified during build
var DefaultLog = "INFO"

// InterfaceNames - List of all interfaces names that was used by p2p historically. These interfaces may not present in the system anymore
var InterfaceNames []string

// OutboundIP is an outbound IP address detected by STUN
var OutboundIP net.IP

var SignalChannel chan os.Signal

var ReadyToServe bool

var StartTime time.Time

// StartProfiling will create a .prof file to analyze p2p app performance
func StartProfiling(profile string) {
	pwd, err := os.Getwd()
	if err != nil {
		ptp.Log(ptp.Error, "Getwd() error : %v", err)
		return
	}

	timeStr := "cpu"
	if profile == "cpu" {
		fileName := fmt.Sprintf("%s/%s.prof", pwd, timeStr)
		f, err := os.Create(fileName)
		if err != nil {
			ptp.Log(ptp.Error, "Create cpu_prof file failed. %v", err)
			return
		}
		ptp.Log(ptp.Info, "Start cpu profiling to file %s", fileName)
		pprof.StartCPUProfile(f)
	} else if profile == "memory" {
		_, err := os.Create(fmt.Sprintf("%s/%s.p2p_mem_prof", pwd, timeStr))
		if err != nil {
			ptp.Log(ptp.Error, "Create mem_prof file failed. %v", err)
			return
		}
	}
}

func main() {
	// Command-line flags
	var (
		SaveFile       string // Save file where p2p will store data about instances
		RPCPort        int    // Port that p2p is daemon is listening to
		Profiling      string // Profiling type
		Syslog         string // Syslog socket
		Infohash       string // Infohash of a swarm
		IP             string // IP address of local p2p interface
		Mac            string // Hardware address of p2p interface
		InterfaceName  string // Name of p2p interface
		Keyfile        string // Path to a file with crypto key
		Key            string // AES key
		Until          string // Until date this key will be active in Unix timestamp
		Ports          string // Ports range for an instance
		UDPPort        int    // Specific UDP port for an instance
		UseForwarders  bool   // Whether or not p2p should force usage of proxy servers for this instance
		ShowInterfaces bool   // Whether or not p2p show command should return information about interfaces in use
		ShowAll        bool   //
		ShowBind       bool   // used with show --interfaces
		LogLevel       string // Log level
		RemoveService  bool   // If yes - service will be removed (used with service)
		InstallService bool   // If yes - service will be installed (used with service)
		MTU            int    // MTU for p2p interface
		ShowMTU        bool   // Show MTU value
		PMTU           bool   // Whether or not PMTU capabilities should be used
	)

	app := cli.NewApp()
	app.Name = "p2p"
	app.Version = AppVersion
	app.Authors = []cli.Author{
		cli.Author{
			Name: "subutai.io",
		},
	}
	app.Description = "Subutai P2P creates private mesh network used by PeerOS. Visit https://subutai.io for more information. " +
		"To get help visit our Slack at https://slack.subutai.io/ or create an issue on GitHub: https://github.com/subutai-io/p2p"
	app.Usage = "Subutai P2P daemon/client application"
	app.Copyright = "Copyright 2018 Subutai.io"

	app.Commands = []cli.Command{
		{
			Name:  "daemon",
			Usage: "Run p2p in daemon mode",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "target",
					Usage:       "Comma-separated list of endpoints",
					Value:       TargetURL,
					Destination: &TargetURL,
				},
				cli.StringFlag{
					Name:        "save",
					Usage:       "Path to save/restore instance data file",
					Value:       "",
					Destination: &SaveFile,
				},
				cli.StringFlag{
					Name:        "profile",
					Usage:       "Run p2p in profiling mode. Possible value: mem, cpu",
					Value:       "",
					Destination: &Profiling,
				},
				cli.StringFlag{
					Name:        "syslog",
					Usage:       "Specify syslog socket",
					Value:       "",
					Destination: &Syslog,
				},
				cli.IntFlag{
					Name:        "mtu",
					Usage:       "Specify global MTU value that will be set on p2p interfaces",
					Value:       ptp.DefaultMTU,
					Destination: &MTU,
				},
				cli.StringFlag{
					Name:        "log",
					Usage:       "Log level. Available levels: trace, debug, info, warning, error",
					Value:       "",
					Destination: &LogLevel,
				},
				cli.BoolFlag{
					Name:        "pmtu",
					Usage:       "When specified - enables PMTU capabilities",
					Destination: &PMTU,
				},
			},
			Action: func(c *cli.Context) error {
				ExecDaemon(RPCPort, TargetURL, SaveFile, Profiling, Syslog, LogLevel, MTU, PMTU)
				return nil
			},
		},
		{
			Name:  "service",
			Usage: "[Windows Only] Run Windows Service",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:        "install",
					Usage:       "If set - windows service will be installed",
					Destination: &InstallService,
				},
				cli.BoolFlag{
					Name:        "remove",
					Usage:       "If set - service will be removed if it's already present in the system",
					Destination: &RemoveService,
				},
			},
			Action: func(c *cli.Context) error {
				if InstallService {
					ptp.SetupPlatform(false)
					return nil
				}
				if RemoveService {
					ptp.SetupPlatform(true)
					return nil
				}
				return ExecService()
			},
		},
		{
			Name:  "start",
			Usage: "Start new p2p instance",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "hash",
					Usage:       "Infohash of p2p swarm",
					Value:       "",
					Destination: &Infohash,
				},
				cli.StringFlag{
					Name:        "ip",
					Usage:       "IP Address of p2p interface. Can be specified in CIDR format or use \"dhcp\" to pick free unused IP",
					Value:       "dhcp",
					Destination: &IP,
				},
				cli.StringFlag{
					Name:        "mac",
					Usage:       "Hardware address of a p2p interface",
					Value:       "",
					Destination: &Mac,
				},
				cli.StringFlag{
					Name:        "dev",
					Usage:       "Name of the p2p interface",
					Value:       "",
					Destination: &InterfaceName,
				},
				// cli.StringFlag{
				// 	Name:        "dht",
				// 	Usage:       "[Deprecated] Comman-separated list of DHT bootstrap nodes",
				// 	Value:       "",
				// 	Destination: &DHTRouters,
				// },
				cli.StringFlag{
					Name:        "keyfile",
					Usage:       "Path to a file containing crypto-key",
					Value:       "",
					Destination: &Keyfile,
				},
				cli.StringFlag{
					Name:        "key",
					Usage:       "AES crypto key",
					Value:       "",
					Destination: &Key,
				},
				cli.StringFlag{
					Name:        "ttl, until",
					Usage:       "Time until specified key will be active",
					Value:       "",
					Destination: &Until,
				},
				cli.StringFlag{
					Name:        "ports",
					Usage:       "Ports range that should be used by p2p in a START-END format",
					Value:       "",
					Destination: &Ports,
				},
				cli.IntFlag{
					Name:        "port",
					Usage:       "UDP port for current instance",
					Value:       0,
					Destination: &UDPPort,
				},
				cli.BoolFlag{
					Name:        "fwd",
					Usage:       "Force proxy servers usage",
					Destination: &UseForwarders,
				},
			},
			Action: func(c *cli.Context) error {
				CommandStart(RPCPort, IP, Infohash, Mac, InterfaceName, Keyfile, Key, Until, UseForwarders, UDPPort)
				return nil
			},
		},
		{
			Name:  "stop",
			Usage: "Shutdown p2p instance",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "hash",
					Usage:       "Infohash of instance that needs to be shutdown",
					Value:       "",
					Destination: &Infohash,
				},
				cli.StringFlag{
					Name:        "dev",
					Usage:       "Specify interface name that needs to be removed from interface history",
					Value:       "",
					Destination: &InterfaceName,
				},
			},
			Action: func(c *cli.Context) error {
				CommandStop(RPCPort, Infohash, InterfaceName)
				return nil
			},
		},
		{
			Name:  "show",
			Usage: "Display different information about p2p daemon or instances",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "hash",
					Usage:       "Display information about specific instance",
					Value:       "",
					Destination: &Infohash,
				},
				cli.StringFlag{
					Name:        "check, ip",
					Usage:       "Check if integration with specified IP has been completed",
					Value:       "",
					Destination: &IP,
				},
				cli.BoolFlag{
					Name:        "interfaces",
					Usage:       "List interfaces used by p2p",
					Destination: &ShowInterfaces,
				},
				cli.BoolFlag{
					Name:        "bind",
					Usage:       "Show swarm names along with interfaces",
					Destination: &ShowBind,
				},
				cli.BoolFlag{
					Name:        "all",
					Usage:       "In combination with -interfaces this will show all interfaces used by p2p, even those that is already not in use",
					Destination: &ShowAll,
				},
				cli.BoolFlag{
					Name:        "mtu",
					Usage:       "Display current MTU value in P2P",
					Destination: &ShowMTU,
				},
			},
			Action: func(c *cli.Context) error {
				CommandShow(RPCPort, Infohash, IP, ShowInterfaces, ShowAll, ShowBind, ShowMTU)
				return nil
			},
		},
		{
			Name:  "set",
			Usage: "Modify daemon or instance",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "log",
					Usage:       "Log level. Available levels: trace, debug, info, warning, error",
					Value:       "",
					Destination: &LogLevel,
				},
				cli.StringFlag{
					Name:        "key",
					Usage:       "Append specified key to a list of crypto keys. Must be used with combination of -until",
					Value:       "",
					Destination: &Key,
				},
				cli.StringFlag{
					Name:        "ttl, until",
					Usage:       "Specify until what time this key should work",
					Value:       "",
					Destination: &Until,
				},
				cli.StringFlag{
					Name:        "hash",
					Usage:       "Specify infohash of instance, that should be modified",
					Value:       "",
					Destination: &Infohash,
				},
			},
			Action: func(c *cli.Context) error {
				CommandSet(RPCPort, LogLevel, Infohash, "", Key, Until)
				return nil
			},
		},
		{
			Name:  "debug",
			Usage: "Display debug information",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
			},
			Action: func(c *cli.Context) error {
				CommandDebug(RPCPort)
				return nil
			},
		},
		{
			Name:  "status",
			Usage: "Display connectivity status",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:        "rpc-port",
					Usage:       "RPC port",
					Value:       52523,
					Destination: &RPCPort,
				},
				cli.StringFlag{
					Name:        "hash",
					Usage:       "Limit results to specified instance",
					Value:       "",
					Destination: &Infohash,
				},
			},
			Action: func(c *cli.Context) error {
				CommandStatus(RPCPort, Infohash)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
