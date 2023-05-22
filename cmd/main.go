package main

import (
	"bufio"
	"sync"
	"syscall/js"

	"github.com/codeation/impress/joint/drawrecv"
	"github.com/codeation/impress/joint/eventsend"
	"github.com/codeation/impress/joint/rpc"

	"github.com/codeation/canvas/jsw"
	"github.com/codeation/canvas/jsw/clientsocket"
	"github.com/codeation/canvas/webapi"
	"github.com/codeation/canvas/webevent"
)

func host() string {
	return js.Global().Get(jsw.Window).Get(jsw.Location).Get(jsw.Host).String()
}

func main() {
	streamSocket := clientsocket.Dial("ws://" + host() + "/stream")
	defer streamSocket.Close()
	syncSocket := clientsocket.Dial("ws://" + host() + "/sync")
	defer syncSocket.Close()
	eventSocket := clientsocket.Dial("ws://" + host() + "/event")
	defer eventSocket.Close()

	eventPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(eventSocket), nil)
	streamPipe := rpc.NewPipe(rpc.WithoutMutex(), nil, streamSocket)
	syncPipe := rpc.NewPipe(rpc.WithoutMutex(), bufio.NewWriter(syncSocket), syncSocket)

	api := webapi.New()
	d := drawrecv.New(api, streamPipe, syncPipe)
	eventSend := eventsend.New(eventPipe)
	webEvent := webevent.New(eventSend)
	defer webEvent.Done()

	d.Wait()
}
