package main

import (
	"fmt"
)

func UsageDaemon() {
	fmt.Printf("p2p network running on the same machine are controlled by daemon. \n" +
		"Daemon should be started by privileged user, because process will attempt to \n" +
		"create new virtual network interfaces (tap). \n" +
		"When running p2p in daemon mode it will listen to a particular port specified by optional -port \n" +
		"argument (Default: 52523) for local RPC connection and wait for commands from p2p client (same \n" +
		"application, but without daemon command)\n\n")
	fmt.Printf("Usage: p2p daemon [OPTIONS]:\n")
}

func UsageStart() {
	fmt.Printf("start command allows user to run new p2p instance. This command executes start procedure in a daemon.\n\n")
	fmt.Printf("Usage: p2p start [-ip IP] [-hash HASH] [OPTIONS]:\n")
}

func UsageStop() {
	fmt.Printf("Usage: p2p stop -hash HASH:\n")
}

func UsageShow() {
	fmt.Printf("Usage: p2p show [-hash HASH]:\n")
}

func UsageSet() {
	fmt.Printf("Usage: p2p set [OPTIONS]:\n")
}
