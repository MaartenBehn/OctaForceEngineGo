package V2

import (
	"log"
	"sync"
)

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
			fit := false
			for i := 0; i < len(repeatingPlan); i++ {
				if doesTaskFitInPlan(i, task) {
					repeatingPlan[i] = append(repeatingPlan[i], task)
					fit = true
				}
			}

			if !fit {
				repeatingPlan = append(repeatingPlan, []*Task{task})
			}
		}
	}

	oneTimePlan = make([][]*Task, len(repeatingPlan))

	for _, task := range oneTimeTasks {
		fit := false
		for i := 0; i < len(repeatingPlan); i++ {
			if doesTaskFitInPlan(i, task) {
				oneTimePlan[i] = append(oneTimePlan[i], task)
				fit = true
			}
		}
		if !fit {
			oneTimePlan = append(oneTimePlan, []*Task{task})
		}
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

var wg = sync.WaitGroup{}

func dispatchPlan() {
	maxIndex := len(oneTimePlan)
	for i := 0; i < maxIndex; i++ {
		wg.Add(len(oneTimePlan[i]))
		for _, task := range oneTimePlan[i] {
			dispatchTask(task)
		}

		wg.Add(len(repeatingPlan[i]))
		if len(repeatingPlan) > i {
			for _, task := range repeatingPlan[i] {
				dispatchTask(task)
			}
		}

		wg.Wait()
	}
	_, fps := frameEnd()
	log.Printf("%f \r", fps)
}
func dispatchTask(task *Task) {
	added := false
	for !added {
		for _, worker := range taskWorkers {
			if worker.tryAddTask(func() {
				task.function()
				wg.Done()
			}) {
				added = true
				break
			}
		}
	}
}
