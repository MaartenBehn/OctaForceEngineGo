package OctaForce

import (
	"log"
	"runtime"
)

func Init(gameStartFunc func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	maxUPS = 60
	maxFPS = 30
	running = true

	initState()
	initActiveMeshesData()
	initActiveCamera()
	initDispatcher()

	initRender()

	gameStartFunc()

	go runDispatcher()
	runRender()
}
