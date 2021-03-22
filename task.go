package OctaForce

type Task struct {
	function  func()
	repeating bool
	start     chan bool
	done      chan bool
	raceTasks []*Task
}

func NewTask(function func()) *Task {
	return &Task{
		function:  function,
		repeating: false,
		start:     make(chan bool),
		done:      make(chan bool, 1),
	}
}

func (t *Task) SetRepeating(repeating bool) {
	t.repeating = repeating
}
func (t *Task) SetRaceTask(tasks ...*Task) {
	t.raceTasks = tasks
}
func (t *Task) run() {
	for running {
		<-t.start
		t.function()

		t.done <- true
		if !t.repeating {
			break
		}
	}
}
