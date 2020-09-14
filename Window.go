package OF

import (
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	windowWidth  = 1920
	windowHeight = 1080
)

var window *glfw.Window

func setUpWindow() {
	var err error

	// Setting up Window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err = glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
}

func updateWindow() {
	window.SwapBuffers()
	glfw.PollEvents()
	if window.ShouldClose() {
		running = false
	}
}
