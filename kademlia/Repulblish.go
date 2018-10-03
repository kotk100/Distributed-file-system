package kademlia

type RepublishTask struct {
	task *Task
}

func (republishTask *RepublishTask) execute() {
	// Get filehash of the file to be republished from task
	filehash := republishTask.task.id

	// Make sure the file still exsists or if it was deleted
	if checkFileExistsHash(filehash) {
		// Start republishing of the file by doing the STORE procedure
		store := CreateNewStoreForRepublish(filehash)
		store.StartStore()
	}
}

func (republishTask *RepublishTask) setTask(task *Task) {
	republishTask.task = task
}

type FileExpirationTask struct {
	task *Task
}

func (fileExpirationTask *FileExpirationTask) execute() {
	// Delete file from filesystem
	removeFileByHash(fileExpirationTask.task.id)

	// TODO do not remove if pin is set

	// Remove the republishing task as the file is deleted
	task := &Task{}
	task.id = fileExpirationTask.task.id
	task.taskType = RepublishFile
	PeriodicTasksReference.removeTask(task)
}

func (fileExpirationTask *FileExpirationTask) setTask(task *Task) {
	fileExpirationTask.task = task
}
