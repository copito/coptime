package window

import (
	"errors"
	"fmt"
	"iter"
	"time"

	"github.com/copito/coptime/common"
	"github.com/copito/coptime/interval"
	rules "github.com/copito/coptime/rules"
)

type Windower struct {
	opt WindowOption
}

func New(option WindowOption) Windower {
	return Windower{
		opt: option,
	}
}

func FromRRULE(rruleString string) (*Windower, error) {
	option, err := convertRRULEtoWindowOption(rruleString)
	if err != nil {
		return nil, err
	}

	if option == nil {
		return nil, errors.New("options for Intervaler are empty")
	}

	iv := New(*option)
	return &iv, nil
}

func (w *Windower) estimateIntervalSize(startTime time.Time, endTime time.Time, freq interval.Frequency, intervalValue int32) int32 {
	var realStartTime, realEndTime time.Time
	if startTime.Equal(endTime) {
		return 1
	} else if startTime.Before(endTime) {
		realStartTime = startTime
		realEndTime = endTime
	} else if startTime.After(endTime) {
		realStartTime = endTime
		realEndTime = startTime
	}

	diff := realEndTime.Sub(realStartTime)
	var periods int32

	switch freq {
	case interval.FrequencyNanoSecond:
		diffRef := diff.Nanoseconds()
		periods = int32(diffRef/int64(intervalValue)) + 1
	case interval.FrequencyMilliSecond:
		diffRef := diff.Milliseconds()
		periods = int32(diffRef/int64(intervalValue)) + 1
	case interval.FrequencySecond:
		diffRef := diff.Seconds()
		periods = int32(diffRef/float64(intervalValue)) + 1
	case interval.FrequencyMinute:
		diffRef := diff.Minutes()
		periods = int32(diffRef/float64(intervalValue)) + 1
	case interval.FrequencyHour:
		diffRef := diff.Hours()
		periods = int32(diffRef/float64(intervalValue)) + 1
	case interval.FrequencyDay:
		diffRef := diff.Hours() / 24
		periods = int32(diffRef/float64(intervalValue)) + 1
	case interval.FrequencyWeek:
		diffRef := diff.Hours() / 24 / 7
		periods = int32(diffRef/float64(intervalValue)) + 1
	case interval.FrequencyMonth:
		diffRef := diff.Hours() / 24 / 30
		periods = int32(diffRef/float64(intervalValue)) + 4
	case interval.FrequencyQuarter:
		diffRef := diff.Hours() / 24 / 30 / 3
		periods = int32(diffRef/float64(intervalValue)) + 4
	case interval.FrequencyYear:
		diffRef := diff.Hours() / 24 / 365
		periods = int32(diffRef/float64(intervalValue)) + 4
	default:
		periods = 10
	}

	return periods
}

func (w *Windower) Iterate(direction interval.Direction, maxAttempt *int32) (iter.Seq[common.WindowResult], error) {
	tz := w.opt.AnchorDate.Location()
	if tz == nil {
		tz = time.UTC
	}

	// create interval iterator
	intervalOption := interval.IntervalOption{
		AnchorDate:    w.opt.AnchorDate,
		StartDate:     w.opt.StartDate,
		EndDate:       w.opt.EndDate,
		FrequencyUnit: w.opt.FrequencyUnit,
		IntervalValue: w.opt.IntervalValue,
	}

	// Create the interval iterator
	iv := interval.New(intervalOption)
	iterator, err := iv.Iterate(direction, maxAttempt)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create iterator for window"), err)
	}

	// Handle rules
	includes := rules.FilterRules(w.opt.Rules, rules.RuleTypeInclusion)
	excludes := rules.FilterRules(w.opt.Rules, rules.RuleTypeExclusion)

	return func(yield func(common.WindowResult) bool) {
		next, _ := iter.Pull(iterator)
		// defer stop()

		previousTime, ok := next()
		if !ok {
			return
		}

		for {
			nextTime, ok := next()
			if !ok {
				return
			}

			window := common.WindowResult{
				Start: previousTime,
				End:   nextTime,
			}

			// Apply Rules
			var subWindows []common.SubWindowResult

			// 1. Apply inclusion rules
			var additiveSubWindows []common.SubWindowResult
			if len(includes) == 0 {
				additiveSubWindows = []common.SubWindowResult{{Start: previousTime, End: nextTime}}
			} else {
				additiveSubWindows = rules.GenerateSubWindowsForRuleType(direction, w.opt.FrequencyUnit, previousTime, nextTime, includes, tz)
			}

			// 2. Apply exclusion rules
			subtractiveSubWindows := rules.GenerateSubWindowsForRuleType(direction, w.opt.FrequencyUnit, previousTime, nextTime, excludes, tz)

			// 3. Combine additive and subtractive sub-windows
			subWindows = rules.SubtractSubwindowsFromAdditives(additiveSubWindows, subtractiveSubWindows)

			if len(subWindows) == 0 {
				// Do not add that sessions if no coditions are met
				// (e.g. session window is between 2024-01-01 to 2024-01-02 and the rule are only for 2025-12-25)

				// Move to the next element
				previousTime = nextTime
				continue // skip to the next iteration if no sub-windows are created
			}

			window.SubWindow = subWindows

			// yield the window
			ok = yield(window)
			if !ok {
				return
			}

			previousTime = nextTime
		}
	}, nil
}

func (w *Windower) All(direction interval.Direction, maxAttempt *int32) ([]common.WindowResult, error) {
	var estimatePeriods int32
	if w.opt.StartDate != nil && w.opt.EndDate != nil {
		estimatePeriods = w.estimateIntervalSize(*w.opt.StartDate, *w.opt.EndDate, w.opt.FrequencyUnit, int32(w.opt.IntervalValue))
	} else {
		estimatePeriods = 1000
	}

	if estimatePeriods >= MAX_LIST_SIZE {
		return []common.WindowResult{}, fmt.Errorf("maximum size is %v - please select a smaller range", MAX_LIST_SIZE)
	}

	list := make([]common.WindowResult, 0, estimatePeriods)
	iterator, err := w.Iterate(direction, maxAttempt)
	if err != nil {
		return []common.WindowResult{}, err
	}

	for value := range iterator {
		list = append(list, value)
	}

	return list, nil
}

func (w *Windower) isWindowBetweenTimes(window common.WindowResult, startTime time.Time, endTime time.Time) bool {
	// Entire window must be inside the start and end time
	// Optional: normalize zones to avoid Location diffs
	ws, we := window.Start.In(time.UTC), window.End.In(time.UTC)
	s, e := startTime.In(time.UTC), endTime.In(time.UTC)

	if e.Before(s) { // bad range
		return false
	}

	return (ws.After(s) || ws.Equal(s)) && (we.Before(e) || we.Equal(e))
}

func (w *Windower) Between(direction interval.Direction, startTime time.Time, endTime time.Time, maxAttempt *int32) ([]common.WindowResult, error) {
	estimatePeriods := w.estimateIntervalSize(startTime, endTime, w.opt.FrequencyUnit, int32(w.opt.IntervalValue))

	if estimatePeriods >= MAX_LIST_SIZE {
		return []common.WindowResult{}, fmt.Errorf("maximum size is %v - please select a smaller range", MAX_LIST_SIZE)
	}

	list := make([]common.WindowResult, 0, estimatePeriods*2)
	iterator, err := w.Iterate(direction, maxAttempt)
	if err != nil {
		return []common.WindowResult{}, err
	}

	for value := range iterator {
		if w.isWindowBetweenTimes(value, startTime, endTime) {
			list = append(list, value)
		}
	}

	return list, nil
}
