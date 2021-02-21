package OctaForceEngine

import "sync"

const (
	ComponentTransform   = 1
	ComponentMesh        = 2
	ComponentMeshInstant = 3
	ComponentCamera      = 4
	componentCount       = 5
)
const (
	componentFuncAdd    = 0
	componentFuncUpdate = 1
	componentFuncGet    = 2
	componentFuncSet    = 3
	componentFuncRemove = 4
	componentFuncCount  = 5
)

type componentContainer struct {
	id        int
	entityId  int
	component interface{}
}

var components [][]componentContainer

var componentsMutex sync.Mutex
var entityWithComponents [][]int

var dependencyTable [][]int
var funcTable [][]func(component interface{}, entityId int) interface{}

func setUpComponentTables() {
	componentsMutex = sync.Mutex{}
	componentsMutex.Lock()
	defer componentsMutex.Unlock()

	components = make([][]componentContainer, 0)
	entityWithComponents = make([][]int, componentCount)

	dependencyTable = make([][]int, componentCount)
	dependencyTable[ComponentTransform] = []int{}
	dependencyTable[ComponentMesh] = []int{ComponentTransform}
	dependencyTable[ComponentMeshInstant] = []int{ComponentTransform}
	dependencyTable[ComponentCamera] = []int{ComponentTransform}

	funcTable = make([][]func(component interface{}, entityId int) interface{}, componentCount)
	funcTable[ComponentTransform] = make([]func(component interface{}, entityId int) interface{}, componentFuncCount)
	funcTable[ComponentTransform][componentFuncAdd] = setUpTransform
	funcTable[ComponentTransform][componentFuncSet] = setTransformMatrix

	funcTable[ComponentMesh] = make([]func(component interface{}, entityId int) interface{}, componentFuncCount)
	funcTable[ComponentMesh][componentFuncAdd] = setUpMesh
	funcTable[ComponentMesh][componentFuncRemove] = deleteMesh

	funcTable[ComponentMeshInstant] = make([]func(component interface{}, entityId int) interface{}, componentFuncCount)
	funcTable[ComponentMeshInstant][componentFuncAdd] = setUpMeshInstant
	funcTable[ComponentMeshInstant][componentFuncSet] = addMeshInstant
	funcTable[ComponentMeshInstant][componentFuncRemove] = removeMeshInstant

	funcTable[ComponentCamera] = make([]func(component interface{}, entityId int) interface{}, componentFuncCount)
	funcTable[ComponentCamera][componentFuncAdd] = setUpCamera
}

var freedIds []int

// CreateEntity creates a new entity and returns its id.
func CreateEntity() int {
	var id int
	if len(freedIds) > 0 {
		id = freedIds[0]
		freedIds = append(freedIds[:0], freedIds[1:]...)
	} else {

		componentsMutex.Lock()
		components = append(components, make([]componentContainer, 0))
		componentsMutex.Unlock()

		id = len(components) - 1
	}
	return id
}

// DeleteEntity deletes the entity of the given id.
func DeleteEntity(id int) {

	componentsMutex.Lock()
	components[id] = make([]componentContainer, 0)
	componentsMutex.Unlock()

	freedIds = append(freedIds, id)
}

// HasEntity returns true if entity of id exists.
func HasEntity(entityId int) bool {

	componentsMutex.Lock()
	hasEntity := len(components) >= entityId && components[entityId] != nil
	componentsMutex.Unlock()

	return hasEntity
}

// GetAllEntitiesWithComponent returns List of all entities with the given component.
func GetAllEntitiesWithComponent(id int) []int {
	return entityWithComponents[id]
}

// AddComponent adds the given component to the given entity.
// Also adds any component the given component is dependent on, when they aren't already added.
func AddComponent(entityId int, componentId int) interface{} {

	if HasComponent(entityId, componentId) {
		return nil
	}
	component := funcTable[componentId][componentFuncAdd](nil, entityId)

	componentsMutex.Lock()
	components[entityId] = append(components[entityId], componentContainer{
		id:        componentId,
		entityId:  entityId,
		component: component,
	})
	componentsMutex.Unlock()

	for _, dependency := range dependencyTable[componentId] {
		if !HasComponent(entityId, dependency) {
			AddComponent(entityId, dependency)
		}
	}

	entityWithComponents[componentId] = append(entityWithComponents[componentId], entityId)

	return component
}

// RemoveComponent removes given component from given entity.
// Will not check if any components are dependent on it before it is removed.
func RemoveComponent(entityId int, componentId int) {
	id := -1

	componentsMutex.Lock()
	component := components[entityId][componentId].component
	for _, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			id = componentContainer.id
		}
	}

	if id <= 0 {

		components[entityId] = append(components[entityId][:id], components[entityId][id+1:]...)

		if funcTable[componentId][componentFuncRemove] != nil {
			funcTable[componentId][componentFuncRemove](component, entityId)
		}
	}
	componentsMutex.Unlock()

	entityIndex := -1
	for i, id := range entityWithComponents[componentId] {
		if id == entityId {
			entityIndex = i
		}
	}
	entityWithComponents[componentId] = append(
		entityWithComponents[componentId][:entityIndex],
		entityWithComponents[componentId][entityIndex+1:]...)
}

// SetComponent sets the value of given component on given entity.
func SetComponent(entityId int, componentId int, component interface{}) {
	if funcTable[componentId][componentFuncSet] != nil {
		component = funcTable[componentId][componentFuncSet](component, entityId)
	}

	componentsMutex.Lock()
	for i, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			components[entityId][i].component = component
			break
		}
	}
	componentsMutex.Unlock()
}

// HasComponent return true if given component on given entity exists.
func HasComponent(entityId int, componentId int) bool {

	hasComponent := false

	componentsMutex.Lock()
	for _, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			hasComponent = true
			break
		}
	}
	componentsMutex.Unlock()

	return hasComponent
}

// GetComponent returns copy of the value of given component on given entity.
func GetComponent(entityId int, componentId int) interface{} {

	var component interface{}
	componentsMutex.Lock()
	for _, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			component = componentContainer.component
			break
		}
	}
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncGet] != nil {
		newComponent := funcTable[componentId][componentFuncGet](component, entityId)

		if newComponent != component {
			component = newComponent

			componentsMutex.Lock()
			for i, componentContainer := range components[entityId] {
				if componentContainer.id == componentId {
					components[entityId][i].component = component
				}
			}
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
		if components[entityId][id].component != nil {
			componentSet = append(componentSet, components[entityId][id].component)
		}
	}
	componentsMutex.Unlock()

	return componentSet
}

func updateAllComponents() {

	componentsMutex.Lock()
	for i, entity := range components {
		for j, componentContainer := range entity {
			if funcTable[componentContainer.id][componentFuncUpdate] != nil {
				componentContainer.component =
					funcTable[componentContainer.id][componentFuncUpdate](componentContainer.component, i)
				components[i][j] = componentContainer
			}
		}
	}
	componentsMutex.Unlock()
}
