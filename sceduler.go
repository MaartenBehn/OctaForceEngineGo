package V2

var repeatingTasks []*Task
var repatingTaskChanged bool
var oneTimeTasks []*Task

func AddTask(task *Task) {
	addTask := NewTask(func() {
		if task.repeating {
			repeatingTasks = append(repeatingTasks, task)
		} else {
			oneTimeTasks = append(oneTimeTasks, task)
		}
		repatingTaskChanged = true
	})
	workers[workerSceduler].addTask(addTask)
}
