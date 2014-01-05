package main

import gs "github.com/phaikawl/gosui/native"
import skia "github.com/phaikawl/gosui/native/skia"

func main() {
	window := gs.NewWindow(new(skia.Backend), 800, 600, "Gosui test app")
	root := window.RootElement()
	w, h := window.Size()
	bg := gs.NewRectElement(root, gs.MakeRectWH(0, 0, w, h))
	bg.FillColor = gs.Color{255, 255, 255, 255}
	bg.SetZIndex(-1000000)
	rect := gs.NewRectElement(root, gs.MakeRectWH(10, 10, 100, 100))
	rect.FillColor = gs.Color{255, 0, 0, 255}
	rect.RectShape().SetAllRadii(15)

	window.Loop()
}
