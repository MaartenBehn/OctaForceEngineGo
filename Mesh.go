package OctaForce

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position mgl32.Vec3
}

type Mesh struct {
	Vertices []Vertex
	meshData []float32
}

var activeMeshes []Mesh

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
	var allMeshData []float32
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
	gl.DrawArrays(gl.TRIANGLES, 0, int32(len(activeMeshes)))
}
