package OctaForce

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
	"unsafe"
)

type guiData struct {
	context *imgui.Context
	io      imgui.IO

	fontTexture            uint32
	shaderHandle           uint32
	vertHandle             uint32
	fragHandle             uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32
}

var gui guiData

func initGui() {
	gui := guiData{}

	gui.context = imgui.CreateContext(nil)
	gui.io = imgui.CurrentIO()

	// Init Gui Render
	var err error
	gui.shaderHandle = gl.CreateProgram()
	gui.vertHandle, err = compileShader(absPath+"/shader/guiVertexShader.shader", gl.VERTEX_SHADER)
	if err != nil {
		panic(err)
	}
	gui.fragHandle, err = compileShader(absPath+"/shader/guiFragmentShader.shader", gl.FRAGMENT_SHADER)
	if err != nil {
		panic(err)
	}

	gl.AttachShader(gui.shaderHandle, gui.vertHandle)
	gl.AttachShader(gui.shaderHandle, gui.fragHandle)
	gl.LinkProgram(gui.shaderHandle)

	gui.attribLocationTex = gl.GetUniformLocation(gui.shaderHandle, gl.Str("Texture"+"\x00"))
	gui.attribLocationProjMtx = gl.GetUniformLocation(gui.shaderHandle, gl.Str("ProjMtx"+"\x00"))
	gui.attribLocationPosition = gl.GetAttribLocation(gui.shaderHandle, gl.Str("Position"+"\x00"))
	gui.attribLocationUV = gl.GetAttribLocation(gui.shaderHandle, gl.Str("UV"+"\x00"))
	gui.attribLocationColor = gl.GetAttribLocation(gui.shaderHandle, gl.Str("Color"+"\x00"))

	gl.GenBuffers(1, &gui.vboHandle)
	gl.GenBuffers(1, &gui.elementsHandle)

	// Init Gui Font
	image := gui.io.Fonts().TextureDataAlpha8()
	gl.GenTextures(1, &gui.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, gui.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height),
		0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	// Store our identifier
	gui.io.Fonts().SetTextureID(imgui.TextureID(gui.fontTexture))

	//mapGuiKeys()
}

func mapGuiKeys() {
	gui.io.KeyMap(imgui.KeyTab, int(KeyTab))
	gui.io.KeyMap(imgui.KeyLeftArrow, int(KeyLeft))
	gui.io.KeyMap(imgui.KeyRightArrow, int(KeyRight))
	gui.io.KeyMap(imgui.KeyUpArrow, int(KeyUp))
	gui.io.KeyMap(imgui.KeyDownArrow, int(KeyDown))
	gui.io.KeyMap(imgui.KeyPageUp, int(KeyPageUp))
	gui.io.KeyMap(imgui.KeyPageDown, int(KeyPageDown))
	// TODO all other keys
}

func renderGui() {

	drawData := imgui.RenderedDrawData()
	if !drawData.Valid() {
		return
	}

	displayWidth, displayHeight := window.GetSize()
	fbWidth, fbHeight := window.GetFramebufferSize()

	drawData.ScaleClipRects(imgui.Vec2{
		X: float32(fbWidth) / float32(displayWidth),
		Y: float32(fbHeight) / float32(displayHeight),
	})

	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [4][4]float32{
		{2.0 / float32(displayWidth), 0.0, 0.0, 0.0},
		{0.0, 2.0 / float32(-displayHeight), 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	gl.UseProgram(gui.shaderHandle)
	gl.Uniform1i(gui.attribLocationTex, 0)
	gl.UniformMatrix4fv(gui.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	var vaoHandle uint32
	gl.GenVertexArrays(1, &vaoHandle)
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(gl.ARRAY_BUFFER, gui.vboHandle)
	gl.EnableVertexAttribArray(uint32(gui.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(gui.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(gui.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribPointer(uint32(gui.attribLocationPosition), 2, gl.FLOAT, false, int32(vertexSize), gl.PtrOffset(vertexOffsetPos))
	gl.VertexAttribPointer(uint32(gui.attribLocationUV), 2, gl.FLOAT, false, int32(vertexSize), gl.PtrOffset(vertexOffsetUv))
	gl.VertexAttribPointer(uint32(gui.attribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), gl.PtrOffset(vertexOffsetCol))
	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	const bytesPerUint32 = 4
	if indexSize == bytesPerUint32 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		var indexBufferOffset uintptr

		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, gui.vboHandle)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, gui.elementsHandle)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElements(gl.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType), unsafe.Pointer(indexBufferOffset))
			}
			indexBufferOffset += uintptr(cmd.ElementCount() * indexSize)
		}
	}

}
