package kademlia

import (
	log "github.com/sirupsen/logrus"
)

type Bootstrap struct {

}

func (bootstrap *Bootstrap)processKClosest(KClosestOfTarget []LookupNodeContact){
	log.Info("Bootstrap completed")
	//for start:= time.Now();time.Since(start)<10*time.Second;{
//		runtime.Gosched()
//	}
//	MyRoutingTable.Print()
}

func (bootstrap *Bootstrap)pingResult(receivedAnswer bool){
	if receivedAnswer {
		log.Info("Receive answer from bootstraping node.")
		lookupNode:=NewLookupNode(MyRoutingTable.GetMe(),bootstrap)
		lookupNode.Start()
	} else {
		log.Info("No response received, boostraping process failed")

	}
}
