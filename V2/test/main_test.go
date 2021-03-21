package test

import (
	"github.com/pkg/profile"
	"testing"
	"time"
)
import of "OctaForceEngineGo/V2"

func TestOctaForce(t *testing.T) {
	defer profile.Start().Stop()
	of.Init(start, stop, "Test")
}

func start() {
	i := 0
	task := of.NewTask(func() {
		i++
	})
	task.SetRepeating(true)
	of.AddTask(task)

	start := time.Now()
	task = of.NewTask(func() {
		sum := 0
		for i := 0; i < int(time.Since(start).Milliseconds()); i++ {
			sum += i
		}
	})
	task.SetRepeating(true)
	for i := 0; i < 400; i++ {
		of.AddTask(task)
	}
}

func stop() {

}
