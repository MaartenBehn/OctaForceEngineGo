package OctaForceEngine

const (
	ComponentTransform   = 1
	ComponentMesh        = 2
	ComponentMeshInstant = 3
	ComponentCamera      = 4
	componentMax         = 5
)
const (
	componentFuncAdd    = 0
	componentFuncUpdate = 1
	componentFuncSet    = 2
	componentFuncRemove = 3
	componentFuncMax    = 4
)

type World struct {
	entities              []*Entity
	entitiesWithComponent [][]int
	writeFuncs            chan func(world *World)
}

func setUpWorld() {
	world.entitiesWithComponent = make([][]int, componentMax)
	world.writeFuncs = make(chan func(world *World), 1000)
}

func pushWriteFunc(function func(world *World)) {
	world.writeFuncs <- function
}

var world World

type Entity struct {
	components []interface{}
}

var dependencyTable [][]int
var funcTable [][]func(component interface{}, entityId int) interface{}

func setUpComponentTables() {
	dependencyTable = make([][]int, componentMax)
	dependencyTable[ComponentTransform] = []int{}
	dependencyTable[ComponentMesh] = []int{ComponentTransform}
	dependencyTable[ComponentMeshInstant] = []int{ComponentTransform}
	dependencyTable[ComponentCamera] = []int{ComponentTransform}

	funcTable = make([][]func(component interface{}, entityId int) interface{}, componentMax)
	funcTable[ComponentTransform] = make([]func(component interface{}, entityId int) interface{}, componentFuncMax)
	funcTable[ComponentTransform][componentFuncAdd] = setUpTransform
	funcTable[ComponentTransform][componentFuncSet] = setTransformMatrix

	funcTable[ComponentMesh] = make([]func(component interface{}, entityId int) interface{}, componentFuncMax)
	funcTable[ComponentMesh][componentFuncAdd] = setUpMesh
	funcTable[ComponentMesh][componentFuncRemove] = deleteMesh

	funcTable[ComponentMeshInstant] = make([]func(component interface{}, entityId int) interface{}, componentFuncMax)
	funcTable[ComponentMeshInstant][componentFuncAdd] = setUpMeshInstant
	funcTable[ComponentMeshInstant][componentFuncSet] = addMeshInstant
	funcTable[ComponentMeshInstant][componentFuncRemove] = removeMeshInstant

	funcTable[ComponentCamera] = make([]func(component interface{}, entityId int) interface{}, componentFuncMax)
	funcTable[ComponentCamera][componentFuncAdd] = setUpCamera
}

var freeEntityIds []int

func CreateEntity() int {
	entity := &Entity{make([]interface{}, componentMax)}

	var id int
	if len(freeEntityIds) > 0 {
		id = freeEntityIds[0]

		pushWriteFunc(func(world *World) {
			freeEntityIds = append(freeEntityIds[:0], freeEntityIds[1:]...)
			world.entities[id] = entity
		})
	} else {
		pushWriteFunc(func(world *World) {
			world.entities = append(world.entities, entity)
		})
		id = len(world.entities) - 1
	}
	return id
}

func DeleteEntity(id int) {
	pushWriteFunc(func(world *World) {
		for i := 0; i < len(world.entitiesWithComponent); i++ {
			for j := 0; j < len(world.entitiesWithComponent[i]); j++ {
				if world.entitiesWithComponent[i][j] == id {
					world.entitiesWithComponent[i] = append(world.entitiesWithComponent[i][:j], world.entitiesWithComponent[i][j+1:]...)
					break
				}
			}
		}

		world.entities[id] = nil
		freeEntityIds = append(freeEntityIds, id)
	})
}

func AddComponent(entityId int, componentId int, component interface{}) {
	if !HasComponent(entityId, componentId) {
		return
	}

	pushWriteFunc(func(world *World) {
		world.entities[entityId].components[componentId] = funcTable[componentId][componentFuncAdd](component, entityId)

		for _, dependency := range dependencyTable[componentId] {
			if !HasComponent(entityId, dependency) {
				AddComponent(entityId, dependency, nil)
			}
		}

		world.entitiesWithComponent[componentId] = append(world.entitiesWithComponent[componentId], entityId)
	})
}

func HasComponent(entityId int, componentId int) bool {
	return world.entities[entityId].components[componentId] != nil
}

func GetComponent(entityId int, componentId int) interface{} {
	return world.entities[entityId].components[componentId]
}

func SetComponent(entityId int, componentId int, component interface{}) {
	if !HasComponent(entityId, componentId) {
		AddComponent(entityId, componentId, nil)
	}

	pushWriteFunc(func(world *World) {
		if funcTable[componentId][componentFuncSet] != nil {
			component = funcTable[componentId][componentFuncSet](component, entityId)
		}
		world.entities[entityId].components[componentId] = component
	})
}
func setComponentInternal(entityId int, componentId int, component interface{}) {
	if funcTable[componentId][componentFuncSet] != nil {
		component = funcTable[componentId][componentFuncSet](component, entityId)
	}
	world.entities[entityId].components[componentId] = component
}

func RemoveComponent(entityId int, componentId int) {
	pushWriteFunc(func(world *World) {
		funcTable[componentId][componentFuncRemove](world.entities[entityId].components[componentId], entityId)

		for i := 0; i < len(world.entitiesWithComponent[componentId]); i++ {
			if world.entitiesWithComponent[componentId][i] == entityId {
				world.entitiesWithComponent[componentId] = append(world.entitiesWithComponent[componentId][:i], world.entitiesWithComponent[componentId][i+1:]...)
				break
			}
		}

		world.entities[entityId].components[componentId] = nil
	})
}

func (w *World) update() {
	for i := range w.entities {
		updateEnitity(i)
	}
}
func updateEnitity(entityId int) {

	entity := world.entities[entityId]
	for i, component := range entity.components {
		if HasComponent(entityId, i) {
			funcTable[i][componentFuncUpdate](component, entityId)
		}
	}
}

func (w *World) write() {
	for writeFunc := range w.writeFuncs {
		writeFunc(w)
	}
}

func (w *World) getAllEntityIdsWithComponent(componentId int) []int {
	return w.entitiesWithComponent[componentId]
}
