package kademlia

import (
	"./protocol"
	"crypto/sha1"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	filename   string
	filehash   *[20]byte
	fileLength int64
}

func handleError(lenght int, err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"lenght bytes": lenght,
			"Error":        err,
		}).Error("Error creating filename.")
	}
}

// Creates a filename under which the file will be stored
func (store *Store) createFilename(file *[]byte, password []byte, pinned bool, originalFilename string) {
	hash := strings.Builder{}

	str := hashToString(store.filehash[:])
	handleError(hash.WriteString(str))
	handleError(hash.WriteString(":"))

	if password != nil {
		passHash := sha1.Sum(password)
		str = hashToString(passHash[:])
		handleError(hash.WriteString(str))
	}

	handleError(hash.WriteString(":"))
	handleError(hash.WriteString(originalFilename))

	// File is not pinned by default
	if pinned {
		handleError(hash.WriteString(":1"))
	} else {
		handleError(hash.WriteString(":0"))
	}

	store.filename = hash.String()
}

func createFileHash(file *[]byte) *[20]byte {
	hash := sha1.Sum(*file)
	return &hash
}

func CreateNewStore(file *[]byte, password []byte, originalFilename string) *Store {
	store := &Store{}
	store.filehash = createFileHash(file)
	store.createFilename(file, password, false, originalFilename)
	store.fileLength = int64(len(*file))

	// Store file localy with pin set
	// TODO how to store file on first node that is the publisher so that the whole file doesn't have to be in memory
	filename := store.filename[:len(store.filename)-1] + "1"

	fileWriter := createFileWriter(filename)
	if fileWriter == nil {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed to create a fileWriter.")

		return nil
	}

	err := fileWriter.StoreFileChunk(file)
	fileWriter.close()
	if err {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed writing to file.")

		return nil
	}

	// Create periodic task for republishing
	republishTask := &Task{}
	republishTask.taskType = RepublishFile
	republishTask.id = hashToString(store.filehash[:])

	timeString := os.Getenv("ORG_REPUBLISH_TIME")
	timeSeconds, errb := strconv.Atoi(timeString)
	if errb != nil {
		log.WithFields(log.Fields{
			"Error": errb,
		}).Error("Failer to convert string to int")
	}
	republishTask.executeEvery = time.Duration(timeSeconds) * time.Second

	repTask := &RepublishTask{}
	republishTask.executor = repTask
	timeToExecute := time.Now().Add(republishTask.executeEvery)

	PeriodicTasksReference.addTask(&timeToExecute, republishTask)

	// Create periodic task for file expiration
	// TODO expiry time should be calculated dynamicaly, recalculate when a new node is added into the routing table
	expiryTask := &Task{}
	expiryTask.taskType = ExpireFile
	expiryTask.id = republishTask.id

	timeString = os.Getenv("EXPARATION_TIME")
	timeSeconds, errb = strconv.Atoi(timeString)
	if errb != nil {
		log.WithFields(log.Fields{
			"Error": errb,
		}).Error("Failer to convert string to int")
	}
	expiryTask.executeEvery = time.Duration(timeSeconds) * time.Second

	expTask := &FileExpirationTask{}
	expiryTask.executor = expTask
	timeToExecute = time.Now().Add(expiryTask.executeEvery)

	PeriodicTasksReference.addTask(&timeToExecute, expiryTask)

	return store
}

func CreateNewStoreForRepublish(fileHash string) *Store {
	store := &Store{}

	hash := StringToHash(fileHash)
	store.filehash = hash

	store.filename = getFilenameFromHash(fileHash)
	store.fileLength = getFileLength(store.filename)

	return store
}

func (store *Store) GetHash() string {
	return hashToString(store.filehash[:])
}

func (store *Store) StartStore() {
	log.WithFields(log.Fields{
		"file hash": hashToString(store.filehash[:]),
	}).Info("Started STORE procedure.")

	target := Contact{}
	target.ID = KademliaIDFromSlice(store.filehash[:])

	// Find k closest nodes to store file on
	// Returns k closest to callback processKClosest in this file
	lookupNode := NewLookupNode(target, store)
	lookupNode.Start()
}

// Send store command to all k closest nodes
func (store *Store) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("Find_node finished for STORE procedure.")
	for _, contact := range KClosestOfTarget {
		storeEx := &StoreRequestExecutor{}
		storeEx.store = store
		storeEx.contact = (contact.contact)
		createRoutine(storeEx)
	}
}

// Parse message inside RPC
func parseStoreRPC(rpc *protocol.RPC) *protocol.Store {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_STORE {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. STORE expected.")

		return nil
	}

	// Parse message as Store
	store := &protocol.Store{}
	if err := proto.Unmarshal(rpc.Message, store); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Ping message.")
	}

	return store
}

// Parse message inside RPC
func parseStoreAnswerRPC(rpc *protocol.RPC) *protocol.StoreAnswer {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_STORE_ANSWER {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. STORE ANSWER expected.")

		return nil
	}

	// Parse message as Store
	store := &protocol.StoreAnswer{}
	if err := proto.Unmarshal(rpc.Message, store); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Store answer message.")
	}

	return store
}

func errorStoreAnswerWraper(err bool) {
	if err {
		log.Error("Failed to send Store Answer message.")
	}
}

func answerStoreRequest(rpc *protocol.RPC) {
	// Parse message and refresh routing table
	store := parseStoreRPC(rpc)
	other := createContactFromRPC(rpc)
	MyRoutingTable.AddContactAsync(*other)

	network := &Network{}
	id := messageID{}
	copy(id[:], rpc.MessageID[0:20])

	log.WithFields(log.Fields{
		"Filename": store.Filename,
		"Contact":  other,
	}).Info("Received STORE request.")

	// Check if file already exsists
	if checkFileExists(store.Filename) {
		errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_ALREADY_STORED, other, id, &rpc.OriginalSender))

		// Update periodic task for republishing and deletion
		task := &Task{}
		task.id = getHashFromFilename(store.Filename)
		task.taskType = RepublishFile
		PeriodicTasksReference.updateTask(task)

		task2 := &Task{}
		task2.id = getHashFromFilename(store.Filename)
		task2.taskType = ExpireFile
		PeriodicTasksReference.updateTask(task2)

	} else {
		// Start listening on port and save file after connecting
		portStr := os.Getenv("FILE_TRANSFER_PORT")
		err := network.RecieveFile(portStr, store.Filename, other, id, store.FileSize, &rpc.OriginalSender)

		// Send confirmation if file was saved correctly
		if err {
			log.WithFields(log.Fields{
				"Contact": other,
			}).Error("Failed to receive file and save it.")

			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_ERROR, other, id, &rpc.OriginalSender))
		} else {
			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_OK, other, id, &rpc.OriginalSender))

			// Create periodic task for republishing
			republishTask := &Task{}
			republishTask.taskType = RepublishFile
			republishTask.id = getHashFromFilename(store.Filename)

			timeString := os.Getenv("REPUBLISH_TIME")
			timeSeconds, err := strconv.Atoi(timeString)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Failer to convert string to int")
			}
			republishTask.executeEvery = time.Duration(timeSeconds) * time.Second

			repTask := &RepublishTask{}
			republishTask.executor = repTask
			timeToExecute := time.Now().Add(republishTask.executeEvery)

			PeriodicTasksReference.addTask(&timeToExecute, republishTask)

			// Create periodic tasks for expiration
			expiryTask := &Task{}
			expiryTask.taskType = ExpireFile
			expiryTask.id = republishTask.id

			timeString = os.Getenv("EXPARATION_TIME")
			timeSeconds, err = strconv.Atoi(timeString)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Failer to convert string to int")
			}
			expiryTask.executeEvery = time.Duration(timeSeconds) * time.Second

			expTask := &FileExpirationTask{}
			expiryTask.executor = expTask
			timeToExecute = time.Now().Add(expiryTask.executeEvery)

			PeriodicTasksReference.addTask(&timeToExecute, expiryTask)
		}
	}
}
