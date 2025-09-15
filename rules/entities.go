package rules

import "time"

type RuleType int32

const (
	RuleTypeInclusion RuleType = iota
	RuleTypeExclusion
)

type TimeReference struct {
	Hour        int
	Minute      int
	Second      int
	MilliSecond int
	NanoSecond  int
}

type TimeRange struct {
	StartTimeReference TimeReference
	EndTimeReference   TimeReference
}

type DayOfWeek int32

const (
	_                            = iota // zero'th term so the week days align
	DayOfWeekMonday    DayOfWeek = iota
	DayOfWeekTuesday   DayOfWeek = iota
	DayOfWeekWednesday DayOfWeek = iota
	DayOfWeekThursday  DayOfWeek = iota
	DayOfWeekFriday    DayOfWeek = iota
	DayOfWeekSaturday  DayOfWeek = iota
	DayOfWeekSunday    DayOfWeek = iota
)

type Rule struct {
	IntervalType RuleType
	TimeRange    *TimeRange
	DayOfWeeks   []time.Weekday
	MonthDays    []uint32
	Months       []time.Month
	Years        []uint32
	SetPos       []int
}
