package main

import (
	"fmt"
	"time"

	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
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

	window, err := glfw.CreateWindow(640, 480, "Testing", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
	w, h := window.GetSize()
	initGL(w, h)
	b := new(skia.Backend)
	b.Init(w, h)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	for !window.ShouldClose() {
		//Do OpenGL stuff
		b.StartLoop()
		b.DrawButton(0, 0, 50, 20, 20)
		b.EndLoop()
		window.SwapBuffers()
		glfw.PollEvents()
		time.Sleep(200 * time.Millisecond)
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
