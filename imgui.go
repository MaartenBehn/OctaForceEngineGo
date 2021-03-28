package OctaForce

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/inkyblackness/imgui-go"
)

type imguiData struct {
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

var gui imguiData

func initImGui() {
	gui = imguiData{}

	gui.context = imgui.CreateContext(nil)
	gui.io = imgui.CurrentIO()
}

func initImGuiKeyMapping() {
	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	gui.io.KeyMap(imgui.KeyTab, int(glfw.KeyTab))
	gui.io.KeyMap(imgui.KeyLeftArrow, int(glfw.KeyLeft))
	gui.io.KeyMap(imgui.KeyRightArrow, int(glfw.KeyRight))
	gui.io.KeyMap(imgui.KeyUpArrow, int(glfw.KeyUp))
	gui.io.KeyMap(imgui.KeyDownArrow, int(glfw.KeyDown))
	gui.io.KeyMap(imgui.KeyPageUp, int(glfw.KeyPageUp))
	gui.io.KeyMap(imgui.KeyPageDown, int(glfw.KeyPageDown))
	gui.io.KeyMap(imgui.KeyHome, int(glfw.KeyHome))
	gui.io.KeyMap(imgui.KeyEnd, int(glfw.KeyEnd))
	gui.io.KeyMap(imgui.KeyInsert, int(glfw.KeyInsert))
	gui.io.KeyMap(imgui.KeyDelete, int(glfw.KeyDelete))
	gui.io.KeyMap(imgui.KeyBackspace, int(glfw.KeyBackspace))
	gui.io.KeyMap(imgui.KeySpace, int(glfw.KeySpace))
	gui.io.KeyMap(imgui.KeyEnter, int(glfw.KeyEnter))
	gui.io.KeyMap(imgui.KeyEscape, int(glfw.KeyEscape))
	gui.io.KeyMap(imgui.KeyA, int(glfw.KeyA))
	gui.io.KeyMap(imgui.KeyC, int(glfw.KeyC))
	gui.io.KeyMap(imgui.KeyV, int(glfw.KeyV))
	gui.io.KeyMap(imgui.KeyX, int(glfw.KeyX))
	gui.io.KeyMap(imgui.KeyY, int(glfw.KeyY))
	gui.io.KeyMap(imgui.KeyZ, int(glfw.KeyZ))
}

func initImGuiGLBuffers() {

	// Backup GL state
	var lastTexture int32
	var lastArrayBuffer int32
	var lastVertexArray int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)

	gui.shaderHandle = gl.CreateProgram()
	var err error
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

	createFontsTexture()

	// Restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindVertexArray(uint32(lastVertexArray))

}

func createFontsTexture() {
	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()

	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &gui.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, gui.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height),
		0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	// Store our identifier
	io.Fonts().SetTextureID(imgui.TextureID(gui.fontTexture))

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
}
