package main

import (
	"bufio"
	"fmt"
	"io"
	"sync"

	"github.com/codeation/canvas/link"
	"github.com/codeation/canvas/webapi"
	"github.com/codeation/impress/joint/iosplit"
	"github.com/codeation/impress/joint/remote"
	"github.com/codeation/impress/joint/rpc"
)

func run() error {

	streamR, streamW := io.Pipe()
	_ = link.NewReader("stream", streamW)
	requestR, requestW := io.Pipe()
	_ = link.NewReader("request", requestW)
	responseLink := link.NewWriter("response", iosplit.NewBufferSplit().WithEternal())
	eventLink := link.NewWriter("event", iosplit.NewBufferSplit().WithEternal())

	eventPipe := rpc.NewPipe(new(sync.Mutex), bufio.NewWriter(eventLink), nil)
	streamPipe := rpc.NewPipe(rpc.WithoutMutex(), nil, streamR)
	syncPipe := rpc.NewPipe(rpc.WithoutMutex(), bufio.NewWriter(responseLink), requestR)

	w := webapi.New(remote.NewCallBacks(eventPipe))
	_ = remote.NewServer(w, streamPipe, syncPipe)

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()

	return nil
}

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
	}
}
