package skia

// #cgo LDFLAGS: -lskia_go_renderer
// #include "gosui/skia.h"
import "C"
import (
	"image/color"
)

type Color color.RGBA

func (c Color) toSkColor() C.Color {
	c.A *= 255
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

func (b *Backend) DrawButton(x, y, w, h, radii int) {
	var rect C.Rect
	rect.x, rect.y, rect.w, rect.h = C.int(x), C.int(y), C.int(w), C.int(h)
	var paint C.Paint
	paint.fillColor = Color{0, 0, 255, 1}.toSkColor()
	paint.strokeColor = Color{255, 0, 0, 1}.toSkColor()
	C.DrawRect(b.r, paint, rect, C.int(radii))
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

func (b *Backend) StartLoop() {
	b.saveCnt = C.Save(b.r)
	b.Clear()
}

func (b *Backend) EndLoop() {
	C.Restore(b.r, b.saveCnt)
	b.Flush()
}
