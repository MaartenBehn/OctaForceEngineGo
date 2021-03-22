package OctaForce

type Data interface{}

type taskTyp int

const (
	RenderTask       taskTyp = 0
	WindowUpdateTask taskTyp = 1
	taskMax                  = 2
)

var engineTasks []*task

func initState() {
	engineTasks = make([]*task, taskMax)
	for i := range engineTasks {
		engineTasks[i] = NewTask(func() {})
	}
}

func GetEngineTask(id taskTyp) *task {
	return engineTasks[id]
}
