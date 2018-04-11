package main

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"testing"
)

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
		t.Errorf("Failed to get (1/1): get returned unexpected map")
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
		t.Errorf("Failed to get instance (1/2): getInstance returned nil, but instance exists")
	}
	if instanceList.getInstance("non-instance") != nil {
		t.Errorf("Failed to get instance (2/2): getInstance returned an instance, but instance does not exist")
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
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	_, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to encode instances (1/1): %v", err)
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
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("instance", P2Pinstance)
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to decode instances (1/2): %v", err)
	}
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances (2/2): %v", err)
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
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to save instances (1/3): %v", err)
	}
	file, err := os.OpenFile("test", os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		t.Errorf("Failed to save instances (2/3): %v", err)
	}
	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		t.Errorf("Failed to save instances (3/3): %v", err)
	}
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
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to load instances (1/6): %v", err)
	}
	file, err := os.OpenFile("/tmp/test", os.O_CREATE|os.O_RDWR, 0700)
	defer func() {
		os.Remove("/tmp/test")
	}()
	if err != nil {
		t.Errorf("Failed to load instances (2/6): %v", err)
	}
	_, err = file.Write(data)
	if err != nil {
		t.Errorf("Failed to load instances (3/6): %v", err)
	}
	file.Close()
	file, err = os.Open("/tmp/test")
	if err != nil {
		t.Errorf("Failed to load instances (4/6): %v", err)
	}
	data = make([]byte, 100000)
	_, err = file.Read(data)
	if err != nil {
		t.Errorf("Failed to load instances (5/6): %v", err)
	}
	file.Close()
	data = bytes.Trim(data, "\x00")
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to load instances (6/6): %v", err)
	}
}
