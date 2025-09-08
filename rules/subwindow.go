package rules

import (
	"slices"
	"sort"
	"time"

	"github.com/copito/coptime/common"
	"github.com/copito/coptime/interval"
)

func timesToWindows(times []time.Time) []common.SubWindowResult {
	if len(times) == 0 {
		return []common.SubWindowResult{}
	}

	if len(times) == 1 {
		return []common.SubWindowResult{
			{
				Start: times[0],
				End:   times[0],
			},
		}
	}

	var windows []common.SubWindowResult
	for i := 0; i < len(times)-1; i++ {
		windows = append(windows, common.SubWindowResult{
			Start: times[i],
			End:   times[i+1],
		})
	}
	return windows
}

func matchesRuleToWindow(r Rule, start time.Time, end time.Time) *common.SubWindowResult {
	// loc := getTimezone(r.Timezone)
	// localTime := startTime.In(loc)

	// Match year
	// if no years are specified, match all years
	if len(r.Years) > 0 {
		startYear := start.Year()

		if !slices.Contains(r.Years, uint32(startYear)) {
			return nil
		}
	}

	// Match month
	// If no months are specified, match all months
	if len(r.Months) > 0 {
		startMonth := start.Month()

		if !slices.Contains(r.Months, startMonth) {
			return nil
		}
	}

	// Match day of month
	// If no days of month are specified, match all days of month
	if len(r.MonthDays) > 0 {
		startDay := uint32(start.Day())
		if !slices.Contains(r.MonthDays, startDay) {
			return nil
		}
	}

	// Match day of week
	// If no days of week are specified, match all days of week
	if len(r.DayOfWeeks) > 0 {
		startDayOfWeek := start.Weekday()
		if !slices.Contains(r.DayOfWeeks, startDayOfWeek) {
			return nil
		}
	}

	// Match time range
	var timeRangeStart time.Time
	var timeRangeEnd time.Time
	if r.TimeRange != nil {
		timeRangeStart = time.Date(start.Year(), start.Month(), start.Day(),
			r.TimeRange.StartTimeReference.Hour,
			r.TimeRange.StartTimeReference.Minute,
			r.TimeRange.StartTimeReference.Second,
			r.TimeRange.StartTimeReference.NanoSecond+(r.TimeRange.StartTimeReference.MilliSecond*1_000_000),
			start.Location(),
		)

		// Special case: if end time is 00:00:00, it means the end of the day
		if r.TimeRange.EndTimeReference.Hour == 0 &&
			r.TimeRange.EndTimeReference.Minute == 0 &&
			r.TimeRange.EndTimeReference.Second == 0 &&
			r.TimeRange.EndTimeReference.MilliSecond == 0 &&
			r.TimeRange.EndTimeReference.NanoSecond == 0 {
			timeRangeEnd = time.Date(start.Year(), start.Month(), start.Day()+1, 0, 0, 0, 0, start.Location())
		} else {
			timeRangeEnd = time.Date(start.Year(), start.Month(), start.Day(),
				r.TimeRange.EndTimeReference.Hour,
				r.TimeRange.EndTimeReference.Minute,
				r.TimeRange.EndTimeReference.Second,
				r.TimeRange.EndTimeReference.NanoSecond+(r.TimeRange.EndTimeReference.MilliSecond*1_000_000),
				start.Location(),
			)
		}

		// Check if not inside the time range
		if start.After(timeRangeEnd) || end.Before(timeRangeStart) {
			return nil
		}

		// Adjust start and end to be within the time range
		if start.Before(timeRangeStart) {
			start = timeRangeStart
		}
		if end.After(timeRangeEnd) {
			end = timeRangeEnd
		}
	} else {
		// If no time range is specified, the whole window is valid
		timeRangeStart = start
		timeRangeEnd = end
	}

	return &common.SubWindowResult{
		Start: timeRangeStart,
		End:   timeRangeEnd,
	}
}

func adjustFrequencyUnitForRuleEvaluation(unit interval.Frequency) interval.Frequency {
	// When the frequency is less than a day, we evaluate rules on a daily basis
	if unit < interval.FrequencyDay {
		return interval.FrequencyHour
	}
	return interval.FrequencyDay
}

func GenerateSubWindowsForRuleType(direction interval.Direction, unit interval.Frequency, start time.Time, end time.Time, rules []Rule, loc *time.Location) []common.SubWindowResult {
	evaluatedFrequencyUnit := adjustFrequencyUnitForRuleEvaluation(unit)

	// Generate a day window for each subwindow
	iv := interval.New(interval.IntervalOption{
		AnchorDate:    start,
		StartDate:     &start,
		EndDate:       &end,
		FrequencyUnit: evaluatedFrequencyUnit,
		IntervalValue: 1,
	})

	times, err := iv.Between(direction, start, end, nil)
	if err != nil {
		return []common.SubWindowResult{}
	}

	evaluationWindow := timesToWindows(times)

	subWindows := make([]common.SubWindowResult, 0, len(evaluationWindow)*len(rules))
	for _, w := range evaluationWindow {
		for _, r := range rules {
			subWindowFilter := matchesRuleToWindow(r, w.Start, w.End)

			if subWindowFilter == nil {
				continue // skip rule if the rule does not match
			}

			if subWindowFilter.Start.IsZero() || subWindowFilter.End.IsZero() {
				continue // skip rule if the rule does not match
			}

			if subWindowFilter.Start.Equal(subWindowFilter.End) || subWindowFilter.Start.After(subWindowFilter.End) {
				continue // skip invalid windows
			}

			subWindows = append(subWindows, *subWindowFilter)
		}
	}

	subWindows = mergeSubWindows(subWindows)
	return subWindows
}

func mergeSubWindows(subWindows []common.SubWindowResult) []common.SubWindowResult {
	if len(subWindows) == 0 {
		return subWindows
	}

	// Sort sub-windows by start time
	sort.Slice(subWindows, func(i, j int) bool {
		return subWindows[i].Start.Before(subWindows[j].Start)
	})

	merged := []common.SubWindowResult{subWindows[0]}

	for i := 1; i < len(subWindows); i++ {
		last := &merged[len(merged)-1]
		current := subWindows[i]

		// if overlapping or touching, merge them
		if !current.Start.After(last.End) {
			if current.End.After(last.End) {
				last.End = current.End
			}
		} else {
			merged = append(merged, current)
		}
	}

	return merged
}

func SubtractSubwindowsFromAdditives(adds []common.SubWindowResult, subs []common.SubWindowResult) []common.SubWindowResult {
	if len(adds) == 0 {
		return []common.SubWindowResult{}
	}

	if len(subs) == 0 {
		return adds
	}

	var result []common.SubWindowResult
	for _, add := range adds {
		// Start with the current additivie window
		remaining := []common.SubWindowResult{add}

		// For each subtractive window, remove its overlap from all remaining windows
		for _, sub := range subs {
			var temp []common.SubWindowResult
			for _, rem := range remaining {
				// Check if the subtractive window overlaps with the remaining window
				if rem.Start.Before(sub.End) && rem.End.After(sub.Start) {
					// There is an overlap, split the remaining window if necessary
					if rem.Start.Before(sub.Start) {
						temp = append(temp, common.SubWindowResult{
							Start: rem.Start,
							End:   sub.Start,
						})
					}
					// if sub.End.Before(rem.End) {
					// 	temp = append(temp, SubWindowResult{
					// 		Start: sub.End,
					// 		End:   rem.End,
					// 	})
					// }
					if rem.End.After(sub.End) {
						temp = append(temp, common.SubWindowResult{
							Start: sub.End,
							End:   rem.End,
						})
					}

				} else {
					// No overlap, keep the remaining window as is
					temp = append(temp, rem)
				}
			}
			remaining = temp
		}

		// Add all remaining parts of the additive window to the result
		result = append(result, remaining...)
	}
	return result
}

func CreateCompleteInclusionRule(start time.Time, end time.Time) Rule {
	startYear := uint32(start.Year())
	endYear := uint32(end.Year())

	minYear := min(startYear, endYear)
	maxYear := max(startYear, endYear)

	years := []uint32{}
	for y := minYear; y <= maxYear; y++ {
		years = append(years, y)
	}

	return Rule{
		IntervalType: RuleTypeInclusion,
		TimeRange: &TimeRange{
			StartTimeReference: TimeReference{
				Hour:        0,
				Minute:      0,
				Second:      0,
				MilliSecond: 0,
				NanoSecond:  0,
			},
			EndTimeReference: TimeReference{
				Hour:        0,
				Minute:      0,
				Second:      0,
				MilliSecond: 0,
				NanoSecond:  0,
			},
		},
		DayOfWeeks: []time.Weekday{time.Sunday, time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday},
		MonthDays:  []uint32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31},
		Months:     []time.Month{time.January, time.February, time.March, time.April, time.May, time.June, time.July, time.August, time.September, time.October, time.November, time.December},
		Years:      years,
	}
}
