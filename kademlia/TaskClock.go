package kademlia

import (
	log "github.com/sirupsen/logrus"
	"time"
)

type TaskClock struct {
	taskClockStop chan bool
	timeToWait    time.Duration
}

func NewTaskClock(timeToWait time.Duration, ch chan bool) *TaskClock {
	taskClock := &TaskClock{}
	taskClock.timeToWait = timeToWait
	taskClock.taskClockStop = ch
	return taskClock
}

func (taskClock *TaskClock) run() {
	select {
	case res := <-taskClock.taskClockStop:
		log.WithFields(log.Fields{
			"res": res,
		}).Debug("Recieved message from channel")
		// Non blocking write
		select {
		case taskClock.taskClockStop <- false:
		default:
		}

		return
	case <-time.After(taskClock.timeToWait):
		// Non blocking write
		select {
		case taskClock.taskClockStop <- true:
		default:
		}
		return
	}
}
