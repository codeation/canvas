package clientsocket

import (
	"encoding/base64"
	"io"
	"sync"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/canvas/jsw/eventlist"
)

type ClientSocket struct {
	ws          js.Value
	eventList   *eventlist.EventList
	isConnected chan struct{}
	waitOnce    sync.Once
	doneOnce    sync.Once
	recvR       io.Reader
	recvW       *io.PipeWriter
	recvPipeR   *io.PipeReader
}

func Dial(address string) *ClientSocket {
	pipeR, pipeW := io.Pipe()
	ws := js.Global().Get(jsw.WebSocket).New(address)
	c := &ClientSocket{
		ws:          ws,
		eventList:   eventlist.NewEventList(ws),
		isConnected: make(chan struct{}),
		recvR:       base64.NewDecoder(base64.StdEncoding, pipeR),
		recvW:       pipeW,
		recvPipeR:   pipeR,
	}
	c.eventList.Add(jsw.Open, c.onOpen)
	c.eventList.Add(jsw.Error, c.onError)
	c.eventList.Add(jsw.Close, c.onClose)
	c.eventList.Add(jsw.Message, c.onMessage)
	return c
}

func (c *ClientSocket) Close() error {
	c.doneOnce.Do(c.doneConnect)
	c.recvPipeR.Close()
	c.eventList.Done()
	c.ws.Call(jsw.Close)
	return nil
}

func (c *ClientSocket) Write(data []byte) (int, error) {
	c.waitOnce.Do(c.waitConnect)
	c.ws.Call(jsw.Send, base64.StdEncoding.EncodeToString(data))
	return len(data), nil
}

func (c *ClientSocket) Read(data []byte) (int, error) {
	return c.recvR.Read(data)
}

func (c *ClientSocket) waitConnect() {
	<-c.isConnected
}

func (c *ClientSocket) doneConnect() {
	close(c.isConnected)
}

func (c *ClientSocket) onOpen(this js.Value, args []js.Value) any {
	c.doneOnce.Do(c.doneConnect)
	return js.ValueOf(true)
}

func (c *ClientSocket) onError(this js.Value, args []js.Value) any {
	c.doneOnce.Do(c.doneConnect)
	return js.ValueOf(true)
}

func (c *ClientSocket) onClose(this js.Value, args []js.Value) any {
	c.Close()
	return js.ValueOf(true)
}

func (c *ClientSocket) onMessage(this js.Value, args []js.Value) any {
	message := args[0]
	c.recvW.Write([]byte(message.Get(jsw.Data).String()))
	return js.ValueOf(true)
}
