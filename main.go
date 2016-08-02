package main

import (
	"io"
	"log"
	"time"

	"github.com/denderello/oom-stats/aggregation"
	"github.com/denderello/oom-stats/filter"
	"github.com/denderello/oom-stats/kernel"
	"github.com/denderello/oom-stats/metric"
	"github.com/denderello/oom-stats/oom"
)

func main() {
	r, w := io.Pipe()
	results := make(chan []byte)
	oomAggr := aggregation.NewOOM()

	go func() {
		doneChan := make(<-chan time.Time)

		jf := kernel.JournaldLogs{}
		if err := jf.Follow(doneChan, w); err != nil {
			log.Fatalf("Could not read from journal with error: %s", err)
		}
	}()

	go func() {
		of := filter.OOM{}
		of.Filter(r, results)
	}()

	go func() {
		m := &metric.Stdout{
			Interval: 5 * time.Second,
			OOM:      oomAggr,
		}
		m.Run()
	}()

	for {
		select {
		case result := <-results:
			report, err := oom.ParseReportFromKernelLogs(result)
			if err != nil {
				log.Printf("Error parsing kernel logs: %s", err)
			}
			oomAggr.Aggregate(report)
		}
	}
}
