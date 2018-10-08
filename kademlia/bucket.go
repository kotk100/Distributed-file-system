package kademlia

import (
	"container/list"
	log "github.com/sirupsen/logrus"
	"sync"
)

// bucket definition
// contains a List
type bucket struct {
	list            *list.List
	muxAdd          sync.Mutex
	muxAccessBucket sync.Mutex
	contactToInsert Contact
	networkAPI      NetworkAPI
}

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.networkAPI = &Network{}
	return bucket
}

// AddContactAsync adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	if MyRoutingTable.GetMe().ID.Equals(contact.ID) {
		return
	}

	bucket.muxAdd.Lock()
	bucket.muxAccessBucket.Lock()

	log.WithFields(log.Fields{
		"Contact": contact,
	}).Debug("Updating bucket.")
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

			bucket.muxAccessBucket.Unlock()
			bucket.muxAdd.Unlock()
		} else {
			bucket.contactToInsert = contact
			contact := bucket.list.Back().Value.(Contact)
			bucket.muxAccessBucket.Unlock()
			pingBucketRequestExecutor := PingBucketRequestExecutor{}
			pingBucketRequestExecutor.contact = contact
			pingBucketRequestExecutor.bucket = bucket
			createRoutine(&pingBucketRequestExecutor)
		}
	} else {
		bucket.list.MoveToFront(element)

		bucket.muxAccessBucket.Unlock()
		bucket.muxAdd.Unlock()
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact
	bucket.muxAccessBucket.Lock()
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	bucket.muxAccessBucket.Unlock()
	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}

func (bucket *bucket) print() {
	bucket.muxAccessBucket.Lock()
	log.WithFields(log.Fields{
		"contents": bucket.list,
	}).Info("")
	bucket.muxAccessBucket.Unlock()
}
