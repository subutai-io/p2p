/*
Generated Test_setupRESTHandlers
Generated Test_sendRequest
Generated Test_sendRequestRaw
Generated Test_getJSON
Generated Test_getResponse
Generated Test_handleMarshalError
*/
package main

import (
	"io"
	"net/http"
	"reflect"
	"testing"
)

func Test_setupRESTHandlers(t *testing.T) {
	type args struct {
		port int
		d    *Daemon
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setupRESTHandlers(tt.args.port, tt.args.d)
		})
	}
}

func Test_sendRequest(t *testing.T) {
	type args struct {
		port    int
		command string
		args    *DaemonArgs
	}
	tests := []struct {
		name    string
		args    args
		want    *RESTResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sendRequest(tt.args.port, tt.args.command, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sendRequestRaw(t *testing.T) {
	type args struct {
		port    int
		command string
		r       *request
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := sendRequestRaw(tt.args.port, tt.args.command, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("sendRequestRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sendRequestRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getJSON(t *testing.T) {
	type args struct {
		body io.ReadCloser
		args *DaemonArgs
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := getJSON(tt.args.body, tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("getJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getResponse(t *testing.T) {
	type args struct {
		exitCode      int
		outputMessage string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getResponse(tt.args.exitCode, tt.args.outputMessage)
			if (err != nil) != tt.wantErr {
				t.Errorf("getResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleMarshalError(t *testing.T) {
	type args struct {
		err error
		w   http.ResponseWriter
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := handleMarshalError(tt.args.err, tt.args.w); (err != nil) != tt.wantErr {
				t.Errorf("handleMarshalError() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
