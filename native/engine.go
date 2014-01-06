package native

import (
	"container/list"
	"image"
	"image/color"
	"sort"

	"github.com/go-gl/gl"
)

const maxInt = int(^uint(0) >> 1)

// Color is actually image/color.RGBA
type Color color.RGBA

// Paint holds infomation about the stroke (border) and the colors
type Paint struct {
	FillColor   Color
	StrokeWidth int
	StrokeColor Color
}

// DrawBackend is the one that actually draws things on the window
type DrawBackend interface {
	DrawRect(image.Rectangle, [4]image.Point, Paint)
}

// MakeRect receives coordinates of the min and max points (x1, y1, x2, y2) and returns a Rectangle of type image.Rectangle
func MakeRect(x1, y1, x2, y2 int) image.Rectangle {
	return image.Rectangle{image.Point{x1, y1}, image.Point{x2, y2}}
}

// MakeRectWH receives X, Y, width, height and returns an image.Rectangle
func MakeRectWH(x, y, w, h int) image.Rectangle {
	return MakeRect(x, y, x+w, y+h)
}

// NoStroke returns a Paint with no stroke (and has the specified fillColor)
func NoStroke(fillColor Color) (p Paint) {
	p.FillColor = fillColor
	p.StrokeWidth = 0
	return p
}

// The render interface for renderable things
type shape interface {
	render(IElement, DrawBackend)
}

// RectShape holds information specific to displayed rectangles (rounded).
// Currently it holds border radius of 4 corners
type RectShape struct {
	borderRadiis [4]image.Point
}

// Element holds information about an element in the window
type Element struct {
	parent  *AbstractElement
	treeLev int // The number that represents this element's depth in the tree, the root element has tier=0
	area    image.Rectangle
	zIndex  float32
	Paint
	_alg algData
}

// IElement is the common interface for AbstractElement and ConcreteElement
type IElement interface {
	BaseElement() *Element // Returns the Element subcomponent
	IsConcrete() bool
	AllConcreteDescns() *list.List // Get all concrete descendants
}

// ConcreteElement is the type of Element that can have a shape (apperance)
// and be rendered.
// It cannot have children.
type ConcreteElement struct {
	Element
	shape shape
}

// AbstractElement is for grouping ConcreteElement's.
// It doesn't have a shape and cannot be displayed.
type AbstractElement struct {
	Element
	children []IElement
}

// RectShape method is a helper that casts element's shape to a RectShape and return it
func (e *ConcreteElement) RectShape() *RectShape {
	return e.shape.(*RectShape)
}

// ZIndex gets the z-index of the Element
func (e *Element) ZIndex() float32 {
	return e.zIndex
}

// SetZIndex on an AbstractElement means increasing all its descendants's zindex by that value
func (e *AbstractElement) SetZIndex(z float32) {
	li := e.AllConcreteDescns()
	for o := li.Front(); o != nil; o = o.Next() {
		o.Value.(*ConcreteElement).zIndex += z
	}
}

// SetZIndex on a ConcreteElement
func (e *ConcreteElement) SetZIndex(z float32) {
	e.zIndex = z
}

// IsConcrete on AbstractElement returns false
func (e *AbstractElement) IsConcrete() bool {
	return false
}

// IsConcrete on ConcreteElement returns true
func (e *ConcreteElement) IsConcrete() bool {
	return true
}

func (r *RectShape) render(ei IElement, backend DrawBackend) {
	e := ei.(*ConcreteElement)
	backend.DrawRect(e.area, r.borderRadiis, e.Paint)
}

// X method returns element's top-left x coordinate
func (e *Element) X() int { return e.area.Min.X }

// Y method returns element's top-left y coordinate
func (e *Element) Y() int { return e.area.Min.Y }

// W method returns element's width
func (e *Element) W() int { return e.area.Max.X - e.area.Min.X }

// H method returns element's height
func (e *Element) H() int { return e.area.Max.Y - e.area.Min.Y }

// IsBehind checks whether the element is behind the other element
func (e *ConcreteElement) IsBehind(e2 *ConcreteElement) bool {
	if e.zIndex == e2.zIndex {
		return e.treeLev < e2.treeLev
	}
	return e.zIndex < e2.zIndex
}

// BaseElement returns the Element subcomponent of the element
func (e *AbstractElement) BaseElement() *Element {
	return &e.Element
}

// BaseElement returns the Element subcomponent of the element
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

// AllConcreteDescns fetches and returns a list.List of all the element's descendants
func (e *AbstractElement) AllConcreteDescns() (li *list.List) {
	li = list.New()
	li = e.fetchConcreteDescns(li)
	return li
}

// AllConcreteDescns on a ConcreteElement just returns itself in a list
func (e *ConcreteElement) AllConcreteDescns() (li *list.List) {
	li = list.New()
	li.PushBack(e)
	return li
}

// Data for algorithm purposes
type algData struct {
	addedToRedraw bool
}

// OverlappedAlgorithm encapsulates the algorithm to get overlapping elements, used for lazy redrawing
type OverlappedAlgorithm struct{}

// InitOverlappedAlgorithm initializes the algorithm.
// All changes by the algorithm are reset.
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

// NewRectElement creates a new rsectangle concrete element
func NewRectElement(parent *AbstractElement, area image.Rectangle) *ConcreteElement {
	e := new(ConcreteElement)
	parent.AddChild(e)
	e.shape = new(RectShape)
	e.area = area
	return e
}

// SetAll4RadiiXY sets all 4 corners of the RectShape to have the radius x, y
func (r *RectShape) SetAll4RadiiXY(x, y int) {
	for i := 0; i < 4; i++ {
		r.borderRadiis[i] = image.Point{x, y}
	}
}

// SetAllRadii sets all 4 corners of the RectShape to have the same radius rad
func (r *RectShape) SetAllRadii(rad int) {
	r.SetAll4RadiiXY(rad, rad)
}

// NewRootElement creates and returns the root element
func NewRootElement() (r *AbstractElement) {
	r = new(AbstractElement)
	r.treeLev = 0
	r.zIndex = 0
	r.area = MakeRect(0, 0, maxInt, maxInt)
	return r
}

// NewAbstractElement creates and returns a new AbstractElement.
// It adds the element as a child.
func NewAbstractElement(parent *AbstractElement) (r *AbstractElement) {
	r = new(AbstractElement)
	parent.AddChild(r)
	return r
}

// AddChild makes an element the child of an AbstractElement
func (e *AbstractElement) AddChild(child IElement) {
	e.children = append(e.children, child)
	c := child.BaseElement()
	c.parent = e
	c.treeLev = e.treeLev + 1
	c.zIndex = 0
}

// Draw the element and all its descendants
func (e *AbstractElement) Draw(backend DrawBackend) {
	l := makeDrawPriorityList(e.AllConcreteDescns())
	sort.Sort(l)
	for _, o := range l {
		o.Draw(backend)
	}
}

// Draw the element
func (e *ConcreteElement) Draw(backend DrawBackend) {
	e.shape.render(e, backend)
}

type drawPriorityList [](*ConcreteElement)

func makeDrawPriorityList(li *list.List) drawPriorityList {
	l := make([](*ConcreteElement), li.Len())
	for o, i := li.Front(), 0; o != nil; o, i = o.Next(), i+1 {
		l[i] = o.Value.(*ConcreteElement)
	}
	return drawPriorityList(l)
}

func (l drawPriorityList) Len() int      { return len(l) }
func (l drawPriorityList) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l drawPriorityList) Less(i, j int) bool {
	return l[i].IsBehind(l[j])
}

// Redraw the element
func Redraw(e IElement, backend DrawBackend, root *AbstractElement) {
	alg := InitOverlappedAlgorithm(root)
	itemsToRedraw := list.New()
	d := e.AllConcreteDescns()
	for o := d.Front(); o != nil; o = o.Next() {
		alg.addToRedrawList(o.Value.(*ConcreteElement), itemsToRedraw)
	}
	alg.fetchOverlappingConcreteElems(e, root, root, itemsToRedraw)
	l := makeDrawPriorityList(itemsToRedraw)
	sort.Sort(l)
	gl.Enable(gl.SCISSOR_TEST)
	for _, o := range l {
		gl.Scissor(o.X(), o.Y(), o.W(), o.H())
		o.Draw(backend)
	}
	gl.Disable(gl.SCISSOR_TEST)
}
