package V2

const (
	workerTypNone = 0
	workerTypTask = 1
	workerTypFunc = 2
)

type worker struct {
	tasks    chan func()
	sync     chan bool
	syncDone chan bool
}

func newWorker() *worker {
	return &worker{
		tasks:    make(chan func(), 1),
		sync:     make(chan bool),
		syncDone: make(chan bool),
	}
}

func (w *worker) run() {
	for running {
		select {
		case task := <-w.tasks:
			task()
		case <-w.sync:
			<-w.syncDone
		}
	}
}
func (w *worker) addTask(function func()){
	w.tasks <- function
}
func (w *worker) tryAddTask(function func()) bool{
	select {
	case w.tasks <- function:
		return true
	default:
		return false
	}
}

const (
	workerRender = 0
	workerSceduler = 1
	workerFixedMax = 2
)
var (
	workers []*worker
	dynamicWorkerAmmount = 5
)

func initWorkers() {
	workers = make([]*worker, workerFixedMax + dynamicWorkerAmmount)

	for i, worker := range workers {
		worker = newWorker()
		workers[i] = worker

		if i == workerRender {continue}
		go worker.run()
	}
}