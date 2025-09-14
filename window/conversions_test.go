package window

import (
	"math"
	"testing"
	"time"

	"github.com/copito/coptime/helper"
	"github.com/copito/coptime/interval"
	rules "github.com/copito/coptime/rules"
	"github.com/stretchr/testify/assert"
	"github.com/teambition/rrule-go"
)

func TestRRULEConversion(t *testing.T) {
	// chicago, err := time.LoadLocation("America/Chicago")
	// if err != nil {
	// 	t.Fatal(err)
	// }

	now := time.Now()
	inf := math.Inf(1)

	testTable := []struct {
		name                 string
		rruleString          string
		expectedWindowOption WindowOption
		expectedError        bool
	}{
		{
			name:        "Invalid RRULE",
			rruleString: "Lemon Pie",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencySecond,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: true,
		},
		{
			name:        "Simple Second Interval (1)",
			rruleString: "FREQ=SECONDLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencySecond,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Minutely Interval (1)",
			rruleString: "FREQ=MINUTELY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMinute,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Hourly Interval (1)",
			rruleString: "FREQ=HOURLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyHour,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Daily Interval (1)",
			rruleString: "FREQ=DAILY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Weekly Interval (1)",
			rruleString: "FREQ=WEEKLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyWeek,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Monthly Interval (1)",
			rruleString: "FREQ=MONTHLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMonth,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Yearly Interval (1)",
			rruleString: "FREQ=YEARLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyYear,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Yearly Interval (1)",
			rruleString: "FREQ=YEARLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyYear,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Daily Interval (2)",
			rruleString: "FREQ=DAILY;INTERVAL=2",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: uint32(2),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Daily Interval (12)",
			rruleString: "FREQ=DAILY;INTERVAL=12",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: uint32(12),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Minutely Interval (18)",
			rruleString: "FREQ=MINUTELY;INTERVAL=18",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    now,
					StartDate:     helper.ToPointer(now),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMinute,
					IntervalValue: uint32(18),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Daily Interval (1) With DTstart (20250101000000)",
			rruleString: "DTSTART=20240101T090000Z;FREQ=DAILY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Monthly Interval (1) With DTstart (20250101000000)",
			rruleString: "DTSTART=20240101T090000Z;FREQ=MONTHLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMonth,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Monthly Interval (1) With DTstart (20250908201000)",
			rruleString: "DTSTART=20250908T201000Z;FREQ=MONTHLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC)),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMonth,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Monthly Interval (1) With DTstart (20250908201000)",
			rruleString: "DTSTART=20250908T201000Z;FREQ=MONTHLY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC)),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyMonth,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		{
			name:        "Simple Daily Interval (1) With DTstart (20250908201000)",
			rruleString: "DTSTART=20250908T201000Z;FREQ=DAILY;INTERVAL=1",
			expectedWindowOption: WindowOption{
				IntervalOption: interval.IntervalOption{
					AnchorDate:    time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC),
					StartDate:     helper.ToPointer(time.Date(2025, 9, 8, 20, 10, 0, 0, time.UTC)),
					EndDate:       nil,
					Size:          &inf,
					FrequencyUnit: interval.FrequencyDay,
					IntervalValue: uint32(1),
				},
				Rules: []rules.Rule{},
			},
			expectedError: false,
		},
		// {
		// 	name:        "Simple Daily Interval (1) With DTstart (20250101000000)",
		// 	rruleString: "DTSTART=20240101T090000Z;FREQ=DAILY;INTERVAL=1",
		// 	expectedWindowOption: WindowOption{
		// 		IntervalOption: interval.IntervalOption{
		// 			AnchorDate:    time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		// 			StartDate:     helper.ToPointer(time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)),
		// 			EndDate:       nil,
		// 			Size:          &inf,
		// 			FrequencyUnit: interval.FrequencyDay,
		// 			IntervalValue: uint32(1),
		// 		},
		// 		Rules: []rules.Rule{},
		// 	},
		// 	expectedError: nil,
		// },
		// {
		// 	name: "Simple Daily Interval (1) MondayStart",
		// 	// rruleString: "DTSTART;TZID=America/Chicago:20250908T000000 RRULE:FREQ=DAILY;INTERVAL=1",
		// 	rruleString: "DTSTART:20250908T000000Z RRULE:FREQ=DAILY;WKST=MO",
		// 	expectedWindowOption: WindowOption{
		// 		IntervalOption: interval.IntervalOption{
		// 			AnchorDate:    time.Date(2025, 9, 8, 0, 0, 0, 0, chicago),
		// 			StartDate:     helper.ToPointer(time.Date(2025, 9, 8, 0, 0, 0, 0, chicago)),
		// 			EndDate:       nil,
		// 			Size:          nil,
		// 			FrequencyUnit: interval.FrequencyDay,
		// 			IntervalValue: uint32(1),
		// 		},
		// 		Rules: []rules.Rule{},
		// 	},
		// 	expectedError: nil,
		// },
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			hasDTStart := false
			rz, err := rrule.StrToRRule(tt.rruleString)
			if err != nil {
				if tt.expectedError {
					assert.NotNil(t, err, "expecting error from parsing rrule")
				} else {
					t.Fatal("expecting no error from parsing rrule")
				}
			} else {
				if !rz.OrigOptions.Dtstart.IsZero() {
					hasDTStart = true
				}
			}

			calculatedActual, err := convertRRULEtoWindowOption(tt.rruleString)
			if tt.expectedError {
				assert.NotNil(t, err, "should be error evaluating rrule")
				return
			} else {
				assert.Nil(t, err, "should be nil for error evaluating rrule")
			}

			// Assert that all the attributes and rules are the same
			if hasDTStart {
				assert.True(t, tt.expectedWindowOption.AnchorDate.Equal(calculatedActual.AnchorDate), "anchor date should be equal")
				if tt.expectedWindowOption.StartDate != nil {
					assert.NotNil(t, calculatedActual.StartDate, "start date is not nil")
					assert.True(t, tt.expectedWindowOption.StartDate.Equal(*calculatedActual.StartDate), "start date should be equal")
				} else {
					assert.Nil(t, calculatedActual.StartDate, "start date should be nil")
				}

				if tt.expectedWindowOption.EndDate != nil {
					assert.NotNil(t, calculatedActual.EndDate, "end date is not nil")
					assert.True(t, tt.expectedWindowOption.EndDate.Equal(*calculatedActual.EndDate), "end date should be equal")
				} else {
					assert.Nil(t, calculatedActual.EndDate, "end date should be nil")
				}
			} else {
				assert.True(t, tt.expectedWindowOption.AnchorDate.After(now.Add(-1*time.Hour)), "anchor date should be equal")
				assert.True(t, tt.expectedWindowOption.AnchorDate.Before(now.Add(1*time.Hour)), "anchor date should be equal")
				if tt.expectedWindowOption.StartDate != nil {
					assert.NotNil(t, calculatedActual.StartDate, "start date is not nil")
					assert.True(t, tt.expectedWindowOption.StartDate.After(now.Add(-1*time.Hour)), "start date should be equal")
					assert.True(t, tt.expectedWindowOption.StartDate.Before(now.Add(1*time.Hour)), "start date should be equal")
				} else {
					assert.Nil(t, calculatedActual.StartDate, "start date should be nil")
				}

				// if tt.expectedWindowOption.EndDate != nil {
				// 	assert.NotNil(t, calculatedActual.EndDate, "end date is not nil")
				// 	assert.True(t, tt.expectedWindowOption.EndDate.After(*calculatedActual.EndDate), "end date should be equal")
				// } else {
				// 	assert.Nil(t, calculatedActual.EndDate, "end date should be nil")
				// }
			}

			if tt.expectedWindowOption.Size != nil {
				assert.NotNil(t, calculatedActual.Size, "size is not nil")
				assert.Equal(t, tt.expectedWindowOption.Size, calculatedActual.Size, "size should be equal")
			} else {
				assert.Nil(t, calculatedActual.Size, "size is nil")
			}

			assert.Equal(t, tt.expectedWindowOption.FrequencyUnit, calculatedActual.FrequencyUnit, "frequency unit should be equal")
			assert.Equal(t, tt.expectedWindowOption.IntervalValue, calculatedActual.IntervalValue, "interval should be equal")

			assert.Equal(t, len(tt.expectedWindowOption.Rules), len(calculatedActual.Rules), "rule count must be the same")
			for i, rule := range tt.expectedWindowOption.Rules {
				assert.ElementsMatch(t, rule.DayOfWeeks, calculatedActual.Rules[i].DayOfWeeks, "day of week matches")
				assert.ElementsMatch(t, rule.Years, calculatedActual.Rules[i].Years, "years matches")
				assert.ElementsMatch(t, rule.Months, calculatedActual.Rules[i].Months, "months matches")
				assert.ElementsMatch(t, rule.MonthDays, calculatedActual.Rules[i].MonthDays, "month days matches")

				// Time Range comparison
				if rule.TimeRange != nil {
					assert.NotNil(t, rule.TimeRange, "time range is not nil")
					assert.Equal(t, rule.TimeRange.StartTimeReference.Hour, calculatedActual.Rules[i].TimeRange.StartTimeReference.Hour, "start time - hour ref should be equal")
					assert.Equal(t, rule.TimeRange.StartTimeReference.Minute, calculatedActual.Rules[i].TimeRange.StartTimeReference.Minute, "start time - minute ref should be equal")
					assert.Equal(t, rule.TimeRange.StartTimeReference.Second, calculatedActual.Rules[i].TimeRange.StartTimeReference.Second, "start time - second ref should be equal")
					assert.Equal(t, rule.TimeRange.StartTimeReference.MilliSecond, calculatedActual.Rules[i].TimeRange.StartTimeReference.MilliSecond, "start time - millisecond ref should be equal")
					assert.Equal(t, rule.TimeRange.StartTimeReference.NanoSecond, calculatedActual.Rules[i].TimeRange.StartTimeReference.NanoSecond, "start time - nanosecond ref should be equal")

					assert.Equal(t, rule.TimeRange.EndTimeReference.Hour, calculatedActual.Rules[i].TimeRange.EndTimeReference.Hour, "end time - hour ref should be equal")
					assert.Equal(t, rule.TimeRange.EndTimeReference.Minute, calculatedActual.Rules[i].TimeRange.EndTimeReference.Minute, "end time - minute ref should be equal")
					assert.Equal(t, rule.TimeRange.EndTimeReference.Second, calculatedActual.Rules[i].TimeRange.EndTimeReference.Second, "end time - second ref should be equal")
					assert.Equal(t, rule.TimeRange.EndTimeReference.MilliSecond, calculatedActual.Rules[i].TimeRange.EndTimeReference.MilliSecond, "end time - millisecond ref should be equal")
					assert.Equal(t, rule.TimeRange.EndTimeReference.NanoSecond, calculatedActual.Rules[i].TimeRange.EndTimeReference.NanoSecond, "end time - nanosecond ref should be equal")
				} else {
					assert.Nil(t, rule.TimeRange, "time range is nil")
				}

			}
		})
	}
}

// func TestHelperDebugConversionRRULE(t *testing.T) {
// 	rruleString := "DTSTART=20250908T201000Z;FREQ=WEEKLY;INTERVAL=1;WKST=MO;BYDAY=WE"
// 	calculatedActual, err := convertRRULEtoWindowOption(rruleString)
// 	if err != nil {
// 		t.Fatal("error")
// 	}

// 	fmt.Println(calculatedActual)
// 	ww := New(*calculatedActual)
// 	windows, _ := ww.Between(
// 		interval.DirectionForward,
// 		time.Date(2025, 9, 8, 0, 0, 0, 0, time.UTC),
// 		time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
// 		nil,
// 	)

// 	for _, value := range windows {
// 		fmt.Printf("value.Start: %s\n", value.Start)
// 		fmt.Printf("value.End: %s\n", value.End)
// 		fmt.Printf("value.SubWindow: %s\n", value.SubWindow)
// 	}
// }
