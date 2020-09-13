package OctaForce

import "github.com/go-gl/mathgl/mgl32"

const (
	COMPONENT_Transform = 1
	COMPONENT_Mesh      = 2
	COMPONENT_Camera    = 3
)
const (
	component_func_Add    = 0
	component_func_Update = 1
	component_func_Get    = 2
	component_func_Set    = 3
	component_func_Remove = 4
)

type component struct {
	id         int
	data       interface{}
	dependency []int
	dataFuncs  map[int]func(data interface{}) interface{}
}

type Transform struct {
	Position mgl32.Vec3
	Rotation mgl32.Vec3
	Scale    mgl32.Vec3
}
type Camera struct {
}

var idCounter int
var components map[int]map[int]component

var dataTable map[int]interface{}
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(data interface{}) interface{}

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

	funcTable = map[int]map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Transform] = map[int]func(data interface{}) interface{}{}

	funcTable[COMPONENT_Mesh] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Mesh][component_func_Set] = updateMeshData

	funcTable[COMPONENT_Camera] = map[int]func(data interface{}) interface{}{}
}

func CreateEntity() int {
	idCounter++
	components[idCounter] = map[int]component{}
	return idCounter
}
func DeleteEntity(id int) {
	delete(components, id)
}

func HasEntity(entityId int) bool {
	return components[entityId] != nil
}

func updateAllComponents() {
	for i, entity := range components {
		for j, component := range entity {
			if component.dataFuncs[component_func_Update] != nil {
				component.data = component.dataFuncs[component_func_Update](component.data)
				components[i][j] = component
			}
		}
	}
}

func AddComponent(entityId int, componentId int) interface{} {
	component := component{
		id:         componentId,
		data:       dataTable[componentId],
		dependency: dependencyTable[componentId],
		dataFuncs:  funcTable[componentId],
	}

	if component.dataFuncs[component_func_Add] != nil {
		component.data = component.dataFuncs[component_func_Add](component.data)
	}
	components[entityId][componentId] = component
	return component.data
}

func RemoveComponent(entityId int, componentId int) {
	component := components[entityId][componentId]
	if component.dataFuncs[component_func_Remove] != nil {
		component.data = component.dataFuncs[component_func_Remove](component.data)
		components[entityId][componentId] = component
	}
	delete(components[entityId], componentId)
}

func HasComponent(entityId int, componentId int) bool {
	return components[entityId][componentId].data != nil
}

func GetComponent(entityId int, componentId int) interface{} {
	component := components[entityId][componentId]
	if component.dataFuncs[component_func_Get] != nil {
		component.data = component.dataFuncs[component_func_Get](component.data)
		components[entityId][componentId] = component
	}
	return component.data
}

func SetComponent(entityId int, componentId int, data interface{}) {
	component := components[entityId][componentId]
	component.data = data
	components[entityId][componentId] = component

	if component.dataFuncs[component_func_Set] != nil {
		component.data = component.dataFuncs[component_func_Set](component.data)
		components[entityId][componentId] = component
	}
}

func GetAllComponentsOfId(id int) []interface{} {
	var datas []interface{}
	for entityId, _ := range components {
		datas = append(datas, components[entityId][id].data)
	}
	return datas
}
