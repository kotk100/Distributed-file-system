package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

type LookupValueCallback interface {
	contactWithFile(contact Contact, findValueRpc *protocol.FindValue, contacts []LookupNodeContact)
	fileContents(fileContents []byte,stringPath string)
	noContactWithFileFound(contacts []Contact)
}

type FindValueRequestCallback interface {
	errorRequest(contact Contact)
	successRequest(contact Contact, contacts []Contact, findValueRpc *protocol.FindValue)
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

func (lookupValue *LookupValue) start() {
	if checkFileExistsHash(hashToString(lookupValue.fileHash)) {
		stringPath := getPathOfFileFromHash(hashToString(lookupValue.fileHash))
		file, error := os.Open(stringPath)
		if error != nil {
			log.WithFields(log.Fields{
				"file string": stringPath,
			}).Info("Error open file")
		} else {
			defer file.Close()
			fileInfo, error := file.Stat()
			if error != nil {
				log.WithFields(log.Fields{
					"file string": stringPath,
				}).Info("Error to get file info")
			} else {
				fileContents := make([]byte, fileInfo.Size())
				_, err := file.Read(fileContents)
				if err != nil {
					log.WithFields(log.Fields{
						"file string": stringPath,
					}).Info("Error to get file contents")
				} else {
					(*lookupValue.lookupValueCallback).fileContents(fileContents,stringPath)
				}
			}
		}
	} else {
		lookupValue.lookupNode.Start()
	}
}

func (lookupValue *LookupValue) errorRequest(contact Contact) {
	lookupValue.lookupNode.errorRequest(contact)
}

func (lookupValue *LookupValue) successRequest(contact Contact, contacts []Contact, findValueRpc *protocol.FindValue) {
	lookupValue.muxSuccessRequest.Lock()
	if !lookupValue.hasBeenFound {
		if findValueRpc.HaveTheFile {
			log.WithFields(log.Fields{
				"Contact": contact,
			}).Info("Found contact which contains the file")
			lookupValue.lookupNode.stop()
			lookupValue.hasBeenFound = true
			(*lookupValue.lookupValueCallback).contactWithFile(contact, findValueRpc, lookupValue.lookupNode.shortlist)
		} else {
			lookupValue.lookupNode.successRequest(contact, contacts)
		}
	}
	lookupValue.muxSuccessRequest.Unlock()
}

func (lookupValue *LookupValue) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("lookup value no node with file found")
	contacts := make([]Contact, 0)
	for _, v := range KClosestOfTarget {
		contacts = append(contacts, v.contact)
	}
	(*lookupValue.lookupValueCallback).noContactWithFileFound(contacts)
}

func (lookupValue *LookupValue) sendLookNode(target *KademliaID, contact *Contact) {
	log.WithFields(log.Fields{
		"lookupValue.fileHash":         lookupValue.fileHash,
		"lookupValue.fileHash string ": hashToString(lookupValue.fileHash),
	}).Info("Send find value")
	SendAndReceiveFindValue(lookupValue, lookupValue.fileHash, *contact)
}
