package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"os"
)

type LookupValueManager struct {
	fileHash []byte
}

func NewLookupValueManager(fileHash []byte) *LookupValueManager {
	testLookupValue := &LookupValueManager{}
	testLookupValue.fileHash = fileHash
	return testLookupValue
}

func (lookupValueManager *LookupValueManager) FindValue() {
	log.WithFields(log.Fields{
		"file hash": hashToString(lookupValueManager.fileHash),
	}).Info("Start find value")
	lookupValue := NewLookupValue(lookupValueManager.fileHash, lookupValueManager)
	lookupValue.start()
}

func (lookupValueManager *LookupValueManager) contactWithFile(contact Contact, findValueRpc *protocol.FindValue, contacts []Contact) {
	log.WithFields(log.Fields{
		"Contact":    contact,
		"find value": findValueRpc,
	}).Info("Test find value : contactWithFile")

	network := Network{}
	portStr := os.Getenv("FILE_TRANSFER_PORT")
	originalSender := MyRoutingTable.me.ID[:]
	if len(contacts) > 0 && !contact.ID.Equals(contacts[0].ID) {
		network.retrieveFile(portStr, lookupValueManager.fileHash, findValueRpc.FileName, &contact, findValueRpc.FileSize, &originalSender, &contacts[0])
	} else {
		network.retrieveFile(portStr, lookupValueManager.fileHash, findValueRpc.FileName, &contact, findValueRpc.FileSize, &originalSender, nil)

	}
}

func (lookupValueManager *LookupValueManager) noContactWithFileFound(contacts []Contact) {
	log.Info("Test find value : noContactWithFileFound")
}

func (lookupValueManager *LookupValueManager) fileContents(fileContents []byte) {
	log.WithFields(log.Fields{
		"value": string(fileContents),
	}).Info("RECEIVED FILE VALUE (saved on the current node).------")
}
