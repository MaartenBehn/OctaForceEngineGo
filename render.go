package OctaForce

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/inkyblackness/imgui-go"
	"path/filepath"
	"runtime"
	"time"
)

var absPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	absPath = filepath.Dir(b)
}

func initRender() {
	initImGui()

	initGLFW()
	initImGuiKeyMapping()

	initOpenGL()
	initImGuiGLBuffers()
}

var (
	maxFPS           float64
	fps              float64
	renderFrameStart time.Time
	renderDeltaTime  float64
)

func runRender() {

	clearColor := [3]float32{0.0, 0.0, 0.0}

	wait := time.Duration(1.0 / maxFPS * 1000000000)
	for running {
		renderFrameStart = time.Now()

		processEvents()

		if window.ShouldClose() {
			running = false
		}

		newFrame()
		imgui.NewFrame()

		imgui.Text(fmt.Sprintf("FPS : %.0f", fps))
		imgui.Text(fmt.Sprintf("UPS : %.0f", ups))

		imgui.Render()

		preRender(clearColor)

		for _, programmData := range programmDatas {
			gl.UseProgram(programmData.id)

			// Creating inverted Camera pos
			view := activeCamera.Transform.getMatrix().Inv()
			gl.UniformMatrix4fv(1, 1, false, &view[0])
			gl.UniformMatrix4fv(0, 1, false, &activeCamera.projection[0])

			programmData.renderFunc()
		}

		render(DisplaySize(), FramebufferSize(), imgui.RenderedDrawData())
		postRender()

		diff := time.Since(renderFrameStart)
		if diff > 0 {
			fps = (wait.Seconds() / diff.Seconds()) * maxFPS
		} else {
			fps = 10000
		}

		if diff < wait {
			renderDeltaTime = wait.Seconds()
			time.Sleep(wait - diff)
		} else {
			renderDeltaTime = diff.Seconds()
		}
	}
}
