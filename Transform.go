package OctaForce

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
	Matrix   mgl32.Mat4
}

func setUpTransform(data interface{}) interface{} {
	transform := Transform{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1}}

	transform = setTransformMatrix(transform).(Transform)
	return transform
}
func setTransformMatrix(data interface{}) interface{} {
	transform := data.(Transform)
	transform.Matrix = mgl32.Ident4()
	transform.Matrix = transform.Matrix.Add(mgl32.Translate3D(
		transform.Position.X(),
		transform.Position.Y(),
		transform.Position.Z()))
	//transform.Matrix = glTransform.Mul4(mgl32.HomogRotate3DX(transform.Rotation.X()))
	//transform.Matrix = glTransform.Mul4(mgl32.HomogRotate3DX(transform.Rotation.Y()))
	//transform.Matrix = glTransform.Mul4(mgl32.HomogRotate3DX(transform.Rotation.Z()))
	return transform
}
