package test

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/profile"
	"math"
	"path/filepath"
	"runtime"
	"testing"
	"time"
)
import of "OctaForceEngineGo"

var absPath string

func init() {
	_, b, _, _ := runtime.Caller(0)
	absPath = filepath.Dir(b)
}

func TestOctaForce(t *testing.T) {
	defer profile.Start().Stop()
	of.Init(start, stop, "Test")
}

func start() {
	camera := of.NewCamera()
	camera.Transform = of.NewTransform()
	camera.Transform.SetPosition(mgl32.Vec3{1000, 1000, 2000})

	of.ActiveCameraData.Camera = camera

	mesh := of.NewMesh()
	mesh.LoadOBJ(absPath+"/mesh/LowPolySphere.obj", false)
	mesh.Material = of.Material{DiffuseColor: [3]float32{1, 1, 1}}
	of.ActiveMeshesData.AddMesh(mesh)

	for i := 0; i < 200; i++ {
		var instants []*of.MeshInstant
		for j := 0; j < 1000; j++ {
			meshInstant := of.NewMeshInstant(mesh, &of.Material{DiffuseColor: [3]float32{1, 0, 1}})
			meshInstant.Transform.SetPosition(mgl32.Vec3{float32(i) * 10, float32(j) * 10, 0})
			instants = append(instants, meshInstant)
		}
		task := of.NewTask(func() {
			for _, instant := range instants {
				instant.Transform.MoveRelative(mgl32.Vec3{0, 0, float32(math.Sin(float64(time.Now().Second())))})
			}
		})
		task.SetRepeating(true)
		of.AddTask(task)
	}

}

func stop() {

}
