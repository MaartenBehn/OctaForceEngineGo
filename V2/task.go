package V2

type Task struct {
	id        int
	readData  []int
	writeData []int
	function  func()
}

func (t *Task) SetReadData(data ...int) {
	t.readData = data
}
func (t *Task) SetWriteData(data ...int) {
	t.writeData = data
}
func (t *Task) SetFunction(function func()) {
	t.function = function
}
