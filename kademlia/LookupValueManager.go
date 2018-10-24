package kademlia

import (
	"./protocol"
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strings"
)

type LookupValueManager struct {
	fileHash []byte
	w        *http.ResponseWriter
	r        *http.Request
	c     chan bool
}

type RequestError struct {
	Error string
}

func NewLookupValueManager(fileHash []byte, w *http.ResponseWriter, r *http.Request) *LookupValueManager {
	lookupValueManager := &LookupValueManager{}
	lookupValueManager.fileHash = fileHash
	lookupValueManager.w = w
	lookupValueManager.r = r
	lookupValueManager.c = make(chan bool)
	return lookupValueManager
}

func (lookupValueManager *LookupValueManager) FindValue() {

	log.Error("Start find value")

	log.WithFields(log.Fields{
		"file hash": hashToString(lookupValueManager.fileHash),
	}).Info("Start find value")
	lookupValue := NewLookupValue(lookupValueManager.fileHash, lookupValueManager)
	lookupValue.start()
	<-lookupValueManager.c
}

func (lookupValueManager *LookupValueManager) contactWithFile(contact Contact, findValueRpc *protocol.FindValue, contacts []LookupNodeContact) {
	log.WithFields(log.Fields{
		"Contact":    contact,
		"find value": findValueRpc,
	}).Info("Test find value : contactWithFile")

	network := Network{}
	portStr := os.Getenv("FILE_TRANSFER_PORT")
	originalSender := MyRoutingTable.me.ID[:]
	contactDidntSendValue := getClosestContactSeenWhichDidntSendFile(contact, contacts)
	error, pathFile := network.retrieveFile(portStr, lookupValueManager.fileHash, findValueRpc.FileName, &contact, findValueRpc.FileSize, &originalSender, contactDidntSendValue)
	if error {
		RequestError := RequestError{"error to retrieve file"}
		json.NewEncoder(*lookupValueManager.w).Encode(RequestError)
	} else {
		(*lookupValueManager.w).Header().Set("file_name", getDownloadFileName(pathFile))
		http.ServeFile(*lookupValueManager.w, lookupValueManager.r, pathFile)
	}
	lookupValueManager.c <- true
}

func getClosestContactSeenWhichDidntSendFile(conctatWhichSendValue Contact, contacts []LookupNodeContact) *Contact {
	if len(contacts) > 0 && !conctatWhichSendValue.ID.Equals(contacts[0].contact.ID) {
		return &contacts[0].contact
	} else if len(contacts) > 1 && !conctatWhichSendValue.ID.Equals(contacts[1].contact.ID) {
		return &contacts[1].contact
	}
	return nil
}

func (lookupValueManager *LookupValueManager) noContactWithFileFound(contacts []Contact) {
	log.Info("Test find value : noContactWithFileFound")
	(*lookupValueManager.w).Header().Set("file_name", "NULL")
	RequestError := RequestError{"error to retrieve file"}
	json.NewEncoder(*lookupValueManager.w).Encode(RequestError)
	lookupValueManager.c <- true
}

func (lookupValueManager *LookupValueManager) fileContents(fileContents []byte, stringPath string) {
	log.WithFields(log.Fields{
		"value": string(fileContents),
	}).Info("RECEIVED FILE VALUE (saved on the current node).------")
	(*lookupValueManager.w).Header().Set("file_name", getDownloadFileName(stringPath))
	http.ServeFile(*lookupValueManager.w, lookupValueManager.r, stringPath)
	lookupValueManager.c <- true
}

func getDownloadFileName(stringPath string) string {
	filePathPart := strings.Split(stringPath, "/")
	fileNamePart := strings.Split(filePathPart[len(filePathPart)-1], ":")
	return fileNamePart[len(fileNamePart)-2]
}
