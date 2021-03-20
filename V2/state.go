package V2

type state struct {
	datas []Data
}

func (s *state) copy() state {
	copy := state{datas: make([]Data, len(s.datas))}
	for i, d := range s.datas {
		copy.datas[i] = d.copy()
	}
	return copy
}

var readState state
var writeState state

func syncState() {
	readState = writeState.copy()
}

func GetData(id int) Data {
	return readState.datas[id]
}

func SetData(id int, data Data) {
	writeState.datas[id] = data
}

var freeDataIds []int

func AddData(data Data) {
	globalChanges <- func() {
		if len(freeDataIds) > 0 {
			data.SetId(freeDataIds[0])
			writeState.datas[data.GetId()] = data
			freeDataIds = freeTaskIds[1:]
		} else {
			data.SetId(len(writeState.datas))
			writeState.datas = append(writeState.datas, data)
		}
	}
}
func RemoveData(data Data) {
	globalChanges <- func() {
		writeState.datas[data.GetId()] = nil
		freeDataIds = append(freeDataIds, data.GetId())
	}
}
