package window

func defaultMaxAttempts(value *int32) int32 {
	if value == nil {
		return DEFAULT_MAX_ATTEMPTS
	}

	if *value <= 0 {
		return DEFAULT_MAX_ATTEMPTS
	}

	return *value
}
