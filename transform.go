package OctaForce

import "github.com/go-gl/mathgl/mgl32"

type Transform struct {
	position mgl32.Vec3
	rotation mgl32.Vec3
	Scale    mgl32.Vec3

	rotationMatrix          mgl32.Mat4
	needsRotateMatrixUpdate bool

	matrix            mgl32.Mat4
	needsMatrixUpdate bool

	updateTask *task
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
func (t *Transform) GetPosition() mgl32.Vec3 {
	return t.position
}
func (t *Transform) SetPosition(position mgl32.Vec3) {
	t.position = position
	t.needsMatrixUpdate = true
	t.updateMatrix()
}
func (t *Transform) MoveRelative(vec3 mgl32.Vec3) {
	if t.rotation[0] != 0 || t.rotation[1] != 0 || t.rotation[2] != 0 {
		vec3 = mgl32.TransformCoordinate(vec3, t.getRotateMatrix())
	}

	t.position = mgl32.Vec3{
		t.position[0] + vec3[0],
		t.position[1] + vec3[1],
		t.position[2] + vec3[2]}

	t.needsMatrixUpdate = true
	t.updateMatrix()
}
func (t *Transform) GetRotation() mgl32.Vec3 {
	return mgl32.Vec3{
		mgl32.RadToDeg(t.rotation[0]),
		mgl32.RadToDeg(t.rotation[1]),
		mgl32.RadToDeg(t.rotation[2])}
}
func (t *Transform) SetRotaion(vec3 mgl32.Vec3) {
	t.rotation = mgl32.Vec3{
		mgl32.DegToRad(vec3[0]),
		mgl32.DegToRad(vec3[1]),
		mgl32.DegToRad(vec3[2])}
	t.needsRotateMatrixUpdate = true
	t.needsMatrixUpdate = true
	t.updateMatrix()
}
func (t *Transform) Rotate(vec3 mgl32.Vec3) {
	t.rotation = mgl32.Vec3{
		t.rotation.X() + mgl32.DegToRad(vec3.X()),
		t.rotation.Y() + mgl32.DegToRad(vec3.Y()),
		t.rotation.Z() + mgl32.DegToRad(vec3.Z())}
	t.needsRotateMatrixUpdate = true
	t.needsMatrixUpdate = true
	t.updateMatrix()
}

func (t *Transform) updateRotationMatrix() {
	if t.needsRotateMatrixUpdate {
		t.rotationMatrix = mgl32.Ident4()

		t.rotationMatrix = t.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			t.rotation.Y(), mgl32.Vec3{0, 1, 0}))

		t.rotationMatrix = t.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			t.rotation.X(), mgl32.Vec3{1, 0, 0}))

		t.rotationMatrix = t.rotationMatrix.Mul4(mgl32.HomogRotate3D(
			t.rotation.Z(), mgl32.Vec3{0, 0, 1}))

		t.needsRotateMatrixUpdate = false
	}
}
func (t *Transform) getRotateMatrix() mgl32.Mat4 {
	t.updateRotationMatrix()
	return t.rotationMatrix
}

func (t *Transform) updateMatrix() {
	if t.needsMatrixUpdate {
		t.matrix = mgl32.Translate3D(
			t.position.X(),
			t.position.Y(),
			t.position.Z())

		t.matrix = t.matrix.Mul4(t.getRotateMatrix())

		t.needsMatrixUpdate = false
	}
}
func (t *Transform) getMatrix() mgl32.Mat4 {
	return t.matrix
}
