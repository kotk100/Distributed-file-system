package kademlia

import (
	"./protocol"
	"time"
)

type Timeout struct {
	ch          chan *protocol.RPC
	messageID   messageID
	timeoutStop chan bool
	isTimeout   bool
}

func NewTimeout(messageID messageID, ch chan *protocol.RPC) *Timeout {
	timeout := &Timeout{}
	timeout.ch = ch
	timeout.messageID = messageID
	timeout.timeoutStop = make(chan bool)
	timeout.isTimeout = false
	return timeout
}

func (timeout *Timeout) start() {
	go timeout.run()
}

func (timeout *Timeout) run() {
	select {
	case <-timeout.timeoutStop:
		return
	case <-time.After(5 * time.Second):
		if !timeout.isTimeout {
			timeout.isTimeout = true
			timeout.ch <- nil
		}
		return
	}
}

func (timeout *Timeout) stop() {
	if !timeout.isTimeout {
		timeout.timeoutStop <- true
	}
}