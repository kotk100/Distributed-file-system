package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"math/rand"
)

type rpcFunc func(chan *protocol.RPC, messageID, *Contact)
type messageID [20]byte

var m = make(map[messageID]chan *protocol.RPC)

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
	m[messageID] = c

	log.WithFields(log.Fields{
		"Map":       m,
		"messageID": messageID,
	}).Debug("Creating a new routine.")

	executor.setChannel(c)
	executor.setMessageId(messageID)
	// Call function with channel
	go executor.execute()
}


func destroyRoutine(id messageID){
	delete(m, id)
}

func sendMessageToRoutine(msg *protocol.RPC) {
	// TODO what if message is not a response
	// Read messageID from message
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])
	originalSender :=KademliaIDFromSlice(msg.OriginalSender)

	// This RPC message is not a response, handle the request
	if MyRoutingTable.me.ID.Equals(originalSender) {
		log.WithFields(log.Fields{
			"ID": id,
		}).Debug("Recieved message is a response to a request.")
		// Get channel
		c := m[id]
		if c!=nil {
			// Write message to channel
			c <- msg
		}
	}else{
		log.WithFields(log.Fields{
			"ID": id,
		}).Debug("Recieved message is a request.")

		switch msgType := msg.MessageType; msgType {
		case protocol.RPC_PING:
			answerPingRequest(msg)
		case protocol.RPC_STORE:
			//TODO
		case protocol.RPC_FIND_NODE:
			answerFindNodeRequest(msg)
		case protocol.RPC_FIND_VALUE:
			//TODO
		case protocol.RPC_PIN:
			//TODO
		case protocol.RPC_UNPIN:
			//TODO
		default:
			log.WithFields(log.Fields{
				"Message": msg,
			}).Error("Failed to parse incomming RPC message.")
		}
	}
}
