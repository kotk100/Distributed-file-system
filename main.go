package main

import (
	"./kademlia"
	"bufio"
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
	dns_name := os.Getenv("BOOTSTRAP_ADDR")
	port := os.Getenv("LISTEN_PORT")

	kademlia.InitMyInformation(port)

	log.Info(dns_name)

	go kademlia.Listen(port)
	log.Info("Started listening.")

	// If node is not BN start bootstraping process
	if dns_name != "" {
		contact := &kademlia.Contact{}
		contact.Address = dns_name

		// Bootstap this node onto the network
		bootstrap := kademlia.Bootstrap{}
		bootstrap.BootstrapNode = *contact
		kademlia.SendAndRecievePing(contact, &bootstrap)
	}

	dRT := kademlia.DisplayRoutingTableClock{}
	go dRT.Display()

	// Test saving files
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		buffer := scanner.Bytes()

		store := kademlia.CreateNewStore(&buffer, nil, "example.txt")
		store.StartStore()
	}

	// Wait forever
	for {
		runtime.Gosched()
	}
}
