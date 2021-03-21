package V2

const (
	workerTypNone = 0
	workerTypTask = 1
	workerTypFunc = 2
)

type iworker interface {
	getTyp() int
	run()
	sync()
	syncDone()
}

type worker struct {
	typ          int
	syncChan     chan bool
	syncDoneChan chan bool
	fixed        bool
}

func newWorker(fixed bool, typ int) *worker {
	return &worker{
		typ:          typ,
		syncChan:     make(chan bool),
		syncDoneChan: make(chan bool),
		fixed:        fixed,
	}
}

func (w *worker) getTyp() int {
	return w.typ
}
func (w *worker) run() {
	for range w.syncChan {
		<-w.syncDoneChan
	}
}
func (w *worker) sync() {
	w.syncChan <- true
}
func (w *worker) syncDone() {
	w.syncDoneChan <- true
}

type taskWorker struct {
	*worker
	tasksIn chan func()
}

func newTaskWorker(fixed bool) *taskWorker {
	return &taskWorker{
		worker:  newWorker(fixed, workerTypTask),
		tasksIn: make(chan func()),
	}
}

func (w *taskWorker) run() {
	for {
		select {
		case task := <-w.tasksIn:
			task()
		case <-w.syncChan:
			<-w.syncDoneChan
		}
	}
}
func (w *taskWorker) addTask(function func()) {
	w.tasksIn <- function
}
func (w *taskWorker) tryAddTask(function func()) bool {
	select {
	case w.tasksIn <- function:
		return true
	default:
		return false
	}
}

type funcWorker struct {
	*worker
	function func()
}

func newFuncWorker(fixed bool, function func()) *funcWorker {
	return &funcWorker{
		worker:   newWorker(fixed, workerTypFunc),
		function: function,
	}
}

func (w *funcWorker) run() {
	for {
		w.function()

		<-w.syncChan
		<-w.syncDoneChan
	}

}

var workers []iworker
var taskWorkers []*taskWorker

const (
	workerScedulerOut = 0
	workerScedulerIn  = 1
	workerRender      = 2
	workerState       = 3
	workerFixedMax    = 4
)

var dynamicWorkerAmmount = 10

func initFixedWorkers() {
	workers = make([]iworker, workerFixedMax+dynamicWorkerAmmount)

	workers[workerScedulerIn] = newTaskWorker(true)
	workers[workerScedulerOut] = newFuncWorker(true, dispatchPlan)

	workers[workerState] = newTaskWorker(true)

	for i := workerFixedMax; i < workerFixedMax+dynamicWorkerAmmount; i++ {
		workers[i] = newTaskWorker(false)
	}

	for i, worker := range workers {
		if worker == nil {
			worker = newWorker(true, workerTypNone)
			workers[i] = worker
		}

		if worker.getTyp() == workerTypTask {
			taskWorkers = append(taskWorkers, worker.(*taskWorker))
		}

		go worker.run()
	}
}

func synceWorkers() {

	for _, worker := range workers {
		worker.sync()
	}
}

func releaseWorkers() {
	for _, worker := range workers {
		worker.syncDone()
	}
}
