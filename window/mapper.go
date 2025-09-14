package window

import (
	"fmt"
	"time"

	"github.com/copito/coptime/interval"
	"github.com/teambition/rrule-go"
)

func MapRRuleWeekToWeekday(wkday rrule.Weekday) time.Weekday {
	var wkst time.Weekday
	switch wkday {
	case rrule.MO:
		wkst = time.Monday
	case rrule.TU:
		wkst = time.Tuesday
	case rrule.WE:
		wkst = time.Wednesday
	case rrule.TH:
		wkst = time.Thursday
	case rrule.FR:
		wkst = time.Friday
	case rrule.SA:
		wkst = time.Saturday
	case rrule.SU:
		wkst = time.Sunday
	}

	return wkst
}

func MapRRULEFrequencyToFrequency(frequency rrule.Frequency) (interval.Frequency, error) {
	var freq interval.Frequency
	switch frequency {
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
		return interval.FrequencyDay, fmt.Errorf("unsupported frequency: %v", frequency)
	}

	return freq, nil
}
