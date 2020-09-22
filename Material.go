package OctaForceEngine

import (
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strings"
)

type Material struct {
	DiffuseColor mgl32.Vec3
}

func LoadMtl(path string) []Material {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	lines := strings.Split(string(content), "\n")

	materials := []Material{}
	var currentMaterial int
	for _, line := range lines {
		values := strings.Split(line, " ")
		values[len(values)-1] = strings.Replace(values[len(values)-1], "\r", "", 1)

		switch values[0] {
		case "newmtl":
			currentMaterial = len(materials)
			materials = append(materials, Material{})
			break
		case "Kd":
			materials[currentMaterial].DiffuseColor = mgl32.Vec3{
				ParseFloat(values[1]),
				ParseFloat(values[2]),
				ParseFloat(values[3])}
			break
		}
	}

	return materials
}
