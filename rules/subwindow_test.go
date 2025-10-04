package rules

import (
	"testing"
	"time"

	"github.com/copito/coptime/interval"
)

func TestGenerateSubWindowsForRuleTypeBidirectional(t *testing.T) {
	base := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	rule := Rule{
		IntervalType: RuleTypeInclusion,
		MonthDays:    []uint32{31},
	}

	tests := []struct {
		name      string
		direction interval.Direction
		start     time.Time
		end       time.Time
	}{
		{
			name:      "forward interval",
			direction: interval.DirectionForward,
			start:     base,
			end:       base.AddDate(0, 1, 0),
		},
		{
			name:      "backward interval",
			direction: interval.DirectionBackward,
			start:     base.AddDate(0, 1, 0),
			end:       base,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			windows := GenerateSubWindowsForRuleType(tt.direction, interval.FrequencyMonth, tt.start, tt.end, []Rule{rule}, time.UTC)

			if len(windows) != 1 {
				t.Fatalf("expected 1 window, got %d", len(windows))
			}

			expectedStart := time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC)
			if !windows[0].Start.Equal(expectedStart) {
				t.Fatalf("expected window to start at %s, got %s", expectedStart, windows[0].Start)
			}

			expectedEnd := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)
			if !windows[0].End.Equal(expectedEnd) {
				t.Fatalf("expected window to end at %s, got %s", expectedEnd, windows[0].End)
			}
		})
	}
}

func TestGenerateSubWindowsForRuleAppliesFilter(t *testing.T) {
	base := time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC)
	rule := Rule{
		IntervalType: RuleTypeInclusion,
		MonthDays:    []uint32{28, 29, 30, 31},
		Filter:       NewOccurrenceFilter([]int{-1}),
	}

	windows := GenerateSubWindowsForRule(interval.DirectionForward, interval.FrequencyMonth, base, base.AddDate(0, 1, 0), rule, time.UTC)
	if len(windows) != 1 {
		t.Fatalf("expected filter to reduce to one window, got %d", len(windows))
	}

	expectedStart := time.Date(2024, time.January, 31, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC)

	if !windows[0].Start.Equal(expectedStart) {
		t.Fatalf("expected filtered start %s, got %s", expectedStart, windows[0].Start)
	}

	if !windows[0].End.Equal(expectedEnd) {
		t.Fatalf("expected filtered end %s, got %s", expectedEnd, windows[0].End)
	}
}
