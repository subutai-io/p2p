package main

import (
	"fmt"
	"net/http"
	"os"

	ptp "github.com/subutai-io/p2p/lib"
)

// CommandStop will terminate P2P instance
// Function will send a request to the /stop/ REST endpoint with
// specified hash that is needed to stop or interface name that's
// needed to be removed from saved interfaces list
func CommandStop(rpcPort int, hash, dev string) {
	args := &DaemonArgs{}
	if hash != "" {
		args.Hash = hash
		args.Dev = ""
	} else if dev != "" {
		args.Dev = dev
		args.Hash = ""
	} else {
		fmt.Printf("Not enough parameters for stop command")
		return
	}
	out, err := sendRequest(rpcPort, "stop", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
}

func (d *Daemon) execRESTStop(w http.ResponseWriter, r *http.Request) {
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
	d.Stop(&DaemonArgs{
		Hash: args.Hash,
		Dev:  args.Dev,
	}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

// Stop is used to terminate a specific P2P instance
func (p *Daemon) Stop(args *DaemonArgs, resp *Response) error {
	resp.ExitCode = 0
	if args.Hash != "" {
		inst := p.Instances.GetInstance(args.Hash)
		if inst == nil {
			resp.ExitCode = 1
			resp.Output = "Instance with hash " + args.Hash + " was not found"
			return nil
		} else {
			ip := inst.PTP.Interface.GetIP().String()
			resp.Output = "Shutting down " + args.Hash
			inst.PTP.Close()
			p.Instances.Delete(args.Hash)
			p.Instances.SaveInstances(p.SaveFile)
			k := 0
			for k, i := range usedIPs {
				if i != ip {
					usedIPs[k] = i
					k++
				}
			}
			usedIPs = usedIPs[:k]
			bootstrap.unregisterInstance(args.Hash)
			return nil
		}
	} else if args.Dev != "" {
		resp.Output = "Removing " + args.Dev
		instances := p.Instances.Get()
		for i, inf := range InterfaceNames {
			if inf == args.Dev {
				for _, instance := range instances {
					if instance.PTP.Interface.GetName() == args.Dev {
						resp.ExitCode = 12
						resp.Output += "Can't remove interface: In use"
						return nil
					}
				}
				InterfaceNames = append(InterfaceNames[:i], InterfaceNames[i+1:]...)
				resp.ExitCode = 0
				resp.Output += "Removed " + args.Dev
				return nil
			}
		}
		resp.ExitCode = 1
		resp.Output += "Interface was not found"
		return nil
	}
	resp.ExitCode = 2
	resp.Output = "Not enough parameters for stop"
	return nil
}
