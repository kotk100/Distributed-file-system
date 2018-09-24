package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
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
	error := network.SendPingMessage(MyRoutingTable.me.ID,contact, id)

	//if the channel return nil then there was error
	if error {
		log.Info("Error to send ping.")
		destroyRoutine(id)
	}else{
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
			}).Error("Ping response from wrong Node recieved. Occurs when a new node joins the network.")
		}

		MyRoutingTable.AddContact(*contactSender)
	}
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

	MyRoutingTable.AddContact(*contact)
}
