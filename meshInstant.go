package V2

// MeshInstant is a Component that creates an Instant of the Mesh Component of MeshEntity.
// Transform and Material are individual changeable.
type MeshInstant struct {
	mesh     *Mesh
	Material *Material

	Transform *Transform
}

func NewMeshInstant(mesh *Mesh, material *Material) *MeshInstant {
	meshInstant := &MeshInstant{
		mesh:      mesh,
		Material:  material,
		Transform: NewTransform(),
	}
	mesh.instances = append(mesh.instances, meshInstant)
	mesh.needsInstanceUpdate = true
	return meshInstant
}

func (m *MeshInstant) Delete() {
	mesh := m.mesh
	for i := len(mesh.instances) - 1; i >= 0; i-- {
		if mesh.instances[i] == m {
			mesh.instances = append(mesh.instances[:i], mesh.instances[i+1:]...)
		}
	}
	mesh.needsInstanceUpdate = true
}
