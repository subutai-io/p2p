package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	ptp "github.com/subutai-io/p2p/lib"
)

// ShowOutput is a JSON object for output from `show` command
type ShowOutput struct {
	ID              string `json:"id"`
	IP              string `json:"ip"`
	Endpoint        string `json:"endpoint"`
	HardwareAddress string `json:"hw"`
	Error           string `json:"error"`
	Text            string `json:"text"`
	Code            int    `json:"code"`
	InterfaceName   string `json:"interface"`
	Hash            string `json:"hash"`
}

// Show outputs information about P2P instances and interfaces
func CommandShow(queryPort int, hash, ip string, interfaces, all, bind bool) {
	req := &request{}
	if hash != "" {
		req.Hash = hash
	} else {
		req.Hash = ""
	}
	req.IP = ip
	req.Interfaces = interfaces
	req.All = all
	req.Bind = bind

	out, err := sendRequestRaw(queryPort, "show", req)
	if err != nil {
		fmt.Println(err.Error())
		if err == errorFailedToMarshal {
			os.Exit(112)
		} else if err == errorFailedToCreatePOSTRequest {
			os.Exit(113)
		} else if err == errorFailedToExecuteRequest {
			os.Exit(115)
		}
		os.Exit(1)
	}
	show := []ShowOutput{}
	err = json.Unmarshal(out, &show)
	if err != nil {
		fmt.Printf("Failed to unmarshal JSON. Error %s\n", err)
		os.Exit(99)
	}

	if req.Hash != "" {
		if req.IP != "" {
			for _, m := range show {
				if m.Code != 0 {
					fmt.Println(m.Error)
				} else {
					fmt.Println(m.Text)
				}
				os.Exit(m.Code)
			}
			fmt.Println("No data available")
			os.Exit(102)
		} else {
			fmt.Println("< Peer ID >\t< IP >\t< Endpoint >\t< HW >")
			for _, m := range show {
				if m.Code != 0 {
					fmt.Println(m.Error)
					os.Exit(m.Code)
				}
				fmt.Printf("%s\t%s\t%s\t%s\n", m.ID, m.IP, m.Endpoint, m.HardwareAddress)
			}
			os.Exit(0)
		}
	}
	if req.Interfaces {
		if req.Bind {
			for _, m := range show {
				fmt.Printf("%s|%s\n", m.Hash, m.InterfaceName)
			}
			os.Exit(0)
		} else {
			for _, m := range show {
				if m.Code != 0 {
					fmt.Println(m.Error)
					os.Exit(m.Code)
				}
				fmt.Println(m.InterfaceName)
			}
			os.Exit(0)
		}
	}

	for _, m := range show {
		if m.Code != 0 {
			fmt.Println(m.Error)
			os.Exit(m.Code)
		}
		fmt.Printf("%s\t%s\t%s\n", m.HardwareAddress, m.IP, m.Hash)
	}
	os.Exit(0)
}

func (d *Daemon) execRESTShow(w http.ResponseWriter, r *http.Request) {
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	output, err := d.Show(&request{
		Hash:       args.Hash,
		IP:         args.IP,
		Interfaces: args.Interfaces,
		Bind:       args.Bind,
		All:        args.All,
	})
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(output)
}

// Show is used to output information about instances
func (d *Daemon) Show(args *request) ([]byte, error) {
	if !ReadyToServe {
		out := []ShowOutput{ShowOutput{Error: "P2P Daemon is in initialization mode. Can't handle request", Code: 105}}
		return d.showOutput(out)
	}
	if !bootstrap.isActive {
		out := []ShowOutput{ShowOutput{Error: "Not connected to DHT nodes", Code: 106}}
		return d.showOutput(out)
	}
	if bootstrap.ip == "" {
		out := []ShowOutput{ShowOutput{Error: "Didn't received outbound IP yet", Code: 107}}
		return d.showOutput(out)
	}

	if args.Hash != "" {
		inst := d.Instances.getInstance(args.Hash)
		if inst != nil {
			if args.IP != "" {
				out, err := d.showIP(args.IP, inst)
				return out, err
			}
			out, err := d.showHash(inst)
			return out, err
		}
		return d.showOutput([]ShowOutput{
			ShowOutput{
				Error: "Specified environment was not found",
				Code:  15,
			},
		})
	} else if args.Interfaces {
		if args.All {
			return d.showAllInterfaces()
		}
		if args.Bind {
			return d.showBindInterfaces()
		}
		return d.showInterfaces()
	} else {
		return d.showInstances()
	}
	return nil, nil
}

func (d *Daemon) showOutput(data []ShowOutput) ([]byte, error) {
	return json.Marshal(data)
}

func (d *Daemon) showIP(ip string, instance *P2PInstance) ([]byte, error) {
	peers := instance.PTP.Peers.Get()
	for _, peer := range peers {
		if peer.PeerLocalIP.String() == ip {
			if peer.State == ptp.PeerStateConnected {
				out := []ShowOutput{
					ShowOutput{
						Text: "Integrated with " + ip,
						Code: 0,
					},
				}
				return d.showOutput(out)
			}
		}
	}
	out := []ShowOutput{
		ShowOutput{
			Error: "Not yet integrated with " + ip,
			Code:  12,
		},
	}
	return d.showOutput(out)
}

func (d *Daemon) showHash(instance *P2PInstance) ([]byte, error) {
	peers := instance.PTP.Peers.Get()
	out := []ShowOutput{}
	for _, peer := range peers {
		s := ShowOutput{
			ID:              peer.ID,
			IP:              peer.PeerLocalIP.String(),
			Endpoint:        peer.Endpoint.String(),
			HardwareAddress: peer.PeerHW.String(),
		}
		out = append(out, s)
	}
	return d.showOutput(out)
}

func (d *Daemon) showInterfaces() ([]byte, error) {
	instances := d.Instances.get()
	out := []ShowOutput{}
	for _, inst := range instances {
		if inst.PTP != nil {
			s := ShowOutput{InterfaceName: inst.PTP.Interface.GetName()}
			out = append(out, s)
		}
	}
	return d.showOutput(out)
}

func (d *Daemon) showAllInterfaces() ([]byte, error) {
	out := []ShowOutput{}
	for _, inf := range InterfaceNames {
		s := ShowOutput{InterfaceName: inf}
		out = append(out, s)
	}
	return d.showOutput(out)
}

func (d *Daemon) showBindInterfaces() ([]byte, error) {
	instances := d.Instances.get()
	out := []ShowOutput{}
	for _, inst := range instances {
		if inst.PTP != nil && inst.PTP.Interface != nil {
			s := ShowOutput{Hash: inst.PTP.Dht.NetworkHash, InterfaceName: inst.PTP.Interface.GetName()}
			out = append(out, s)
		}
	}
	return d.showOutput(out)
}

func (d *Daemon) showInstances() ([]byte, error) {
	instances := d.Instances.get()
	out := []ShowOutput{}
	for key, inst := range instances {
		if inst.PTP != nil {
			s := ShowOutput{
				HardwareAddress: inst.PTP.Interface.GetHardwareAddress().String(),
				IP:              inst.PTP.Interface.GetIP().String(),
				Hash:            key,
			}
			out = append(out, s)
		} else {
			s := ShowOutput{
				HardwareAddress: "Unknown",
				IP:              "Unknown",
				Hash:            key,
			}
			out = append(out, s)
		}
	}
	return d.showOutput(out)
}
