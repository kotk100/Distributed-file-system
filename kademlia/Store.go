package kademlia

import (
	"./protocol"
	"crypto/sha1"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

type Store struct {
	filename string
	filehash *[20]byte
	file     *[]byte
}

func handleError(lenght int, err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error creating filename.")
	}
}

// Creates a filename under which the file will be stored
func (store *Store) createFilename(file *[]byte, password []byte, pinned bool, originalFilename string) {
	hash := strings.Builder{}

	str := hashToString(store.filehash[:])
	handleError(hash.WriteString(str))
	handleError(hash.WriteString(":"))

	if password != nil {
		passHash := sha1.Sum(password)
		str = hashToString(passHash[:])
		handleError(hash.WriteString(str))
	}

	handleError(hash.WriteString(":"))
	handleError(hash.WriteString(originalFilename))

	// File is not pinned by default
	if pinned {
		handleError(hash.WriteString(":1"))
	} else {
		handleError(hash.WriteString(":0"))
	}

	store.filename = hash.String()
}

func createFileHash(file *[]byte) *[20]byte {
	hash := sha1.Sum(*file)
	return &hash
}

func CreateNewStore(file *[]byte, password []byte, originalFilename string) *Store {
	store := &Store{}
	store.filehash = createFileHash(file)
	store.createFilename(file, password, false, originalFilename)
	store.file = file

	return store
}

func (store *Store) StartStore() {
	log.Info("Started STORE procedure.")

	// TODO Store file localy with pin set, should we pin file on publisher?
	filename := store.filename[:len(store.filename)-1] + "1"

	fileWriter := createFileWriter(filename)
	if fileWriter == nil {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed to create a fileWriter.")
	}

	err := fileWriter.StoreFileChunk(store.file)
	fileWriter.close()
	if err {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed writing to file.")
	}

	target := Contact{}
	target.ID = KademliaIDFromSlice(store.filehash[:])

	// Find k closest nodes to store file on
	// Returns k closest to callback processKClosest in this file
	lookupNode := NewLookupNode(target, store)
	lookupNode.Start()
}

// Send store command to all k closest nodes
func (store *Store) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("Find_node finished for STORE procedure.")
	for _, contact := range KClosestOfTarget {
		storeEx := &StoreRequestExecutor{}
		storeEx.store = store
		storeEx.contact = (contact.contact)
		createRoutine(storeEx)
	}
}

// Parse message inside RPC
func parseStoreRPC(rpc *protocol.RPC) *protocol.Store {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_STORE {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. STORE expected.")

		return nil
	}

	// Parse message as Store
	store := &protocol.Store{}
	if err := proto.Unmarshal(rpc.Message, store); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Ping message.")
	}

	return store
}

// Parse message inside RPC
func parseStoreAnswerRPC(rpc *protocol.RPC) *protocol.StoreAnswer {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_STORE_ANSWER {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type recieved. STORE ANSWER expected.")

		return nil
	}

	// Parse message as Store
	store := &protocol.StoreAnswer{}
	if err := proto.Unmarshal(rpc.Message, store); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incomming Store answer message.")
	}

	return store
}

func errorStoreAnswerWraper(err bool) {
	if err {
		log.Error("Failed to send Store Answer message.")
	}
}

func answerStoreRequest(rpc *protocol.RPC) {
	// Parse message and refresh routing table
	store := parseStoreRPC(rpc)
	other := createContactFromRPC(rpc)
	MyRoutingTable.AddContactAsync(*other)

	network := &Network{}
	id := messageID{}
	copy(id[:], rpc.MessageID[0:20])

	// Check if file already exsists
	if checkFileExists(store.Filename) {
		errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_ALREADY_STORED, other, id, &rpc.OriginalSender))
	} else {
		// Start listening on port and save file after connecting
		portStr := os.Getenv("FILE_TRANSFER_PORT")
		error := network.RecieveFile(portStr, store.Filename, other, id, store.FileSize, &rpc.OriginalSender)

		// Send confirmation if file was saved correctly
		if error {
			log.WithFields(log.Fields{
				"Contact": other,
			}).Error("Failed to receive file and save it.")

			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_ERROR, other, id, &rpc.OriginalSender))
		} else {
			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_OK, other, id, &rpc.OriginalSender))
		}
	}
}

// TODO use locks to lock file access?
type FileWriter struct {
	file *os.File
}

// Open file for writing
func createFileWriter(filename string) *FileWriter {
	f, err := os.Create("/var/File_storage/" + filename)

	if err != nil {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed opening file for writing.")

		return nil
	}

	writer := &FileWriter{}
	writer.file = f
	return writer
}

// Close file
func (fileWriter *FileWriter) close() {
	fileWriter.file.Close()
}

// TODO how to do it on first node that is the publisher?
// Writes a chunk of the file returning error
func (fileWriter *FileWriter) StoreFileChunk(file *[]byte) bool {
	n, err := fileWriter.file.Write(*file)

	if err != nil || n != len(*file) {
		log.WithFields(log.Fields{
			"Error":        err,
			"Writen bytes": n,
		}).Error("Failed writing a chunk of a file.")
		return true
	}

	fileWriter.file.Sync()
	return false
}

type FileReader struct {
	file *os.File
}

// Open file for writing
func createFileReader(filehash *[20]byte) *FileReader {
	str := hex.EncodeToString((*filehash)[:])

	// Find file by the hash if it exists
	matches, err := filepath.Glob("/var/File_storage/" + str + ":*")
	if err != nil {
		log.WithFields(log.Fields{
			"filehash": filehash,
		}).Error("Error getting filenames.")

		return nil
	}
	if len(matches) < 1 {
		log.WithFields(log.Fields{
			"filehash": filehash,
		}).Error("No file found with the given hash.")

		return nil
	}

	// Ignore all other matches, use only first one
	f, err := os.Open(matches[0])
	if err != nil {
		log.WithFields(log.Fields{
			"filename": matches[0],
		}).Error("Failed opening file for reading.")

		return nil
	}

	reader := &FileReader{}
	reader.file = f
	return reader
}

func hashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}

func checkFileExists(filename string) bool {
	s := strings.Split(filename, ":")
	filehash := s[0]
	return checkFileExistsHash(filehash)
}

func checkFileExistsHash(filehash string) bool {
	// Find file by the hash if it exists
	matches, err := filepath.Glob("/var/File_storage/" + filehash + ":*")
	if err != nil {
		log.WithFields(log.Fields{
			"filehash": filehash,
		}).Error("Error getting filenames.")

		return false
	}

	if len(matches) < 1 {
		return false
	}

	return true
}

// Close file
func (fileReader *FileReader) close() {
	fileReader.file.Close()
}

// Writes a chunk of the file returning error
func (fileReader *FileReader) ReadFileChunk(file *[]byte) (int, bool) {
	n, err := fileReader.file.Read(*file)

	if err != nil {
		log.WithFields(log.Fields{
			"Error":      err,
			"Read bytes": n,
		}).Error("Failed reading a chunk of a file.")
		return n, true
	}

	return n, false
}
