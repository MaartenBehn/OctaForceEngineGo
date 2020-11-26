package OctaForceEngine

import (
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strings"
)

type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	UVCord   mgl32.Vec2
}
type Mesh struct {
	Vertices      []Vertex
	Indices       []uint32
	needNewBuffer bool
	vao           uint32
	vertexVBO     uint32
	instanceVBO   uint32
	ebo           uint32

	Instants []int
	Material Material
}

func setUpMesh(_ interface{}) interface{} {
	return Mesh{}
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

	mesh.needNewBuffer = true

	return mesh
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
		mesh.Instants = append(mesh.Instants, meshInstant.OwnEntity)
		mesh.needNewBuffer = true
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
	for i := len(mesh.Instants); i > 0; i-- {
		if mesh.Instants[i] == meshInstant.currentlySetEntity {
			mesh.Instants = append(mesh.Instants[:i], mesh.Instants[i+1:]...)
		}
	}
}
