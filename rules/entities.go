package rules

type IntervalRuleType int32

const (
	IntervalRuleTypeInclusion IntervalRuleType = iota
	IntervalRuleTypeExclusion IntervalRuleType = iota
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

type Month int32

const (
	_                    = iota // 0 th term so the other months align
	MonthJanuary   Month = iota
	MonthFebruary  Month = iota
	MonthMarch     Month = iota
	MonthApril     Month = iota
	MonthMay       Month = iota
	MonthJune      Month = iota
	MonthJuly      Month = iota
	MonthAugust    Month = iota
	MonthSeptember Month = iota
	MonthOctober   Month = iota
	MonthNovember  Month = iota
	MonthDecember  Month = iota
)

type Rules struct {
	IntervalType IntervalRuleType
	TimeRange    TimeRange
	DayOfWeeks   []DayOfWeek
	Months       []Month
	MonthDays    []uint32
	Years        []uint32
}
