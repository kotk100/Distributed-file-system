package kademlia

import (
	"sort"
	"sync"
)

type FindNodeRequestCallback interface {
	successRequest(contactAsked Contact,KClosestOfTarget []Contact)
	errorRequest(contactAsked Contact)
}

type LookupNode struct{
	target Contact
	shortlist []LookupNodeContact
	mux sync.Mutex
}

func NewLookupNode(target Contact) *LookupNode {
	lookup := &LookupNode{}
	lookup.shortlist=make([]LookupNodeContact,0)
	lookup.target = target
	return lookup
}

//give channel or callback to get K closest result
func (lookupNode *LookupNode) start(){
	contactsCloseOfTarget := MyRoutingTable.FindClosestContacts(lookupNode.target.ID,bucketSize)
	for _, contact := range contactsCloseOfTarget{
		lookupNode.shortlist = append(lookupNode.shortlist,NewLookupNodeContact(contact))
	}
	lookupNode.sendFindNodeRequest()
	lookupNode.sendFindNodeRequest()
	lookupNode.sendFindNodeRequest()
}

//callback for request executor
func (lookupNode *LookupNode) successRequest(contactAsked Contact,KClosestOfTarget []Contact){

}

func (lookupNode *LookupNode) errorRequest(contactAsked Contact){

}

func (lookupNode *LookupNode) sendFindNodeRequest(){
	lookupNode.mux.Lock()
	contactToAsk := lookupNode.getLookupNodeContactUnAsked()
	lookupNode.changeLookupNodeContactState(contactToAsk,ASKING)
	lookupNode.mux.Unlock()
	SendAndReceiveFindNode(lookupNode,lookupNode.target.ID,&contactToAsk.contact)
}



func (lookupNode *LookupNode) handleNewContact(contact Contact){
	contact.distance = contact.ID.CalcDistance(lookupNode.target.ID)
	lookupNode.mux.Lock()
	if !lookupNode.contactIsInShortlist(&contact){
		lookupNode.updateShortlistIfCloser(&contact)
	}
	lookupNode.mux.Unlock()
}

func (lookupNode *LookupNode) contactIsInShortlist(contact *Contact) bool{
	for _, v := range lookupNode.shortlist {
		if v.contact.ID.Equals(contact.ID){
			return true
		}
	}
	return false
}

func (lookupNode *LookupNode) updateShortlistIfCloser(contact *Contact){
	index := sort.Search(len(lookupNode.shortlist), func(i int) bool { return contact.distance.Less(lookupNode.shortlist[i].contact.distance)})
	if index<21 {
		if len(lookupNode.shortlist)<20 {
			lookupNode.shortlist = append(lookupNode.shortlist, LookupNodeContact{})
		} else{
			lookupNode.shortlist[19]=LookupNodeContact{}
			if index==20{
				index=19
			}
		}
		copy(lookupNode.shortlist[index+1:], lookupNode.shortlist[index:])
		lookupNode.shortlist[index] = NewLookupNodeContact(*contact)
	}
}

func (lookupNode *LookupNode) isKClosestContactHasBeenFound() bool{
	for _, v := range lookupNode.shortlist {
		if v.lookupState != ASKED {
			return false
		}
	}
	return true
}

//can return nil
func (lookupNode *LookupNode) getLookupNodeContactUnAsked() *LookupNodeContact{
	for i:= range lookupNode.shortlist {
		if lookupNode.shortlist[i].lookupState == UNASKED {
			return &lookupNode.shortlist[i]
		}
	}
	return nil
}

func (lookupNode *LookupNode) changeLookupNodeContactState(contact *LookupNodeContact,state LookupNodeContactState){
	for i := range lookupNode.shortlist {
		if lookupNode.shortlist[i].contact.ID == contact.contact.ID {
			lookupNode.shortlist[i].lookupState = state
			break
		}
	}
}




