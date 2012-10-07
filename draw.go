package wego

import (
	"math"
	"code.google.com/p/draw2d/draw2dgl"
)

type Drawer struct {
	Gc *draw2dgl.GraphicContext
}

const (
	degree = math.Pi / 180.0
)

func (d *Drawer) DrawRoundedRect(x, y, width, height, radius float64) {
	d.Gc.ArcTo(x+width-radius, y+radius, radius, radius, -90*degree, 90*degree)
	d.Gc.ArcTo(x+width-radius, y+height-radius, radius, radius, 0*degree, 90*degree)
	d.Gc.ArcTo(x+radius, y+height-radius, radius, radius, 90*degree, 90*degree)
	d.Gc.ArcTo(x+radius, y+radius, radius, radius, 180*degree, 90*degree)
	d.Gc.Close()
	d.Gc.FillStroke()
}
