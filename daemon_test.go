package main

import (
	"testing"
)

func TestValidateDHT(t *testing.T) {
	if validateDHT("") != errEmptyDHTEndpoint {
		t.Fatalf("DHT Validate: providing empty list doesn't generate proper error")
	}
	if validateDHT("google.com") != errBadDHTEndpoint {
		t.Fatalf("Providing URL without port doesn't generate expected error")
	}
	if validateDHT("iamnotexist.atall:80") != errBadDHTEndpoint {
		t.Fatalf("Providing non existing URL doesn't generate expected error")
	}
	if validateDHT("google.com:80") != nil {
		t.Fatalf("Providing correct endpoint generates error")
	}
	if validateDHT("google.com:80,yandex.ru:80") != nil {
		t.Fatalf("Providing correct endpoints generates error")
	}
}
/*
Generated TestExecDaemon
Generated Test_validateDHT
package main

import (
	"testing"
)
*/

/*
func TestValidateDHT(t *testing.T) {
	if validateDHT("") != errEmptyDHTEndpoint {
		t.Fatalf("DHT Validate: providing empty list doesn't generate proper error")
	}
	if validateDHT("google.com") != errBadDHTEndpoint {
		t.Fatalf("Providing URL without port doesn't generate expected error")
	}
	if validateDHT("iamnotexist.atall:80") != errBadDHTEndpoint {
		t.Fatalf("Providing non existing URL doesn't generate expected error")
	}
	if validateDHT("google.com:80") != nil {
		t.Fatalf("Providing correct endpoint generates error")
	}
	if validateDHT("google.com:80,yandex.ru:80") != nil {
		t.Fatalf("Providing correct endpoints generates error")
	}
}
*/

func TestExecDaemon(t *testing.T) {
	type args struct {
		port      int
		dht       string
		sFile     string
		profiling string
		syslog    string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ExecDaemon(tt.args.port, tt.args.dht, tt.args.sFile, tt.args.profiling, tt.args.syslog)
		})
	}
}

func Test_validateDHT(t *testing.T) {
	type args struct {
		dht string
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
			if err := validateDHT(tt.args.dht); (err != nil) != tt.wantErr {
				t.Errorf("validateDHT() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
