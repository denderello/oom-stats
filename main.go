package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/coreos/go-systemd/sdjournal"
)

func main() {
	c := sdjournal.JournalReaderConfig{
		NumFromTail: 30,
		Matches: []sdjournal.Match{{
			Field: sdjournal.SD_JOURNAL_FIELD_TRANSPORT,
			Value: "kernel",
		}},
		Path: "/run/log/journal",
	}

	r, err := sdjournal.NewJournalReader(c)
	if err != nil {
		log.Fatalf("Could not open journal with error: %s", err)
	}

	if r == nil {
		log.Fatal("Could not open journal. Got an invalid reader")
	}

	defer r.Close()

	readBuf, err := ioutil.ReadAll(r)
	if err != nil {
		log.Fatalf("Could not read from journal with error: %s", err)
	}

	journalBuf := bytes.NewBuffer(readBuf)

	fmt.Printf("Read %d bytes from Kernel journal\n", journalBuf.Len())
	if journalBuf.Len() > 0 {
		fmt.Printf("Content of journal:\n\n%s", journalBuf)
	}
}
