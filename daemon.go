package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"
	"strings"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
)

var (
	errEmptyDHTEndpoint = errors.New("DHT endpoint wasn't specified")
	errBadDHTEndpoint   = errors.New("Endpoint have wrong format")
)

type DaemonArgs struct {
	IP         string `json:"ip"`
	Mac        string `json:"mac"`
	Dev        string `json:"dev"`
	Hash       string `json:"hash"`
	Dht        string `json:"dht"`
	Keyfile    string `json:"keyfile"`
	Key        string `json:"key"`
	TTL        string `json:"ttl"`
	Fwd        bool   `json:"fwd"`
	Port       int    `json:"port"`
	Interfaces bool   `json:"interfaces"` // show only
	All        bool   `json:"all"`        // show only
	Command    string `json:"command"`
	Args       string `json:"args"`
	Log        string `json:"log"`
	Bind       bool   `json:"bind"`
}

var bootstrap DHTConnection

// ExecDaemon starts P2P daemon
func ExecDaemon(port int, dht, sFile, profiling, syslog string) {
	if validateDHT(dht) != nil {
		os.Exit(213)
	}
	if syslog != "" {
		ptp.SetSyslogSocket(syslog)
	}
	StartProfiling(profiling)
	err := ptp.InitPlatform()
	if err != nil {
		ptp.Log(ptp.Error, "An error occurred while initializing the platform: %v", err)
		os.Exit(1717)
	}
	ptp.InitErrors()
	if DefaultLog == "TRACE" {
		ptp.SetMinLogLevel(ptp.Trace)
	} else if DefaultLog == "DEBUG" {
		ptp.SetMinLogLevel(ptp.Debug)
	} else if DefaultLog == "INFO" {
		ptp.SetMinLogLevel(ptp.Info)
	} else if DefaultLog == "WARNING" {
		ptp.SetMinLogLevel(ptp.Warning)
	} else if DefaultLog == "ERROR" {
		ptp.SetMinLogLevel(ptp.Error)
	}

	if !ptp.CheckPermissions() {
		os.Exit(1)
	}
	StartTime = time.Now()

	ReadyToServe = false

	err := bootstrap.init(dht)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to initialize bootstrap node connection")
		os.Exit(152)
	}
	go bootstrap.run()
	go func() {
		for _, r := range bootstrap.routers {
			if r != nil {
				go r.run()
				go r.keepAlive()
			}
		}
		for !bootstrap.isActive {
			for _, r := range bootstrap.routers {
				if r.running && r.handshaked {
					bootstrap.isActive = true
					break
				}
			}
			time.Sleep(time.Millisecond * 100)
		}
		for bootstrap.ip == "" {
			time.Sleep(time.Millisecond * 100)
		}
		OutboundIP = net.ParseIP(bootstrap.ip)
	}()

	proc := new(Daemon)
	proc.Initialize(sFile)
	setupRESTHandlers(port, proc)

	go func() {
		for !bootstrap.isActive {
			time.Sleep(100 * time.Millisecond)
		}
		if sFile != "" {
			ptp.Log(ptp.Info, "Restore file provided")
			// Try to restore from provided file
			instances, err := proc.Instances.loadInstances(proc.SaveFile)
			if err != nil {
				ptp.Log(ptp.Error, "Failed to load instances: %v", err)
			} else {
				ptp.Log(ptp.Info, "%d instances were loaded from file", len(instances))
				for _, inst := range instances {
					proc.run(&inst, new(Response))
				}
			}
		}
	}()

	ReadyToServe = true

	SignalChannel = make(chan os.Signal, 1)
	signal.Notify(SignalChannel, os.Interrupt)

	go func() {
		for {
			active := 0
			for _, r := range bootstrap.routers {
				if !r.stop {
					active++
				}
			}
			if active == 0 {
				ptp.Log(ptp.Info, "No active bootstrap nodes")
				os.Exit(0)
			}
			time.Sleep(time.Millisecond * 100)
		}
	}()

	go func() {
		for sig := range SignalChannel {
			fmt.Println("Received signal: ", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()
	select {}
}

func validateDHT(dht string) error {
	if dht == "" {
		ptp.Log(ptp.Error, "Empty bootstrap list")
		return errEmptyDHTEndpoint
	}
	eps := strings.Split(dht, ",")
	for _, ep := range eps {
		_, err := net.ResolveTCPAddr("tcp4", ep)
		if err != nil {
			ptp.Log(ptp.Error, "Bootstrap %s have bad format or wrong address: %s", ep, err)
			return errBadDHTEndpoint
		}
	}
	return nil
}
