package OctaForce

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"path/filepath"
	"runtime"
	"time"
)

var absPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	absPath = filepath.Dir(b)
}

func StartUp(gameStartUpFunc func(), gameStopFunc func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	maxFPS = 60
	maxUPS = 30
	running = true
	needAllMeshUpdate = false

	// Setting up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	startUpWindow()
	setUpComponentTables()

	gameStartUpFunc()

	go runUpdate()
	runRender()

	gameStopFunc()
}

var running bool

var fps float64

func GetFPS() float64 {
	return fps
}
func GetCappedFPS() float64 {
	if fps > maxFPS {
		fps = maxFPS
	}
	return fps
}

var maxFPS float64

func SetMaxFPS(maxfps float64) {
	maxFPS = maxfps
}

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

func GetUPS() float64 {
	return ups
}
func GetCappedUPS() float64 {
	if ups > maxUPS {
		ups = maxUPS
	}
	return ups
}

var maxUPS float64

func SetMaxUPS(maxups float64) {
	maxUPS = maxups
}

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

var gameUpadteFuncs [20]func()

func AddUpdateCallback(newGameUpdatefunc func()) int {
	for i, gameUpdatefunc := range gameUpadteFuncs {
		if gameUpdatefunc == nil {
			gameUpdatefunc = newGameUpdatefunc
			return i
		}
	}

	return -1
}

func RemoveUpdateUpCallback(i int) {
	gameUpadteFuncs[i] = nil
}

func update() {
	updateAllComponents()

	if needAllMeshUpdate {
		updateAllMeshData()
	}

	for _, gameUpadteFunc := range gameUpadteFuncs {
		if gameUpadteFunc != nil {
			gameUpadteFunc()
		}
	}
}
