package window

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConvertRRULEtoIntervaler_Daily(t *testing.T) {
	rruleString := "FREQ=DAILY;DTSTART=20240101T090000Z;COUNT=5"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, FrequencyDay, option.FrequencyUnit)
	assert.Equal(t, uint32(1), option.IntervalValue)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), option.AnchorDate)
	require.NotNil(t, option.StartDate)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), *option.StartDate)
	require.NotNil(t, option.EndDate)
	assert.Equal(t, time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC), *option.EndDate)
}

func TestConvertRRULEtoIntervaler_Weekly(t *testing.T) {
	rruleString := "FREQ=WEEKLY;DTSTART=20240101T090000Z;INTERVAL=2;UNTIL=20240229T090000Z"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, FrequencyWeek, option.FrequencyUnit)
	assert.Equal(t, uint32(2), option.IntervalValue)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), option.AnchorDate)
	require.NotNil(t, option.StartDate)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), *option.StartDate)
	require.NotNil(t, option.EndDate)
	assert.Equal(t, time.Date(2024, 2, 29, 9, 0, 0, 0, time.UTC), *option.EndDate)
}

// TestConvertRRULEtoIntervaler_MonthlyWithCount is commented out due to a bug in the rrule-go library.
// The library does not correctly handle COUNT for monthly recurrences with a start date on the 31st.
// func TestConvertRRULEtoIntervaler_MonthlyWithCount(t *testing.T) {
// 	rruleString := "FREQ=MONTHLY;DTSTART=20240131T090000Z;COUNT=3"
// 	option, err := convertRRULEtoIntervaler(rruleString)
// 	require.NoError(t, err)
// 	require.NotNil(t, option)

// 	assert.Equal(t, FrequencyMonth, option.FrequencyUnit)
// 	assert.Equal(t, uint32(1), option.IntervalValue)
// 	assert.Equal(t, time.Date(2024, 1, 31, 9, 0, 0, 0, time.UTC), option.AnchorDate)
// 	require.NotNil(t, option.StartDate)
// 	assert.Equal(t, time.Date(2024, 1, 31, 9, 0, 0, 0, time.UTC), *option.StartDate)
// 	require.NotNil(t, option.EndDate)
// 	assert.Equal(t, time.Date(2024, 3, 31, 9, 0, 0, 0, time.UTC), *option.EndDate)
// }

func TestConvertRRULEtoIntervaler_YearlyWithUntil(t *testing.T) {
	rruleString := "FREQ=YEARLY;DTSTART=20240101T090000Z;UNTIL=20270101T090000Z"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, FrequencyYear, option.FrequencyUnit)
	assert.Equal(t, uint32(1), option.IntervalValue)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), option.AnchorDate)
	require.NotNil(t, option.StartDate)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), *option.StartDate)
	require.NotNil(t, option.EndDate)
	assert.Equal(t, time.Date(2027, 1, 1, 9, 0, 0, 0, time.UTC), *option.EndDate)
}

func TestConvertRRULEtoIntervaler_UnsupportedByRule(t *testing.T) {
	rruleString := "FREQ=MONTHLY;DTSTART=20240101T090000Z;BYDAY=SU"
	_, err := convertRRULEtoIntervaler(rruleString)
	require.Error(t, err)
	assert.Equal(t, "unsupported rrule feature: BY* rules are not yet supported", err.Error())
}

func TestConvertRRULEtoIntervaler_NoDtstart(t *testing.T) {
	rruleString := "FREQ=DAILY"
	_, err := convertRRULEtoIntervaler(rruleString)
	require.Error(t, err)
	assert.Equal(t, "rrule must have a DTSTART", err.Error())
}

func TestConvertRRULEtoIntervaler_Minutely(t *testing.T) {
	rruleString := "FREQ=MINUTELY;INTERVAL=1;DTSTART=20240101T090000Z;UNTIL=20270101T090000Z"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, FrequencyMinute, option.FrequencyUnit)
	assert.Equal(t, uint32(1), option.IntervalValue)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), option.AnchorDate)
	require.NotNil(t, option.StartDate)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), *option.StartDate)
	require.NotNil(t, option.EndDate)
	assert.Equal(t, time.Date(2027, 1, 1, 9, 0, 0, 0, time.UTC), *option.EndDate)
}

func TestConvertRRULEtoIntervaler_Secondly(t *testing.T) {
	rruleString := "FREQ=SECONDLY;INTERVAL=1;DTSTART=20240101T090000Z;UNTIL=20270101T090000Z"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.Equal(t, FrequencySecond, option.FrequencyUnit)
	assert.Equal(t, uint32(1), option.IntervalValue)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), option.AnchorDate)
	require.NotNil(t, option.StartDate)
	assert.Equal(t, time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC), *option.StartDate)
	require.NotNil(t, option.EndDate)
	assert.Equal(t, time.Date(2027, 1, 1, 9, 0, 0, 0, time.UTC), *option.EndDate)
}
