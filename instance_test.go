package main

import (
	"bytes"
	"fmt"
	"runtime"
	"testing"
	"os"
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
