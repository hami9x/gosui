package gosui

import (
	draw2d "code.google.com/p/draw2d/draw2dgl"
	"github.com/go-gl/gl"
	"image"
	"image/color"
	"math"
)

type Drawer struct {
	Gc    *draw2d.GraphicContext
	Image *image.RGBA
}

const (
	degree = math.Pi / 180.0
)

func MakeDrawer(gc *draw2d.GraphicContext, image *image.RGBA) *Drawer {
	drawer := &Drawer{gc, image}
	drawer.Init()
	return drawer
}

func (d *Drawer) Init() {
	// draw2d.SetFontFolder("font/")
}

func (d *Drawer) DrawRoundedRect(x, y, width, height, radius float64) {
	gc := d.Gc
	gc.ArcTo(x+width-radius, y+radius, radius, radius, -90*degree, 90*degree)
	gc.ArcTo(x+width-radius, y+height-radius, radius, radius, 0*degree, 90*degree)
	gc.ArcTo(x+radius, y+height-radius, radius, radius, 90*degree, 90*degree)
	gc.ArcTo(x+radius, y+radius, radius, radius, 180*degree, 90*degree)
	gc.Close()
	gc.FillStroke()
}

// func (d *Drawer) DrawText(str string, x, y float64, fontSize float64, fontData draw2d.FontData) {
// 	gc := d.Gc
// 	gc.FillStroke()
// 	gc.SetFontSize(fontSize)
// 	gc.SetFontData(fontData)
// 	gc.FillStringAt(str, x, y)
// }

//Render OpenGL vertices from go standard Image
func (d *Drawer) RenderToGL() {
	img := d.Image
	b := img.Bounds()
	gl.Begin(gl.POINTS)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			col := img.At(x, y)
			r, g, b, a := GLColorFromRGBA(col.(color.RGBA))
			gl.Color4f(r, g, b, a)
			gl.Vertex2i(x, y)
		}
	}
	gl.End()
}

//Convert RGBA color to an OpenGL-compatible value
func GLColorFromRGBA(col color.RGBA) (float32, float32, float32, float32) {
	r, g, b, _a := col.RGBA()
	if _a == 0 {
		return 0, 0, 0, 0
	}
	to1 := func(c uint32) float32 {
		return float32(c) / 65535
	}
	a := to1(_a)
	return to1(r) / a, to1(g) / a, to1(b) / a, a
}
