// Package implements an internal mechanism to communicate with an impress terminal.
package link

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

type LinkReader struct {
	linkURL string
	rw      io.Writer
}

func NewReader(linkURL string, rw io.Writer) *LinkReader {
	l := &LinkReader{
		linkURL: linkURL,
		rw:      rw,
	}

	go l.Run()

	return l
}

func (l *LinkReader) Run() {
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

func (l *LinkReader) run() error {
	req, err := http.NewRequest(http.MethodPost, l.linkURL, nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("client.Do: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp.StatusCode: %d", resp.StatusCode)
	}

	if _, err := io.Copy(l.rw, resp.Body); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	return nil
}
