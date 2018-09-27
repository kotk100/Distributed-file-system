package kademlia

import "time"

type LookupNodeParallelism struct {
	lookupNode *LookupNode
	chRestartClock chan bool
	isRunning bool
}


func NewLookupNodeParallelism(lookupNode *LookupNode) *LookupNodeParallelism{
	lookupNodeParallelism:=&LookupNodeParallelism{}
	lookupNodeParallelism.lookupNode=lookupNode
	lookupNodeParallelism.chRestartClock = make(chan bool)
	lookupNodeParallelism.isRunning = false
	return lookupNodeParallelism
}

func (lookupNodeParallelism *LookupNodeParallelism) start(){
	lookupNodeParallelism.isRunning = true
	go lookupNodeParallelism.run()
}

func (lookupNodeParallelism *LookupNodeParallelism) stop(){
	lookupNodeParallelism.isRunning = false
	lookupNodeParallelism.chRestartClock <-true
}

func (lookupNodeParallelism *LookupNodeParallelism) restartClock(){
	lookupNodeParallelism.chRestartClock <-true
}

func (lookupNodeParallelism *LookupNodeParallelism) run(){
	for lookupNodeParallelism.isRunning {
		select {
		case <-lookupNodeParallelism.chRestartClock:
			break
		case <-time.After(10 * time.Second):
			lookupNodeParallelism.lookupNode.sendFindNodeRequest()
		break
		}
	}
}

