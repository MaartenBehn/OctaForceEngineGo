package OctaForce

import "github.com/go-gl/mathgl/mgl32"

const (
	COMPONENT_Transform = 1
	COMPONENT_Mesh      = 2
	COMPONENT_Camera    = 3
)

type component struct {
	id         int
	data       interface{}
	dependency []int
}

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}
type Mesh struct {
}
type Camera struct {
}

var idCounter int
var components map[int]map[int]component

var dataTable map[int]interface{}
var dependencyTable map[int][]int

func setUpComponentTables() {
	components = map[int]map[int]component{}

	dataTable = map[int]interface{}{}
	dataTable[COMPONENT_Transform] = Transform{
		Position: mgl32.Vec3{0, 0, 0},
		Rotation: mgl32.Vec3{0, 0, 0},
		Scale:    mgl32.Vec3{1, 1, 1}}
	dataTable[COMPONENT_Mesh] = Mesh{}
	dataTable[COMPONENT_Camera] = Camera{}

	dependencyTable = map[int][]int{}
	dependencyTable[COMPONENT_Transform] = []int{}
	dependencyTable[COMPONENT_Mesh] = []int{}
	dependencyTable[COMPONENT_Camera] = []int{}
}

func CreateEntity() int {
	idCounter++
	components[idCounter] = map[int]component{}
	return idCounter
}
func DeleteEntity(id int) {
	delete(components, id)
}

func AddComponent(entityId int, componentId int) {
	components[entityId][componentId] = component{
		id:         componentId,
		data:       dataTable[componentId],
		dependency: dependencyTable[componentId]}
}

func RemoveComponent(entityId int, componentId int) {
	delete(components[entityId], componentId)
}

func GetComponent(entityId int, componentId int) interface{} {
	return components[entityId][componentId].data
}

func SetComponent(entityId int, componentId int, data interface{}) {
	component := components[entityId][componentId]
	component.data = data
	components[entityId][componentId] = component
}
