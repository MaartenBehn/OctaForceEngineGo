package V2

const (
	workerTypNone = 0
	workerTypTask = 1
	workerTypFunc = 2
)

type worker struct {
	tasks    chan *Task
	sync     chan bool
	syncDone chan bool
}

func newWorker() *worker {
	return &worker{
		tasks:    make(chan *Task, 1),
		sync:     make(chan bool),
		syncDone: make(chan bool),
	}
}

func (w *worker) run() {
	for running {
		select {
		case task := <-w.tasks:
			task.function()
		case task := <-globalTasks:
			task.function()
		case <-w.sync:
			<-w.syncDone
		}
	}
}
func (w *worker) addTask(task *Task) {
	w.tasks <- task
}
func (w *worker) tryAddTask(task *Task) bool {
	select {
	case w.tasks <- task:
		return true
	default:
		return false
	}
}

const (
	workerRender   = 0
	workerSceduler = 1
	workerFixedMax = 2
)

var (
	workers              []*worker
	dynamicWorkerAmmount = 9
)

func initWorkers() {
	workers = make([]*worker, workerFixedMax+dynamicWorkerAmmount)

	for i, worker := range workers {
		worker = newWorker()
		workers[i] = worker

		if i == workerRender {
			continue
		}
		go worker.run()
	}
}
