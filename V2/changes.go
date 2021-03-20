package V2

var (
	globalChanges        chan func()
	applyChangesFlag     chan bool
	applyChangesDoneFlag chan bool
)

func initChangesBuffer() {
	globalChanges = make(chan func(), 1)
	applyChangesFlag = make(chan bool)
	applyChangesDoneFlag = make(chan bool)

	go runChangeBuffer()
}

func runChangeBuffer() {
	var stateChangeList []func()

	for {
		select {
		case change := <-globalChanges:
			stateChangeList = append(stateChangeList, change)
		case <-applyChangesFlag:

			for i := 0; i < len(stateChangeList); i++ {
				stateChangeList[i]()
			}
			applyChangesDoneFlag <- true
		}
	}
}

func applyChangeBuffer() {
	applyChangesFlag <- true
	<-applyChangesDoneFlag
}
