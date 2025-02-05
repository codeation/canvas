// Package implements an internal mechanism to communicate with an impress terminal.
package webapi

import (
	"sync"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
)

type webAPI struct {
	document js.Value
	window   js.Value
	frames   map[int]*frame
	windows  map[int]*window
	fonts    map[int]*font
	mutex    sync.RWMutex
}

func New() *webAPI {
	w := &webAPI{
		document: js.Global().Get(jsw.Document),
		window:   js.Global().Get(jsw.Window),
		frames:   map[int]*frame{},
		windows:  map[int]*window{},
		fonts:    map[int]*font{},
	}

	w.document.Get(jsw.Body).Get(jsw.Style).Set(jsw.Margin, 0)
	w.document.Get(jsw.Body).Get(jsw.Style).Set(jsw.Padding, 0)

	return w
}

func (w *webAPI) ApplicationTitle(title string)           {}
func (w *webAPI) ApplicationSize(x, y, width, height int) {}
func (w *webAPI) ApplicationExit()                        {}
func (w *webAPI) ApplicationVersion() string              { return canvasAPIVersion }
func (w *webAPI) ClipboardGet(typeID int)                 {}
func (w *webAPI) ClipboardPut(typeID int, data []byte)    {}

func (w *webAPI) Sync() {}
