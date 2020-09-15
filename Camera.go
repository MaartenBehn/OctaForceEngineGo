package OctaForceEngine

import "github.com/go-gl/mathgl/mgl32"

type Camera struct {
	projection mgl32.Mat4
}

func setUpCamera(_ interface{}) interface{} {
	return Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(45.0),
			float32(windowWidth)/windowHeight,
			0.1,
			100000.0),
	}
}
