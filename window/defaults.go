package window

import "github.com/copito/coptime/interval"

func defaultMaxAttempts(value *int32) int32 {
	if value == nil {
		return DEFAULT_MAX_ATTEMPTS
	}

	if *value <= 0 {
		return DEFAULT_MAX_ATTEMPTS
	}

	return *value
}

func adjustFrequencyUnitForRuleEvaluation(unit interval.Frequency) interval.Frequency {
	// When the frequency is less than a day, we evaluate rules on a daily basis
	if unit < interval.FrequencyDay {
		return interval.FrequencyHour
	}
	return interval.FrequencyDay
}
