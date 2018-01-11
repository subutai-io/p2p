package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/ccding/go-stun/stun"
	ptp "github.com/subutai-io/p2p/lib"
	"github.com/urfave/cli"
)

// AppVersion is a Version of P2P
var AppVersion = "Unknown"
var BuildID = "Unknown"
var DefaultDHT = "mdht.subut.ai:6881"

// InterfaceNames - List of all interfaces names that was used by p2p historically. These interfaces may not present in the system anymore
var InterfaceNames []string

// OutboundIP is an outbound IP address detected by STUN
var OutboundIP net.IP

var SignalChannel chan os.Signal

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
		DHTRouters     string // Comma-separated list of DHT routers
		Keyfile        string // Path to a file with crypto key
		Key            string // AES key
		Until          string // Until date this key will be active in Unix timestamp
		Ports          string // Ports range for an instance
		UDPPort        int    // Specific UDP port for an instance
		UseForwarders  bool   // Whether or not p2p should force usage of proxy servers for this instance
		ShowInterfaces bool   // Whether or not p2p show command should return information about interfaces in use
		ShowAll        bool   //
		LogLevel       string // Log level
		RemoveService  bool   // If yes - service will be removed (used with service)
		InstallService bool   // If yes - service will be installed (used with service)
	)

	app := cli.NewApp()
	app.Name = "p2p"
	app.Version = AppVersion
	app.Authors = []cli.Author{
		cli.Author{
			Name: "Subutai.io",
		},
	}
	app.Copyright = "Copyright 2017 Subutai.io"

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
			},
			Action: func(c *cli.Context) error {
				ExecDaemon(RPCPort, SaveFile, Profiling, Syslog)
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
				cli.StringFlag{
					Name:        "dht",
					Usage:       "Comman-separated list of DHT bootstrap nodes",
					Value:       "",
					Destination: &DHTRouters,
				},
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
				Start(RPCPort, IP, Infohash, Mac, InterfaceName, DHTRouters, Keyfile, Key, Until, UseForwarders, UDPPort)
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
				Stop(RPCPort, Infohash, InterfaceName)
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
					Name:        "all",
					Usage:       "In combination with -interfaces this will show all interfaces used by p2p, even those that is already not in use",
					Destination: &ShowAll,
				},
			},
			Action: func(c *cli.Context) error {
				Show(RPCPort, Infohash, IP, ShowInterfaces, ShowAll)
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
					Usage:       "Log level. Available levels: TRACE, DEBUG, INFO, WARNING, ERROR, FATAL",
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
				Set(RPCPort, LogLevel, Infohash, "", Key, Until)
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
				Debug(RPCPort)
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
			},
			Action: func(c *cli.Context) error {
				ShowStatus(RPCPort)
				return nil
			},
		},
	}
	app.Run(os.Args)
}

// Dial connects to a local RPC server
func Dial(rpchost string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", rpchost)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to connect to RPC %v", err)
		os.Exit(1)
	}
	return client
}

// Start - begin P2P Instance
func Start(rpcPort int, ip, hash, mac, dev, dht, keyfile, key, ttl string, fwd bool, port int) {
	// client := Dial(fmt.Sprintf("localhost:%d", rpcPort))
	// var response Response

	args := &DaemonArgs{}
	/*if net.ParseIP(ip) == nil {
		fmt.Printf("Bad IP Address specified\n")
		return
	}*/
	args.IP = ip
	if hash == "" {
		fmt.Printf("Hash cannot be empty. Please start new instances with -hash VALUE argument\n")
		return
	}
	args.Hash = hash
	if mac != "" {
		_, err := net.ParseMAC(mac)
		if err != nil {
			fmt.Printf("Invalid MAC address provided\n")
			return
		}
	}
	args.Mac = mac
	args.Dev = dev
	if dht != "" {
		_, err := net.ResolveUDPAddr("udp4", dht)
		if err != nil {
			fmt.Printf("Invalid DHT node address provided. Please specify correct DHT address in form HOST:PORT\n")
			return
		}
	}
	args.Dht = dht
	args.Keyfile = keyfile
	args.Key = key
	args.TTL = ttl
	args.Fwd = fwd
	args.Port = port

	out, err := sendRequest(rpcPort, "start", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)

	// err := client.Call("Daemon.Run", args, &response)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
	// 	return
	// }
	// if response.ExitCode == 0 {
	// 	fmt.Printf("%s\n", response.Output)
	// } else {
	// 	fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	// }
	// os.Exit(response.ExitCode)
}

// Stop will terminate P2P instance
func Stop(rpcPort int, hash, dev string) {
	args := &DaemonArgs{}
	if hash != "" {
		args.Hash = hash
		args.Dev = ""
	} else if dev != "" {
		args.Dev = dev
		args.Hash = ""
	} else {
		fmt.Printf("Not enough parameters for stop command")
		return
	}
	out, err := sendRequest(rpcPort, "stop", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
}

// Show outputs information about P2P instances and interfaces
func Show(queryPort int, hash, ip string, interfaces, all bool) {
	// client := Dial(fmt.Sprintf("localhost:%d", rpcPort))
	// var response Response
	args := &DaemonArgs{}
	if hash != "" {
		args.Hash = hash
	} else {
		args.Hash = ""
	}
	args.IP = ip
	args.Interfaces = interfaces
	args.All = all

	out, err := sendRequestRaw(queryPort, "show", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	show := []ShowOutput{}
	err = json.Unmarshal(out, &show)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON. Error %s\n", err)
		os.Exit(99)
	}

	if args.Hash != "" {
		if args.IP != "" {
			for _, m := range show {
				if m.Code != 0 {
					fmt.Println(m.Error)
				} else {
					fmt.Println(m.Text)
				}
				os.Exit(m.Code)
			}
			fmt.Println("No data available")
			os.Exit(102)
		} else {
			fmt.Println("< Peer ID >\t< IP >\t< Endpoint >\t< HW >")
			for _, m := range show {
				fmt.Printf("%s\t%s\t%s\t%s\n", m.ID, m.IP, m.Endpoint, m.HardwareAddress)
			}
			os.Exit(0)
		}
	}
	if args.Interfaces {
		for _, m := range show {
			fmt.Println(m.InterfaceName)
		}
		os.Exit(0)
	}

	for _, m := range show {
		fmt.Printf("%s\t%s\t%s\n", m.HardwareAddress, m.IP, m.Hash)
	}
	os.Exit(0)

	// err := client.Call("Daemon.Show", args, &response)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
	// 	return
	// }
	// if response.ExitCode == 0 {
	// 	fmt.Printf("%s\n", response.Output)
	// } else {
	// 	fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	// }
	// os.Exit(response.ExitCode)
}

// ShowStatus outputs connectivity status of each peer
func ShowStatus(rpcPort int) {
	// client := Dial(fmt.Sprintf("localhost:%d", rpcPort))
	// var response Response
	args := &DaemonArgs{}
	// err := client.Call("Daemon.Status", args, &response)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
	// 	return
	// }

	out, err := sendRequest(rpcPort, "status", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)

	// if response.ExitCode == 0 {
	// 	fmt.Printf("%s\n", response.Output)
	// } else {
	// 	fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	// }
	// os.Exit(response.ExitCode)
}

// Set modifies different options of P2P daemon
func Set(rpcPort int, log, hash, keyfile, key, ttl string) {
	out, err := sendRequest(rpcPort, "set", &DaemonArgs{Log: log, Keyfile: keyfile, Key: key, TTL: ttl})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
	// client := Dial(fmt.Sprintf("localhost:%d", rpcPort))
	// var response Response
	// var err error
	// if log != "" {
	// 	args := &NameValueArg{"log", log}
	// 	err = client.Call("Daemon.SetLog", args, &response)
	// } else if key != "" {
	// 	args := &RunArgs{}
	// 	args.Key = key
	// 	args.TTL = ttl
	// 	args.Hash = hash
	// 	err = client.Call("Daemon.AddKey", args, &response)
	// }
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
	// 	return
	// }
	// if response.ExitCode == 0 {
	// 	fmt.Printf("%s\n", response.Output)
	// } else {
	// 	fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	// }
	// os.Exit(response.ExitCode)
}

// Debug prints debug information
func Debug(rpcPort int) {
	out, err := sendRequest(rpcPort, "debug", &DaemonArgs{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
	// client := Dial(fmt.Sprintf("localhost:%d", rpcPort))
	// var response Response
	// args := &Args{}
	// err := client.Call("Daemon.Debug", args, &response)
	// if err != nil {
	// 	fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
	// 	return
	// }
	// fmt.Printf("%s\n", response.Output)
	// os.Exit(response.ExitCode)
}

// ExecDaemon starts P2P daemon
func ExecDaemon(port int, sFile, profiling, syslog string) {
	if syslog != "" {
		ptp.SetSyslogSocket(syslog)
	}
	StartProfiling(profiling)
	go ptp.InitPlatform()
	ptp.InitErrors()

	if !ptp.CheckPermissions() {
		os.Exit(1)
	}

	ptp.Log(ptp.Info, "Determining outbound IP")
	nat, host, err := stun.NewClient().Discover()
	if err != nil {
		ptp.Log(ptp.Error, "Failed to discover outbound IP: %s", err)
		OutboundIP = nil
	} else {
		OutboundIP = net.ParseIP(host.IP())
		ptp.Log(ptp.Info, "Public IP is %s. %s", OutboundIP.String(), nat)
	}

	proc := new(Daemon)
	proc.Initialize(sFile)
	setupRESTHandlers(port, proc)

	if sFile != "" {
		ptp.Log(ptp.Info, "Restore file provided")
		// Try to restore from provided file
		instances, err := proc.Instances.LoadInstances(proc.SaveFile)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to load instances: %v", err)
		} else {
			ptp.Log(ptp.Info, "%d instances were loaded from file", len(instances))
			for _, inst := range instances {
				proc.Run(&inst, new(Response))
			}
		}
	}

	SignalChannel = make(chan os.Signal, 1)
	signal.Notify(SignalChannel, os.Interrupt)

	go func() {
		for sig := range SignalChannel {
			fmt.Println("Received signal: ", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	select {}
}

func setupRESTHandlers(port int, d *Daemon) {
	http.HandleFunc("/rest/v1/start", d.execRESTStart)
	http.HandleFunc("/rest/v1/stop", d.execRESTStop)
	http.HandleFunc("/rest/v1/show", d.execRESTShow)
	http.HandleFunc("/rest/v1/status", d.execRESTStatus)
	http.HandleFunc("/rest/v1/debug", d.execRESTDebug)
	http.HandleFunc("/rest/v1/set", d.execRESTSet)

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func sendRequest(port int, command string, args *DaemonArgs) (*RESTResponse, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/rest/v1/%s", port, command), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute command. Check if p2p daemon is running.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	out := &RESTResponse{}
	err = json.Unmarshal(body, out)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response: %s", err)
	}
	return out, nil
}

func sendRequestRaw(port int, command string, args *DaemonArgs) ([]byte, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/rest/v1/%s", port, command), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute command. Check if p2p daemon is running.")
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
