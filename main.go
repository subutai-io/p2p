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
		argDirect   string
		argHash     string
		argDht      string
		argKeyfile  string
		argKey      string
		argTTL      string
		argLog      string
		argSaveFile string
		argFwd      bool
		argVersion  bool

		// Daemon configuration and commands
		argDaemon     bool
		argRPCPort    string
		CommandList   bool
		CommandShow   string
		CommandRun    bool
		CommandStop   bool
		CommandSet    bool
		CommandAddKey bool
		CommandDebug  bool
		argProfyle    string
	)

	flag.BoolVar(&argDaemon, "daemon", false, "Starts PTP package in daemon mode")
	flag.StringVar(&argIp, "ip", "none", "IP Address to be used")
	// TODO: Parse this properly
	flag.StringVar(&argMask, "mask", "255.255.255.0", "Network mask")
	flag.StringVar(&argMac, "mac", "", "MAC Address for a TUN/TAP interface")
	flag.StringVar(&argDev, "dev", "", "TUN/TAP interface name")
	// TODO: Direct connection is not implemented yet
	flag.StringVar(&argDirect, "direct", "none", "IP to connect to directly")
	flag.StringVar(&argHash, "hash", "none", "Infohash for environment")
	flag.StringVar(&argDht, "dht", "", "Specify DHT bootstrap node address")
	flag.StringVar(&argKeyfile, "keyfile", "", "Path to yaml file containing crypto key")
	flag.StringVar(&argKey, "key", "", "AES crypto key")
	flag.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")
	flag.StringVar(&argLog, "log", "", "Log level")
	flag.StringVar(&argSaveFile, "save", "", "Path to restore file")
	flag.BoolVar(&argFwd, "fwd", false, "Force traffic forwarding through proxy")
	flag.BoolVar(&argVersion, "version", false, "Show current version")

	// RPC
	flag.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC Communication")
	flag.BoolVar(&CommandList, "list", false, "Lists environments running on this machine")
	flag.BoolVar(&CommandRun, "start", false, "Starts new P2P instance")
	flag.BoolVar(&CommandStop, "stop", false, "Stops P2P instance")
	flag.BoolVar(&CommandSet, "set", false, "Modify p2p behaviour by changing it's options")
	flag.BoolVar(&CommandAddKey, "add-key", false, "Add new key to the list of keys for a specified hash")
	flag.BoolVar(&CommandDebug, "debug", false, "Shows debug info")
	flag.StringVar(&CommandShow, "show", "", "Show known participants of a swarm")

	//profyle
	flag.StringVar(&argProfyle, "profyle", "", "Starts PTP package with profiling. Possible values : memory, cpu")

	flag.Parse()

	if argVersion {
		fmt.Printf("p2p Cloud project %s\n", VERSION)
		os.Exit(0)
	}

	if argDaemon {
		start_profyle(argProfyle)
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
		listen, e := net.Listen("tcp", "localhost:"+argRPCPort)
		if e != nil {
			ptp.Log(ptp.ERROR, "Cannot start RPC listener %v", err)
			os.Exit(1)
		}

		if argSaveFile != "" {
			SaveFile = argSaveFile
			ptp.Log(ptp.INFO, "Restore file provided")
			// Try to restore from provided file
			instances, err := LoadInstances(argSaveFile)
			if err != nil {
				ptp.Log(ptp.ERROR, "%v", err)
			} else {
				ptp.Log(ptp.INFO, "%d instances were loaded from file", len(instances))
				for _, inst := range instances {
					resp := new(Response)
					proc.Run(&inst, resp)
				}
			}
		}

		ptp.Log(ptp.INFO, "Starting RPC Listener")
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
	//if not daemon

	client, err := rpc.DialHTTP("tcp", "localhost:"+argRPCPort)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to connect to RPC %v", err)
		os.Exit(1)
	}
	var response Response
	if CommandList {
		args := &Args{"List", ""}
		err = client.Call("Procedures.List", args, &response)
	} else if CommandSet {
		if argLog != "" {
			args := &NameValueArg{"log", argLog}
			err = client.Call("Procedures.Set", args, &response)
		}
	} else if CommandRun {
		var ok bool = true

		args := &RunArgs{}
		// TODO: Parse ARGS here
		args.Hash = argHash
		args.IP = argIp
		if net.ParseIP(args.IP) == nil {
			fmt.Printf("Bad IP Address specified")
			ok = false
		}
		args.Mask = argMask
		args.Mac = argMac
		args.Dev = argDev
		args.Hash = argHash
		args.Dht = argDht
		args.Keyfile = argKeyfile
		args.Key = argKey
		args.TTL = argTTL
		args.Fwd = argFwd
		if ok {
			err = client.Call("Procedures.Run", args, &response)
		}
	} else if CommandStop {
		args := &StopArgs{}
		args.Hash = argHash
		err = client.Call("Procedures.Stop", args, &response)
	} else if CommandAddKey {
		args := &RunArgs{}
		args.Key = argKey
		args.TTL = argTTL
		args.Hash = argHash
		err = client.Call("Procedures.AddKey", args, &response)
	} else if CommandShow != "" {
		args := &Args{}
		args.Command = CommandShow
		args.Args = "0"
		err = client.Call("Procedures.Show", args, &response)
	} else if CommandDebug {
		args := &Args{}
		err = client.Call("Procedures.Debug", args, &response)
	} else {
		args := &Args{"RandomCommand", "someeeeee"}
		err = client.Call("Procedures.Execute", args, &response)
		if err != nil {
			ptp.Log(ptp.ERROR, "Failed to execute remote procedure %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s\n", response.Output)
	ptp.Log(ptp.DEBUG, "%v", response)
}
