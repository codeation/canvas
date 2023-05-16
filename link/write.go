package link

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LinkWriter struct {
	linkURL string
	rw      io.ReadWriter
}

func NewWriter(linkURL string, rw io.ReadWriter) *LinkWriter {
	l := &LinkWriter{
		linkURL: linkURL,
		rw:      rw,
	}

	go l.Run()

	return l
}

func (l *LinkWriter) Run() {
	var timeout time.Duration
	for {
		if err := l.run(); err != nil {
			log.Println(err)
			timeout += 2 * time.Second
			time.Sleep(timeout)
			continue
		}

		timeout = 0
	}
}

func (l *LinkWriter) Write(p []byte) (int, error) {
	return l.rw.Write(p)
}

func (l *LinkWriter) run() error {
	req, err := http.NewRequest(http.MethodPost, l.linkURL, l.rw)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
	}

	return nil
}
