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
	MTU        bool   `json:"mtu"`
}

var bootstrap DHTConnection
var UsePMTU bool

// ExecDaemon starts P2P daemon
func ExecDaemon(port int, targetURL, sFile, profiling, syslog, logLevel string, mtu int, pmtu bool) {
	if logLevel == "" {
		ptp.SetMinLogLevelString(DefaultLog)
	} else {
		ptp.SetMinLogLevelString(logLevel)
	}
	if targetURL == "" {
		targetURL = "subutai.io"
	}
	if syslog != "" {
		ptp.SetSyslogSocket(syslog)
	}
	StartProfiling(profiling)
	ptp.InitPlatform()
	ptp.InitErrors()
	ptp.UsePMTU = pmtu

	if !ptp.CheckPermissions() {
		os.Exit(1)
	}
	StartTime = time.Now()

	ptp.GlobalMTU = mtu

	ReadyToServe = false

	err := bootstrap.init(targetURL)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to initialize bootstrap node connection")
		os.Exit(152)
	}
	go bootstrap.run()
	go waitOutboundIP()

	proc := new(Daemon)
	proc.Initialize(sFile)
	setupRESTHandlers(port, proc)

	go restoreInstances(proc)

	ReadyToServe = true

	SignalChannel = make(chan os.Signal, 1)
	signal.Notify(SignalChannel, os.Interrupt)

	go waitActiveBootstrap()

	go func() {
		for sig := range SignalChannel {
			fmt.Println("Received signal: ", sig)
			pprof.StopCPUProfile()
			os.Exit(0)
		}
	}()

	// main loop
	for {
		for id, inst := range proc.Instances.get() {
			if inst == nil || inst.PTP == nil {
				continue
			}
			if inst.PTP.ReadyToStop {
				err := proc.Stop(&DaemonArgs{Hash: id}, &Response{})
				if err != nil {
					ptp.Log(ptp.Error, "Failed to stop instance: %s", err)
				}
			}
		}
		time.Sleep(time.Millisecond * 100)
	}
	//select {}
}

func waitOutboundIP() {
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
}

func waitActiveBootstrap() {
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
}

func restoreInstances(daemon *Daemon) {
	for !bootstrap.isActive {
		time.Sleep(100 * time.Millisecond)
	}
	if daemon.SaveFile != "" {
		ptp.Log(ptp.Info, "Restore file provided")
		// Try to restore from provided file
		instances, err := daemon.Instances.loadInstances(daemon.SaveFile)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to load instances: %v", err)
		} else {
			ptp.Log(ptp.Info, "%d instances were loaded from file", len(instances))
			for _, inst := range instances {
				daemon.run(&inst, new(Response))
			}
		}
	}
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
