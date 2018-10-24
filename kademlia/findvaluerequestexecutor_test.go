package kademlia

import (
	"./protocol"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
)

type MockedCallbackFVR struct {
	mock.Mock
}

func (mockedCallback *MockedCallbackFVR) successRequest(contact Contact, contacts []Contact, findValueRpc *protocol.FindValue) {
	mockedCallback.Called(contact, contacts, findValueRpc)
}

func (mockedCallback *MockedCallbackFVR) errorRequest(contact Contact) {
	mockedCallback.Called(contact)
}

func linkCallbackFVR(executor *FindValueRequestExecutor, findNodeRequestCallback FindValueRequestCallback) {
	executor.callback = &findNodeRequestCallback
}

func TestExecutionErrorFVR(t *testing.T) {
	InitMyInformation("test")
	executor := &FindValueRequestExecutor{}

	ch := make(chan *protocol.RPC)
	executor.ch = ch

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	executor.id = messageID

	executor.contact = Contact{}

	executor.fileHash = make([]byte, 0)

	callback := &MockedCallbackFVR{}
	callback.On("errorRequest", executor.contact)
	linkCallbackFVR(executor, callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindDataMessage", executor.fileHash, &executor.contact, make([]Contact, 0), executor.id, MyRoutingTable.me.ID, false,"", int64(0)).Return(true)
	executor.networkAPI = testObj

	executor.execute()

	callback.AssertCalled(t, "errorRequest", executor.contact)

}

func TestExecutionTimeoutFVR(t *testing.T) {
	InitMyInformation("test")
	executor := &FindValueRequestExecutor{}

	ch := make(chan *protocol.RPC)
	executor.ch = ch

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	executor.id = messageID

	executor.contact = Contact{}

	executor.fileHash = make([]byte, 0)

	callback := &MockedCallbackFVR{}
	callback.On("errorRequest", executor.contact)
	linkCallbackFVR(executor, callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindDataMessage", executor.fileHash, &executor.contact, make([]Contact, 0), executor.id, MyRoutingTable.me.ID, false, "", int64(0)).Return(false)
	executor.networkAPI = testObj

	executor.execute()

	time.Sleep(1 * time.Second)
	callback.AssertCalled(t, "errorRequest", executor.contact)

}

func TestSuccesFVRs(t *testing.T) {
	InitMyInformation("test")
	executor := &FindValueRequestExecutor{}

	ch := make(chan *protocol.RPC)
	executor.ch = ch

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	executor.id = messageID

	executor.contact = Contact{}
	executor.contact.ID = NewRandomKademliaID()
	executor.contact.Address = "0.0.0.0"

	executor.fileHash = make([]byte, 0)

	contacts := make([]Contact,0)
	fileName := "file test"
	haveTheFile := true
	fileHash := []byte("my hash")
	out, _ := createFindValueToByte(fileHash, contacts, haveTheFile, fileName, 13)
	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_VALUE
	rpc.Message = out
	rpc.KademliaID = executor.contact.ID[:]
	rpc.IPaddress = executor.contact.Address
	findValue := parseFindValueRequest(&rpc)

	callback := &MockedCallbackFVR{}
	callback.On("successRequest", executor.contact, make([]Contact, 0),findValue)
	linkCallbackFVR(executor, callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindDataMessage", executor.fileHash, &executor.contact, make([]Contact, 0), executor.id, MyRoutingTable.me.ID, false, "", int64(0)).Return(false)
	executor.networkAPI = testObj

	go executor.execute()
	time.Sleep(1 * time.Second)

	ch <- &rpc
	time.Sleep(2 * time.Second)
	callback.AssertCalled(t, "successRequest", executor.contact, make([]Contact, 0),findValue)

}
