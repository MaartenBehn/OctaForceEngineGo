package V2

type Task struct {
	function  func()
	repeating bool
	worker int
}

func NewTask(function func()) *Task {
	return &Task{
		function:  function,
		repeating: false,
		worker: -1,
	}
}

func (t *Task) SetRepeating(repeating bool) {
	t.repeating = repeating
}
func (t *Task) SetWorker(worker int) {
	t.worker = worker
}
