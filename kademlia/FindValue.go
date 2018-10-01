package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func SendAndReceiveFindValue(callback FindValueRequestCallback, fileHash []byte, contact *Contact) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending FindValue message to node.")
	findValueRequestExecutor := FindValueRequestExecutor{}
	findValueRequestExecutor.contact = contact
	findValueRequestExecutor.fileHash = fileHash
	findValueRequestExecutor.callback = &callback
	createRoutine(&findValueRequestExecutor)
}

func parseFindValueRequest(rpc *protocol.RPC) *protocol.FindValue {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_FIND_VALUE {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type received. FIND_VALUE expected.")
	}

	findValue := &protocol.FindValue{}
	if err := proto.Unmarshal(rpc.Message, findValue); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming FindValue message.")
	}

	return findValue
}


func answerFindValueRequest(msg *protocol.RPC) {

	// Create message ID
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])

	// Parse findNode message and create contact
	findValue := parseFindValueRequest(msg)
	sender := createContactFromFindValue(msg)
	fileHashKademlia := KademliaIDFromSlice(findValue.FileHash)
	net := &Network{}
	originalSender := KademliaIDFromSlice(msg.OriginalSender)

	//check if node has the file set boolean
	haveTheFile:=false
	contacts := MyRoutingTable.FindClosestContacts(fileHashKademlia, bucketSize)

	net.SendFindDataMessage(findValue.FileHash, sender, contacts,id, originalSender, haveTheFile)
	MyRoutingTable.AddContactAsync(*sender)
}

func createContactFromFindValue(rpc *protocol.RPC) *Contact {
	// Create contact
	contact := &Contact{}
	contact.Address = rpc.IPaddress
	contact.ID = KademliaIDFromSlice(rpc.KademliaID)
	return contact
}

func contactToFindValue_Contact(contacts []Contact) []*protocol.FindValue_Contact {
	findValueContacts := make([]*protocol.FindValue_Contact, 0)
	for _, contact := range contacts {
		protocolFindValueContact := &protocol.FindValue_Contact{}
		protocolFindValueContact.KademliaID = contact.ID[:]
		protocolFindValueContact.IPaddress = contact.Address
		findValueContacts = append(findValueContacts, protocolFindValueContact)
	}

	return findValueContacts
}

func FindValue_ContactToContact(findValueContacts []*protocol.FindValue_Contact) []Contact {
	contacts := make([]Contact, 0)
	for _, findValueContact := range findValueContacts {
		contact := Contact{}
		contact.ID = KademliaIDFromSlice(findValueContact.KademliaID)
		contact.Address = findValueContact.IPaddress
		contacts = append(contacts, contact)
	}
	return contacts
}

func createFindValueToByte(fileHash []byte, contacts []Contact, haveTheFile bool) ([]byte, error) {
	findValue := &protocol.FindValue{}
	findValue.FileHash = fileHash
	findValue.Contacts = contactToFindValue_Contact(contacts)
	findValue.HaveTheFile = haveTheFile
	return proto.Marshal(findValue)
}
