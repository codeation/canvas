package eventlist

import (
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type eventFunc struct {
	event  string
	funcOf js.Func
}

type EventList struct {
	original   js.Value
	eventFuncs []eventFunc
}

func NewEventList(original js.Value) *EventList {
	return &EventList{
		original: original,
	}
}

func (el *EventList) Add(event string, f func(this js.Value, args []js.Value) any) {
	funcOf := js.FuncOf(f)
	el.original.Call(jsw.AddEventListener, event, funcOf)
	el.eventFuncs = append(el.eventFuncs, eventFunc{
		event:  event,
		funcOf: funcOf,
	})
}

func (el *EventList) Done() {
	for _, e := range el.eventFuncs {
		el.original.Call(jsw.RemoveEventListener, e.event, e.funcOf)
		e.funcOf.Release()
	}
}
