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
)

type Args struct {
	Command string
	Args    string
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

func (p *Procedures) Execute(args *Args, resp *Response) error {
	log.Log(log.INFO, "Received %v", args)
	resp.ExitCode = 0
	resp.Output = "Command executed"
	return nil
}

func (p *Procedures) Run(args *RunArgs, resp *Response) error {
	ptp := p2pmain(args.IP, args.Mask, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "")
	go ptp.Run()
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash
	return nil
}

func (p *Procedures) Stop(args *StopArgs, resp *Response) error {
	return nil

}

func (p *Procedures) List(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
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
	)

	flag.BoolVar(&argDaemon, "daemon", false, "Starts PTP package in daemon mode")
	flag.StringVar(&argIp, "ip", "none", "IP Address to be used")
	// TODO: Parse this properly
	flag.StringVar(&argMask, "mask", "255.255.255.0", "Network mask")
	flag.StringVar(&argMac, "mac", "none", "MAC Address for a TUN/TAP interface")
	flag.StringVar(&argDev, "dev", "", "TUN/TAP interface name")
	// TODO: Direct connection is not implemented yet
	flag.StringVar(&argDirect, "direct", "none", "IP to connect to directly")
	flag.StringVar(&argHash, "hash", "none", "Infohash for environment")
	flag.StringVar(&argDht, "dht", "", "Specify DHT bootstrap node address")
	flag.StringVar(&argKeyfile, "keyfile", "", "Path to yaml file containing crypto key")
	flag.StringVar(&argKey, "key", "", "AES crypto key")
	flag.StringVar(&argTTL, "ttl", "", "Time until specified key will be available")
	flag.StringVar(&argLog, "log", "INFO", "Log level")

	// RPC
	flag.StringVar(&argRPCPort, "rpc", "52523", "Port for RPC Communication")
	flag.BoolVar(&argList, "list", false, "Lists environments running on this machine")
	flag.BoolVar(&argRun, "start", false, "Starts new P2P instance")
	flag.BoolVar(&argStop, "stop", false, "Stops P2P instance")

	flag.Parse()
	if argDaemon {
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
		} else if argRun {
			args := &RunArgs{}
			// TODO: Parse ARGS here
			args.Hash = argHash
			args.IP = argIp
			err = client.Call("Procedures.Run", args, &response)
		} else if argStop {

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
