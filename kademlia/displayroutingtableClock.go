package kademlia

import "time"

type DisplayRoutingTableClock struct {
}

func (displayRoutingTableClock *DisplayRoutingTableClock) Display() {
	for {
		select {
		case <-time.After(60 * time.Second):
			MyRoutingTable.Print()
			break
		}
	}
}
