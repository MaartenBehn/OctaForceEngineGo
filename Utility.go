package OctaForce

// Slices
func AddFuncToSlice(slicePtr *[]func(), elementPtr func()) int {
	slice := *slicePtr
	slice = append(slice, elementPtr)
	return len(slice) - 1
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
