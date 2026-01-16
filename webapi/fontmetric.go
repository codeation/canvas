package webapi

import (
	"log"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

func (w *webAPI) FontMetricNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int, int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	cssValue := fontValue(height, style, variant, weight, stretch, family)
	canvas := w.document.Call(jsw.CreateElement, jsw.Canvas)
	canvasCtx := canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Set(jsw.Font, cssValue)
	metrics := canvasCtx.Call(jsw.MeasureText, "Gg")
	ascent := int(metrics.Get(jsw.FontBoundingBoxAscent).Float())
	descent := int(metrics.Get(jsw.FontBoundingBoxDescent).Float())
	lineheight := height
	baseline := ascent

	w.metricFonts[fontID] = &font{
		id:       fontID,
		canvas:   canvas,
		cssValue: cssValue,
		height:   height,
		r:        js.Global().Get(jsw.Range).New(),
		baseline: baseline,
	}

	return lineheight, baseline, ascent, descent
}

func (w *webAPI) FontMetricDrop(fontID int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	font, ok := w.metricFonts[fontID]
	if !ok {
		log.Printf("FontMetricDrop: font not found: %d", fontID)
		return
	}

	font.canvas.Call(jsw.Remove)

	delete(w.metricFonts, fontID)
}

func (w *webAPI) FontMetricSplit(fontID int, text string, edge, indent int) []int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	f, ok := w.metricFonts[fontID]
	if !ok {
		log.Printf("FontMetricSplit: font not found: %d", fontID)
		return nil
	}

	p := w.document.Call(jsw.CreateElement, jsw.P)
	p.Set(jsw.TextContent, text)
	pStyle := p.Get(jsw.Style)
	pStyle.Set(jsw.Visibility, jsw.Hidden)
	pStyle.Set(jsw.Font, f.cssValue)
	pStyle.Set(jsw.Margin, 0)
	pStyle.Set(jsw.Padding, 0)
	pStyle.Set(jsw.Width, px(edge))
	pStyle.Set(jsw.TextIndent, px(indent))
	pStyle.Set(jsw.LineHeight, 1)
	w.document.Get(jsw.Body).Call(jsw.AppendChild, p)

	output := textSplit(p, f, text)

	p.Call(jsw.Remove)

	return output
}

func runePos(p js.Value, f *font, index int) (int, int) {
	f.r.Call(jsw.SetStart, p.Get(jsw.FirstChild), index)
	f.r.Call(jsw.SetEnd, p.Get(jsw.FirstChild), index+1)
	offsets := f.r.Call(jsw.GetBoundingClientRect)
	return offsets.Get(jsw.Left).Int(), offsets.Get(jsw.Top).Int()
}

func textSplit(p js.Value, f *font, text string) []int {
	var output []int
	x0 := 0
	cut := 0
	index := 0
	for i := range text {
		x, _ := runePos(p, f, index)
		index++
		if x < x0 {
			output = append(output, i-cut)
			cut = i
		}
		x0 = x
	}
	output = append(output, len(text)-cut)
	return output
}

func (w *webAPI) FontMetricSize(fontID int, text string) (int, int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	f, ok := w.metricFonts[fontID]
	if !ok {
		log.Printf("FontMetricSize: font not found: %d", fontID)
		return 0, 0
	}

	canvasCtx := f.canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Set(jsw.Font, f.cssValue)
	metrics := canvasCtx.Call(jsw.MeasureText, text)
	return metrics.Get(jsw.Width).Int(), f.height
}
