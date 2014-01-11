package gosui

import (
	"container/list"
	"image"
	"image/color"
	"sort"
)

const (
	maxInt = int(^uint(0) >> 1)
)

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
	DrawRect(image.Rectangle, [4]int, Paint)
	DrawText(image.Point, *TextShape, Paint) (int, int)
}

type RenderBackend interface {
	DrawBackend
	Init(int, int)
	DrawElementsInArea(DrawPriorityList, image.Rectangle)
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
// Currently it holds corner radius of 4 corners
type RectShape struct {
	cornerRadiis [4]int
}

// Element holds information about an element in the window
type Element struct {
	parent  *AbstractElement
	treeLev int // The number that represents this element's depth in the tree, the root element has tier=0
	Area    image.Rectangle
	zIndex  float32
	Paint
	cData map[string]interface{}
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

type FontStyle struct {
	Bold, Italic bool
}

var (
	BoldItalic = FontStyle{true, true}
	Bold       = FontStyle{true, false}
	Italic     = FontStyle{false, true}
	Regular    = FontStyle{false, false}
)

type Font struct {
	Family string
	Size   int
	Style  FontStyle
}

type TextShape struct {
	Content  string
	Font     Font
	Editable bool
}

func (e *ConcreteElement) TextShape() *TextShape {
	return e.shape.(*TextShape)
}

func (e *ConcreteElement) UpdateSize(w, h int) {
	e.Area.Max = image.Point{e.Area.Min.X + w, e.Area.Min.Y + h}
}

// SetData is used by something to add data to the element
// for its own purpose, like for an algorithm
func (e *Element) SetData(key string, data interface{}) {
	if e.cData == nil {
		e.cData = make(map[string]interface{})
	}
	e.cData[key] = data
}

func (e *Element) GetData(key string) interface{} {
	return e.cData[key]
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
	backend.DrawRect(e.Area, r.cornerRadiis, e.Paint)
}

func (s *TextShape) render(ei IElement, backend DrawBackend) {
	e := ei.(*ConcreteElement)
	e.UpdateSize(backend.DrawText(e.Area.Min, s, e.Paint))
}

// X method returns element's top-left x coordinate
func (e *Element) X() int { return e.Area.Min.X }

// Y method returns element's top-left y coordinate
func (e *Element) Y() int { return e.Area.Min.Y }

// W method returns element's width
func (e *Element) W() int { return e.Area.Dx() }

// H method returns element's height
func (e *Element) H() int { return e.Area.Dy() }

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
		e := o.Value.(*ConcreteElement)
		e.BaseElement().SetData("addedToRedraw", false)
	}
	return alg
}

func (alg OverlappedAlgorithm) addToRedrawList(e *ConcreteElement, li *list.List) {
	e.BaseElement().SetData("addedToRedraw", true)
	li.PushBack(e)
}

func (alg OverlappedAlgorithm) hasAdded(e *Element) bool {
	return e.GetData("addedToRedraw").(bool)
}

func (alg OverlappedAlgorithm) fetchOverlappingConcreteElems(target IElement, e *AbstractElement, root *AbstractElement, li *list.List) *list.List {
	for _, oi := range e.children {
		o := oi.BaseElement()
		t := target.BaseElement()
		if t.Area.Overlaps(o.Area) && (o != t) && (!alg.hasAdded(o)) {
			if oi.IsConcrete() {
				alg.addToRedrawList(oi.(*ConcreteElement), li)
			} else {
				li = alg.fetchOverlappingConcreteElems(target, oi.(*AbstractElement), root, li)
			}
		}
	}
	return li
}

func NewTextInputElement(parent *AbstractElement, x, y int, font Font) *ConcreteElement {
	return NewTextElement(parent, x, y, font, true)
}

func NewTextElement(parent *AbstractElement, x, y int, font Font, editable bool) *ConcreteElement {
	e := new(ConcreteElement)
	parent.AddChild(e)
	ts := new(TextShape)
	ts.Editable = editable
	ts.Font = font
	e.shape = ts
	e.Area = image.Rectangle{image.Point{x, y}, image.Point{x, y}}
	return e
}

// NewRectElement creates a new rsectangle concrete element
func NewRectElement(parent *AbstractElement, area image.Rectangle) *ConcreteElement {
	e := new(ConcreteElement)
	parent.AddChild(e)
	e.shape = new(RectShape)
	e.Area = area
	return e
}

// SetAllCornerRadiusTo sets all 4 corners of the RectShape to have same radius rad
func (r *RectShape) SetAllCornerRadiusTo(rad int) {
	for i := 0; i < 4; i++ {
		r.cornerRadiis[i] = rad
	}
}

type RectCornersRad struct {
	TopLeft, TopRight, BotLeft, BotRight int
}

func (r *RectShape) SetCornerRadiis(conf RectCornersRad) {
	r.cornerRadiis[0], r.cornerRadiis[1] = conf.TopLeft, conf.TopRight
	r.cornerRadiis[2], r.cornerRadiis[3] = conf.BotLeft, conf.BotRight
}

// NewRootElement creates and returns the root element
func NewRootElement() (r *AbstractElement) {
	r = new(AbstractElement)
	r.treeLev = 0
	r.zIndex = 0
	r.Area = MakeRect(0, 0, maxInt, maxInt)
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

type DrawPriorityList [](*ConcreteElement)

func makeDrawPriorityList(li *list.List) DrawPriorityList {
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

// Redraw the element
func Redraw(e IElement, backend RenderBackend, root *AbstractElement) {
	alg := InitOverlappedAlgorithm(root)
	itemsToRedraw := list.New()
	d := e.AllConcreteDescns()
	for o := d.Front(); o != nil; o = o.Next() {
		alg.addToRedrawList(o.Value.(*ConcreteElement), itemsToRedraw)
	}
	alg.fetchOverlappingConcreteElems(e, root, root, itemsToRedraw)
	l := makeDrawPriorityList(itemsToRedraw)
	sort.Sort(l)

	backend.DrawElementsInArea(l, e.BaseElement().Area)
}
