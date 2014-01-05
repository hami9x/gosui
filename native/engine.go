package native

import (
	"container/list"
	"image"
	"image/color"
	"sort"

	"github.com/go-gl/gl"
)

const MaxInt = int(^uint(0) >> 1)

type Color color.RGBA

type Paint struct {
	FillColor   Color
	StrokeWidth int
	StrokeColor Color
}

type DrawBackend interface {
	DrawRect(image.Rectangle, [4]image.Point, Paint)
}

func ColorWHITE() Color { return Color{255, 255, 255, 255} }

func MakeRect(x1, y1, x2, y2 int) image.Rectangle {
	return image.Rectangle{image.Point{x1, y1}, image.Point{x2, y2}}
}

func MakeRectWH(x, y, w, h int) image.Rectangle {
	return MakeRect(x, y, x+w, y+h)
}

func NoStroke(fillColor Color) (p Paint) {
	p.FillColor = fillColor
	p.StrokeWidth = 0
	return p
}

type Shape interface {
	render(IElement, DrawBackend)
}

type NoShape struct{}

func (*NoShape) render(e *Element, b DrawBackend) {}

type Element struct {
	parent  *AbstractElement
	treeLev int //The number that represents this element's depth in the tree, the root element has tier=0
	area    image.Rectangle
	zIndex  float32
	Paint
	_alg AlgData
}

type IElement interface {
	BaseElement() *Element
	IsConcrete() bool
	AllConcreteDescns() *list.List //Get all concrete descendants
}

type ConcreteElement struct {
	Element
	shape Shape
}

type AbstractElement struct {
	Element
	children []IElement
}

func (e *Element) ZIndex() float32 {
	return e.zIndex
}

func (e *AbstractElement) setZIndex(z float32) {
	li := e.AllConcreteDescns()
	for o := li.Front(); o != nil; o = o.Next() {
		o.Value.(*ConcreteElement).zIndex += z
	}
}

func (e *ConcreteElement) SetZIndex(z float32) {
	e.zIndex = z
}

func (e *AbstractElement) IsConcrete() bool {
	return false
}

func (e *ConcreteElement) IsConcrete() bool {
	return true
}

func (e *ConcreteElement) RectShape() *RectShape {
	return e.shape.(*RectShape)
}

type RectShape struct {
	borderRadiis [4]image.Point
}

func (r *RectShape) render(ei IElement, backend DrawBackend) {
	e := ei.(*ConcreteElement)
	backend.DrawRect(e.area, r.borderRadiis, e.Paint)
}

func (e *Element) X() int { return e.area.Min.X }

func (e *Element) Y() int { return e.area.Min.Y }

func (e *Element) W() int { return e.area.Max.X - e.area.Min.X }

func (e *Element) H() int { return e.area.Max.Y - e.area.Min.Y }

func (e *ConcreteElement) IsBehind(e2 *ConcreteElement) bool {
	if e.zIndex == e2.zIndex {
		return e.treeLev < e2.treeLev
	}
	return e.zIndex < e2.zIndex
}

func (e *AbstractElement) BaseElement() *Element {
	return &e.Element
}

func (e *ConcreteElement) BaseElement() *Element {
	return &e.Element
}

func (e *AbstractElement) fetchConcreteDescns(li *list.List) *list.List {
	for _, oi := range e.children {
		if o, isConcrete := oi.(*ConcreteElement); isConcrete {
			li.PushBack(o)
			continue
		}
		o := oi.(*AbstractElement)
		o.fetchConcreteDescns(li)
	}
	return li
}

func (e *AbstractElement) AllConcreteDescns() (li *list.List) {
	li = list.New()
	li = e.fetchConcreteDescns(li)
	return li
}

func (e *ConcreteElement) AllConcreteDescns() (li *list.List) {
	li = list.New()
	li.PushBack(e)
	return li
}

//Data for algorithm purposes
type AlgData struct {
	addedToRedraw bool
}

//Algorithm to get overlapping elements, used for lazy redrawing
type OverlappedAlgorithm struct{}

func InitOverlappedAlgorithm(root *AbstractElement) (alg OverlappedAlgorithm) {
	l := root.AllConcreteDescns()
	for o := l.Front(); o != nil; o = o.Next() {
		o.Value.(*ConcreteElement)._alg.addedToRedraw = false
	}
	return alg
}

func (alg OverlappedAlgorithm) addToRedrawList(e *ConcreteElement, li *list.List) {
	e._alg.addedToRedraw = true
	li.PushBack(e)
}

func (alg OverlappedAlgorithm) hasAdded(e *Element) bool {
	return e._alg.addedToRedraw
}

func (alg OverlappedAlgorithm) fetchOverlappingConcreteElems(target IElement, e *AbstractElement, root *AbstractElement, li *list.List) *list.List {
	for _, oi := range e.children {
		o := oi.BaseElement()
		t := target.BaseElement()
		if t.area.Overlaps(o.area) && (o != t) && (!alg.hasAdded(o)) {
			// fmt.Printf("%v :: %v\n", target.area, o.area)
			if oi.IsConcrete() {
				alg.addToRedrawList(oi.(*ConcreteElement), li)
			} else {
				li = alg.fetchOverlappingConcreteElems(target, oi.(*AbstractElement), root, li)
			}
		}
	}
	return li
}

func NewRectElement(parent *AbstractElement, area image.Rectangle) *ConcreteElement {
	e := new(ConcreteElement)
	parent.AddChild(e)
	e.shape = new(RectShape)
	e.area = area
	return e
}

func (r *RectShape) SetAll4RadiiXY(x, y int) {
	for i := 0; i < 4; i += 1 {
		r.borderRadiis[i] = image.Point{x, y}
	}
}

func (r *RectShape) SetAllRadii(rad int) {
	r.SetAll4RadiiXY(rad, rad)
}

func NewRootElement() (r *AbstractElement) {
	r = new(AbstractElement)
	r.treeLev = 0
	r.zIndex = 0
	r.area = MakeRect(0, 0, MaxInt, MaxInt)
	return r
}

func NewAbstractElement(parent *AbstractElement) (r *AbstractElement) {
	r = new(AbstractElement)
	parent.AddChild(r)
	return r
}

func (parent *AbstractElement) AddChild(ei IElement) {
	parent.children = append(parent.children, ei)
	e := ei.BaseElement()
	e.parent = parent
	e.treeLev = parent.treeLev + 1
	e.zIndex = 0
}

//Draw the element and all its descendants
func (e *AbstractElement) Draw(backend DrawBackend) {
	l := MakeDrawPriorityList(e.AllConcreteDescns())
	sort.Sort(l)
	for _, o := range l {
		o.Draw(backend)
	}
}

func (e *ConcreteElement) Draw(backend DrawBackend) {
	e.shape.render(e, backend)
}

func rectNoRounding() (r [4]image.Point) {
	return r
}

func (e *Element) Clear(backend DrawBackend) {
	backend.DrawRect(e.area, rectNoRounding(), Paint{ColorWHITE(), 0, ColorWHITE()})
}

type DrawPriorityList [](*ConcreteElement)

func MakeDrawPriorityList(li *list.List) DrawPriorityList {
	l := make([](*ConcreteElement), li.Len())
	for o, i := li.Front(), 0; o != nil; o, i = o.Next(), i+1 {
		l[i] = o.Value.(*ConcreteElement)
	}
	return DrawPriorityList(l)
}

func (l DrawPriorityList) Len() int      { return len(l) }
func (l DrawPriorityList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l DrawPriorityList) Less(i, j int) bool {
	return l[i].IsBehind(l[j])
}

func Redraw(e IElement, backend DrawBackend, root *AbstractElement) {
	alg := InitOverlappedAlgorithm(root)
	itemsToRedraw := list.New()
	d := e.AllConcreteDescns()
	for o := d.Front(); o != nil; o = o.Next() {
		alg.addToRedrawList(o.Value.(*ConcreteElement), itemsToRedraw)
	}
	alg.fetchOverlappingConcreteElems(e, root, root, itemsToRedraw)
	l := MakeDrawPriorityList(itemsToRedraw)
	sort.Sort(l)
	gl.Enable(gl.SCISSOR_TEST)
	for _, o := range l {
		gl.Scissor(o.X(), o.Y(), o.W(), o.H())
		o.Draw(backend)
	}
	gl.Disable(gl.SCISSOR_TEST)
}
