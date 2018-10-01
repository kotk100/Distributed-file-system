package kademlia

import (
	log "github.com/sirupsen/logrus"
	"sync"
	"./protocol"
)

//TODO : When an iterativeFindValue succeeds, the initiator must store the key/value pair at the closest node seen which did not return the value.

type LookupValueCallback interface {
	contactWithFile(contact Contact,findValueRpc *protocol.FindValue)
	noContactWithFileFound(contacts []Contact)
}

type FindValueRequestCallback interface {
	errorRequest(contact Contact)
	successRequest(contact Contact, contacts []Contact,findValueRpc *protocol.FindValue)
}

type LookupValue struct {
	lookupNode          *LookupNode
	fileHash            []byte
	lookupValueCallback *LookupValueCallback
	contactWithFile     Contact
	hasBeenFound        bool
	muxSuccessRequest   sync.Mutex
}

func NewLookupValue(fileHash []byte, lookupValueCallback LookupValueCallback) *LookupValue {
	lookup := &LookupValue{}
	lookup.fileHash = fileHash
	lookup.lookupValueCallback = &lookupValueCallback
	lookup.hasBeenFound = false
	contact := NewContact(KademliaIDFromSlice(fileHash), "")
	lookup.lookupNode = NewLookupNodeWithCustomSender(contact, lookup, lookup)
	return lookup
}

func (lookupValue *LookupValue) start(){
	lookupValue.lookupNode.Start()
}

func (lookupValue *LookupValue) errorRequest(contact Contact) {
	lookupValue.lookupNode.errorRequest(contact)
}

func (lookupValue *LookupValue) successRequest(contact Contact, contacts []Contact, findValueRpc *protocol.FindValue) {
	lookupValue.muxSuccessRequest.Lock()
	if ! lookupValue.hasBeenFound {
		if findValueRpc.HaveTheFile {
			log.WithFields(log.Fields{
				"Contact":  contact,
			}).Info("Found contact which contains the file")
			lookupValue.lookupNode.stop()
			lookupValue.hasBeenFound = true
			(*lookupValue.lookupValueCallback).contactWithFile(contact,findValueRpc)
		} else {
			lookupValue.lookupNode.successRequest(contact, contacts)
		}
	}
	lookupValue.muxSuccessRequest.Unlock()
}

func (lookupValue *LookupValue) processKClosest(KClosestOfTarget []LookupNodeContact) {
	contacts := make([]Contact,0)
	for _,v := range KClosestOfTarget{
		contacts=append(contacts,v.contact)
	}
	(*lookupValue.lookupValueCallback).noContactWithFileFound(contacts)
}

func (lookupValue *LookupValue) sendLookNode(target *KademliaID, contact *Contact) {
	log.WithFields(log.Fields{
		"lookupValue.fileHash":  lookupValue.fileHash,
		"lookupValue.fileHash string ":hashToString(lookupValue.fileHash),
	}).Info("Send find value")
	SendAndReceiveFindValue(lookupValue, lookupValue.fileHash, contact)
}
