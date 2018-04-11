package main

import (
	"bytes"
	"os"
	"testing"
)

func TestEncodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "vptp1"
	P2Pinstance.Args.Hash = "test"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.instances["hell"] = P2Pinstance
	_, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to encode instances: %v", err)
	}
}

func TestDecodingInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "vptp1"
	P2Pinstance.Args.Hash = "test"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.instances["hell"] = P2Pinstance
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to decode instances: %v", err)
	}
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to decode instances: %v", err)
	}
}

func TestSaveInstances(t *testing.T) {
	file, err := os.OpenFile("test", os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		t.Errorf("Failed to save instances: %v", err)
	}
	defer file.Close()
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "vptp1"
	P2Pinstance.Args.Hash = "test"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.instances["hell"] = P2Pinstance
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Failed to save instances: %v", err)
	}
	_, err = file.Write(data)
	if err != nil {
		t.Errorf("Failed to save instances: %v", err)
	}
}

func TestLoadInstances(t *testing.T) {
	P2Pinstance := new(P2PInstance)
	P2Pinstance.Args.IP = "10.10.10.1"
	P2Pinstance.Args.Mac = "Mac"
	P2Pinstance.Args.Dev = "vptp1"
	P2Pinstance.Args.Hash = "test"
	P2Pinstance.Args.Dht = "Dht"
	P2Pinstance.Args.Keyfile = "Keyfile"
	P2Pinstance.Args.Key = "Key"
	P2Pinstance.Args.TTL = "TTL"
	P2Pinstance.Args.Fwd = false
	P2Pinstance.Args.Port = 0
	instanceList := new(InstanceList)
	instanceList.init()
	instanceList.update("hell", P2Pinstance)
	t.Log(P2Pinstance)
	data, err := instanceList.encodeInstances()
	t.Log(data)
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	file, err := os.OpenFile("test", os.O_CREATE|os.O_RDWR, 0700)
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	_, err = file.Write(data)
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	file.Close()
	file, err = os.Open("test")
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	data = make([]byte, 100000)
	_, err = file.Read(data)
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	file.Close()
	data = bytes.Trim(data, "\x00") // TODO: add more security to this
	t.Log(data)
	t.Log(string(data))
	_, err = instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
}
