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
	ptp.Log(ptp.Info, "Initializing P2P Daemon")
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

	bootstrapConnected := false
	bootstrapLastConnection := time.Unix(0, 0)

	for !bootstrapConnected {
		if time.Since(bootstrapLastConnection) > time.Duration(time.Second*5) {
			bootstrapLastConnection = time.Now()
			err := bootstrap.init(targetURL)
			if err == nil {
				bootstrapConnected = true
			} else {
				ptp.Log(ptp.Error, "Failed to connect to %s", targetURL)
			}
		}
		time.Sleep(time.Millisecond * 100)
	}

	go bootstrap.run()
	go waitOutboundIP()

	proc := new(Daemon)
	proc.init(sFile)
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
	if daemon.Restore != nil && daemon.Restore.isActive() {
		ptp.Log(ptp.Info, "Restore subsystem initialized")

		// loading from restore file
		err := daemon.Restore.load()
		if err != nil {
			ptp.Log(ptp.Error, "Failed to restore from file")
			return
		}

		entries := daemon.Restore.get()
		if len(entries) == 0 {
			return
		}

		ptp.Log(ptp.Info, "Attempt to restore %d instances", len(entries))

		restored := 0

		for _, e := range entries {
			err := daemon.run(&RunArgs{
				IP:      e.IP,
				Mac:     e.Mac,
				Dev:     e.Dev,
				Hash:    e.Hash,
				Keyfile: e.Keyfile,
				Key:     e.Key,
				TTL:     e.TTL,
			}, new(Response))
			if err != nil {
				ptp.Log(ptp.Error, "Failed to start instance %s during restore: %s", e.Hash, err.Error())
				continue
			} else {
				restored++
				daemon.Restore.bumpInstance(e.Hash)
			}
		}
		err = daemon.Restore.save()
		if err != nil {
			ptp.Log(ptp.Error, "Failed to save restore file")
		}
		ptp.Log(ptp.Info, "Restored %d of %d instances", restored, len(entries))
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
