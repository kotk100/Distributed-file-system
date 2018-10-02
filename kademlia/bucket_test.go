package kademlia

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"math/rand"
	"testing"
	"time"
	"./protocol"
)

type MockedNetwork struct{
	mock.Mock
}

func (mockedNetwork *MockedNetwork) SendPingMessage(originalSender *KademliaID,contact *Contact, messageID messageID) bool {
	args := mockedNetwork.Called(originalSender ,contact, messageID )
	return args.Bool(0)
}

func containContact(bucket *bucket,contact *Contact)  bool{
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID
		if (contact).ID.Equals(nodeID) {
			return true
		}
	}
	return false
}

func fillUpBucket(bucket *bucket){
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000000003420000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001036400000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000340000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000003700001000000000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000340000001000000000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00003000000001002300000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000000001100000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000020001002430000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000001110000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000200001534230000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000300200000030000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000200000000004000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001003000000000002300"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF36530000000001000450000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000746000001000340000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000000004003400000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000001000000003423000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00200000000001000000000000007650"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000034001000000000000000000"), "localhost:8001"))
	bucket.AddContact(NewContact(NewKademliaID("FFFFFFFF00222000000001000000000000000000"), "localhost:8001"))
}

func TestBucketIsEmptyAtTheBeginning(t *testing.T){
	bucket:=newBucket()
	assert.Equal(t, 0, bucket.Len(), "size should be 0 at the beginning")
}

func TestBucketShouldAddAContactWhenIsNotFull(t *testing.T){
	InitMyInformation("test")
	bucket:=newBucket()
	contact := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001")
	bucket.AddContact(contact)
	assert.Equal(t, 1, bucket.Len(), "size should be 1 after adding a contact")
	assert.Equal(t, true, containContact(bucket,&contact), "contact should be in the list")
}

func TestContactShouldMoveToFront(t *testing.T){
	InitMyInformation("test")
	bucket:=newBucket()
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000001000000000000000000"), "localhost:8001")
	bucket.AddContact(contact1)
	contact2 := NewContact(NewKademliaID("FFFFFFFF00000000000000100000000000000000"), "localhost:8001")
	bucket.AddContact(contact2)
	contact3 := NewContact(NewKademliaID("FFFFFFFF00000100000000000000000000000000"), "localhost:8001")
	bucket.AddContact(contact3)
	assert.Equal(t, 3, bucket.Len(), "size should be 3")
	assert.Equal(t, contact3.ID, bucket.list.Front().Value.(Contact).ID, "ID should be the same")
	bucket.AddContact(contact1)
	assert.Equal(t, 3, bucket.Len(), "size should be 3")
	assert.Equal(t, contact1.ID, bucket.list.Front().Value.(Contact).ID, "ID should be the same")
}

func TestInsertContactWhenBucketIsFullAndOldestContactTimeOut(t *testing.T){
	InitMyInformation("test")
	bucket:=newBucket()

	ch := make(chan *protocol.RPC)

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}

	fillUpBucket(bucket)

	contactToInsert := NewContact(NewKademliaID("FFFFFFFF00000100000000000000000000001234"), "localhost:8001")
	oldestContact:=NewContact(NewKademliaID("FFFFFFFF00000000000001000000003420000000"), "localhost:8001")

	testObj := new(MockedNetwork)
	testObj.On("SendPingMessage", MyRoutingTable.me.ID,&oldestContact,messageID).Return(false)
	bucket.networkAPI = testObj

	bucket.contactToInsert = contactToInsert
	bucket.muxAdd.Lock()

	pingBucketRequestExecutor:= PingBucketRequestExecutor{}
	pingBucketRequestExecutor.contact = &oldestContact
	pingBucketRequestExecutor.bucket = bucket
	pingBucketRequestExecutor.id = messageID
	pingBucketRequestExecutor.ch=ch
	go pingBucketRequestExecutor.execute()

	time.Sleep(11 * time.Second)

	assert.Equal(t, 20, bucket.Len(), "size should be 20")
	assert.Equal(t, contactToInsert.ID, bucket.list.Front().Value.(Contact).ID, "ID should be the same")
}


func TestInsertContactWhenBucketIsFullAndErrorToSendPing(t *testing.T){
	InitMyInformation("test")
	bucket:=newBucket()

	ch := make(chan *protocol.RPC)

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}

	fillUpBucket(bucket)

	contactToInsert := NewContact(NewKademliaID("FFFFFFFF00000100000000000000000000001234"), "localhost:8001")
	oldestContact:=NewContact(NewKademliaID("FFFFFFFF00000000000001000000003420000000"), "localhost:8001")

	testObj := new(MockedNetwork)
	testObj.On("SendPingMessage", MyRoutingTable.me.ID,&oldestContact,messageID).Return(true)
	bucket.networkAPI = testObj

	bucket.contactToInsert = contactToInsert
	bucket.muxAdd.Lock()
	pingBucketRequestExecutor:= PingBucketRequestExecutor{}
	pingBucketRequestExecutor.contact = &oldestContact
	pingBucketRequestExecutor.bucket = bucket
	pingBucketRequestExecutor.id = messageID
	pingBucketRequestExecutor.ch=ch
	go pingBucketRequestExecutor.execute()

	time.Sleep(2 * time.Second)

	assert.Equal(t, 20, bucket.Len(), "size should be 20")
	assert.Equal(t, contactToInsert.ID, bucket.list.Front().Value.(Contact).ID, "ID should be the same")
}


func TestInsertContactWhenBucketIsFullAndOldestContactAnswer(t *testing.T){
	InitMyInformation("test")
	bucket:=newBucket()

	ch := make(chan *protocol.RPC)

	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}

	fillUpBucket(bucket)

	bucket.contactToInsert = NewContact(NewKademliaID("FFFFFFFF00000100000000000000000000001234"), "localhost:8001")
	oldestContact:=NewContact(NewKademliaID("FFFFFFFF00000000000001000000003420000000"), "localhost:8001")

	testObj := new(MockedNetwork)
	testObj.On("SendPingMessage", MyRoutingTable.me.ID,&oldestContact,messageID).Return(false)
	bucket.networkAPI = testObj

	bucket.muxAdd.Lock()
	pingBucketRequestExecutor:= PingBucketRequestExecutor{}
	pingBucketRequestExecutor.contact = &oldestContact
	pingBucketRequestExecutor.bucket = bucket
	pingBucketRequestExecutor.id = messageID
	pingBucketRequestExecutor.ch=ch
	go pingBucketRequestExecutor.execute()

	ch <-&protocol.RPC{}
	time.Sleep(2 * time.Second)

	assert.Equal(t, 20, bucket.Len(), "size should be 20")
	assert.Equal(t, oldestContact.ID, bucket.list.Front().Value.(Contact).ID, "ID should be the same")
}

