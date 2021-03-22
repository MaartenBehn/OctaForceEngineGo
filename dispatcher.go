package OctaForce

import (
	"log"
	"time"
)

var (
	running    bool
	FPS        float64
	MaxFPS     float64
	frameStart time.Time
	DeltaTime  float64
)

func initDispatcher() {
	addTask = make(chan *task, 1)
	removeTask = make(chan *task, 1)

	updateTasksListSync = make(chan bool)
	updateTasksListRelease = make(chan bool)

	go updateTasksList()
}

func runDispatcher() {
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		frameStart = time.Now()

		updateTasksListSync <- true

		copyTaskSlices()

		updateTasksListRelease <- true

		dispatchTasks()

		diff := time.Since(frameStart)
		if diff > 0 {
			FPS = (wait.Seconds() / diff.Seconds()) * MaxFPS
		} else {
			FPS = 10000
		}
		log.Print(FPS)

		if diff < wait {
			DeltaTime = wait.Seconds()
			time.Sleep(wait - diff)
		} else {
			DeltaTime = diff.Seconds()
		}
	}
}

var addTask chan *task

func AddTask(task *task) {
	go task.run()
	addTask <- task
}

var removeTask chan *task

func RemoveTask(task *task) {
	removeTask <- task
}

var repeatingTasks []*task
var oneTimeTasks []*task
var updateTasksListSync chan bool
var updateTasksListRelease chan bool

func updateTasksList() {
	for running {
		select {
		case <-updateTasksListSync:
			<-updateTasksListRelease

		case task := <-addTask:
			if task.repeating {
				repeatingTasks = append(repeatingTasks, task)
			} else {
				oneTimeTasks = append(oneTimeTasks, task)
			}

		case task := <-removeTask:
			if !task.repeating {
				break
			}

			for i, repeatingTask := range repeatingTasks {
				if repeatingTask == task {
					repeatingTasks = append(repeatingTasks[:i], repeatingTasks[i+1:]...)
					break
				}
			}
		}
	}
}

var tasks []*task

func copyTaskSlices() {
	tasks = nil

	for _, task := range repeatingTasks {
		tasks = append(tasks, task)
	}

	for _, task := range oneTimeTasks {
		tasks = append(tasks, task)
	}
	oneTimeTasks = nil
}

var done bool
var aktiveTasks []*task

func dispatchTasks() {
	done = false
	for !done {
		done = true

		for i := len(aktiveTasks) - 1; i >= 0; i-- {
			done = false
			select {
			case <-aktiveTasks[i].done:
				aktiveTasks = append(aktiveTasks[:i], aktiveTasks[i+1:]...)

			default:
			}
		}

		for i := len(tasks) - 1; i >= 0; i-- {
			done = false
			task := tasks[i]
			if canStartTask(task) {
				select {
				case task.start <- true:
					tasks = append(tasks[:i], tasks[i+1:]...)
					aktiveTasks = append(aktiveTasks, task)
				default:
				}
			}
		}
	}
}

func canStartTask(task *task) bool {
	for _, aktiveTask := range aktiveTasks {
		for _, raceTask := range task.raceTasks {
			if aktiveTask == raceTask {
				return false
			}
		}

		for _, raceTask := range aktiveTask.raceTasks {
			if task == raceTask {
				return false
			}
		}
	}
	return true
}
