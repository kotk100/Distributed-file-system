package main

import (
	"./kademlia"
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"runtime/pprof"
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

	// Create a task scheduler
	periodicTasks := kademlia.CreatePeriodicTasks()
	kademlia.PeriodicTasksReference = periodicTasks

	// Create a Unpin handler
	kademlia.CreateUnpinOperationsStruct()

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
		kademlia.SendAndRecievePing(*contact, &bootstrap)
	}

	/*dRT := kademlia.DisplayRoutingTableClock{}
	go dRT.Display()*/

	// TODO infinite loop for find_node where two nodes think the other one is the original sender and constantly send eachother responses

	// Test saving files
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Select command: \n1. for STORE\n2. for FIND_VALUE\n3. for Unpin\n4. Print stack traces\n5. Print routing table")
		scanner.Scan()
		input := scanner.Text()

		switch input {
		case "1":
			fmt.Println("Write file contents:")
			scanner.Scan()
			buffer := scanner.Bytes()
			fmt.Println("Write password:")
			scanner.Scan()
			password := scanner.Bytes()

			store := kademlia.CreateNewStore(&buffer, password, "example.txt")
			store.StartStore()
			fmt.Print("File hash: ")
			fmt.Println(store.GetHash())
		case "2":
			fmt.Println("Write file hash:")
			scanner.Scan()
			hash := scanner.Text()

			testFindValue := kademlia.NewLookupValueManager(kademlia.StringToHash(hash)[:])
			go testFindValue.FindValue()
		case "3":
			fmt.Println("Write file hash:")
			scanner.Scan()
			hash := scanner.Text()
			fmt.Println("Write password:")
			scanner.Scan()
			password := scanner.Bytes()

			unpin := kademlia.CreateUnpinExecutor(hash, password)
			unpin.StartUnpin()
		case "4":
			pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
		case "5":
			kademlia.MyRoutingTable.Print()
		default:
			continue
		}
	}

	// Wait forever
	for {
		runtime.Gosched()
	}
}
