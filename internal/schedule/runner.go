package schedule

import (
	"context"
	"time"
)

type TriggerFunc func(float32) error
type ReportFunc func(error)

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
			if meal != nil {
				err := r.trigger(meal.Weight)
				if err != nil || r.sched.ReportOnSuccess {
					r.report(err)
				}
			}
		}
	}
}
