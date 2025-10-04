package window

import (
	"testing"
	"time"

	"github.com/copito/coptime/common"
	"github.com/copito/coptime/interval"
	rulespkg "github.com/copito/coptime/rules"
)

func TestBuildRuleBindingsFiltersByType(t *testing.T) {
	opt := WindowOption{
		Rules: []rulespkg.Rule{
			{IntervalType: rulespkg.RuleTypeInclusion},
			{IntervalType: rulespkg.RuleTypeExclusion},
		},
	}

	includes := buildRuleBindings(opt, rulespkg.RuleTypeInclusion)
	if len(includes) != 1 {
		t.Fatalf("expected 1 inclusion binding, got %d", len(includes))
	}

	excludes := buildRuleBindings(opt, rulespkg.RuleTypeExclusion)
	if len(excludes) != 1 {
		t.Fatalf("expected 1 exclusion binding, got %d", len(excludes))
	}
}

func TestGenerateSubWindowsForBindingsAppliesRuleFilters(t *testing.T) {
	start := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)

	binding := ruleBinding{
		rule: rulespkg.Rule{
			IntervalType: rulespkg.RuleTypeInclusion,
			MonthDays:    []uint32{28, 29, 30, 31},
			Filter:       rulespkg.NewOccurrenceFilter([]int{-1}),
		},
	}

	windows := generateSubWindowsForBindings(interval.DirectionForward, interval.FrequencyMonth, start, end, []ruleBinding{binding}, time.UTC)
	if len(windows) != 1 {
		t.Fatalf("expected a single filtered window, got %d", len(windows))
	}

	expectedStart := time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC)
	expectedEnd := expectedStart.AddDate(0, 0, 1)

	if !windows[0].Start.Equal(expectedStart) {
		t.Fatalf("expected start %s, got %s", expectedStart, windows[0].Start)
	}

	if !windows[0].End.Equal(expectedEnd) {
		t.Fatalf("expected end %s, got %s", expectedEnd, windows[0].End)
	}

	// verify filters do not mutate the original slice and still produce sorted output
	raw := []common.SubWindowResult{
		{Start: start.AddDate(0, 0, 27), End: start.AddDate(0, 0, 28)},
		{Start: start.AddDate(0, 0, 28), End: start.AddDate(0, 0, 29)},
		{Start: start.AddDate(0, 0, 29), End: start.AddDate(0, 0, 30)},
		{Start: start.AddDate(0, 0, 30), End: start.AddDate(0, 0, 31)},
	}

	filtered := binding.rule.Filter(raw)
	if len(filtered) != 1 {
		t.Fatalf("expected filtered result length 1, got %d", len(filtered))
	}

	if !filtered[0].Start.Equal(expectedStart) {
		t.Fatalf("expected filtered start %s, got %s", expectedStart, filtered[0].Start)
	}
}
