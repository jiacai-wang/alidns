package main

import (
	"fmt"
	"net"
	"testing"
)

func TestGtPubIp(t *testing.T) {
	ip := getPubIp()
	if net.ParseIP(ip) == nil {
		fmt.Printf("IP Address: %s - Invalid\n", ip)
		t.Fail()
	} else {
		fmt.Printf("IP Address: %s - Valid\n", ip)
	}

}
