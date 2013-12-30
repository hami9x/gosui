package main

import (
	"fmt"
	"image"
	"image/color"

	draw2d "code.google.com/p/draw2d/draw2dgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"github.com/phaikawl/gosui"
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
	for !window.ShouldClose() {
		//Do OpenGL stuff
		resetGL()
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		gc := draw2d.NewGraphicContext(w, h)
		gc.SetFillColor(color.RGBA{0x80, 0x80, 0xFF, 0xFF})
		gc.SetStrokeColor(color.RGBA{0x80, 0, 0, 0x80})
		// draw2d.SetFontFolder("font/")
		drawer := gosui.MakeDrawer(gc, img)
		drawer.DrawRoundedRect(0, 0, 100, 30, 10)
		gc.SetFillColor(color.RGBA{0, 0, 0, 255})
		gc.SetStrokeColor(color.RGBA{0, 0, 0, 255})
		// fontData := draw2d.FontData{"luxi", draw2d.FontFamilyMono, draw2d.FontStyleNormal}
		// drawer.DrawText("Awesome!", 10, 10, 10, fontData)
		// drawer.RenderToGL()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func initGL(w, h int) {
	gl.Viewport(0, 0, w, h)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, float64(w), float64(h), 0, 0, 1)

	gl.ShadeModel(gl.SMOOTH)
	gl.ClearColor(1, 1, 1, 0)
	gl.Enable(gl.SMOOTH)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.DEPTH_TEST)
	gl.Hint(gl.LINE_SMOOTH_HINT|gl.LINE_SMOOTH_HINT, gl.NICEST)
}

func resetGL() {
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// Reset the matrix
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
}
