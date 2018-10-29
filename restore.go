package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/subutai-io/p2p/lib"

	yaml "gopkg.in/yaml.v2"
)

// Restore represents dump/restore subsystem for instances
// This class keeps a list of so-called "save entries", which is
// a binding for instance information
// Restore system saves entries in a YAML-formatted save file, specified
// as an argument on daemon launch using `--save`
type Restore struct {
	entries  []saveEntry
	filepath string
	lock     sync.RWMutex
	active   bool
}

// saveEntry is a YAML binding for data save file
type saveEntry struct {
	IP          string `yaml:"ip"`
	Mac         string `yaml:"mac"`
	Dev         string `yaml:"dev"`
	Hash        string `yaml:"hash"`
	Keyfile     string `yaml:"keyfile"`
	Key         string `yaml:"key"`
	TTL         string `yaml:"ttl"`
	LastSuccess string `yaml:"last_success"`
	Enabled     bool
}

// init will initialize restore subsystem by checking if
// file exists and can be modified
func (r *Restore) init(filepath string) error {
	if filepath == "" {
		r.active = false
		return nil
	}
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		return err
	}
	file.Close()
	r.filepath = filepath
	r.active = true
	return nil
}

// save will write dump of entries into a save file
func (r *Restore) save() error {
	r.lock.Lock()
	file, err := os.OpenFile(r.filepath, os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		r.lock.Unlock()
		return err
	}
	data, err := r.encode()
	if err != nil {
		r.lock.Unlock()
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		r.lock.Unlock()
		return err
	}
	r.lock.Unlock()
	return nil
}

// load will read save file and unmarshal saved entries
func (r *Restore) load() error {
	r.lock.Lock()
	data, err := ioutil.ReadFile(r.filepath)
	if err != nil {
		r.lock.Unlock()
		return err
	}
	r.lock.Unlock()
	data = bytes.Trim(data, "\x00") // TODO: add more security to this
	if len(data) == 0 {
		return nil
	}
	if string(data[0]) == "-" || string(data[0]) == "[" {
		r.decode(data)
		return nil
	}
	// TODO: This code is deprecated and must be removed in version 9
	return r.decodeInstances(data)
}

// addInstance will create new save file entry from instance
func (r *Restore) addInstance(inst *P2PInstance) error {
	ls, _ := inst.Args.LastSuccess.MarshalText()
	return r.addEntry(saveEntry{
		IP:          inst.Args.IP,
		Mac:         inst.Args.Mac,
		Dev:         inst.Args.Dev,
		Hash:        inst.Args.Hash,
		Keyfile:     inst.Args.Keyfile,
		Key:         inst.Args.Key,
		TTL:         inst.Args.TTL,
		LastSuccess: string(ls),
		Enabled:     true,
	})
}

// addEntry will create new save entry if it's unique by hash
func (r *Restore) addEntry(entry saveEntry) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for _, e := range r.entries {
		if e.Hash == entry.Hash {
			return fmt.Errorf("Instance %s already in list of saved entries", entry.Hash)
		}
	}

	r.entries = append(r.entries, entry)
	return nil
}

func (r *Restore) removeEntry(hash string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, e := range r.entries {
		if e.Hash == hash {
			r.entries = append(r.entries[:i], r.entries[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Can't delete save entry: %s not found", hash)
}

func (r *Restore) bumpInstance(hash string) error {
	r.lock.Lock()
	defer r.lock.Unlock()
	for i, e := range r.entries {
		if e.Hash == hash {
			ls, _ := time.Now().MarshalText()
			r.entries[i].LastSuccess = string(ls)
			return nil
		}
	}
	return fmt.Errorf("Can't update last success date for the instance: %s", hash)
}

func (r *Restore) disableStaleInstances(inst *P2PInstance) error {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for i, e := range r.entries {
		var t time.Time
		err := t.UnmarshalText([]byte(e.LastSuccess))
		if err != nil {
			ptp.Log(ptp.Error, "Failed to unmarshal date for save file entry %s. Disabling it")
			r.entries[i].Enabled = false
			continue
		}
		if time.Since(t) > time.Duration(time.Hour*24*20) {
			ptp.Log(ptp.Warning, "Instance %s was active more than 20 days ago")
			r.entries[i].Enabled = false
		}
	}
	return nil
}

// encode will generate YAML
func (r *Restore) encode() ([]byte, error) {
	if len(r.entries) == 0 {
		return nil, nil
	}
	var data []saveEntry
	for _, e := range r.entries {
		if e.Enabled {
			data = append(data, e)
		}
	}
	output, err := yaml.Marshal(data)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// decode will accept YAML and generate a slice of RunArgs
func (r *Restore) decode(data []byte) error {
	var saved []saveEntry
	err := yaml.Unmarshal(data, &saved)
	if err != nil {
		return err
	}
	r.lock.Lock()
	r.entries = saved
	r.lock.Unlock()
	return nil
}

// decodeInstances is an obsolet variant of instances unmarshal
// TODO: Remove in version 10
func (r *Restore) decodeInstances(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	var args []saveEntry
	b := bytes.Buffer{}
	b.Write(data)
	d := gob.NewDecoder(&b)
	err := d.Decode(&args)
	if err != nil {
		blocksOfInstancesOld := bytes.Split(data, bytes.NewBufferString("|~|").Bytes())
		blocksOfInstances := bytes.Split(data, bytes.NewBufferString("|||").Bytes())
		if len(blocksOfInstancesOld) == len(blocksOfInstances) {
			if len(blocksOfInstancesOld) != 1 || len(blocksOfInstances) != 1 {
				return fmt.Errorf("Unexpected error in decoding process")
			}
		} else {
			if len(blocksOfInstancesOld) > len(blocksOfInstances) {
				blocksOfInstances = blocksOfInstancesOld
			}
		}
		for _, str := range blocksOfInstances {
			blocksOfArguments := bytes.Split(str, bytes.NewBufferString("~").Bytes())
			if len(blocksOfArguments) != 10 {
				return fmt.Errorf("Couldn't decode the instances")
			}
			var item saveEntry
			item.IP = string(blocksOfArguments[0])
			item.Mac = string(blocksOfArguments[1])
			item.Dev = string(blocksOfArguments[2])
			item.Hash = string(blocksOfArguments[3])
			//item.Dht = string(blocksOfArguments[4])
			item.Keyfile = string(blocksOfArguments[5])
			item.Key = string(blocksOfArguments[6])
			item.TTL = string(blocksOfArguments[7])
			//item.Fwd = false
			// if string(blocksOfArguments[8]) == "1" {
			// 	item.Fwd = true
			// }
			//item.Port, err = strconv.Atoi(string(blocksOfArguments[9]))
			// if err != nil {
			// 	return fmt.Errorf("Couldn't decode the Port: %v", err)
			// }
			args = append(args, item)
		}
	}

	r.lock.Lock()
	ptp.Log(ptp.Info, "Decoded %d entries from the old format", len(args))
	r.entries = args
	r.lock.Unlock()
	return nil
}

// get will return slice of entries
func (r *Restore) get() []saveEntry {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.entries
}

func (r *Restore) isActive() bool {
	return r.active
}
