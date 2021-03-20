package OctaForceEngine

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
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
	instances           []int
	needsInstanceUpdate bool
}

func setUpMesh(data interface{}, entityId int) interface{} {
	var mesh Mesh
	if data == nil {
		mesh = Mesh{}
	} else {
		mesh = data.(Mesh)
	}

	gl.GenVertexArrays(1, &mesh.vao)
	return mesh
}
func deleteMesh(component interface{}, entityId int) interface{} {
	mesh := component.(Mesh)
	unUsedVAOs = append(unUsedVAOs, mesh.vao)
	return nil
}

func renderMeshes() {
	entities := world.getAllEntityIdsWithComponent(ComponentMesh)

	for _, entityId := range entities {

		mesh := GetComponent(entityId, ComponentMesh).(Mesh)
		if len(mesh.instances) != 0 {
			continue
		}

		gl.BindVertexArray(mesh.vao)

		if mesh.NeedsVertexUpdate {
			pushVertexData(mesh)
			setComponentInternal(entityId, ComponentMesh, mesh)
		}

		// Transform
		transform := GetComponent(entityId, ComponentTransform).(Transform)
		gl.UniformMatrix4fv(2, 1, false, &transform.matrix[0])

		// Color
		gl.Uniform3f(3, mesh.Material.DiffuseColor[0], mesh.Material.DiffuseColor[1], mesh.Material.DiffuseColor[2])

		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)

	}
}
func renderInstantMeshes() {
	entities := world.getAllEntityIdsWithComponent(ComponentMesh)

	for _, entityId := range entities {

		mesh := GetComponent(entityId, ComponentMesh).(Mesh)
		if len(mesh.instances) == 0 {
			continue
		}

		gl.BindVertexArray(mesh.vao)

		if mesh.NeedsVertexUpdate {
			pushVertexData(mesh)
			setComponentInternal(entityId, ComponentMesh, mesh)
		}

		pushInstanceData(mesh, entityId)

		gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil, int32(len(mesh.instances)+1))
	}
}

const vertexStride int32 = 3 * 4

func pushVertexData(mesh Mesh) {
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

func pushInstanceData(mesh Mesh, entityId int) {
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

		mesh.needsInstanceUpdate = false
		SetComponent(entityId, ComponentMesh, mesh)
	}

	// Set Instance Data
	transform := GetComponent(entityId, ComponentTransform).(Transform)
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
	for _, id := range mesh.instances {
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
	gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.instances)+1)*int(instanceStride), gl.Ptr(instanceData), gl.DYNAMIC_DRAW)
}

// LoadOBJ returns the mesh struct of the given OBJ file.
func (mesh *Mesh) LoadOBJ(path string, loadMaterials bool) {

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
				material = LoadMtl("mesh" + values[1])[0]
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

	mesh.Vertices = make([]vertex, len(vertices))
	mesh.Material = material
	for _, face := range faces {
		for _, values := range face {
			vertexIndex := values[0] - 1
			mesh.Indices = append(mesh.Indices, vertexIndex)
			//goland:noinspection GoNilness
			mesh.Vertices[vertexIndex].Position = vertices[vertexIndex]
			//mesh.Vertices[vertexIndex].UVCord = uvCord[values[1] -1]
			//mesh.Vertices[vertexIndex].Normal = normals[values[2] -1]
		}
	}

	mesh.NeedsVertexUpdate = true
}

// MeshInstant is a Component that creates an Instant of the Mesh Component of MeshEntity.
// Transform and Material are individual changeable.
type MeshInstant struct {
	MeshEntity int
	Material   Material

	ownEntity          int
	currentlySetEntity int
}

func setUpMeshInstant(_ interface{}, entityId int) interface{} {
	return MeshInstant{
		ownEntity: entityId,
	}
}
func addMeshInstant(component interface{}, entityId int) interface{} {
	meshInstant := component.(MeshInstant)

	if meshInstant.ownEntity == 0 || meshInstant.MeshEntity == 0 {
		return component
	}

	if meshInstant.currentlySetEntity != meshInstant.MeshEntity {

		if HasComponent(meshInstant.currentlySetEntity, ComponentMesh) {
			mesh := GetComponent(meshInstant.currentlySetEntity, ComponentMesh).(Mesh)
			mesh.removeMeshInstantFromMesh(meshInstant)
			setComponentInternal(meshInstant.currentlySetEntity, ComponentMesh, mesh)
		}

		mesh := GetComponent(meshInstant.MeshEntity, ComponentMesh).(Mesh)
		mesh.instances = append(mesh.instances, meshInstant.ownEntity)
		mesh.needsInstanceUpdate = true
		setComponentInternal(meshInstant.MeshEntity, ComponentMesh, mesh)

		meshInstant.currentlySetEntity = meshInstant.MeshEntity
	}

	return meshInstant
}
func removeMeshInstant(component interface{}, entityId int) interface{} {
	meshInstant := component.(MeshInstant)

	if HasComponent(meshInstant.currentlySetEntity, ComponentMesh) {
		mesh := GetComponent(meshInstant.currentlySetEntity, ComponentMesh).(Mesh)
		mesh.removeMeshInstantFromMesh(meshInstant)
		setComponentInternal(meshInstant.currentlySetEntity, ComponentMesh, mesh)
	}

	return meshInstant
}

func (mesh *Mesh) removeMeshInstantFromMesh(meshInstant MeshInstant) {
	for i := len(mesh.instances); i > 0; i-- {
		if mesh.instances[i] == meshInstant.currentlySetEntity {
			mesh.instances = append(mesh.instances[:i], mesh.instances[i+1:]...)
		}
	}
}
