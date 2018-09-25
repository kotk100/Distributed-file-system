package kademlia

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"./protocol"
)

func getContactSlice()[]Contact{
	return []Contact{
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000010"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000100"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000001000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000010000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000100000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000001000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000010000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000000100000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000001000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000010000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000000100000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000001000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000010000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000000100000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000001000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000010000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000000100000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000001000000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000010000000000000000000"), "localhost:8001"),
		NewContact(NewKademliaID("FFFFFFFF00000000000100000000000000000000"), "localhost:8001")}
}

func TestConvertFindNodeMessageAndRetrieveIt(t *testing.T) {
	InitMyInformation("test")

	contacts := getContactSlice()
	targetId := NewRandomKademliaID()
	out, _ := createFindNodeToByte(targetId, contacts)

	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_NODE
	rpc.Message = out

	findNode := parseFindNodeRequest(&rpc)

	assert.Equal(t, targetId[:], findNode.TargetID, "ID should be the same")

	contactsRevert := FindNode_ContactToContact(findNode.Contacts)

	assert.Equal(t, 20, len(contactsRevert), "Size should be 20")

	assert.Equal(t, NewKademliaID("FFFFFFFF00000000000000000000000000000010"), contactsRevert[0].ID, "Should be same ID")
	assert.Equal(t, NewKademliaID("FFFFFFFF00000000000100000000000000000000"), contactsRevert[19].ID, "Should be same ID")

}