package schedule

import (
	"context"
	"time"

	"github.com/indeedhat/chonker/internal/types"
)

type TriggerFunc func(float64) types.Report
type ReportFunc func(types.Report)

type Runner struct {
	sched   Schedule
	trigger TriggerFunc
	report  ReportFunc
}

func NewRunner(s Schedule, t TriggerFunc, r ReportFunc) Runner {
	return Runner{s, t, r}
}

func (r Runner) Start(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			meal := r.sched.Now()
			if meal == nil {
				break
			}

			if rep := r.trigger(meal.Weight); rep.Error() != nil || r.sched.ReportOnSuccess {
				r.report(rep)
			}
		}
	}
}
