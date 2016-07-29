package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

type OOMMessageFilter struct{}

func (f OOMMessageFilter) Filter(r io.ReadCloser) {
	defer r.Close()

	buf := bufio.NewReader(r)
	for {
		line, _ := buf.ReadBytes('\n')
		fmt.Printf("Got journal line: %s\n", line)
	}
}

func main() {
	r, w := io.Pipe()

	c := sdjournal.JournalReaderConfig{
		NumFromTail: 30,
		Matches: []sdjournal.Match{{
			Field: sdjournal.SD_JOURNAL_FIELD_TRANSPORT,
			Value: "kernel",
		}},
		Path: "/run/log/journal",
	}

	jr, err := sdjournal.NewJournalReader(c)
	if err != nil {
		log.Fatalf("Could not open journal with error: %s", err)
	}

	if jr == nil {
		log.Fatal("Could not open journal. Got an invalid reader")
	}

	defer jr.Close()

	go func() {
		doneChan := make(<-chan time.Time)
		err = jr.Follow(doneChan, w)
		if err != nil {
			log.Fatalf("Could not read from journal with error: %s", err)
		}
	}()

	f := OOMMessageFilter{}
	go f.Filter(r)

	for {
	}
}
