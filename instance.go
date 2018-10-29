package main

import (
	"fmt"
	"net"
	"sync"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
)

// RunArgs is a list of arguments used at instance startup and
// some other RPC calls
type RunArgs struct {
	IP          string `json:"ip"`
	Mac         string `json:"mac"`
	Dev         string `json:"dev"`
	Hash        string `json:"hash"`
	Dht         string `json:"dht"`
	Keyfile     string `json:"keyfile"`
	Key         string `json:"key"`
	TTL         string `json:"ttl"`
	Fwd         bool   `json:"fwd"`
	Port        int    `json:"port"`
	LastSuccess time.Time
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
	saveFile     string
	usedIPs      []string
	saveFileLock sync.RWMutex
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

// // encodeInstancesDeprecated will be removed in the next major release
// // TODO: Remove this code in version 9
// func (p *InstanceList) encodeInstancesDeprecated() []byte {
// 	var savedInstances []RunArgs
// 	instances := p.get()
// 	for _, instance := range instances {
// 		savedInstances = append(savedInstances, instance.Args)
// 	}
// 	var result bytes.Buffer
// 	flag := false
// 	for _, instance := range savedInstances {
// 		if flag == true {
// 			result.WriteString("|||")
// 		}
// 		result.WriteString(instance.IP + "~")
// 		result.WriteString(instance.Mac + "~")
// 		result.WriteString(instance.Dev + "~")
// 		result.WriteString(instance.Hash + "~")
// 		result.WriteString(instance.Dht + "~")
// 		result.WriteString(instance.Keyfile + "~")
// 		result.WriteString(instance.Key + "~")
// 		result.WriteString(instance.TTL + "~")
// 		var Fwd int
// 		if instance.Fwd == true {
// 			Fwd = 1
// 		}
// 		result.WriteString(strconv.Itoa(Fwd) + "~")
// 		result.WriteString(strconv.Itoa(instance.Port))
// 		flag = true
// 	}
// 	return result.Bytes()
// }

// // encode will generate YAML
// func (p *InstanceList) encode() ([]byte, error) {
// 	var data []saveData
// 	instances := p.get()
// 	if len(instances) == 0 {
// 		return nil, nil
// 	}
// 	for _, instance := range instances {
// 		ls, _ := instance.Args.LastSuccess.MarshalText()
// 		s := saveData{
// 			IP:          instance.Args.IP,
// 			Mac:         instance.Args.Mac,
// 			Dev:         instance.Args.Dev,
// 			Hash:        instance.Args.Hash,
// 			Keyfile:     instance.Args.Keyfile,
// 			Key:         instance.Args.Key,
// 			TTL:         instance.Args.TTL,
// 			LastSuccess: string(ls),
// 		}
// 		data = append(data, s)
// 	}
// 	output, err := yaml.Marshal(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	fmt.Println(string(output))
// 	return output, nil
// }

// // decode will accept YAML and generate a slice of RunArgs
// func (p *InstanceList) decode(data []byte) ([]RunArgs, error) {
// 	fmt.Println(string(data))
// 	var saved []saveData
// 	err := yaml.Unmarshal(data, &saved)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var args []RunArgs
// 	for _, i := range saved {
// 		var t time.Time
// 		err := t.UnmarshalText([]byte(i.LastSuccess))
// 		if err != nil {
// 			t = time.Unix(0, 0)
// 		}
// 		item := RunArgs{
// 			IP:          i.IP,
// 			Mac:         i.Mac,
// 			Dev:         i.Dev,
// 			Hash:        i.Hash,
// 			Keyfile:     i.Keyfile,
// 			Key:         i.Key,
// 			TTL:         i.TTL,
// 			LastSuccess: t,
// 		}
// 		args = append(args, item)
// 	}
// 	return args, nil
// }

// // TODO: Remove in version 9
// func (p *InstanceList) decodeInstances(data []byte) ([]RunArgs, error) {
// 	var args []RunArgs
// 	b := bytes.Buffer{}
// 	b.Write(data)
// 	d := gob.NewDecoder(&b)
// 	err := d.Decode(&args)
// 	if err != nil {
// 		blocksOfInstancesOld := bytes.Split(data, bytes.NewBufferString("|~|").Bytes())
// 		blocksOfInstances := bytes.Split(data, bytes.NewBufferString("|||").Bytes())
// 		if len(blocksOfInstancesOld) == len(blocksOfInstances) {
// 			if len(blocksOfInstancesOld) != 1 || len(blocksOfInstances) != 1 {
// 				return nil, fmt.Errorf("Unexpected error in decoding process")
// 			}
// 		} else {
// 			if len(blocksOfInstancesOld) > len(blocksOfInstances) {
// 				blocksOfInstances = blocksOfInstancesOld
// 			}
// 		}
// 		for _, str := range blocksOfInstances {
// 			blocksOfArguments := bytes.Split(str, bytes.NewBufferString("~").Bytes())
// 			if len(blocksOfArguments) != 10 {
// 				return nil, fmt.Errorf("Couldn't decode the instances")
// 			}
// 			var item RunArgs
// 			item.IP = string(blocksOfArguments[0])
// 			item.Mac = string(blocksOfArguments[1])
// 			item.Dev = string(blocksOfArguments[2])
// 			item.Hash = string(blocksOfArguments[3])
// 			item.Dht = string(blocksOfArguments[4])
// 			item.Keyfile = string(blocksOfArguments[5])
// 			item.Key = string(blocksOfArguments[6])
// 			item.TTL = string(blocksOfArguments[7])
// 			item.Fwd = false
// 			if string(blocksOfArguments[8]) == "1" {
// 				item.Fwd = true
// 			}
// 			item.Port, err = strconv.Atoi(string(blocksOfArguments[9]))
// 			if err != nil {
// 				return nil, fmt.Errorf("Couldn't decode the Port: %v", err)
// 			}
// 			args = append(args, item)
// 		}
// 	}
// 	return args, nil
// }

// // Calls encodeInstances() and saves results into specified file
// // Return number of bytes written and error if any
// func (p *InstanceList) saveInstances(filename string) (int, error) {
// 	saveFileLock.Lock()
// 	defer saveFileLock.Unlock()
// 	file, err := os.Open(filename)
// 	if err == nil {
// 		file.Close()
// 		os.Remove(filename)
// 	}
// 	file, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0700)
// 	if err != nil {
// 		return 0, err
// 	}
// 	defer file.Close()
// 	stat, _ := file.Stat()
// 	if stat.Size() > 0 {
// 		auxiliary := make([]byte, 100000)
// 		len, err := file.Read(auxiliary)
// 		if err != nil {
// 			return 0, err
// 		}
// 		auxiliary = bytes.Trim(auxiliary, "\x00")
// 		return 0, fmt.Errorf("SaveFile was not empty: %+v %+v", len, bytes.NewBuffer(auxiliary).String())
// 	}
// 	data, err := p.encode()
// 	if err != nil {
// 		return 0, fmt.Errorf("Failed to encode instances: %s", err.Error())
// 	}
// 	s, err := file.Write(data)
// 	if err != nil {
// 		return s, err
// 	}
// 	return s, nil
// }

// func (p *InstanceList) loadInstances(filename string) ([]RunArgs, error) {
// 	saveFileLock.RLock()
// 	defer saveFileLock.RUnlock()
// 	file, err := os.Open(filename)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer file.Close()
// 	data := make([]byte, 100000)
// 	_, err = file.Read(data)
// 	if err != nil {
// 		return nil, err
// 	}
// 	data = bytes.Trim(data, "\x00") // TODO: add more security to this
// 	if string(data[0]) == "-" || string(data[0]) == "[" {
// 		return p.decode(data)
// 	}
// 	// TODO: This code is deprecated and must be removed in version 9
// 	return p.decodeInstances(data)
// }

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

// Daemon is an object of RPC procedures
type Daemon struct {
	Instances  *InstanceList
	Restore    *Restore
	OutboundIP net.IP
}

// init will initialize daemon, instnaces and restore subsystems
func (d *Daemon) init(saveFile string) error {
	d.Instances = new(InstanceList)
	d.Instances.init()
	d.Restore = new(Restore)
	err := d.Restore.init(saveFile)
	if err != nil {
		return err
	}
	return nil
}

// Execute is a dummy method used for tests
func (d *Daemon) Execute(args *Args, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = ""
	return nil
}
