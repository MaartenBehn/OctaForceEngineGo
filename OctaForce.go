package OctaForce

func StartUp(gameStartUpFunc func(), gameStopFunc func()) {
	startUp(gameStartUpFunc, gameStopFunc)
}

//region Callback
func AddUpdateCallback(function func()) {
	addUpdateCallback(function)
}
func RemoveUpdateUpCallback(index int) {
	removeUpdateUpCallback(index)
}

//endregion

//region fps ups
func SetMaxFPS(maxfps float64) {
	maxFPS = maxfps
}
func GetFPS() float64 {
	return fps
}
func GetCappedFPS() float64 {
	if fps > maxFPS {
		fps = maxFPS
	}
	return fps
}

func SetMaxUPS(maxups float64) {
	maxUPS = maxups
}
func GetUPS() float64 {
	return ups
}
func GetCappedUPS() float64 {
	if ups > maxUPS {
		ups = maxUPS
	}
	return ups
}

//endregion
