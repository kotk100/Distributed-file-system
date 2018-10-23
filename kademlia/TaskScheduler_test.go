package kademlia

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func createTaskForTest(executeEvery time.Duration) *Task {
	task := &Task{}
	task.executeEvery = executeEvery
	return task
}

func TestSchedulerAddTask(t *testing.T) {
	os.Setenv("TIME_VARIATION", "2")
	periodicTasks := CreatePeriodicTasks()
	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute, task)

	time.Sleep(1 * time.Second)

	_, task2 := periodicTasks.treeMap.Find(createTaskComparator(task))
	assert.Equal(t, *task, task2, "should have the task")
}

func TestSchedulerGetNextTasks(t *testing.T) {
	os.Setenv("TIME_VARIATION", "0")
	periodicTasks := CreatePeriodicTasks()
	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute, task)

	task2 := createTaskForTest(1 * time.Minute)
	timeToExecute2 := time.Now().Add(task2.executeEvery)
	periodicTasks.addTask(&timeToExecute2, task2)

	task3 := createTaskForTest(30 * time.Minute)
	timeToExecute3 := time.Now().Add(task3.executeEvery)
	periodicTasks.addTask(&timeToExecute3, task3)

	nextTaskKey, _ := periodicTasks.getNextTask()

	assert.Equal(t, timeToExecute2, nextTaskKey.(time.Time), "keys should be equals")
}

func TestSchedulerUpdateTask(t *testing.T) {
	os.Setenv("TIME_VARIATION", "0")
	periodicTasks := CreatePeriodicTasks()

	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute, task)

	updated := periodicTasks.updateTask(task)

	assert.Equal(t, true, updated, "should have the task")
}
