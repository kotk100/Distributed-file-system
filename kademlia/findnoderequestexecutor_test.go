package kademlia

import (
	"./protocol"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
)

type MockedCallback struct {
	mock.Mock
}

func (mockedCallback *MockedCallback) successRequest(contactAsked Contact, KClosestOfTarget []Contact) {
	mockedCallback.Called(contactAsked, KClosestOfTarget)
}

func (mockedCallback *MockedCallback) errorRequest(contactAsked Contact) {
	mockedCallback.Called(contactAsked)
}


func linkCallback(executor *FindNodeRequestExecutor,findNodeRequestCallback FindNodeRequestCallback ) {
	executor.callback = &findNodeRequestCallback
}

func TestExecutionError(t *testing.T) {
	InitMyInformation("test")
	executor := &FindNodeRequestExecutor{}

	ch := make(chan *protocol.RPC)
	executor.ch = ch

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	executor.id = messageID

	executor.contact = Contact{}

	executor.target = NewRandomKademliaID()

	callback := &MockedCallback{}
	callback.On("errorRequest",executor.contact)
	linkCallback(executor,callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindContactMessage", executor.target, MyRoutingTable.me.ID, &executor.contact, make([]Contact, 0), executor.id).Return(true)
	executor.networkAPI = testObj

	executor.execute()

	callback.AssertCalled(t,"errorRequest",executor.contact)

}

func TestExecutionTimeout(t *testing.T) {
	InitMyInformation("test")
	executor := &FindNodeRequestExecutor{}

	ch := make(chan *protocol.RPC)
	executor.ch = ch

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	executor.id = messageID

	executor.contact = Contact{}

	executor.target = NewRandomKademliaID()

	callback := &MockedCallback{}
	callback.On("errorRequest",executor.contact)
	linkCallback(executor,callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindContactMessage", executor.target, MyRoutingTable.me.ID, &executor.contact, make([]Contact, 0), executor.id).Return(false)
	executor.networkAPI = testObj

	executor.execute()

	time.Sleep(6 * time.Second)
	callback.AssertCalled(t,"errorRequest",executor.contact)

}

func TestSuccess(t *testing.T) {
	InitMyInformation("test")
	executor := &FindNodeRequestExecutor{}

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

	executor.target = NewRandomKademliaID()

	callback := &MockedCallback{}
	callback.On("successRequest",executor.contact,make([]Contact,0))
	linkCallback(executor,callback)

	testObj := &MockedNetwork{}
	testObj.On("SendFindContactMessage", executor.target, MyRoutingTable.me.ID, &executor.contact, make([]Contact, 0), executor.id).Return(false)
	executor.networkAPI = testObj

	go executor.execute()
	time.Sleep(1 * time.Second)

	contacts := make([]Contact,0)
	targetId := NewRandomKademliaID()
	out, _ := createFindNodeToByte(targetId, contacts)

	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_NODE
	rpc.Message = out
	rpc.KademliaID = executor.contact.ID[:]
	rpc.IPaddress = executor.contact.Address

	ch <- &rpc
	time.Sleep(2 * time.Second)
	callback.AssertCalled(t,"successRequest",executor.contact,make([]Contact,0))

}
