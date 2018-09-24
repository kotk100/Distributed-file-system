package kademlia

import (
	"net"
	"os"
	log "github.com/sirupsen/logrus"
)

type Kademlia struct {
}

func InitMyInformation(port string){
	addrs, err := net.InterfaceAddrs()
	// TODO use logrus
	var IPaddress string
	if err != nil {
		os.Stderr.WriteString("Oops: " + err.Error() + "\n")
		os.Exit(1)
	}else{
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					IPaddress= ipnet.IP.String()
					break
				}
			}
		}
	}
	myKademlia := NewRandomKademliaID()
	me := NewContact(myKademlia,IPaddress + port)
	MyRoutingTable = NewRoutingTable(me)

	log.WithFields(log.Fields{
		"me":  MyRoutingTable.GetMe(),
	}).Info("My setting information (ID and adress.")
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
