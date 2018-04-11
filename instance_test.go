package main

import (
	"fmt"
	"testing"
)

func testEncodingAndDecoding(t *testing.T) {
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
	instanceList.instances["hell"] = P2Pinstance
	data, err := instanceList.encodeInstances()
	if err != nil {
		t.Errorf("Encode failed: %v", err)
	}
	fmt.Println(data)
	args, err := instanceList.decodeInstances(data)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
	fmt.Println(args)
}
