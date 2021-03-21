package V2

type Task struct {
	function  func()
	readData  []Data
	writeData []Data
	repeating bool
}

func NewTask(function func()) *Task {
	return &Task{
		function:  function,
		repeating: false,
	}
}

func (t *Task) SetReadData(data ...Data) {
	t.readData = data
}
func (t *Task) SetWriteData(data ...Data) {
	t.writeData = data
}
func (t *Task) SetRepeating(repeating bool) {
	t.repeating = repeating
}
