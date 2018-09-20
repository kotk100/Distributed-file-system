package main

import (
	"./kademlia"
	"fmt"
	"os"
	"runtime"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log output
  log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
  //log.SetLevel(log.WarnLevel)
}

func main() {
	//kademlia.Listen("10.0.0.1", 4658)
	log.Info("Hello, World!")

	// Get bootstrap node address
	// TODO if empty that means this is the bootstraping node
	dns_name := os.Getenv("BOOTSTRAP_ADDR")
	port := os.Getenv("LISTEN_PORT")

	log.Info(dns_name)
	log.Info("Port: ")
	log.Info(port)

	go kademlia.Listen(port)
	log.Info("Started listening.")

	if dns_name != "" {
		var network kademlia.Network
		network.SendPingMessageTMP(dns_name)
		fmt.Println("Message sent.")
	}

	// Wait forever
	for {
		runtime.Gosched()
	}
}
