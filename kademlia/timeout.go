package kademlia

import (
	"./protocol"
	"sync"
	"time"
)

var timeOutManager = NewTimeOutManager()

type Timeout struct {
	ch          chan *protocol.RPC
	messageID   messageID
	timeoutStop chan bool
}

func NewTimeout(messageID messageID, ch chan *protocol.RPC) *Timeout {
	timeout := &Timeout{}
	timeout.ch = ch
	timeout.messageID = messageID
	timeout.timeoutStop = make(chan bool)
	return timeout
}

func (timeout *Timeout) start() {
	go timeout.run()
}

func (timeout *Timeout) run() {
	select {
	case <-timeout.timeoutStop:
		return
	case <-time.After(10 * time.Second):
		timeOutManager := timeOutManager.tryGetAndRemoveTimeOut(timeout.messageID)
		if timeOutManager != nil {
			timeout.ch <- nil
		}
		return
	}
}

func (timeout *Timeout) stop() {
	timeout.timeoutStop <- true
}

type TimeOutManager struct {
	timeOutPingMap map[messageID]*Timeout
	mu             sync.Mutex
}

func NewTimeOutManager() *TimeOutManager {
	timeOutManager := &TimeOutManager{}
	timeOutManager.timeOutPingMap = make(map[messageID]*Timeout)
	return timeOutManager
}

func (timeOutManager *TimeOutManager) insertAndStart(timeout *Timeout) {
	timeOutManager.mu.Lock()
	timeOutManager.timeOutPingMap[timeout.messageID] = timeout
	timeOutManager.mu.Unlock()
	timeout.start()
}

//return nil if the time out ping is not in the map
func (timeOutManager *TimeOutManager) tryGetAndRemoveTimeOut(messageID messageID) *Timeout {
	timeOutManager.mu.Lock()
	if val, ok := timeOutManager.timeOutPingMap[messageID]; ok {
		delete(timeOutManager.timeOutPingMap, messageID)
		timeOutManager.mu.Unlock()
		return val
	}
	timeOutManager.mu.Unlock()
	return nil
}
