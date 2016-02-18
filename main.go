package main

import (
	"flag"
	"fmt"
	ptp "github.com/subutai-io/p2p/lib"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"os/user"
	"runtime/pprof"
	"time"
)

var VERSION string = "Unknown"

func start_profyle(profyle string) {

	pwd, err := os.Getwd()
	if err != nil {
		ptp.Log(ptp.ERROR, "Getwd() error : %v", err)
		return
	}

	time_str := "cpu"
	if profyle == "cpu" {
		file_name := fmt.Sprintf("%s/%s.prof", pwd, time_str)
		f, err := os.Create(file_name)
		if err != nil {
			ptp.Log(ptp.ERROR, "Create cpu_prof file failed. %v", err)
			return
		}
		ptp.Log(ptp.INFO, "Start cpu profiling to file %s", file_name)
		pprof.StartCPUProfile(f)
	} else if profyle == "memory" {
		_, err := os.Create(fmt.Sprintf("%s/%s.p2p_mem_prof", pwd, time_str))
		if err != nil {
			ptp.Log(ptp.ERROR, "Create mem_prof file failed. %v", err)
			return
		}
	}
}

func main() {

	var (
		argIp       string
		argMask     string
		argMac      string
		argDev      string
		argHash     string
		argDht      string
		argKeyfile  string
		argKey      string
		argTTL      string
		argLog      string
		argSaveFile string
		argFwd      bool
		argRPCPort  string
		argProfile  string
	)

	var Usage = func() {
		fmt.Printf("Usage: p2p <command> [OPTIONS]:\n")
		fmt.Printf("Commands available:\n")
		fmt.Printf("  daemon    Run p2p in daemon mode\n")
		fmt.Printf("  start     Start new p2p instance\n")
		fmt.Printf("  stop      Stop particular p2p instance\n")
		fmt.Printf("  set       Modify p2p options during runtime\n")
		fmt.Printf("  show      Display various information about p2p instances\n")
		fmt.Printf("  debug     Control debugging and profiling options\n")
		fmt.Printf("  version   Display version information\n")
		fmt.Printf("  help      Show this message or detailed information about commands listed above\n")
		fmt.Printf("\n")
		fmt.Printf("Use 'p2p help <command>' to see detailed help information for specified command\n")
	}

	daemon := flag.NewFlagSet("p2p in daemon mode", flag.ContinueOnError)
	daemon.StringVar(&argSaveFile, "save", "", "Path to restore file")
	daemon.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC communication")
	daemon.StringVar(&argProfile, "profyle", "", "Starts PTP package with profiling. Possible values : memory, cpu")

	start := flag.NewFlagSet("Startup options", flag.ContinueOnError)
	start.StringVar(&argIp, "ip", "none", "`IP` address to be used")
	start.StringVar(&argMask, "mask", "255.255.255.0", "Network mask a.k.a. `subnet`")
	start.StringVar(&argMac, "mac", "", "MAC or `Hardware Address` for a TUN/TAP interface")
	start.StringVar(&argDev, "dev", "", "TUN/TAP `interface name`")
	start.StringVar(&argHash, "hash", "", "`Infohash` for environment")
	start.StringVar(&argDht, "dht", "", "Specify DHT bootstrap node address in a form of `HOST:PORT`")
	start.StringVar(&argKeyfile, "keyfile", "", "Path to yaml file containing crypto key")
	start.StringVar(&argKey, "key", "", "AES crypto key")
	start.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")

	stop := flag.NewFlagSet("Shutdown options", flag.ContinueOnError)
	stop.StringVar(&argHash, "hash", "", "Infohash for environment")

	show := flag.NewFlagSet("Show flagset", flag.ContinueOnError)
	show.StringVar(&argHash, "hash", "", "Infohash for environment")

	set := flag.NewFlagSet("Option Setting", flag.ContinueOnError)
	set.StringVar(&argLog, "log", "", "Log level")
	set.StringVar(&argKey, "key", "", "AES crypto key")
	set.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")

	debug := flag.NewFlagSet("Debug and Profiling mode", flag.ContinueOnError)

	if len(os.Args) < 2 {
		os.Args = append(os.Args, "help")
	}

	switch os.Args[1] {
	case "daemon":
		daemon.Parse(os.Args[2:])
		Daemon(argRPCPort, argSaveFile, argProfile)
	case "start":
		start.Parse(os.Args[2:])
		Start(argRPCPort, argIp, argHash, argMask, argMac, argDev, argDht, argKeyfile, argKey, argTTL, argFwd)
	case "stop":
		stop.Parse(os.Args[2:])
		Stop(argRPCPort, argHash)
	case "show":
		show.Parse(os.Args[2:])
		Show(argRPCPort, argHash)
	case "set":
		set.Parse(os.Args[2:])
		Set(argRPCPort, argLog, argHash, argKeyfile, argKey, argTTL)
	case "debug":
		debug.Parse(os.Args[2:])
		Debug(argRPCPort)
	case "version":
		fmt.Printf("p2p Cloud project %s\n", VERSION)
		os.Exit(0)
	case "help":
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "daemon":
				UsageDaemon()
				daemon.PrintDefaults()
			case "start":
				UsageStart()
				start.PrintDefaults()
			case "show":
				UsageShow()
				show.PrintDefaults()
			case "stop":
				UsageStop()
				stop.PrintDefaults()
			case "set":
				UsageSet()
				set.PrintDefaults()
			}

		} else {
			Usage()
		}
		os.Exit(1)
	default:
		Usage()
		os.Exit(1)
	}
}

func Dial(port string) *rpc.Client {
	client, err := rpc.DialHTTP("tcp", "localhost:"+port)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to connect to RPC %v", err)
		os.Exit(1)
	}
	return client
}

func Start(rpcPort, ip, hash, mask, mac, dev, dht, keyfile, key, ttl string, fwd bool) {
	client := Dial(rpcPort)
	var response Response

	args := &RunArgs{}
	if net.ParseIP(ip) == nil {
		fmt.Printf("Bad IP Address specified\n")
		return
	}
	args.IP = ip
	if hash == "" {
		fmt.Printf("Hash cannot be empty. Please start new instances with -hash VALUE argument\n")
		return
	}
	args.Hash = hash
	args.Mask = mask
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
	err := client.Call("Procedures.Run", args, &response)
	if err != nil {
		fmt.Printf("Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
}

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
		fmt.Printf("Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
}

func Show(rpcPort, hash string) {
	client := Dial(rpcPort)
	var response Response
	args := &Args{}
	args.Command = ""
	if hash != "" {
		args.Args = hash
	} else {
		args.Args = ""
	}
	err := client.Call("Procedures.Show", args, &response)
	if err != nil {
		fmt.Printf("Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
}

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
		fmt.Printf("Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
}

func Debug(rpcPort string) {
	client := Dial(rpcPort)
	var response Response
	args := &Args{}
	err := client.Call("Procedures.Debug", args, &response)
	if err != nil {
		fmt.Printf("Failed to run RPC request: %v\n", err)
		return
	}
	fmt.Printf("%s\n", response.Output)
}

func Daemon(port, saveFile, profiling string) {
	start_profyle(profiling)
	Instances = make(map[string]Instance)
	user, err := user.Current()
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to retrieve information about user: %v", err)
	}
	if user.Uid != "0" {
		ptp.Log(ptp.ERROR, "P2P cannot run in daemon mode without root privileges")
		os.Exit(1)
	}

	proc := new(Procedures)
	rpc.Register(proc)
	rpc.HandleHTTP()
	listen, e := net.Listen("tcp", "localhost:"+port)
	if e != nil {
		ptp.Log(ptp.ERROR, "Cannot start RPC listener %v", err)
		os.Exit(1)
	}

	if saveFile != "" {
		ptp.Log(ptp.INFO, "Restore file provided")
		// Try to restore from provided file
		instances, err := LoadInstances(saveFile)
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to load instances: %v", err)
		} else {
			ptp.Log(ptp.INFO, "%d instances were loaded from file", len(instances))
			for _, inst := range instances {
				resp := new(Response)
				proc.Run(&inst, resp)
			}
		}
	}

	ptp.Log(ptp.INFO, "Starting RPC Listener on %s port", port)
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
	for {
		time.Sleep(5 * time.Second)
	}
	return
}
