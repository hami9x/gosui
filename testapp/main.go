//Test app for gosui
package main

import (
	gs "github.com/phaikawl/gosui"

	gsr "github.com/phaikawl/gosui/native/skia" //Native rendering backend

	//gse "github.com/phaikawl/gosui/web" //Web engine

	gse "github.com/phaikawl/gosui/native" //Native engine

	//gsr "github.com/phaikawl/gosui/web/htmlcanvas" //Web rendering backend
)

func main() {
	gse.AddAssetDir(gse.LocalDir("dist"))
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
	input := gs.NewTextInputElement(root, 30, 70, gs.Font{"Arial", 18, gs.BoldItalic})
	input.FillColor = gs.Color{0, 0, 40, 255}
	input.SetZIndex(1000)
	input.TextShape().Content = "Hello world!"
	window.Start()
}
