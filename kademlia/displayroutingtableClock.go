package kademlia

import "time"

type DisplayRoutingTableClock struct {
}

func (displayRoutingTableClock *DisplayRoutingTableClock) Display() {
	time.Sleep(100 * time.Second)
	MyRoutingTable.Print()

	time.Sleep(100 * time.Second)
	MyRoutingTable.Print()
}
