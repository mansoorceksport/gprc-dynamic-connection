package main

import (
	"fmt"
	"testing"
)

func Test_hashAddresses(t *testing.T) {
	ipAddress := []string{"10.244.101.126", "10.244.101.125", "10.244.101.123"}

	s := hashAddresses(ipAddress)
	fmt.Println(s)
	fmt.Println(ipAddress)
}
