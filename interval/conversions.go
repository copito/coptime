package interval

import (
	"errors"
	"fmt"
	"time"

	"github.com/teambition/rrule-go"
)

func rruleWeekdayToTimeWeekday(wd rrule.Weekday) time.Weekday {
	switch wd {
	case rrule.SU:
		return time.Sunday
	case rrule.MO:
		return time.Monday
	case rrule.TU:
		return time.Tuesday
	case rrule.WE:
		return time.Wednesday
	case rrule.TH:
		return time.Thursday
	case rrule.FR:
		return time.Friday
	case rrule.SA:
		return time.Saturday
	}
	return time.Sunday // Should not happen
}

func convertRRULEtoIntervaler(ruleString string) (*IntervalOption, error) {
	ro, err := rrule.StrToROption(ruleString)
	if err != nil {
		return nil, err
	}
	r, err := rrule.NewRRule(*ro)
	if err != nil {
		return nil, err
	}

	var monthEnd bool
	if len(ro.Bymonthday) == 1 && ro.Bymonthday[0] == -1 {
		monthEnd = true
	} else if len(ro.Bymonth) > 0 || len(ro.Byweekno) > 0 || len(ro.Byyearday) > 0 || len(ro.Bymonthday) > 0 || len(ro.Bysetpos) > 0 || len(ro.Byhour) > 0 || len(ro.Byminute) > 0 || len(ro.Bysecond) > 0 {
		return nil, errors.New("unsupported rrule feature: BY* rules other than BYMONTHDAY=-1 and BYDAY are not yet supported")
	}

	var byDay []int
	if len(ro.Byweekday) > 0 {
		for _, wd := range ro.Byweekday {
			byDay = append(byDay, int(rruleWeekdayToTimeWeekday(wd)))
		}
	}

	wkst := rruleWeekdayToTimeWeekday(ro.Wkst)

	var freq Frequency
	switch ro.Freq {
	case rrule.YEARLY:
		freq = FrequencyYear
	case rrule.MONTHLY:
		freq = FrequencyMonth
	case rrule.WEEKLY:
		freq = FrequencyWeek
	case rrule.DAILY:
		freq = FrequencyDay
	case rrule.HOURLY:
		freq = FrequencyHour
	case rrule.MINUTELY:
		freq = FrequencyMinute
	case rrule.SECONDLY:
		freq = FrequencySecond
	default:
		return nil, fmt.Errorf("unsupported frequency: %v", ro.Freq)
	}

	interval := uint32(ro.Interval)
	if interval == 0 {
		interval = 1
	}

	if ro.Dtstart.IsZero() {
		return nil, errors.New("rrule must have a DTSTART")
	}
	anchorDate := ro.Dtstart
	startDate := &ro.Dtstart

	var endDate *time.Time
	if !ro.Until.IsZero() {
		endDate = &ro.Until
	} else if ro.Count > 0 {
		all := r.All()
		if len(all) > 0 {
			var newAll []time.Time
			if len(all) > ro.Count {
				newAll = all[:ro.Count]
			} else {
				newAll = all
			}
			lastDate := newAll[len(newAll)-1]
			endDate = &lastDate
		}
	}

	return &IntervalOption{
		AnchorDate:    anchorDate,
		StartDate:     startDate,
		EndDate:       endDate,
		FrequencyUnit: freq,
		IntervalValue: interval,
		MonthEnd:      monthEnd,
		ByDay:         byDay,
		Wkst:          wkst,
	}, nil
}

func convertCrontoIntervaler(cronString string) (*IntervalOption, error) {
	return nil, errors.New("not implemented")
}
