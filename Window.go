package OctaForce

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

var (
	window  *glfw.Window
	version string

	program uint32
)

func startUpWindow() {
	// Takes in mainStartUpFunc function wich will bw called at the end of this so it has the correct glfw context

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

	// Initialize Gl
	if err = gl.Init(); err != nil {
		panic(err)
	}

	version = gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err = newProgram()
	if err != nil {
		panic(err)
	}
	gl.UseProgram(program)

	mesh := Mesh{
		Vertices: []Vertex{
			{
				Position: mgl32.Vec3{-0.5, -0.5, 0.0},
			},
			{
				Position: mgl32.Vec3{0.5, -0.5, 0.0},
			},
			{
				Position: mgl32.Vec3{0.0, 0.5, 0.0},
			},
		},
	}
	mesh.updateMeshData()
	activeMeshes = append(activeMeshes, mesh)
	updateAllMeshData()

	// Configure global settings
	gl.ClearColor(0, 0, 0, 0)
}

func newProgram() (uint32, error) {

	vertexShader, err := compileShader("H:\\dev\\Go\\src\\OctaForceEngineGo\\shader\\vertexShader.shader", gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader("H:\\dev\\Go\\src\\OctaForceEngineGo\\shader\\fragmentShader.shader", gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func updateWindow() {
	gl.ClearColor(0, 0, 0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	renderMeshes()

	// Maintenance
	window.SwapBuffers()
	glfw.PollEvents()

	if window.ShouldClose() {
		running = false
	}
}
