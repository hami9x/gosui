package native

import (
	"fmt"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	gs "github.com/phaikawl/gosui"
)

//RenderBackend does the real work of drawing and updating graphics
//It may contain code that is specific to a library or backend, like skia
type RenderBackend interface {
	gs.RenderBackend
	Flush()
	UpdateViewportSize(int, int)
	Die()
}

//Window is the application window
type Window struct {
	glw *glfw.Window
	b   RenderBackend

	root *gs.AbstractElement
}

//RootElement gets the root element
func (wn *Window) RootElement() *gs.AbstractElement {
	return wn.root
}

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

//NewWindow creates and return a new Windows
func NewWindow(b RenderBackend, w, h int, title string) *Window {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	window, err := glfw.CreateWindow(w, h, title, nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	b.Init(w, h)
	setupGL(w, h)

	return &Window{window, b, gs.NewRootElement()}
}

//Size of the window
func (wn *Window) Size() (w, h int) {
	return wn.glw.GetSize()
}

//Loop is the main loop for the window, everything happens there
//Call this method to get the window running
func (wn *Window) Start() {
	b := wn.b

	needUpdate := false
	continuousRedraw := true

	defer glfw.Terminate()
	defer wn.b.Die()

	//Will redraw the window again and again in a short while
	//Dirty way to ensure that a we don't get a buggy black window
	time.AfterFunc(500*time.Millisecond, func() {
		continuousRedraw = false
	})

	w, h := wn.Size()
	for !wn.glw.ShouldClose() {
		cw, ch := wn.Size()
		if cw != w || ch != h {
			w, h = cw, ch
			b.UpdateViewportSize(cw, ch)
			needUpdate = true
		}
		if needUpdate || continuousRedraw {
			wn.glw.SwapBuffers()
			gs.Redraw(wn.root, b, wn.root)
			needUpdate = false
			b.Flush()
			wn.glw.SwapBuffers()
		}

		glfw.PollEvents()
	}
}

func setupGL(w, h int) {
	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), float64(h), 0, 0, 1)

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(1, 1, 1, 0)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	gl.Hint(gl.LINE_SMOOTH_HINT|gl.LINE_SMOOTH_HINT, gl.NICEST)
}
