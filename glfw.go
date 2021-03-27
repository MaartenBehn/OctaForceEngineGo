package OctaForce

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/inkyblackness/imgui-go"
	"math"
	"runtime"
)

var (
	window *glfw.Window

	WindowWidth  = 1280
	WindowHeight = 720

	mouseButtonPrimary   = 0
	mouseButtonSecondary = 1
	mouseButtonTertiary  = 2
	mouseButtonCount     = 3

	lastTime         float64
	mouseJustPressed [3]bool
)

func initGLFW() {
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		panic(err)
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, 1)

	window, err = glfw.CreateWindow(WindowWidth, WindowHeight, "Test", nil, nil)
	if err != nil {
		glfw.Terminate()
		panic(err)
	}

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
}

var glfwButtonIndexByID = map[glfw.MouseButton]int{
	glfw.MouseButton1: mouseButtonPrimary,
	glfw.MouseButton2: mouseButtonSecondary,
	glfw.MouseButton3: mouseButtonTertiary,
}

var glfwButtonIDByIndex = map[int]glfw.MouseButton{
	mouseButtonPrimary:   glfw.MouseButton1,
	mouseButtonSecondary: glfw.MouseButton2,
	mouseButtonTertiary:  glfw.MouseButton3,
}

func processEvents() {
	glfw.PollEvents()
}

func DisplaySize() [2]float32 {
	w, h := window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

func FramebufferSize() [2]float32 {
	w, h := window.GetFramebufferSize()
	return [2]float32{float32(w), float32(h)}
}

func newFrame() {
	// Setup display size (every frame to accommodate for window resizing)
	displaySize := DisplaySize()
	gui.io.SetDisplaySize(imgui.Vec2{X: displaySize[0], Y: displaySize[1]})

	// Setup lastTime step
	currentTime := glfw.GetTime()
	if lastTime > 0 {
		gui.io.SetDeltaTime(float32(currentTime - lastTime))
	}
	lastTime = currentTime

	// Setup inputs
	if window.GetAttrib(glfw.Focused) != 0 {
		x, y := window.GetCursorPos()
		gui.io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})
	} else {
		gui.io.SetMousePosition(imgui.Vec2{X: -math.MaxFloat32, Y: -math.MaxFloat32})
	}

	for i := 0; i < len(mouseJustPressed); i++ {
		down := mouseJustPressed[i] || (window.GetMouseButton(glfwButtonIDByIndex[i]) == glfw.Press)
		gui.io.SetMouseButtonDown(i, down)
		mouseJustPressed[i] = false
	}
}

func postRender() {
	window.SwapBuffers()
}
