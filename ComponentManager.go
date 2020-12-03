package OctaForceEngine

import "sync"

const (
	ComponentTransform   = 1
	ComponentMesh        = 2
	ComponentMeshInstant = 3
	ComponentCamera      = 4
)
const (
	componentFuncAdd    = 0
	componentFuncUpdate = 1
	componentFuncGet    = 2
	componentFuncSet    = 3
	componentFuncRemove = 4
)

type ComponentContainer struct {
	id        int
	entityId  int
	component interface{}
}

var components [][]ComponentContainer
var componentsMutex sync.Mutex
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(component interface{}) interface{}

func setUpComponentTables() {
	componentsMutex = sync.Mutex{}
	componentsMutex.Lock()
	defer componentsMutex.Unlock()

	components = make([][]ComponentContainer, 0)

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

var freedIds []int

// CreateEntity creates a new entity and returns its id.
func CreateEntity() int {
	var id int
	if len(freedIds) > 0 {
		id = freedIds[0]
		freedIds = append(freedIds[:0], freedIds[1:]...)
	} else {

		componentsMutex.Lock()
		components = append(components, make([]ComponentContainer, 0))
		componentsMutex.Unlock()

		id = len(components) - 1
	}
	return id
}

// DeleteEntity deletes the entity of the given id.
func DeleteEntity(id int) {

	componentsMutex.Lock()
	components[id] = make([]ComponentContainer, 0)
	componentsMutex.Unlock()

	freedIds = append(freedIds, id)
}

// HasEntity returns true if entity of id exists.
func HasEntity(entityId int) bool {

	componentsMutex.Lock()
	defer componentsMutex.Unlock()
	return len(components) >= entityId && components[entityId] != nil
}

// GetAllEntitiesWithComponent returns List of all entities with the given component.
func GetAllEntitiesWithComponent(id int) []int {
	var entities []int

	componentsMutex.Lock()
	for entityId, _ := range components {
		for _, component := range components[entityId] {
			if component.id == id {
				entities = append(entities, entityId)
			}
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
	components[entityId] = append(components[entityId], ComponentContainer{
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
	return component
}

// RemoveComponent removes given component from given entity.
// Will not check if any components are dependent on it before it is removed.
func RemoveComponent(entityId int, componentId int) {

	componentsMutex.Lock()
	component := components[entityId][componentId].component
	for i, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			components[entityId] = append(components[entityId][:i], components[entityId][i+1:]...)
		}
	}
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
	for i, componentContainer := range components[entityId] {
		if componentContainer.id == componentId {
			components[entityId][i].component = component
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
		}
	}
	componentsMutex.Unlock()

	if funcTable[componentId][componentFuncGet] != nil {
		newComponent := funcTable[componentId][componentFuncGet](component)

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
			componentSet = append(componentSet, components[entityId][id])
		}
	}
	componentsMutex.Unlock()

	return componentSet
}

func updateAllComponents() {

	componentsMutex.Lock()
	for i, entity := range components {
		for j, componentContainer := range entity {
			if funcTable[j][componentFuncUpdate] != nil {
				componentContainer.component = funcTable[j][componentFuncUpdate](componentContainer.component)
				components[i][j] = componentContainer
			}
		}
	}
	componentsMutex.Unlock()
}
