package main

import (
	"./kademlia"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)


func init() {
	// Log output
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}


func main() {
	//kademlia.Listen("10.0.0.1", 4658)
	log.Info("Hello, World!")

	// Get bootstrap node address
	// TODO if empty that means this is the bootstraping node
	dns_name := os.Getenv("BOOTSTRAP_ADDR")
	port := os.Getenv("LISTEN_PORT")

	kademlia.InitMyInformation(port)

	log.Info(dns_name)

	go kademlia.Listen(port)
	log.Info("Started listening.")

	if dns_name != "" {
		contact := &kademlia.Contact{}
		contact.Address = dns_name


		bootstrap:=kademlia.Bootstrap{}
		bootstrap.BootstrapNode = *contact
		kademlia.SendAndRecievePing(contact,&bootstrap)

		fmt.Println("Message sent.")
	}

	dRT:=kademlia.DisplayRoutingTableClock{}
	dRT.Display()

	// Wait forever
	for {
		runtime.Gosched()
	}
}
