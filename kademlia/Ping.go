package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

func SendAndRecievePing(contact Contact, pingCallback PingCallback) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending PING message to node.")
	pingRequestExecutor := PingRequestExecutor{}
	pingRequestExecutor.contact = contact
	pingRequestExecutor.pingCallback = pingCallback
	createRoutine(&pingRequestExecutor)
}


func answerPingRequest(msg *protocol.RPC) {
	// Parse ping message and create contact
	contact := createContactFromRPC(msg)

	// Create message ID
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])

	// Send ping response
	log.WithFields(log.Fields{
		"Contact":   contact,
		"MessageID": id,
	}).Info("Sending PING response message.")
	net := &Network{}
	originalSender := KademliaIDFromSlice(msg.OriginalSender)
	net.SendPingMessage(originalSender, contact, id)

	MyRoutingTable.AddContactAsync(*contact)
}
