package webapi

import (
	"log"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type frame struct {
	id       int
	parentID int
	div      js.Value
	x, y     int
}

func (w *webAPI) FrameNew(frameID int, parentFrameID int, x, y, width, height int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	var parent js.Value
	if parentFrameID == 0 {
		parent = w.document.Get(jsw.Body)
	} else {
		parentFrame, ok := w.frames[parentFrameID]
		if !ok {
			log.Printf("FrameNew: parent frame not found: %d", parentFrameID)
			return
		}
		parent = parentFrame.div
	}

	offsetX, offsetY := w.frameOffset(parentFrameID)

	div := w.document.Call(jsw.CreateElement, jsw.Div)
	divStyle := div.Get(jsw.Style)
	divStyle.Set(jsw.Position, jsw.Fixed)
	divStyle.Set(jsw.Left, px(x+offsetX))
	divStyle.Set(jsw.Top, px(y+offsetY))
	divStyle.Set(jsw.Width, px(width))
	divStyle.Set(jsw.Height, px(height))
	divStyle.Set(jsw.Overflow, jsw.Hidden)

	parent.Call(jsw.AppendChild, div)

	w.frames[frameID] = &frame{
		id:       frameID,
		parentID: parentFrameID,
		div:      div,
		x:        x,
		y:        y,
	}
}

func (w *webAPI) FrameDrop(frameID int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	f, ok := w.frames[frameID]
	if !ok {
		log.Printf("FrameDrop: frame not found: %d", frameID)
		return
	}

	var parent js.Value
	if f.parentID == 0 {
		parent = w.document.Get(jsw.Body)
	} else {
		parentFrame, ok := w.frames[f.parentID]
		if !ok {
			log.Printf("FrameDrop: parent frame not found: %d", f.parentID)
			return
		}
		parent = parentFrame.div
	}

	parent.Call(jsw.RemoveChild, f.div)
	f.div.Call(jsw.Remove)

	delete(w.frames, frameID)
}

func (w *webAPI) FrameSize(frameID int, x, y, width, height int) {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	f, ok := w.frames[frameID]
	if !ok {
		log.Printf("FrameSize: frame not found: %d", frameID)
		return
	}

	offsetX, offsetY := w.frameOffset(f.parentID)

	divStyle := f.div.Get(jsw.Style)
	divStyle.Set(jsw.Left, px(x+offsetX))
	divStyle.Set(jsw.Top, px(y+offsetY))
	divStyle.Set(jsw.Width, px(width))
	divStyle.Set(jsw.Height, px(height))

	f.x = x
	f.y = y

	for _, window := range w.windows {
		if window.frameID != frameID {
			continue
		}

		canvasStyle := window.canvas.Get(jsw.Style)
		canvasStyle.Set(jsw.Left, px(window.x+x+offsetX))
		canvasStyle.Set(jsw.Top, px(window.y+y+offsetY))
	}
}

func (w *webAPI) frameOffset(frameID int) (int, int) {
	var x, y int
	for frameID != 0 {
		frame, ok := w.frames[frameID]
		if !ok {
			log.Printf("frameOffset: frame not found: %d", frameID)
			return 0, 0
		}
		x += frame.x
		y += frame.y
		frameID = frame.parentID
	}
	return x, y
}

func (w *webAPI) FrameRaise(frameID int) {}
