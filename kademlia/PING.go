package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func SendAndRecievePing(contact *Contact,pingCallback PingCallback) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending PING message to node.")
	pingRequestExecutor:= PingRequestExecutor{}
	pingRequestExecutor.contact = contact
	pingRequestExecutor.pingCallback = pingCallback
	createRoutine(&pingRequestExecutor)
}

// Parse message inside RPC
func parsePingRPC(rpc *protocol.RPC) *protocol.Ping {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_PING {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. PING expected.")
	}

	// Parse message as Ping
	ping := &protocol.Ping{}
	if err := proto.Unmarshal(rpc.Message, ping); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Ping message.")
	}

	return ping
}

func createContactFromPing(ping *protocol.Ping, rpc *protocol.RPC) *Contact {

	// Create contact
	contact := &Contact{}
	contact.Address = rpc.IPaddress
	contact.ID = KademliaIDFromSlice(ping.KademliaID)

	return contact
}

func answerPingRequest(msg *protocol.RPC) {
	// Parse ping message and create contact
	ping := parsePingRPC(msg)
	contact := createContactFromPing(ping, msg)

	// Create message ID
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])

	// Send ping response
	log.WithFields(log.Fields{
		"Contact":   contact,
		"MessageID": id,
	}).Info("Sending PING response message.")
	net := &Network{}
	originalSender :=KademliaIDFromSlice(msg.OriginalSender)
	net.SendPingMessage(originalSender,contact, id)

	MyRoutingTable.AddContactAsync(*contact)
}
