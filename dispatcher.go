package V2

import (
	"log"
	"time"
)

var (
	running    bool
	FPS        float64
	MaxFPS     float64
	frameStart time.Time
)

func runDispatcher() {
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	//syncWorkers()
	frameStart = time.Now()
	for running {
		diff := time.Since(frameStart)
		if diff > 0 {
			FPS = (wait.Seconds() / diff.Seconds()) * MaxFPS
		} else {
			FPS = 10000
		}
		log.Print(FPS)

		if diff < wait {
			time.Sleep(wait - diff)
		}
		frameStart = time.Now()

		copyTaskSlices()
		//releaseWorkers()

		dispatchTasks()

		//syncWorkers()
	}
}

func syncWorkers() {
	for _, worker := range workers {
		worker.sync <- true
	}
}
func releaseWorkers() {
	for _, worker := range workers {
		worker.syncDone <- true
	}
}

var workerTasks [][]*Task
var tasks []*Task

func copyTaskSlices() {
	workerTasks = make([][]*Task, workerFixedMax)

	for _, task := range repeatingTasks {
		if task.worker >= 0 {
			workerTasks[task.worker] = append(workerTasks[task.worker], task)
		} else {
			tasks = append(tasks, task)
		}
	}

	for _, task := range oneTimeTasks {
		if task.worker >= 0 {
			workerTasks[task.worker] = append(workerTasks[task.worker], task)
		} else {
			tasks = append(tasks, task)
		}
	}
	oneTimeTasks = nil
}

func dispatchTasks() {
	workerEmpti := false
	for !workerEmpti || len(tasks) > 0 {

		workerEmpti = true
		for i, worker := range workers {
			if i < workerFixedMax && len(workerTasks[i]) > 0 {
				workerEmpti = false

				if worker.tryAddTask(workerTasks[i][0].function) {
					workerTasks[i] = workerTasks[i][1:]
					continue
				}
			}

			if len(tasks) > 0 && worker.tryAddTask(tasks[0].function) {
				tasks = tasks[1:]
			}
		}
	}
}
