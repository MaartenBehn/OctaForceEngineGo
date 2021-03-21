package V2

import (
	"log"
	"runtime"
	"time"
)

var (
	running    bool
	MaxFPS     float64
	Fps        float64
	frameStart time.Time
	wait       time.Duration
)

func Init(gameStartFunc func(), gameStopFunc func(), name string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	MaxFPS = 60
	running = true

	initFixedWorkers()

	gameStartFunc()

	run()

	gameStopFunc()
}

func run() {
	wait = time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		frameStart = time.Now()

		synceWorkers()

		updatePlans()

		releaseWorkers()

		timeLeft, _ := frameEnd()
		time.Sleep(timeLeft)
	}
}

func frameEnd() (timeLeft time.Duration, fps float64) {
	diff := time.Since(frameStart)
	if diff > 0 {
		fps = (wait.Seconds() / diff.Seconds()) * MaxFPS
	} else {
		fps = 10000
	}

	if diff < wait {
		timeLeft = wait - diff
	}
	return timeLeft, fps
}
