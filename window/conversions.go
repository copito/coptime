package window

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/copito/coptime/interval"
	"github.com/teambition/rrule-go"
)

func convertRRULEtoWindowOption(ruleString string) (*WindowOption, error) {
	ro, err := rrule.StrToROption(ruleString)
	if err != nil {
		return nil, err
	}
	rr, err := rrule.NewRRule(*ro)
	if err != nil {
		return nil, err
	}

	if len(ro.Bymonth) > 0 || len(ro.Byweekno) > 0 || len(ro.Byyearday) > 0 || len(ro.Bymonthday) > 0 || len(ro.Byweekday) > 0 || len(ro.Bysetpos) > 0 || len(ro.Byhour) > 0 || len(ro.Byminute) > 0 || len(ro.Bysecond) > 0 {
		return nil, errors.New("unsupported rrule feature: BY* rules are not yet supported")
	}

	var freq interval.Frequency
	switch ro.Freq {
	case rrule.YEARLY:
		freq = interval.FrequencyYear
	case rrule.MONTHLY:
		freq = interval.FrequencyMonth
	case rrule.WEEKLY:
		freq = interval.FrequencyWeek
	case rrule.DAILY:
		freq = interval.FrequencyDay
	case rrule.HOURLY:
		freq = interval.FrequencyHour
	case rrule.MINUTELY:
		freq = interval.FrequencyMinute
	case rrule.SECONDLY:
		freq = interval.FrequencySecond
	default:
		return nil, fmt.Errorf("unsupported frequency: %v", ro.Freq)
	}

	intervalValue := uint32(ro.Interval)
	if intervalValue == 0 {
		intervalValue = 1
	}

	if ro.Dtstart.IsZero() {
		return nil, errors.New("rrule must have a DTSTART")
	}
	anchorDate := ro.Dtstart
	startDate := &ro.Dtstart

	var endDate *time.Time
	if !ro.Until.IsZero() {
		endDate = &ro.Until
	}

	// Handling count given this is a post param on between and all
	var sizeValue float64
	if ro.Count > 0 {
		sizeValue = float64(int32(ro.Count))
	} else {
		sizeValue = math.Inf(1)
	}

	// Calculate the rules based on the rr 
	rr.

	return &WindowOption{
		IntervalOption: interval.IntervalOption{
			AnchorDate:    anchorDate,
			StartDate:     startDate,
			EndDate:       endDate,
			FrequencyUnit: freq,
			IntervalValue: intervalValue,
		},
		Size: sizeValue,
	}, nil
}

func convertCrontoIntervaler(cronString string) (*WindowOption, error) {
	return nil, errors.New("not implemented")
}
