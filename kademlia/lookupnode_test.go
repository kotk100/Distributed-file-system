package kademlia

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)



///////////UTILS
func fillShortList(lookupNode *LookupNode){
	contact:=NewContact(NewKademliaID("1111111101000000000000000000000111000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000110000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111111100000000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101010000000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101001000000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000100000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000010000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000100000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000010000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000001000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000000100000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000000010000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000000001000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000000000100000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000000000000010000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000011000000000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000001100000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000001000000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000100000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
	contact=NewContact(NewKademliaID("1111111101000000000000010000000000000000"), "localhost:8000")
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.updateShortlistIfCloser(&contact)
}

func setStateForAllContact(lookupNode *LookupNode,state LookupNodeContactState){
	for i := range lookupNode.shortlist {
		lookupNode.shortlist[i].lookupState = state
	}
}

//////////////////TESTS

func TestFirstContactIsInShortList(t *testing.T) {
	lookup := NewLookupNode(NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000"),nil)
	contact:=NewContact(NewKademliaID("1111111101000000000000000000000000000000"), "localhost:8000")
	lookup.shortlist=append(lookup.shortlist,NewLookupNodeContact(contact))
	assert.Equal(t,true, lookup.contactIsInShortlist(&contact), "Contact should be in the shortlist")
}

func TestContactIsNotInEmptyShortList(t *testing.T) {
	lookup := NewLookupNode(NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000"),nil)
	contact:=NewContact(NewKademliaID("1111111101000000000000000000000000000000"), "localhost:8000")
	assert.Equal(t,false, lookup.contactIsInShortlist(&contact), "Contact should not be in the shortlist")

}

func TestContactIsInFullShortList(t *testing.T) {
	lookup := NewLookupNode(NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000"),nil)
	fillShortList(lookup)
	contact:=NewContact(NewKademliaID("1111111101000000000000000100000000000000"), "localhost:8000")
	assert.Equal(t,true, lookup.contactIsInShortlist(&contact), "Contact should be in the shortlist")
}

func TestContactIsNotInFullShortList(t *testing.T) {
	lookup := NewLookupNode(NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000"),nil)
	fillShortList(lookup)
	contact:=NewContact(NewKademliaID("1111111101000000000001111100000000000000"), "localhost:8000")
	assert.Equal(t,false, lookup.contactIsInShortlist(&contact), "Contact should be in the shortlist")
}

func TestShortlistIsOrdered(t *testing.T) {
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	assert.Equal(t,20, len(lookup.shortlist), "Shortlist size should be 20")
	if !(lookup.shortlist[0].contact.distance.Less(lookup.shortlist[1].contact.distance) &&
		lookup.shortlist[1].contact.distance.Less(lookup.shortlist[2].contact.distance) &&
		lookup.shortlist[2].contact.distance.Less(lookup.shortlist[3].contact.distance) &&
		lookup.shortlist[3].contact.distance.Less(lookup.shortlist[4].contact.distance) &&
		lookup.shortlist[4].contact.distance.Less(lookup.shortlist[5].contact.distance) &&
		lookup.shortlist[5].contact.distance.Less(lookup.shortlist[6].contact.distance) &&
		lookup.shortlist[6].contact.distance.Less(lookup.shortlist[7].contact.distance) &&
		lookup.shortlist[7].contact.distance.Less(lookup.shortlist[8].contact.distance) &&
		lookup.shortlist[8].contact.distance.Less(lookup.shortlist[9].contact.distance) &&
		lookup.shortlist[9].contact.distance.Less(lookup.shortlist[10].contact.distance) &&
		lookup.shortlist[10].contact.distance.Less(lookup.shortlist[11].contact.distance) &&
		lookup.shortlist[11].contact.distance.Less(lookup.shortlist[12].contact.distance) &&
		lookup.shortlist[12].contact.distance.Less(lookup.shortlist[13].contact.distance) &&
		lookup.shortlist[13].contact.distance.Less(lookup.shortlist[14].contact.distance) &&
		lookup.shortlist[14].contact.distance.Less(lookup.shortlist[15].contact.distance) &&
		lookup.shortlist[15].contact.distance.Less(lookup.shortlist[16].contact.distance) &&
		lookup.shortlist[16].contact.distance.Less(lookup.shortlist[17].contact.distance) &&
		lookup.shortlist[17].contact.distance.Less(lookup.shortlist[18].contact.distance) &&
		lookup.shortlist[18].contact.distance.Less(lookup.shortlist[19].contact.distance)){
		t.Errorf("shortlist should be ordered")
	}
}

func TestAllKClosestHasBeenFound(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	setStateForAllContact(lookup,ASKED)
	assert.Equal(t,true, lookup.isKClosestContactHasBeenFound(), "Should have found k closest")
}

func TestChangeContactStatus(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	contact:=&lookup.shortlist[3]
	lookup.changeLookupNodeContactState(contact,ASKING)
	assert.Equal(t,ASKING, lookup.shortlist[3].lookupState, "Should have ASKING state")
}


func TestUpdateShortList(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	newContact:=NewContact(NewKademliaID("1111111111000000000000000000000000000001"), "localhost:8000")
	newContact.distance = newContact.ID.CalcDistance(target.ID)
	lookup.updateShortlistIfCloser(&newContact)
	assert.Equal(t,newContact.ID, lookup.shortlist[0].contact.ID, "Should have same ID")
}

func TestUpdateLastContactShortList(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	newContact:=NewContact(NewKademliaID("1011111111111111000000000000000000000000"), "localhost:8000")
	newContact.distance = newContact.ID.CalcDistance(target.ID)
	lookup.updateShortlistIfCloser(&newContact)
	assert.Equal(t,20, len(lookup.shortlist), "Size should be 20")
	assert.Equal(t,newContact.ID, lookup.shortlist[19].contact.ID, "Should have same ID")
}
func TestUpdateMiddleContactShortList(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	newContact:=NewContact(NewKademliaID("1111111101000110000000000000000000000000"), "localhost:8000")
	newContact.distance = newContact.ID.CalcDistance(target.ID)
	lookup.updateShortlistIfCloser(&newContact)
	fmt.Println(lookup.shortlist)
	assert.Equal(t,20, len(lookup.shortlist), "Size should be 20")
	assert.Equal(t,newContact.ID, lookup.shortlist[18].contact.ID, "Should have same ID")
}

func TestCleanShortList(t *testing.T){
	target:=NewContact(NewKademliaID("1111111111000000000000000000000000000000"), "localhost:8000")
	lookup := NewLookupNode(target,nil)
	fillShortList(lookup)
	contact:=lookup.shortlist[3]
	lookup.errorRequest(contact.contact)
	lookup.changeLookupNodeContactState(&contact,ASKING)
	assert.Equal(t,19, len(lookup.shortlist), "Size should be 20")
	assert.Equal(t,false, lookup.contactIsInShortlist(&contact.contact), "Contact should be removed")
}
