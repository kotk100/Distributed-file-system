package kademlia

import (
	"sort"
	"sync"
	log "github.com/sirupsen/logrus"
)

const alpha = 3

type FindNodeRequestCallback interface {
	successRequest(contactAsked Contact,KClosestOfTarget []Contact)
	errorRequest(contactAsked Contact)
}

type LookupNodeCallback interface {
	processKClosest(KClosestOfTarget []LookupNodeContact)
}


type LookupNode struct{
	target Contact
	shortlist []LookupNodeContact
	mux sync.Mutex
	lookupNodeParallelism *LookupNodeParallelism
	lookupNodeCallback *LookupNodeCallback
}

func NewLookupNode(target Contact,lookupNodeCallback LookupNodeCallback) *LookupNode {
	lookup := &LookupNode{}
	lookup.shortlist=make([]LookupNodeContact,0)
	lookup.target = target
	lookup.lookupNodeParallelism=NewLookupNodeParallelism(lookup)
	lookup.lookupNodeCallback = &lookupNodeCallback
	return lookup
}

//give channel or callback to get K closest result
//select alpha contacts close of the target id
func (lookupNode *LookupNode) Start(){
	contactsCloseOfTarget := MyRoutingTable.FindClosestContacts(lookupNode.target.ID,alpha)
	log.WithFields(log.Fields{
		"Contacts alpha":   contactsCloseOfTarget,
	}).Info("Start find node.")
	if len(contactsCloseOfTarget)>0 {
		SendAndReceiveFindNode(lookupNode,lookupNode.target.ID,&contactsCloseOfTarget[0])
	}
	if len(contactsCloseOfTarget)>1 {
		SendAndReceiveFindNode(lookupNode,lookupNode.target.ID,&contactsCloseOfTarget[1])
	}
	if len(contactsCloseOfTarget)>2 {
		SendAndReceiveFindNode(lookupNode,lookupNode.target.ID,&contactsCloseOfTarget[2])
	}

	lookupNode.lookupNodeParallelism.start()
}

//callback for request executor
func (lookupNode *LookupNode) successRequest(contactAsked Contact,KClosestOfTarget []Contact){
	lookupNode.changeLookupContactState(contactAsked,ASKED)
	for _, v := range KClosestOfTarget {
		lookupNode.handleNewContact(v)
	}
	log.WithFields(log.Fields{
		"Contacts":   contactAsked,
		"KClosestOfTarget":KClosestOfTarget,
		"shortlist":   lookupNode.shortlist,
	}).Info("Find node get answer.")
	if lookupNode.isKClosestContactHasBeenFound(){
		log.WithFields(log.Fields{
			"KClosestOfTarget":lookupNode.shortlist,
		}).Info("Find node return k closest.")
		lookupNode.lookupNodeParallelism.stop()
		(*lookupNode.lookupNodeCallback).processKClosest(lookupNode.shortlist)
	}else{
		lookupNode.lookupNodeParallelism.restartClock()
		lookupNode.sendFindNodeRequest()
	}
}

func (lookupNode *LookupNode) errorRequest(contactAsked Contact){
	lookupNode.changeLookupContactState(contactAsked,FAILED)
	lookupNode.cleanShortList()
}

func (lookupNode *LookupNode) cleanShortList(){
	lookupNode.mux.Lock()
	for lookupNode.shortListNeedToBeCleaned() {
		failNodeIndex := lookupNode.getFailNodeIndex()
		lookupNode.shortlist =lookupNode.shortlist[:failNodeIndex+copy(lookupNode.shortlist[failNodeIndex:], lookupNode.shortlist[failNodeIndex+1:])]
	}
	lookupNode.mux.Unlock()
}

func (lookupNode *LookupNode) sendFindNodeRequest(){
	lookupNode.mux.Lock()
	contactToAsk := lookupNode.getLookupNodeContactUnAsked()
	if contactToAsk !=nil{
		lookupNode.changeLookupNodeContactState(contactToAsk,ASKING)
		lookupNode.mux.Unlock()
		log.WithFields(log.Fields{
			"contactToAsk":contactToAsk,
		}).Info("Send find node request from lookup node")
		SendAndReceiveFindNode(lookupNode,lookupNode.target.ID,&contactToAsk.contact)
	} else{
		lookupNode.mux.Unlock()
	}

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

func (lookupNode *LookupNode) shortListNeedToBeCleaned() bool{
	for _, v := range lookupNode.shortlist {
		if v.lookupState == FAILED {
			return true
		}
	}
	return false
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

func (lookupNode *LookupNode) getFailNodeIndex() int{
	for i:= range lookupNode.shortlist {
		if lookupNode.shortlist[i].lookupState == FAILED {
			return i
		}
	}
	return -1
}

func (lookupNode *LookupNode) changeLookupNodeContactState(contact *LookupNodeContact,state LookupNodeContactState){
	for i := range lookupNode.shortlist {
		if lookupNode.shortlist[i].contact.ID.Equals(contact.contact.ID) {
			lookupNode.shortlist[i].lookupState = state
			break
		}
	}
}

func (lookupNode *LookupNode) changeLookupContactState(contact Contact,state LookupNodeContactState){
	for i := range lookupNode.shortlist {
		if lookupNode.shortlist[i].contact.ID.Equals(contact.ID) {
			lookupNode.shortlist[i].lookupState = state
			break
		}
	}
}




