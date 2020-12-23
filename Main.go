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

var gameUpdateFunction func()

// StartUp needs to be called in the game main function. It requires a game start function and stop function.
// The game start function is called after StartUp but before the update calls. So do here all initial game engine setup.
// The game stop function is called when the game stops. So do here all stuff you need to do when the game stops.
func StartUp(gameStartUpFunc func(), gameUpdateFunc func(), gameStopFunc func(), name string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	// setting var values
	gameUpdateFunction = gameUpdateFunc
	MaxFPS = 60
	MaxUPS = 30
	running = true

	// Initialising vars
	// Setting up glfw
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	// internal setup calls
	setUpWindow(name)
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
func GetFPS() float64 {
	return fps
}

// GetCappedFPS returns the current frames per second capped to max fps value set.
func GetCappedFPS() float64 {
	if fps > MaxFPS {
		fps = MaxFPS
	}
	return fps
}

var MaxFPS float64

func runRender() {
	startTime := time.Now()
	var startDuration time.Duration
	wait := time.Duration(1.0 / MaxFPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		// All render Calls
		renderRenderer()
		renderWindow()
		printGlErrors("RenderLoop")

		diff := time.Since(startTime) - startDuration
		if diff > 0 {
			fps = (wait.Seconds() / diff.Seconds()) * MaxFPS
		} else {
			fps = 10000
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
	if ups > MaxUPS {
		ups = MaxUPS
	}
	return ups
}

var MaxUPS float64

var updateDeltaTime float64

// GetDeltaTime return the Time (Sek) passed since last Update.
func GetDeltaTime() float64 {
	return updateDeltaTime
}
func runUpdate() {
	startTime := time.Now()
	var startDuration time.Duration
	wait := time.Duration(1.0 / MaxUPS * 1000000000)

	for running {
		startDuration = time.Since(startTime)

		// All update Calls
		updateWindow()
		updateAllComponents()
		gameUpdateFunction()

		diff := time.Since(startTime) - startDuration
		if diff > 0 {
			ups = (wait.Seconds() / diff.Seconds()) * MaxUPS
		} else {
			ups = 10000
		}
		if diff < wait {
			updateDeltaTime = wait.Seconds()
			time.Sleep(wait - diff)
		} else {
			updateDeltaTime = diff.Seconds()
		}
	}
}
