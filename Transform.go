package OctaForceEngine

import "github.com/go-gl/mathgl/mgl32"

// Transform is the Component that hold the position and rotation data of the entity.
type Transform struct {
	Position       mgl32.Vec3
	rotation       mgl32.Vec3
	rotationMatrix mgl32.Mat4
	Scale          mgl32.Vec3
	matrix         mgl32.Mat4
}

func setUpTransform(data interface{}, entityId int) interface{} {
	var transform Transform
	if data == nil {
		transform = Transform{
			Position: mgl32.Vec3{0, 0, 0},
			rotation: mgl32.Vec3{0, 0, 0},
			Scale:    mgl32.Vec3{1, 1, 1}}
	} else {
		transform = data.(Transform)
	}

	transform.calcRotationMatrix()
	transform = setTransformMatrix(transform, entityId).(Transform)
	return transform
}
func setTransformMatrix(component interface{}, entityId int) interface{} {
	transform := component.(Transform)
	transform.matrix = mgl32.Translate3D(
		transform.Position.X(),
		transform.Position.Y(),
		transform.Position.Z())

	transform.matrix = transform.matrix.Mul4(transform.rotationMatrix)

	return transform
}

func (transform *Transform) MoveRelative(vec3 mgl32.Vec3) {
	vec3 = mgl32.TransformCoordinate(vec3, transform.rotationMatrix)
	transform.Position = mgl32.Vec3{
		transform.Position.X() + vec3.X(),
		transform.Position.Y() + vec3.Y(),
		transform.Position.Z() + vec3.Z()}
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

	transform.calcRotationMatrix()
}
func (transform *Transform) Rotate(vec3 mgl32.Vec3) {
	transform.rotation = mgl32.Vec3{
		transform.rotation.X() + mgl32.DegToRad(vec3.X()),
		transform.rotation.Y() + mgl32.DegToRad(vec3.Y()),
		transform.rotation.Z() + mgl32.DegToRad(vec3.Z())}

	transform.calcRotationMatrix()
}
func (transform *Transform) calcRotationMatrix() {
	transform.rotationMatrix = mgl32.Ident4()
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.Y(), mgl32.Vec3{0, 1, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.X(), mgl32.Vec3{1, 0, 0}))
	transform.rotationMatrix = transform.rotationMatrix.Mul4(mgl32.HomogRotate3D(
		transform.rotation.Z(), mgl32.Vec3{0, 0, 1}))
}
