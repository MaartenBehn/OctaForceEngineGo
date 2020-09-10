package OctaForce

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"runtime"
	"time"
)

func startUp(gameStartUpFunc func(), gameStopFunc func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	maxFPS = 60
	maxUPS = 30
	running = true

	gameUpadteFuncs = make([]func(), 0)

	// Setting up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

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
	var wait = time.Duration(1.0 / maxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		updateWindow()

		var diff = time.Since(startTime) - startDuration

		if diff > 0 {
			fps = (wait / diff).Seconds() * maxFPS
		} else {
			fps = maxFPS
		}

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
	var wait = time.Duration(1.0 / maxUPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		update()

		var diff = time.Since(startTime) - startDuration

		if diff > 0 {
			ups = (wait / diff).Seconds() * maxUPS
		} else {
			ups = maxUPS
		}

		if diff < wait {
			time.Sleep(wait - diff)
		}
	}
}

var gameUpadteFuncs []func()

func addUpdateCallback(gameUpdatefunc func()) int {
	var index int
	gameUpadteFuncs, index = AddFuncToSlice(gameUpadteFuncs, gameUpdatefunc)
	return index
}
func removeUpdateUpCallback(i int) {
	RemoveFuncFromSlice(&gameUpadteFuncs, i, false)
}

func update() {
	for _, gameUpadteFunc := range gameUpadteFuncs {
		gameUpadteFunc()
	}
}
