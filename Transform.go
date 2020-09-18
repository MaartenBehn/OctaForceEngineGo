package OctaForceEngine

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	position       mgl32.Vec3
	rotation       mgl32.Vec3
	rotationMatrix mgl32.Mat4
	scale          mgl32.Vec3
	matrix         mgl32.Mat4
}

func setUpTransform(_ interface{}) interface{} {
	transform := Transform{
		position: mgl32.Vec3{0, 0, 0},
		rotation: mgl32.Vec3{0, 0, 0},
		scale:    mgl32.Vec3{1, 1, 1}}

	transform.CalcRotationMatrix()
	transform = setTransformMatrix(transform).(Transform)
	return transform
}
func setTransformMatrix(data interface{}) interface{} {
	transform := data.(Transform)
	transform.matrix = mgl32.Translate3D(
		transform.position.X(),
		transform.position.Y(),
		transform.position.Z())
	transform.matrix = transform.matrix.Mul4(transform.rotationMatrix)
	return transform
}
func (transform *Transform) GetPosition() mgl32.Vec3 {
	return transform.position
}
func (transform *Transform) SetPosition(vec3 mgl32.Vec3) {
	transform.position = vec3
}
func (transform *Transform) Move(vec3 mgl32.Vec3) {
	transform.position = mgl32.Vec3{
		transform.position.X() + vec3.X(),
		transform.position.Y() + vec3.Y(),
		transform.position.Z() + vec3.Z()}
}
func (transform *Transform) MoveRelative(vec3 mgl32.Vec3) {
	vec3 = mgl32.TransformCoordinate(vec3, transform.rotationMatrix)
	transform.position = mgl32.Vec3{
		transform.position.X() + vec3.X(),
		transform.position.Y() + vec3.Y(),
		transform.position.Z() + vec3.Z()}
}

func (transform *Transform) GetRotation() mgl32.Vec3 {
	return mgl32.Vec3{
		mgl32.RadToDeg(transform.rotation.X()),
		mgl32.RadToDeg(transform.rotation.Y()),
		mgl32.RadToDeg(transform.rotation.Z())}
}
func (transform *Transform) SetRotaion(vec3 mgl32.Vec3) {
	transform.rotation = mgl32.Vec3{
		mgl32.DegToRad(vec3.X()),
		mgl32.DegToRad(vec3.Y()),
		mgl32.DegToRad(vec3.Z())}
	transform.CalcRotationMatrix()
}
func (transform *Transform) Rotate(vec3 mgl32.Vec3) {
	transform.rotation = mgl32.Vec3{
		transform.rotation.X() + mgl32.DegToRad(vec3.X()),
		transform.rotation.Y() + mgl32.DegToRad(vec3.Y()),
		transform.rotation.Z() + mgl32.DegToRad(vec3.Z())}
	transform.CalcRotationMatrix()
}
func (transform *Transform) CalcRotationMatrix() {
	transform.rotationMatrix = mgl32.Ident4()
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.Y(), mgl32.Vec3{0, 1, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.X(), mgl32.Vec3{1, 0, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.Z(), mgl32.Vec3{0, 0, 1}))
}
