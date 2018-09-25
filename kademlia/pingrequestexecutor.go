package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type PingRequestExecutor struct{
	ch chan *protocol.RPC
	id messageID
	contact *Contact
}


func (requestExecutorPing *PingRequestExecutor) execute(){
	// Send ping message to other node
	var network Network
	error := network.SendPingMessage(MyRoutingTable.me.ID,requestExecutorPing.contact, requestExecutorPing.id)

	//if the channel return nil then there was error
	if error {
		log.Info("Error to send ping.")
		destroyRoutine(requestExecutorPing.id)
	}else{
		// Recieve response message through channel
		rpc := <-requestExecutorPing.ch
		log.Info("Received PING message response.")

		// Parse ping message and create contact
		ping := parsePingRPC(rpc)
		contactSender := createContactFromPing(ping, rpc)

		if contactSender.ID != requestExecutorPing.contact.ID {
			log.WithFields(log.Fields{
				"contactSender": contactSender,
				"contact":       requestExecutorPing.contact,
			}).Error("Ping response from wrong Node recieved. Occurs when a new node joins the network.")
		}

		MyRoutingTable.AddContact(*contactSender)
	}
}

func (requestExecutorPing *PingRequestExecutor) setChannel(ch chan *protocol.RPC){
	requestExecutorPing.ch = ch
}


func (requestExecutorPing *PingRequestExecutor) setMessageId(id messageID){
	requestExecutorPing.id = id
}