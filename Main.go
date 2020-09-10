package OctaForce

import "time"

func startUp(gameStartUpFunc func(), gameStopFunc func()) {
	maxFPS = 60
	maxUPS = 30

	startUpWindow()

	gameStartUpFunc()

	go runUpdate()
	runRender()

	gameStopFunc()
}

var running bool

var fps float64
var maxFPS float64

func runRender() {
	var startTime = time.Now()
	var startDuration time.Duration
	var wait = time.Duration(1 / maxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		updateWindow()

		var diff = time.Since(startTime) - startDuration
		fps = (wait / diff).Seconds() * maxFPS
		if diff < wait {
			time.Sleep(wait - diff)
		}
	}
}

var ups float64
var maxUPS float64

func runUpdate() {
	var startTime = time.Now()
	var startDuration time.Duration
	var wait = time.Duration(1 / maxUPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		update()

		var diff = time.Since(startTime) - startDuration
		ups = (wait / diff).Seconds() * maxFPS
		if diff < wait {
			time.Sleep(wait - diff)
		}
	}
}

var gameUpadteFuncs []func()

func addUpdateCallback(gameUpdatefunc func()) int {
	return AddFuncToSlice(&gameUpadteFuncs, gameUpdatefunc)
}
func removeUpdateUpCallback(i int) {
	RemoveFuncFromSlice(&gameUpadteFuncs, i, false)
}

func update() {
	for _, gameUpadteFunc := range gameUpadteFuncs {
		gameUpadteFunc()
	}
}
