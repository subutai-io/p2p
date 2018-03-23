package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	ptp "github.com/subutai-io/p2p/lib"
)

type statusResponse struct {
	Instances []*statusInstance `json:"instances"`
	Code      int               `json:"code"`
}

type statusInstance struct {
	ID    string        `json:"id"`
	IP    string        `json:"ip"`
	Peers []*statusPeer `json:"peers"`
}

type statusPeer struct {
	ID        string `json:"id"`
	IP        string `json:"ip"`
	State     string `json:"state"`
	LastError string `json:"lastError"`
}

// CommandStatus outputs connectivity status of each peer
func CommandStatus(restPort int) {
	out, err := sendRequestRaw(restPort, "status", &request{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	response := new(statusResponse)
	err = json.Unmarshal(out, response)
	if err != nil {
		fmt.Printf("Failed to unmarshal status response: %s", err)
		os.Exit(125)
	}

	if response.Code != 0 {
		os.Exit(response.Code)
	}

	for _, instance := range response.Instances {
		fmt.Printf("%s|%s\n", instance.ID, instance.IP)
		for _, peer := range instance.Peers {
			fmt.Printf("%s|%s|State:%s|", peer.ID, peer.IP, peer.State)
			if peer.LastError != "" {
				fmt.Printf("LastError:%s", peer.LastError)
			}
			fmt.Printf("\n")
		}
	}
	os.Exit(0)
}

func (d *Daemon) execRESTStatus(w http.ResponseWriter, r *http.Request) {
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
	response, err := d.Status()
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	output, err := json.Marshal(response)
	if err != nil {
		ptp.Log(ptp.Error, "Failed to marshal status response: %s", err)
		return
	}
	w.Write(output)
}

// Status displays information about instances, peers and their statuses
func (d *Daemon) Status() (*statusResponse, error) {
	response := &statusResponse{}
	if !ReadyToServe {
		response.Code = 105
		return response, nil
	}
	response.Instances = []*statusInstance{}
	instances := d.Instances.Get()
	for _, inst := range instances {
		instance := &statusInstance{
			ID: inst.ID,
			IP: inst.PTP.Interface.GetIP().String(),
		}
		peers := inst.PTP.Peers.Get()
		for _, peer := range peers {
			instance.Peers = append(instance.Peers, &statusPeer{
				ID:        peer.ID,
				IP:        peer.PeerLocalIP.String(),
				State:     ptp.StringifyState(peer.State),
				LastError: peer.LastError,
			})
		}
		response.Instances = append(response.Instances, instance)
	}
	return response, nil
}
