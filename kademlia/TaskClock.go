package kademlia

import "time"

type TaskClockAlarm interface {
	alarm()
}

type TaskClock struct {
	taskClockStop  chan bool
	timeToWait     time.Duration
	taskClockAlarm *TaskClockAlarm
}

func NewTaskClock(timeToWait time.Duration, taskClockAlarm TaskClockAlarm) *TaskClock {
	taskClock := &TaskClock{}
	taskClock.timeToWait = timeToWait
	taskClock.taskClockAlarm = &taskClockAlarm
	return taskClock
}

func (taskClock *TaskClock) start() {
	go taskClock.run()
}

func (taskClock *TaskClock) run() {
	select {
	case <-taskClock.taskClockStop:
		return
	case <-time.After(taskClock.timeToWait):
		(*taskClock.taskClockAlarm).alarm()
		return
	}
}

func (taskClock *TaskClock) stop() {
	taskClock.taskClockStop <- true
}
