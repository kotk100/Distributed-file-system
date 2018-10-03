package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type FindNodeRequestExecutor struct {
	ch       chan *protocol.RPC
	id       messageID
	contact  *Contact
	target   *KademliaID
	callback *FindNodeRequestCallback
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) execute() {
	// Send ping message to other node
	var network Network
	error := network.SendFindContactMessage(findNodeRequestExecutor.target,
		MyRoutingTable.me.ID,
		findNodeRequestExecutor.contact,
		make([]Contact, 0),
		findNodeRequestExecutor.id)
	//if the channel return nil then there was error
	if error {
		log.Info("Error to send FindNode message.")
		destroyRoutine(findNodeRequestExecutor.id)
		if findNodeRequestExecutor.callback != nil {
			(*findNodeRequestExecutor.callback).errorRequest(*findNodeRequestExecutor.contact)
		}
	} else {
		timeout := NewTimeout(findNodeRequestExecutor.id, findNodeRequestExecutor.ch)
		timeout.start()
		// Recieve response message through channel
		rpc := <-findNodeRequestExecutor.ch
		timeout.stop()

		log.Info("Received FindNode message response.")

		if findNodeRequestExecutor.callback != nil {
			if rpc == nil {
				log.Error("find node request time out.")
				(*findNodeRequestExecutor.callback).errorRequest(*findNodeRequestExecutor.contact)
			} else {
				// Parse ping message and create contact
				findNode := parseFindNodeRequest(rpc)
				contactSender := createContactFromRPC(rpc)
				contacts := FindNode_ContactToContact(findNode.Contacts)
				(*findNodeRequestExecutor.callback).successRequest(*contactSender, contacts)

				MyRoutingTable.AddContactAsync(*contactSender)
			}
		}
		destroyRoutine(findNodeRequestExecutor.id)
	}
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) setChannel(ch chan *protocol.RPC) {
	findNodeRequestExecutor.ch = ch
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) setMessageId(id messageID) {
	findNodeRequestExecutor.id = id
}
