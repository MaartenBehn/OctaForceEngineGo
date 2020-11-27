package OctaForceEngine

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	Position       mgl32.Vec3
	Rotation       mgl32.Vec3
	rotationMatrix mgl32.Mat4
	Scale          mgl32.Vec3
	matrix         mgl32.Mat4
}

func setUpTransform(_ interface{}) interface{} {
	transform := Transform{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1}}

	transform = setTransformMatrix(transform).(Transform)
	return transform
}
func setTransformMatrix(component interface{}) interface{} {
	transform := component.(Transform)
	transform.matrix = mgl32.Translate3D(
		transform.Position.X(),
		transform.Position.Y(),
		transform.Position.Z())
	transform.calcRotationMatrix()
	transform.matrix = transform.matrix.Mul4(transform.rotationMatrix)
	return transform
}

func (transform *Transform) MoveRelative(vec3 mgl32.Vec3) {
	transform.calcRotationMatrix()
	vec3 = mgl32.TransformCoordinate(vec3, transform.rotationMatrix)
	transform.Position = mgl32.Vec3{
		transform.Position.X() + vec3.X(),
		transform.Position.Y() + vec3.Y(),
		transform.Position.Z() + vec3.Z()}
}
func (transform *Transform) GetRotationInDegree() mgl32.Vec3 {
	return mgl32.Vec3{
		mgl32.RadToDeg(transform.Rotation.X()),
		mgl32.RadToDeg(transform.Rotation.Y()),
		mgl32.RadToDeg(transform.Rotation.Z())}
}
func (transform *Transform) SetRotaionInDegree(vec3 mgl32.Vec3) {
	transform.Rotation = mgl32.Vec3{
		mgl32.DegToRad(vec3.X()),
		mgl32.DegToRad(vec3.Y()),
		mgl32.DegToRad(vec3.Z())}
}
func (transform *Transform) RotateInDegree(vec3 mgl32.Vec3) {
	transform.Rotation = mgl32.Vec3{
		transform.Rotation.X() + mgl32.DegToRad(vec3.X()),
		transform.Rotation.Y() + mgl32.DegToRad(vec3.Y()),
		transform.Rotation.Z() + mgl32.DegToRad(vec3.Z())}
	transform.calcRotationMatrix()
}
func (transform *Transform) calcRotationMatrix() {
	transform.rotationMatrix = mgl32.Ident4()
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.Rotation.Y(), mgl32.Vec3{0, 1, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.Rotation.X(), mgl32.Vec3{1, 0, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.Rotation.Z(), mgl32.Vec3{0, 0, 1}))
}
