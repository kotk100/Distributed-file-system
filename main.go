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
	"strconv"
)

type StoreFileResponse struct {
	Error    string
	HashFile string
}

type UnpinResponse struct {
	Error string
	Info  string
}

type GetFileResponse struct {
	Error string
	Info  string
}

type GetPinResponse struct {
	Error string
	Info  string
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
		store := kademlia.NewStoreRestRequest(&buffer, []byte(password), header.Filename, &w, r)
		store.StartStore()
	}
}

func GetFile(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	hashValue, err := kademlia.StringToHash(params["hashfile"])
	if err != nil {
		w.Header().Set("file_name", "NULL")
		getFileResponse := GetFileResponse{"error with hash file", ""}
		json.NewEncoder(w).Encode(getFileResponse)
	} else {
		testFindValue := kademlia.NewLookupValueManager(hashValue[:], &w, r)
		testFindValue.FindValue()
	}
}

func UnpinFile(w http.ResponseWriter, r *http.Request) {
	password := r.FormValue("password")
	hashfile := r.FormValue("hashfile")
	unpin := kademlia.CreateUnpinExecutor(hashfile, []byte(password))
	unpin.StartUnpin()
	unpinResponse := UnpinResponse{"", "Unpin executed"}
	json.NewEncoder(w).Encode(unpinResponse)
}

func Pin(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	res := kademlia.PinFileHash(params["hashfile"])
	if res {
		unpinResponse := UnpinResponse{"", "pin executed"}
		json.NewEncoder(w).Encode(unpinResponse)
	} else {
		unpinResponse := UnpinResponse{"error appeared during pin execution", ""}
		json.NewEncoder(w).Encode(unpinResponse)
	}
}

func setRestAPI() {
	c := cors.New(cors.Options{
		AllowCredentials: true,
		AllowedHeaders:   []string{"Content-Type", "Content-Disposition", "Origin", "Accept", "X-Requested-With"},
		AllowedOrigins:   []string{"*"}, // All origins
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "DELETE", "OPTIONS"},
		ExposedHeaders:   []string{"file_name"},
	})
	router := mux.NewRouter() //TODO chnage route name
	router.HandleFunc("/file", StoreFile)
	router.HandleFunc("/unpin", UnpinFile)
	router.HandleFunc("/file/{hashfile}", GetFile)
	router.HandleFunc("/pin/{hashfile}", Pin)
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

	//dRT := kademlia.DisplayRoutingTableClock{}
	//go dRT.Display()

	// Test saving files
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Println("Select command: \n1. for STORE\n2. for FIND_VALUE\n3. for Unpin\n4. Print stack traces\n5. Print routing table\n6. Print specific bucket")
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
			hashValue, _ := kademlia.StringToHash(hash)
			testFindValue := kademlia.NewLookupValueManager(hashValue[:], nil, nil)
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
		case "6":
			fmt.Println("Write bucket number:")
			scanner.Scan()
			num, _ := strconv.Atoi(scanner.Text())
			kademlia.MyRoutingTable.PrintBucket(num)
		default:
			continue
		}
	}

	// Wait forever
	for {
		runtime.Gosched()
	}
}
