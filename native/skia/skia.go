package skia

// #cgo LDFLAGS: -lskia_go_renderer
// #include "gosui/skia.h"
import "C"
import (
	// "fmt"
	"image"

	"github.com/go-gl/gl"
	gs "github.com/phaikawl/gosui"
)

func toSkColor(c gs.Color) C.Color {
	return (C.Color)(C.ColorFromRGBA(C.int(c.R), C.int(c.G), C.int(c.B), C.int(c.A)))
}

type Backend struct {
	w, h    int
	r       C.SkiaRenderer
	saveCnt C.int
}

func (b *Backend) Init(w, h int) {
	b.w = w
	b.h = h
	r := C.Init(C.int(w), C.int(h))
	b.r = r
}

func toCRect(rect image.Rectangle) (crect C.Rect) {
	crect.min, crect.max = toCPoint(rect.Min), toCPoint(rect.Max)
	return crect
}

func toCPoint(sp image.Point) (p C.Point) {
	p.x = C.int(sp.X)
	p.y = C.int(sp.Y)
	return p
}

func (b *Backend) DrawRect(rect image.Rectangle, radiis [4]int, paint gs.Paint) {
	crect := toCRect(rect)
	// fmt.Printf("%v : %v\n", crect.min.x, crect.min.y)

	var cpaint C.Paint
	cpaint.fillColor = toSkColor(paint.FillColor)
	cpaint.strokeColor = toSkColor(paint.StrokeColor)
	cpaint.strokeWidth = C.int(paint.StrokeWidth)
	var cRads [4]C.Point
	for i := 0; i < 4; i += 1 {
		cRads[i] = toCPoint(image.Point{radiis[i], radiis[i]})
	}
	C.DrawRect(b.r, cpaint, crect, (*C.Point)(&cRads[0]))
}

func (b *Backend) Die() {
	C.Die(b.r)
}

func (b *Backend) Flush() {
	C.Flush(b.r)
}

func (b *Backend) Clear() {
	C.Clear(b.r)
}

func (b *Backend) Save() {
	b.saveCnt = C.Save(b.r)
}

func (b *Backend) Restore() {
	C.Restore(b.r, b.saveCnt)
}

func (b *Backend) ClipRect(rect image.Rectangle) {
	C.ClipRect(b.r, toCRect(rect))
}

func (b *Backend) UpdateViewportSize(w, h int) {
	C.UpdateWindowSize(b.r, C.int(w), C.int(h))
}

//DrawElementsInArea is used for redrawing
func (b *Backend) DrawElementsInArea(l gs.DrawPriorityList, area image.Rectangle) {
	gl.Enable(gl.SCISSOR_TEST)
	gl.Scissor(area.Min.X, area.Min.Y, area.Dx(), area.Dy())
	for _, o := range l {
		o.Draw(b)
	}
	gl.Disable(gl.SCISSOR_TEST)
}
