package window

import (
	"testing"
	"time"

	"github.com/copito/coptime/common"
	"github.com/copito/coptime/helper"
	"github.com/copito/coptime/interval"
	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	type testCase struct {
		name           string
		windower       Windower
		direction      interval.Direction
		maxAttempt     *int32
		expectedResult []common.WindowResult
		expectedError  error
	}

	testTable := []testCase{
		{
			name: "Daily interval",
			windower: New(WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       helper.ToPointer(time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)),
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: 1,
				},
			}),
			direction: interval.DirectionForward,
			expectedResult: []common.WindowResult{
				{
					Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
		{
			name: "Hourly interval",
			windower: New(WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       helper.ToPointer(time.Date(2025, 1, 1, 2, 0, 0, 0, time.UTC)),
					FrequencyUnit: interval.FrequencyHour,
					IntervalValue: 1,
				},
			}),
			direction: interval.DirectionForward,
			expectedResult: []common.WindowResult{
				{
					Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 1, 2, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 1, 1, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 1, 2, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
		{
			name: "No end date",
			windower: New(WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: 1,
				},
			}),
			direction:  interval.DirectionForward,
			maxAttempt: helper.ToPointer(int32(2)),
			expectedResult: []common.WindowResult{
				{
					Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.windower.All(tc.direction, tc.maxAttempt)

			assert.Equal(t, tc.expectedError, err, "The error is not correct")
			assert.Equal(t, len(tc.expectedResult), len(result), "The length of the result is not correct")

			for i, expected := range tc.expectedResult {
				assert.Equal(t, expected.Start, result[i].Start, "The start time is not correct")
				assert.Equal(t, expected.End, result[i].End, "The end time is not correct")
				assert.Equal(t, len(expected.SubWindow), len(result[i].SubWindow), "The length of the subwindow is not correct")
				for j, expectedSubWindow := range expected.SubWindow {
					assert.Equal(t, expectedSubWindow.Start, result[i].SubWindow[j].Start, "The start time of the subwindow is not correct")
					assert.Equal(t, expectedSubWindow.End, result[i].SubWindow[j].End, "The end time of the subwindow is not correct")
				}
			}
		})
	}
}

func TestBetween(t *testing.T) {
	type testCase struct {
		name           string
		windower       Windower
		direction      interval.Direction
		startTime      time.Time
		endTime        time.Time
		maxAttempt     *int32
		expectedResult []common.WindowResult
		expectedError  error
	}

	testTable := []testCase{
		{
			name: "Daily interval, full range",
			windower: New(WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       helper.ToPointer(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)),
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: 1,
				},
			}),
			direction: interval.DirectionForward,
			startTime: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
			expectedResult: []common.WindowResult{
				{
					Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
		{
			name: "Daily interval, partial range",
			windower: New(WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)),
					EndDate:       helper.ToPointer(time.Date(2025, 1, 5, 0, 0, 0, 0, time.UTC)),
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: 1,
				},
			}),
			direction: interval.DirectionForward,
			startTime: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
			endTime:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
			expectedResult: []common.WindowResult{
				{
					Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 2, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC)},
					},
				},
				{
					Start: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC),
					End:   time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC),
					SubWindow: []common.SubWindowResult{
						{Start: time.Date(2025, 1, 3, 0, 0, 0, 0, time.UTC), End: time.Date(2025, 1, 4, 0, 0, 0, 0, time.UTC)},
					},
				},
			},
		},
	}

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			result, err := tc.windower.Between(tc.direction, tc.startTime, tc.endTime, tc.maxAttempt)

			assert.Equal(t, tc.expectedError, err, "The error is not correct")
			assert.Equal(t, len(tc.expectedResult), len(result), "The length of the result is not correct")

			for i, expected := range tc.expectedResult {
				assert.Equal(t, expected.Start, result[i].Start, "The start time is not correct")
				assert.Equal(t, expected.End, result[i].End, "The end time is not correct")
				assert.Equal(t, len(expected.SubWindow), len(result[i].SubWindow), "The length of the subwindow is not correct")
				for j, expectedSubWindow := range expected.SubWindow {
					assert.Equal(t, expectedSubWindow.Start, result[i].SubWindow[j].Start, "The start time of the subwindow is not correct")
					assert.Equal(t, expectedSubWindow.End, result[i].SubWindow[j].End, "The end time of the subwindow is not correct")
				}
			}
		})
	}
}
