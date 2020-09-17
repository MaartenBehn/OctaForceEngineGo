package OctaForceEngine

const (
	EventCallbackKey         = 1
	EventCallbackMouseButton = 2
)

var eventCallbackMap map[int]func()

// SetEventCallback set the given function in a map.
// The function will be called when the engine recives the given event.
func SetEventCallback(eventId int, callbackFunc func()) {
	eventCallbackMap[eventId] = callbackFunc
}
