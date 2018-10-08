package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
)

func SendAndReceiveFindNode(callback FindNodeRequestCallback, target *KademliaID, contact Contact) {
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Info("Sending FindNode message to node.")
	findNodeRequestExecutor := FindNodeRequestExecutor{}
	findNodeRequestExecutor.contact = contact
	findNodeRequestExecutor.target = target
	findNodeRequestExecutor.callback = &callback
	createRoutine(&findNodeRequestExecutor)
}

func parseFindNodeRequest(rpc *protocol.RPC) *protocol.FindNode {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_FIND_NODE {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. FIND_NODE expected.")
	}

	findNode := &protocol.FindNode{}
	if err := proto.Unmarshal(rpc.Message, findNode); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming FindeNode message.")
	}

	return findNode
}

func answerFindNodeRequest(msg *protocol.RPC) {

	// Create message ID
	id := messageID{}
	copy(id[:], msg.MessageID[0:20])

	// Parse findNode message and create contact
	findNode := parseFindNodeRequest(msg)
	sender := createContactFromRPC(msg)
	targetId := KademliaIDFromSlice(findNode.TargetID)
	contacts := MyRoutingTable.FindClosestContacts(targetId, bucketSize)
	net := &Network{}
	originalSender := KademliaIDFromSlice(msg.OriginalSender)

	net.SendFindContactMessage(targetId, originalSender, sender, contacts, id)

	MyRoutingTable.AddContactAsync(*sender)
}

func contactToFindNode_Contact(contacts []Contact) []*protocol.FindNode_Contact {
	findNodeContacts := make([]*protocol.FindNode_Contact, 0)
	for _, contact := range contacts {
		protocolFindNodeContact := &protocol.FindNode_Contact{}
		protocolFindNodeContact.KademliaID = contact.ID[:]
		protocolFindNodeContact.IPaddress = contact.Address
		findNodeContacts = append(findNodeContacts, protocolFindNodeContact)
	}

	return findNodeContacts
}

func FindNode_ContactToContact(findNodeContacts []*protocol.FindNode_Contact) []Contact {
	contacts := make([]Contact, 0)
	for _, findNodeContact := range findNodeContacts {
		contact := Contact{}
		contact.ID = KademliaIDFromSlice(findNodeContact.KademliaID)
		contact.Address = findNodeContact.IPaddress
		contacts = append(contacts, contact)
	}

	return contacts
}

func createFindNodeToByte(targetId *KademliaID, contacts []Contact) ([]byte, error) {
	findNode := &protocol.FindNode{}
	findNode.TargetID = targetId[:]
	findNode.Contacts = contactToFindNode_Contact(contacts)
	return proto.Marshal(findNode)
}
