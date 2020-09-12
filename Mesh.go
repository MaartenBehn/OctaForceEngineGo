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
	Position   mgl32.Vec3
	Normal     mgl32.Vec3
	TextVertex mgl32.Vec2
}

type Mesh struct {
	Vertices []Vertex
	meshData []float32
}

var activeMeshes []Mesh
var allMeshData []float32

func (mesh *Mesh) updateMeshData() {
	mesh.meshData = []float32{}
	for _, vertex := range mesh.Vertices {
		mesh.meshData = append(mesh.meshData, []float32{
			vertex.Position.X(),
			vertex.Position.Y(),
			vertex.Position.Z(),
		}...)
	}
}

var vao uint32
var vbo uint32

func updateAllMeshData() {
	allMeshData = []float32{}
	for _, mesh := range activeMeshes {
		allMeshData = append(allMeshData, mesh.meshData...)
	}

	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(allMeshData)*4, gl.Ptr(allMeshData), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.BindVertexArray(0)
}

func renderMeshes() {

	gl.BindVertexArray(vao)
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(allMeshData)))
}

func LoadOBJ(path string) Mesh {

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	var vertices []mgl32.Vec3
	var normals []mgl32.Vec3
	var textVertices []mgl32.Vec2
	var faces [][3][3]float32
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
			textVertices = append(textVertices, mgl32.Vec2{parseFloat(values[1]), parseFloat(values[2])})
			break
		case "f":
			var face [3][3]float32
			for j, value := range values {
				if j == 0 {
					continue
				}

				number := strings.Split(value, "/")
				face[j-1][0] = parseFloat(number[0])
				face[j-1][1] = parseFloat(number[1])
				face[j-1][2] = parseFloat(number[2])
			}
			faces = append(faces, face)
			break
		}
	}

	mesh := Mesh{}
	for i, _ := range vertices {
		vertex := Vertex{
			Position: vertices[i],
		}
		mesh.Vertices = append(mesh.Vertices, vertex)
	}
	return mesh
}

func parseFloat(number string) float32 {
	float, _ := strconv.ParseFloat(number, 32)
	return float32(float)
}
