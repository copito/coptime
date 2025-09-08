package rules

import "time"

type IntervalRuleType int32

const (
	IntervalRuleTypeInclusion IntervalRuleType = iota
	IntervalRuleTypeExclusion
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
	IntervalType IntervalRuleType
	TimeRange    *TimeRange
	DayOfWeeks   []time.Weekday
	MonthDays    []uint32
	Months       []time.Month
	Years        []uint32
}
