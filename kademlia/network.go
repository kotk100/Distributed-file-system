package kademlia

import (
	"fmt"
	"net"
)

type Network struct {
}

// Open read/write connection to given Contact
// TODO add error handling, close connection
func OpenConnection(contact *Contact) net.Conn {
	conn, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Println(err)
	}
	return conn
}

func handleConnection() {

}

func Listen( /*ip string,*/ port string) {
	// TODO read port from ENV
	ln, err := net.ListenPacket("udp", port)
	if err != nil {
		//TODO handle error
		fmt.Println(err)
	}
	buffer := make([]byte, 1024)
	ln.ReadFrom(buffer)
	fmt.Println("Message recieved.")
	fmt.Printf("%s\n", buffer)
}

// TODO delete
func (network *Network) SendPingMessageTMP(addr string) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		fmt.Println(err)
	}

	n, err := conn.Write([]byte("Hello from Node"))
	if err != nil {
		fmt.Print(n)
		fmt.Print("bytes   ")
		fmt.Println(err)
	}
	fmt.Println("Message writen to conn.")
	conn.Close()
}

func (network *Network) SendPingMessage(contact *Contact) {
	conn := OpenConnection(contact)
	conn.Write([]byte("Hello from Node"))
	conn.Close()

	// TODO close connection, error handling
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
