package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type FindNodeRequestExecutor struct {
	ch         chan *protocol.RPC
	id         messageID
	contact    Contact
	target     *KademliaID
	callback   *FindNodeRequestCallback
	networkAPI NetworkAPI
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) execute() {
	// Send ping message to other node
	error := findNodeRequestExecutor.networkAPI.SendFindContactMessage(findNodeRequestExecutor.target,
		MyRoutingTable.me.ID,
		&findNodeRequestExecutor.contact,
		make([]Contact, 0),
		findNodeRequestExecutor.id)
	//if the channel return nil then there was error
	if error {
		log.Info("Error to send FindNode message.")
		destroyRoutine(findNodeRequestExecutor.id)
		if findNodeRequestExecutor.callback != nil {
			(*findNodeRequestExecutor.callback).errorRequest(findNodeRequestExecutor.contact)
		}
	} else {
		timeout := NewTimeout(findNodeRequestExecutor.id, findNodeRequestExecutor.ch)
		timeout.start()
		// Recieve response message through channel
		rpc := <-findNodeRequestExecutor.ch
		timeout.stop()
		destroyRoutine(findNodeRequestExecutor.id)

		if findNodeRequestExecutor.callback != nil {
			if rpc == nil {
				log.WithFields(log.Fields{
					"contact": findNodeRequestExecutor.contact,
				}).Error("find node request time out.")
				(*findNodeRequestExecutor.callback).errorRequest(findNodeRequestExecutor.contact)
			} else {
				// Parse ping message and create contact
				findNode := parseFindNodeRequest(rpc)
				contactSender := createContactFromRPC(rpc)
				contacts := FindNode_ContactToContact(findNode.Contacts)
				log.WithFields(log.Fields{
					"contact sender": *contactSender,
				}).Info("FIND node message received.")
				(*findNodeRequestExecutor.callback).successRequest(*contactSender, contacts)

				MyRoutingTable.AddContactAsync(*contactSender)
			}
		}
	}
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) setChannel(ch chan *protocol.RPC) {
	findNodeRequestExecutor.ch = ch
}

func (findNodeRequestExecutor *FindNodeRequestExecutor) setMessageId(id messageID) {
	findNodeRequestExecutor.id = id
}
