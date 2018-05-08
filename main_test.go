package main

import (
	"os"
	"testing"
)

func TestStateRestore(t *testing.T) {
	daemon := new(Daemon)
	daemon.Initialize("t.file")

	i1 := new(P2PInstance)
	i2 := new(P2PInstance)
	i1.Args.IP = "10.10.10.10"
	i1.Args.Dev = "vptp1"
	daemon.Instances.update("1", i1)
	i2.Args.IP = "127.0.0.1"
	i2.Args.Dev = "vptp2"
	daemon.Instances.update("2", i2)

	_, err := daemon.Instances.saveInstances("t.file")
	if err != nil {
		t.Errorf("%v", err)
	}

	loaded, err := daemon.Instances.loadInstances("t.file")
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
/*
Generated TestStartProfiling
Generated Test_main
package main

import (
	"os"
	"testing"
)
*/

/*
func TestStateRestore(t *testing.T) {
	daemon := new(Daemon)
	daemon.Initialize("t.file")

	i1 := new(P2PInstance)
	i2 := new(P2PInstance)
	i1.Args.IP = "10.10.10.10"
	i1.Args.Dev = "vptp1"
	daemon.Instances.update("1", i1)
	i2.Args.IP = "127.0.0.1"
	i2.Args.Dev = "vptp2"
	daemon.Instances.update("2", i2)

	_, err := daemon.Instances.saveInstances("t.file")
	if err != nil {
		t.Errorf("%v", err)
	}

	loaded, err := daemon.Instances.loadInstances("t.file")
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
*/

func TestStartProfiling(t *testing.T) {
	type args struct {
		profile string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StartProfiling(tt.args.profile)
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}
