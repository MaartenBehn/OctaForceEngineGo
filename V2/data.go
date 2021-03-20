package V2

type Data interface {
	GetId() int
	SetId(int)
	copy() Data
}

type data struct {
	id int
}

func (d data) GetId() int {
	return d.id
}
func (d *data) SetId(id int) {
	d.id = id
}
func (d *data) copy() Data {
	return &data{id: d.id}
}
