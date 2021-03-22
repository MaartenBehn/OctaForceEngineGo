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
)

func initDispatcher() {
	addTask = make(chan *Task, 1)
	removeTask = make(chan *Task, 1)

	go updateTasksList()
}

func runDispatcher() {
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		frameStart = time.Now()

		dispatchTasks()

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
	}
}

var addTask chan *Task

func AddTask(task *Task) {
	go task.run()
	addTask <- task
}

var removeTask chan *Task

func RemoveTask(task *Task) {
	removeTask <- task
}

var repeatingTasks []*Task
var repeatingTasksStarted []bool
var oneTimeTasks []*Task

func updateTasksList() {
	for running {
		select {
		case task := <-addTask:
			if task.repeating {
				repeatingTasks = append(repeatingTasks, task)
				repeatingTasksStarted = append(repeatingTasksStarted, false)
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
					repeatingTasksStarted = append(repeatingTasksStarted[:i], repeatingTasksStarted[i+1:]...)
					break
				}
			}
		}
	}
}

var done bool
var aktiveTasks []*Task

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

		for i := len(oneTimeTasks) - 1; i >= 0; i-- {
			done = false
			task := oneTimeTasks[i]
			if canStartTask(task) {
				select {
				case task.start <- true:
					oneTimeTasks = append(oneTimeTasks[:i], oneTimeTasks[i+1:]...)
					aktiveTasks = append(aktiveTasks, task)
				default:
				}
			}
		}

		for i, task := range repeatingTasks {
			if !repeatingTasksStarted[i] && canStartTask(task) {
				done = false
				select {
				case task.start <- true:
					repeatingTasksStarted[i] = true
					aktiveTasks = append(aktiveTasks, task)
				default:
				}
			}
		}
	}

	for i := range repeatingTasksStarted {
		repeatingTasksStarted[i] = false
	}
}

func canStartTask(task *Task) bool {
	return true
	for _, aktiveTask := range aktiveTasks {
		for _, dependency := range aktiveTask.dependencies {
			for _, data := range task.dependencies {
				if data.checkDependency(dependency) {
					return false
				}
			}
		}
	}
	return true
}
