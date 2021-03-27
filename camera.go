package OctaForce

import (
	"github.com/go-gl/mathgl/mgl32"
)

type Camera struct {
	projection mgl32.Mat4
	Transform  *Transform
}

func NewCamera() *Camera {
	return &Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(45.0),
			float32(WindowWidth)/float32(WindowHeight),
			0.1,
			100000.0),
	}
}

func initActiveCamera() {
	activeCamera = NewCamera()
}

var activeCamera *Camera

func GetActiveCamera() *Camera {
	return activeCamera
}
func SetActiveCamera(camera *Camera) {
	activeCamera = camera
}
