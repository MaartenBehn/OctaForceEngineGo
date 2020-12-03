package OctaForceEngine

import (
	"github.com/go-gl/mathgl/mgl32"
	"io/ioutil"
	"log"
	"strings"
)

// Material is a Struct with is needed by the Mesh Component to set the Color of an Mesh.
type Material struct {
	DiffuseColor mgl32.Vec3
}

// LoadMtl loads the material file of an OBJ File.
// But you can also just set the loadMaterials bool in the LoadOBJ function of the Mesh Component to true.
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
