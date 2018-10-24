package kademlia

import (
	"./protocol"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type PingBucketRequestExecutor struct {
	bucket       *bucket
	ch           chan *protocol.RPC
	id           messageID
	contactToAdd Contact
}

/*
bucket.contactupdates.Acquire(context.TODO(), 1)
	bucket.readers.Acquire(context.TODO(), allowedReaders)
	bucket.bucketsLocked += 1;
*/

func (pingBucketRequestExecutor *PingBucketRequestExecutor) execute() {
	// Reserve next bucket for pinging
	pingBucketRequestExecutor.bucket.readers.Acquire(context.TODO(), 1)
	pingBucketRequestExecutor.bucket.contactupdates.Acquire(context.TODO(), 1)
	pingBucketRequestExecutor.bucket.readers.Release(1)

	// Dosn't matter if the contact is changed in between these two statements
	// Do it in two steps so that if all the contacts are locked a dead lock does not occur

	// Get write access and update the number of locked buckets and get the bucket to be updated
	pingBucketRequestExecutor.bucket.readers.Acquire(context.TODO(), allowedReaders)
	elementToReplace := pingBucketRequestExecutor.bucket.list.Back()
	for i := 0; i < pingBucketRequestExecutor.bucket.bucketsLocked; i++ {
		elementToReplace = elementToReplace.Prev()
	}

	pingBucketRequestExecutor.bucket.bucketsLocked += 1
	pingBucketRequestExecutor.bucket.readers.Release(allowedReaders)
	contactToReplace := elementToReplace.Value.(Contact)

	// Send ping message to other node
	error := pingBucketRequestExecutor.bucket.networkAPI.SendPingMessage(MyRoutingTable.me.ID, &contactToReplace, pingBucketRequestExecutor.id)

	if error {
		log.Error("Error to send bucket ping.")
		pingBucketRequestExecutor.bucket.readers.Acquire(context.TODO(), allowedReaders)

		pingBucketRequestExecutor.bucket.list.Remove(elementToReplace)
		pingBucketRequestExecutor.bucket.list.PushFront(pingBucketRequestExecutor.contactToAdd)

		pingBucketRequestExecutor.bucket.contactupdates.Release(1)
		pingBucketRequestExecutor.bucket.bucketsLocked -= 1
		pingBucketRequestExecutor.bucket.readers.Release(allowedReaders)
		destroyRoutine(pingBucketRequestExecutor.id)
	} else {
		// Receive response message through channel
		timeout := NewTimeout(pingBucketRequestExecutor.id, pingBucketRequestExecutor.ch)
		timeout.start()
		rpc := <-pingBucketRequestExecutor.ch
		timeout.stop()
		destroyRoutine(pingBucketRequestExecutor.id)

		pingBucketRequestExecutor.bucket.readers.Acquire(context.TODO(), allowedReaders)
		if rpc == nil {
			log.Info("PING bucket time out.")
			pingBucketRequestExecutor.bucket.list.Remove(elementToReplace)
			pingBucketRequestExecutor.bucket.list.PushFront(pingBucketRequestExecutor.contactToAdd)
		} else {
			log.Info("Received PING bucket message response.")
			pingBucketRequestExecutor.bucket.list.MoveToFront(elementToReplace)
		}

		pingBucketRequestExecutor.bucket.contactupdates.Release(1)
		pingBucketRequestExecutor.bucket.bucketsLocked -= 1
		pingBucketRequestExecutor.bucket.readers.Release(allowedReaders)
	}
}

func (pingBucketRequestExecutor *PingBucketRequestExecutor) setChannel(ch chan *protocol.RPC) {
	pingBucketRequestExecutor.ch = ch
}

func (pingBucketRequestExecutor *PingBucketRequestExecutor) setMessageId(id messageID) {
	pingBucketRequestExecutor.id = id
}
