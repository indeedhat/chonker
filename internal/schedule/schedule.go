package schedule

import (
	"time"

	"github.com/indeedhat/icl"
)

const schedulePath = ".schedule.icl"

type Schedule struct {
	ReportOnSuccess bool `icl:"report_on_success"`

	Default   Day  `icl:"default"`
	Monday    *Day `icl:"monday"`
	Tuesday   *Day `icl:"tuesday"`
	Wednesday *Day `icl:"wednesday"`
	Thursday  *Day `icl:"thursday"`
	Friday    *Day `icl:"friday"`
	Saturday  *Day `icl:"saturday"`
	Sunday    *Day `icl:"sunday"`
}

func (s Schedule) Now() *Meal {
	day := s.day()
	if day != nil && day.Skip {
		return nil
	} else if day == nil {
		day = &s.Default
	}

	return day.Now()
}

func (s Schedule) day() *Day {
	switch time.Now().Weekday() {
	case time.Sunday:
		return s.Sunday
	case time.Monday:
		return s.Monday
	case time.Tuesday:
		return s.Tuesday
	case time.Wednesday:
		return s.Wednesday
	case time.Thursday:
		return s.Thursday
	case time.Friday:
		return s.Friday
	case time.Saturday:
		return s.Saturday
	default:
		return nil
	}

}

type Day struct {
	Skip  bool   `icl:"skip"`
	Meals []Meal `icl:"meal"`
}

func (d Day) Now() *Meal {
	for _, m := range d.Meals {
		if m.TimeOfDay == time.Now().Format("15:04") {
			return &m
		}
	}

	return nil
}

type Meal struct {
	TimeOfDay string  `icl:".param"`
	Weight    float32 `icl:"weight_g"`
	Silent    bool    `icl:"silent"`
}

func Load() (Schedule, error) {
	var sched Schedule

	err := icl.UnMarshalFile(schedulePath, &sched)

	return sched, err
}

func Save(sched Schedule) error {
	return icl.MarshalFile(sched, schedulePath)
}
