package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strings"
)

type NetworkAPI interface {
	SendPingMessage(originalSender *KademliaID, contact *Contact, messageID messageID) bool
}

type Network struct {
}

func handleIncomingMessage(buf []byte) {
	// Parse incoming message
	rpc := &protocol.RPC{}
	if err := proto.Unmarshal(buf, rpc); err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to parse incomming RPC message.")
	}

	log.WithFields(log.Fields{
		"RPC": rpc,
	}).Debug("Recieved incoming message.")

	// Forward message to the right routine
	sendMessageToRoutine(rpc)
}

func Listen(port string) {
	// Prepare an address at any interface at port 10001*/
	addr, err := net.ResolveUDPAddr("udp", port)
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to prepare UDP port.")
	}

	// Listen at selected port
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":    err,
			"UDP port": port,
		}).Error("Failed to start listening on UDP port.")
	}
	defer conn.Close()

	buf := make([]byte, 2048)

	// Listen for messages
	for {
		n, addr, err := conn.ReadFromUDP(buf)

		if err != nil {
			log.WithFields(log.Fields{
				"Error":   err,
				"Bytes":   n,
				"Address": addr,
			}).Error("Failed to recieve a message.")
		}

		handleIncomingMessage(buf[0:n])
	}
}

// Create RPC message wrapper and return bytes to send
func (network *Network) GetRPCMessage(message []byte, messageType protocol.RPCMessageTYPE, messageID []byte, originalSender []byte) (output []byte) {
	rpc := &protocol.RPC{}
	rpc.MessageType = messageType
	rpc.MessageID = messageID
	rpc.Message = message
	rpc.IPaddress = MyRoutingTable.me.Address
	rpc.OriginalSender = originalSender
	rpc.KademliaID = MyRoutingTable.me.ID[:]

	out, err := proto.Marshal(rpc)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to encode RPC message:")
	}

	return out
}

// Send ping message to another node
//if error false is return
func (network *Network) SendPingMessage(originalSender *KademliaID, contact *Contact, messageID messageID) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)

	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		return true
	} else {
		// Close connection
		defer conn.Close()

		// Create ping message
		ping := &protocol.Ping{}
		out, err := proto.Marshal(ping)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode PING message:")
			return true
		} else {
			// Wrap ping message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_PING, messageID[:], originalSender[:])
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				return true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
	}
	return false
}

func (network *Network) SendFindContactMessage(targetId *KademliaID, originalSender *KademliaID, contact *Contact, contacts []Contact, messageID messageID) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		return true
	} else {
		// Close connection
		defer conn.Close()

		out, err := createFindNodeToByte(targetId, contacts)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode FindNode message:")
			return true
		} else {
			message := network.GetRPCMessage(out, protocol.RPC_FIND_NODE, messageID[:], originalSender[:])
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				return true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
	}

	return false
}

func (network *Network) SendFindDataMessage(fileHash []byte, contact *Contact, contacts []Contact, messageID messageID, originalSender *KademliaID, haveTheFile bool, fileName string, fileSize int64) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		return true
	} else {
		// Close connection
		defer conn.Close()

		out, err := createFindValueToByte(fileHash, contacts, haveTheFile, fileName, fileSize)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode FindValue message:")
			return true
		} else {
			message := network.GetRPCMessage(out, protocol.RPC_FIND_VALUE, messageID[:], originalSender[:])
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				return true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
	}

	return false
}

func (network *Network) SendStoreMessage(filename string, lenght int64, contact *Contact, id messageID, originalSender *[]byte) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		return true
	} else {
		// Close connection
		defer conn.Close()

		store := &protocol.Store{}
		store.Filename = filename
		store.FileSize = lenght
		out, err := proto.Marshal(store)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode STORE message:")
			return true
		} else {
			// Wrap store message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_STORE, id[:], *originalSender)
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				return true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
	}

	return false
}

func (network *Network) SendStoreAnswerMessage(answer protocol.StoreAnswerStoreAnswer, contact *Contact, id messageID, originalSender *[]byte) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		return true
	} else {
		// Close connection
		defer conn.Close()

		store := &protocol.StoreAnswer{}
		store.Answer = answer
		out, err := proto.Marshal(store)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode STORE ANSWER message:")
			return true
		} else {
			// Wrap store message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_STORE_ANSWER, id[:], *originalSender)
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				return true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
	}

	return false
}

// Buffer size for sending files over TCP
var BUFFERSIZE int64 = 1300

// port format;  ":<port>"
func (network *Network) SendFile(filehash *[20]byte, contact *Contact, port string, fileLenght int64) bool {
	// Create address to connect to and connect to it
	s := strings.Split(contact.Address, ":")
	connection, err := net.Dial("tcp", s[0]+port)
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to make a TCP connection to contact.")

		return true
	}
	defer connection.Close()
	log.WithFields(log.Fields{
		"Contact": contact,
	}).Debug("Connected via TCP to Contact for sending file.")

	fileReader := createFileReader(filehash)
	if fileReader == nil {
		log.WithFields(log.Fields{
			"Error":    err,
			"Filehash": filehash,
		}).Error("Failed to open file for reading.")

		return true
	}
	defer fileReader.close()

	sendBuffer := make([]byte, BUFFERSIZE)
	var totalReadBytes int64

	for fileLenght > totalReadBytes {
		readBytes, err := fileReader.ReadFileChunk(&sendBuffer)
		totalReadBytes += int64(readBytes)

		if err {
			log.WithFields(log.Fields{
				"Error":    err,
				"Filehash": filehash,
			}).Error("Failed to read from file.")
			return true
		}
		connection.Write(sendBuffer)
	}

	log.WithFields(log.Fields{
		"Contact": contact,
	}).Debug("Sent file over TCP to contact.")
	return false
}

// port format;  ":<port>"
// Returns true if an error occured
func (network *Network) RecieveFile(port string, filename string, contact *Contact, id messageID, fileSize int64, originalSender *[]byte) bool {
	// Start listening on specified port
	server, err := net.Listen("tcp", port)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to start listening for a TCP connection.")

		return true
	}

	log.WithFields(log.Fields{
		"Port":    port,
		"Contact": contact,
	}).Debug("Started listening for a TCP connection.")

	// Close connection when exiting function
	defer server.Close()

	// Send response indicating that we started listening and are ready for the file
	errBool := network.SendStoreAnswerMessage(protocol.StoreAnswer_OK, contact, id, originalSender)
	if errBool {
		log.WithFields(log.Fields{
			"Contact": contact,
		}).Error("Failed to send Store Answer message.")

		return true
	}

	// Connect to incoming TCP connection
	connection, err := server.Accept()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to accept a TCP connection.")

		return true
	}
	log.WithFields(log.Fields{
		"Port":    port,
		"Contact": contact,
	}).Debug("Client connected.")

	// Open file for writing
	fileWriter := createFileWriter(filename)
	if fileWriter == nil {
		log.WithFields(log.Fields{
			"filename": filename,
		}).Error("Failed to create file for writing.")
		return true
	}

	defer fileWriter.close()

	var receivedBytes int64 = 0
	buffer := make([]byte, BUFFERSIZE)

	for {
		// Read packet
		n, err := connection.Read(buffer)
		if err != nil {
			log.WithFields(log.Fields{
				"Error":      err,
				"bytes read": n,
			}).Error("Failed to read bytes from TCP connection.")

			return true
		}

		// Store last part of file
		if (fileSize - receivedBytes) < BUFFERSIZE {
			endBuffer := buffer[0:(fileSize - receivedBytes)]
			err := fileWriter.StoreFileChunk(&endBuffer)
			if err {
				log.WithFields(log.Fields{
					"buffer": &buffer,
				}).Error("Failed to write bytes to file.")

				return true
			}

			break
		} else {
			err := fileWriter.StoreFileChunk(&buffer)
			if err {
				log.WithFields(log.Fields{
					"buffer": &buffer,
				}).Error("Failed to write bytes to file.")

				return true
			}
			receivedBytes += BUFFERSIZE

			if fileSize == receivedBytes {
				break
			}
		}
	}

	log.WithFields(log.Fields{
		"Filename": filename,
	}).Info("Received and stored file.")
	return false
}

func (network *Network) retrieveFile(port string, fileHash []byte, filename string, contact *Contact, fileSize int64, originalSender *[]byte) bool {
	server, err := net.Listen("tcp", port)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to start listening for a TCP connection.")

		return true
	}

	log.WithFields(log.Fields{
		"Port":    port,
		"Contact": contact,
	}).Debug("Started listening for a TCP connection.")

	// Close connection when exiting function
	defer server.Close()

	error := network.sendSendFileRequest(contact, fileHash, fileSize)
	if error {
		log.WithFields(log.Fields{
			"Contact": contact,
		}).Error("Failed to send SendFile request.")

		return true
	}

	// Connect to incoming TCP connection
	connection, err := server.Accept()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to accept a TCP connection.")

		return true
	}
	log.WithFields(log.Fields{
		"Port":    port,
		"Contact": contact,
	}).Debug("Client connected.")

	var receivedBytes int64 = 0
	buffer := make([]byte, BUFFERSIZE)

	for {
		// Read packet
		n, err := connection.Read(buffer)
		if err != nil {
			log.WithFields(log.Fields{
				"Error":      err,
				"bytes read": n,
			}).Error("Failed to read bytes from TCP connection.")

			return true
		}
		log.WithFields(log.Fields{
			"fileSize":      fileSize,
			"receivedBytes": receivedBytes,
		}).Info("receive send file answer from client")
		// Store last part of file
		if (fileSize - receivedBytes) < BUFFERSIZE {
			endBuffer := buffer[0:(fileSize - receivedBytes)]
			log.WithFields(log.Fields{
				"value": string(endBuffer),
			}).Info("RECEIVED FILE VALUE.------")
			break
		} else {
			log.WithFields(log.Fields{
				"value": string(buffer),
			}).Info("RECEIVED FILE VALUE.------")
			receivedBytes += BUFFERSIZE
			if fileSize == receivedBytes {
				break
			}
		}
	}

	log.WithFields(log.Fields{
		"Filename": filename,
	}).Info("END FILE VALUE-----")
	return false
}

func (network *Network) sendSendFileRequest(contact *Contact, fileHash []byte, fileSize int64) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	error := false
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		error = true
	} else {
		sendFile := &protocol.SendFile{}
		sendFile.FileHash = fileHash
		sendFile.IPaddress = MyRoutingTable.me.Address
		sendFile.KademliaID = MyRoutingTable.me.ID[:]
		sendFile.OriginalSender = MyRoutingTable.me.ID[:]
		sendFile.FileSize = fileSize
		out, err := proto.Marshal(sendFile)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to create send file request.")
			error = true
		} else {
			messageID := messageID{}
			for i := 0; i < 20; i++ {
				messageID[i] = uint8(rand.Intn(256))
			}
			message := network.GetRPCMessage(out, protocol.RPC_SEND_FILE, messageID[:], MyRoutingTable.me.ID[:])
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error = true
			} else {
				log.Info("Send file request sent.")
			}
		}
		// Close connection
		conn.Close()
	}
	return error
}
