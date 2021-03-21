package V2

var repeatingTasks []*Task
var repatingTaskChanged bool
var oneTimeTasks []*Task

func AddTask(task *Task) {
	workers[workerScedulerIn].(*taskWorker).addTask(func() {
		if task.repeating {
			repeatingTasks = append(repeatingTasks, task)
		} else {
			oneTimeTasks = append(oneTimeTasks, task)
		}
		repatingTaskChanged = true
	})
}

var repeatingPlan [][]*Task
var oneTimePlan [][]*Task

func updatePlans() {
	if repatingTaskChanged {
		repatingTaskChanged = false

		repeatingPlan = nil
		for _, task := range repeatingTasks {
			for i := 0; i < len(repeatingPlan); i++ {
				if doesTaskFitInPlan(i, task) {
					repeatingPlan[i] = append(repeatingPlan[i], task)
				}
			}
			repeatingPlan = append(repeatingPlan, []*Task{task})
		}
	}

	oneTimePlan = make([][]*Task, len(repeatingPlan))

	for _, task := range oneTimeTasks {
		for i := 0; i < len(repeatingPlan); i++ {
			if doesTaskFitInPlan(i, task) {
				oneTimePlan[i] = append(oneTimePlan[i], task)
			}
		}
		oneTimePlan = append(oneTimePlan, []*Task{task})
	}
	oneTimeTasks = nil
}
func doesTaskFitInPlan(i int, task *Task) bool {
	fits := true
	for j := 0; j < len(repeatingPlan[i]); j++ {
		testTask := repeatingPlan[i][j]

		for _, data := range task.writeData {
			for _, testdata := range testTask.writeData {
				if data == testdata {
					fits = false
				}
			}
		}
	}
	return fits
}

func dispatchPlan() {
	maxIndex := len(oneTimePlan)
	for i := 0; i < maxIndex; i++ {
		for _, task := range oneTimePlan[i] {
			dispatchTask(task)
		}

		if len(repeatingPlan) > i {
			for _, task := range repeatingPlan[i] {
				dispatchTask(task)
			}
		}
	}
}
func dispatchTask(task *Task) {
	added := false
	for !added {
		for _, worker := range taskWorkers {
			if worker.tryAddTask(task.function) {
				added = true
				break
			}
		}
	}
}
