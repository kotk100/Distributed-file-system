package kademlia

import (
	"./protocol"
	"crypto/sha1"
	"encoding/hex"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Store struct {
	filename   string
	filehash   *[20]byte
	fileLength int64
}

func handleError(lenght int, err error) {
	if err != nil {
		log.WithFields(log.Fields{
			"lenght bytes": lenght,
			"Error":        err,
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
	store.fileLength = int64(len(*file))

	// Store file localy with pin set
	// TODO how to store file on first node that is the publisher so that the whole file doesn't have to be in memory
	filename := store.filename[:len(store.filename)-1] + "1"

	fileWriter := createFileWriter(filename)
	if fileWriter == nil {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed to create a fileWriter.")

		return nil
	}

	err := fileWriter.StoreFileChunk(file)
	fileWriter.close()
	if err {
		log.WithFields(log.Fields{
			"Filename": filename,
		}).Error("Failed writing to file.")

		return nil
	}

	return store
}

func CreateNewStoreForRepublish(fileHash string) *Store {
	store := &Store{}

	hash := StringToHash(fileHash)
	store.filehash = hash

	store.filename = getFilenameFromHash(fileHash)
	store.fileLength = getFileLength(store.filename)

	return store
}

func (store *Store) GetHash() string {
	return hashToString(store.filehash[:])
}

func (store *Store) StartStore() {
	log.Info("Started STORE procedure.")

	target := Contact{}
	target.ID = KademliaIDFromSlice(store.filehash[:])

	// Create periodic task for republishing
	republishTask := &Task{}
	republishTask.taskType = RefreshBucket
	republishTask.id = hashToString(store.filehash[:])

	timeString := os.Getenv("ORG_REPUBLISH_TIME")
	timeSeconds, err := strconv.Atoi(timeString)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failer to convert string to int")
	}
	republishTask.executeEvery = time.Duration(timeSeconds) * time.Second

	repTask := &RepublishTask{}
	republishTask.executor = repTask
	timeToExecute := time.Now().Add(republishTask.executeEvery)

	PeriodicTasksReference.addTask(&timeToExecute, republishTask)

	// Create periodic task for file expiration
	// TODO expiry time should be calculated dynamicaly, recalculate when a new node is added into the routing table
	expiryTask := &Task{}
	expiryTask.taskType = ExpireFile
	expiryTask.id = republishTask.id

	timeString = os.Getenv("EXPARATION_TIME")
	timeSeconds, err = strconv.Atoi(timeString)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failer to convert string to int")
	}
	expiryTask.executeEvery = time.Duration(timeSeconds) * time.Second

	expTask := &FileExpirationTask{}
	expiryTask.executor = expTask
	timeToExecute = time.Now().Add(expiryTask.executeEvery)

	PeriodicTasksReference.addTask(&timeToExecute, expiryTask)

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

		// Update periodic task for republishing and deletion
		task := &Task{}
		task.id = getHashFromFilename(store.Filename)
		task.taskType = RepublishFile
		PeriodicTasksReference.updateTask(task)

		task.taskType = ExpireFile
		PeriodicTasksReference.updateTask(task)

	} else {
		// Start listening on port and save file after connecting
		portStr := os.Getenv("FILE_TRANSFER_PORT")
		err := network.RecieveFile(portStr, store.Filename, other, id, store.FileSize, &rpc.OriginalSender)

		// Send confirmation if file was saved correctly
		if err {
			log.WithFields(log.Fields{
				"Contact": other,
			}).Error("Failed to receive file and save it.")

			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_ERROR, other, id, &rpc.OriginalSender))
		} else {
			errorStoreAnswerWraper(network.SendStoreAnswerMessage(protocol.StoreAnswer_OK, other, id, &rpc.OriginalSender))

			// Create periodic task for republishing
			republishTask := &Task{}
			republishTask.taskType = RefreshBucket
			republishTask.id = getHashFromFilename(store.Filename)

			timeString := os.Getenv("REPUBLISH_TIME")
			timeSeconds, err := strconv.Atoi(timeString)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Failer to convert string to int")
			}
			republishTask.executeEvery = time.Duration(timeSeconds) * time.Second

			repTask := &RepublishTask{}
			republishTask.executor = repTask
			timeToExecute := time.Now().Add(republishTask.executeEvery)

			PeriodicTasksReference.addTask(&timeToExecute, republishTask)

			// Create periodic tasks for expiration
			expiryTask := &Task{}
			expiryTask.taskType = ExpireFile
			expiryTask.id = republishTask.id

			timeString = os.Getenv("EXPARATION_TIME")
			timeSeconds, err = strconv.Atoi(timeString)
			if err != nil {
				log.WithFields(log.Fields{
					"Error": err,
				}).Error("Failer to convert string to int")
			}
			expiryTask.executeEvery = time.Duration(timeSeconds) * time.Second

			expTask := &FileExpirationTask{}
			expiryTask.executor = expTask
			timeToExecute = time.Now().Add(expiryTask.executeEvery)

			PeriodicTasksReference.addTask(&timeToExecute, expiryTask)
		}
	}
}

// TODO use locks to lock file access, so that a file isn't deleted while it is being read
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

func getFileLength(filename string) int64 {
	fi, err := os.Stat("/var/File_storage/" + filename)

	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Error("Error getting file length.")
	}
	// get the size
	return fi.Size()
}

func hashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}

func StringToHash(fileHash string) *[20]byte {
	hash, error := hex.DecodeString(fileHash)

	if error != nil {
		log.WithFields(log.Fields{
			"error": error,
		}).Error("Failed converting string back to hash.")

		return nil
	}

	var byteHash [20]byte
	copy(byteHash[:], hash)

	return &byteHash
}

func getHashFromFilename(filename string) string {
	s := strings.Split(filename, ":")
	return s[0]
}

func checkFileExists(filename string) bool {
	return checkFileExistsHash(getHashFromFilename(filename))
}

func getFilenameFromHash(fileHash string) string {
	filePath := getPathOfFileFromHash(fileHash)
	filePathPart := strings.Split(filePath, "/")

	return filePathPart[len(filePathPart)-1]
}

func getPathOfFileFromHash(fileHash string) string {
	matches, err := filepath.Glob("/var/File_storage/" + fileHash + ":*")
	if err != nil {
		log.WithFields(log.Fields{
			"filehash": fileHash,
		}).Error("Error getting filenames.")
		return ""
	}
	if len(matches) == 0 {
		return ""
	}
	return matches[0]
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

func removeFileByHash(hash string) bool {
	log.WithFields(log.Fields{
		"FileHash": hash,
	}).Info("Removing file.")
	return removeFile(getFilenameFromHash(hash))
}

func removeFile(filename string) bool {
	err := os.Remove("/var/File_storage/" + filename)

	if err != nil {
		log.WithFields(log.Fields{
			"Error":    err,
			"filename": filename,
		}).Error("Failed removing file.")
		return false
	}

	return true
}
