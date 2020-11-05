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

const vertexStride int32 = 3 * 4
const instanceStride int32 = 19 * 4

func renderEntity(entityId int) {
	mesh := GetComponent(entityId, ComponentMesh).(Mesh)
	transform := GetComponent(entityId, ComponentTransform).(Transform)

	if mesh.needNewBuffer {
		var vertexData []float32
		for _, vertex := range mesh.Vertices {
			vertexData = append(vertexData, []float32{
				vertex.Position.X(),
				vertex.Position.Y(),
				vertex.Position.Z(),
			}...)
		}

		// VAO
		gl.GenVertexArrays(1, &mesh.vao)
		gl.BindVertexArray(mesh.vao)

		// Vertex VBO
		gl.GenBuffers(1, &mesh.vertexVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vertexVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

		vertexPositionAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertexPosition\x00")))
		gl.EnableVertexAttribArray(vertexPositionAttrib)
		gl.VertexAttribPointer(vertexPositionAttrib, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(0))

		// EBO
		gl.GenBuffers(1, &mesh.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

		// Instance VBO
		gl.GenBuffers(1, &mesh.instanceVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.instanceVBO)
		gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.Instants)+1)*int(instanceStride), gl.Ptr(nil), gl.DYNAMIC_DRAW)

		colorAttrib := uint32(gl.GetAttribLocation(program, gl.Str("instanceColor\x00")))
		gl.EnableVertexAttribArray(colorAttrib)
		gl.VertexAttribPointer(colorAttrib, 3, gl.FLOAT, false, instanceStride, gl.PtrOffset(0))
		gl.VertexAttribDivisor(colorAttrib, 1)

		transformXAttrib := uint32(gl.GetAttribLocation(program, gl.Str("transformX\x00")))
		gl.EnableVertexAttribArray(transformXAttrib)
		gl.VertexAttribPointer(transformXAttrib, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(3*4))
		gl.VertexAttribDivisor(transformXAttrib, 1)

		transformYAttrib := uint32(gl.GetAttribLocation(program, gl.Str("transformY\x00")))
		gl.EnableVertexAttribArray(transformYAttrib)
		gl.VertexAttribPointer(transformYAttrib, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(7*4))
		gl.VertexAttribDivisor(transformYAttrib, 1)

		transformZAttrib := uint32(gl.GetAttribLocation(program, gl.Str("transformZ\x00")))
		gl.EnableVertexAttribArray(transformZAttrib)
		gl.VertexAttribPointer(transformZAttrib, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(11*4))
		gl.VertexAttribDivisor(transformZAttrib, 1)

		transformSAttrib := uint32(gl.GetAttribLocation(program, gl.Str("transformS\x00")))
		gl.EnableVertexAttribArray(transformSAttrib)
		gl.VertexAttribPointer(transformSAttrib, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(15*4))
		gl.VertexAttribDivisor(transformSAttrib, 1)

		mesh.needNewBuffer = false

	} else {
		gl.BindVertexArray(mesh.vao)
	}

	var instanceData = []float32{
		mesh.Material.DiffuseColor[0],
		mesh.Material.DiffuseColor[1],
		mesh.Material.DiffuseColor[2],

		transform.matrix[0],
		transform.matrix[1],
		transform.matrix[2],
		transform.matrix[3],

		transform.matrix[4],
		transform.matrix[5],
		transform.matrix[6],
		transform.matrix[7],

		transform.matrix[8],
		transform.matrix[9],
		transform.matrix[10],
		transform.matrix[11],

		transform.matrix[12],
		transform.matrix[13],
		transform.matrix[14],
		transform.matrix[15],
	}

	for _, id := range mesh.Instants {
		meshInstant := GetComponent(id, ComponentMeshInstant).(MeshInstant)
		instantTransform := GetComponent(id, ComponentTransform).(Transform)

		instanceData = append(instanceData, []float32{
			meshInstant.Material.DiffuseColor[0],
			meshInstant.Material.DiffuseColor[1],
			meshInstant.Material.DiffuseColor[2],

			instantTransform.matrix[0],
			instantTransform.matrix[1],
			instantTransform.matrix[2],
			instantTransform.matrix[3],

			instantTransform.matrix[4],
			instantTransform.matrix[5],
			instantTransform.matrix[6],
			instantTransform.matrix[7],

			instantTransform.matrix[8],
			instantTransform.matrix[9],
			instantTransform.matrix[10],
			instantTransform.matrix[11],

			instantTransform.matrix[12],
			instantTransform.matrix[13],
			instantTransform.matrix[14],
			instantTransform.matrix[15],
		}...)
	}

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.instanceVBO)
	gl.BufferSubData(gl.ARRAY_BUFFER, 0, (len(mesh.Instants)+1)*int(instanceStride), gl.Ptr(instanceData))

	SetComponent(entityId, ComponentMesh, mesh)

	gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil, int32(len(mesh.Instants)+1))
}

var unUsedVAOs []uint32

func deleteUnUsedVAOs() {
	for _, vao := range unUsedVAOs {
		gl.DeleteVertexArrays(1, &vao)
	}
	unUsedVAOs = []uint32{}
}
