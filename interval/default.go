package interval

import (
	"slices"
	"time"
)

func defaultAnchorTime(dt time.Time) time.Time {
	if dt.IsZero() {
		now := time.Now()
		newTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
		return newTime
	}

	return dt
}

func defaultStartTime(dt *time.Time, anchorTime time.Time) *time.Time {
	if dt == nil || dt.IsZero() {
		newTime := time.Date(anchorTime.Year(), anchorTime.Month(), anchorTime.Day(), anchorTime.Hour(), anchorTime.Minute(), anchorTime.Second(), anchorTime.Nanosecond(), anchorTime.Location())
		return &newTime
	}

	return dt
}

func defaultEndTime(dt *time.Time, direction Direction) time.Time {
	if (dt == nil || dt.IsZero()) && direction == DirectionForward {
		newTime := time.Date(99999, 12, 31, 23, 59, 59, 59, time.UTC)
		return newTime
	}

	if (dt == nil || dt.IsZero()) && direction == DirectionBackward {
		newTime := time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)
		return newTime
	}

	if dt == nil || dt.IsZero() {
		return time.Time{}
	}

	return *dt
}

func defaultFrequencyUnit(freq Frequency) Frequency {
	if !slices.Contains(possibleFrequencies, freq) {
		return FrequencyDay
	}
	return freq
}

func defaultIntervalValue(value uint32) uint32 {
	if value == 0 {
		return 1
	}

	return value
}

func defaultMaxAttempts(value *int32) int32 {
	if value == nil {
		return DEFAULT_MAX_ATTEMPTS
	}

	if *value <= 0 {
		return DEFAULT_MAX_ATTEMPTS
	}

	return *value
}
