package OctaForceEngine

import "github.com/go-gl/mathgl/mgl32"

// An Entity with an Camera Component can be set as the used as a Camera with SetActiveCameraEntity.
type Camera struct {
	projection mgl32.Mat4
}

func setUpCamera(_ interface{}, entityId int) interface{} {
	return Camera{
		projection: mgl32.Perspective(mgl32.DegToRad(45.0),
			float32(windowWidth)/windowHeight,
			0.1,
			100000.0),
	}
}

// SetActiveCameraEntity sets the given entity as the camera.
func SetActiveCameraEntity(entityId int) {
	if HasComponent(entityId, ComponentCamera) {
		cameraEntityId = entityId
	}
}
