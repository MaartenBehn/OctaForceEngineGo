package V2

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Position       mgl32.Vec3
	rotation       mgl32.Vec3
	rotationMatrix mgl32.Mat4
	Scale          mgl32.Vec3
	matrix         mgl32.Mat4
}

func NewTransform() *Transform {
	return &Transform{
		Position: mgl32.Vec3{0, 0, 0},
		rotation: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1},
	}
}
