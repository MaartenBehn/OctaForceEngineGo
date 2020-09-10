package OctaForce

// Slices
func AddFuncToSlice(slice []func(), elementPtr func()) ([]func(), int) {
	slice = append(slice, elementPtr)
	return slice, len(slice) - 1
}
func RemoveFuncFromSlice(slicePtr *[]func(), elementIndex int, keepOrder bool) {
	slice := *slicePtr
	if keepOrder {

		copy(slice[elementIndex:], slice[elementIndex+1:])

	} else {

		slice[elementIndex] = slice[len(slice)-1]
	}
	slice[elementIndex] = nil
	slice = slice[:len(slice)-1]
}
