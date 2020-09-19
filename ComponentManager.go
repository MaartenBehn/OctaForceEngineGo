package OctaForceEngine

import "sync"

const (
	ComponentTransform = 1
	ComponentMesh      = 2
	ComponentCamera    = 3
)
const (
	componentFuncAdd    = 0
	componentFuncUpdate = 1
	componentFuncGet    = 2
	componentFuncSet    = 3
	componentFuncRemove = 4
)

type component struct {
	id   int
	data interface{}
}

var components map[int]map[int]component
var componentsMutex sync.Mutex
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(data interface{}) interface{}

func setUpComponentTables() {
	componentsMutex = sync.Mutex{}
	componentsMutex.Lock()
	defer componentsMutex.Unlock()

	components = map[int]map[int]component{}

	dependencyTable = map[int][]int{}
	dependencyTable[ComponentTransform] = []int{}
	dependencyTable[ComponentMesh] = []int{ComponentTransform}
	dependencyTable[ComponentCamera] = []int{ComponentTransform}

	funcTable = map[int]map[int]func(data interface{}) interface{}{}
	funcTable[ComponentTransform] = map[int]func(data interface{}) interface{}{}
	funcTable[ComponentTransform][componentFuncAdd] = setUpTransform
	funcTable[ComponentTransform][componentFuncSet] = setTransformMatrix

	funcTable[ComponentMesh] = map[int]func(data interface{}) interface{}{}
	funcTable[ComponentMesh][componentFuncAdd] = setUpMesh
	funcTable[ComponentMesh][componentFuncRemove] = deleteMesh

	funcTable[ComponentCamera] = map[int]func(data interface{}) interface{}{}
	funcTable[ComponentCamera][componentFuncAdd] = setUpCamera
}

var idCounter int

// CreateEntity creates a new entity and returns its id.
func CreateEntity() int {
	idCounter++
	componentsMutex.Lock()
	components[idCounter] = map[int]component{}
	componentsMutex.Unlock()
	return idCounter
}

// DeleteEntity deletes the entity of the given id.
func DeleteEntity(id int) {
	componentsMutex.Lock()
	defer componentsMutex.Unlock()
	delete(components, id)
}

// HasEntity returns true if entity of id exists.
func HasEntity(entityId int) bool {
	componentsMutex.Lock()
	defer componentsMutex.Unlock()
	return components[entityId] != nil
}

// GetAllEntitiesWithComponent returns List of all entities with the given component.
func GetAllEntitiesWithComponent(id int) []int {
	var entities []int

	componentsMutex.Lock()
	for entityId, _ := range components {
		if components[entityId][id].data != nil {
			entities = append(entities, entityId)
		}
	}
	componentsMutex.Unlock()

	return entities
}

// AddComponent adds the given component to the given entity.
// Also adds any component the given component is dependent on, when they aren't already added.
func AddComponent(entityId int, componentId int) interface{} {
	component := component{
		id:   componentId,
		data: funcTable[componentId][componentFuncAdd](nil),
	}

	componentsMutex.Lock()
	components[entityId][componentId] = component
	componentsMutex.Unlock()

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

	componentsMutex.Lock()
	component := components[entityId][componentId]
	delete(components[entityId], componentId)
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncRemove] != nil {
		funcTable[componentId][componentFuncRemove](component.data)
	}
}

// SetComponent sets the value of given component on given entity.
func SetComponent(entityId int, componentId int, data interface{}) {

	componentsMutex.Lock()
	component := components[entityId][componentId]
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncSet] != nil {
		component.data = funcTable[componentId][componentFuncSet](data)
	} else {
		component.data = data
	}

	componentsMutex.Lock()
	components[entityId][componentId] = component
	componentsMutex.Unlock()
}

// HasComponent return true if given component on given entity exists.
func HasComponent(entityId int, componentId int) bool {
	componentsMutex.Lock()
	defer componentsMutex.Unlock()
	return components[entityId][componentId].data != nil
}

// GetComponent returns copy of the value of given component on given entity.
func GetComponent(entityId int, componentId int) interface{} {

	componentsMutex.Lock()
	component := components[entityId][componentId]
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncGet] != nil {
		data := funcTable[componentId][componentFuncGet](component.data)

		if data != component.data {
			component.data = data

			componentsMutex.Lock()
			components[entityId][componentId] = component
			componentsMutex.Unlock()
		}
	}
	return component.data
}

// GetAllComponentsOfId returns a slice of all copied values of components of given id.
func GetAllComponentsOfId(id int) []interface{} {
	var dataset []interface{}

	componentsMutex.Lock()
	for entityId, _ := range components {
		if components[entityId][id].data != nil {
			dataset = append(dataset, components[entityId][id].data)
		}
	}
	componentsMutex.Unlock()

	return dataset
}

func updateAllComponents() {

	componentsMutex.Lock()
	for i, entity := range components {
		for j, component := range entity {
			if funcTable[j][componentFuncUpdate] != nil {
				component.data = funcTable[j][componentFuncUpdate](component.data)
				components[i][j] = component
			}
		}
	}
	componentsMutex.Unlock()
}
