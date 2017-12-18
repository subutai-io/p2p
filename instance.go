package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
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

type ShowArgs struct {
	Hash       string
	IP         string
	Interfaces bool
	All        bool
}

// P2PInstance is a holder for P2P instances started by daemon
type P2PInstance struct {
	PTP  *ptp.PeerToPeer
	ID   string
	Args RunArgs
}

var (
	saveFile string
	usedIPs  []string
)

type InstOperation int

// Type of instance operations
const (
	InstWrite  InstOperation = 0
	InstDelete InstOperation = 1
)

type InstanceList struct {
	instances map[string]*P2PInstance
	lock      sync.RWMutex
}

func (p *InstanceList) operate(action InstOperation, id string, inst *P2PInstance) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	if action == InstWrite {
		p.instances[id] = inst
	} else if action == InstDelete {
		_, exists := p.instances[id]
		if !exists {
			return fmt.Errorf("Specified instance has not been found")
		}
		delete(p.instances, id)
	}
	return nil
}

func (p *InstanceList) Init() {
	p.instances = make(map[string]*P2PInstance)
}

func (p *InstanceList) Update(id string, inst *P2PInstance) error {
	return p.operate(InstWrite, id, inst)
}

func (p *InstanceList) Delete(id string) error {
	return p.operate(InstDelete, id, nil)
}

func (p *InstanceList) Get() map[string]*P2PInstance {
	result := make(map[string]*P2PInstance)
	p.lock.RLock()
	for id, inst := range p.instances {
		result[id] = inst
	}
	p.lock.RUnlock()
	return result
}

func (p *InstanceList) GetInstance(id string) *P2PInstance {
	p.lock.RLock()
	inst, exists := p.instances[id]
	p.lock.RUnlock()
	if !exists {
		return nil
	}
	return inst
}

func (p *InstanceList) EncodeInstances() ([]byte, error) {
	var savedInstances []RunArgs
	instances := p.Get()
	for _, inst := range instances {
		savedInstances = append(savedInstances, inst.Args)
	}
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(savedInstances)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to encode instances: %v", err)
		return []byte(""), err
	}
	return b.Bytes(), nil
}

func (p *InstanceList) DecodeInstances(data []byte) ([]RunArgs, error) {
	var args []RunArgs
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(&args)
	return args, err
}

// Calls encodeInstances() and saves results into specified file
// Return number of bytes written and error if any
func (p *InstanceList) SaveInstances(filename string) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	data, err := p.EncodeInstances()
	if err != nil {
		return 0, err
	}
	s, err := file.Write(data)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (p *InstanceList) LoadInstances(filename string) ([]RunArgs, error) {
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

	loadedInstances, err = p.DecodeInstances(data)
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
type Daemon struct {
	Instances  *InstanceList
	SaveFile   string
	OutboundIP net.IP
}

func (p *Daemon) Initialize(saveFile string) {
	p.Instances = new(InstanceList)
	p.Instances.Init()
	p.SaveFile = saveFile
}

// SetLog modifies specific option
func (p *Daemon) SetLog(args *NameValueArg, resp *Response) error {
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
func (p *Daemon) AddKey(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	if args.Hash == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified hash"
	}
	if args.Key == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified key"
	}
	inst := p.Instances.GetInstance(args.Hash)
	if inst == nil {
		resp.ExitCode = 1
		resp.Output = "No instances with specified hash were found"
	}
	if resp.ExitCode == 0 {
		resp.Output = "New key added"
		var newKey ptp.CryptoKey

		newKey = inst.PTP.Crypter.EnrichKeyValues(newKey, args.Key, args.TTL)
		inst.PTP.Crypter.Keys = append(inst.PTP.Crypter.Keys, newKey)
		p.Instances.Update(args.Hash, inst)
	}
	return nil
}

// Execute is a dummy method used for tests
func (p *Daemon) Execute(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
	return nil
}

// Run starts a P2P instance
func (p *Daemon) Run(args *RunArgs, resp *Response) error {
	args.Dht = DefaultDHT
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash + "\n"

	// Validate if interface name is unique
	if args.Dev != "" {
		instances := p.Instances.Get()
		for _, inst := range instances {
			if inst.PTP.Interface.Name == args.Dev {
				resp.ExitCode = 1
				resp.Output = "Device name is already in use"
				return errors.New(resp.Output)
			}
		}
	}

	inst := p.Instances.GetInstance(args.Hash)
	if inst == nil {
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

		newInst := new(P2PInstance)
		newInst.ID = args.Hash
		newInst.Args = *args
		newInst.PTP = ptp.New(args.IP, args.Mac, args.Dev, "", args.Hash, args.Dht, args.Keyfile, args.Key, args.TTL, "", args.Fwd, args.Port, usedIPs, OutboundIP)
		if newInst.PTP == nil {
			resp.Output = resp.Output + "Failed to create P2P Instance"
			resp.ExitCode = 1
			return errors.New("Failed to create P2P Instance")
		}

		// Saving interface name
		infFound := false
		for _, inf := range InterfaceNames {
			if inf == newInst.PTP.Interface.Name {
				infFound = true
			}
		}
		if !infFound && newInst.PTP.Interface.Name != "" {
			InterfaceNames = append(InterfaceNames, newInst.PTP.Interface.Name)
		}

		usedIPs = append(usedIPs, newInst.PTP.Interface.IP.String())
		ptp.Log(ptp.Info, "Instance created")
		p.Instances.Update(args.Hash, newInst)

		go newInst.PTP.Run()
		if p.SaveFile != "" {
			resp.Output = resp.Output + "Saving instance into file"
			p.Instances.SaveInstances(p.SaveFile)
		}
	} else {
		resp.Output = resp.Output + "Hash already in use\n"
	}
	return nil
}

// Stop is used to terminate a specific P2P instance
func (p *Daemon) Stop(args *StopArgs, resp *Response) error {
	resp.ExitCode = 0
	inst := p.Instances.GetInstance(args.Hash)
	if inst == nil {
		resp.ExitCode = 1
		resp.Output = "Instance with hash " + args.Hash + " was not found"
	} else {
		ip := inst.PTP.Interface.IP.String()
		resp.Output = "Shutting down " + args.Hash
		inst.PTP.StopInstance()
		p.Instances.Delete(args.Hash)
		p.Instances.SaveInstances(p.SaveFile)
		k := 0
		for k, i := range usedIPs {
			if i != ip {
				usedIPs[k] = i
				k++
			}
		}
		usedIPs = usedIPs[:k]
	}
	return nil
}

// Show is used to output information about instances
func (p *Daemon) Show(args *ShowArgs, resp *Response) error {
	if args.Hash != "" {
		inst := p.Instances.GetInstance(args.Hash)
		resp.ExitCode = 0
		if inst != nil {
			peers := inst.PTP.Peers.Get()
			if args.IP != "" {
				for _, peer := range peers {
					if peer.PeerLocalIP.String() == args.IP {
						if peer.State == ptp.PeerStateConnected {
							resp.ExitCode = 0
							resp.Output = "Integrated with " + args.IP
							return nil
						}
					}
				}
				resp.ExitCode = 1
				resp.Output = "Not yet integrated with " + args.IP
				return nil
			}
			resp.Output = "< Peer ID >\t< IP >\t< Endpoint >\t< HW >\n"
			for _, peer := range peers {
				resp.Output = resp.Output + peer.ID + "\t"
				resp.Output = resp.Output + peer.PeerLocalIP.String() + "\t"
				resp.Output = resp.Output + peer.Endpoint.String() + "\t"
				resp.Output = resp.Output + peer.PeerHW.String() + "\n"
			}
		} else {
			resp.Output = "Specified environment was not found: " + args.Hash
			resp.ExitCode = 1
		}
	} else if args.Interfaces {
		if !args.All {
			instances := p.Instances.Get()
			for _, inst := range instances {
				if inst.PTP != nil {
					resp.Output = resp.Output + inst.PTP.Interface.Name
				}
				resp.Output = resp.Output + "\n"
			}
		} else {
			for _, inf := range InterfaceNames {
				resp.Output = resp.Output + inf + "\n"
			}
		}
	} else {
		resp.ExitCode = 0
		instances := p.Instances.Get()
		instLen := len(instances)
		if instLen == 0 {
			resp.Output = "No instances was found"
		}
		for key, inst := range instances {
			if inst.PTP != nil {
				resp.Output = resp.Output + "\t" + inst.PTP.Interface.Mac.String() + "\t" + inst.PTP.Interface.IP.String() + "\t" + key
			} else {
				resp.Output = resp.Output + "\tUnknown\tUnknown\t" + key
			}
			resp.Output = resp.Output + "\n"
		}
	}
	return nil
}

// Debug output debug information
func (p *Daemon) Debug(args *Args, resp *Response) error {
	resp.Output = "DEBUG INFO:\n"
	resp.Output = fmt.Sprintf("Version: %s\n", AppVersion)
	resp.Output += fmt.Sprintf("Number of gouroutines: %d\n", runtime.NumGoroutine())
	resp.Output += fmt.Sprintf("Instances information:\n")
	instances := p.Instances.Get()
	for _, inst := range instances {
		resp.Output += fmt.Sprintf("Bootstrap nodes:\n")
		for _, conn := range inst.PTP.Dht.Connections {
			resp.Output += fmt.Sprintf("\t%s\n", conn.RemoteAddr().String())
		}
		resp.Output += fmt.Sprintf("Hash: %s\n", inst.ID)
		resp.Output += fmt.Sprintf("ID: %s\n", inst.PTP.Dht.ID)
		resp.Output += fmt.Sprintf("UDP Port: %d\n", inst.PTP.UDPSocket.GetPort())
		resp.Output += fmt.Sprintf("Interface %s, HW Addr: %s, IP: %s\n", inst.PTP.Interface.Name, inst.PTP.Interface.Mac.String(), inst.PTP.Interface.IP.String())
		resp.Output += fmt.Sprintf("Proxies:\n")
		if len(inst.PTP.Proxies) == 0 {
			resp.Output += fmt.Sprintf("\tNo proxies in use\n")
		}
		for _, proxy := range inst.PTP.Proxies {
			resp.Output += fmt.Sprintf("\tProxy address: %s\n", proxy.Addr.String())
		}
		resp.Output += fmt.Sprintf("Peers:\n")

		peers := inst.PTP.Peers.Get()
		for _, peer := range peers {
			resp.Output += fmt.Sprintf("\t--- %s ---\n", peer.ID)
			if peer.PeerLocalIP == nil {
				resp.Output += "\t\tNo IP assigned\n"

			} else if peer.PeerHW == nil {
				resp.Output += "\t\tNo MAC assigned\n"
			} else {
				resp.Output += fmt.Sprintf("\t\tHWAddr: %s\n", peer.PeerHW.String())
				resp.Output += fmt.Sprintf("\t\tIP: %s\n", peer.PeerLocalIP.String())
				resp.Output += fmt.Sprintf("\t\tEndpoint: %s\n", peer.Endpoint)
				resp.Output += fmt.Sprintf("\t\tPeer Address: %s\n", peer.PeerAddr.String())
				proxyInUse := "No"
				if peer.IsUsingTURN {
					proxyInUse = "Yes"
				}
				resp.Output += fmt.Sprintf("\t\tUsing proxy: %s\n", proxyInUse)
				//resp.Output += fmt.Sprintf("\t\tProxy ID: %d\n", peer.ProxyID)
			}
			resp.Output += fmt.Sprintf("\t--- End of %s ---\n", peer.ID)
		}
	}
	return nil
}

// Status displays information about instances, peers and their statuses
func (p *Daemon) Status(args *RunArgs, resp *Response) error {
	instances := p.Instances.Get()
	for _, inst := range instances {
		resp.Output += inst.ID + " | " + inst.PTP.Interface.IP.String() + "\n"
		peers := inst.PTP.Peers.Get()
		for _, peer := range peers {
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
