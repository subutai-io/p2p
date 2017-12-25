package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	ptp "github.com/subutai-io/p2p/lib"
)

type RESTResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

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
}

func (d *Daemon) execRESTStart(w http.ResponseWriter, r *http.Request) {
	ptp.Log(ptp.Info, "Start request")
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	response := new(Response)
	d.Run(&RunArgs{
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

func (d *Daemon) execRESTStop(w http.ResponseWriter, r *http.Request) {
	ptp.Log(ptp.Info, "Stop request")
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	response := new(Response)
	d.Stop(&StopArgs{
		Hash: args.Hash,
	}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

func (d *Daemon) execRESTShow(w http.ResponseWriter, r *http.Request) {
	ptp.Log(ptp.Info, "Show request")
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	response := new(Response)
	d.Show(&ShowArgs{
		Hash:       args.Hash,
		IP:         args.IP,
		Interfaces: args.Interfaces,
		All:        args.All,
	}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

func (d *Daemon) execRESTStatus(w http.ResponseWriter, r *http.Request) {
	ptp.Log(ptp.Info, "Status request")
	args := new(DaemonArgs)
	err := getJSON(r.Body, args)
	if handleMarshalError(err, w) != nil {
		return
	}
	response := new(Response)
	d.Status(&RunArgs{}, response)
	resp, err := getResponse(response.ExitCode, response.Output)
	if err != nil {
		ptp.Log(ptp.Error, "Internal error: %s", err)
		return
	}
	w.Write(resp)
}

func (d *Daemon) execRESTDebug(w http.ResponseWriter, r *http.Request) {
	ptp.Log(ptp.Info, "Debug request")
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
	ptp.Log(ptp.Info, "RESPONSE: %s", string(resp))
	w.Write(resp)
}

func (d *Daemon) execRESTSet(w http.ResponseWriter, r *http.Request) {

}

func getJSON(body io.ReadCloser, args *DaemonArgs) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	//args := new(RunArgs)
	err = json.Unmarshal(data, args)
	if err != nil {
		return err
	}
	return nil
}

func getResponse(exitCode int, outputMessage string) ([]byte, error) {
	resp := &RESTResponse{
		Code:    exitCode,
		Message: outputMessage,
	}
	out, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal response: %s", err)
	}
	return out, nil
}

func handleMarshalError(err error, w http.ResponseWriter) error {
	if err != nil {
		errText := fmt.Sprintf("Failed to read request body: %s", err)
		ptp.Log(ptp.Error, "%s", errText)
		resp, err := getResponse(1, errText)
		if err != nil {
			ptp.Log(ptp.Error, "Internal error: %s", err)
			return err
		}
		w.Write(resp)
		return err
	}
	return nil
}
