package OctaForce

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
)

var (
	version string
	program uint32
)

type programmData struct {
	vertexPath   string
	fragmentPath string
	id           uint32
	renderFunc   func()
}

var programmDatas []programmData

func initRenderer() {

	var err error
	// Initialize Gl
	if err = gl.Init(); err != nil {
		panic(err)
	}
	version = gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	gl.Enable(gl.DEBUG_OUTPUT)
	gl.DebugMessageCallback(nil, nil)

	programmDatas = make([]programmData, 2)
	programmDatas[0] = programmData{
		vertexPath:   "/shader/vertexShader.shader",
		fragmentPath: "/shader/fragmentShader.shader",
		renderFunc:   renderMeshes,
	}
	programmDatas[1] = programmData{
		vertexPath:   "/shader/vertexShaderInstancing.shader",
		fragmentPath: "/shader/fragmentShader.shader",
		renderFunc:   renderInstantMeshes,
	}

	for i, programmData := range programmDatas {
		// Configure the vertex and fragment shaders
		vertexShader, err := compileShader(absPath+programmData.vertexPath, gl.VERTEX_SHADER)
		if err != nil {
			panic(err)
		}
		fragmentShader, err := compileShader(absPath+programmData.fragmentPath, gl.FRAGMENT_SHADER)
		if err != nil {
			panic(err)
		}

		program = gl.CreateProgram()
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

			panic(fmt.Errorf("failed to link program: %v", log))
		}
		gl.DeleteShader(vertexShader)
		gl.DeleteShader(fragmentShader)

		// Output data Flag
		gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

		programmData.id = program
		programmDatas[i] = programmData
	}

	// Global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 0)
}

func runRender() {
	*engineTasks[RenderTask] = *NewTask(func() {
		renderRenderer()
		renderGui()
		renderWindow()
		printGlErrors()
	})
	engineTasks[RenderTask].SetRepeating(true)
	addTask <- engineTasks[RenderTask]
	engineTasks[RenderTask].run()
}

func renderRenderer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	for _, programmData := range programmDatas {
		gl.UseProgram(programmData.id)

		// Creating inverted Camera pos
		view := ActiveCamera.Transform.getMatrix().Inv()
		gl.UniformMatrix4fv(1, 1, false, &view[0])
		gl.UniformMatrix4fv(0, 1, false, &ActiveCamera.projection[0])

		programmData.renderFunc()

	}

	gl.BindVertexArray(0)
	deleteUnUsedVAOs()
}

var unUsedVAOs []uint32

func deleteUnUsedVAOs() {
	for _, vao := range unUsedVAOs {
		gl.DeleteVertexArrays(1, &vao)
	}
	unUsedVAOs = []uint32{}
}
func printGlErrors() {
	glerror := gl.GetError()
	if glerror == gl.NO_ERROR {
		return
	}

	fmt.Printf("GLError ")

	switch glerror {
	case gl.INVALID_ENUM:
		fmt.Printf("GL_INVALID_ENUM")
	case gl.INVALID_VALUE:
		fmt.Printf("GL_INVALID_VALUE")
	case gl.INVALID_OPERATION:
		fmt.Printf("GL_INVALID_OPERATION")
	case gl.STACK_OVERFLOW:
		fmt.Printf("GL_STACK_OVERFLOW")
	case gl.STACK_UNDERFLOW:
		fmt.Printf("GL_STACK_UNDERFLOW")
	case gl.OUT_OF_MEMORY:
		fmt.Printf("GL_OUT_OF_MEMORY")
	default:
		fmt.Printf("%d", glerror)
	}

	fmt.Printf("\n")
}
