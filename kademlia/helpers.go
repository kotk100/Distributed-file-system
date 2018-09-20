package kademlia

import (
	"./protocol"
	"fmt"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"runtime"
)

type Helper struct {
}

type rpcFunc func(chan *protocol.RPC, messageID, *Contact)
type messageID [20]byte

var m = make(map[messageID]chan *protocol.RPC)

// Create a routine for the provided function and a way to send message responses back to it
func (helper *Helper) createRoutine(fn rpcFunc, contact *Contact){
	// Create channel for sending RPC responses
	c := make(chan *protocol.RPC)

	// Generate random messageID
	id := messageID{}
	for i := 0; i < 20; i++ {
		id[i] = uint8(rand.Intn(256))
	}

	// Save (ID, channel) in map
	m[id] = c

	// Call function with channel
	go fn(c, id, contact)
}

func (helper *Helper) sendMessageToRoutine(msg *protocol.RPC){
	// TODO what if message is not a response
	// Read messageID from message
	var id messageID
	copy(id[:], msg.MessageID[0:20])

	// Get channel
	c := m[id]

	// This RPC message is not a response, handle the request
	if c == nil {
		switch msgType := msg.MessageType; msgType {
		case protocol.RPC_PING:
			ping.answerPingRequest()
		case protocol.RPC_STORE:

		default:
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to parse incomming RPC message.")
		}
	}

	// Write message to channel
	c <- msg
}

func sendAndRecievePing(c chan *protocol.RPC, id messageID, contact *Contact){
	var network Network

	network.SendPingMessage(contact, id);
	fmt.Println("Message sent.")

	// TODO recieve response message through channel
}

