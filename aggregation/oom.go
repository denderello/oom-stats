package aggregation

import "github.com/denderello/oom-stats/oom"

type OOM struct {
	ProcessCounter map[string]uint
}

func NewOOM() *OOM {
	return &OOM{
		ProcessCounter: map[string]uint{},
	}
}

func (o *OOM) Aggregate(r *oom.Report) {
	o.ProcessCounter[r.ProcessName]++
}
