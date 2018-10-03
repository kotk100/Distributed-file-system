package kademlia

import (
	"./protocol"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestStopTimeout(t *testing.T) {
	// Generate random messageID
	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	c := make(chan *protocol.RPC)

	timeout := NewTimeout(messageID, c)
	timeOutManager := NewTimeOutManager()
	timeOutManager.insertAndStart(timeout)

	time.Sleep(4 * time.Second)

	timeoutEnd := timeOutManager.tryGetAndRemoveTimeOut(messageID)
	timeoutEnd.stop()

	assert.Equal(t, timeout, timeoutEnd, "time out should be not removed")
}

func TestTimeout(t *testing.T) {
	// Generate random messageID
	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	c := make(chan *protocol.RPC)

	timeout := NewTimeout(messageID, c)
	timeOutManager := NewTimeOutManager()
	timeOutManager.insertAndStart(timeout)

	time.Sleep(12 * time.Second)

	timeoutEnd := timeOutManager.tryGetAndRemoveTimeOut(messageID)

	assert.Equal(t, nil, timeoutEnd, "time out should be removed")
}
