//Test app for gosui
package main

import gs "github.com/phaikawl/gosui"

//import gsr "github.com/phaikawl/gosui/native/skia"

import gse "github.com/phaikawl/gosui/web"

//import gse "github.com/phaikawl/gosui/native"

import gsr "github.com/phaikawl/gosui/web/htmlcanvas"

func main() {
	window := gse.NewWindow(new(gsr.Backend), 800, 600, "Gosui test app")
	root := window.RootElement()
	w, h := window.Size()
	bg := gs.NewRectElement(root, gs.MakeRectWH(0, 0, w, h))
	bg.FillColor = gs.Color{255, 255, 255, 255}
	bg.SetZIndex(-1000000)
	rect := gs.NewRectElement(root, gs.MakeRectWH(10, 10, 100, 100))
	rect.FillColor = gs.Color{255, 0, 0, 255}
	rect.RectShape().SetAllCornerRadiusTo(10)
	rect2 := gs.NewRectElement(root, gs.MakeRectWH(50, 0, 120, 400))
	rect2.FillColor = gs.Color{120, 0, 30, 100}
	rect2.StrokeColor = gs.Color{0, 155, 20, 255}
	rect2.RectShape().SetCornerRadiis(gs.RectCornersRad{10, 30, 5, 0})
	window.Start()
}
