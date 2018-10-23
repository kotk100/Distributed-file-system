package kademlia

import (
	"./protocol"
	"crypto/sha1"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

type UnpinOperationsStruct struct {
	pastUnpinRequests map[string]bool
	lock              sync.Mutex
}

var unpinReference *UnpinOperationsStruct

func CreateUnpinOperationsStruct() {
	unpinReference = &UnpinOperationsStruct{}
	unpinReference.pastUnpinRequests = make(map[string]bool)
}

type UnpinExecutor struct {
	ch       chan *protocol.RPC
	id       messageID
	filehash string
	password string
}

func CreateUnpinExecutor(filehash string, password []byte) *UnpinExecutor {
	unpin := &UnpinExecutor{}
	unpin.filehash = filehash

	passHash := sha1.Sum(password)
	str := hashToString(passHash[:])
	unpin.password = str

	return unpin
}

func (unpinExecutor *UnpinExecutor) StartUnpin() {
	createRoutine(unpinExecutor)
}

func (unpinExecutor *UnpinExecutor) execute() {
	// Remove file localy
	success := checkAndRemoveFile(unpinExecutor.filehash, unpinExecutor.password)

	if success {
		// Add to set
		unpinReference.lock.Lock()
		unpinReference.pastUnpinRequests[unpinExecutor.filehash] = true
		unpinReference.lock.Unlock()

		// Create a task to remove hash
		deleteUnpinRequestTask := &Task{}
		deleteUnpinRequestTask.taskType = DeleteUnpinRequest
		deleteUnpinRequestTask.id = unpinExecutor.filehash
		deleteUnpinRequestTask.executeEvery = 0

		delTask := &DeleteUnpinRequestTask{}
		deleteUnpinRequestTask.executor = delTask
		timeToExecute := time.Now().Add(100 * time.Second)
		PeriodicTasksReference.addTask(&timeToExecute, deleteUnpinRequestTask)

		log.WithFields(log.Fields{
			"Filehash": unpinExecutor.filehash,
		}).Info("Removed file, sending unpin requests.")

		// Send request to all other nodes
		net := &Network{}
		contacts := MyRoutingTable.getAllContacts()
		for _, contact := range *contacts {
			net.SendUnpinMessage(unpinExecutor.filehash, unpinExecutor.password, &contact, unpinExecutor.id[:])
		}
	} else {
		log.WithFields(log.Fields{
			"File": unpinExecutor.filehash,
		}).Error("Failed to remove file. Unpin operation aborted.")
	}

	destroyRoutine(unpinExecutor.id)
}

func (unpinExecutor *UnpinExecutor) setChannel(ch chan *protocol.RPC) {
	unpinExecutor.ch = ch
}

func (unpinExecutor *UnpinExecutor) setMessageId(id messageID) {
	unpinExecutor.id = id
}

// Parse message inside RPC
func parseUnpinRPC(rpc *protocol.RPC) *protocol.UnPin {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_UNPIN {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. STORE ANSWER expected.")

		return nil
	}

	// Parse message as Store
	store := &protocol.UnPin{}
	if err := proto.Unmarshal(rpc.Message, store); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Store answer message.")
	}

	return store
}

func (unpinOperationsStruct *UnpinOperationsStruct) answerUnpinRequest(rpc *protocol.RPC) {
	// Parse message
	unpin := parseUnpinRPC(rpc)

	log.WithFields(log.Fields{
		"Filehash": unpin.FileHash,
	}).Info("Received unpin request.")

	// Check if request already serviced
	unpinOperationsStruct.lock.Lock()
	serviced := unpinOperationsStruct.pastUnpinRequests[unpin.FileHash]
	unpinOperationsStruct.lock.Unlock()

	if serviced {
		// Do nothing
	} else {
		// Add to set
		unpinOperationsStruct.lock.Lock()
		unpinOperationsStruct.pastUnpinRequests[unpin.FileHash] = true
		unpinOperationsStruct.lock.Unlock()

		// Create a task to remove hash
		deleteUnpinRequestTask := &Task{}
		deleteUnpinRequestTask.taskType = DeleteUnpinRequest
		deleteUnpinRequestTask.id = unpin.FileHash
		deleteUnpinRequestTask.executeEvery = 0

		delTask := &DeleteUnpinRequestTask{}
		deleteUnpinRequestTask.executor = delTask
		timeToExecute := time.Now().Add(100 * time.Second)
		PeriodicTasksReference.addTask(&timeToExecute, deleteUnpinRequestTask)

		// Remove file
		if checkAndRemoveFile(unpin.FileHash, unpin.PasswordHash) {
			// Create message ID
			id := messageID{}
			copy(id[:], rpc.MessageID[0:20])

			// Send Unpin messages to all other nodes
			log.WithFields(log.Fields{
				"Filehash": unpin.FileHash,
			}).Info("Removed file, sending unpin requests.")

			net := &Network{}
			contacts := MyRoutingTable.getAllContacts()
			for _, contact := range *contacts {
				net.SendUnpinMessage(unpin.FileHash, unpin.PasswordHash, &contact, rpc.MessageID)
			}
		} else {
			log.WithFields(log.Fields{
				"Filehash": unpin.FileHash,
			}).Error("Password incorrect, not removing file.")
		}
	}

	// Parse Unpin message and create contact
	contact := createContactFromRPC(rpc)
	MyRoutingTable.AddContactAsync(*contact)
}

// TODO no feedback if file is removed to the user
func checkAndRemoveFile(filehash string, passwordHash string) bool {
	if checkFileExistsHash(filehash) {
		filename := getFilenameFromHash(filehash)
		s := strings.Split(filename, ":")

		// Check if the password hash matches the actual file and the password is set, otherwise don't allow deletion
		if len(s) >= 2 && s[1] != "" && s[1] == passwordHash {
			removeFileAndPeriodicTasks(filehash)
			return true
		} else {
			return false
		}
	} else {
		return true
	}
}

func removeFileAndPeriodicTasks(filehash string) {
	// Remove file from node
	removeFileByHash(filehash)

	// Remove republish task
	task := Task{}
	task.id = filehash
	task.taskType = RepublishFile
	PeriodicTasksReference.removeTask(&task)

	// Remove file deletion task
	task2 := Task{}
	task2.id = filehash
	task2.taskType = ExpireFile
	PeriodicTasksReference.removeTask(&task2)
}

type DeleteUnpinRequestTask struct {
	task *Task
}

func (deleteUnpinRequestTask *DeleteUnpinRequestTask) execute() {
	// Get filehash of the entry to
	filehash := deleteUnpinRequestTask.task.id

	unpinReference.lock.Lock()
	delete(unpinReference.pastUnpinRequests, filehash)
	unpinReference.lock.Unlock()
}

func (deleteUnpinRequestTask *DeleteUnpinRequestTask) setTask(task *Task) {
	deleteUnpinRequestTask.task = task
}
