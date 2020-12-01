package OctaForceEngine

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	UVCord   mgl32.Vec2
}
type Mesh struct {
	Vertices    []Vertex
	Indices     []uint32
	vao         uint32
	vertexVBO   uint32
	instanceVBO uint32
	ebo         uint32

	instants []int
	Material Material

	needsMeshUpdate     bool
	needsInstanceUpdate bool
}

func setUpMesh(_ interface{}) interface{} {
	mesh := Mesh{}
	gl.GenVertexArrays(1, &mesh.vao)
	return mesh
}
func deleteMesh(component interface{}) interface{} {
	mesh := component.(Mesh)
	unUsedVAOs = append(unUsedVAOs, mesh.vao)
	return nil
}

// LoadOBJ returns the mesh struct of the given OBJ file.
func LoadOBJ(path string, loadMaterials bool) Mesh {

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

	mesh := Mesh{}
	mesh.Vertices = make([]Vertex, len(vertices))
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

	mesh.needsMeshUpdate = true

	return mesh
}

func setUpMeshVAO(entityId int) {
	mesh := GetComponent(entityId, ComponentMesh).(Mesh)
	transform := GetComponent(entityId, ComponentTransform).(Transform)

	gl.BindVertexArray(mesh.vao)
	if mesh.needsMeshUpdate {
		var vertexData []float32
		for _, vertex := range mesh.Vertices {
			vertexData = append(vertexData, []float32{
				vertex.Position.X(),
				vertex.Position.Y(),
				vertex.Position.Z(),
			}...)
		}

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

		mesh.needsMeshUpdate = false
	}

	if len(mesh.instants) > 0 {
		if mesh.needsInstanceUpdate {
			// Instance VBO
			gl.GenBuffers(1, &mesh.instanceVBO)
			gl.BindBuffer(gl.ARRAY_BUFFER, mesh.instanceVBO)
			gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.instants)+1)*int(instanceStride), gl.Ptr(nil), gl.DYNAMIC_DRAW)

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

			mesh.needsInstanceUpdate = false
		}

		// Set Instance Data
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
		for _, id := range mesh.instants {
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
		gl.BufferData(gl.ARRAY_BUFFER, (len(mesh.instants)+1)*int(instanceStride), gl.Ptr(instanceData), gl.DYNAMIC_DRAW)
	}
	SetComponent(entityId, ComponentMesh, mesh)
}

type MeshInstant struct {
	OwnEntity          int
	MeshEntity         int
	currentlySetEntity int

	Material Material
}

func setUpMeshInstant(_ interface{}) interface{} {
	return MeshInstant{}
}
func addMeshInstant(component interface{}) interface{} {
	meshInstant := component.(MeshInstant)

	if meshInstant.OwnEntity == 0 || meshInstant.MeshEntity == 0 {
		return component
	}

	if meshInstant.currentlySetEntity != meshInstant.MeshEntity {

		if HasComponent(meshInstant.currentlySetEntity, ComponentMesh) {
			mesh := GetComponent(meshInstant.currentlySetEntity, ComponentMesh).(Mesh)
			mesh.removeMeshInstantFromMesh(meshInstant)
			SetComponent(meshInstant.currentlySetEntity, ComponentMesh, mesh)
		}

		mesh := GetComponent(meshInstant.MeshEntity, ComponentMesh).(Mesh)
		mesh.instants = append(mesh.instants, meshInstant.OwnEntity)
		mesh.needsInstanceUpdate = true
		SetComponent(meshInstant.MeshEntity, ComponentMesh, mesh)

		meshInstant.currentlySetEntity = meshInstant.MeshEntity
	}

	return meshInstant
}

func removeMeshInstant(component interface{}) interface{} {
	meshInstant := component.(MeshInstant)

	if HasComponent(meshInstant.currentlySetEntity, ComponentMesh) {
		mesh := GetComponent(meshInstant.currentlySetEntity, ComponentMesh).(Mesh)
		mesh.removeMeshInstantFromMesh(meshInstant)
		SetComponent(meshInstant.currentlySetEntity, ComponentMesh, mesh)
	}

	return meshInstant
}

func (mesh *Mesh) removeMeshInstantFromMesh(meshInstant MeshInstant) {
	for i := len(mesh.instants); i > 0; i-- {
		if mesh.instants[i] == meshInstant.currentlySetEntity {
			mesh.instants = append(mesh.instants[:i], mesh.instants[i+1:]...)
		}
	}
}
