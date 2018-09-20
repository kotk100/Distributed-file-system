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
func createRoutine(fn rpcFunc, contact *Contact) {
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
	}).Info("Creating a new routine.")

	// Call function with channel
	go fn(c, messageID, contact)
}

func sendMessageToRoutine(msg *protocol.RPC) {
	// TODO what if message is not a response
	// Read messageID from message
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])

	// Get channel
	c := m[id]

	// This RPC message is not a response, handle the request
	if c == nil {
		log.WithFields(log.Fields{
			"ID": id,
		}).Info("Recieved message is a request.")

		switch msgType := msg.MessageType; msgType {
		case protocol.RPC_PING:
			answerPingRequest(msg)
		case protocol.RPC_STORE:
			//TODO
		case protocol.RPC_FIND_NODE:
			//TODO
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
	} else {
		log.WithFields(log.Fields{
			"ID": id,
		}).Info("Recieved message is a response to a request.")

		// Write message to channel
		c <- msg
	}
}
