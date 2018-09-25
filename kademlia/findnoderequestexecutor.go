package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type FindNodeRequestExecutor struct{
	ch chan *protocol.RPC
	id messageID
	contact *Contact
	target *KademliaID
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) execute(){
	// Send ping message to other node
	var network Network
	error := network.SendFindContactMessage(findNodeRequestExecutor.target,
		MyRoutingTable.me.ID,
		findNodeRequestExecutor.contact,
		make([]Contact,0),
		findNodeRequestExecutor.id)
	//if the channel return nil then there was error
	if error {
		log.Info("Error to send FindNode message.")
		destroyRoutine(findNodeRequestExecutor.id)
	}else{
		// Recieve response message through channel
		rpc := <-findNodeRequestExecutor.ch
		log.Info("Received FindNode message response.")
		// Parse ping message and create contact
		findNode := parseFindNodeRequest(rpc)
		contactSender := createContactFromFindNode(findNode, rpc)

		MyRoutingTable.AddContact(*contactSender)
	}
}


func (findNodeRequestExecutor *FindNodeRequestExecutor) setChannel(ch chan *protocol.RPC){
	findNodeRequestExecutor.ch = ch
}


func (findNodeRequestExecutor *FindNodeRequestExecutor) setMessageId(id messageID){
	findNodeRequestExecutor.id = id
}