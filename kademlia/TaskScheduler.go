package kademlia

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// TODO priority que of periodic tasks
// TODO each task has a time at which it needs to be exectuted at
// TODO each task has a function to be called to complete the tast
// TODO the task time can be refreshed, use locks to access the priority queue (if store for key recieved reset timer)
// TODO the tastks include:
// TODO 	- refresh bucket
// TODO 	- delete file
// TODO 	- republish file

var PeriodicTasksReference *PeriodicTasks

type PeriodicTasks struct {
	treeMap *treemap.Map
	mapLock sync.Mutex
	ch      chan bool
}

func CreatePeriodicTasks() *PeriodicTasks {
	periodicTasks := &PeriodicTasks{}
	periodicTasks.treeMap = treemap.NewWith(utils.TimeComparator)
	periodicTasks.ch = make(chan bool,1)

	// Create go routine that will service tasks
	go periodicTasks.handleTasks()

	log.Info("Created routine for servicing periodic tasks.")

	return periodicTasks
}

func castToTimeAndTask(key, task interface{}) (time.Time, Task) {
	return key.(time.Time), task.(Task)
}

func (periodicTasks *PeriodicTasks) getNextTask() (key, task interface{}) {
	// Get the task that needs to be executed first
	periodicTasks.mapLock.Lock()
	key, taskRet := periodicTasks.treeMap.Min()
	periodicTasks.mapLock.Unlock()

	return key, taskRet
}

// Handle all tasks and execute them at a appropriate time
func (periodicTasks *PeriodicTasks) handleTasks() {
	log.Info("Start handling.")
	// Repeat forever
	for {
		key, taskRet := periodicTasks.getNextTask()

		// If tree is empty wait to be woken up
		if key == nil || taskRet == nil {
			log.Info("No new tasks, waiting.")

			taskClock := NewTaskClock(24*3600*time.Second, periodicTasks.ch)
			taskClock.run()

			log.Info("stop waiting.")
		} else {
			// Cast returned results
			timeToExecute, task := castToTimeAndTask(key, taskRet)

			// Wait for task or interrupt on channel
			timeToWait := timeToExecute.Sub(time.Now())
			log.WithFields(log.Fields{
				"timeToWait": timeToWait,
			}).Info("Waiting for next task.")

			taskClock := NewTaskClock(timeToWait, periodicTasks.ch)
			taskClock.run()

			// Check if the task is ready to execute or if another routine woke up this task
			isTimeout := <-periodicTasks.ch
			if !isTimeout {
				log.Info("HandleTask routine woken up by another.")

				// Check if it is still the earliest task otherwise restart loop, the result should never be nil
				key2, taskRet2 := periodicTasks.getNextTask()
				timeToExecute2, _ := castToTimeAndTask(key2, taskRet2)

				if timeToExecute2.Before(timeToExecute) {
					continue
				}
			}

			// TODO execute task
			task.executor.setTask(&task)
			go task.executor.execute()
			log.WithFields(log.Fields{
				"task": task,
			}).Info("Task executed.")

			// Remove task or reschedule
			if task.executeEvery.Nanoseconds() != 0 {
				newTime := time.Now().Add(task.executeEvery)
				periodicTasks.updateTask(&task, &newTime)
			} else {
				periodicTasks.mapLock.Lock()
				periodicTasks.treeMap.Remove(timeToExecute)
				periodicTasks.mapLock.Unlock()
			}
		}

	}
}

func createTask(taskType TaskType, id string, executor TaskExecutor) *Task {
	task := &Task{}
	task.taskType = taskType
	task.id = id
	task.executor = executor

	return task
}

func (periodicTasks *PeriodicTasks) addTask(timeToExecute *time.Time, task *Task) {
	periodicTasks.mapLock.Lock()

	// All keys need to be different in map so check if the key already exists
	for _, found := periodicTasks.treeMap.Get(*timeToExecute); found; _, found = periodicTasks.treeMap.Get(*timeToExecute) {
		// Change time a little and retry
		timeToExecute.Add(time.Nanosecond)
	}

	// Make sure if there are no tasks or if the added task is before the first to wake up the handleTasks routine
	needToWakeHandleTaskRoutine := false
	key, taskRet := periodicTasks.treeMap.Min()
	if key == nil || taskRet == nil {
		needToWakeHandleTaskRoutine = true
	} else {
		timeToExecuteMin, _ := castToTimeAndTask(key, taskRet)
		if timeToExecute.Before(timeToExecuteMin) {
			needToWakeHandleTaskRoutine = true
		}
	}

	periodicTasks.treeMap.Put(*timeToExecute, *task)

	periodicTasks.mapLock.Unlock()

	// Wake up handleTask routine
	if needToWakeHandleTaskRoutine {
		periodicTasks.ch <- false
	}

	log.WithFields(log.Fields{
		"task": task,
	}).Info("Added new periodic task to be executed.")
}

func createTaskComparator(dummyTask *Task) func(interface{}, interface{}) bool {
	return dummyTask.taskComparator
}

func (task *Task) taskComparator(key interface{}, value interface{}) bool {
	_, taskToCompare := castToTimeAndTask(key, value)

	if task.taskType == taskToCompare.taskType && task.id == taskToCompare.id {
		return true
	}

	return false
}

func (periodicTasks *PeriodicTasks) updateTask(dummyTask *Task, newTimeToExecute *time.Time) bool {
	periodicTasks.mapLock.Lock()
	// Find the task
	timeToExecute, task := periodicTasks.treeMap.Find(createTaskComparator(dummyTask))

	if timeToExecute == nil || task == nil {
		return false
	}

	// Remove the task
	periodicTasks.treeMap.Remove(timeToExecute)

	taskC := task.(Task)
	periodicTasks.mapLock.Unlock()

	// Add the task
	periodicTasks.addTask(newTimeToExecute, &taskC)

	log.WithFields(log.Fields{
		"task": task,
	}).Info("Task execution time updated.")

	return true
}
