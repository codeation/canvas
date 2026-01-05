package clientsocket

import (
	"encoding/base64"
	"errors"
	"io"
	"syscall/js"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/canvas/jsw/eventlist"
)

type ClientSocket struct {
	ws          js.Value
	listeners   *eventlist.EventListeners
	waitOpening chan struct{}
	reader      io.Reader
	writer      *io.PipeWriter
	closer      *io.PipeReader
}

func Dial(address string) *ClientSocket {
	pipeReader, pipeWriter := io.Pipe()
	ws := js.Global().Get(jsw.WebSocket).New(address)
	c := &ClientSocket{
		ws:          ws,
		listeners:   eventlist.NewEventListeners(ws),
		waitOpening: make(chan struct{}),
		reader:      base64.NewDecoder(base64.StdEncoding, pipeReader),
		writer:      pipeWriter,
		closer:      pipeReader,
	}
	c.listeners.Add(jsw.Open, c.onOpen)
	c.listeners.Add(jsw.Error, c.onError)
	c.listeners.Add(jsw.Close, c.onClose)
	c.listeners.Add(jsw.Message, c.onMessage)
	return c
}

func (c *ClientSocket) Close() error {
	c.closer.Close()
	c.listeners.Done()
	c.ws.Call(jsw.Close)
	return nil
}

func (c *ClientSocket) Write(data []byte) (int, error) {
	<-c.waitOpening
	c.ws.Call(jsw.Send, base64.StdEncoding.EncodeToString(data))
	return len(data), nil
}

func (c *ClientSocket) Read(data []byte) (int, error) {
	return c.reader.Read(data)
}

func (c *ClientSocket) onOpen(this js.Value, args []js.Value) any {
	close(c.waitOpening)
	return nil
}

func (c *ClientSocket) onError(this js.Value, args []js.Value) any {
	c.closer.CloseWithError(errors.New(args[0].Get(jsw.Type).String()))
	return nil
}

func (c *ClientSocket) onClose(this js.Value, args []js.Value) any {
	c.Close()
	return nil
}

func (c *ClientSocket) onMessage(this js.Value, args []js.Value) any {
	message := args[0]
	if _, err := c.writer.Write([]byte(message.Get(jsw.Data).String())); err != nil {
		c.writer.CloseWithError(err)
	}
	return nil
}
