package main

import (
	"./kademlia"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
)

type StoreFileResponse struct {
	Error    string
	HashFile string
}

type UnpinResponse struct {
	Error string
	Info string
}

func init() {
	// Log output
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func StoreFile(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	password := r.FormValue("password")
	if err != nil {
		storeFileResponse := StoreFileResponse{"error to access to the file", ""}
		json.NewEncoder(w).Encode(storeFileResponse)
	} else {
		buffer := make([]byte, header.Size)
		file.Read(buffer)
		defer file.Close()
		store := kademlia.CreateNewStore(&buffer, []byte(password), header.Filename)
		store.StartStore()
		storeFileResponse := StoreFileResponse{"", store.GetHash()}
		json.NewEncoder(w).Encode(storeFileResponse)
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	testFindValue := kademlia.NewLookupValueManager(kademlia.StringToHash(params["hashfile"])[:], w, r)
	testFindValue.FindValue()
}

func UnpinFile(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	hashfile := r.FormValue("hashfile")
	unpin := kademlia.CreateUnpinExecutor(hashfile, []byte(password))
	unpin.StartUnpin()
	unpinResponse := UnpinResponse{"","Unpin executed"}
	json.NewEncoder(w).Encode(unpinResponse)
}

func setRestAPI() {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Content-Disposition", "Origin", "Accept", "X-Requested-With"},
		AllowedOrigins:   []string{"*"}, // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposedHeaders:   []string{"file_name"},
	})
	router := mux.NewRouter()
	router.HandleFunc("/file", StoreFile)
	router.HandleFunc("/unpin", UnpinFile)
	router.HandleFunc("/file/{hashfile}", GetFile)
	log.Fatal(http.ListenAndServe(":3000", c.Handler(router)))
}

func main() {
	//kademlia.Listen("10.0.0.1", 4658)
	log.Info("Hello, World!")

	go setRestAPI()

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

	dRT := kademlia.DisplayRoutingTableClock{}
	go dRT.Display()

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
			//hash := scanner.Text()

			//testFindValue := kademlia.NewLookupValueManager(kademlia.StringToHash(hash)[:])
			//go testFindValue.FindValue()
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
