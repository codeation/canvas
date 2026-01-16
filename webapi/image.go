package webapi

import (
	"log"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type image struct {
	canvas js.Value
	width  int
	height int
}

func (w *webAPI) ImageNew(imageID int, width, height int, bitmap []byte) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	uint8Array := js.Global().Get(jsw.Uint8ClampedArray).New(len(bitmap))
	js.CopyBytesToJS(uint8Array, bitmap)
	imageData := js.Global().Get(jsw.ImageData).New(uint8Array, width, height)
	canvas := js.Global().Get(jsw.OffscreenCanvas).New(width, height)
	canvasCtx := canvas.Call(jsw.GetContext, jsw.Context2D)
	canvasCtx.Call(jsw.PutImageData, imageData, 0, 0)

	w.images[imageID] = &image{
		canvas: canvas,
		width:  width,
		height: height,
	}
}

func (w *webAPI) ImageDrop(imageID int) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	image, ok := w.images[imageID]
	if !ok {
		log.Printf("ImageDrop: image not found: %d", imageID)
		return
	}

	image.canvas.Call(jsw.Remove)

	delete(w.images, imageID)
}
