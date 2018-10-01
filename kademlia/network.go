package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"net"
	"strings"
)

type NetworkAPI interface {
	SendPingMessage(originalSender *KademliaID, contact *Contact, messageID messageID) bool
}

type Network struct {
}

// TODO
//TODO routing messages
func handleIncomingMessage(buf []byte, addr *net.UDPAddr) {
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

//TODO parse message
func parseMessage(buf []byte) {

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

	// TODO size of buffer
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

		handleIncomingMessage(buf[0:n], addr)
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
	error := false

	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		error = true
	} else {
		// Create ping message
		ping := &protocol.Ping{}
		out, err := proto.Marshal(ping)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode PING message:")
			error = true
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
				error = true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
		// Close connection
		conn.Close()
	}
	return error
}

func (network *Network) SendFindContactMessage(targetId *KademliaID, originalSender *KademliaID, contact *Contact, contacts []Contact, messageID messageID) bool {
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
		out, err := createFindNodeToByte(targetId, contacts)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode FindNode message:")
			error = true
		} else {
			message := network.GetRPCMessage(out, protocol.RPC_FIND_NODE, messageID[:], originalSender[:])
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error = true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
		// Close connection
		conn.Close()
	}

	return error
}

func (network *Network) SendFindDataMessage(fileHash []byte, contact *Contact, contacts []Contact, messageID messageID, originalSender *KademliaID, haveTheFile bool) bool {
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
		out, err := createFindValueToByte(fileHash, contacts, haveTheFile)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode FindValue message:")
			error = true
		} else {
			message := network.GetRPCMessage(out, protocol.RPC_FIND_VALUE, messageID[:], originalSender[:])
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error = true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
		// Close connection
		conn.Close()
	}

	return error
}

func (network *Network) SendStoreMessage(filename string, lenght int64, contact *Contact, id messageID) bool {
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
		store := &protocol.Store{}
		store.Filename = filename
		store.FileSize = lenght
		out, err := proto.Marshal(store)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode STORE message:")
			error = true
		} else {
			// Wrap store message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_STORE, id[:], MyRoutingTable.GetMe().ID[:])
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error = true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
		// Close connection
		conn.Close()
	}

	return error
}

func (network *Network) SendStoreAnswerMessage(answer protocol.StoreAnswerStoreAnswer, contact *Contact, id messageID) bool {
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
		store := &protocol.StoreAnswer{}
		store.Answer = answer
		out, err := proto.Marshal(store)

		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode STORE ANSWER message:")
			error = true
		} else {
			// Wrap store message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_STORE_ANSWER, id[:], MyRoutingTable.GetMe().ID[:])
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error = true
			} else {
				log.Debug("Message writen to conn.")
			}
		}
		// Close connection
		conn.Close()
	}

	return error
}

// Buffer size for sending files over TCP
var BUFFERSIZE int64 = 1300

// port format;  ":<port>"
func (network *Network) SendFile(filehash *[20]byte, contact *Contact, port string) bool {
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
	for {
		_, err := fileReader.ReadFileChunk(&sendBuffer)
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
func (network *Network) RecieveFile(port string, filename string, contact *Contact, id messageID, fileSize int64) bool {
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
	error := network.SendStoreAnswerMessage(protocol.StoreAnswer_OK, contact, id)
	if error {
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
	defer fileWriter.close()

	var receivedBytes int64
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
			// TODO make sure that correct part of slice is used
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
