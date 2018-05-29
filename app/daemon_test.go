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
