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
	timeout.start()

	time.Sleep(4 * time.Second)
	timeout.stop()
	time.Sleep(2 * time.Second)

	assert.Equal(t, false, timeout.isTimeout, "should not be time out")
}

func TestTimeout(t *testing.T) {
	// Generate random messageID
	messageID := messageID{}
	for i := 0; i < 20; i++ {
		messageID[i] = uint8(rand.Intn(256))
	}
	c := make(chan *protocol.RPC)

	timeout := NewTimeout(messageID, c)
	timeout.start()

	time.Sleep(12 * time.Second)

	assert.Equal(t, true, timeout.isTimeout, "should be time out")
}
