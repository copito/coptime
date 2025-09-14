package interval

import (
	"fmt"
	"testing"
	"time"

	"github.com/copito/coptime/helper"
	"github.com/stretchr/testify/assert"
)

func TestIntervalerNew(t *testing.T) {
	now := time.Now()
	testTable := []struct {
		name     string
		opt      IntervalOption
		expected Intervaler
	}{
		{
			name: "no value",
			opt:  IntervalOption{},
			expected: Intervaler{
				opt: IntervalOption{
					AnchorDate:    time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)),
					EndDate:       nil,
					FrequencyUnit: FrequencyNanoSecond,
					IntervalValue: 1,
				},
			},
		},
		{
			name: "default values",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     nil,
				EndDate:       nil,
				FrequencyUnit: FrequencyDay,
				IntervalValue: 1,
			},
			expected: Intervaler{
				opt: IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       nil,
					FrequencyUnit: FrequencyDay,
					IntervalValue: 1,
				},
			},
		},
		{
			name: "default values with zero interval",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     nil,
				EndDate:       nil,
				FrequencyUnit: FrequencyDay,
				IntervalValue: 0,
			},
			expected: Intervaler{
				opt: IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       nil,
					FrequencyUnit: FrequencyDay,
					IntervalValue: 1,
				},
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.opt)
			assert.Equal(t, tt.expected.opt.AnchorDate, i.opt.AnchorDate, "checking anchor date")
			assert.Equal(t, tt.expected.opt.StartDate, i.opt.StartDate, "checking start date")
			assert.Equal(t, tt.expected.opt.EndDate, i.opt.EndDate, "checking end date")
			assert.Equal(t, tt.expected.opt.FrequencyUnit, i.opt.FrequencyUnit, "checking frequency unit")
			assert.Equal(t, tt.expected.opt.IntervalValue, i.opt.IntervalValue, "checking interval value")
		})
	}
}

func TestCalculateNext(t *testing.T) {
	testTable := []struct {
		name         string
		opt          IntervalOption
		previousTime time.Time
		direction    Direction
		expected     time.Time
	}{
		{
			name: "daily 1 with 2025-01-01 anchor with previous 2025-01-02",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyDay,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "daily 1 with 2025-01-01 anchor with previous 2025-01-31",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyDay,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "hourly 1 with 2025-01-01 anchor with previous 2025-01-01 10:00:00",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyHour,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 1, 11, 0, 0, 0, time.UTC),
		},
		{
			name: "minutely 1 with 2025-01-01 anchor with previous 2025-01-01 00:00:00",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyMinute,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 1, 0, 1, 0, 0, time.UTC),
		},
		{
			name: "weekly 1 with 2025-01-01 anchor with previous 2025-01-01",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyWeek,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 8, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly 1 with 2025-01-02 anchor with previous 2025-01-04",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyWeek,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly 1 with 2025-01-02 anchor with previous 2025-01-05",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyWeek,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly 1 with 2025-01-02 anchor with previous 2025-01-09",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyWeek,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 1, 9, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 1, 16, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "month 1 with 2025-01-02 anchor with previous 2025-02-02",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyMonth,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "month 1 with 2025-01-02 anchor with previous 2025-02-10",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyMonth,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 2, 10, 0, 0, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "weekly 1 with 2025-09-08 20:10:00 anchor with previous 2025-01-01",
			opt: IntervalOption{
				AnchorDate:    time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC),
				StartDate:     helper.ToPointer(time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC)),
				EndDate:       nil,
				FrequencyUnit: FrequencyWeek,
				IntervalValue: 1,
			},
			previousTime: time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC),
			direction:    DirectionForward,
			expected:     time.Date(2025, 9, 15, 20, 10, 0, 0, time.UTC),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			i := New(tt.opt)
			nextTime := i.calculateNext(tt.previousTime, tt.direction)
			assert.Equal(t, tt.expected, nextTime, fmt.Sprintf("expected: %s and got %s", tt.expected, nextTime))
		})
	}
}
