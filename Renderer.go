package OctaForce

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"strings"
)

var (
	version string
	program uint32

	projection        mgl32.Mat4
	projectionUniform int32

	camera        mgl32.Mat4
	cameraUniform int32

	glTransform      mgl32.Mat4
	transformUniform int32
)

func setUpRenderer() {
	var err error
	// Initialize Gl
	if err = gl.Init(); err != nil {
		panic(err)
	}
	version = gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	vertexShader, err := compileShader(absPath+"/shader/vertexShader.shader", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	fragmentShader, err := compileShader(absPath+"/shader/fragmentShader.shader", gl.FRAGMENT_SHADER)
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
	gl.UseProgram(program)

	// Perspective Projection matrix
	projection = mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	// Camera matrix
	camera = mgl32.LookAtV(mgl32.Vec3{0, 0, 10}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	// Object glTransform matrix
	glTransform = mgl32.Ident4()
	transformUniform = gl.GetUniformLocation(program, gl.Str("transform\x00"))
	gl.UniformMatrix4fv(transformUniform, 1, false, &glTransform[0])

	// Output Data Flag
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Global settings
	gl.ClearColor(0, 0, 0, 0)
}

func updateRenderer() {
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.UseProgram(program)
	renderMeshes()
}
