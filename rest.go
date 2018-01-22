package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	ptp "github.com/subutai-io/p2p/lib"
)

type request struct {
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
	Interfaces bool   `json:"interfaces"` // Used for show request
	All        bool   `json:"all"`        // Used for show request
}

type RESTResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ErrorOutput is a JSON object for output
type ErrorOutput struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func setupRESTHandlers(port int, d *Daemon) {
	http.HandleFunc("/rest/v1/start", d.execRESTStart)
	http.HandleFunc("/rest/v1/stop", d.execRESTStop)
	http.HandleFunc("/rest/v1/show", d.execRESTShow)
	http.HandleFunc("/rest/v1/status", d.execRESTStatus)
	http.HandleFunc("/rest/v1/debug", d.execRESTDebug)
	http.HandleFunc("/rest/v1/set", d.execRESTSet)

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func sendRequest(port int, command string, args *DaemonArgs) (*RESTResponse, error) {
	data, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/rest/v1/%s", port, command), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute command. Check if p2p daemon is running.")
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	out := &RESTResponse{}
	err = json.Unmarshal(body, out)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshal response: %s", err)
	}
	return out, nil
}

func sendRequestRaw(port int, command string, r *request) ([]byte, error) {
	data, err := json.Marshal(r)
	if err != nil {
		return nil, fmt.Errorf("Failed to marshal request: %s", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("http://localhost:%d/rest/v1/%s", port, command), bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Couldn't execute command. Check if p2p daemon is running.")
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func getJSON(body io.ReadCloser, args *DaemonArgs) error {
	data, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	//args := new(RunArgs)
	if len(data) == 0 {
		return nil
	}
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
