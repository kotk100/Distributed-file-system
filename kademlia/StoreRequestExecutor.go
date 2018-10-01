package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"os"
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
	senderID := MyRoutingTable.GetMe().ID[:]
	error := network.SendStoreMessage(storeRequestExecutor.store.filename, int64(len(*storeRequestExecutor.store.file)), &storeRequestExecutor.contact, storeRequestExecutor.id, &senderID)
	log.WithFields(log.Fields{
		"Contact":  storeRequestExecutor.contact,
		"Filename": storeRequestExecutor.store.filename,
	}).Info("Send STORE message to contact.")

	//if the channel return nil then there was error
	if error {
		log.Error("Error sending Store RPC.")
		destroyRoutine(storeRequestExecutor.id)
	} else {
		timeout := NewTimeout(storeRequestExecutor.id, storeRequestExecutor.ch)
		timeOutManager.insertAndStart(timeout)
		// Recieve response message through channel
		rpc := <-storeRequestExecutor.ch
		if optionalTimeout := timeOutManager.tryGetAndRemoveTimeOut(storeRequestExecutor.id); optionalTimeout != nil {
			optionalTimeout.stop()
		}

		if rpc == nil {
			log.Error("Store request time out.")
		} else {
			log.WithFields(log.Fields{
				"Contact":  storeRequestExecutor.contact,
				"Filename": storeRequestExecutor.store.filename,
			}).Info("Received STORE ANSWER message response.")

			// Parse store message and create contact for adding to routing table
			store := parseStoreAnswerRPC(rpc)
			other := createContactFromRPC(rpc)
			MyRoutingTable.AddContactAsync(*other)

			switch answer := store.Answer; answer {
			// Other node is ready to start receiving the file
			case protocol.StoreAnswer_OK:
				// Send data (file) in chunks to other node
				port := os.Getenv("FILE_TRANSFER_PORT")

				log.WithFields(log.Fields{
					"Contact":  storeRequestExecutor.contact,
					"Filename": storeRequestExecutor.store.filename,
				}).Info("Send file to contact.")
				error = network.SendFile(storeRequestExecutor.store.filehash, &storeRequestExecutor.contact, port, int64(len(*storeRequestExecutor.store.file)))
				if error {
					log.WithFields(log.Fields{
						"Other node": storeRequestExecutor.contact,
						"Filehash":   storeRequestExecutor.store.filehash,
					}).Error("Error sending file.")
				}

				log.WithFields(log.Fields{
					"Contact":  storeRequestExecutor.contact,
					"Filename": storeRequestExecutor.store.filename,
				}).Info("Sending file finished.")

				// Wait for response to see if saving file succeded
				// Recieve response message through channel
				timeout := NewTimeout(storeRequestExecutor.id, storeRequestExecutor.ch)
				timeOutManager.insertAndStart(timeout)
				rpc := <-storeRequestExecutor.ch
				if optionalTimeout := timeOutManager.tryGetAndRemoveTimeOut(storeRequestExecutor.id); optionalTimeout != nil {
					optionalTimeout.stop()
				}

				if rpc == nil {
					log.Error("Sending file time out. No response recieved.")
				} else {
					log.WithFields(log.Fields{
						"Contact":  storeRequestExecutor.contact,
						"Filename": storeRequestExecutor.store.filename,
					}).Info("Received response after storing file.")

					// Parse store message and create contact for adding to routing table
					store := parseStoreAnswerRPC(rpc)
					other := createContactFromRPC(rpc)
					MyRoutingTable.AddContactAsync(*other)

					switch answer := store.Answer; answer {
					// Other node successfully saved file
					case protocol.StoreAnswer_OK:
						// Do nothing
					case protocol.StoreAnswer_ERROR:
						// TODO what hapens if a file cant be stored on node? (for example if there is not enough space)
						log.WithFields(log.Fields{
							"Other node": other,
							"FileName":   storeRequestExecutor.store.filename,
						}).Error("Failed to store file on node.")
					default:
						log.WithFields(log.Fields{
							"Answer": answer,
						}).Error("Wrong option from store answer response.")
					}
				}

			case protocol.StoreAnswer_ALREADY_STORED:
				// File already stored on node, don't have to do anything
			case protocol.StoreAnswer_ERROR:
				// TODO what hapens if a file cant be stored on node? (for example if there is not enough space)
				log.WithFields(log.Fields{
					"Other node": other,
					"FileName":   storeRequestExecutor.store.filename,
				}).Error("Failed to store file on node.")
			default:
				log.WithFields(log.Fields{
					"Answer": answer,
				}).Error("Wrong option from store answer response.")
			}
		}
		destroyRoutine(storeRequestExecutor.id)
	}
}

func (storeRequestExecutor *StoreRequestExecutor) setChannel(ch chan *protocol.RPC) {
	storeRequestExecutor.ch = ch
}

func (storeRequestExecutor *StoreRequestExecutor) setMessageId(id messageID) {
	storeRequestExecutor.id = id
}
