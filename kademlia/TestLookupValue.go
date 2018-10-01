package kademlia

import (
	log "github.com/sirupsen/logrus"
	"./protocol"
	"os"
)

type TestLookupValue struct {
	fileHash []byte
}

func NewTestLookupValue(fileHash []byte) *TestLookupValue {
	testLookupValue := &TestLookupValue{}
	testLookupValue.fileHash = fileHash
	return testLookupValue
}

func (testLookupValue *TestLookupValue) StartTest() {
	lookupValue := NewLookupValue(testLookupValue.fileHash, testLookupValue)
	lookupValue.start()
}

func (testlookupValue *TestLookupValue) contactWithFile(contact Contact,findValueRpc *protocol.FindValue) {
	log.WithFields(log.Fields{
		"Contact": contact,
		"find value":findValueRpc,
	}).Info("Test find value : contactWithFile")

	network:=Network{}
	portStr := os.Getenv("FILE_TRANSFER_PORT")
	originalSender := MyRoutingTable.me.ID[:]
	network.retrieveFile(portStr,testlookupValue.fileHash,findValueRpc.FileName,&contact,findValueRpc.FileSize,&originalSender)
}

func (testlookupValue *TestLookupValue) noContactWithFileFound(contacts []Contact) {
	log.WithFields(log.Fields{
	}).Info("Test find value : noContactWithFileFound")
}
