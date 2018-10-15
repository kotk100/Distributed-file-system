package kademlia

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"time"
)

type Bootstrap struct {
	BootstrapNode Contact
}

// Returned nodes from Find_node lookup
// All nodes should already be in the routing table as they have all been contacted and entered when/if they responded
func (bootstrap *Bootstrap) processKClosest(KClosestOfTarget []LookupNodeContact) {
	log.Info("Bootstrap completed.")

	// Refresh empty buckets
	//TODO too many find_node requests, stops working because of it.
	log.Info("Refreshing empty buckets.")

	for i := 0; i < IDLength*8; i++ {
		task := &Task{}
		task.id = strconv.Itoa(i)
		task.taskType = RefreshBucket

		// TODO only one find node at a time?
		// Check if bucket already has contacts
		if MyRoutingTable.buckets[i].Len() < 1 {
			timeToExecute := time.Now().Add(time.Duration(i) * 20 * time.Second)
			PeriodicTasksReference.updateTaskWithTime(task, &timeToExecute)
		}
	}
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
		SendAndRecievePing(bootstrap.BootstrapNode, bootstrap)
	}
}
