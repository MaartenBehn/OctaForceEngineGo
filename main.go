package V2

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

func Init(gameStartFunc func(), gameStopFunc func(), name string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	MaxFPS = 60
	running = true
	windowName = name

	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	initWorkers()
	initWindow()
	initRenderer()

	gameStartFunc()

	go runDispatcher()
	workers[workerRender].run()

	gameStopFunc()
}
