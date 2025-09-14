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
	if len(rr.OrigOptions.Byweekday) != 0 {
		var wkDays []time.Weekday
		for _, wk := range rr.OrigOptions.Byweekday {
			wkDays = append(wkDays, MapRRuleWeekToWeekday(wk))
		}

		rulesList = append(rulesList, rules.Rule{
			IntervalType: rules.RuleTypeInclusion,
			DayOfWeeks:   wkDays,
		})
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
