package kernel

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/coreos/go-systemd/sdjournal"
)

const DefaultJournaldPath = "/run/log/journal"

type JournaldLogs struct {
	Path string
}

func (jls JournaldLogs) Follow(until <-chan time.Time, writer io.Writer) error {
	path := DefaultJournaldPath
	if jls.Path != "" {
		path = jls.Path
	}

	c := sdjournal.JournalReaderConfig{
		Since: 1 * time.Nanosecond,
		Matches: []sdjournal.Match{{
			Field: sdjournal.SD_JOURNAL_FIELD_TRANSPORT,
			Value: "kernel",
		}},
		Path: path,
	}

	jr, err := sdjournal.NewJournalReader(c)
	if err != nil {
		return fmt.Errorf("Could not open journal with error: %s", err)
	}

	if jr == nil {
		return errors.New("Could not open journal. Got an invalid reader")
	}

	defer jr.Close()

	if err := jr.Follow(until, writer); err != nil {
		return fmt.Errorf("Could not read from journal with error: %s", err)
	}
	return nil
}
