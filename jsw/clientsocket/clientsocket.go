package clientsocket

import (
	"errors"
	"fmt"
	"io"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/canvas/jsw/eventlist"
)

type ClientSocket struct {
	ws          js.Value
	listeners   *eventlist.EventListeners
	waitOpening chan struct{}
	pipeReader  *io.PipeReader
	pipeWriter  *io.PipeWriter
}

func Dial(address string) *ClientSocket {
	pipeReader, pipeWriter := io.Pipe()
	ws := js.Global().Get(jsw.WebSocket).New(address)
	ws.Set(jsw.BinaryType, jsw.ArrayBuffer)
	c := &ClientSocket{
		ws:          ws,
		listeners:   eventlist.NewEventListeners(ws),
		waitOpening: make(chan struct{}),
		pipeReader:  pipeReader,
		pipeWriter:  pipeWriter,
	}
	c.listeners.Add(jsw.Open, c.onOpen)
	c.listeners.Add(jsw.Error, c.onError)
	c.listeners.Add(jsw.Close, c.onClose)
	c.listeners.Add(jsw.Message, c.onMessage)
	return c
}

func (c *ClientSocket) Close() error {
	c.listeners.Done()
	c.pipeWriter.Close()
	c.ws.Call(jsw.Close)
	return nil
}

func (c *ClientSocket) Write(data []byte) (int, error) {
	<-c.waitOpening
	uint8Array := js.Global().Get(jsw.Uint8Array).New(len(data))
	js.CopyBytesToJS(uint8Array, data)
	c.ws.Call(jsw.Send, uint8Array)
	return len(data), nil
}

func (c *ClientSocket) Read(data []byte) (int, error) {
	return c.pipeReader.Read(data)
}

func (c *ClientSocket) onOpen(this js.Value, args []js.Value) any {
	close(c.waitOpening)
	return nil
}

func (c *ClientSocket) onError(this js.Value, args []js.Value) any {
	err := errors.New("WebSocket error event")
	c.pipeWriter.CloseWithError(err)
	return nil
}

func (c *ClientSocket) onClose(this js.Value, args []js.Value) any {
	err := fmt.Errorf("WebSocket closed: %d", args[0].Get(jsw.Code).Int())
	c.pipeWriter.CloseWithError(err)
	return nil
}

func (c *ClientSocket) onMessage(this js.Value, args []js.Value) any {
	data := args[0].Get(jsw.Data)
	if data.Type() == js.TypeString {
		if _, err := c.pipeWriter.Write([]byte(data.String())); err != nil {
			c.pipeWriter.CloseWithError(err)
		}
		return nil
	}
	uint8Array := js.Global().Get(jsw.Uint8Array).New(data)
	buffer := make([]byte, uint8Array.Get(jsw.ByteLength).Int())
	js.CopyBytesToGo(buffer, uint8Array)
	if _, err := c.pipeWriter.Write(buffer); err != nil {
		c.pipeWriter.CloseWithError(err)
	}
	return nil
}
