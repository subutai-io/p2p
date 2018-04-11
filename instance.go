package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"sync"

	"strconv"

	ptp "github.com/subutai-io/p2p/lib"
)

// RunArgs is a list of arguments used at instance startup and
// some other RPC calls
type RunArgs struct {
	IP      string `json:"ip"`
	Mac     string `json:"mac"`
	Dev     string `json:"dev"`
	Hash    string `json:"hash"`
	Dht     string `json:"dht"`
	Keyfile string `json:"keyfile"`
	Key     string `json:"key"`
	TTL     string `json:"ttl"`
	Fwd     bool   `json:"fwd"`
	Port    int    `json:"port"`
}

type ShowArgs struct {
	Hash       string `json:"hash"`
	IP         string `json:"ip"`
	Interfaces bool   `json:"interfaces"`
	All        bool   `json:"all"`
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

func (p *InstanceList) init() {
	p.instances = make(map[string]*P2PInstance)
}

func (p *InstanceList) update(id string, inst *P2PInstance) error {
	return p.operate(InstWrite, id, inst)
}

func (p *InstanceList) delete(id string) error {
	return p.operate(InstDelete, id, nil)
}

func (p *InstanceList) get() map[string]*P2PInstance {
	result := make(map[string]*P2PInstance)
	p.lock.RLock()
	for id, inst := range p.instances {
		result[id] = inst
	}
	p.lock.RUnlock()
	return result
}

func (p *InstanceList) getInstance(id string) *P2PInstance {
	p.lock.RLock()
	inst, exists := p.instances[id]
	p.lock.RUnlock()
	if !exists {
		return nil
	}
	return inst
}

func (p *InstanceList) encodeInstances() ([]byte, error) {
	var savedInstances []RunArgs
	instances := p.get()
	for _, instance := range instances {
		savedInstances = append(savedInstances, instance.Args)
	}
	var result bytes.Buffer
	flag := false
	for _, instance := range savedInstances {
		if flag == true {
			result.WriteString("|~|")
		}
		result.WriteString(instance.IP + "~")
		result.WriteString(instance.Mac + "~")
		result.WriteString(instance.Dev + "~")
		result.WriteString(instance.Hash + "~")
		result.WriteString(instance.Dht + "~")
		result.WriteString(instance.Keyfile + "~")
		result.WriteString(instance.Key + "~")
		result.WriteString(instance.TTL + "~")
		fmt.Println(result.String())
		var Fwd int
		if instance.Fwd == true {
			Fwd = 1
		}
		result.WriteString(strconv.Itoa(Fwd) + "~")
		result.WriteString(strconv.Itoa(instance.Port))
		flag = true
	}
	return result.Bytes(), nil
}

func (p *InstanceList) decodeInstances(data []byte) ([]RunArgs, error) {
	var args []RunArgs
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(&args)
	if err != nil {
		blocksOfInstances := bytes.Split(data, bytes.NewBufferString("|~|").Bytes())
		for _, str := range blocksOfInstances {
			blocksOfArguments := bytes.Split(str, bytes.NewBufferString("~").Bytes())
			if len(blocksOfArguments) != 10 {
				return nil, fmt.Errorf("Couldn't decode the instances")
			}
			var item RunArgs
			item.IP = string(blocksOfArguments[0])
			item.Mac = string(blocksOfArguments[1])
			item.Dev = string(blocksOfArguments[2])
			item.Hash = string(blocksOfArguments[3])
			item.Dht = string(blocksOfArguments[4])
			item.Keyfile = string(blocksOfArguments[5])
			item.Key = string(blocksOfArguments[6])
			item.TTL = string(blocksOfArguments[7])
			item.Fwd = false
			if string(blocksOfArguments[8]) == "true" {
				item.Fwd = true
			}
			item.Port, err = strconv.Atoi(string(blocksOfArguments[9]))
			if err != nil {
				return nil, fmt.Errorf("Couldn't decode the Port: %v", err)
			}
			args = append(args, item)
		}
	}
	return args, nil
}

// Calls encodeInstances() and saves results into specified file
// Return number of bytes written and error if any
func (p *InstanceList) saveInstances(filename string) (int, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	data, err := p.encodeInstances()
	if err != nil {
		return 0, err
	}
	s, err := file.Write(data)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (p *InstanceList) loadInstances(filename string) ([]RunArgs, error) {
	var loadedInstances []RunArgs
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	data := make([]byte, 100000)
	_, err = file.Read(data)
	data = bytes.Trim(data, "\x00") // TODO: add more security to this
	if err != nil {
		return loadedInstances, err
	}
	loadedInstances, err = p.decodeInstances(data)
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
	p.Instances.init()
	p.SaveFile = saveFile
}

// Execute is a dummy method used for tests
func (p *Daemon) Execute(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
	return nil
}
