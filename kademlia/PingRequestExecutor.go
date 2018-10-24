package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
)

type PingCallback interface {
	pingResult(receivedAnswer bool)
}

type PingRequestExecutor struct {
	ch           chan *protocol.RPC
	id           messageID
	contact      Contact
	pingCallback PingCallback
}

func (requestExecutorPing *PingRequestExecutor) execute() {
	// Send ping message to other node
	var network Network
	error := network.SendPingMessage(MyRoutingTable.me.ID, &requestExecutorPing.contact, requestExecutorPing.id)

	//if the channel return nil then there was error
	if error {
		log.Info("Error to send ping.")
		destroyRoutine(requestExecutorPing.id)
		requestExecutorPing.pingCallback.pingResult(false)
	} else {
		timeout := NewTimeout(requestExecutorPing.id, requestExecutorPing.ch)
		timeout.start()
		// Recieve response message through channel
		rpc := <-requestExecutorPing.ch
		timeout.stop()
		destroyRoutine(requestExecutorPing.id)

		log.Info("Received PING message response.")
		if rpc == nil {
			log.Error("ping request time out.")
			if requestExecutorPing.pingCallback != nil {
				requestExecutorPing.pingCallback.pingResult(false)
			}
		} else {
			// Parse ping message and create contact
			contactSender := createContactFromRPC(rpc)

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

func (requestExecutorPing *PingRequestExecutor) setChannel(ch chan *protocol.RPC) {
	requestExecutorPing.ch = ch
}

func (requestExecutorPing *PingRequestExecutor) setMessageId(id messageID) {
	requestExecutorPing.id = id
}
