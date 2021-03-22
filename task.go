package OctaForce

type Task struct {
	function     func()
	repeating    bool
	start        chan bool
	done         chan bool
	dependencies []Data
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
func (t *Task) SetDependencies(dependencies ...Data) {
	t.dependencies = dependencies
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
