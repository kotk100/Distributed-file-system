package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type StoreRequestExecutor struct {
	ch      chan *protocol.RPC
	id      messageID
	store   *Store
	contact Contact
}

func (storeRequestExecutor *StoreRequestExecutor) execute() {
	// Send store message to other node

	var network Network
	error := network.SendPingMessage(MyRoutingTable.me.ID, requestExecutorPing.contact, requestExecutorPing.id)

	//if the channel return nil then there was error
	if error {
		log.Info("Error to send ping.")
		destroyRoutine(requestExecutorPing.id)
		requestExecutorPing.pingCallback.pingResult(false)
	} else {
		timeout := NewTimeout(requestExecutorPing.id, requestExecutorPing.ch)
		timeOutManager.insertAndStart(timeout)
		// Recieve response message through channel
		rpc := <-requestExecutorPing.ch
		if optionalTimeout := timeOutManager.tryGetAndRemoveTimeOut(requestExecutorPing.id); optionalTimeout != nil {
			optionalTimeout.stop()
		}
		log.Info("Received PING message response.")
		if rpc == nil {
			log.Error("ping request time out.")
			if requestExecutorPing.pingCallback != nil {
				requestExecutorPing.pingCallback.pingResult(false)
			}
		} else {
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

			if requestExecutorPing.pingCallback != nil {
				requestExecutorPing.pingCallback.pingResult(true)
			}
		}
	}
}

func (storeRequestExecutor *StoreRequestExecutor) setChannel(ch chan *protocol.RPC) {
	storeRequestExecutor.ch = ch
}

func (storeRequestExecutor *StoreRequestExecutor) setMessageId(id messageID) {
	storeRequestExecutor.id = id
}
