package kademlia

import (
	"./protocol"
	"container/list"
	log "github.com/sirupsen/logrus"
)

type PingBucketRequestExecutor struct{
	bucket *bucket
	ch chan *protocol.RPC
	id messageID
	contact *Contact
}

func (pingBucketRequestExecutor *PingBucketRequestExecutor) execute(){
	// Send ping message to other node
	error := pingBucketRequestExecutor.bucket.networkAPI.SendPingMessage(MyRoutingTable.me.ID,pingBucketRequestExecutor.contact, pingBucketRequestExecutor.id)

	var element *list.Element
	for e := pingBucketRequestExecutor.bucket.list.Front(); e != nil; e = e.Next() {
		nodeID := e.Value.(Contact).ID

		if (pingBucketRequestExecutor.contact).ID.Equals(nodeID) {
			element = e
		}
	}
	// Receive response message through channel
	if error {
		log.Info("Error to send bucket ping.")
		pingBucketRequestExecutor.bucket.list.Remove(element)
		pingBucketRequestExecutor.bucket.list.PushFront(pingBucketRequestExecutor.bucket.contactToInsert)
		pingBucketRequestExecutor.bucket.mux.Unlock()
		destroyRoutine(pingBucketRequestExecutor.id)
	}else{
		timeout:=NewTimeout(pingBucketRequestExecutor.id,pingBucketRequestExecutor.ch)
		timeOutManager.insertAndStart(timeout)
		rpc := <-pingBucketRequestExecutor.ch
		timeOutManager.tryGetAndRemoveTimeOut(pingBucketRequestExecutor.id)
		//timeout
		if element !=nil{
			if rpc == nil{
				log.Info("PING bucket time out.")
				pingBucketRequestExecutor.bucket.list.Remove(element)
				pingBucketRequestExecutor.bucket.list.PushFront(pingBucketRequestExecutor.bucket.contactToInsert)
			} else{
				log.Info("Received PING bucket message response.")
				pingBucketRequestExecutor.bucket.list.MoveToFront(element)
			}
		}
		pingBucketRequestExecutor.bucket.mux.Unlock()
		destroyRoutine(pingBucketRequestExecutor.id)
	}
}

func (pingBucketRequestExecutor *PingBucketRequestExecutor) setChannel(ch chan *protocol.RPC){
	pingBucketRequestExecutor.ch = ch
}

func (pingBucketRequestExecutor *PingBucketRequestExecutor) setMessageId(id messageID){
	pingBucketRequestExecutor.id = id
}