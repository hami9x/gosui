package gosui

import (
	"container/list"
	// "fmt"
	"image"
	"image/color"
	"sort"
)

const MaxInt = int(^uint(0) >> 1)

type Color color.RGBA

type Paint struct {
	FillColor   Color
	StrokeWidth int
	StrokeColor Color
}

type BackendPtr interface {
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

type Appearance interface {
	render(*Element, BackendPtr)
}

type NoAppearance struct{}

func (*NoAppearance) render(e *Element, b BackendPtr) {}

type Element struct {
	parent   *Element
	children [](*Element)
	treeLev  int //The number that represents this element's depth in the tree, the root element has tier=0
	area     image.Rectangle
	appr     Appearance
	Paint
	ZIndex float32
	_alg   AlgData
}

func (e *Element) IsBehind(e2 *Element) bool {
	if e.ZIndex == e2.ZIndex {
		return e.treeLev <= e2.treeLev
	}
	return e.ZIndex < e2.ZIndex
}

func (e *Element) fetchDescendants(li *list.List) *list.List {
	for _, o := range e.children {
		li.PushBack(o)
		li = o.fetchDescendants(li)
	}
	return li
}

func (e *Element) AllDescendants() (li *list.List) {
	li = list.New()
	li = e.fetchDescendants(li)
	return li
}

func (e *Element) IsLeaf() bool {
	if len(e.children) == 0 {
		return true
	}
	return false
}

//Data for algorithm purposes
type AlgData struct {
	addedToRedraw bool
}

//Algorithm to get overlapping elements, used for lazy redrawing
type OverlappedAlgorithm struct{}

func InitOverlappedAlgorithm(root *Element) (alg OverlappedAlgorithm) {
	l := root.AllDescendants()
	for o := l.Front(); o != nil; o = o.Next() {
		o.Value.(*Element)._alg.addedToRedraw = false
	}
	return alg
}

func (alg OverlappedAlgorithm) addToRedrawList(e *Element, li *list.List) {
	e._alg.addedToRedraw = true
	li.PushBack(e)
}

func (alg OverlappedAlgorithm) hasAdded(e *Element) bool {
	return e._alg.addedToRedraw
}

func (alg OverlappedAlgorithm) fetchOverlappingLeafElems(target *Element, e *Element, root *Element, li *list.List) *list.List {
	for _, o := range e.children {
		if target.area.Overlaps(o.area) && (o != target) && (!alg.hasAdded(o)) {
			// fmt.Printf("%v :: %v\n", target.area, o.area)
			if o.IsLeaf() && target.IsBehind(o) {
				alg.addToRedrawList(o, li)
				li = alg.fetchOverlappingLeafElems(o, root, root, li)
			} else {
				li = alg.fetchOverlappingLeafElems(target, o, root, li)
			}
		}
	}
	return li
}

type RectShape struct {
	borderRadiis [4]image.Point
}

func NewRectElement(parent *Element, area image.Rectangle) *Element {
	e := new(Element).AddTo(parent)
	e.appr = new(RectShape)
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

func (r *RectShape) render(e *Element, backend BackendPtr) {
	backend.DrawRect(e.area, r.borderRadiis, e.Paint)
}

func NewRootElement() (r *Element) {
	r = new(Element)
	r.treeLev = 0
	r.ZIndex = 0
	r.appr = &NoAppearance{}
	r.area = MakeRect(0, 0, MaxInt, MaxInt)
	return r
}

func (e *Element) AddTo(parent *Element) *Element {
	parent.children = append(parent.children, e)
	e.appr = &NoAppearance{}
	e.parent = parent
	e.treeLev = parent.treeLev + 1
	e.ZIndex = 0
	return e
}

//Draw the element and all its descendants
func (e *Element) Draw(backend BackendPtr) {
	li := e.AllDescendants()
	li.PushFront(e)
	l := MakeDrawPriorityList(li)
	sort.Sort(l)
	for _, o := range l {
		o.appr.render(o, backend)
	}
}

func rectNoRounding() (r [4]image.Point) {
	return r
}

func (e *Element) Clear(backend BackendPtr) {
	backend.DrawRect(e.area, rectNoRounding(), Paint{ColorWHITE(), 0, ColorWHITE()})
}

type DrawPriorityList [](*Element)

func MakeDrawPriorityList(li *list.List) DrawPriorityList {
	l := make([](*Element), li.Len())
	for o, i := li.Front(), 0; o != nil; o, i = o.Next(), i+1 {
		l[i] = o.Value.(*Element)
	}
	return DrawPriorityList(l)
}

func (l DrawPriorityList) Len() int      { return len(l) }
func (l DrawPriorityList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l DrawPriorityList) Less(i, j int) bool {
	return l[i].IsBehind(l[j])
}

func (e *Element) Redraw(backend BackendPtr, root *Element) {
	alg := InitOverlappedAlgorithm(root)
	itemsToRedraw := list.New()
	alg.addToRedrawList(e, itemsToRedraw)
	alg.fetchOverlappingLeafElems(e, root, root, itemsToRedraw)
	l := MakeDrawPriorityList(itemsToRedraw)
	sort.Sort(l)
	for _, o := range l {
		o.Clear(backend)
		o.Draw(backend)
	}
}
