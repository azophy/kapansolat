package main

import (
	"testing"
)

const (
	ipKemenag = "103.7.13.247" // kemenag.go.id
)

// func TestGetIpInfoLocalhost(t *testing.T) {
// 	_, err := getIpInfo("127.0.0.1")
// 	if err == nil {
// 		t.Error("expecting error, got success")
// 	}
// }

func TestGetIpInfoKemenag(t *testing.T) {
	res, _ := getIpInfo(ipKemenag)
	if res.Country != "Indonesia" {
		t.Errorf("expecting country Indonesia, got %v", res.Country)
	}
}
