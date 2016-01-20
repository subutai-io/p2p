package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"os/user"
	log "p2p/p2p_log"
	"p2p/udpcs"
	"runtime/pprof"
	"time"
)

type Instance struct {
	PTP *PTPCloud
	ID  string
}

var (
	Instances map[string]Instance
)

type Args struct {
	Command string
	Args    string
}

type NameValueArg struct {
	Name  string
	Value string
}

type RunArgs struct {
	IP      string
	Mask    string
	Mac     string
	Dev     string
	Hash    string
	Dht     string
	Keyfile string
	Key     string
	TTL     string
	Fwd     bool
}

type StopArgs struct {
	Hash string
}

type Response struct {
	ExitCode int
	Output   string
}

type Procedures int

func (p *Procedures) Set(args *NameValueArg, resp *Response) error {
	log.Log(log.INFO, "Setting option %s to %s", args.Name, args.Value)
	resp.ExitCode = 0
	if args.Name == "log" {
		resp.Output = "Logging level has switched to " + args.Value + " level"
		if args.Value == "DEBUG" {
			log.SetMinLogLevel(log.DEBUG)
		} else if args.Value == "INFO" {
			log.SetMinLogLevel(log.INFO)
		} else if args.Value == "TRACE" {
			log.SetMinLogLevel(log.TRACE)
		} else if args.Value == "WARNING" {
			log.SetMinLogLevel(log.WARNING)
		} else if args.Value == "ERROR" {
			log.SetMinLogLevel(log.ERROR)
		} else {
			resp.ExitCode = 1
			resp.Output = "Unknown log level was specified. Supported log levels is:\n"
			resp.Output = resp.Output + "TRACE\n"
			resp.Output = resp.Output + "DEBUG\n"
			resp.Output = resp.Output + "INFO\n"
			resp.Output = resp.Output + "WARNING\n"
			resp.Output = resp.Output + "ERROR\n"
		}
	}
	return nil
}

func (p *Procedures) AddKey(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	if args.Hash == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified hash"
	}
	if args.Key == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified key"
	}
	_, exists := Instances[args.Hash]
	if !exists {
		resp.ExitCode = 1
		resp.Output = "No instances with specified hash were found"
	}
	if resp.ExitCode == 0 {
		resp.Output = "New key added"
		var newKey udpcs.CryptoKey
		newKey = Instances[args.Hash].PTP.Crypter.EncrichKeyValues(newKey, args.Key, args.TTL)
		Instances[args.Hash].PTP.Crypter.Keys = append(Instances[args.Hash].PTP.Crypter.Keys, newKey)
	}
	return nil
}

func (p *Procedures) Execute(args *Args, resp *Response) error {
	log.Log(log.INFO, "Received %v", args)
	resp.ExitCode = 0
	resp.Output = "Command executed"
	return nil
}

func (p *Procedures) Run(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash + "\n"
	var exists bool
	_, exists = Instances[args.Hash]
	if !exists {
		resp.Output = resp.Output + "Lookup finished\n"
		key := []byte(args.Key)
		if len(key) > udpcs.BLOCK_SIZE {
			key = key[:udpcs.BLOCK_SIZE]
		} else {
			zeros := make([]byte, udpcs.BLOCK_SIZE-len(key))
			key = append([]byte(key), zeros...)
		}
		args.Key = string(key)

		ptp := p2pmain(args.IP, args.Mask, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "", args.Fwd)
		var newInst Instance
		newInst.ID = args.Hash
		newInst.PTP = ptp
		Instances[args.Hash] = newInst
		go ptp.Run()
	} else {
		resp.Output = resp.Output + "Hash already in use\n"
	}
	return nil
}

func (p *Procedures) Stop(args *StopArgs, resp *Response) error {
	resp.ExitCode = 0
	var exists bool
	_, exists = Instances[args.Hash]
	if !exists {
		resp.ExitCode = 1
		resp.Output = "Instance with hash " + args.Hash + " was not found"
	} else {
		Instances[args.Hash].PTP.Shutdown = true
		resp.Output = "Shutting down " + args.Hash
		delete(Instances, args.Hash)
	}
	return nil
}

func (p *Procedures) List(args *Args, resp *Response) error {
	resp.ExitCode = 0
	if len(Instances) == 0 {
		resp.Output = "No instances was found"
	}
	for key, inst := range Instances {
		resp.Output = resp.Output + "\t" + inst.PTP.Mac + "\t" + inst.PTP.IP + "\t" + key
		resp.Output = resp.Output + "\n"
	}
	return nil
}

func start_profyle(profyle string) {

	pwd, err := os.Getwd()
	if err != nil {
		log.Log(log.ERROR, "Getwd() error : %v", err)
		return
	}

	time_str := "cpu"
	if profyle == "cpu" {
		file_name := fmt.Sprintf("%s/%s.prof", pwd, time_str)
		f, err := os.Create(file_name)
		if err != nil {
			log.Log(log.ERROR, "Create cpu_prof file failed. %v", err)
			return
		}
		log.Log(log.INFO, "Start cpu profiling to file %s", file_name)
		pprof.StartCPUProfile(f)
	} else if profyle == "memory" {
		_, err := os.Create(fmt.Sprintf("%s/%s.p2p_mem_prof", pwd, time_str))
		if err != nil {
			log.Log(log.ERROR, "Create mem_prof file failed. %v", err)
			return
		}
	}
}

func main() {

	var (
		argIp      string
		argMask    string
		argMac     string
		argDev     string
		argDirect  string
		argHash    string
		argDht     string
		argKeyfile string
		argKey     string
		argTTL     string
		argLog     string
		argFwd     bool

		// Daemon configuration and commands
		argDaemon     bool
		argRPCPort    string
		CommandList   bool
		CommandRun    bool
		CommandStop   bool
		CommandSet    bool
		CommandAddKey bool
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
	flag.BoolVar(&argFwd, "fwd", false, "Force traffic forwarding through proxy")

	// RPC
	flag.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC Communication")
	flag.BoolVar(&CommandList, "list", false, "Lists environments running on this machine")
	flag.BoolVar(&CommandRun, "start", false, "Starts new P2P instance")
	flag.BoolVar(&CommandStop, "stop", false, "Stops P2P instance")
	flag.BoolVar(&CommandSet, "set", false, "Modify p2p behaviour by changing it's options")
	flag.BoolVar(&CommandAddKey, "add-key", false, "Add new key to the list of keys for a specified hash")

	//profyle
	flag.StringVar(&argProfyle, "profyle", "", "Starts PTP package with profiling. Possible values : memory, cpu")

	flag.Parse()

	if argDaemon {
		start_profyle(argProfyle)
		Instances = make(map[string]Instance)
		user, err := user.Current()
		if err != nil {
			log.Log(log.ERROR, "Failed to retrieve information about user: %v", err)
		}
		if user.Uid != "0" {
			log.Log(log.ERROR, "P2P cannot run in daemon mode without root privileges")
			os.Exit(1)
		}

		proc := new(Procedures)
		rpc.Register(proc)
		rpc.HandleHTTP()
		listen, e := net.Listen("tcp", "localhost:"+argRPCPort)
		if e != nil {
			log.Log(log.ERROR, "Cannot start RPC listener %v", err)
			os.Exit(1)
		}
		log.Log(log.INFO, "Starting RPC Listener")
		go http.Serve(listen, nil)
		//p2pmain(argIp, argMask, argMac, argDev, argDirect, argHash, argDht, argKeyfile, argKey, argTTL, argLog)
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
		log.Log(log.ERROR, "Failed to connect to RPC %v", err)
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

		args := &RunArgs{}
		// TODO: Parse ARGS here
		args.Hash = argHash
		args.IP = argIp
		args.Mask = argMask
		args.Mac = argMac
		args.Dev = argDev
		args.Hash = argHash
		args.Dht = argDht
		args.Keyfile = argKeyfile
		args.Key = argKey
		args.TTL = argTTL
		args.Fwd = argFwd
		err = client.Call("Procedures.Run", args, &response)
	} else if CommandStop {
		args := &StopArgs{}
		args.Hash = argHash
		err = client.Call("Procedures.Stop", args, &response)
	} else if CommandAddKey {
		args := &RunArgs{}
		args.Key = argKey
		args.TTL = argTTL
		args.Hash = argHash
		err = client.Call("Procedure.AddKey", args, &response)
	} else {
		args := &Args{"RandomCommand", "someeeeee"}
		err = client.Call("Procedures.Execute", args, &response)
		if err != nil {
			log.Log(log.ERROR, "Failed to execute remote procedure %v", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s\n", response.Output)
	log.Log(log.DEBUG, "%v", response)
}
