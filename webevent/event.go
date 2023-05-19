package webevent

import (
	"syscall/js"
	"unicode/utf8"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/impress/event"
	"github.com/codeation/impress/joint/iface"
)

type webEvent struct {
	callbacks iface.CallbackSet
	window    js.Value
}

func New(callbacks iface.CallbackSet) {
	w := &webEvent{
		callbacks: callbacks,
		window:    js.Global().Get(jsw.Window),
	}

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

}

func (w *webEvent) onResize(this js.Value, args []js.Value) any {
	w.callbacks.EventConfigure(
		w.window.Get(jsw.OuterWidth).Int(), w.window.Get(jsw.OuterHeight).Int(),
		w.window.Get(jsw.InnerWidth).Int(), w.window.Get(jsw.InnerHeight).Int())
	return js.ValueOf(true)
}

func (w *webEvent) onButton(this js.Value, args []js.Value, action int) any {
	if len(args) < 1 {
		return js.ValueOf(false)
	}
	pointerEvent := args[0]
	w.callbacks.EventButton(
		action,
		pointerEvent.Get(jsw.Button).Int()+1,
		pointerEvent.Get(jsw.ClientX).Int(),
		pointerEvent.Get(jsw.ClientY).Int())
	return js.ValueOf(true)
}

func (w *webEvent) onPointerDown(this js.Value, args []js.Value) any {
	return w.onButton(this, args, event.ButtonActionPress)
}

func (w *webEvent) onPointerUp(this js.Value, args []js.Value) any {
	return w.onButton(this, args, event.ButtonActionRelease)
}

func (w *webEvent) onDoubleClick(this js.Value, args []js.Value) any {
	return w.onButton(this, args, event.ButtonActionDouble)
}

func (w *webEvent) onContextMenu(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(false)
	}
	args[0].Call(jsw.PreventDefault)
	return js.ValueOf(false)
}

func (w *webEvent) onMousemove(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(false)
	}
	mouseEvent := args[0]
	w.callbacks.EventMotion(
		mouseEvent.Get(jsw.ClientX).Int(),
		mouseEvent.Get(jsw.ClientY).Int(),
		mouseEvent.Get(jsw.ShiftKey).Bool(),
		mouseEvent.Get(jsw.CtrlKey).Bool(),
		mouseEvent.Get(jsw.AltKey).Bool(),
		mouseEvent.Get(jsw.MetaKey).Bool())
	return js.ValueOf(true)
}

func (w *webEvent) onWheel(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(false)
	}
	wheelEvent := args[0]
	w.callbacks.EventScroll(
		event.ScrollSmooth,
		wheelEvent.Get(jsw.DeltaX).Int(),
		wheelEvent.Get(jsw.DeltaY).Int())
	return js.ValueOf(false)
}

func (w *webEvent) onKeyDown(this js.Value, args []js.Value) any {
	if len(args) < 1 {
		return js.ValueOf(false)
	}
	keyboardEvent := args[0]
	key := keyboardEvent.Get(jsw.Key).String()
	r, length := utf8.DecodeRuneInString(key)
	if length != len(key) {
		r = 0
	}
	w.callbacks.EventKeyboard(r,
		keyboardEvent.Get(jsw.ShiftKey).Bool(),
		keyboardEvent.Get(jsw.CtrlKey).Bool(),
		keyboardEvent.Get(jsw.AltKey).Bool(),
		keyboardEvent.Get(jsw.MetaKey).Bool(),
		keyboardEvent.Get(jsw.Code).String())
	return js.ValueOf(true)
}

func (w *webEvent) onUnload(this js.Value, args []js.Value) any {
	w.callbacks.EventGeneral(event.DestroyEvent.Event)
	return js.ValueOf(true)
}
