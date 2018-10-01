package kademlia

import (
	log "github.com/sirupsen/logrus"
	"net"
	"os"
	"strings"
)

type Kademlia struct {
}

func InitMyInformation(port string) {
	var IPaddress string

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failed to get network interfaces.")

		os.Exit(1)
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil && strings.HasPrefix(ipnet.IP.String(), "10") {
					IPaddress = ipnet.IP.String()
					break
				}
			}
		}
	}
	myKademlia := NewRandomKademliaID()
	me := NewContact(myKademlia, IPaddress+port)
	MyRoutingTable = NewRoutingTable(me)

	log.WithFields(log.Fields{
		"me": MyRoutingTable.GetMe(),
	}).Info("My setting information (ID and adress.")
}
