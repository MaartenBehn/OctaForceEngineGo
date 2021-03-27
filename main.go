package OctaForce

import (
	"log"
	"runtime"
)

func Init(gameStartFunc func()) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	runtime.LockOSThread()

	maxFPS = 60
	running = true

	initState()
	initActiveMeshesData()
	initActiveCamera()
	initDispatcher()

	InitRender()

	gameStartFunc()

	go runDispatcher()
	RunRender()
}
