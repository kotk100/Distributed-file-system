package kademlia

import (
	"container/list"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	syncSem "golang.org/x/sync/semaphore"
)

// bucket definition
// contains a List
type bucket struct {
	list           *list.List
	readers        *syncSem.Weighted
	contactupdates *syncSem.Weighted
	bucketsLocked  int
	networkAPI     NetworkAPI
}

const allowedReaders = 100

// newBucket returns a new instance of a bucket
func newBucket() *bucket {
	bucket := &bucket{}
	bucket.list = list.New()
	bucket.networkAPI = &Network{}
	bucket.readers = syncSem.NewWeighted(allowedReaders)
	bucket.contactupdates = syncSem.NewWeighted(bucketSize)
	bucket.bucketsLocked = 0
	return bucket
}

// AddContactAsync adds the Contact to the front of the bucket
// or moves it to the front of the bucket if it already existed
func (bucket *bucket) AddContact(contact Contact) {
	if MyRoutingTable.GetMe().ID.Equals(contact.ID) {
		return
	}
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Debug("Updating bucket.")

	// Find if concact exsits
	bucket.readers.Acquire(context.TODO(), allowedReaders)
	var element *list.Element
	elementPosition := 0 //From back of list
	for e := bucket.list.Back(); e != nil; e = e.Prev() {
		if (contact).ID.Equals(e.Value.(Contact).ID) {
			element = e
			break
		}
		elementPosition += 1
	}

	if element == nil {
		if bucket.list.Len() < bucketSize {
			//add to the front of the list
			bucket.list.PushFront(contact)
			bucket.readers.Release(allowedReaders)
		} else {
			bucket.readers.Release(allowedReaders)
			pingBucketRequestExecutor := PingBucketRequestExecutor{}
			pingBucketRequestExecutor.contactToAdd = contact
			pingBucketRequestExecutor.bucket = bucket
			createRoutine(&pingBucketRequestExecutor)
		}
	} else {
		// If bucket is locked because we are waiting for ping do nothing as the other routine will either move it to the front or replace it if it fails to respond otherwise move it to front
		if bucket.bucketsLocked <= elementPosition {
			bucket.list.MoveToFront(element)
		}
		bucket.readers.Release(allowedReaders)
	}
}

// GetContactAndCalcDistance returns an array of Contacts where
// the distance has already been calculated
func (bucket *bucket) GetContactAndCalcDistance(target *KademliaID) []Contact {
	var contacts []Contact
	bucket.readers.Acquire(context.TODO(), 1)
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contact.CalcDistance(target)
		contacts = append(contacts, contact)
	}
	bucket.readers.Release(1)
	return contacts
}

// Len return the size of the bucket
func (bucket *bucket) Len() int {
	return bucket.list.Len()
}

func (bucket *bucket) GetContacts() *[]Contact {
	contacts := []Contact{}
	bucket.readers.Acquire(context.TODO(), 1)
	for elt := bucket.list.Front(); elt != nil; elt = elt.Next() {
		contact := elt.Value.(Contact)
		contacts = append(contacts, contact)
	}
	bucket.readers.Release(1)
	return &contacts
}

func (bucket *bucket) print() {
	bucket.readers.Acquire(context.TODO(), 1)
	log.WithFields(log.Fields{
		"contents": bucket.list,
	}).Info("")
	bucket.readers.Release(1)
}

func (bucket *bucket) printContacts() {
	bucket.readers.Acquire(context.TODO(), 1)
	i := 0
	for e := bucket.list.Front(); e != nil; e = e.Next() {
		log.WithFields(log.Fields{
			"Contact number": i,
			"contents":       e,
		}).Info("")
		i++
	}
	bucket.readers.Release(1)
}
