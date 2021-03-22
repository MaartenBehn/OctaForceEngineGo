package V2

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Vec3
	Scale    mgl32.Vec3

	rotationMatrix          mgl32.Mat4
	needsRotateMatrixUpdate bool

	matrix            mgl32.Mat4
	needsMatrixUpdate bool

	updateTask *Task
}

func NewTransform() *Transform {
	transform := &Transform{
		position:                mgl32.Vec3{0, 0, 0},
		rotation:                mgl32.Vec3{0, 0, 0},
		Scale:                   mgl32.Vec3{1, 1, 1},
		needsRotateMatrixUpdate: true,
		needsMatrixUpdate:       true,
	}
	transform.updateMatrix()
	return transform
}
func (transform *Transform) GetPosition() mgl32.Vec3 {
	return transform.position
}
func (transform *Transform) SetPosition(position mgl32.Vec3) {
	transform.position = position
	transform.needsMatrixUpdate = true
	transform.updateMatrix()
}
func (transform *Transform) MoveRelative(vec3 mgl32.Vec3) {
	if transform.rotation[0] != 0 && transform.rotation[1] != 0 && transform.rotation[2] != 0 {
		vec3 = mgl32.TransformCoordinate(vec3, transform.getRotateMatrix())
	}

	transform.position = mgl32.Vec3{
		transform.position[0] + vec3[0],
		transform.position[1] + vec3[1],
		transform.position[2] + vec3[2]}

	transform.needsMatrixUpdate = true
	transform.updateMatrix()
}
func (transform *Transform) GetRotation() mgl32.Vec3 {
	return mgl32.Vec3{
		mgl32.RadToDeg(transform.rotation[0]),
		mgl32.RadToDeg(transform.rotation[1]),
		mgl32.RadToDeg(transform.rotation[2])}
}
func (transform *Transform) SetRotaion(vec3 mgl32.Vec3) {
	transform.rotation = mgl32.Vec3{
		mgl32.DegToRad(vec3[0]),
		mgl32.DegToRad(vec3[1]),
		mgl32.DegToRad(vec3[2])}
	transform.needsRotateMatrixUpdate = true
	transform.updateMatrix()
}
func (transform *Transform) Rotate(vec3 mgl32.Vec3) {
	transform.rotation = mgl32.Vec3{
		transform.rotation.X() + mgl32.DegToRad(vec3.X()),
		transform.rotation.Y() + mgl32.DegToRad(vec3.Y()),
		transform.rotation.Z() + mgl32.DegToRad(vec3.Z())}
	transform.needsRotateMatrixUpdate = true
	transform.updateMatrix()
}

func (transform *Transform) updateRotationMatrix() {
	if transform.needsRotateMatrixUpdate {
		transform.rotationMatrix = mgl32.Ident4()

		transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			transform.rotation.Y(), mgl32.Vec3{0, 1, 0}))

		transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			transform.rotation.X(), mgl32.Vec3{1, 0, 0}))

		transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			transform.rotation.Z(), mgl32.Vec3{0, 0, 1}))

		transform.needsRotateMatrixUpdate = false
	}
}
func (transform *Transform) getRotateMatrix() mgl32.Mat4 {
	transform.updateRotationMatrix()
	return transform.rotationMatrix
}

func (transform *Transform) updateMatrix() {
	if transform.needsMatrixUpdate {
		transform.matrix = mgl32.Translate3D(
			transform.position.X(),
			transform.position.Y(),
			transform.position.Z())

		transform.matrix = transform.matrix.Mul4(transform.getRotateMatrix())

		transform.needsMatrixUpdate = false
	}
}
func (transform *Transform) getMatrix() mgl32.Mat4 {
	transform.updateMatrix()
	return transform.matrix
}
