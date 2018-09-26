package kademlia

type LookupNodeContactState int32

const (
	FAILED    LookupNodeContactState = 0
	UNASKED   LookupNodeContactState = 1
	ASKING    LookupNodeContactState = 2
	ASKED     LookupNodeContactState = 3
)

type LookupNodeContact struct{
	contact Contact
	lookupState LookupNodeContactState
}

func NewLookupNodeContact(contact Contact) LookupNodeContact {
	lookupContact := LookupNodeContact{}
	lookupContact.contact = contact
	lookupContact.lookupState = UNASKED
	return lookupContact
}
