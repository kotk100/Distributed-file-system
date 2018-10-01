package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type FindValueRequestExecutor struct{
	ch       chan *protocol.RPC
	id       messageID
	contact  *Contact
	fileHash []byte
	callback *FindValueRequestCallback
}

func (findValueRequestExecutor *FindValueRequestExecutor) execute() {
	// Send ping message to other node
	var network Network
	error := network.SendFindDataMessage(findValueRequestExecutor.fileHash,
		findValueRequestExecutor.contact,
		make([]Contact, 0),
		findValueRequestExecutor.id,
		MyRoutingTable.me.ID)
	//if the channel return nil then there was error
	if error {
		log.Info("Error to send FindValue message.")
		destroyRoutine(findValueRequestExecutor.id)
		if findValueRequestExecutor.callback != nil {
			(*findValueRequestExecutor.callback).errorRequest(*findValueRequestExecutor.contact)
		}
	} else {
		timeout := NewTimeout(findValueRequestExecutor.id, findValueRequestExecutor.ch)
		timeOutManager.insertAndStart(timeout)
		// Recieve response message through channel
		rpc := <-findValueRequestExecutor.ch
		if optionalTimeout := timeOutManager.tryGetAndRemoveTimeOut(findValueRequestExecutor.id); optionalTimeout != nil {
			optionalTimeout.stop()
		}
		log.Info("Received FindNode message response.")

		if findValueRequestExecutor.callback != nil {
			if rpc == nil {
				log.Error("find value request time out.")
				(*findValueRequestExecutor.callback).errorRequest(*findValueRequestExecutor.contact)
			} else {
				// Parse ping message and create contact
				findValue := parseFindValueRequest(rpc)
				contactSender := createContactFromFindValue(rpc)
				contacts := FindValue_ContactToContact(findValue.Contacts)
				(*findValueRequestExecutor.callback).successRequest(*contactSender, contacts,findValue.HaveTheFile)

				MyRoutingTable.AddContactAsync(*contactSender)
			}
		}
		destroyRoutine(findValueRequestExecutor.id)
	}
}

func (findValueRequestExecutor *FindValueRequestExecutor) setChannel(ch chan *protocol.RPC) {
	findValueRequestExecutor.ch = ch
}

func (findValueRequestExecutor *FindValueRequestExecutor) setMessageId(id messageID) {
	findValueRequestExecutor.id = id
}
