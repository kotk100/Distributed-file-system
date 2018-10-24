package kademlia

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type StoreFileResponse struct {
	Error      string
	FileHash   string
	IpWithFile []string
}

type StoreRestRequest struct {
	store *Store
	w     *http.ResponseWriter
	r     *http.Request
	c     chan bool
}

func NewStoreRestRequest(file *[]byte, password []byte, originalFilename string, w *http.ResponseWriter, r *http.Request) *StoreRestRequest {
	storeRestRequest := &StoreRestRequest{}
	storeRestRequest.store = CreateNewStore(file, password, originalFilename)
	storeRestRequest.w = w
	storeRestRequest.r = r
	storeRestRequest.c = make(chan bool)
	return storeRestRequest
}

func (storeRestRequest *StoreRestRequest) StartStore() {
	log.WithFields(log.Fields{
		"file hash": hashToString(storeRestRequest.store.filehash[:]),
	}).Info("Started STORE procedure.")

	target := Contact{}
	target.ID = KademliaIDFromSlice(storeRestRequest.store.filehash[:])

	lookupNode := NewLookupNode(target, storeRestRequest)
	lookupNode.Start()
	<-storeRestRequest.c
}

func (storeRestRequest *StoreRestRequest) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("Find_node finished for StoreRestRequest.")
	ipWithFile := make([]string, 0)
	for _, contact := range KClosestOfTarget {
		ipWithFile = append(ipWithFile, contact.contact.Address)
	}
	storeFileResponse := StoreFileResponse{"", storeRestRequest.store.GetHash(), ipWithFile}
	json.NewEncoder(*storeRestRequest.w).Encode(storeFileResponse)
	storeRestRequest.c <- true
	storeRestRequest.store.processKClosest(KClosestOfTarget)
}
