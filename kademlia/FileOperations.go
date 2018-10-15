package kademlia

import (
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
)

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

func isFilePinned(hash string) bool {
	return isFilePinnedFilename(getFilenameFromHash(hash))
}

func isFilePinnedFilename(filename string) bool {
	s := strings.Split(filename, ":")

	if s[len(s)-1] == "1" {
		return true
	} else {
		return false
	}
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

func pinFileHash(hash string) bool {
	return pinFileFilename(getFilenameFromHash(hash))
}

func getPinnedFilename(filename string) string {
	return filename[:len(filename)-1] + "1"
}

func pinFileFilename(filename string) bool {
	if isFilePinnedFilename(filename) {
		return true
	}

	err := os.Rename(filename, getPinnedFilename(filename))

	if err != nil {
		log.WithFields(log.Fields{
			"Error":    err,
			"filename": filename,
		}).Error("Failed to pin/rename file.")
		return false
	}
	return true
}
