package main

import (
	"./kademlia"
	"fmt"
	"os"
	"runtime"
)

func main() {
	//kademlia.Listen("10.0.0.1", 4658)
	fmt.Println("Hello, World!")

	// Get bootstrap node address
	// TODO if empty that means this is the bootstraping node
	dns_name := os.Getenv("BOOTSTRAP_ADDR")
	port := os.Getenv("LISTEN_PORT")

	fmt.Println(dns_name)
	fmt.Print("Port: ")
	fmt.Println(port)

	go kademlia.Listen(port)
	fmt.Println("Started listening.")

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
