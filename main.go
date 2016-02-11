package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	ptp "github.com/subutai-io/p2p/lib"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"os/signal"
	"os/user"
	"runtime"
	"runtime/pprof"
	"time"
)

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

type Instance struct {
	PTP  *PTPCloud
	ID   string
	Args RunArgs
}

var (
	Instances map[string]Instance
	SaveFile  string
)

func EncodeInstances() ([]byte, error) {
	var savedInstances []RunArgs

	for _, inst := range Instances {
		savedInstances = append(savedInstances, inst.Args)
	}
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(savedInstances)
	if err != nil {
		ptp.Log(ptp.ERROR, "Failed to encode instances: %v", err)
		return []byte(""), err
	}
	return b.Bytes(), nil
}

func DecodeInstances(data []byte) ([]RunArgs, error) {
	var args []RunArgs
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(&args)
	if err != nil {
		return args, err
	}
	return args, nil
}

// Calls EncodeInstances() and saves results into specified file
// Return number of bytes written and error if any
func SaveInstances(filename string) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return 0, err
	}

	data, err := EncodeInstances()
	if err != nil {
		return 0, err
	}

	s, err := file.Write(data)
	if err != nil {
		return s, err
	}
	file.Close()
	return s, nil
}

func LoadInstances(filename string) ([]RunArgs, error) {
	var loadedInstances []RunArgs
	file, err := os.Open(filename)
	data := make([]byte, 100000)
	_, err = file.Read(data)
	if err != nil {
		return loadedInstances, err
	}

	loadedInstances, err = DecodeInstances(data)
	if err != nil {
		return loadedInstances, err
	}

	return loadedInstances, nil
}

type Args struct {
	Command string
	Args    string
}

type NameValueArg struct {
	Name  string
	Value string
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
	ptp.Log(ptp.INFO, "Setting option %s to %s", args.Name, args.Value)
	resp.ExitCode = 0
	if args.Name == "log" {
		resp.Output = "Logging level has switched to " + args.Value + " level"
		if args.Value == "DEBUG" {
			ptp.SetMinLogLevel(ptp.DEBUG)
		} else if args.Value == "INFO" {
			ptp.SetMinLogLevel(ptp.INFO)
		} else if args.Value == "TRACE" {
			ptp.SetMinLogLevel(ptp.TRACE)
		} else if args.Value == "WARNING" {
			ptp.SetMinLogLevel(ptp.WARNING)
		} else if args.Value == "ERROR" {
			ptp.SetMinLogLevel(ptp.ERROR)
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
		var newKey ptp.CryptoKey
		newKey = Instances[args.Hash].PTP.Crypter.EnrichKeyValues(newKey, args.Key, args.TTL)
		Instances[args.Hash].PTP.Crypter.Keys = append(Instances[args.Hash].PTP.Crypter.Keys, newKey)
	}
	return nil
}

func (p *Procedures) Execute(args *Args, resp *Response) error {
	ptp.Log(ptp.INFO, "Received %v", args)
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
		if args.Key != "" {
			key := []byte(args.Key)
			if len(key) > ptp.BLOCK_SIZE {
				key = key[:ptp.BLOCK_SIZE]
			} else {
				zeros := make([]byte, ptp.BLOCK_SIZE-len(key))
				key = append([]byte(key), zeros...)
			}
			args.Key = string(key)
		}

		ptpInstance := p2pmain(args.IP, args.Mask, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "", args.Fwd)
		var newInst Instance
		newInst.ID = args.Hash
		newInst.PTP = ptpInstance
		newInst.Args = *args
		Instances[args.Hash] = newInst
		go ptpInstance.Run()
		if SaveFile != "" {
			SaveInstances(SaveFile)
		}
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
		resp.Output = "Shutting down " + args.Hash
		Instances[args.Hash].PTP.StopInstance()
		delete(Instances, args.Hash)
	}
	return nil
}

func (p *Procedures) Show(args *Args, resp *Response) error {
	swarm, exists := Instances[args.Command]
	resp.ExitCode = 0
	if exists {
		resp.Output = "< Peer ID >\t< IP >\t< Endpoint >\t< HW >\n"
		for _, peer := range swarm.PTP.NetworkPeers {
			resp.Output = resp.Output + peer.ID + "\t"
			resp.Output = resp.Output + peer.PeerLocalIP.String() + "\t"
			resp.Output = resp.Output + peer.Endpoint + "\t"
			resp.Output = resp.Output + peer.PeerHW.String() + "\n"
		}
	} else {
		resp.Output = "Specified environment was not found"
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

func (p *Procedures) Debug(args *Args, resp *Response) error {
	resp.Output = "DEBUG INFO:\n"
	resp.Output += fmt.Sprintf("Number of gouroutines: %d\n", runtime.NumGoroutine())
	resp.Output += fmt.Sprintf("Instances information:\n")
	for _, ins := range Instances {
		resp.Output += fmt.Sprintf("Hash: %s\n", ins.ID)
		resp.Output += fmt.Sprintf("ID: %s\n", ins.PTP.dht.ID)
		resp.Output += fmt.Sprintf("Interface %s, HW Addr: %s, IP: %s\n", ins.PTP.DeviceName, ins.PTP.Mac, ins.PTP.IP)
		resp.Output += fmt.Sprintf("Peers:\n")
		// TODO: Rewrite this part
		for _, id := range ins.PTP.IPIDTable {
			resp.Output += fmt.Sprintf("\t--- %s ---\n", id)
			peer, exists := ins.PTP.NetworkPeers[id]
			if !exists {
				resp.Output += fmt.Sprintf("\tPeer was not integrated into network\n")
			} else {
				resp.Output += fmt.Sprintf("\t\tHWAddr: %s\n", peer.PeerHW.String())
				resp.Output += fmt.Sprintf("\t\tIP: %s\n", peer.PeerLocalIP.String())
				resp.Output += fmt.Sprintf("\t\tEndpoint: %s\n", peer.Endpoint)
				resp.Output += fmt.Sprintf("\t\tPeer Address: %s\n", peer.PeerAddr.String())
				resp.Output += fmt.Sprintf("\t\tProxy ID: %d\n", peer.ProxyID)
			}
			resp.Output += fmt.Sprintf("\t--- End of %s ---\n", id)
		}
	}
	return nil
}

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
