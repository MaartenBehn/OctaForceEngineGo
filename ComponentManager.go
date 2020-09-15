package OctaForceEngine

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
	dependencyTable[COMPONENT_Camera] = []int{COMPONENT_Transform}

	funcTable = map[int]map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Transform] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Transform][component_func_Add] = setUpTransform
	funcTable[COMPONENT_Transform][component_func_Set] = setTransformMatrix

	funcTable[COMPONENT_Mesh] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Mesh][component_func_Add] = setUpMesh
	funcTable[COMPONENT_Mesh][component_func_Remove] = deleteMesh

	funcTable[COMPONENT_Camera] = map[int]func(data interface{}) interface{}{}
	funcTable[COMPONENT_Camera][component_func_Add] = setUpCamera
}

var idCounter int

// CreateEntity creates a new entity and returns its id.
func CreateEntity() int {
	idCounter++
	components[idCounter] = map[int]component{}
	return idCounter
}

// DeleteEntity deletes the entity of the given id.
func DeleteEntity(id int) {
	delete(components, id)
}

// HasEntity returns true if entity of id exists.
func HasEntity(entityId int) bool {
	return components[entityId] != nil
}

// GetAllEntitiesWithComponent returns List of all entities with the given component.
func GetAllEntitiesWithComponent(id int) []int {
	var entities []int
	for entityId, _ := range components {
		if HasComponent(entityId, id) {
			entities = append(entities, entityId)
		}
	}
	return entities
}

// AddComponent adds the given component to the given entity.
// Also adds any component the given component is dependent on, when they aren't already added.
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

// RemoveComponent removes given component from given entity.
// Will not check if any components are dependent on it before it is removed.
func RemoveComponent(entityId int, componentId int) {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Remove] != nil {
		component.data = funcTable[componentId][component_func_Remove](component.data)
		components[entityId][componentId] = component
	}
	delete(components[entityId], componentId)
}

// SetComponent sets the value of given component on given entity.
func SetComponent(entityId int, componentId int, data interface{}) {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Set] != nil {
		component.data = funcTable[componentId][component_func_Set](data)
	} else {
		component.data = data
	}
	components[entityId][componentId] = component
}

// HasComponent return true if given component on given entity exists.
func HasComponent(entityId int, componentId int) bool {
	return components[entityId][componentId].data != nil
}

// GetComponent returns copy of the value of given component on given entity.
func GetComponent(entityId int, componentId int) interface{} {
	component := components[entityId][componentId]
	if funcTable[componentId][component_func_Get] != nil {
		component.data = funcTable[componentId][component_func_Get](component.data)
		components[entityId][componentId] = component
	}
	return component.data
}

// GetAllComponentsOfId returns a slice of all copied values of components of given id.
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
