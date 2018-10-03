package kademlia

import "time"

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
	case <-taskClock.taskClockStop:
		return
	case <-time.After(taskClock.timeToWait):
		taskClock.taskClockStop <- true
		return
	}
}
