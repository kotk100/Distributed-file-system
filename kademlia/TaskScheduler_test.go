package kademlia

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createTaskForTest(executeEvery time.Duration) *Task{
	task := &Task{}
	task.executeEvery = executeEvery
	return task
}


func TestAddTask(t *testing.T){
	periodicTasks := CreatePeriodicTasks()
	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute,task)

	time.Sleep(3 * time.Second)

	_,found := periodicTasks.treeMap.Get(timeToExecute)
	assert.Equal(t, true, found, "should have the task")
}

func TestGetNextTasks(t *testing.T){
	periodicTasks := CreatePeriodicTasks()
	fmt.Println("test 1")
	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute,task)
	fmt.Println("test 2")

	task2 := createTaskForTest(1 * time.Minute)
	timeToExecute2 := time.Now().Add(task2.executeEvery)
	periodicTasks.addTask(&timeToExecute2,task2)
	fmt.Println("test 3")

	task3 := createTaskForTest(30 * time.Minute)
	timeToExecute3 := time.Now().Add(task3.executeEvery)
	periodicTasks.addTask(&timeToExecute3,task3)

	nextTaskKey , _ := periodicTasks.getNextTask()

	assert.Equal(t,timeToExecute2 , nextTaskKey.(time.Time), "keys should be equals")
}

func TestUpdateTask(t *testing.T){
	periodicTasks := CreatePeriodicTasks()

	task := createTaskForTest(60 * time.Minute)
	timeToExecute := time.Now().Add(task.executeEvery)
	periodicTasks.addTask(&timeToExecute,task)

	newTimeToExecute := time.Now().Add(1 * time.Minute)
	updated := periodicTasks.updateTask(task,&newTimeToExecute)

	assert.Equal(t, true, updated, "should have the task")
	_,found := periodicTasks.treeMap.Get(newTimeToExecute)
	assert.Equal(t, true, found, "should have the task")

}
