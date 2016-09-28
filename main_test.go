package main

import (
	"os"
	"testing"
)

func TestStateRestore(t *testing.T) {
	instances = make(map[string]instance)
	var (
		i1 instance
		i2 instance
	)
	i1.Args.IP = "10.10.10.10"
	i1.Args.Dev = "vptp1"
	instances["1"] = i1
	i2.Args.IP = "127.0.0.1"
	i2.Args.Dev = "vptp2"
	instances["2"] = i2
	_, err := saveInstances("t.file")
	if err != nil {
		t.Errorf("%v", err)
	}

	loaded, err := loadInstances("t.file")
	if err != nil {
		t.Errorf("Failed to load instances: %v", err)
	}
	if len(loaded) != 2 {
		t.Errorf("Resulting instances size doesn't match saved. Expecting 2, Received: %d", len(loaded))
	}
	if loaded[0].IP != "10.10.10.10" && loaded[0].IP != "127.0.0.1" {
		t.Errorf("Loaded IP doesn't match saved: %s", loaded[0].IP)
	}
	if loaded[1].IP != "127.0.0.1" && loaded[1].IP != "10.10.10.10" {
		t.Errorf("Loaded IP doesn't match saved: %s", loaded[1].IP)
	}
	if loaded[0].Dev != "vptp1" && loaded[0].Dev != "vptp2" {
		t.Errorf("Loaded device name doesn't match saved: %s", loaded[0].Dev)
	}
	if loaded[1].Dev != "vptp2" && loaded[1].Dev != "vptp1" {
		t.Errorf("Loaded device name doesn't match saved: %s", loaded[1].Dev)
	}
	os.Remove("t.file")
}
