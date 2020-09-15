package OctaForceEngine

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
var dependencyTable map[int][]int
var funcTable map[int]map[int]func(data interface{}) interface{}

func setUpComponentTables() {
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
		data: funcTable[componentId][componentFuncAdd](nil),
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
	if funcTable[componentId][componentFuncRemove] != nil {
		component.data = funcTable[componentId][componentFuncRemove](component.data)
		components[entityId][componentId] = component
	}
	delete(components[entityId], componentId)
}

// SetComponent sets the value of given component on given entity.
func SetComponent(entityId int, componentId int, data interface{}) {
	component := components[entityId][componentId]
	if funcTable[componentId][componentFuncSet] != nil {
		component.data = funcTable[componentId][componentFuncSet](data)
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
	if funcTable[componentId][componentFuncGet] != nil {
		component.data = funcTable[componentId][componentFuncGet](component.data)
		components[entityId][componentId] = component
	}
	return component.data
}

// GetAllComponentsOfId returns a slice of all copied values of components of given id.
func GetAllComponentsOfId(id int) []interface{} {
	var dataset []interface{}
	for entityId, _ := range components {
		if HasComponent(entityId, id) {
			dataset = append(dataset, components[entityId][id].data)
		}
	}
	return dataset
}

func updateAllComponents() {
	for i, entity := range components {
		for j, component := range entity {
			if funcTable[j][componentFuncUpdate] != nil {
				component.data = funcTable[j][componentFuncUpdate](component.data)
				components[i][j] = component
			}
		}
	}
}
