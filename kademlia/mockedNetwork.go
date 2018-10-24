package kademlia

import "github.com/stretchr/testify/mock"

type MockedNetwork struct {
	mock.Mock
}

func (mockedNetwork *MockedNetwork) SendPingMessage(originalSender *KademliaID, contact *Contact, messageID messageID) bool {
	args := mockedNetwork.Called(originalSender, contact, messageID)
	return args.Bool(0)
}

func (mockedNetwork *MockedNetwork) SendFindContactMessage(targetId *KademliaID, originalSender *KademliaID, contact *Contact, contacts []Contact, messageID messageID) bool {
	args := mockedNetwork.Called(targetId, originalSender, contact, contacts, messageID)
	return args.Bool(0)
}

func (mockedNetwork *MockedNetwork) SendFindDataMessage(fileHash []byte, contact *Contact, contacts []Contact, messageID messageID, originalSender *KademliaID, haveTheFile bool, fileName string, fileSize int64) bool {
	args := mockedNetwork.Called(fileHash, contact, contacts, messageID, originalSender, haveTheFile, fileName, fileSize)
	return args.Bool(0)
}
