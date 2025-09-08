package window

import (
	"errors"
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

	// Create a include window with the entire range if no inclusion rules are defined
	if len(includes) == 0 {
		completeInclusion := rules.CreateCompleteInclusionRule(*w.opt.StartDate, *w.opt.EndDate)
		includes = append(includes, completeInclusion)
	}

	// Safe-guard: avoid infinite loop
	calculatedMaxAttempts := defaultMaxAttempts(maxAttempt)
	maxCounter := int32(0)

	return func(yield func(common.WindowResult) bool) {
		next, _ := iter.Pull(iterator)
		// defer stop()

		previousTime, ok := next()
		if !ok {
			return
		}

		for {
			// Stop if we reached the max attempts
			if maxCounter >= calculatedMaxAttempts {
				return // stop iteration
			}

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
			additiveSubWindows := rules.GenerateSubWindowsForInclusion(direction, w.opt.FrequencyUnit, previousTime, nextTime, includes, tz)

			// 2. Apply exclusion rules
			subtractiveSubWindows := rules.GenerateSubWindowsForExclusion(direction, w.opt.FrequencyUnit, previousTime, nextTime, excludes, tz)

			// 3. Combine additive and subtractive sub-windows
			subWindows = rules.SubtractSubwindowsFromAdditives(additiveSubWindows, subtractiveSubWindows)

			if len(subWindows) == 0 {
				// Do not add that sessions if no coditions are met
				// (e.g. session window is between 2024-01-01 to 2024-01-02 and the rule are only for 2025-12-25)

				// Move to the next element
				previousTime = nextTime
				maxCounter += 1
				continue // skip to the next iteration if no sub-windows are created
			}

			window.SubWindow = subWindows

			// yield the window
			ok = yield(window)
			if !ok {
				return
			}

			maxCounter = 0
			previousTime = nextTime
		}
	}, nil
}

func (w *Windower) All(direction interval.Direction, maxAttempt *int32) ([]common.WindowResult, error) {
	// panic("unimplemented")
	return nil, nil
}

func (w *Windower) Between(direction interval.Direction, startTime time.Time, endTime time.Time, maxAttempt *int32) ([]common.WindowResult, error) {
	// panic("unimplemented")
	return nil, nil
}
