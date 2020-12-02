package OctaForceEngine

import (
	"io/ioutil"
	"log"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	UVCord   mgl32.Vec2
}

type Mesh struct {
	Vertices        []Vertex
	Indices         []uint32
	vao             uint32
	vbo             uint32
	ebo             uint32
	Material        Material
	needsMeshUpdate bool
}

func setUpMesh(_ interface{}) interface{} {
	mesh := Mesh{}

	printGlErrors("Mesh VAO")
	return mesh
}
func deleteMesh(component interface{}) interface{} {
	mesh := component.(Mesh)
	unUsedVAOs = append(unUsedVAOs, mesh.vao)
	return nil
}

func renderMesh(entityId int) {

	mesh := GetComponent(entityId, ComponentMesh).(Mesh)

	if mesh.needsMeshUpdate {

		var vertexData []float32
		for _, vertex := range mesh.Vertices {
			vertexData = append(vertexData, []float32{
				vertex.Position.X(),
				vertex.Position.Y(),
				vertex.Position.Z(),
			}...)
		}
		gl.GenVertexArrays(1, &mesh.vao)
		gl.BindVertexArray(mesh.vao)

		// Vertex VBO
		gl.GenBuffers(1, &mesh.vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(0))

		// EBO
		gl.GenBuffers(1, &mesh.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

		mesh.needsMeshUpdate = false
		printGlErrors("Mesh SetUp")

		SetComponent(entityId, ComponentMesh, mesh)
	} else {
		gl.BindVertexArray(mesh.vao)
	}

	transform := GetComponent(entityId, ComponentTransform).(Transform)
	gl.UniformMatrix4fv(2, 1, false, &transform.matrix[0])

	// Color
	gl.Uniform3f(3, mesh.Material.DiffuseColor[0], mesh.Material.DiffuseColor[1], mesh.Material.DiffuseColor[2])

	gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, nil)
	printGlErrors("Mesh Render")
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
}

type MeshInstantRoot struct {
	Vertices            []Vertex
	Indices             []uint32
	vao                 uint32
	vertexVBO           uint32
	instanceVBO         uint32
	ebo                 uint32
	instances           []int
	Material            Material
	needsMeshUpdate     bool
	needsInstanceUpdate bool
}

func setUpMeshInstanceRoot(_ interface{}) interface{} {
	meshInstantRoot := MeshInstantRoot{}
	//gl.GenVertexArrays(1, &meshInstantRoot.vao)
	printGlErrors("Mesh Instance VAO")
	return meshInstantRoot
}
func deleteMeshInstanceRoot(component interface{}) interface{} {
	meshInstantRoot := component.(MeshInstantRoot)
	unUsedVAOs = append(unUsedVAOs, meshInstantRoot.vao)
	return nil
}

func renderMeshIstantRoot(entityId int) {
	meshInstantRoot := GetComponent(entityId, ComponentMeshInstantRoot).(MeshInstantRoot)
	transform := GetComponent(entityId, ComponentTransform).(Transform)
	gl.BindVertexArray(meshInstantRoot.vao)

	if meshInstantRoot.needsMeshUpdate {
		var vertexData []float32
		for _, vertex := range meshInstantRoot.Vertices {
			vertexData = append(vertexData, []float32{
				vertex.Position.X(),
				vertex.Position.Y(),
				vertex.Position.Z(),
			}...)
		}

		// Vertex VBO
		gl.GenBuffers(1, &meshInstantRoot.vertexVBO)
		gl.BindBuffer(gl.ARRAY_BUFFER, meshInstantRoot.vertexVBO)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertexData)*4, gl.Ptr(vertexData), gl.STATIC_DRAW)

		gl.EnableVertexAttribArray(0)
		gl.VertexAttribPointer(0, 3, gl.FLOAT, false, vertexStride, gl.PtrOffset(0))

		// EBO
		gl.GenBuffers(1, &meshInstantRoot.ebo)
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, meshInstantRoot.ebo)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(meshInstantRoot.Indices)*4, gl.Ptr(meshInstantRoot.Indices), gl.STATIC_DRAW)

		meshInstantRoot.needsMeshUpdate = false
	}
	if len(meshInstantRoot.instances) > 0 {
		if meshInstantRoot.needsInstanceUpdate {
			// Instance VBO
			gl.GenBuffers(1, &meshInstantRoot.instanceVBO)
			gl.BindBuffer(gl.ARRAY_BUFFER, meshInstantRoot.instanceVBO)
			gl.BufferData(gl.ARRAY_BUFFER, (len(meshInstantRoot.instances)+1)*int(instanceStride), gl.Ptr(nil), gl.DYNAMIC_DRAW)

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

			meshInstantRoot.needsInstanceUpdate = false
		}

		// Set Instance Data
		var instanceData = []float32{
			meshInstantRoot.Material.DiffuseColor[0],
			meshInstantRoot.Material.DiffuseColor[1],
			meshInstantRoot.Material.DiffuseColor[2],

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
		for _, id := range meshInstantRoot.instances {
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

		gl.BindBuffer(gl.ARRAY_BUFFER, meshInstantRoot.instanceVBO)
		gl.BufferData(gl.ARRAY_BUFFER, (len(meshInstantRoot.instances)+1)*int(instanceStride), gl.Ptr(instanceData), gl.DYNAMIC_DRAW)
	}
	SetComponent(entityId, ComponentMeshInstantRoot, meshInstantRoot)

	gl.DrawElementsInstanced(gl.TRIANGLES, int32(len(meshInstantRoot.Indices)), gl.UNSIGNED_INT, nil, int32(len(meshInstantRoot.instances)+1))
}

// LoadOBJ returns the mesh struct of the given OBJ file.
func (mesh *MeshInstantRoot) LoadOBJ(path string, loadMaterials bool) {

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

		if HasComponent(meshInstant.currentlySetEntity, ComponentMeshInstantRoot) {
			meshInstantRoot := GetComponent(meshInstant.currentlySetEntity, ComponentMeshInstantRoot).(MeshInstantRoot)
			meshInstantRoot.removeMeshInstantFromMesh(meshInstant)
			SetComponent(meshInstant.currentlySetEntity, ComponentMeshInstantRoot, meshInstantRoot)
		}

		meshInstantRoot := GetComponent(meshInstant.MeshEntity, ComponentMeshInstantRoot).(MeshInstantRoot)
		meshInstantRoot.instances = append(meshInstantRoot.instances, meshInstant.OwnEntity)
		meshInstantRoot.needsInstanceUpdate = true
		SetComponent(meshInstant.MeshEntity, ComponentMeshInstantRoot, meshInstantRoot)

		meshInstant.currentlySetEntity = meshInstant.MeshEntity
	}

	return meshInstant
}
func removeMeshInstant(component interface{}) interface{} {
	meshInstant := component.(MeshInstant)

	if HasComponent(meshInstant.currentlySetEntity, ComponentMesh) {
		meshInstantRoot := GetComponent(meshInstant.currentlySetEntity, ComponentMesh).(MeshInstantRoot)
		meshInstantRoot.removeMeshInstantFromMesh(meshInstant)
		SetComponent(meshInstant.currentlySetEntity, ComponentMeshInstantRoot, meshInstantRoot)
	}

	return meshInstant
}

func (meshInstantRoot *MeshInstantRoot) removeMeshInstantFromMesh(meshInstant MeshInstant) {
	for i := len(meshInstantRoot.instances); i > 0; i-- {
		if meshInstantRoot.instances[i] == meshInstant.currentlySetEntity {
			meshInstantRoot.instances = append(meshInstantRoot.instances[:i], meshInstantRoot.instances[i+1:]...)
		}
	}
}
