package kademlia

import (
	"./protocol"
	"container/list"
	"sync"
	log "github.com/sirupsen/logrus"
)

// bucket definition
// contains a List
type bucket struct {
	list *list.List
	mux sync.Mutex
	contactToInsert Contact
	networkAPI NetworkAPI
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.networkAPI = &Network{}
	return bucket
}

// AddContact adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	bucket.mux.Lock()
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Updating bucket.")
	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			bucket.list.PushFront(contact)
			bucket.mux.Unlock()
		} else{
			bucket.contactToInsert = contact
			contact:=bucket.list.Back().Value.(Contact)
			createRoutine(bucket.xxx_sendAndReceiveCheckBucketPing,&contact )
		}
	} else {
		bucket.list.MoveToFront(element)
		bucket.mux.Unlock()
	}
}

func (bucket *bucket) xxx_sendAndReceiveCheckBucketPing(ch chan *protocol.RPC, id messageID, contact *Contact){
	// Send ping message to other node
	error := bucket.networkAPI.SendPingMessage(MyRoutingTable.me.ID,contact, id)

	var element *list.Element
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (contact).ID.Equals(nodeID) {
			element = e
		}
	}
	// Receive response message through channel
	if error {
		log.Info("Error to send bucket ping.")
		bucket.list.Remove(element)
		bucket.list.PushFront(bucket.contactToInsert)
		bucket.mux.Unlock()
		destroyRoutine(id)
	}else{
		timeout:=NewTimeout(id,ch)
		timeOutManager.insertAndStart(timeout)
		rpc := <-ch
		timeOutManager.tryGetAndRemoveTimeOut(id)
		//timeout
		if element !=nil{
			if rpc == nil{
				log.Info("PING bucket time out.")
				bucket.list.Remove(element)
				bucket.list.PushFront(bucket.contactToInsert)
			} else{
				log.Info("Received PING bucket message response.")
				bucket.list.MoveToFront(element)
			}
		}
		bucket.mux.Unlock()
		destroyRoutine(id)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact

	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}

	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}
