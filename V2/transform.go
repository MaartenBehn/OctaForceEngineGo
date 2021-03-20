package V2

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	data
	Position       mgl32.Vec3
	rotation       mgl32.Vec3
	rotationMatrix mgl32.Mat4
	Scale          mgl32.Vec3
	matrix         mgl32.Mat4
}
