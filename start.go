package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	ptp "github.com/subutai-io/p2p/lib"
)

// CommandStart will create new P2P instance
func CommandStart(restPort int, ip, hash, mac, dev, keyfile, key, ttl string, fwd bool, port int) {
	args := &DaemonArgs{}
	args.IP = ip
	if hash == "" {
		fmt.Printf("Hash cannot be empty. Please start new instances with -hash VALUE argument\n")
		os.Exit(12)
	}
	if strings.Index(hash, "~") != -1 {
		fmt.Printf("Hash cannot contain the ~. Please start new instances with hash value that doesn't contain it.\n")
		os.Exit(17)
	}
	args.Hash = hash
	if mac != "" {
		_, err := net.ParseMAC(mac)
		if err != nil {
			fmt.Printf("Invalid MAC address provided\n")
			os.Exit(13)
		}
	}
	args.Mac = mac
	args.Dev = dev
	args.Keyfile = keyfile
	args.Key = key
	args.TTL = ttl
	args.Fwd = fwd
	args.Port = port

	out, err := sendRequest(restPort, "start", args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
}

func (d *Daemon) execRESTStart(w http.ResponseWriter, r *http.Request) {
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
	d.run(&RunArgs{
		IP:      args.IP,
		Mac:     args.Mac,
		Dev:     args.Dev,
		Hash:    args.Hash,
		Dht:     args.Dht,
		Keyfile: args.Keyfile,
		Key:     args.Key,
		TTL:     args.TTL,
		Fwd:     args.Fwd,
		Port:    args.Port,
	}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

// Run starts a P2P instance
func (d *Daemon) run(args *RunArgs, resp *Response) error {
	resp.ExitCode = 0
	resp.Output = "Running new P2P instance for " + args.Hash + "\n"

	ptp.Log(ptp.Trace, "Requested new P2P instance: %+v", args)

	// Validate if interface name is unique
	if args.Dev != "" {
		instances := d.Instances.get()
		for _, inst := range instances {
			if inst.PTP.Interface.GetName() == args.Dev {
				resp.ExitCode = 1
				resp.Output = "Device with specified name is already in use"
				return errors.New(resp.Output)
			}
		}
	}

	inst := d.Instances.getInstance(args.Hash)
	if inst == nil {
		resp.Output = resp.Output + "Lookup finished\n"
		if args.Key != "" {
			if len(args.Key) < 16 {
				args.Key += "0000000000000000"[:16-len(args.Key)]
			} else if len(args.Key) > 16 && len(args.Key) < 24 {
				args.Key += "000000000000000000000000"[:24-len(args.Key)]
			} else if len(args.Key) > 24 && len(args.Key) < 32 {
				args.Key += "00000000000000000000000000000000"[:32-len(args.Key)]
			} else if len(args.Key) > 32 {
				args.Key = args.Key[:32]
			}
		}

		newInst := new(P2PInstance)
		newInst.ID = args.Hash
		newInst.Args = *args
		newInst.PTP = ptp.New(args.Mac, args.Hash, args.Keyfile, args.Key, args.TTL, TargetURL, args.Fwd, args.Port, OutboundIP)
		if newInst.PTP == nil {
			resp.Output = resp.Output + "Failed to create P2P Instance"
			resp.ExitCode = 1
			return errors.New("Failed to create P2P Instance")
		}

		err := bootstrap.registerInstance(newInst.ID, newInst)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to register instance with bootstrap nodes: %s", err.Error())
			if newInst.PTP != nil {
				newInst.PTP.Close()
				newInst.PTP = nil
			}
			resp.Output = resp.Output + "Failed to register instance: %s" + err.Error()
			resp.ExitCode = 601
			return errors.New("Failed to register instance")
		}

		go newInst.PTP.ReadDHT()
		newInst.PTP.Dht.LocalPort = newInst.PTP.UDPSocket.GetPort()
		newInst.PTP.FindNetworkAddresses()
		err = newInst.PTP.Dht.Connect(newInst.PTP.LocalIPs, newInst.PTP.ProxyManager.GetList())
		if err != nil {
			if newInst.PTP != nil {
				newInst.PTP.Close()
				newInst.PTP = nil
			}
			bootstrap.unregisterInstance(newInst.ID)
			resp.Output = resp.Output + err.Error()
			resp.ExitCode = 602
			return err
		}

		err = newInst.PTP.PrepareInterfaces(args.IP, args.Dev)
		if err != nil {
			ptp.Log(ptp.Error, "Failed to configure network interface: %s", err)
			if newInst.PTP != nil {
				newInst.PTP.Close()
				newInst.PTP = nil
			}
			bootstrap.unregisterInstance(newInst.ID)
			resp.Output = resp.Output + "Failed to configure network: " + err.Error()
			resp.ExitCode = 603
			return errors.New("Failed to configure network interface")
		}
		go newInst.PTP.ListenInterface()

		// Saving interface name
		infFound := false
		for _, inf := range InterfaceNames {
			if inf == newInst.PTP.Interface.GetName() {
				infFound = true
			}
		}
		if !infFound && newInst.PTP.Interface.GetName() != "" {
			InterfaceNames = append(InterfaceNames, newInst.PTP.Interface.GetName())
		}

		usedIPs = append(usedIPs, newInst.PTP.Interface.GetIP().String())
		ptp.Log(ptp.Info, "Instance created")
		d.Instances.update(args.Hash, newInst)

		go newInst.PTP.Run()
		if d.SaveFile != "" {
			resp.Output = resp.Output + "Saving instance into file"
			d.Instances.saveInstances(d.SaveFile)
		}
	} else {
		resp.Output = resp.Output + "Hash already in use\n"
	}
	return nil
}
