package window

import (
	"errors"
	"math"
	"time"

	"github.com/copito/coptime/interval"
	"github.com/copito/coptime/rules"
	"github.com/teambition/rrule-go"
)

func convertRRULEtoWindowOption(ruleString string) (*WindowOption, error) {
	rr, err := rrule.StrToRRule(ruleString)
	if err != nil {
		return nil, err
	}

	freq, err := MapRRULEFrequencyToFrequency(rr.Options.Freq)
	if err != nil {
		return nil, err
	}

	intervalValue := uint32(rr.Options.Interval)
	if intervalValue == 0 {
		intervalValue = 1
	}

	if rr.Options.Dtstart.IsZero() {
		rr.Options.Dtstart = time.Now().UTC()
	}

	// Add rules to bridge the gap between rrule and interval
	var rulesList []rules.Rule

	// Weekstart reference
	anchorDate := rr.Options.Dtstart
	startDate := rr.Options.Dtstart
	if freq == interval.FrequencyWeek && rr.Options.Wkst != rrule.MO && len(rr.Options.Byweekday) != 0 {
		// move anchor date to first weekday that matches
		wkst := time.Monday
		wkst = MapRRuleWeekToWeekday(rr.Options.Wkst)
		// If different than the anchor day, then move to the next weekday that matches
		if rr.Options.Dtstart.Weekday() != wkst {
			missingDaysTillWeekstart := time.Hour * 24 * 1 // TODO: determine the 1 value to find the next weekstart
			startDate = startDate.Add(missingDaysTillWeekstart)
		}

	}

	var endDate *time.Time
	if !rr.OrigOptions.Until.IsZero() {
		endDate = &rr.Options.Until
	}

	// Handling count given this is a post param on between and all
	var sizeValue float64
	if rr.OrigOptions.Count > 0 {
		sizeValue = float64(int32(rr.OrigOptions.Count))
	} else {
		sizeValue = math.Inf(1)
	}

	// Rules specific logic
	rule := rules.Rule{IntervalType: rules.RuleTypeInclusion}
	hasRule := false

	if len(rr.OrigOptions.Byweekday) != 0 {
		var wkDays []time.Weekday
		for _, wk := range rr.OrigOptions.Byweekday {
			wkDays = append(wkDays, MapRRuleWeekToWeekday(wk))
		}
		rule.DayOfWeeks = wkDays
		hasRule = true
	}

	if len(rr.OrigOptions.Bymonth) != 0 {
		var months []time.Month
		for _, m := range rr.OrigOptions.Bymonth {
			months = append(months, (time.Month)(m))
		}
		rule.Months = months
		hasRule = true
	}

	if len(rr.OrigOptions.Bymonthday) != 0 {
		var monthDays []uint32
		for _, day := range rr.OrigOptions.Bymonthday {
			monthDays = append(monthDays, (uint32)(day))
		}
		rule.MonthDays = monthDays
		hasRule = true
	}

	var occurrenceIndexes []int
	if len(rr.OrigOptions.Bysetpos) != 0 {
		occurrenceIndexes = append(occurrenceIndexes, rr.OrigOptions.Bysetpos...)
		hasRule = true
	}

	// For time, we can only represent a single time or a single range.
	// RRULE can have multiple values, e.g., BYHOUR=8,18.
	// The current TimeRange struct does not support this.
	// For now, we will only take the first value for each.
	timeRef := rules.TimeReference{}
	hasTimeRef := false
	if len(rr.OrigOptions.Byhour) > 0 {
		timeRef.Hour = rr.OrigOptions.Byhour[0]
		hasTimeRef = true
	}
	if len(rr.OrigOptions.Byminute) > 0 {
		timeRef.Minute = rr.OrigOptions.Byminute[0]
		hasTimeRef = true
	}
	if len(rr.OrigOptions.Bysecond) > 0 {
		timeRef.Second = rr.OrigOptions.Bysecond[0]
		hasTimeRef = true
	}

	if hasTimeRef {
		rule.TimeRange = &rules.TimeRange{
			StartTimeReference: timeRef,
			// For now, we assume a single point in time, so start and end are the same.
			EndTimeReference: timeRef,
		}
		hasRule = true
	}

	if hasRule {
		if len(occurrenceIndexes) > 0 {
			rule.Filter = rules.NewOccurrenceFilter(occurrenceIndexes)
		}
		rulesList = append(rulesList, rule)
	}

	return &WindowOption{
		IntervalOption: interval.IntervalOption{
			AnchorDate:    anchorDate,
			StartDate:     &startDate,
			EndDate:       endDate,
			Size:          &sizeValue,
			FrequencyUnit: freq,
			IntervalValue: intervalValue,
		},
		Rules: rulesList,
	}, nil
}

func convertCrontoIntervaler(cronString string) (*WindowOption, error) {
	return nil, errors.New("not implemented")
}
