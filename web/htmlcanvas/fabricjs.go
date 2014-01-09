package htmlcanvas

import (
	"encoding/json"
	"fmt"
	"image"

	gs "github.com/phaikawl/gosui"
)

type Backend struct{}

type FabricObject struct {
	Left        int    `json:"left"`
	Top         int    `json:"top"`
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Fill        string `json:"fill"`
	Stroke      string `json:"stroke"`
	StrokeWidth int    `json:"strokeWidth"`
	Angle       int    `json:"angle"`
}

type FabricRectObj struct {
	FabricObject
	CornerRadiis [4]int `json:"cornerRadiis"`
}

func toHtmlColor(color gs.Color) string {
	return fmt.Sprintf("rgba(%d, %d, %d, %v)", color.R, color.G, color.B, float32(color.A)/255)
}

func makeFabricObject(area image.Rectangle, paint gs.Paint) FabricObject {
	return FabricObject{
		Left:        area.Min.X,
		Top:         area.Min.Y,
		Width:       area.Dx(),
		Height:      area.Dy(),
		Fill:        toHtmlColor(paint.FillColor),
		Stroke:      toHtmlColor(paint.StrokeColor),
		StrokeWidth: paint.StrokeWidth,
	}
}

func iDrawRect(spec string) {}

const js_iDrawRect = `fabricDrawRect(JSON.parse(spec));`

func jsLog(msg string) {}

const js_jsLog = `console.log(msg)`

func (b *Backend) DrawRect(rect image.Rectangle, radiis [4]int, paint gs.Paint) {
	jsObj, err := json.Marshal(FabricRectObj{
		FabricObject: makeFabricObject(rect, paint),
		CornerRadiis: radiis,
	})
	if err != nil {
		panic(err.Error())
	}
	iDrawRect(string(jsObj[:]))
}

func jsInit(w, h int) {}

const js_jsInit = `fabricCanvasResize(w, h)`

func (b *Backend) Init(w, h int) {
	jsInit(w, h)
}

func (b *Backend) DrawElementsInArea(l gs.DrawPriorityList, area image.Rectangle) {
}
