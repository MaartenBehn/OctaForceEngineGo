package V2

import (
	"log"
	"runtime"
	"time"
)

var (
	running bool
	MaxFPS  float64
	Fps     float64
)

func Init(gameStartFunc func(), gameStopFunc func(), name string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	MaxFPS = 1
	running = true

	initFixedWorkers()

	gameStartFunc()

	run()

	gameStopFunc()
}

func run() {
	startTime := time.Now()
	var startDuration time.Duration
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		synceWorkers()

		updatePlans()

		releaseWorkers()

		log.Print(Fps)

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
