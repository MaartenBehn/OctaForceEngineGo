package OctaForce

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	UVCord   mgl32.Vec2
}

// Mesh holds all data needed to render a 3D Object.
// When you change Vertices or Indices buy your self don't forget to set NeedsVertexUpdate to true. Otherwise the changes
// will not be applied.
type Mesh struct {
	Vertices          []vertex
	Indices           []uint32
	vao               uint32
	vertexVBO         uint32
	ebo               uint32
	Material          Material
	NeedsVertexUpdate bool

	instanceVBO         uint32
	instances           []*MeshInstant
	instanceData        []float32
	needsInstanceUpdate bool

	Transform *Transform
}

func NewMesh() *Mesh {
	return &Mesh{Transform: NewTransform()}
}

func renderMeshes() {
	for _, mesh := range globalActiveMeshesData.meshes {

		if len(mesh.instances) != 0 {
			continue
		}

		gl.BindVertexArray(mesh.vao)

		if mesh.NeedsVertexUpdate {
			pushVertexData(mesh)
		}

		matrix := mesh.Transform.getMatrix()
		// Transform
		gl.UniformMatrix4fv(2, 1, false, &matrix[0])

		// Color
		gl.Uniform3f(3, mesh.Material.DiffuseColor[0], mesh.Material.DiffuseColor[1], mesh.Material.DiffuseColor[2])

		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)

	}
}
func renderInstantMeshes() {
	for _, mesh := range globalActiveMeshesData.meshes {

		if len(mesh.instances) == 0 {
			continue
		}

		gl.BindVertexArray(mesh.vao)

		if mesh.NeedsVertexUpdate {
			pushVertexData(mesh)
		}

		pushInstanceData(mesh)

		gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil, int32(len(mesh.instances)+1))
	}
}

const vertexStride int32 = 3 * 4

func pushVertexData(mesh *Mesh) {
	var vertexData []float32
	for _, vertex := range mesh.Vertices {
		vertexData = append(vertexData, []float32{
			vertex.Position.X(),
			vertex.Position.Y(),
			vertex.Position.Z(),
		}...)
	}

	// vertex VBO
	gl.GenBuffers(1, &mesh.vertexVBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vertexVBO)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(0))

	// EBO
	gl.GenBuffers(1, &mesh.ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

	mesh.NeedsVertexUpdate = false
}

const instanceStride int32 = 19 * 4

func pushInstanceData(mesh *Mesh) {
	if mesh.needsInstanceUpdate {
		// Instance VBO
		gl.GenBuffers(1, &mesh.instanceVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.instanceVBO)
		gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.instances)+1)*int(instanceStride), gl.Ptr(nil), gl.DYNAMIC_DRAW)

		gl.EnableVertexAttribArray(1)
		gl.VertexAttribPointer(1, 3, gl.FLOAT, false, instanceStride, gl.PtrOffset(0))
		gl.VertexAttribDivisor(1, 1)

		gl.EnableVertexAttribArray(2)
		gl.VertexAttribPointer(2, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(3*4))
		gl.VertexAttribDivisor(2, 1)

		gl.EnableVertexAttribArray(3)
		gl.VertexAttribPointer(3, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(7*4))
		gl.VertexAttribDivisor(3, 1)

		gl.EnableVertexAttribArray(4)
		gl.VertexAttribPointer(4, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(11*4))
		gl.VertexAttribDivisor(4, 1)

		gl.EnableVertexAttribArray(5)
		gl.VertexAttribPointer(5, 4, gl.FLOAT, false, instanceStride, gl.PtrOffset(15*4))
		gl.VertexAttribDivisor(5, 1)

		mesh.instanceData = make([]float32, (1+len(mesh.instances))*19)

		mesh.needsInstanceUpdate = false
	}

	// Set Instance Data
	transform := mesh.Transform
	instanceData := mesh.instanceData

	instanceData[0] = mesh.Material.DiffuseColor[0]
	instanceData[1] = mesh.Material.DiffuseColor[1]
	instanceData[2] = mesh.Material.DiffuseColor[2]

	matrix := transform.getMatrix()
	instanceData[3] = matrix[0]
	instanceData[4] = matrix[1]
	instanceData[5] = matrix[2]
	instanceData[6] = matrix[3]

	instanceData[7] = matrix[4]
	instanceData[8] = matrix[5]
	instanceData[9] = matrix[6]
	instanceData[10] = matrix[7]

	instanceData[11] = matrix[8]
	instanceData[12] = matrix[9]
	instanceData[13] = matrix[10]
	instanceData[14] = matrix[11]

	instanceData[15] = matrix[12]
	instanceData[16] = matrix[13]
	instanceData[17] = matrix[14]
	instanceData[18] = matrix[15]

	for i, meshInstant := range mesh.instances {
		instantTransform := meshInstant.Transform
		index := (1 + i) * 19

		instanceData[index] = meshInstant.Material.DiffuseColor[0]
		instanceData[index+1] = meshInstant.Material.DiffuseColor[1]
		instanceData[index+2] = meshInstant.Material.DiffuseColor[2]

		matrix = instantTransform.getMatrix()
		instanceData[index+3] = matrix[0]
		instanceData[index+4] = matrix[1]
		instanceData[index+5] = matrix[2]
		instanceData[index+6] = matrix[3]

		instanceData[index+7] = matrix[4]
		instanceData[index+8] = matrix[5]
		instanceData[index+9] = matrix[6]
		instanceData[index+10] = matrix[7]

		instanceData[index+11] = matrix[8]
		instanceData[index+12] = matrix[9]
		instanceData[index+13] = matrix[10]
		instanceData[index+14] = matrix[11]

		instanceData[index+15] = matrix[12]
		instanceData[index+16] = matrix[13]
		instanceData[index+17] = matrix[14]
		instanceData[index+18] = matrix[15]
	}
	mesh.instanceData = instanceData

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.instanceVBO)
	gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.instances)+1)*int(instanceStride), gl.Ptr(instanceData), gl.DYNAMIC_DRAW)
}

// LoadOBJ returns the Mesh struct of the given OBJ file.
func (m *Mesh) LoadOBJ(path string, loadMaterials bool) {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	var vertices []mgl32.Vec3
	var normals []mgl32.Vec3
	var uvCord []mgl32.Vec2
	var faces [][3][3]uint32
	var material Material
	for _, line := range lines {
		values := strings.Split(line, " ")
		values[len(values)-1] = strings.Replace(values[len(values)-1], "\r", "", 1)

		switch values[0] {
		case "mtllib":
			if loadMaterials {
				material = LoadMtl("Mesh" + values[1])[0]
			}
			break
		case "v":
			vertices = append(vertices, mgl32.Vec3{ParseFloat(values[1]), ParseFloat(values[2]), ParseFloat(values[3])})
			break
		case "vn":
			normals = append(normals, mgl32.Vec3{ParseFloat(values[1]), ParseFloat(values[2]), ParseFloat(values[3])})
			break
		case "vt":
			uvCord = append(uvCord, mgl32.Vec2{ParseFloat(values[1]), ParseFloat(values[2])})
			break
		case "f":
			var face [3][3]uint32
			for j, value := range values {
				if j == 0 {
					continue
				}

				number := strings.Split(value, "/")
				face[j-1][0] = uint32(ParseInt(number[0]))
				face[j-1][1] = uint32(ParseInt(number[1]))
				face[j-1][2] = uint32(ParseInt(number[2]))
			}
			faces = append(faces, face)
			break
		}
	}

	m.Vertices = make([]vertex, len(vertices))
	m.Material = material
	for _, face := range faces {
		for _, values := range face {
			vertexIndex := values[0] - 1
			m.Indices = append(m.Indices, vertexIndex)
			//goland:noinspection GoNilness
			m.Vertices[vertexIndex].Position = vertices[vertexIndex]
			//Mesh.Vertices[vertexIndex].UVCord = uvCord[values[1] -1]
			//Mesh.Vertices[vertexIndex].Normal = normals[values[2] -1]
		}
	}

	m.NeedsVertexUpdate = true
}

type activeMeshesData struct {
	meshes []*Mesh
}

func initActiveMeshesData() {
	globalActiveMeshesData = &activeMeshesData{}
}

var globalActiveMeshesData *activeMeshesData

func GetActiveMeshes() *activeMeshesData {
	return globalActiveMeshesData
}

func (a *activeMeshesData) AddMesh(mesh *Mesh) {
	gl.GenVertexArrays(1, &mesh.vao)
	a.meshes = append(a.meshes, mesh)
}
func (a *activeMeshesData) RemoveMesh(mesh *Mesh) {

	for i := len(a.meshes) - 1; i >= 0; i-- {
		if a.meshes[i] == mesh {
			a.meshes = append(a.meshes[:i], a.meshes[i+1:]...)
		}
	}

	//unUsedVAOs = append(unUsedVAOs, mesh.vao)
}
