package main

import gexpect "gopkg.in/andrewgoktepe/etcd.v2/Godeps/_workspace/src/github.com/coreos/gexpect"
import "log"

func main() {
	log.Printf("Testing Ping interact... \n")

	child, err := gexpect.Spawn("ping -c8 127.0.0.1")
	if err != nil {
		panic(err)
	}
	child.Interact()
	log.Printf("Success\n")
}
