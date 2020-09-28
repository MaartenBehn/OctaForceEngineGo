package OctaForceEngine

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

var (
	version string
	program uint32

	cameraEntityId         int
	cameraTransformUniform int32
	projectionUniform      int32

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
	projectionUniform = gl.GetUniformLocation(program, gl.Str("projection\x00"))
	cameraTransformUniform = gl.GetUniformLocation(program, gl.Str("camera\x00"))
	transformUniform = gl.GetUniformLocation(program, gl.Str("transform\x00"))

	// Output data Flag
	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0, 0, 0, 0)
}

// SetActiveCameraEntity sets the given entity as the camera. The given entity must have a camera component.
// This function does not check that, so be careful.
func SetActiveCameraEntity(entityId int) {
	cameraEntityId = entityId
}

func renderRenderer() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(program)

	cameraTransform := GetComponent(cameraEntityId, ComponentTransform).(Transform)
	// Creating inverted Camera pos
	view := cameraTransform.matrix.Inv()
	gl.UniformMatrix4fv(cameraTransformUniform, 1, false, &view[0])

	camera := GetComponent(cameraEntityId, ComponentCamera).(Camera)
	gl.UniformMatrix4fv(projectionUniform, 1, false, &camera.projection[0])

	entities := GetAllEntitiesWithComponent(ComponentMesh)
	for _, entity := range entities {
		renderEntity(entity)
	}

	gl.BindVertexArray(0)

	deleteUnUsedVAOs()
}

const stride int32 = 6 * 4

func renderEntity(entityId int) {
	mesh := GetComponent(entityId, ComponentMesh).(Mesh)

	if mesh.NeedsRenderUpdate {
		var vertexData []float32
		for _, vertex := range mesh.Vertices {
			vertexData = append(vertexData, []float32{
				vertex.Position.X(),
				vertex.Position.Y(),
				vertex.Position.Z(),

				mesh.Material.DiffuseColor.X(),
				mesh.Material.DiffuseColor.Y(),
				mesh.Material.DiffuseColor.Z(),
			}...)
		}

		gl.GenVertexArrays(1, &mesh.vao)
		gl.BindVertexArray(mesh.vao)

		gl.GenBuffers(1, &mesh.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

		gl.GenBuffers(1, &mesh.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

		vertexPositionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertexPosition\x00")))
		gl.EnableVertexAttribArray(vertexPositionAttrib)
		gl.VertexAttribPointer(vertexPositionAttrib, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))

		vertexColorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertexColor\x00")))
		gl.EnableVertexAttribArray(vertexColorAttrib)
		gl.VertexAttribPointer(vertexColorAttrib, 3, gl.FLOAT, false, stride, gl.PtrOffset(3*4))

		mesh.NeedsRenderUpdate = false

		SetComponent(entityId, ComponentMesh, mesh)

	} else {
		gl.BindVertexArray(mesh.vao)
	}

	transform := GetComponent(entityId, ComponentTransform).(Transform)
	gl.UniformMatrix4fv(transformUniform, 1, false, &transform.matrix[0])

	gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
}

var unUsedVAOs []uint32

func deleteUnUsedVAOs() {
	for _, vao := range unUsedVAOs {
		gl.DeleteVertexArrays(1, &vao)
	}
	unUsedVAOs = []uint32{}
}
