package main

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
	"os"
	"sync"
	"reflect"
	"net"
)

func TestOperate(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.operate(InstWrite, "instance", P2Pinstance)
	instances := instanceList.get()
	if len(instances) != 1 || instances["instance"] != P2Pinstance {
		t.Errorf("Failed to operate (1): operate didn't add an instance")
	}
	instanceList.operate(InstDelete, "instance", P2Pinstance)
	instances = instanceList.get()
	if len(instances) > 0 {
		t.Errorf("Failed to operate (2): operate didn't delete the instance")
	}
}

func TestInit(t *testing.T) {
	instanceList := new(InstanceList)
	instanceList.init()
	if instanceList.instances == nil {
		t.Errorf("Failed to init (1): init didn't initialize instances map")
	}
}

func TestUpdate(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	instances := instanceList.get()
	if len(instances) != 1 || instances["instance"] != P2Pinstance {
		t.Errorf("Failed to update (1): update didn't add an instance")
	}
}

func TestDelete(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	err := instanceList.delete("instance")
	if err != nil {
		t.Errorf("Failed to delete (1): %v", err)
	}
	err = instanceList.delete("instance")
	if err == nil {
		t.Errorf("Failed to delete (2): must have returned non-nil but returned nil")
	}
}

func TestGet(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data := instanceList.get()
	if len(data) != 1 {
		t.Errorf("Failed to get (1): get returned unexpected map")
	}
}

func TestGetInstance(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if instanceList.getInstance("instance") == nil {
		t.Errorf("Failed to get instance (1): getInstance returned nil, but instance exists")
	}
	if instanceList.getInstance("non-instance") != nil {
		t.Errorf("Failed to get instance (2): getInstance returned an instance, but instance does not exist")
	}
}

func TestEncodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "10.10.10.1~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0" {
		t.Errorf("Failed to encode instances (1): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0" {
		t.Errorf("Failed to encode instances (2): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "~~~~~~~~0~0" {
		t.Errorf("Failed to encode instances (3): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "10.10.10.1~~~~~~~~0~0" {
		t.Errorf("Failed to encode instances (4): encodedInstances incorrectly encoded the instanceList")
	}
	P2PinstanceFull := new(P2PInstance)
	P2PinstanceFull.Args.IP = "10.10.10.2"
	P2PinstanceFull.Args.Mac = "Mac"
	P2PinstanceFull.Args.Dev = "Dev"
	P2PinstanceFull.Args.Hash = "Hash"
	P2PinstanceFull.Args.Dht = "Dht"
	P2PinstanceFull.Args.Keyfile = "Keyfile"
	P2PinstanceFull.Args.Key = "Key"
	P2PinstanceFull.Args.TTL = "TTL"
	P2PinstanceFull.Args.Fwd = true
	P2PinstanceFull.Args.Port = 0
	instanceList.update("instanceFull", P2PinstanceFull)
	P2PinstanceSemi := new(P2PInstance)
	P2PinstanceSemi.Args.IP = "10.10.10.3"
	P2PinstanceSemi.Args.Mac = "Mac"
	P2PinstanceSemi.Args.Dev = "Dev"
	P2PinstanceSemi.Args.Hash = "Hash"
	P2PinstanceSemi.Args.Fwd = false
	P2PinstanceSemi.Args.Port = 0
	instanceList.update("instanceSemi", P2PinstanceSemi)
	encodedInstances := bytes.NewBuffer(instanceList.encodeInstances())
	parts := bytes.Split(encodedInstances.Bytes(), bytes.NewBufferString("|||").Bytes())
	set := make(map[string]bool)
	for i := 0; i < 3; i++ {
		set[bytes.NewBuffer(parts[i]).String()] = true
	}
	instanceString := "10.10.10.1~~~~~~~~0~0"
	instanceFullString := "10.10.10.2~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0"
	instanceSemiString := "10.10.10.3~Mac~Dev~Hash~~~~~0~0"
	if !set[instanceString] || !set[instanceFullString] || !set[instanceSemiString] {
		t.Errorf("Failed to encode instances (5): encodedInstances incorrectly encoded the instanceList")
	}
}
func TestDecodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data := instanceList.encodeInstances()
	_, err := instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (1): %v", err)
	}
	data[len(data)-1] = 65
	_, err = instanceList.decodeInstances(data)
	if err == nil {
		t.Errorf("Failed to decode instances (2): must have returned non-nil but returned nil")
	}
	data = make([]byte, 0)
	_, err = instanceList.decodeInstances(data)
	if err == nil {
		t.Errorf("Failed to decode instances (3): must have returned non-nil but returned nil")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (4): %v", err)
	}
	P2Pinstance = new(P2PInstance)
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (5): %v", err)
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (6): %v", err)
	}
}

func TestSaveInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	_, err := instanceList.saveInstances("/")
	if err == nil {
		t.Errorf("Failed to save instances (1): must have returned non-nil but returned nil")
	}
	_, err = instanceList.saveInstances("save-4.save")
	defer os.Remove("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (2): %v", err)
	}
	file, err := os.Open("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (3): %v", err)
	}
	auxiliary := make([]byte, 100000)
	file.Read(auxiliary)
	file.Close()
	auxiliary = bytes.Trim(auxiliary, "\x00")
	t.Log(bytes.NewBuffer(auxiliary).String())
	P2PinstanceSecond := new(P2PInstance)
	P2PinstanceSecond.Args.IP = "10.10.10.2"
	P2PinstanceSecond.Args.Mac = "Mac"
	P2PinstanceSecond.Args.Dev = "Dev"
	P2PinstanceSecond.Args.Hash = "Hash"
	P2PinstanceSecond.Args.Dht = "Dht"
	P2PinstanceSecond.Args.Keyfile = "Keyfile"
	P2PinstanceSecond.Args.Key = "Key"
	P2PinstanceSecond.Args.TTL = "TTL"
	P2PinstanceSecond.Args.Fwd = false
	P2PinstanceSecond.Args.Port = 0
	instanceList.update("instanceSecond", P2PinstanceSecond)
	_, err = instanceList.saveInstances("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (4): %v", err)
	}
	file, err = os.Open("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (5): %v", err)
	}
	auxiliary = make([]byte, 100000)
	file.Read(auxiliary)
	file.Close()
	auxiliary = bytes.Trim(auxiliary, "\x00")
	t.Log(bytes.NewBuffer(auxiliary).String())
}

func TestLoadInstances(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Println("This test is not supported on Windows")
	}
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	_, err := instanceList.loadInstances("/non-existing-file")
	if err == nil {
		t.Errorf("Failed to load instances (1): must have returned non-nil but returned nil")
	}
}

func TestInitialize(t *testing.T) {
	daemon := new(Daemon)
	daemon.Initialize("saveFile")
	if daemon.SaveFile != "saveFile" {
		t.Errorf("Failed to initialize (1): daemon couldn't initialize")
	}
}

func TestExecute(t *testing.T) {
	daemon := new(Daemon)
	args := new(Args)
	resp := new(Response)
	err := daemon.Execute(args, resp)
	if err != nil {
		t.Errorf("Failed to execute (1): %v", err)
	}
}
/*
Generated TestInstanceList_operate
Generated TestInstanceList_init
Generated TestInstanceList_update
Generated TestInstanceList_delete
Generated TestInstanceList_get
Generated TestInstanceList_getInstance
Generated TestInstanceList_encodeInstances
Generated TestInstanceList_decodeInstances
Generated TestInstanceList_saveInstances
Generated TestInstanceList_loadInstances
Generated TestDaemon_Initialize
Generated TestDaemon_Execute
package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"reflect"
	"runtime"
	"sync"
	"testing"
)
*/

/*
func TestOperate(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.operate(InstWrite, "instance", P2Pinstance)
	instances := instanceList.get()
	if len(instances) != 1 || instances["instance"] != P2Pinstance {
		t.Errorf("Failed to operate (1): operate didn't add an instance")
	}
	instanceList.operate(InstDelete, "instance", P2Pinstance)
	instances = instanceList.get()
	if len(instances) > 0 {
		t.Errorf("Failed to operate (2): operate didn't delete the instance")
	}
}

func TestInit(t *testing.T) {
	instanceList := new(InstanceList)
	instanceList.init()
	if instanceList.instances == nil {
		t.Errorf("Failed to init (1): init didn't initialize instances map")
	}
}

func TestUpdate(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	instances := instanceList.get()
	if len(instances) != 1 || instances["instance"] != P2Pinstance {
		t.Errorf("Failed to update (1): update didn't add an instance")
	}
}

func TestDelete(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	err := instanceList.delete("instance")
	if err != nil {
		t.Errorf("Failed to delete (1): %v", err)
	}
	err = instanceList.delete("instance")
	if err == nil {
		t.Errorf("Failed to delete (2): must have returned non-nil but returned nil")
	}
}

func TestGet(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data := instanceList.get()
	if len(data) != 1 {
		t.Errorf("Failed to get (1): get returned unexpected map")
	}
}

func TestGetInstance(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if instanceList.getInstance("instance") == nil {
		t.Errorf("Failed to get instance (1): getInstance returned nil, but instance exists")
	}
	if instanceList.getInstance("non-instance") != nil {
		t.Errorf("Failed to get instance (2): getInstance returned an instance, but instance does not exist")
	}
}

func TestEncodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "10.10.10.1~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0" {
		t.Errorf("Failed to encode instances (1): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0" {
		t.Errorf("Failed to encode instances (2): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "~~~~~~~~0~0" {
		t.Errorf("Failed to encode instances (3): encodedInstances incorrectly encoded the instanceList")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	if bytes.NewBuffer(instanceList.encodeInstances()).String() != "10.10.10.1~~~~~~~~0~0" {
		t.Errorf("Failed to encode instances (4): encodedInstances incorrectly encoded the instanceList")
	}
	P2PinstanceFull := new(P2PInstance)
	P2PinstanceFull.Args.IP = "10.10.10.2"
	P2PinstanceFull.Args.Mac = "Mac"
	P2PinstanceFull.Args.Dev = "Dev"
	P2PinstanceFull.Args.Hash = "Hash"
	P2PinstanceFull.Args.Dht = "Dht"
	P2PinstanceFull.Args.Keyfile = "Keyfile"
	P2PinstanceFull.Args.Key = "Key"
	P2PinstanceFull.Args.TTL = "TTL"
	P2PinstanceFull.Args.Fwd = true
	P2PinstanceFull.Args.Port = 0
	instanceList.update("instanceFull", P2PinstanceFull)
	P2PinstanceSemi := new(P2PInstance)
	P2PinstanceSemi.Args.IP = "10.10.10.3"
	P2PinstanceSemi.Args.Mac = "Mac"
	P2PinstanceSemi.Args.Dev = "Dev"
	P2PinstanceSemi.Args.Hash = "Hash"
	P2PinstanceSemi.Args.Fwd = false
	P2PinstanceSemi.Args.Port = 0
	instanceList.update("instanceSemi", P2PinstanceSemi)
	encodedInstances := bytes.NewBuffer(instanceList.encodeInstances())
	parts := bytes.Split(encodedInstances.Bytes(), bytes.NewBufferString("|||").Bytes())
	set := make(map[string]bool)
	for i := 0; i < 3; i++ {
		set[bytes.NewBuffer(parts[i]).String()] = true
	}
	instanceString := "10.10.10.1~~~~~~~~0~0"
	instanceFullString := "10.10.10.2~Mac~Dev~Hash~Dht~Keyfile~Key~TTL~1~0"
	instanceSemiString := "10.10.10.3~Mac~Dev~Hash~~~~~0~0"
	if !set[instanceString] || !set[instanceFullString] || !set[instanceSemiString] {
		t.Errorf("Failed to encode instances (5): encodedInstances incorrectly encoded the instanceList")
	}
}
func TestDecodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data := instanceList.encodeInstances()
	_, err := instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (1): %v", err)
	}
	data[len(data)-1] = 65
	_, err = instanceList.decodeInstances(data)
	if err == nil {
		t.Errorf("Failed to decode instances (2): must have returned non-nil but returned nil")
	}
	data = make([]byte, 0)
	_, err = instanceList.decodeInstances(data)
	if err == nil {
		t.Errorf("Failed to decode instances (3): must have returned non-nil but returned nil")
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = true
	P2Pinstance.Args.Port = 0
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (4): %v", err)
	}
	P2Pinstance = new(P2PInstance)
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (5): %v", err)
	}
	P2Pinstance = new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	instanceList = new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data = instanceList.encodeInstances()
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (6): %v", err)
	}
}

func TestSaveInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	_, err := instanceList.saveInstances("/")
	if err == nil {
		t.Errorf("Failed to save instances (1): must have returned non-nil but returned nil")
	}
	_, err = instanceList.saveInstances("save-4.save")
	defer os.Remove("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (2): %v", err)
	}
	file, err := os.Open("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (3): %v", err)
	}
	auxiliary := make([]byte, 100000)
	file.Read(auxiliary)
	file.Close()
	auxiliary = bytes.Trim(auxiliary, "\x00")
	t.Log(bytes.NewBuffer(auxiliary).String())
	P2PinstanceSecond := new(P2PInstance)
	P2PinstanceSecond.Args.IP = "10.10.10.2"
	P2PinstanceSecond.Args.Mac = "Mac"
	P2PinstanceSecond.Args.Dev = "Dev"
	P2PinstanceSecond.Args.Hash = "Hash"
	P2PinstanceSecond.Args.Dht = "Dht"
	P2PinstanceSecond.Args.Keyfile = "Keyfile"
	P2PinstanceSecond.Args.Key = "Key"
	P2PinstanceSecond.Args.TTL = "TTL"
	P2PinstanceSecond.Args.Fwd = false
	P2PinstanceSecond.Args.Port = 0
	instanceList.update("instanceSecond", P2PinstanceSecond)
	_, err = instanceList.saveInstances("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (4): %v", err)
	}
	file, err = os.Open("save-4.save")
	if err != nil {
		t.Errorf("Failed to save instances (5): %v", err)
	}
	auxiliary = make([]byte, 100000)
	file.Read(auxiliary)
	file.Close()
	auxiliary = bytes.Trim(auxiliary, "\x00")
	t.Log(bytes.NewBuffer(auxiliary).String())
}

func TestLoadInstances(t *testing.T) {
	if runtime.GOOS == "windows" {
		fmt.Println("This test is not supported on Windows")
	}
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "Dev"
	P2Pinstance.Args.Hash = "Hash"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	_, err := instanceList.loadInstances("/non-existing-file")
	if err == nil {
		t.Errorf("Failed to load instances (1): must have returned non-nil but returned nil")
	}
}

func TestInitialize(t *testing.T) {
	daemon := new(Daemon)
	daemon.Initialize("saveFile")
	if daemon.SaveFile != "saveFile" {
		t.Errorf("Failed to initialize (1): daemon couldn't initialize")
	}
}

func TestExecute(t *testing.T) {
	daemon := new(Daemon)
	args := new(Args)
	resp := new(Response)
	err := daemon.Execute(args, resp)
	if err != nil {
		t.Errorf("Failed to execute (1): %v", err)
	}
}
*/

func TestInstanceList_operate(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		action InstOperation
		id     string
		inst   *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if err := p.operate(tt.args.action, tt.args.id, tt.args.inst); (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.operate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceList_init(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			p.init()
		})
	}
}

func TestInstanceList_update(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		id   string
		inst *P2PInstance
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if err := p.update(tt.args.id, tt.args.inst); (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceList_delete(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if err := p.delete(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInstanceList_get(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*P2PInstance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if got := p.get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceList.get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceList_getInstance(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		id string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *P2PInstance
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if got := p.getInstance(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceList.getInstance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceList_encodeInstances(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			if got := p.encodeInstances(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceList.encodeInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceList_decodeInstances(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		data []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []RunArgs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			got, err := p.decodeInstances(tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.decodeInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceList.decodeInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceList_saveInstances(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			got, err := p.saveInstances(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.saveInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InstanceList.saveInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInstanceList_loadInstances(t *testing.T) {
	type fields struct {
		instances map[string]*P2PInstance
		lock      sync.RWMutex
	}
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []RunArgs
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &InstanceList{
				instances: tt.fields.instances,
				lock:      tt.fields.lock,
			}
			got, err := p.loadInstances(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("InstanceList.loadInstances() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InstanceList.loadInstances() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaemon_Initialize(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		saveFile string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			d.Initialize(tt.args.saveFile)
		})
	}
}

func TestDaemon_Execute(t *testing.T) {
	type fields struct {
		Instances  *InstanceList
		SaveFile   string
		OutboundIP net.IP
	}
	type args struct {
		args *Args
		resp *Response
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &Daemon{
				Instances:  tt.fields.Instances,
				SaveFile:   tt.fields.SaveFile,
				OutboundIP: tt.fields.OutboundIP,
			}
			if err := d.Execute(tt.args.args, tt.args.resp); (err != nil) != tt.wantErr {
				t.Errorf("Daemon.Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
