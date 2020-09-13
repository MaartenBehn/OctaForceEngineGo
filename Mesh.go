package OctaForce

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
)

type Vertex struct {
	Position mgl32.Vec3
	Normal   mgl32.Vec3
	UVCord   mgl32.Vec2
}
type Mesh struct {
	Vertices   []Vertex
	Indices    []uint32
	vertexData []float32
}

func setUpMesh(data interface{}) interface{} {
	return Mesh{}
}
func updateMeshData(data interface{}) interface{} {
	mesh := data.(Mesh)
	mesh.vertexData = []float32{}
	for _, vertex := range mesh.Vertices {
		mesh.vertexData = append(mesh.vertexData, []float32{
			vertex.Position.X(),
			vertex.Position.Y(),
			vertex.Position.Z(),

			vertex.Normal.X(),
			vertex.Normal.Y(),
			vertex.Normal.Z(),

			vertex.UVCord.X(),
			vertex.UVCord.Y(),
		}...)
	}

	needAllMeshUpdate = true
	return mesh
}

var allVertexData []float32
var activeMeshes []int
var vao uint32
var vbo uint32
var ebo uint32

const stride int32 = 8 * 4

var needAllMeshUpdate bool

func updateAllMeshData() {
	activeMeshes = GetAllEntitiesWithComponent(COMPONENT_Mesh)
	allVertexData = []float32{}
	var allIndexData []uint32
	for _, meshId := range activeMeshes {
		mesh := GetComponent(meshId, COMPONENT_Mesh).(Mesh)
		allVertexData = append(allVertexData, mesh.vertexData...)
		allIndexData = append(allIndexData, mesh.Indices...)
	}

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(allVertexData)*4, gl.Ptr(allVertexData), gl.STATIC_DRAW)

	gl.GenBuffers(1, &ebo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(allIndexData)*4, gl.Ptr(allIndexData), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, stride, gl.PtrOffset(0))

	gl.BindVertexArray(0)

	needAllMeshUpdate = false
}

func renderMeshes() {
	if needAllMeshUpdate {
		updateAllMeshData()
	}

	gl.BindVertexArray(vao)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)

	indexCounter := 0
	for _, entityId := range activeMeshes {
		glTransform = GetComponent(entityId, COMPONENT_Transform).(Transform).Matrix
		gl.UniformMatrix4fv(transformUniform, 1, false, &glTransform[0])

		mesh := GetComponent(entityId, COMPONENT_Mesh).(Mesh)
		gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(indexCounter*4))
		indexCounter += len(mesh.Indices)
	}

}

func LoadOBJ(path string) Mesh {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	var vertices []mgl32.Vec3
	var normals []mgl32.Vec3
	var uvCord []mgl32.Vec2
	var faces [][3][3]uint32
	for _, line := range lines {

		values := strings.Split(line, " ")
		switch values[0] {
		case "v":
			vertices = append(vertices, mgl32.Vec3{parseFloat(values[1]), parseFloat(values[2]), parseFloat(values[3])})
			break
		case "vn":
			normals = append(normals, mgl32.Vec3{parseFloat(values[1]), parseFloat(values[2]), parseFloat(values[3])})
			break
		case "vt":
			uvCord = append(uvCord, mgl32.Vec2{parseFloat(values[1]), parseFloat(values[2])})
			break
		case "f":
			var face [3][3]uint32
			for j, value := range values {
				if j == 0 {
					continue
				}

				number := strings.Split(value, "/")
				face[j-1][0] = parseInt(number[0])
				face[j-1][1] = parseInt(number[1])
				face[j-1][2] = parseInt(number[2])
			}
			faces = append(faces, face)
			break
		}
	}

	mesh := Mesh{}
	mesh.Vertices = make([]Vertex, len(vertices))
	for _, face := range faces {
		for _, values := range face {
			vertexIndex := values[0] - 1
			mesh.Indices = append(mesh.Indices, vertexIndex)
			mesh.Vertices[vertexIndex].Position = vertices[vertexIndex]
			//mesh.Vertices[vertexIndex].UVCord = uvCord[values[1] -1]
			//mesh.Vertices[vertexIndex].Normal = normals[values[2] -1]
		}
	}
	return mesh
}

func parseFloat(number string) float32 {
	float, _ := strconv.ParseFloat(number, 32)
	return float32(float)
}
func parseInt(number string) uint32 {
	int, _ := strconv.ParseInt(number, 10, 32)
	return uint32(int)
}
