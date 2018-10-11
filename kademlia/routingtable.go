package kademlia

import (
	log "github.com/sirupsen/logrus"
	"os"
	"strconv"
	"time"
)

const bucketSize = 20

var MyRoutingTable *RoutingTable

// RoutingTable definition
// keeps a refrence contact of me and an array of buckets
type RoutingTable struct {
	me      Contact
	buckets [IDLength * 8]*bucket
}

// NewRoutingTable returns a new instance of a RoutingTable
func NewRoutingTable(me Contact) *RoutingTable {
	routingTable := &RoutingTable{}
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i] = newBucket()

		// Create periodic task for refreshing buckets
		refreshBucketTask := &Task{}
		refreshBucketTask.taskType = RefreshBucket
		refreshBucketTask.id = strconv.Itoa(i)

		timeString := os.Getenv("BUCKET_REFRESH_TIME")
		timeSeconds, errb := strconv.Atoi(timeString)
		if errb != nil {
			log.WithFields(log.Fields{
				"Error": errb,
			}).Error("Failer to convert string to int")
		}
		refreshBucketTask.executeEvery = time.Duration(timeSeconds) * time.Second

		refTask := &RefreshBucketTask{}
		refreshBucketTask.executor = refTask
		timeToExecute := time.Now().Add(refreshBucketTask.executeEvery)

		PeriodicTasksReference.addTask(&timeToExecute, refreshBucketTask)
	}
	routingTable.me = me
	return routingTable
}

// AddContactAsync add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContactAsync(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]

	// Update periodic task for refreshing buckets
	task := &Task{}
	task.id = strconv.Itoa(bucketIndex)
	task.taskType = RefreshBucket
	go PeriodicTasksReference.updateTask(task)

	go bucket.AddContact(contact)
}

// AddContactAsync add a new contact to the correct Bucket
func (routingTable *RoutingTable) AddContact(contact Contact) {
	bucketIndex := routingTable.getBucketIndex(contact.ID)
	bucket := routingTable.buckets[bucketIndex]

	// Update periodic task for refreshing buckets
	task := &Task{}
	task.id = strconv.Itoa(bucketIndex)
	task.taskType = RefreshBucket
	go PeriodicTasksReference.updateTask(task)

	bucket.AddContact(contact)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable
func (routingTable *RoutingTable) FindClosestContacts(target *KademliaID, count int) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.Append(bucket.GetContactAndCalcDistance(target))

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.Append(bucket.GetContactAndCalcDistance(target))
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// FindClosestContacts finds the count closest Contacts to the target in the RoutingTable which are not in the slice
func (routingTable *RoutingTable) FindClosestContactsNotInTheSlice(target *KademliaID, count int, contacts []Contact) []Contact {
	var candidates ContactCandidates
	bucketIndex := routingTable.getBucketIndex(target)
	bucket := routingTable.buckets[bucketIndex]

	candidates.AppendFilter(bucket.GetContactAndCalcDistance(target), contacts)

	for i := 1; (bucketIndex-i >= 0 || bucketIndex+i < IDLength*8) && candidates.Len() < count; i++ {
		if bucketIndex-i >= 0 {
			bucket = routingTable.buckets[bucketIndex-i]
			candidates.AppendFilter(bucket.GetContactAndCalcDistance(target), contacts)
		}
		if bucketIndex+i < IDLength*8 {
			bucket = routingTable.buckets[bucketIndex+i]
			candidates.AppendFilter(bucket.GetContactAndCalcDistance(target), contacts)
		}
	}

	candidates.Sort()

	if count > candidates.Len() {
		count = candidates.Len()
	}

	return candidates.GetContacts(count)
}

// getBucketIndex get the correct Bucket index for the KademliaID
func (routingTable *RoutingTable) getBucketIndex(id *KademliaID) int {
	distance := id.CalcDistance(routingTable.me.ID)
	for i := 0; i < IDLength; i++ {
		for j := 0; j < 8; j++ {
			if (distance[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return IDLength*8 - 1
}

// TODO make sure the actual id is not changed and a new one is created
func (routingTable *RoutingTable) getRandomIDForBucket(bucketID int) *KademliaID {
	id := *routingTable.me.ID
	j := bucketID % 8

	// Change ID so that it falls into the right bucket
	id[bucketID/8] ^= 0x1 << uint8(7-j)

	// Randomize it a bit, ignores up to 7 bits when doing it depending on the bucketID
	rand := NewRandomKademliaID()
	for i := (bucketID / 8) + 1; i < IDLength; i++ {
		id[i] ^= rand[i]
	}

	return &id
}

func (routingTable *RoutingTable) getBucketContactNumber(bucketID int) int {
	return routingTable.buckets[bucketID].list.Len()
}

func (routingTable *RoutingTable) GetMe() Contact {
	return routingTable.me
}

func (routingTable *RoutingTable) Print() {
	log.Info("routing table contents")
	for i := 0; i < IDLength*8; i++ {
		routingTable.buckets[i].print()
	}
}
