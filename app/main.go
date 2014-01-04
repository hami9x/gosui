package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	gs "github.com/phaikawl/gosui"
	skia "github.com/phaikawl/gosui/skia"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

func main() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}
	defer glfw.Terminate()

	w, h := 800, 600
	window, err := glfw.CreateWindow(w, h, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	initGL(w, h)
	b := new(skia.Backend)
	b.Init(w, h)

	root := gs.NewRootElement()
	bg := gs.NewRectElement(root, gs.MakeRectWH(0, 0, w, h))
	bg.FillColor = gs.Color{255, 255, 255, 255}
	bg.ZIndex = -1000000
	rect := gs.NewRectElement(root, gs.MakeRectWH(10, 10, 100, 100))
	rect.FillColor = gs.Color{255, 0, 0, 255}
	rect.Appr.(*gs.RectShape).SetAllRadii(15)
	needUpdate := false
	continuousRedraw := true
	time.AfterFunc(200*time.Millisecond, func() {
		continuousRedraw = false
	})

	for !window.ShouldClose() {
		cw, ch := window.GetSize()
		if cw != w || ch != h {
			w, h = cw, ch
			b.UpdateWindowSize(cw, ch)
			needUpdate = true
		}
		if needUpdate || continuousRedraw {
			window.SwapBuffers()
			root.Redraw(b, root)
			needUpdate = false
			b.Flush()
			window.SwapBuffers()
		}
		//Do OpenGL stuff
		glfw.PollEvents()
	}
	b.Die()
}

func initGL(w, h int) {
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
