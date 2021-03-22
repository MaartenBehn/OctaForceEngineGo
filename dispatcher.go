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

	syncWorkers()
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
		releaseWorkers()

		dispatchTasks()

		syncWorkers()
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
var workerWithTasks []int
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

	for i := 0; i < workerFixedMax; i++ {
		if len(workerTasks[i]) > 0 {
			workerWithTasks = append(workerWithTasks, i)
		}
	}

	oneTimeTasks = nil
}

var globalTasks = make(chan *Task, 1)

func dispatchTasks() {
	for len(workerWithTasks) > 0 || len(tasks) > 0 {

		if len(tasks) > 0 {
			select {
			case globalTasks <- tasks[0]:
				tasks = tasks[1:]
			default:
			}
		}

		for i := range workerWithTasks {
			if workers[i].tryAddTask(workerTasks[i][0]) {
				workerTasks[i] = workerTasks[i][1:]
				if len(workerTasks[i]) == 0 {
					workerWithTasks = append(workerWithTasks[:i], workerWithTasks[i+1:]...)
				}
			}
		}
	}
}
