package OctaForceEngine

import "C"
import (
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	windowWidth  = 1280
	windowHeight = 720
)

var window *glfw.Window

func setUpWindow(name string) {
	var err error

	// Setting up Window
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err = glfw.CreateWindow(windowWidth, windowHeight, name, nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
}

func renderWindow() {
	window.SwapBuffers()
	glfw.PollEvents()
}

func updateWindow() {
	if window.ShouldClose() {
		running = false
	}

	// Mouse Info
	mouseX, mouseY := window.GetCursorPos()
	mouseMovement = mgl32.Vec2{float32(mouseX) - mousePos.X(), float32(mouseY) - mousePos.Y()}
	mousePos = mgl32.Vec2{float32(mouseX), float32(mouseY)}
}

// KeyPressed returns true when the Key is pressed.
func KeyPressed(key int) bool {
	return window.GetKey(glfw.Key(key)) == glfw.Press
}

const (
	KeyApostrophe   int = int(glfw.KeyApostrophe)
	KeyComma        int = int(glfw.KeyComma)
	KeyMinus        int = int(glfw.KeyMinus)
	KeyPeriod       int = int(glfw.KeyPeriod)
	KeySlash        int = int(glfw.KeySlash)
	Key0            int = int(glfw.Key0)
	Key1            int = int(glfw.Key1)
	Key2            int = int(glfw.Key2)
	Key3            int = int(glfw.Key3)
	Key4            int = int(glfw.Key4)
	Key5            int = int(glfw.Key5)
	Key6            int = int(glfw.Key6)
	Key7            int = int(glfw.Key7)
	Key8            int = int(glfw.Key8)
	Key9            int = int(glfw.Key9)
	KeySemicolon    int = int(glfw.KeySemicolon)
	KeyEqual        int = int(glfw.KeyEqual)
	KeyA            int = int(glfw.KeyA)
	KeyB            int = int(glfw.KeyB)
	KeyC            int = int(glfw.KeyC)
	KeyD            int = int(glfw.KeyD)
	KeyE            int = int(glfw.KeyE)
	KeyF            int = int(glfw.KeyF)
	KeyG            int = int(glfw.KeyG)
	KeyH            int = int(glfw.KeyH)
	KeyI            int = int(glfw.KeyI)
	KeyJ            int = int(glfw.KeyJ)
	KeyK            int = int(glfw.KeyK)
	KeyL            int = int(glfw.KeyL)
	KeyM            int = int(glfw.KeyM)
	KeyN            int = int(glfw.KeyN)
	KeyO            int = int(glfw.KeyO)
	KeyP            int = int(glfw.KeyP)
	KeyQ            int = int(glfw.KeyQ)
	KeyR            int = int(glfw.KeyR)
	KeyS            int = int(glfw.KeyS)
	KeyT            int = int(glfw.KeyT)
	KeyU            int = int(glfw.KeyU)
	KeyV            int = int(glfw.KeyV)
	KeyW            int = int(glfw.KeyW)
	KeyX            int = int(glfw.KeyX)
	KeyY            int = int(glfw.KeyY)
	KeyZ            int = int(glfw.KeyZ)
	KeyLeftBracket  int = int(glfw.KeyLeftBracket)
	KeyBackslash    int = int(glfw.KeyBackslash)
	KeyRightBracket int = int(glfw.KeyRightBracket)
	KeyGraveAccent  int = int(glfw.KeyGraveAccent)
	KeyWorld1       int = int(glfw.KeyWorld1)
	KeyWorld2       int = int(glfw.KeyWorld2)
	KeyEscape       int = int(glfw.KeyEscape)
	KeyEnter        int = int(glfw.KeyEnter)
	KeyTab          int = int(glfw.KeyTab)
	KeyBackspace    int = int(glfw.KeyBackspace)
	KeyInsert       int = int(glfw.KeyInsert)
	KeyDelete       int = int(glfw.KeyDelete)
	KeyRight        int = int(glfw.KeyRight)
	KeyLeft         int = int(glfw.KeyLeft)
	KeyDown         int = int(glfw.KeyDown)
	KeyUp           int = int(glfw.KeyUp)
	KeyPageUp       int = int(glfw.KeyPageUp)
	KeyPageDown     int = int(glfw.KeyPageDown)
	KeyHome         int = int(glfw.KeyHome)
	KeyEnd          int = int(glfw.KeyEnd)
	KeyCapsLock     int = int(glfw.KeyCapsLock)
	KeyScrollLock   int = int(glfw.KeyScrollLock)
	KeyNumLock      int = int(glfw.KeyNumLock)
	KeyPrintScreen  int = int(glfw.KeyPrintScreen)
	KeyPause        int = int(glfw.KeyPause)
	KeyF1           int = int(glfw.KeyF1)
	KeyF2           int = int(glfw.KeyF2)
	KeyF3           int = int(glfw.KeyF3)
	KeyF4           int = int(glfw.KeyF4)
	KeyF5           int = int(glfw.KeyF5)
	KeyF6           int = int(glfw.KeyF6)
	KeyF7           int = int(glfw.KeyF7)
	KeyF8           int = int(glfw.KeyF8)
	KeyF9           int = int(glfw.KeyF9)
	KeyF10          int = int(glfw.KeyF10)
	KeyF11          int = int(glfw.KeyF11)
	KeyF12          int = int(glfw.KeyF12)
	KeyF13          int = int(glfw.KeyF13)
	KeyF14          int = int(glfw.KeyF14)
	KeyF15          int = int(glfw.KeyF15)
	KeyF16          int = int(glfw.KeyF16)
	KeyF17          int = int(glfw.KeyF17)
	KeyF18          int = int(glfw.KeyF17)
	KeyF19          int = int(glfw.KeyF19)
	KeyF20          int = int(glfw.KeyF20)
	KeyF21          int = int(glfw.KeyF21)
	KeyF22          int = int(glfw.KeyF22)
	KeyF23          int = int(glfw.KeyF23)
	KeyF24          int = int(glfw.KeyF24)
	KeyF25          int = int(glfw.KeyF25)
	KeyKP0          int = int(glfw.KeyF25)
	KeyKP1          int = int(glfw.KeyKP1)
	KeyKP2          int = int(glfw.KeyKP2)
	KeyKP3          int = int(glfw.KeyKP3)
	KeyKP4          int = int(glfw.KeyKP4)
	KeyKP5          int = int(glfw.KeyKP5)
	KeyKP6          int = int(glfw.KeyKP6)
	KeyKP7          int = int(glfw.KeyKP7)
	KeyKP8          int = int(glfw.KeyKP8)
	KeyKP9          int = int(glfw.KeyKP9)
	KeyKPDecimal    int = int(glfw.KeyKPDecimal)
	KeyKPDivide     int = int(glfw.KeyKPDivide)
	KeyKPMultiply   int = int(glfw.KeyKPMultiply)
	KeyKPSubtract   int = int(glfw.KeyKPSubtract)
	KeyKPAdd        int = int(glfw.KeyKPAdd)
	KeyKPEnter      int = int(glfw.KeyKPEnter)
	KeyKPEqual      int = int(glfw.KeyKPEqual)
	KeyLeftShift    int = int(glfw.KeyLeftShift)
	KeyLeftControl  int = int(glfw.KeyLeftControl)
	KeyLeftAlt      int = int(glfw.KeyLeftAlt)
	KeyLeftSuper    int = int(glfw.KeyLeftSuper)
	KeyRightShift   int = int(glfw.KeyRightShift)
	KeyRightControl int = int(glfw.KeyRightControl)
	KeyRightAlt     int = int(glfw.KeyRightAlt)
	KeyRightSuper   int = int(glfw.KeyRightSuper)
	KeyMenu         int = int(glfw.KeyMenu)
	KeyLast         int = int(glfw.KeyLast)
)

// MouseButtonPressed returns true when the mouse key is pressed.
func MouseButtonPressed(button int) bool {
	return window.GetMouseButton(glfw.MouseButton(button)) == glfw.Press
}

const (
	MouseButton2      int = int(glfw.MouseButton2)
	MouseButton3      int = int(glfw.MouseButton3)
	MouseButton4      int = int(glfw.MouseButton4)
	MouseButton5      int = int(glfw.MouseButton5)
	MouseButton6      int = int(glfw.MouseButton6)
	MouseButton7      int = int(glfw.MouseButton7)
	MouseButton8      int = int(glfw.MouseButton8)
	MouseButtonLast   int = int(glfw.MouseButtonLast)
	MouseButtonLeft   int = int(glfw.MouseButtonLeft)
	MouseButtonRight  int = int(glfw.MouseButtonRight)
	MouseButtonMiddle int = int(glfw.MouseButtonMiddle)
)

var mousePos mgl32.Vec2

// GetMousePos returns the position of the Mouse relative to the screen.
func GetMousePos() mgl32.Vec2 {
	return mousePos
}

var mouseMovement mgl32.Vec2

// GetMouseMovement returns the relative movement since the update.
func GetMouseMovement() mgl32.Vec2 {
	return mouseMovement
}
