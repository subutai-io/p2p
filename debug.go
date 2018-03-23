package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	ptp "github.com/subutai-io/p2p/lib"
)

// Debug prints debug information
func CommandDebug(restPort int) {
	out, err := sendRequest(restPort, "debug", &DaemonArgs{})
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println(out.Message)
	os.Exit(out.Code)
}

func (d *Daemon) execRESTDebug(w http.ResponseWriter, r *http.Request) {
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
	d.Debug(&Args{
		Command: args.Command,
		Args:    args.Args,
	}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

// Debug output debug information
func (p *Daemon) Debug(args *Args, resp *Response) error {
	resp.Output = fmt.Sprintf("Version: %s Build: %s\n", AppVersion, BuildID)
	resp.Output += fmt.Sprintf("Uptime: %d h %d m %d s\n", int(time.Since(StartTime).Hours()), int(time.Since(StartTime).Minutes())%60, int(time.Since(StartTime).Seconds())%60)
	resp.Output += fmt.Sprintf("Number of gouroutines: %d\n", runtime.NumGoroutine())
	resp.Output += fmt.Sprintf("Bootstrap nodes information:\n")
	for _, node := range bootstrap.routers {
		if node != nil {
			resp.Output += fmt.Sprintf("  %s Rx: %d Tx: %d\n", node.addr.String(), node.rx, node.tx)
		}
	}
	resp.Output += fmt.Sprintf("Instances information:\n")
	instances := p.Instances.Get()
	for _, inst := range instances {
		resp.Output += fmt.Sprintf("Bootstrap nodes:\n")
		for _, conn := range inst.PTP.Dht.Connections {
			resp.Output += fmt.Sprintf("\t%s\n", conn.RemoteAddr().String())
		}
		resp.Output += fmt.Sprintf("Hash: %s\n", inst.ID)
		resp.Output += fmt.Sprintf("ID: %s\n", inst.PTP.Dht.ID)
		resp.Output += fmt.Sprintf("UDP Port: %d\n", inst.PTP.UDPSocket.GetPort())
		resp.Output += fmt.Sprintf("Network interfaces: ")
		for _, ip := range inst.PTP.LocalIPs {
			resp.Output += fmt.Sprintf("%s ", ip.String())
		}
		resp.Output += "\n"
		resp.Output += fmt.Sprintf("P2P Interface %s, HW Addr: %s, IP: %s\n", inst.PTP.Interface.GetName(), inst.PTP.Interface.GetHardwareAddress().String(), inst.PTP.Interface.GetIP().String())
		resp.Output += fmt.Sprintf("Proxies: ")
		proxyList := inst.PTP.ProxyManager.GetList()
		if len(proxyList) == 0 {
			resp.Output += fmt.Sprintf("No proxies in use")
		}
		for _, proxy := range proxyList {
			resp.Output += fmt.Sprintf("%s/%d ", proxy.Addr.String(), proxy.Endpoint.Port)
		}
		resp.Output += "\n"
		resp.Output += fmt.Sprintf("Peers:\n")

		peers := inst.PTP.Peers.Get()
		for _, peer := range peers {
			resp.Output += fmt.Sprintf("\t--- %s ---\n", peer.ID)
			resp.Output += fmt.Sprintf("\tStates: %s | %s\n", ptp.StringifyState(peer.State), ptp.StringifyState(peer.RemoteState))
			if peer.PeerLocalIP == nil {
				resp.Output += "\tNo IP assigned\n"
			} else if peer.PeerHW == nil {
				resp.Output += "\tNo MAC assigned\n"
			} else {
				resp.Output += fmt.Sprintf("\tNetwork: %s %s\n", peer.PeerLocalIP.String(), peer.PeerHW.String())
				resp.Output += fmt.Sprintf("\tEndpoint: %s\n", peer.Endpoint)
				resp.Output += fmt.Sprintf("\tAll Endpoints: ")
				for _, ep := range peer.EndpointsActive {
					resp.Output += fmt.Sprintf("%s ", ep.Addr.String())
				}
				resp.Output += "\n"
			}
			resp.Output += fmt.Sprintf("\tEndpoints pool: ")
			pool := []*net.UDPAddr{}
			pool = append(pool, peer.KnownIPs...)
			pool = append(pool, peer.Proxies...)
			for _, v := range pool {
				resp.Output += fmt.Sprintf("%s ", v.String())
			}
			resp.Output += "\n"
			resp.Output += fmt.Sprintf("\t--- End of %s ---\n", peer.ID)
		}
	}
	return nil
}
