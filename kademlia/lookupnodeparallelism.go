package kademlia

import "time"

type LookupNodeParallelism struct {
	lookupNode     *LookupNode
	chRestartClock chan bool
}

func NewLookupNodeParallelism(lookupNode *LookupNode) *LookupNodeParallelism {
	lookupNodeParallelism := &LookupNodeParallelism{}
	lookupNodeParallelism.lookupNode = lookupNode
	lookupNodeParallelism.chRestartClock = make(chan bool)
	return lookupNodeParallelism
}

func (lookupNodeParallelism *LookupNodeParallelism) start() {
	go lookupNodeParallelism.run()
}

func (lookupNodeParallelism *LookupNodeParallelism) stop() {
	lookupNodeParallelism.chRestartClock <- true
}

func (lookupNodeParallelism *LookupNodeParallelism) restartClock() {
	lookupNodeParallelism.chRestartClock <- true
	lookupNodeParallelism.start()
}

func (lookupNodeParallelism *LookupNodeParallelism) run() {
	select {
	case <-lookupNodeParallelism.chRestartClock:
		break
	case <-time.After(10 * time.Second):
		lookupNodeParallelism.lookupNode.parallelismSendFindNodeRequest()
		break
	}
}
