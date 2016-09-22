package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"runtime"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
)

var InstanceLock bool = false

func WaitLock() {
	for InstanceLock {
		time.Sleep(100 * time.Microsecond)
	}
}

func Lock() {
	InstanceLock = true
}

func Unlock() {
	InstanceLock = false
}

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

type Instance struct {
	PTP  *ptp.PTPCloud
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
	return args, err
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
	return loadedInstances, err
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

func (p *Procedures) SetLog(args *NameValueArg, resp *Response) error {
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
	WaitLock()
	Lock()
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
	Unlock()
	return nil
}

func (p *Procedures) Execute(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
	return nil
}

func (p *Procedures) Run(args *RunArgs, resp *Response) error {
	WaitLock()
	Lock()
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash + "\n"
	defer Unlock()

	// Validate if interface name is unique
	if args.Dev != "" {
		for _, inst := range Instances {
			if inst.PTP.DeviceName == args.Dev {
				resp.ExitCode = 1
				resp.Output = "Device name is already in use"
				Unlock()
				return errors.New(resp.Output)
			}
		}
	}

	var exists bool
	_, exists = Instances[args.Hash]
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

		var newInst Instance
		newInst.ID = args.Hash
		newInst.Args = *args
		Instances[args.Hash] = newInst
		ptpInstance := ptp.StartP2PInstance(args.IP, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "", args.Fwd, args.Port)
		if ptpInstance == nil {
			delete(Instances, args.Hash)
			resp.Output = resp.Output + "Failed to create P2P Instance"
			resp.ExitCode = 1
			Unlock()
			return errors.New("Failed to create P2P Instance")
		}
		newInst.PTP = ptpInstance
		Instances[args.Hash] = newInst
		go ptpInstance.Run()
		if SaveFile != "" {
			resp.Output = resp.Output + "Saving instance into file"
			SaveInstances(SaveFile)
		}
	} else {
		resp.Output = resp.Output + "Hash already in use\n"
	}
	Unlock()
	return nil
}

func (p *Procedures) Stop(args *StopArgs, resp *Response) error {
	WaitLock()
	Lock()
	defer Unlock()
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
		SaveInstances(SaveFile)
	}
	Unlock()
	return nil
}

func (p *Procedures) Show(args *RunArgs, resp *Response) error {
	if args.Hash != "" {
		swarm, exists := Instances[args.Hash]
		resp.ExitCode = 0
		if exists {
			if args.IP != "" {
				swarm.PTP.PeersLock.Lock()
				for _, peer := range swarm.PTP.NetworkPeers {
					if peer.PeerLocalIP.String() == args.IP {
						if peer.State == ptp.P_CONNECTED {
							resp.ExitCode = 0
							resp.Output = "Integrated with " + args.IP
							swarm.PTP.PeersLock.Unlock()
							runtime.Gosched()
							return nil
						}
					}
				}
				swarm.PTP.PeersLock.Unlock()
				runtime.Gosched()
				resp.ExitCode = 1
				resp.Output = "Not yet integrated with " + args.IP
				return nil
			} else {
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
			}
		} else {
			resp.Output = "Specified environment was not found: " + args.Hash
			resp.ExitCode = 1
		}
	} else {
		resp.ExitCode = 0
		if len(Instances) == 0 {
			resp.Output = "No instances was found"
		}
		for key, inst := range Instances {
			if inst.PTP != nil {
				resp.Output = resp.Output + "\t" + inst.PTP.Mac + "\t" + inst.PTP.IP + "\t" + key
			} else {
				resp.Output = resp.Output + "\tUnknown\tUnknown\t" + key
			}
			resp.Output = resp.Output + "\n"
		}
	}
	return nil
}

func (p *Procedures) Debug(args *Args, resp *Response) error {
	resp.Output = "DEBUG INFO:\n"
	resp.Output += fmt.Sprintf("Number of gouroutines: %d\n", runtime.NumGoroutine())
	resp.Output += fmt.Sprintf("Instances information:\n")
	for _, ins := range Instances {
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
	return nil
}

func (p *Procedures) Status(args *RunArgs, resp *Response) error {
	for _, ins := range Instances {
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
	return nil
}

func StringifyState(state ptp.PeerState) string {
	switch state {
	case ptp.P_INIT:
		return "Initializing"
	case ptp.P_REQUESTED_IP:
		return "Waiting for IP"
	case ptp.P_CONNECTING_DIRECTLY:
		return "Trying direct connection"
	case ptp.P_CONNECTED:
		return "Connected"
	case ptp.P_HANDSHAKING:
		return "Handshaking"
	case ptp.P_HANDSHAKING_FAILED:
		return "Handshaking failed"
	case ptp.P_WAITING_FORWARDER:
		return "Waiting forwarder IP"
	case ptp.P_HANDSHAKING_FORWARDER:
		return "Handshaking forwarder"
	case ptp.P_DISCONNECT:
		return "Disconnected"
	case ptp.P_STOP:
		return "Stopped"
	}
	return "Unknown"
}
