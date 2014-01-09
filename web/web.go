package web

import (
	"image"

	gs "github.com/phaikawl/gosui"
)

//Window is the application window
type Window struct {
	b gs.RenderBackend

	area image.Rectangle
	root *gs.AbstractElement
}

func (wn *Window) Size() (w, h int) {
	return wn.area.Dx(), wn.area.Dy()
}

func NewWindow(b gs.RenderBackend, w, h int, title string) *Window {
	b.Init(w, h)
	return &Window{
		b:    b,
		area: gs.MakeRectWH(0, 0, w, h),
		root: gs.NewRootElement(),
	}
}

func (wn *Window) RootElement() *gs.AbstractElement {
	return wn.root
}

func (wn *Window) Start() {
	wn.root.Draw(wn.b)
}
