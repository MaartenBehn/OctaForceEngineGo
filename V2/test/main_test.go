package test

import (
	"log"
	"testing"
)
import of "OctaForceEngineGo/V2"

func TestOctaForce(t *testing.T) {
	of.Init(start, stop, "Test")
}

func start() {
	i := 0
	task := of.NewTask(func() {
		log.Print(i)
		i++
	})
	task.SetRepeating(true)
	of.AddTask(task)
}

func stop() {

}
