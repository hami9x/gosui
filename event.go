package gosui

import "image"

const (
	MouseButtonLeft   = 0
	MouseButtonRight  = 1
	MouseButtonMiddle = 2
	EventPress        = 0
	EventRelease      = 1
	EventRepeat       = 2
)

type Modifiers struct {
	Control, Shift, Alt, Super bool
}

type MouseButton int

type EventAction int

type mouseHandler interface {
	OnMouseEvent(*MouseEvent) bool
}

type MouseEvent struct {
	Pos    image.Point
	Button MouseButton
	Mod    Modifiers
	Action EventAction
}

func HandleMouse(evt *MouseEvent, e *AbstractElement) {
	for _, child := range e.children {
		if !evt.Pos.In(child.BaseElement().Area) {
			continue
		}
		propagate := true
		if handler, ok := child.BaseElement().Handler.(mouseHandler); ok {
			propagate = handler.OnMouseEvent(evt)
		}
		if !child.IsConcrete() && propagate {
			HandleMouse(evt, child.(*AbstractElement))
		}
	}
}
