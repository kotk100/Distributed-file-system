package kademlia

import (
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/utils"
	log "github.com/sirupsen/logrus"
	rand2 "math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var PeriodicTasksReference *PeriodicTasks

type PeriodicTasks struct {
	treeMap  *treemap.Map
	mapLock  sync.Mutex
	ch       chan bool
	variance int64
}

func CreatePeriodicTasks() *PeriodicTasks {
	periodicTasks := &PeriodicTasks{}
	periodicTasks.treeMap = treemap.NewWith(utils.TimeComparator)
	periodicTasks.ch = make(chan bool, 1)

	// Get randomization time interval
	timeString := os.Getenv("TIME_VARIATION")
	timeSeconds, errb := strconv.Atoi(timeString)
	if errb != nil {
		log.WithFields(log.Fields{
			"Error": errb,
		}).Error("Failer to convert string to int")
	}
	periodicTasks.variance = int64(timeSeconds) * int64(time.Second)

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

			log.Info("HandleTask: Waiting complete, checking for new tasks.")
			// Read from channel, discard value
			_ = <-periodicTasks.ch
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

			log.Info("Waiting for task complete.")

			// Check if the task is ready to execute or if another routine woke up this task
			isTimeout := <-periodicTasks.ch
			if !isTimeout {
				log.Info("HandleTask routine woken up by another.")
				continue
			}

			task.executor.setTask(&task)
			go task.executor.execute()
			log.WithFields(log.Fields{
				"task": task,
			}).Info("Task executed.")

			// Remove task or reschedule
			if task.executeEvery.Nanoseconds() != 0 && task.taskType != ExpireFile {
				log.Info("Task being updated.")
				periodicTasks.updateTaskWithoutWakeUp(&task)
			} else {
				log.Info("Task removed.")
				periodicTasks.mapLock.Lock()
				periodicTasks.treeMap.Remove(timeToExecute)
				periodicTasks.mapLock.Unlock()
			}
		}
	}
}

func (periodicTasks *PeriodicTasks) addTask(timeToExecute *time.Time, task *Task) {
	periodicTasks.addTaskInternal(timeToExecute, task, false)
}

func (periodicTasks *PeriodicTasks) addTaskInternal(timeToExecute *time.Time, task *Task, forceWakeUp bool) {
	// Randomize timeToExecute
	rand := rand2.Int63n(periodicTasks.variance*2) - periodicTasks.variance
	timeToExecuteRand := timeToExecute.Add(time.Duration(rand))

	periodicTasks.mapLock.Lock()

	// All keys need to be different in map so check if the key already exists
	for _, found := periodicTasks.treeMap.Get(timeToExecuteRand); found; _, found = periodicTasks.treeMap.Get(timeToExecuteRand) {
		// Change time a little and retry
		timeToExecuteRand = timeToExecuteRand.Add(time.Nanosecond)
	}

	// Make sure if there are no tasks or if the added task is before the first to wake up the handleTasks routine
	needToWakeHandleTaskRoutine := false
	key, taskRet := periodicTasks.treeMap.Min()
	if key == nil || taskRet == nil {
		needToWakeHandleTaskRoutine = true
	} else {
		timeToExecuteMin, _ := castToTimeAndTask(key, taskRet)
		if timeToExecuteRand.Before(timeToExecuteMin) {
			needToWakeHandleTaskRoutine = true
		}
	}

	periodicTasks.treeMap.Put(timeToExecuteRand, *task)

	periodicTasks.mapLock.Unlock()

	// Wake up handleTask routine
	if needToWakeHandleTaskRoutine || forceWakeUp {
		// Non blocking write
		select {
		case periodicTasks.ch <- false:
		default:
		}
	}

	log.WithFields(log.Fields{
		"task": task,
	}).Debug("Added new periodic task to be executed.")
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

func (periodicTasks *PeriodicTasks) updateTask(dummyTask *Task) bool {
	return periodicTasks.updateTaskWithTime(dummyTask, nil)
}

func (periodicTasks *PeriodicTasks) updateTaskWithTime(dummyTask *Task, timeToExecuteParam *time.Time) bool {
	periodicTasks.mapLock.Lock()
	// Find the task
	timeToExecute, task := periodicTasks.treeMap.Find(createTaskComparator(dummyTask))
	if timeToExecute == nil || task == nil {
		periodicTasks.mapLock.Unlock()
		return false
	}
	taskC := task.(Task)

	// If this is the next routine in line to be executed the handleTasks routine needs to be informed
	needToWakeHandleTaskRoutine := false
	if taskC.taskComparator(periodicTasks.treeMap.Min()) {
		log.WithFields(log.Fields{
			"task": task,
		}).Debug("HandleTasks needs to be woken up.")
		needToWakeHandleTaskRoutine = true
	}

	// Remove the task
	periodicTasks.treeMap.Remove(timeToExecute)
	periodicTasks.mapLock.Unlock()

	var nextTimeToExecute time.Time
	if timeToExecuteParam == nil {
		nextTimeToExecute = time.Now().Add(taskC.executeEvery)
	} else {
		nextTimeToExecute = *timeToExecuteParam
	}

	// Add the task
	periodicTasks.addTaskInternal(&nextTimeToExecute, &taskC, needToWakeHandleTaskRoutine)
	log.WithFields(log.Fields{
		"task": task,
	}).Debug("Task execution time updated.")

	return true
}

func (periodicTasks *PeriodicTasks) updateTaskWithoutWakeUp(dummyTask *Task) bool {
	periodicTasks.mapLock.Lock()
	// Find the task
	timeToExecute, task := periodicTasks.treeMap.Find(createTaskComparator(dummyTask))
	if timeToExecute == nil || task == nil {
		periodicTasks.mapLock.Unlock()
		return false
	}
	taskC := task.(Task)

	// Remove the task
	periodicTasks.treeMap.Remove(timeToExecute)
	periodicTasks.mapLock.Unlock()

	nextTimeToExecute := time.Now().Add(taskC.executeEvery)

	// Add the task
	periodicTasks.addTaskInternal(&nextTimeToExecute, &taskC, false)
	log.WithFields(log.Fields{
		"task": task,
	}).Debug("Task rescheduled for next operation.")

	return true
}

func (periodicTasks *PeriodicTasks) removeTask(dummyTask *Task) bool {
	periodicTasks.mapLock.Lock()
	// Find the task
	timeToExecute, task := periodicTasks.treeMap.Find(createTaskComparator(dummyTask))

	if timeToExecute == nil || task == nil {
		log.WithFields(log.Fields{
			"task": dummyTask,
		}).Error("Task cannot be removed because it was not found.")

		periodicTasks.mapLock.Unlock()
		return false
	}

	taskC := task.(Task)

	// If this is the next routine in line to be executed the handleTasks routine needs to be informed
	needToWakeHandleTaskRoutine := false
	if taskC.taskComparator(periodicTasks.treeMap.Min()) {
		needToWakeHandleTaskRoutine = true
	}

	// Remove the task
	periodicTasks.treeMap.Remove(timeToExecute)
	periodicTasks.mapLock.Unlock()

	// Wake up handleTask routine
	if needToWakeHandleTaskRoutine {
		// Non blocking write
		select {
		case periodicTasks.ch <- false:
		default:
		}
	}

	log.WithFields(log.Fields{
		"task": task,
	}).Info("Task removed.")

	return true
}
