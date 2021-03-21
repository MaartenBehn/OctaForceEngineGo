package V2

type Data interface {
}

var datas []Data
var freeDataIds []int

func AddData(data Data) {
	workers[workerState].(*taskWorker).addTask(func() {
		if len(freeDataIds) > 0 {
			datas[freeDataIds[0]] = data
			freeDataIds = freeDataIds[1:]
		} else {

			datas = append(datas, data)
		}
	})
}
