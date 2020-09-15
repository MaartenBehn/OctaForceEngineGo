package OctaForceEngine

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

// StartUp needs to be called in the game main function. It requires a game start function and stop function.
// The game start function is called after StartUp but before the update calls. So do here all initial game engine setup.
// The game stop function is called when the game stops. So do here all stuff you need to do when the game stops.
func StartUp(gameStartUpFunc func(), gameStopFunc func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	// setting var values
	maxFPS = 60
	maxUPS = 30
	running = true

	// Initialising vars
	// Setting up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()
	gameUpdateFuncMap = map[int]func(){}

	// internal setup calls
	setUpWindow()
	setUpRenderer()
	setUpComponentTables()

	// parsed game setup call
	gameStartUpFunc() //(I need to do it in that way because somehow the glfw context only applies if a func is in the stack higher than the init call.)

	go runUpdate() // Running the update calls on sprat thread.
	runRender()    // The render calls needs run on the main tread so the glfw init call system still applies.

	gameStopFunc()
}

var running bool

var fps float64

// GetFPS returns the current frames per second.
// 0 is an edge case and can mean that they are 0 or actually infinite.
func GetFPS() float64 {
	return fps
}

// GetCappedFPS returns the current frames per second capped to max fps value set.
// 0 is an edge case and can mean that they are 0 or actually infinite.
func GetCappedFPS() float64 {
	if fps > maxFPS {
		fps = maxFPS
	}
	return fps
}

var maxFPS float64

func SetMaxFPS(_maxFPS float64) {
	maxFPS = _maxFPS
}
func runRender() {
	var startTime = time.Now()
	var startDuration time.Duration
	var wait = time.Duration(1.0 / maxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		// All render Calls
		updateRenderer()
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

// GetUPS returns the current updates per second.
// 0 is an edge case and can mean that they are 0 or actually infinite.
func GetUPS() float64 {
	return ups
}

// GetCappedUPS returns the current updates per second capped to max ups value set.
// 0 is an edge case and can mean that they are 0 or actually infinite.
func GetCappedUPS() float64 {
	if ups > maxUPS {
		ups = maxUPS
	}
	return ups
}

var maxUPS float64

func SetMaxUPS(_maxUPS float64) {
	maxUPS = _maxUPS
}
func runUpdate() {
	var startTime = time.Now()
	var startDuration time.Duration
	var wait = time.Duration(1.0 / maxUPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		// All update Calls
		updateAllComponents()
		performGameUpdateFunctions()

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

var gameUpdateFuncMap map[int]func()
var gameUpdateFuncCounter int

// AddUpdateCallback adds the given function to a map with a random int id.
// The given function will be called every engine update.
// The returned int is the id of the function in the map.
func AddUpdateCallback(newGameUpdateFunc func()) int {
	gameUpdateFuncCounter++
	gameUpdateFuncMap[gameUpdateFuncCounter] = newGameUpdateFunc
	return gameUpdateFuncCounter
}

// RemoveUpdateCallback removes the function at the given int id.
func RemoveUpdateCallback(id int) {
	gameUpdateFuncMap[id] = nil
}
func performGameUpdateFunctions() {
	for _, gameUpdateFunc := range gameUpdateFuncMap {
		gameUpdateFunc()
	}
}
