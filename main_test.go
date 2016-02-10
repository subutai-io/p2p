// +build ignore
package main

import (
	"os"
	"testing"
)

func TestStateRestore(t *testing.T) {
	Instances = make(map[string]Instance)
	var (
		i1 Instance
		i2 Instance
	)
	i1.Args.IP = "10.10.10.10"
	i1.Args.Dev = "vptp1"
	Instances["1"] = i1
	i2.Args.IP = "127.0.0.1"
	i2.Args.Dev = "vptp2"
	Instances["2"] = i2
	_, err := SaveInstances("t.file")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	loaded, _ := LoadInstances("t.file")
	if len(loaded) != 2 {
		t.Errorf("Resulting instances size doesn't match saved. Expecting 2, Received: %d", len(loaded))
	}
	if loaded[0].IP != "10.10.10.10" {
		t.Errorf("Loaded IP doesn't match saved")
	}
	if loaded[1].IP != "127.0.0.1" {
		t.Errorf("Loaded IP doesn't match saved")
	}
	if loaded[0].Dev != "vptp1" {
		t.Errorf("Loaded device name doesn't match saved")
	}
	if loaded[1].Dev != "vptp2" {
		t.Errorf("Loaded device name doesn't match saved")
	}
	os.Remove("t.file")
}
