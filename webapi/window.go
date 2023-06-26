package webapi

import (
	"log"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type window struct {
	id        int
	frameID   int
	x, y      int
	canvas    js.Value
	canvasCtx js.Value
}

func (w *webAPI) WindowNew(windowID int, frameID int, x, y, width, height int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	frame, ok := w.frames[frameID]
	if !ok {
		log.Printf("WindowNew: frame not found: %d", frameID)
		return
	}

	offsetX, offsetY := w.frameOffset(frameID)
	x += offsetX
	y += offsetY

	canvas := w.document.Call(jsw.CreateElement, jsw.Canvas)
	canvasStyle := canvas.Get(jsw.Style)
	canvasStyle.Set(jsw.Position, jsw.Fixed)
	canvasStyle.Set(jsw.Left, px(x))
	canvasStyle.Set(jsw.Top, px(y))
	canvasStyle.Set(jsw.Width, px(width))
	canvasStyle.Set(jsw.Height, px(height))
	canvas.Set(jsw.Width, width*2)
	canvas.Set(jsw.Height, height*2)
	canvasCtx := canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Call(jsw.Scale, 2, 2)

	frame.div.Call(jsw.AppendChild, canvas)

	w.windows[windowID] = &window{
		id:        windowID,
		frameID:   frameID,
		canvas:    canvas,
		canvasCtx: canvasCtx,
		x:         x,
		y:         y,
	}
}

func (w *webAPI) WindowDrop(windowID int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	window, ok := w.windows[windowID]
	if !ok {
		log.Printf("WindowDrop: window not found: %d", windowID)
		return
	}

	frame, ok := w.frames[window.frameID]
	if !ok {
		log.Printf("WindowDrop: frame not found: %d", window.frameID)
		return
	}

	frame.div.Call(jsw.RemoveChild, window.canvas)
	window.canvas.Call(jsw.Remove)

	delete(w.windows, windowID)
}

func (w *webAPI) WindowRaise(windowID int) {}
func (w *webAPI) WindowClear(windowID int) {}
func (w *webAPI) WindowShow(windowID int)  {}

func (w *webAPI) WindowSize(windowID int, x, y, width, height int) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	window, ok := w.windows[windowID]
	if !ok {
		log.Printf("WindowSize: window not found: %d", windowID)
		return
	}
	window.x = x
	window.y = y

	offsetX, offsetY := w.frameOffset(window.frameID)
	x += offsetX
	y += offsetY

	canvasStyle := window.canvas.Get(jsw.Style)
	canvasStyle.Set(jsw.Left, px(x))
	canvasStyle.Set(jsw.Top, px(y))
	if window.canvas.Get(jsw.Width).Int() != width*2 ||
		window.canvas.Get(jsw.Height).Int() != height*2 {
		canvasStyle.Set(jsw.Width, px(width))
		canvasStyle.Set(jsw.Height, px(height))
		window.canvas.Set(jsw.Width, width*2)
		window.canvas.Set(jsw.Height, height*2)
		window.canvasCtx.Call(jsw.Scale, 2, 2)
	}
}

func (w *webAPI) WindowFill(windowID int, x, y, width, height int, r, g, b, a uint16) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	window, ok := w.windows[windowID]
	if !ok {
		log.Printf("WindowFill: window not found: %d", windowID)
		return
	}

	window.canvasCtx.Set(jsw.FillStyle, color(r, g, b, a))
	window.canvasCtx.Call(jsw.FillRect, x, y, width, height)
}

func (w *webAPI) WindowLine(windowID int, x0, y0, x1, y1 int, r, g, b, a uint16) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	window, ok := w.windows[windowID]
	if !ok {
		log.Printf("WindowLine: window not found: %d", windowID)
		return
	}

	window.canvasCtx.Set(jsw.StrokeStyle, color(r, g, b, a))
	window.canvasCtx.Set(jsw.LineWidth, 1)
	window.canvasCtx.Call(jsw.BeginPath)
	window.canvasCtx.Call(jsw.MoveTo, float64(x0)+0.5, float64(y0)+0.5)
	window.canvasCtx.Call(jsw.LineTo, float64(x1)+0.5, float64(y1)+0.5)
	window.canvasCtx.Call(jsw.Stroke)
}

func (w *webAPI) WindowText(windowID int, x, y int, r, g, b, a uint16, fontID int, height int, text string) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	window, ok := w.windows[windowID]
	if !ok {
		log.Printf("WindowText: window not found: %d", windowID)
		return
	}

	font, ok := w.fonts[fontID]
	if !ok {
		log.Printf("WindowText: font not found %d", fontID)
		return
	}

	window.canvasCtx.Set(jsw.Font, font.cssValue)
	window.canvasCtx.Set(jsw.FillStyle, color(r, g, b, a))
	window.canvasCtx.Call(jsw.FillText, text, x, y+font.baseline)
}

func (w *webAPI) WindowImage(windowID int, x, y, width, height int, imageID int) {}
