package kademlia

import (
	"crypto/sha1"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"strings"
)

type Store struct {
	filename string
	filehash *[20]byte
	file     *[]byte
}

func handleError(lenght int, err error) {
	log.WithFields(log.Fields{
		"Error": err,
	}).Error("Error creating filename.")
}

// Creates a filename under which the file will be stored
func (store *Store) createFilename(file *[]byte, password []byte) {
	hash := strings.Builder{}
	handleError(hash.Write((*store.filehash)[:]))
	handleError(hash.WriteString(":"))
	handleError(hash.Write(sha1.Sum(password)[:]))

	// File is not pined by default
	handleError(hash.WriteString(":0"))

	store.filename = hash.String()
}

func createFileHash(file *[]byte) *[20]byte {
	hash := sha1.Sum(*file)
	return &hash
}

func (store *Store) StoreFile(file []byte, filename string) bool {
	err := ioutil.WriteFile("/var/File_storage/"+filename, file, 0644)

	if err != nil {
		log.WithFields(log.Fields{
			"Error":    err,
			"Filename": filename,
		}).Error("Failed to write file to local storage.")

		return false
	}

	return true
}

func createNewStore(file *[]byte, password []byte) *Store {
	store := &Store{}
	store.filehash = createFileHash(file)
	store.createFilename(file, password)
	store.file = file

	return store
}

func (store *Store) StartStore() {
	// Find k closest nodes to store file on
	// Returns k closest to callback processKClosest in this file
	lookupNode := NewLookupNode(MyRoutingTable.GetMe(), store)
	lookupNode.Start()
}

// Send store command to all k closest nodes
func (store *Store) processKClosest(KClosestOfTarget []LookupNodeContact) {
	for _, contact := range KClosestOfTarget {
		storeEx := &StoreRequestExecutor{}
		storeEx.store = store
		storeEx.contact = contact.contact
		createRoutine(storeEx)
	}
}
