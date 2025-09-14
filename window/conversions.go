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
	rr, err := rrule.StrToRRule(ruleString)
	if err != nil {
		return nil, err
	}

	// if len(rr.Options.Bymonth) > 0 || len(rr.Options.Byweekno) > 0 || len(rr.Options.Byyearday) > 0 ||
	// 	len(rr.Options.Bymonthday) > 0 || len(rr.Options.Byweekday) > 0 || len(rr.Options.Bysetpos) > 0 ||
	// 	len(rr.Options.Byhour) > 0 || len(rr.Options.Byminute) > 0 || len(rr.Options.Bysecond) > 0 {
	// 	return nil, errors.New("unsupported rrule feature: BY* rules are not yet supported")
	// }

	var freq interval.Frequency
	switch rr.Options.Freq {
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
		return nil, fmt.Errorf("unsupported frequency: %v", rr.Options.Freq)
	}

	intervalValue := uint32(rr.Options.Interval)
	if intervalValue == 0 {
		intervalValue = 1
	}

	if rr.Options.Dtstart.IsZero() {
		return nil, errors.New("rrule must have a DTSTART")
	}
	anchorDate := rr.Options.Dtstart
	startDate := &rr.Options.Dtstart

	var endDate *time.Time
	if !rr.Options.Until.IsZero() {
		endDate = &rr.Options.Until
	}

	// Handling count given this is a post param on between and all
	var sizeValue float64
	if rr.Options.Count > 0 {
		sizeValue = float64(int32(rr.Options.Count))
	} else {
		sizeValue = math.Inf(1)
	}

	// Calculate the rules based on the rr
	// rr.Options.
	_ = rr.Options.Byeaster

	return &WindowOption{
		IntervalOption: interval.IntervalOption{
			AnchorDate:    anchorDate,
			StartDate:     startDate,
			EndDate:       endDate,
			Size:          &sizeValue,
			FrequencyUnit: freq,
			IntervalValue: intervalValue,
		},
	}, nil
}

func convertCrontoIntervaler(cronString string) (*WindowOption, error) {
	return nil, errors.New("not implemented")
}
