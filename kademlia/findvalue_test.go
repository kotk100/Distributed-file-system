package kademlia

import (
	"./protocol"
	"github.com/stretchr/testify/assert"
	"testing"
)

func getContactFindValueSlice() []Contact {
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

func TestConvertFindValueMessageWithContactsAndRetrieveIt(t *testing.T) {
	InitMyInformation("test")

	contacts := getContactFindValueSlice()
	fileName := "file test"
	haveTheFile := false
	fileHash := []byte("my hash")
	out, _ := createFindValueToByte(fileHash, contacts, haveTheFile, fileName, 13)

	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_VALUE
	rpc.Message = out

	findValue := parseFindValueRequest(&rpc)

	assert.Equal(t, false, findValue.HaveTheFile, "Should not have the file")
	assert.Equal(t, "file test", findValue.FileName, "File name should be the same")
	assert.Equal(t, int64(13), findValue.FileSize, "File size should be the same")
	assert.Equal(t, fileHash, findValue.FileHash, "File hash should be the same")

	contactsRevert := FindValue_ContactToContact(findValue.Contacts)

	assert.Equal(t, 20, len(contactsRevert), "Size should be 20")
	assert.Equal(t, NewKademliaID("FFFFFFFF00000000000000000000000000000010"), contactsRevert[0].ID, "Should be same ID")
	assert.Equal(t, NewKademliaID("FFFFFFFF00000000000100000000000000000000"), contactsRevert[19].ID, "Should be same ID")

}

func TestConvertFindValueMessageWithoutContactsAndRetrieveIt(t *testing.T) {
	InitMyInformation("test")

	contacts := make([]Contact,0)
	fileName := "file test"
	haveTheFile := false
	fileHash := []byte("my hash")
	out, _ := createFindValueToByte(fileHash, contacts, haveTheFile, fileName, 13)

	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_VALUE
	rpc.Message = out

	findValue := parseFindValueRequest(&rpc)

	assert.Equal(t, false, findValue.HaveTheFile, "Should not have the file")
	assert.Equal(t, "file test", findValue.FileName, "File name should be the same")
	assert.Equal(t, int64(13), findValue.FileSize, "File size should be the same")
	assert.Equal(t, fileHash, findValue.FileHash, "File hash should be the same")

	contactsRevert := FindValue_ContactToContact(findValue.Contacts)

	assert.Equal(t, 0, len(contactsRevert), "Size should be 0")

}

func TestConvertFindValueMessageWhereFileHasBeenFoundAndRetrieveIt(t *testing.T){

	contacts := make([]Contact,0)
	fileName := "file test"
	haveTheFile := true
	fileHash := []byte("my hash")
	out, _ := createFindValueToByte(fileHash, contacts, haveTheFile, fileName, 13)

	rpc := protocol.RPC{}
	rpc.MessageType = protocol.RPC_FIND_VALUE
	rpc.Message = out

	findValue := parseFindValueRequest(&rpc)

	assert.Equal(t, true, findValue.HaveTheFile, "Should have the file")
	assert.Equal(t, "file test", findValue.FileName, "File name should be the same")
	assert.Equal(t, int64(13), findValue.FileSize, "File size should be the same")
	assert.Equal(t, fileHash, findValue.FileHash, "File hash should be the same")

	contactsRevert := FindValue_ContactToContact(findValue.Contacts)

	assert.Equal(t, 0, len(contactsRevert), "Size should be 0")
}

