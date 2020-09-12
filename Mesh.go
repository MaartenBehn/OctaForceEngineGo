package OctaForce

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position mgl32.Vec3
	Normals  mgl32.Vec3
	UV       mgl32.Vec2
}

type Mesh struct {
	Vertices      []Vertex
	Indices       []uint32
	Textures      []Texture
	vao, vbo, ebo uint32
}

func (mesh *Mesh) setUpMesh() {

	gl.GenVertexArrays(1, &mesh.vao)
	gl.GenBuffers(1, &mesh.vbo)
	gl.GenBuffers(1, &mesh.ebo)

	gl.BindVertexArray(mesh.vao)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.ebo)

	gl.BufferData(gl.ARRAY_BUFFER, len(mesh.Vertices)*8*4, gl.Ptr(mesh.Vertices), gl.STATIC_DRAW)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(mesh.Indices)*4, gl.Ptr(mesh.Indices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("position\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))

	normalAttrib := uint32(gl.GetAttribLocation(program, gl.Str("normal\x00")))
	gl.EnableVertexAttribArray(normalAttrib)
	gl.VertexAttribPointer(normalAttrib, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("texCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
}

func (mesh *Mesh) renderMesh() {

	for i, texture := range mesh.Textures {
		gl.ActiveTexture(uint32(gl.TEXTURE0 + i))
		gl.BindTexture(gl.TEXTURE_2D, texture.Id)
	}
	gl.BindVertexArray(mesh.vao)
	gl.DrawElements(gl.TRIANGLES, int32(len(mesh.Indices)), gl.UNSIGNED_INT, gl.PtrOffset(0))
}
