package kademlia

import (
	"./protocol"
	"github.com/golang/protobuf/proto"
	log "github.com/sirupsen/logrus"
	"net"
)

type NetworkAPI interface {
	SendPingMessage(originalSender *KademliaID,contact *Contact, messageID messageID) bool
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
	}).Info("Recieved incoming message.")

	// Forward message to the right routine
	sendMessageToRoutine(rpc)
}

//TODO parse message
func parseMessage(buf []byte) {

}

func Listen( /*ip string,*/ port string) {
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
func (network *Network) SendPingMessage(originalSender *KademliaID,contact *Contact, messageID messageID) bool {
	// Open connection
	conn, err := net.Dial("udp", contact.Address)
	error := false
	if err != nil {
		log.WithFields(log.Fields{
			"Error":   err,
			"Contact": contact,
		}).Error("Failed to dial UDP address.")
		error=true
	}else{
		// Create ping message
		// TODO use correct ID
		ping := &protocol.Ping{}
		ping.KademliaID = MyRoutingTable.me.ID[:]
		out, err := proto.Marshal(ping)
		if err != nil {
			log.WithFields(log.Fields{
				"Error": err,
			}).Error("Failed to encode PING message:")
			error=true
		}else{
			//TODO Message ID generation
			// Wrap ping message and get bytes to send
			message := network.GetRPCMessage(out, protocol.RPC_PING, messageID[:],originalSender[:])
			// Write message to the connection (send to other node)
			n, err := conn.Write(message)
			if err != nil {
				log.WithFields(log.Fields{
					"Error":           err,
					"Number of bytes": n,
				}).Error("Failed to write message to connection.")
				error=true
			} else {
				log.Info("Message writen to conn.")
			}
		}
	}
	// Close connection
	conn.Close()
	return error
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
