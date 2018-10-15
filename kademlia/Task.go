package kademlia

import "time"

type TaskExecutor interface {
	execute()
	setTask(task *Task)
}

type TaskType int

const (
	RepublishFile      TaskType = 0
	ExpireFile         TaskType = 1
	RefreshBucket      TaskType = 2
	DeleteUnpinRequest TaskType = 3
)

type Task struct {
	taskType     TaskType
	id           string // Either a file hash or bucket ID
	executor     TaskExecutor
	executeEvery time.Duration
}

func (task *Task) equals(task2 *Task) bool {
	if task.taskType == task2.taskType && task.id == task.id {
		return true
	} else {
		return false
	}
}
