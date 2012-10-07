package main

import (
	"fmt"
	"os"
	"image/color"

	"github.com/banthar/gl"
	"code.google.com/p/draw2d/draw2dgl"
	"github.com/jteeuwen/glfw"
	"phaikawl/wego"
)

func main() {
	var err error
	if err = glfw.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	// Ensure glfw is cleanly terminated on exit.
	defer glfw.Terminate()

	const (
		w = 500
		h = 500
	)

	if err = glfw.OpenWindow(w, h, 8, 8, 8, 0, 0, 0, glfw.Windowed); err != nil {
		fmt.Fprintf(os.Stderr, "[e] %v\n", err)
		return
	}

	// Ensure window is cleanly closed on exit.
	defer glfw.CloseWindow()

	// Enable vertical sync on cards that support it.
	glfw.SetSwapInterval(1)

	// Set window title
	glfw.SetWindowTitle("Wego testing")

	// Hook some events to demonstrate use of callbacks.
	// These are not necessary if you don't need them.
	glfw.SetWindowSizeCallback(onResize)
	glfw.SetWindowCloseCallback(onClose)
	glfw.SetMouseButtonCallback(onMouseBtn)
	glfw.SetMouseWheelCallback(onMouseWheel)
	glfw.SetKeyCallback(onKey)
	glfw.SetCharCallback(onChar)

	initGL(w, h)
	// Start loop
	running := true
	for running {
		// OpenGL rendering goes here.
		resetGL()
		w, h := glfw.WindowSize()
		gc := draw2dgl.NewGraphicContext(w, h)
		gc.SetFillColor(color.RGBA{0x80, 0x80, 0xFF, 0xFF})
		gc.SetStrokeColor(color.RGBA{0x80, 0, 0, 0x80})
		drawer := wego.Drawer{gc}
		drawer.DrawRoundedRect(10, 10, 100, 30, 10)
		// Swap front and back rendering buffers. This also implicitly calls
		// glfw.PollEvents(), so we have valid key/mouse/joystick states after
		// this. This behavior can be disabled by calling glfw.Disable with the
		// argument glfw.AutoPollEvents. You must be sure to manually call
		// PollEvents() or WaitEvents() in this case.
		glfw.SwapBuffers()

		// Break out of loop when Escape key is pressed, or window is closed.
		running = glfw.Key(glfw.KeyEsc) == 0 && glfw.WindowParam(glfw.Opened) == 1
	}
}

func onResize(w, h int) {
	initGL(w, h)
	fmt.Printf("resized: %dx%d\n", w, h)
}

func onClose() int {
	fmt.Println("closed")
	return 1 // return 0 to keep window open.
}

func onMouseBtn(button, state int) {
	fmt.Printf("mouse button: %d, %d\n", button, state)
}

func onMouseWheel(delta int) {
	fmt.Printf("mouse wheel: %d\n", delta)
}

func onKey(key, state int) {
	fmt.Printf("key: %d, %d\n", key, state)
}

func onChar(key, state int) {
	fmt.Printf("char: %d, %d\n", key, state)
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
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA);
	gl.Disable(gl.DEPTH_TEST)
	gl.Hint(gl.LINE_SMOOTH_HINT|gl.LINE_SMOOTH_HINT, gl.NICEST)
}

func resetGL() {
	gl.Clear(gl.COLOR_BUFFER_BIT);

	// Reset the matrix
	gl.MatrixMode(gl.MODELVIEW);
	gl.LoadIdentity();
}
