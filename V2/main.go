package V2

import "time"

var (
	running bool
	MaxFPS  float64
	Fps     float64
)

func Init() {
	initChangesBuffer()

	run()
}

func run() {
	startTime := time.Now()
	var startDuration time.Duration
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		runPlan()
		applyChangeBuffer()

		diff := time.Since(startTime) - startDuration
		if diff > 0 {
			Fps = (wait.Seconds() / diff.Seconds()) * MaxFPS
		} else {
			Fps = 10000
		}
		if diff < wait {
			time.Sleep(wait - diff)
		}
	}

}
