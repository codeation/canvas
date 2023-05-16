// Package implements an internal mechanism to communicate with an impress terminal.
package webapi

import (
	"sync"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/impress/joint/iface"
)

type webAPI struct {
	document  js.Value
	window    js.Value
	frames    map[int]*frame
	windows   map[int]*window
	fonts     map[int]*font
	mutex     sync.RWMutex
	callbacks iface.CallbackSet
}

func New(callbacks iface.CallbackSet) *webAPI {
	w := &webAPI{
		document:  js.Global().Get(jsw.Document),
		window:    js.Global().Get(jsw.Window),
		frames:    map[int]*frame{},
		windows:   map[int]*window{},
		fonts:     map[int]*font{},
		callbacks: callbacks,
	}

	w.document.Get(jsw.Body).Get(jsw.Style).Set(jsw.Margin, 0)
	w.document.Get(jsw.Body).Get(jsw.Style).Set(jsw.Padding, 0)

	w.window.Call(jsw.AddEventListener, jsw.Resize, js.FuncOf(w.onResize))
	w.onResize(js.ValueOf(nil), nil)

	w.window.Call(jsw.AddEventListener, jsw.Pointerup, js.FuncOf(w.onPointerUp))
	w.window.Call(jsw.AddEventListener, jsw.Pointerdown, js.FuncOf(w.onPointerDown))
	w.window.Call(jsw.AddEventListener, jsw.Dblclick, js.FuncOf(w.onDoubleClick))
	w.window.Call(jsw.AddEventListener, jsw.Contextmenu, js.FuncOf(w.onContextMenu))
	w.window.Call(jsw.AddEventListener, jsw.Mousemove, js.FuncOf(w.onMousemove))
	w.window.Call(jsw.AddEventListener, jsw.Wheel, js.FuncOf(w.onWheel))

	w.window.Call(jsw.AddEventListener, jsw.Keydown, js.FuncOf(w.onKeyDown))

	w.window.Call(jsw.AddEventListener, jsw.Unload, js.FuncOf(w.onUnload))
	w.window.Call(jsw.AddEventListener, jsw.Beforeunload, js.FuncOf(w.onUnload))

	return w
}

func (w *webAPI) ApplicationTitle(title string)           {}
func (w *webAPI) ApplicationSize(x, y, width, height int) {}
func (w *webAPI) ApplicationExit() string                 { return "" }
func (w *webAPI) ApplicationVersion() string              { return "" }

func (w *webAPI) Sync() {}
