package OctaForce

type Data interface {
	checkDependency(data Data) bool
}
