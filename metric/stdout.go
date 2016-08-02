package metric

import (
	"fmt"
	"time"

	"github.com/denderello/oom-stats/aggregation"
)

type Stdout struct {
	Interval time.Duration
	OOM      *aggregation.OOM
}

func (stdout *Stdout) Run() {
	for range time.Tick(stdout.Interval) {
		if stdout.OOM != nil {
			fmt.Println("Killed processes by name:")
			for process, count := range stdout.OOM.ProcessCounter {
				fmt.Printf(" - %s: %d\n", process, count)
			}
		}
	}
}
