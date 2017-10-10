package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"runtime/pprof"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
)

// AppVersion is a Version of P2P
var AppVersion = "Unknown"

// InterfaceNames - List of all interfaces names that was used by p2p historically. These interfaces may not present in the system anymore
var InterfaceNames []string

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

	var (
		argIP                string
		argMac               string
		argDev               string
		argHash              string
		argDht               string
		argKeyfile           string
		argKey               string
		argTTL               string
		argLog               string
		argsaveFile          string
		argFwd               bool
		argRPCPort           string
		argProfile           string
		argSyslog            string
		argPort              int
		argType              bool
		argShowInterfaces    bool
		argShowInterfacesAll bool
	)

	var Usage = func() {
		fmt.Printf("Usage: p2p <command> [OPTIONS]:\n")
		fmt.Printf("Commands available:\n")
		fmt.Printf("  daemon    Run p2p in daemon mode\n")
		fmt.Printf("  start     Start new p2p instance\n")
		fmt.Printf("  stop      Stop particular p2p instance\n")
		fmt.Printf("  set       Modify p2p options during runtime\n")
		fmt.Printf("  show      Display various information about p2p instances\n")
		fmt.Printf("  status    Show detailed status about connectivity with each peer\n")
		fmt.Printf("  debug     Control debugging and profiling options\n")
		fmt.Printf("  version   Display version information\n")
		fmt.Printf("  help      Show this message or detailed information about commands listed above\n")
		fmt.Printf("\n")
		fmt.Printf("Use 'p2p help <command>' to see detailed help information for specified command\n")
	}

	daemon := flag.NewFlagSet("p2p in daemon mode", flag.ContinueOnError)
	daemon.StringVar(&argsaveFile, "save", "", "Path to restore file")
	daemon.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC communication")
	daemon.StringVar(&argSyslog, "syslog", "", "Syslog socket: 127.0.0.1:1514")
	daemon.StringVar(&argProfile, "profile", "", "Starts PTP package with profiling. Possible values : memory, cpu")

	start := flag.NewFlagSet("Startup options", flag.ContinueOnError)
	start.StringVar(&argIP, "ip", "dhcp", "`IP` address to be used in local system. Should be specified in CIDR format or `dhcp` is used by default to receive free unused IP")
	start.StringVar(&argMac, "mac", "", "MAC or `Hardware Address` for a TUN/TAP interface")
	start.StringVar(&argDev, "dev", "", "TUN/TAP `interface name`")
	start.StringVar(&argHash, "hash", "", "`Infohash` for environment")
	start.StringVar(&argDht, "dht", "", "Specify DHT bootstrap node address in a form of `HOST:PORT`")
	start.StringVar(&argKeyfile, "keyfile", "", "Path to yaml file containing crypto key")
	start.StringVar(&argKey, "key", "", "AES crypto key")
	start.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")
	start.StringVar(&argTTL, "ports", "", "Ports range")
	start.IntVar(&argPort, "port", 0, "`Port` that will be used for p2p communication. Random port number will be generated if no port were specified")
	start.BoolVar(&argFwd, "fwd", false, "If specified, only external routing schemes will be used with use of proxy servers")

	stop := flag.NewFlagSet("Shutdown options", flag.ContinueOnError)
	stop.StringVar(&argHash, "hash", "", "Infohash for environment")

	show := flag.NewFlagSet("Show flagset", flag.ContinueOnError)
	show.StringVar(&argHash, "hash", "", "Infohash for environment")
	show.StringVar(&argIP, "check", "", "Check if integration with specified IP is finished")
	show.BoolVar(&argShowInterfaces, "interfaces", false, "Show interface names")
	show.BoolVar(&argShowInterfacesAll, "all", false, "Show all interfaces")

	set := flag.NewFlagSet("Option Setting", flag.ContinueOnError)
	set.StringVar(&argLog, "log", "", "Log level")
	set.StringVar(&argKey, "key", "", "AES crypto key")
	set.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")
	set.StringVar(&argHash, "hash", "", "Infohash of environment")

	debug := flag.NewFlagSet("Debug and Profiling mode", flag.ContinueOnError)

	version := flag.NewFlagSet("Version output", flag.ContinueOnError)
	version.BoolVar(&argType, "n", false, "Prints numeric variant of the version")

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "help")
	}

	switch os.Args[1] {
	case "daemon":
		daemon.Parse(os.Args[2:])
		if argSyslog != "" {
			ptp.SetSyslogSocket(argSyslog)
		}
		Daemon(argRPCPort, argsaveFile, argProfile)
	case "start":
		start.Parse(os.Args[2:])
		Start(argRPCPort, argIP, argHash, argMac, argDev, argDht, argKeyfile, argKey, argTTL, argFwd, argPort)
	case "stop":
		stop.Parse(os.Args[2:])
		Stop(argRPCPort, argHash)
	case "show":
		show.Parse(os.Args[2:])
		Show(argRPCPort, argHash, argIP, argShowInterfaces, argShowInterfacesAll)
	case "set":
		set.Parse(os.Args[2:])
		Set(argRPCPort, argLog, argHash, argKeyfile, argKey, argTTL)
	case "debug":
		debug.Parse(os.Args[2:])
		Debug(argRPCPort)
	case "version":
		version.Parse(os.Args[2:])
		if argType {
			var macro, minor, micro int
			fmt.Sscanf(AppVersion, "%d.%d.%d", &macro, &minor, &micro)
			fmt.Printf("%d.%d.%d\n", macro, minor, micro)
		} else {
			fmt.Printf("p2p Cloud project %s. Packet version: %s\n", AppVersion, ptp.PacketVersion)
		}
		os.Exit(0)
	case "stop-packet":
		net.DialTimeout("tcp", os.Args[2], 2*time.Second)
		os.Exit(0)
	case "status":
		ShowStatus(argRPCPort)
	case "help":
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "daemon":
				usageDaemon()
				daemon.PrintDefaults()
			case "start":
				usageStart()
				start.PrintDefaults()
			case "show":
				usageShow()
				show.PrintDefaults()
			case "stop":
				usageStop()
				stop.PrintDefaults()
			case "set":
				usageSet()
				set.PrintDefaults()
			}

		} else {
			Usage()
		}
		os.Exit(0)
	default:
		Usage()
		os.Exit(0)
	}
}

// Dial connects to a local RPC server
func Dial(port string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "localhost:"+port)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to connect to RPC %v", err)
		os.Exit(1)
	}
	return client
}

// Start - begin P2P Instance
func Start(rpcPort, ip, hash, mac, dev, dht, keyfile, key, ttl string, fwd bool, port int) {
	client := Dial(rpcPort)
	var response Response

	args := &RunArgs{}
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
		_, err := net.ResolveUDPAddr("udp", dht)
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
	err := client.Call("Procedures.Run", args, &response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	if response.ExitCode == 0 {
		fmt.Printf("%s\n", response.Output)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	}
	os.Exit(response.ExitCode)
}

// Stop will terminate P2P instance
func Stop(rpcPort, hash string) {
	client := Dial(rpcPort)
	var response Response
	args := &StopArgs{}
	if hash == "" {
		fmt.Printf("Specify a hash of instance with -hash argument\n")
		return
	}
	args.Hash = hash
	err := client.Call("Procedures.Stop", args, &response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	if response.ExitCode == 0 {
		fmt.Printf("%s\n", response.Output)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	}
	os.Exit(response.ExitCode)
}

// Show outputs information about P2P instances and interfaces
func Show(rpcPort, hash, ip string, interfaces, all bool) {
	client := Dial(rpcPort)
	var response Response
	args := &ShowArgs{}
	if hash != "" {
		args.Hash = hash
	} else {
		args.Hash = ""
	}
	args.IP = ip
	args.Interfaces = interfaces
	args.All = all
	err := client.Call("Procedures.Show", args, &response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	if response.ExitCode == 0 {
		fmt.Printf("%s\n", response.Output)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	}
	os.Exit(response.ExitCode)
}

// ShowStatus outputs connectivity status of each peer
func ShowStatus(rpcPort string) {
	client := Dial(rpcPort)
	var response Response
	args := &RunArgs{}
	err := client.Call("Procedures.Status", args, &response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	if response.ExitCode == 0 {
		fmt.Printf("%s\n", response.Output)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	}
	os.Exit(response.ExitCode)
}

// Set modifies different options of P2P daemon
func Set(rpcPort, log, hash, keyfile, key, ttl string) {
	client := Dial(rpcPort)
	var response Response
	var err error
	if log != "" {
		args := &NameValueArg{"log", log}
		err = client.Call("Procedures.SetLog", args, &response)
	} else if key != "" {
		args := &RunArgs{}
		args.Key = key
		args.TTL = ttl
		args.Hash = hash
		err = client.Call("Procedures.AddKey", args, &response)
	}
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	if response.ExitCode == 0 {
		fmt.Printf("%s\n", response.Output)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", response.Output)
	}
	os.Exit(response.ExitCode)
}

// Debug prints debug information
func Debug(rpcPort string) {
	client := Dial(rpcPort)
	var response Response
	args := &Args{}
	err := client.Call("Procedures.Debug", args, &response)
	if err != nil {
		fmt.Printf("[ERROR] Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
	os.Exit(response.ExitCode)
}

// Daemon starts P2P daemon
func Daemon(port, sFile, profiling string) {
	StartProfiling(profiling)
	ptp.InitPlatform()
	instances = make(map[string]instance)
	ptp.InitErrors()

	if !ptp.CheckPermissions() {
		os.Exit(1)
	}

	proc := new(Procedures)
	rpc.Register(proc)
	rpc.HandleHTTP()
	listen, err := net.Listen("tcp", "localhost:"+port)
	if err != nil {
		ptp.Log(ptp.Error, "Cannot start RPC listener %v", err)
		os.Exit(1)
	}

	if sFile != "" {
		saveFile = sFile
		ptp.Log(ptp.Info, "Restore file provided")
		// Try to restore from provided file
		instances, err := loadInstances(saveFile)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to load instances: %v", err)
		} else {
			ptp.Log(ptp.Info, "%d instances were loaded from file", len(instances))
			for _, inst := range instances {
				resp := new(Response)
				proc.Run(&inst, resp)
			}
		}
	}

	ptp.Log(ptp.Info, "Starting RPC Listener on %s port", port)
	go http.Serve(listen, nil)

	// Capture SIGINT
	// This is used for development purposes only, but later we should consider updating
	// this code to handle signals
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for sig := range c {
			fmt.Println("Received signal: ", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	select {}
}
