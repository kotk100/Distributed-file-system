package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

func sendAndRecievePing(contact *Contact) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending PING message to node.")

	createRoutine(XXX_sendAndRecievePing, contact)
}

func XXX_sendAndRecievePing(ch chan *protocol.RPC, id messageID, contact *Contact) {
	// Send ping message to other node
	var network Network
	network.SendPingMessage(contact, id)

	// Recieve response message through channel
	rpc := <-ch
	log.Info("Received PING message response.")

	// Parse message
	if rpc.MessageType != protocol.RPC_PING {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. PING expected.")
	}

	// TODO update routingTable
	log.Info("Updating routing table because of PING response.")
}
