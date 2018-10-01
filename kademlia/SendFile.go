package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"os"
)


func parseSendFileRequest(rpc *protocol.RPC) *protocol.SendFile {
	// Check type is correct
	if rpc.MessageType != protocol.RPC_SEND_FILE {
		log.WithFields(log.Fields{
			"Message": rpc,
		}).Error("Wrong message type received. SEND FILE expected.")
	}

	sendFile := &protocol.SendFile{}
	if err := proto.Unmarshal(rpc.Message, sendFile); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
			"Msg":   rpc.Message,
		}).Error("Failed to parse incoming SendFile message.")
	}

	return sendFile
}

func answerSendFileRequest(rpc *protocol.RPC) {
	log.WithFields(log.Fields{
	}).Info("Receive Send file request .")
	// Parse findNode message and create contact
	sendFile := parseSendFileRequest(rpc)
	sender := createContactFromRPC(rpc)
	port := os.Getenv("FILE_TRANSFER_PORT")
	network:= Network{}
	var fileHash [20]byte
	copy(fileHash[:],sendFile.FileHash)
	log.WithFields(log.Fields{
		"fileHash":  fileHash,
	}).Info("Answer to Send file request .")
	network.SendFile(&fileHash,sender,port,sendFile.FileSize)
}
