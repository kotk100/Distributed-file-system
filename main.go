package main

import (
	"./kademlia"
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"time"
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

		time.Sleep(10 * time.Second)

		testFindValue:=kademlia.NewTestLookupValue(store.GetHashFile())
		testFindValue.StartTest()
		//check node if is not contains file
		//if not launch find value
		//display contact which contain file

	}

	// Wait forever
	for {
		runtime.Gosched()
	}
}
