package test

import (
	"github.com/go-gl/mathgl/mgl32"
	"log"
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
	//defer profile.Start().Stop()
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

	start := time.Now()
	task := of.NewTask(func() {
		sum := 0
		for i := 0; i < int(time.Since(start).Milliseconds()); i++ {
			sum += i
		}
	})
	task.SetRepeating(true)
	for i := 0; i < 400; i++ {
		//of.AddTask(task)
	}

	for i := 0; i < 100; i++ {
		var instants []*of.MeshInstant
		for j := 0; j < 10000; j++ {
			meshInstant := of.NewMeshInstant(mesh, &of.Material{DiffuseColor: [3]float32{1, 0, 1}})
			meshInstant.Transform.SetPosition(mgl32.Vec3{float32(i) * 10, float32(j) * 10, 0})
			instants = append(instants, meshInstant)
		}

		newTask(instants)
	}
}

func newTask(instants []*of.MeshInstant) {
	task := of.NewTask(func() {
		z := float32(math.Sin(float64(time.Now().Second())))
		for _, instant := range instants {
			instant.Transform.MoveRelative(mgl32.Vec3{0, 0, z})
		}
	})
	task.SetRepeating(true)
	of.AddTask(task)
}

func stop() {

}

var globalInt int

func TestGoRotine(t *testing.T) {

	for i := 0; i < 100; i++ {
		go multiFunc(i)
	}
	multiFunc(101)

}
func multiFunc(id int) {
	for {
		globalInt++
		log.Printf("%d %d", id, globalInt)
	}
}
