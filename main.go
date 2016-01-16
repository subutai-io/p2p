package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/user"
	log "p2p/p2p_log"
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
		if len(key) > 32 {
			key = key[:32]
		} else {
			zeros := make([]byte, 32-len(key))
			key = append([]byte(key), zeros...)
		}
		args.Key = string(key)

		ptp := p2pmain(args.IP, args.Mask, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "")
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

func main() {

	var (
		argDaemon  bool
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

		// RPC
		argRPCPort string
		argList    bool
		argRun     bool
		argStop    bool
		argSet     bool
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

	// RPC
	flag.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC Communication")
	flag.BoolVar(&argList, "list", false, "Lists environments running on this machine")
	flag.BoolVar(&argRun, "start", false, "Starts new P2P instance")
	flag.BoolVar(&argStop, "stop", false, "Stops P2P instance")
	flag.BoolVar(&argSet, "set", false, "Modify p2p behaviour by changing it's options")

	flag.Parse()
	if argDaemon {
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
		for {
			time.Sleep(5 * time.Second)
		}
	} else {
		client, err := rpc.DialHTTP("tcp", "localhost:"+argRPCPort)
		if err != nil {
			log.Log(log.ERROR, "Failed to connect to RPC %v", err)
			os.Exit(1)
		}
		var response Response
		if argList {
			args := &Args{"List", ""}
			err = client.Call("Procedures.List", args, &response)
		} else if argSet {
			if argLog != "" {
				args := &NameValueArg{"log", argLog}
				err = client.Call("Procedures.Set", args, &response)
			}
		} else if argRun {

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
			err = client.Call("Procedures.Run", args, &response)
		} else if argStop {
			args := &StopArgs{}
			args.Hash = argHash
			err = client.Call("Procedures.Stop", args, &response)
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
}
