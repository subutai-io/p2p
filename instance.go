package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"runtime"
	"sync"

	ptp "github.com/subutai-io/p2p/lib"
)

// RunArgs is a list of arguments used at instance startup and
// some other RPC calls
type RunArgs struct {
	IP      string
	Mac     string
	Dev     string
	Hash    string
	Dht     string
	Keyfile string
	Key     string
	TTL     string
	Fwd     bool
	Port    int
}

type instance struct {
	PTP  *ptp.PeerToPeer
	ID   string
	Args RunArgs
}

var (
	instances     map[string]instance
	saveFile      string
	instances_mut sync.Mutex
)

func encodeInstances() ([]byte, error) {
	var savedInstances []RunArgs
	instances_mut.Lock()
	for _, inst := range instances {
		savedInstances = append(savedInstances, inst.Args)
	}
	instances_mut.Unlock()
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(savedInstances)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to encode instances: %v", err)
		return []byte(""), err
	}
	return b.Bytes(), nil
}

func decodeInstances(data []byte) ([]RunArgs, error) {
	var args []RunArgs
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(&args)
	return args, err
}

// Calls encodeInstances() and saves results into specified file
// Return number of bytes written and error if any
func saveInstances(filename string) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	data, err := encodeInstances()
	if err != nil {
		return 0, err
	}

	s, err := file.Write(data)
	if err != nil {
		return s, err
	}
	return s, nil
}

func loadInstances(filename string) ([]RunArgs, error) {
	var loadedInstances []RunArgs
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]byte, 100000)
	_, err = file.Read(data)
	if err != nil {
		return loadedInstances, err
	}

	loadedInstances, err = decodeInstances(data)
	return loadedInstances, err
}

// Args is a simple name-value RPC arguments
type Args struct {
	Command string
	Args    string
}

// NameValueArg is a simple name-value RPC arguments
type NameValueArg struct {
	Name  string
	Value string
}

// StopArgs is an arguments used with Stop RPC call
type StopArgs struct {
	Hash string
}

// Response represents a result of RPC call
type Response struct {
	ExitCode int
	Output   string
}

// Procedures is an object of RPC procedures
type Procedures int

// SetLog modifies specific option
func (p *Procedures) SetLog(args *NameValueArg, resp *Response) error {
	ptp.Log(ptp.Info, "Setting option %s to %s", args.Name, args.Value)
	resp.ExitCode = 0
	if args.Name == "log" {
		resp.Output = "Logging level has switched to " + args.Value + " level"
		if args.Value == "DEBUG" {
			ptp.SetMinLogLevel(ptp.Debug)
		} else if args.Value == "INFO" {
			ptp.SetMinLogLevel(ptp.Info)
		} else if args.Value == "TRACE" {
			ptp.SetMinLogLevel(ptp.Trace)
		} else if args.Value == "WARNING" {
			ptp.SetMinLogLevel(ptp.Warning)
		} else if args.Value == "ERROR" {
			ptp.SetMinLogLevel(ptp.Error)
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

// AddKey adds a new crypto-key
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
	instances_mut.Lock()
	_, exists := instances[args.Hash]
	if !exists {
		resp.ExitCode = 1
		resp.Output = "No instances with specified hash were found"
	}
	if resp.ExitCode == 0 {
		resp.Output = "New key added"
		var newKey ptp.CryptoKey
		newKey = instances[args.Hash].PTP.Crypter.EnrichKeyValues(newKey, args.Key, args.TTL)
		instances[args.Hash].PTP.Crypter.Keys = append(instances[args.Hash].PTP.Crypter.Keys, newKey)
	}
	instances_mut.Unlock()
	return nil
}

// Execute is a dummy method used for tests
func (p *Procedures) Execute(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
	return nil
}

// Run starts a P2P instance
func (p *Procedures) Run(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash + "\n"

	// Validate if interface name is unique
	if args.Dev != "" {
		instances_mut.Lock()
		for _, inst := range instances {
			if inst.PTP.DeviceName == args.Dev {
				resp.ExitCode = 1
				resp.Output = "Device name is already in use"
				instances_mut.Unlock()
				return errors.New(resp.Output)
			}
		}
		instances_mut.Unlock()
	}

	var exists bool
	instances_mut.Lock()
	_, exists = instances[args.Hash]
	instances_mut.Unlock()
	if !exists {
		resp.Output = resp.Output + "Lookup finished\n"
		if args.Key != "" {
			if len(args.Key) < 16 {
				args.Key += "0000000000000000"[:16-len(args.Key)]
			} else if len(args.Key) > 16 && len(args.Key) < 24 {
				args.Key += "000000000000000000000000"[:24-len(args.Key)]
			} else if len(args.Key) > 24 && len(args.Key) < 32 {
				args.Key += "00000000000000000000000000000000"[:32-len(args.Key)]
			} else if len(args.Key) > 32 {
				args.Key = args.Key[:32]
			}
		}

		var newInst instance
		newInst.ID = args.Hash
		newInst.Args = *args
		ptpInstance := ptp.StartP2PInstance(args.IP, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "", args.Fwd, args.Port)
		if ptpInstance == nil {
			resp.Output = resp.Output + "Failed to create P2P Instance"
			resp.ExitCode = 1
			return errors.New("Failed to create P2P Instance")
		}
		ptp.Log(ptp.Info, "Instance created")
		newInst.PTP = ptpInstance
		instances_mut.Lock()
		instances[args.Hash] = newInst
		instances_mut.Unlock()
		go ptpInstance.Run()
		if saveFile != "" {
			resp.Output = resp.Output + "Saving instance into file"
			saveInstances(saveFile)
		}
	} else {
		resp.Output = resp.Output + "Hash already in use\n"
	}
	return nil
}

// Stop is used to terminate a specific P2P instance
func (p *Procedures) Stop(args *StopArgs, resp *Response) error {
	resp.ExitCode = 0
	var exists bool
	instances_mut.Lock()
	_, exists = instances[args.Hash]
	if !exists {
		resp.ExitCode = 1
		resp.Output = "Instance with hash " + args.Hash + " was not found"
		instances_mut.Unlock()
	} else {
		resp.Output = "Shutting down " + args.Hash
		instances[args.Hash].PTP.StopInstance()
		delete(instances, args.Hash)
		instances_mut.Unlock()
		saveInstances(saveFile)
	}
	return nil
}

// Show is used to output information about instances
func (p *Procedures) Show(args *RunArgs, resp *Response) error {
	if args.Hash != "" {
		instances_mut.Lock()
		swarm, exists := instances[args.Hash]
		resp.ExitCode = 0
		if exists {
			if args.IP != "" {
				swarm.PTP.PeersLock.Lock()
				for _, peer := range swarm.PTP.NetworkPeers {
					if peer.PeerLocalIP.String() == args.IP {
						if peer.State == ptp.PeerStateConnected {
							resp.ExitCode = 0
							resp.Output = "Integrated with " + args.IP
							swarm.PTP.PeersLock.Unlock()
							instances_mut.Unlock()
							runtime.Gosched()
							return nil
						}
					}
				}
				swarm.PTP.PeersLock.Unlock()
				instances_mut.Unlock()
				runtime.Gosched()
				resp.ExitCode = 1
				resp.Output = "Not yet integrated with " + args.IP
				return nil
			}
			resp.Output = "< Peer ID >\t< IP >\t< Endpoint >\t< HW >\n"
			swarm.PTP.PeersLock.Lock()
			for _, peer := range swarm.PTP.NetworkPeers {
				resp.Output = resp.Output + peer.ID + "\t"
				resp.Output = resp.Output + peer.PeerLocalIP.String() + "\t"
				resp.Output = resp.Output + peer.Endpoint.String() + "\t"
				resp.Output = resp.Output + peer.PeerHW.String() + "\n"
			}
			swarm.PTP.PeersLock.Unlock()
			runtime.Gosched()
		} else {
			resp.Output = "Specified environment was not found: " + args.Hash
			resp.ExitCode = 1
		}
		instances_mut.Unlock()
	} else {
		resp.ExitCode = 0
		instances_mut.Lock()
		inst_len := len(instances)
		instances_mut.Unlock()
		if inst_len == 0 {
			resp.Output = "No instances was found"
		}
		instances_mut.Lock()
		for key, inst := range instances {
			if inst.PTP != nil {
				resp.Output = resp.Output + "\t" + inst.PTP.Mac + "\t" + inst.PTP.IP + "\t" + key
			} else {
				resp.Output = resp.Output + "\tUnknown\tUnknown\t" + key
			}
			resp.Output = resp.Output + "\n"
		}
		instances_mut.Unlock()
	}
	return nil
}

// Debug output debug information
func (p *Procedures) Debug(args *Args, resp *Response) error {
	resp.Output = "DEBUG INFO:\n"
	resp.Output += fmt.Sprintf("Number of gouroutines: %d\n", runtime.NumGoroutine())
	resp.Output += fmt.Sprintf("Instances information:\n")
	instances_mut.Lock()
	for _, ins := range instances {
		resp.Output += fmt.Sprintf("Hash: %s\n", ins.ID)
		resp.Output += fmt.Sprintf("ID: %s\n", ins.PTP.Dht.ID)
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
	instances_mut.Unlock()
	return nil
}

// Status displays information about instances, peers and their statuses
func (p *Procedures) Status(args *RunArgs, resp *Response) error {
	instances_mut.Lock()
	for _, ins := range instances {
		resp.Output += ins.ID + " | " + ins.PTP.IP + "\n"
		for _, peer := range ins.PTP.NetworkPeers {
			resp.Output += peer.ID + "|"
			resp.Output += peer.PeerLocalIP.String() + "|"
			resp.Output += "State:" + StringifyState(peer.State) + "|"
			if peer.LastError != "" {
				resp.Output += "LastError:" + peer.LastError
			}
			resp.Output += "\n"
		}
	}
	instances_mut.Unlock()
	return nil
}

// StringifyState extracts human-readable word that represents a peer status
func StringifyState(state ptp.PeerState) string {
	switch state {
	case ptp.PeerStateInit:
		return "Initializing"
	case ptp.PeerStateRequestedIP:
		return "Waiting for IP"
	case ptp.PeerStateConnectingDirectly:
		return "Trying direct connection"
	case ptp.PeerStateConnected:
		return "Connected"
	case ptp.PeerStateHandshaking:
		return "Handshaking"
	case ptp.PeerStateHandshakingFailed:
		return "Handshaking failed"
	case ptp.PeerStateWaitingForwarder:
		return "Waiting forwarder IP"
	case ptp.PeerStateHandshakingForwarder:
		return "Handshaking forwarder"
	case ptp.PeerStateDisconnect:
		return "Disconnected"
	case ptp.PeerStateStop:
		return "Stopped"
	}
	return "Unknown"
}
