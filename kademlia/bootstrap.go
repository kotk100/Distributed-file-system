package kademlia

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type Bootstrap struct {
	BootstrapNode Contact
}

// Returned nodes from Find_node lookup
// All nodes should already be in the routing table as they have all been contacted and entered when/if they responded
func (bootstrap *Bootstrap) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("Bootstrap completed")
}

// Callback for pinging bootstraping node at start
func (bootstrap *Bootstrap) pingResult(receivedAnswer bool) {
	if receivedAnswer {
		log.Info("Receive answer from bootstraping node.")
		lookupNode := NewLookupNode(MyRoutingTable.GetMe(), bootstrap)
		lookupNode.Start()
	} else {
		log.Info("No response received, boostraping process failed")
		time.Sleep(30 * time.Second)
		SendAndRecievePing(&bootstrap.BootstrapNode, bootstrap)
	}
}
