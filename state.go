package OctaForce

type Data interface{}

const (
	RenderTask       = 0
	WindowUpdateTask = 1
	taskMax          = 2
)

var engineTasks []*Task

func initState() {
	engineTasks = make([]*Task, taskMax)
	for i := range engineTasks {
		engineTasks[i] = NewTask(func() {})
	}
}

func GetEngineTask(id int) *Task {
	return engineTasks[id]
}
