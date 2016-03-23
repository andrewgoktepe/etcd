package main

import (
	"fmt"
	"os"

	"gopkg.in/andrewgoktepe/etcd.v2/Godeps/_workspace/src/github.com/bgentry/speakeasy"
)

func main() {
	password, err := speakeasy.Ask("Please enter a password: ")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("Password result: %q\n", password)
	fmt.Printf("Password len: %d\n", len(password))
}
