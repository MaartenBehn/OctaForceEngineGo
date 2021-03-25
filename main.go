package OctaForce

import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"log"
	"path/filepath"
	"runtime"
)

var absPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	absPath = filepath.Dir(b)
}

var stopCallback func()

func SetStopCallback(function func()) {
	stopCallback = function
}

func Init(start func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	maxFPS = 60
	running = true

	initState()
	initActiveMeshesData()
	initActiveCamera()
	initDispatcher()

	initWindow()
	initRenderer()
	//initGui()

	start()

	go runDispatcher()
	runRender()

	if stopCallback != nil {
		stopCallback()
	}

	glfw.Terminate()
	gui.context.Destroy()
}
