package OctaForce

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
	id   int
	data interface{}
}

var components map[int]map[int]component
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(data interface{}) interface{}

func setUpComponentTables() {
	components = map[int]map[int]component{}

	dependencyTable = map[int][]int{}
	dependencyTable[COMPONENT_Transform] = []int{}
	dependencyTable[COMPONENT_Mesh] = []int{COMPONENT_Transform}
	dependencyTable[COMPONENT_Camera] = []int{}

	funcTable = map[int]map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Transform] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Transform][component_func_Add] = setUpTransform
	funcTable[COMPONENT_Transform][component_func_Set] = setTransformMatrix

	funcTable[COMPONENT_Mesh] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Mesh][component_func_Add] = setUpMesh
	funcTable[COMPONENT_Mesh][component_func_Set] = updateMeshData

	funcTable[COMPONENT_Camera] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Camera][component_func_Add] = setUpCamera
}

var idCounter int

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
func GetAllEntitiesWithComponent(id int) []int {
	var entities []int
	for entityId, _ := range components {
		if HasComponent(entityId, id) {
			entities = append(entities, entityId)
		}
	}
	return entities
}

func AddComponent(entityId int, componentId int) interface{} {
	component := component{
		id:   componentId,
		data: funcTable[componentId][component_func_Add](nil),
	}
	components[entityId][componentId] = component

	for _, dependency := range dependencyTable[componentId] {
		if !HasComponent(entityId, dependency) {
			AddComponent(entityId, dependency)
		}
	}
	return component.data
}
func RemoveComponent(entityId int, componentId int) {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Remove] != nil {
		component.data = funcTable[componentId][component_func_Remove](component.data)
		components[entityId][componentId] = component
	}
	delete(components[entityId], componentId)
}
func SetComponent(entityId int, componentId int, data interface{}) {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Set] != nil {
		component.data = funcTable[componentId][component_func_Set](data)
	} else {
		component.data = data
	}
	components[entityId][componentId] = component
}
func HasComponent(entityId int, componentId int) bool {
	return components[entityId][componentId].data != nil
}
func GetComponent(entityId int, componentId int) interface{} {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Get] != nil {
		component.data = funcTable[componentId][component_func_Get](component.data)
		components[entityId][componentId] = component
	}
	return component.data
}
func GetAllComponentsOfId(id int) []interface{} {
	var datas []interface{}
	for entityId, _ := range components {
		if HasComponent(entityId, id) {
			datas = append(datas, components[entityId][id].data)
		}
	}
	return datas
}

func updateAllComponents() {
	for i, entity := range components {
		for j, component := range entity {
			if funcTable[j][component_func_Update] != nil {
				component.data = funcTable[j][component_func_Update](component.data)
				components[i][j] = component
			}
		}
	}
}
