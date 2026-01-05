package eventlist

import (
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type eventFunc struct {
	event  string
	funcOf js.Func
}

type EventListeners struct {
	original   js.Value
	eventFuncs []eventFunc
}

func NewEventListeners(original js.Value) *EventListeners {
	return &EventListeners{
		original: original,
	}
}

func (el *EventListeners) Add(event string, f func(this js.Value, args []js.Value) any) {
	funcOf := js.FuncOf(f)
	el.original.Call(jsw.AddEventListener, event, funcOf)
	el.eventFuncs = append(el.eventFuncs, eventFunc{
		event:  event,
		funcOf: funcOf,
	})
}

func (el *EventListeners) Done() {
	for _, e := range el.eventFuncs {
		el.original.Call(jsw.RemoveEventListener, e.event, e.funcOf)
		e.funcOf.Release()
	}
}
