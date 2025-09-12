package interval

import "time"

type Frequency int32

const (
	FrequencyNanoSecond Frequency = iota
	FrequencyMilliSecond
	FrequencySecond
	FrequencyMinute
	FrequencyHour
	FrequencyDay
	FrequencyWeek
	FrequencyMonth
	FrequencyQuarter // suggar syntax (3 months)
	FrequencyYear
)

var possibleFrequencies = []Frequency{
	FrequencyNanoSecond,
	FrequencyMilliSecond,
	FrequencySecond,
	FrequencyMinute,
	FrequencyHour,
	FrequencyDay,
	FrequencyWeek,
	FrequencyMonth,
	FrequencyQuarter,
	FrequencyYear,
}

type IntervalOption struct {
	AnchorDate    time.Time
	StartDate     *time.Time
	EndDate       *time.Time
	FrequencyUnit Frequency
	IntervalValue uint32
}

type Direction int32

const (
	DirectionForward Direction = iota
	DirectionBackward
)

const (
	DEFAULT_MAX_ATTEMPTS = 50_000
	MAX_LIST_SIZE        = 25_000
)
