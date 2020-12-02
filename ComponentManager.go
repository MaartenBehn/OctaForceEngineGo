package OctaForceEngine

import "sync"

const (
	ComponentTransform   = 1
	ComponentMesh        = 2
	ComponentMeshInstant = 4
	ComponentCamera      = 5
)
const (
	componentFuncAdd    = 0
	componentFuncUpdate = 1
	componentFuncGet    = 2
	componentFuncSet    = 3
	componentFuncRemove = 4
)

var components map[int]map[int]interface{}
var componentsMutex sync.Mutex
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(component interface{}) interface{}

func setUpComponentTables() {
	componentsMutex = sync.Mutex{}
	componentsMutex.Lock()
	defer componentsMutex.Unlock()

	components = map[int]map[int]interface{}{}

	dependencyTable = map[int][]int{}
	dependencyTable[ComponentTransform] = []int{}
	dependencyTable[ComponentMesh] = []int{ComponentTransform}
	dependencyTable[ComponentMeshInstant] = []int{ComponentTransform}
	dependencyTable[ComponentCamera] = []int{ComponentTransform}

	funcTable = map[int]map[int]func(component interface{}) interface{}{}
	funcTable[ComponentTransform] = map[int]func(component interface{}) interface{}{}
	funcTable[ComponentTransform][componentFuncAdd] = setUpTransform
	funcTable[ComponentTransform][componentFuncSet] = setTransformMatrix

	funcTable[ComponentMesh] = map[int]func(component interface{}) interface{}{}
	funcTable[ComponentMesh][componentFuncAdd] = setUpMesh
	funcTable[ComponentMesh][componentFuncRemove] = deleteMesh

	funcTable[ComponentMeshInstant] = map[int]func(component interface{}) interface{}{}
	funcTable[ComponentMeshInstant][componentFuncAdd] = setUpMeshInstant
	funcTable[ComponentMeshInstant][componentFuncSet] = addMeshInstant
	funcTable[ComponentMeshInstant][componentFuncRemove] = removeMeshInstant

	funcTable[ComponentCamera] = map[int]func(component interface{}) interface{}{}
	funcTable[ComponentCamera][componentFuncAdd] = setUpCamera
}

var idCounter int
var freedIds []int

// CreateEntity creates a new entity and returns its id.
func CreateEntity() int {
	var id int
	if len(freedIds) > 0 {
		id = freedIds[0]
		freedIds = append(freedIds[:0], freedIds[1:]...)
	} else {
		idCounter++
		id = idCounter
	}

	componentsMutex.Lock()
	components[id] = map[int]interface{}{}
	componentsMutex.Unlock()

	return id
}

// DeleteEntity deletes the entity of the given id.
func DeleteEntity(id int) {

	componentsMutex.Lock()
	delete(components, id)
	componentsMutex.Unlock()

	freedIds = append(freedIds, id)
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
		if components[entityId][id] != nil {
			entities = append(entities, entityId)
		}
	}
	componentsMutex.Unlock()

	return entities
}

// AddComponent adds the given component to the given entity.
// Also adds any component the given component is dependent on, when they aren't already added.
func AddComponent(entityId int, componentId int) interface{} {

	component := funcTable[componentId][componentFuncAdd](nil)

	componentsMutex.Lock()
	components[entityId][componentId] = component
	componentsMutex.Unlock()

	for _, dependency := range dependencyTable[componentId] {
		if !HasComponent(entityId, dependency) {
			AddComponent(entityId, dependency)
		}
	}
	return component
}

// RemoveComponent removes given component from given entity.
// Will not check if any components are dependent on it before it is removed.
func RemoveComponent(entityId int, componentId int) {

	componentsMutex.Lock()
	component := components[entityId][componentId]
	delete(components[entityId], componentId)
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncRemove] != nil {
		funcTable[componentId][componentFuncRemove](component)
	}
}

// SetComponent sets the value of given component on given entity.
func SetComponent(entityId int, componentId int, component interface{}) {
	if funcTable[componentId][componentFuncSet] != nil {
		component = funcTable[componentId][componentFuncSet](component)
	}

	componentsMutex.Lock()
	components[entityId][componentId] = component
	componentsMutex.Unlock()
}

// HasComponent return true if given component on given entity exists.
func HasComponent(entityId int, componentId int) bool {
	componentsMutex.Lock()
	defer componentsMutex.Unlock()
	return components[entityId][componentId] != nil
}

// GetComponent returns copy of the value of given component on given entity.
func GetComponent(entityId int, componentId int) interface{} {

	componentsMutex.Lock()
	component := components[entityId][componentId]
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncGet] != nil {
		newComponent := funcTable[componentId][componentFuncGet](component)

		if newComponent != component {
			component = newComponent

			componentsMutex.Lock()
			components[entityId][componentId] = component
			componentsMutex.Unlock()
		}
	}
	return component
}

// GetAllComponentsOfId returns a slice of all copied values of components of given id.
func GetAllComponentsOfId(id int) []interface{} {
	var componentSet []interface{}

	componentsMutex.Lock()
	for entityId, _ := range components {
		if components[entityId][id] != nil {
			componentSet = append(componentSet, components[entityId][id])
		}
	}
	componentsMutex.Unlock()

	return componentSet
}

func updateAllComponents() {

	componentsMutex.Lock()
	for i, entity := range components {
		for j, component := range entity {
			if funcTable[j][componentFuncUpdate] != nil {
				component = funcTable[j][componentFuncUpdate](component)
				components[i][j] = component
			}
		}
	}
	componentsMutex.Unlock()
}
