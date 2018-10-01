package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sync"
)

type rpcFunc func(chan *protocol.RPC, messageID, *Contact)
type messageID [20]byte

var m = make(map[messageID]chan *protocol.RPC)
var mutexMap sync.Mutex

// Create a routine for the provided function and a way to send message responses back to it
// TODO logging
func createRoutine(executor RequestExecutor) {
	// Create channel for sending RPC responses
	c := make(chan *protocol.RPC)

	// Generate random messageID
	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}

	// Save (ID, channel) in map
	mutexMap.Lock()
	m[messageID] = c
	mutexMap.Unlock()

	log.WithFields(log.Fields{
		"Map":       m,
		"messageID": messageID,
	}).Debug("Creating a new routine.")

	executor.setChannel(c)
	executor.setMessageId(messageID)
	// Call function with channel
	go executor.execute()
}

func destroyRoutine(id messageID) {
	mutexMap.Lock()
	delete(m, id)
	mutexMap.Unlock()
}

func sendMessageToRoutine(msg *protocol.RPC) {
	// TODO what if message is not a response
	// Read messageID from message
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])
	originalSender := KademliaIDFromSlice(msg.OriginalSender)

	// This RPC message is not a response, handle the request
	if MyRoutingTable.me.ID.Equals(originalSender) {
		log.WithFields(log.Fields{
			"ID": id,
		}).Debug("Recieved message is a response to a request.")
		// Get channel
		mutexMap.Lock()
		c := m[id]
		if c != nil {
			// Write message to channel
			c <- msg
		}
		mutexMap.Unlock()
	} else {
		log.WithFields(log.Fields{
			"ID": id,
		}).Debug("Recieved message is a request.")

		switch msgType := msg.MessageType; msgType {
		case protocol.RPC_PING:
			go answerPingRequest(msg)
		case protocol.RPC_STORE:
			go answerStoreRequest(msg)
		case protocol.RPC_FIND_NODE:
			go answerFindNodeRequest(msg)
		case protocol.RPC_FIND_VALUE:
			go answerFindValueRequest(msg)
		case protocol.RPC_PIN:
			//TODO
		case protocol.RPC_UNPIN:
			//TODO
		default:
			log.WithFields(log.Fields{
				"Message": msg,
			}).Error("Message type not supported.")
		}
	}
}
