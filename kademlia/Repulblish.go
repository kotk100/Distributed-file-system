package kademlia

type RepublishTask struct {
	task *Task
}

func (republishTask *RepublishTask) execute() {
	// Get filehash of the file to be republished from task
	filehash := republishTask.task.id

	// Start republishing of the file by doing the STORE procedure
	store := CreateNewStoreForRepublish(filehash)
	store.StartStore()
}

func (republishTask *RepublishTask) setTask(task *Task) {
	republishTask.task = task
}
