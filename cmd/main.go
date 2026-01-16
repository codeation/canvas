package main

import (
	"bufio"
	"sync"

	"github.com/codeation/impress/joint/drawrecv"
	"github.com/codeation/impress/joint/eventsend"
	"github.com/codeation/impress/joint/rpc"

	"github.com/codeation/canvas/jsw/clientsocket"
	"github.com/codeation/canvas/webapi"
	"github.com/codeation/canvas/webevent"
)

func main() {
	asyncSocket := clientsocket.Dial("async")
	defer asyncSocket.Close()
	syncSocket := clientsocket.Dial("sync")
	defer syncSocket.Close()

	eventPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(asyncSocket), nil)
	streamPipe := rpc.NewPipe(rpc.WithoutMutex(), nil, asyncSocket)
	syncPipe := rpc.NewPipe(rpc.WithoutMutex(), bufio.NewWriter(syncSocket), syncSocket)

	api := webapi.New()
	d := drawrecv.New(api, streamPipe, syncPipe)
	eventSend := eventsend.New(eventPipe)
	webEvent := webevent.New(eventSend)
	defer webEvent.Done()

	d.Wait()
}
