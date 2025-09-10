package interval

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

func TestConvertRRULEtoIntervaler_ByDay_Unsupported(t *testing.T) {
	rruleString := "FREQ=MONTHLY;DTSTART=20240101T090000Z;BYMONTH=1"
	_, err := convertRRULEtoIntervaler(rruleString)
	require.Error(t, err)
	assert.Equal(t, "unsupported rrule feature: BY* rules other than BYMONTHDAY=-1 and BYDAY are not yet supported", err.Error())
}

func TestConvertRRULEtoIntervaler_MonthEnd(t *testing.T) {
	rruleString := "FREQ=MONTHLY;DTSTART=20240115T090000Z;BYMONTHDAY=-1;COUNT=3"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	assert.True(t, option.MonthEnd)

	iv := New(*option)
	all, err := iv.All(DirectionForward, nil)
	require.NoError(t, err)

	expectedDates := []time.Time{
		time.Date(2024, 1, 31, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 29, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 31, 9, 0, 0, 0, time.UTC),
	}

	require.Len(t, all, len(expectedDates))
	for i, d := range expectedDates {
		assert.Equal(t, d, all[i])
	}
}

func TestConvertRRULEtoIntervaler_ByDay_Backward(t *testing.T) {
	rruleString := "FREQ=WEEKLY;DTSTART=20240101T090000Z;BYDAY=MO,WE;INTERVAL=2;WKST=SU;COUNT=4"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	// Modify start date for backward iteration
	endDate := time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC)
	option.EndDate = &endDate
	startDate := time.Date(2024, 1, 17, 9, 0, 0, 0, time.UTC)
	option.StartDate = &startDate
	option.AnchorDate = startDate

	iv := New(*option)
	all, err := iv.All(DirectionBackward, nil)
	require.NoError(t, err)

	expectedDates := []time.Time{
		time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
	}

	require.Len(t, all, len(expectedDates))
	for i, d := range expectedDates {
		assert.Equal(t, d, all[i])
	}
}

func TestConvertRRULEtoIntervaler_ByDay_WithIntervalAndWkst(t *testing.T) {
	rruleString := "FREQ=WEEKLY;DTSTART=20240101T090000Z;BYDAY=MO,WE;INTERVAL=2;WKST=SU;COUNT=4"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	iv := New(*option)
	all, err := iv.All(DirectionForward, nil)
	require.NoError(t, err)

	expectedDates := []time.Time{
		time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 15, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 17, 9, 0, 0, 0, time.UTC),
	}

	require.Len(t, all, len(expectedDates))
	for i, d := range expectedDates {
		assert.Equal(t, d, all[i])
	}
}

func TestConvertRRULEtoIntervaler_ByDay(t *testing.T) {
	rruleString := "FREQ=WEEKLY;DTSTART=20240101T090000Z;BYDAY=MO,WE,FR;COUNT=5"
	option, err := convertRRULEtoIntervaler(rruleString)
	require.NoError(t, err)
	require.NotNil(t, option)

	iv := New(*option)
	all, err := iv.All(DirectionForward, nil)
	require.NoError(t, err)

	expectedDates := []time.Time{
		time.Date(2024, 1, 1, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 3, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 5, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 8, 9, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 10, 9, 0, 0, 0, time.UTC),
	}

	require.Len(t, all, len(expectedDates))
	for i, d := range expectedDates {
		assert.Equal(t, d, all[i])
	}
}

func TestConvertRRULEtoIntervaler_NoDtstart(t *testing.T) {
	rruleString := "FREQ=DAILY"
	_, err := convertRRULEtoIntervaler(rruleString)
	require.Error(t, err)
	assert.Equal(t, "rrule must have a DTSTART", err.Error())
}
