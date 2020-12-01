package OctaForceEngine

import (
	"fmt"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
)

var (
	version string
	program uint32

	cameraEntityId         int
	cameraTransformUniform int32
	projectionUniform      int32
)

type ProgrammData struct {
	vertexPath   string
	fragmentPath string
	id           int
}

var programmDatas []ProgrammData

func setUpRenderer() {

	var err error
	// Initialize Gl
	if err = gl.Init(); err != nil {
		panic(err)
	}
	version = gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

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

		// Perspective Projection matrix
		projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
		cameraTransformUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))

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

// SetActiveCameraEntity sets the given entity as the camera. The given entity must have a camera component.
// This function does not check that, so be careful.
func SetActiveCameraEntity(entityId int) {
	if HasComponent(entityId, ComponentCamera) {
		cameraEntityId = entityId
	}
}

func renderRenderer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	cameraTransform := GetComponent(cameraEntityId, ComponentTransform).(Transform)
	// Creating inverted Camera pos
	view := cameraTransform.matrix.Inv()
	gl.UniformMatrix4fv(cameraTransformUniform, 1, false, &view[0])

	camera := GetComponent(cameraEntityId, ComponentCamera).(Camera)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &camera.projection[0])

	entities := GetAllEntitiesWithComponent(ComponentMesh)
	for _, entity := range entities {
		renderMesh(entity)
	}

	gl.BindVertexArray(0)

	deleteUnUsedVAOs()
}

const vertexStride int32 = 3 * 4
const instanceStride int32 = 19 * 4

var unUsedVAOs []uint32

func deleteUnUsedVAOs() {
	for _, vao := range unUsedVAOs {
		gl.DeleteVertexArrays(1, &vao)
	}
	unUsedVAOs = []uint32{}
}
