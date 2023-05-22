package webapi

import (
	"log"
	"sort"
	"syscall/js"
	"unicode/utf8"

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

func (w *webAPI) FontNew(fontID int, height int, style, variant, weight, stretch int, family string) (int, int, int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	cssValue := fontValue(height, style, variant, weight, stretch, family)
	canvas := w.document.Call(jsw.CreateElement, jsw.Canvas)
	canvasCtx := canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Set(jsw.Font, cssValue)
	metrics := canvasCtx.Call(jsw.MeasureText, "Gg")
	ascent := int(metrics.Get(jsw.FontBoundingBoxAscent).Float())
	descent := int(metrics.Get(jsw.FontBoundingBoxDescent).Float())
	baseline := ascent

	w.fonts[fontID] = &font{
		id:       fontID,
		canvas:   canvas,
		cssValue: cssValue,
		height:   height,
		r:        js.Global().Get(jsw.Range).New(),
		baseline: baseline,
	}

	return baseline, ascent, descent
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

func (w *webAPI) FontSplit(fontID int, text string, edge, indent int) []int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	f, ok := w.fonts[fontID]
	if !ok {
		log.Printf("FontSplit: font not found: %d", fontID)
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

	length := utf8.RuneCountInString(text)
	positions := findZeros(p, f.r, 1, length-1, f.height)
	sort.Slice(positions, func(i, j int) bool { return positions[i] < positions[j] })

	w.document.Get(jsw.Body).Call(jsw.RemoveChild, p)
	p.Call(jsw.Remove)

	return posToLengths(text, positions)
}

func posToLengths(text string, positions []int) []int {
	output := []int{}
	count := 0
	offset := 0
	for i := range text {
		if len(positions) > 0 && count >= positions[0] {
			output = append(output, i-offset)
			positions = positions[1:]
			offset = i
		}
		count++
	}
	if len(text)-offset > 0 {
		output = append(output, len(text)-offset)
	}

	return output

}

func pos(p js.Value, r js.Value, no int) (int, int) {
	r.Call(jsw.SetStart, p.Get(jsw.FirstChild), no)
	r.Call(jsw.SetEnd, p.Get(jsw.FirstChild), no)
	offsets := r.Call(jsw.GetBoundingClientRect)
	return offsets.Get(jsw.Left).Int(), offsets.Get(jsw.Top).Int()
}

func findZeros(p js.Value, r js.Value, left, right int, height int) []int {
	if left > right {
		return nil
	}

	if left == right {
		x, _ := pos(p, r, left)
		if x == 0 {
			return []int{left}
		}
		return nil
	}

	x0, y0 := pos(p, r, left)
	x1, y1 := pos(p, r, right)

	output := []int{}
	if x0 == 0 {
		output = append(output, left)
	}
	if x1 == 0 {
		output = append(output, right)
	}

	if y0 == y1 {
		return output
	}

	if y1-y0 <= height*3/2 && x1 == 0 {
		return output
	}

	output = append(output, findZeros(p, r, left+1, (left+right)/2, height)...)
	output = append(output, findZeros(p, r, (left+right)/2+1, right-1, height)...)
	return output
}

func (w *webAPI) FontSize(fontID int, text string) (int, int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	f, ok := w.fonts[fontID]
	if !ok {
		log.Printf("FontSize: font not found: %d", fontID)
		return 0, 0
	}

	canvasCtx := f.canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Set(jsw.Font, f.cssValue)
	metrics := canvasCtx.Call(jsw.MeasureText, text)
	return metrics.Get(jsw.Width).Int(), f.height
}
