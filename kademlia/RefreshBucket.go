package kademlia

import (
	log "github.com/sirupsen/logrus"
	"strconv"
	"sync"
)

type RefreshBucketTask struct {
	task *Task
	//orgTime time.Duration
}

// Only allow one find_node procedure for refreshing at once
var refreshBucketLock sync.Mutex

func (refreshBucketTask *RefreshBucketTask) execute() {
	// Get bucket id
	bucketIDString := refreshBucketTask.task.id
	bucketID, err := strconv.Atoi(bucketIDString)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failer to convert string (bucketID) to int.")

		return
	}

	// Get random ID from bucket
	randomID := MyRoutingTable.getRandomIDForBucket(bucketID)
	target := Contact{}
	target.ID = randomID

	// Make node_find for that ID
	refreshBucketLock.Lock()
	lookupNode := NewLookupNode(target, refreshBucketTask)
	lookupNode.Start()
}

func (refreshBucketTask *RefreshBucketTask) setTask(task *Task) {
	refreshBucketTask.task = task
	//refreshBucketTask.orgTime = task.executeEvery
}

func (refreshBucketTask *RefreshBucketTask) processKClosest(KClosestOfTarget []LookupNodeContact) {
	/*bucketIDString := refreshBucketTask.task.id
	bucketID, err := strconv.Atoi(bucketIDString)

	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Failer to convert string (bucketID) to int.")

		return;
	}

	// If no contacts found for the bucket increase time until next refresh.
	if MyRoutingTable.getBucketContactNumber(bucketID) < 1 {
		refreshBucketTask.task.executeEvery *= 2
	} else {
		refreshBucketTask.task.executeEvery = refreshBucketTask.orgTime
	}*/
	refreshBucketLock.Unlock()
}
