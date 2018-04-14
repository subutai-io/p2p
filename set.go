package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	ptp "github.com/subutai-io/p2p/lib"
)

// Set modifies different options of P2P daemon
func CommandSet(rpcPort int, log, hash, keyfile, key, ttl string) {
	out, err := sendRequest(rpcPort, "set", &DaemonArgs{Log: log, Keyfile: keyfile, Key: key, TTL: ttl})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
}

func (d *Daemon) execRESTSet(w http.ResponseWriter, r *http.Request) {
	if !ReadyToServe {
		resp, _ := getResponse(105, "P2P Daemon is in initialization state")
		w.Write(resp)
		return
	}
	if !bootstrap.isActive {
		resp, _ := getResponse(106, "Not connected to DHT nodes")
		w.Write(resp)
		return
	}
	if bootstrap.ip == "" {
		resp, _ := getResponse(107, "Didn't received outbound IP yet")
		w.Write(resp)
		return
	}
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	response := new(Response)
	if args.Log != "" {
		d.SetLog(&NameValueArg{
			Name:  "log",
			Value: args.Log,
		}, response)
	} else {
		response.ExitCode = 0
		response.Output = "Unknown command"
	}
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	ptp.Log(ptp.Info, "RESPONSE: %s", string(resp))
	w.Write(resp)
}

// SetLog modifies specific option
func (d *Daemon) SetLog(args *NameValueArg, resp *Response) error {
	args.Value = strings.ToLower(args.Value)
	ptp.Log(ptp.Info, "Setting option %s to %s", args.Name, args.Value)
	resp.ExitCode = 0
	if args.Name == "log" {
		resp.Output = "Logging level has switched to " + args.Value + " level"
		if args.Value == "debug" {
			ptp.SetMinLogLevel(ptp.Debug)
		} else if args.Value == "info" {
			ptp.SetMinLogLevel(ptp.Info)
		} else if args.Value == "trace" {
			ptp.SetMinLogLevel(ptp.Trace)
		} else if args.Value == "warning" {
			ptp.SetMinLogLevel(ptp.Warning)
		} else if args.Value == "error" {
			ptp.SetMinLogLevel(ptp.Error)
		} else {
			resp.ExitCode = 1
			resp.Output = "Unknown log level was specified. Supported log levels is:\n"
			resp.Output = resp.Output + "TRACE\n"
			resp.Output = resp.Output + "DEBUG\n"
			resp.Output = resp.Output + "INFO\n"
			resp.Output = resp.Output + "WARNING\n"
			resp.Output = resp.Output + "ERROR\n"
		}
	}
	return nil
}

// AddKey adds a new crypto-key
func (p *Daemon) AddKey(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	if args.Hash == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified hash"
	}
	if args.Key == "" {
		resp.ExitCode = 1
		resp.Output = "You have not specified key"
	}
	inst := p.Instances.getInstance(args.Hash)
	if inst == nil {
		resp.ExitCode = 1
		resp.Output = "No instances with specified hash were found"
	}
	if resp.ExitCode == 0 {
		resp.Output = "New key added"
		var newKey ptp.CryptoKey

		newKey = inst.PTP.Crypter.EnrichKeyValues(newKey, args.Key, args.TTL)
		inst.PTP.Crypter.Keys = append(inst.PTP.Crypter.Keys, newKey)
		p.Instances.update(args.Hash, inst)
	}
	return nil
}
