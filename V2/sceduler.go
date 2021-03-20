package V2

import "sync"

var tasks []*Task
var freeTaskIds []int

func AddTask(task *Task) {
	for _, testTask := range tasks {
		if testTask == task {
			return
		}
	}

	globalChanges <- func() {
		if len(freeTaskIds) > 0 {
			task.id = freeTaskIds[0]
			tasks[task.id] = task
			freeTaskIds = freeTaskIds[1:]
		} else {
			task.id = len(tasks)
			tasks = append(tasks, task)
		}
		updateTaskPlan()
	}
}
func RemoveTask(task *Task) {
	globalChanges <- func() {
		tasks[task.id] = nil
		freeTaskIds = append(freeTaskIds, task.id)
		updateTaskPlan()
	}
}

var taskPlan [][]*Task

func updateTaskPlan() {
	taskPlan = nil
	for _, task := range tasks {

		for i := 0; i < len(taskPlan); i++ {
			fits := true
			for j := 0; j < len(taskPlan[i]); j++ {
				testTask := taskPlan[i][j]

				for _, dataId := range task.writeData {
					for _, testDataId := range testTask.readData {
						if testDataId == dataId {
							fits = false
							break
						}
					}
					for _, testDataId := range testTask.writeData {
						if testDataId == dataId {
							fits = false
							break
						}
					}
				}
			}

			if fits {
				taskPlan[i] = append(taskPlan[i], task)
			}
		}

		taskPlan = append(taskPlan, []*Task{task})
	}
}

func runPlan() {
	for i := 0; i < len(taskPlan); i++ {

		wg := sync.WaitGroup{}
		for j := 0; j < len(taskPlan[i]); j++ {
			wg.Add(1)
			go func() {
				taskPlan[i][j].function()
				wg.Done()
			}()
		}
		wg.Wait()
		syncState()
	}
}
