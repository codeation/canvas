package webapi

import (
	"log"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type font struct {
	id       int
	canvas   js.Value
	cssValue string
	height   int
	r        js.Value
	baseline int
}

func (w *webAPI) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	cssValue := fontValue(height, style, variant, weight, stretch, family)
	canvas := w.document.Call(jsw.CreateElement, jsw.Canvas)
	canvasCtx := canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Set(jsw.Font, cssValue)
	metrics := canvasCtx.Call(jsw.MeasureText, "Gg")
	ascent := int(metrics.Get(jsw.FontBoundingBoxAscent).Float())
	baseline := ascent

	w.fonts[fontID] = &font{
		id:       fontID,
		canvas:   canvas,
		cssValue: cssValue,
		height:   height,
		r:        js.Global().Get(jsw.Range).New(),
		baseline: baseline,
	}
}

func (w *webAPI) FontDrop(fontID int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	font, ok := w.fonts[fontID]
	if !ok {
		log.Printf("FontDrop: font not found: %d", fontID)
		return
	}

	font.canvas.Call(jsw.Remove)

	delete(w.fonts, fontID)
}
