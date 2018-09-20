package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
)

func SendAndRecievePing(contact *Contact) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending PING message to node.")

	createRoutine(xxx_sendAndRecievePing, contact)
}

func xxx_sendAndRecievePing(ch chan *protocol.RPC, id messageID, contact *Contact) {
	// Send ping message to other node
	var network Network
	network.SendPingMessage(contact, id)

	// Recieve response message through channel
	rpc := <-ch
	log.Info("Received PING message response.")

	// Parse ping message and create contact
	ping := parsePingRPC(rpc)
	contactSender := createContactFromPing(ping, rpc)

	if contactSender.ID != contact.ID {
		log.WithFields(log.Fields{
			"contactSender": contactSender,
			"contact":       contact,
		}).Error("Ping response from wrong Node recieved.")
	}

	// TODO update routingTable
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Updating routing table because of PING response.")
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
	// Read the port number nodes are listening on
	port := os.Getenv("LISTEN_PORT")

	// Create contact
	contact := &Contact{}
	contact.Address = rpc.IPaddress + port
	contact.ID = KademliaIDFromSlice(ping.KademliaID)

	return contact
}

func answerPingRequest(msg *protocol.RPC) {
	// Parse ping message and create contact
	ping := parsePingRPC(msg)
	contact := createContactFromPing(ping, msg)

	// Create message ID
	id := messageID{}
	copy(id[:], msg.MessageID[0:19])

	// Send ping response
	log.WithFields(log.Fields{
		"Contact":   contact,
		"MessageID": id,
	}).Info("Sending PING response message.")
	net := &Network{}
	net.SendPingMessage(contact, id)

	// TODO ADD contact to routing table
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Updating routing table because of PING response.")
}
