package interval

import (
	"errors"
	"iter"
	"time"
)

type Intervaler struct {
	opt IntervalOption
}

func New(option IntervalOption) Intervaler {
	anchorDate := defaultAnchorTime(option.AnchorDate)
	startDate := defaultStartTime(option.StartDate, anchorDate)
	frequencyUnit := defaultFrequencyUnit(option.FrequencyUnit)
	intervalValue := defaultIntervalValue(option.IntervalValue)
	monthEnd := option.MonthEnd

	// TODO: Validate

	if monthEnd {
		year, month, _ := anchorDate.Date()
		firstOfNextMonth := time.Date(year, month+1, 1, anchorDate.Hour(), anchorDate.Minute(), anchorDate.Second(), anchorDate.Nanosecond(), anchorDate.Location())
		lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
		anchorDate = time.Date(lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day(), anchorDate.Hour(), anchorDate.Minute(), anchorDate.Second(), anchorDate.Nanosecond(), anchorDate.Location())
		startDate = &anchorDate
	}

	intervalOption := IntervalOption{
		AnchorDate:    anchorDate,
		StartDate:     startDate,
		EndDate:       option.EndDate,
		FrequencyUnit: frequencyUnit,
		IntervalValue: intervalValue,
		MonthEnd:      monthEnd,
		ByDay:         option.ByDay,
		Wkst:          option.Wkst,
	}

	return Intervaler{
		opt: intervalOption,
	}
}

func FromRRULE(rruleString string) (*Intervaler, error) {
	option, err := convertRRULEtoIntervaler(rruleString)
	if err != nil {
		return nil, err
	}

	if option == nil {
		return nil, errors.New("options for Intervaler are empty")
	}

	iv := New(*option)
	return &iv, nil
}

// func (i *Intervaler) calculateNext(previousTime time.Time, direction Direction) time.Time {
// 	var nextTime time.Time

// 	step := 1
// 	if direction == DirectionBackward {
// 		step = -1
// 	}

// 	intervalValue := int64(i.opt.IntervalValue) * int64(step)
// 	duration := time.Duration(intervalValue)

// 	switch i.opt.FrequencyUnit {
// 	case FrequencyNanoSecond:
// 		nextTime = previousTime.Add(time.Nanosecond * duration)
// 	case FrequencyMilliSecond:
// 		nextTime = previousTime.Add(time.Millisecond * duration)
// 	case FrequencySecond:
// 		nextTime = previousTime.Add(time.Second * duration)
// 	case FrequencyMinute:
// 		nextTime = previousTime.Add(time.Minute * duration)
// 	case FrequencyHour:
// 		nextTime = previousTime.Add(time.Hour * duration)
// 	case FrequencyDay:
// 		nextTime = previousTime.AddDate(0, 0, int(intervalValue))
// 	case FrequencyWeek:
// 		nextTime = previousTime.AddDate(0, 0, int(intervalValue)*7)
// 	case FrequencyMonth:
// 		nextTime = previousTime.AddDate(0, int(intervalValue), 0)
// 	case FrequencyQuarter:
// 		nextTime = previousTime.AddDate(0, int(intervalValue)*3, 0)
// 	case FrequencyYear:
// 		nextTime = previousTime.AddDate(int(intervalValue), 0, 0)
// 	default:
// 		nextTime = time.Time{}
// 	}

// 	return nextTime
// }

func (i *Intervaler) calculateNext(previousTime time.Time, direction Direction) time.Time {
	anchor := i.opt.AnchorDate
	iv := int(i.opt.IntervalValue)
	if iv <= 0 {
		iv = 1
	}
	step := 1
	if direction == DirectionBackward {
		step = -1
	}

	switch i.opt.FrequencyUnit {
	case FrequencyNanoSecond, FrequencyMilliSecond, FrequencySecond, FrequencyMinute, FrequencyHour:
		// Sub-day units: stepping from previousTime is fine
		intervalValue := time.Duration(int64(i.opt.IntervalValue) * int64(step))
		switch i.opt.FrequencyUnit {
		case FrequencyNanoSecond:
			return previousTime.Add(time.Nanosecond * intervalValue)
		case FrequencyMilliSecond:
			return previousTime.Add(time.Millisecond * intervalValue)
		case FrequencySecond:
			return previousTime.Add(time.Second * intervalValue)
		case FrequencyMinute:
			return previousTime.Add(time.Minute * intervalValue)
		case FrequencyHour:
			return previousTime.Add(time.Hour * intervalValue)
		}

	case FrequencyDay:
		// Align to anchor day cadence
		return nextFromAnchorDays(anchor, previousTime, iv, step)

	case FrequencyWeek:
		// Align to anchor weekday/time; weeks are multiples of 7 days from anchor
		return nextFromAnchorWeeks(anchor, previousTime, iv, step, i.opt.ByDay, i.opt.Wkst)

	case FrequencyMonth:
		return nextFromAnchorMonths(anchor, previousTime, iv, step, i.opt.MonthEnd)

	case FrequencyQuarter:
		return nextFromAnchorMonths(anchor, previousTime, 3*iv, step, false)

	case FrequencyYear:
		return nextFromAnchorYears(anchor, previousTime, iv, step)
	}

	return time.Time{}
}

func (i *Intervaler) estimateIntervalSize(startTime time.Time, endTime time.Time, freq Frequency, interval int32) int32 {
	var realStartTime, realEndTime time.Time
	if startTime.Equal(endTime) {
		return 1
	} else if startTime.Before(endTime) {
		realStartTime = startTime
		realEndTime = endTime
	} else if startTime.After(endTime) {
		realStartTime = endTime
		realEndTime = startTime
	}

	diff := realEndTime.Sub(realStartTime)
	var periods int32

	switch freq {
	case FrequencyNanoSecond:
		diffRef := diff.Nanoseconds()
		periods = int32(diffRef/int64(interval)) + 1
	case FrequencyMilliSecond:
		diffRef := diff.Milliseconds()
		periods = int32(diffRef/int64(interval)) + 1
	case FrequencySecond:
		diffRef := diff.Seconds()
		periods = int32(diffRef/float64(interval)) + 1
	case FrequencyMinute:
		diffRef := diff.Minutes()
		periods = int32(diffRef/float64(interval)) + 1
	case FrequencyHour:
		diffRef := diff.Hours()
		periods = int32(diffRef/float64(interval)) + 1
	case FrequencyDay:
		diffRef := diff.Hours() / 24
		periods = int32(diffRef/float64(interval)) + 1
	case FrequencyWeek:
		diffRef := diff.Hours() / 24 / 7
		periods = int32(diffRef/float64(interval)) + 1
	case FrequencyMonth:
		diffRef := diff.Hours() / 24 / 30
		periods = int32(diffRef/float64(interval)) + 4
	case FrequencyQuarter:
		diffRef := diff.Hours() / 24 / 30 / 3
		periods = int32(diffRef/float64(interval)) + 4
	case FrequencyYear:
		diffRef := diff.Hours() / 24 / 365
		periods = int32(diffRef/float64(interval)) + 4
	default:
		periods = 10
	}

	return periods
}

func (i *Intervaler) Iterate(direction Direction, maxAttempt *int32) (iter.Seq[time.Time], error) {
	endDate := defaultEndTime(i.opt.EndDate, direction)

	// high level check
	// TODO: revalidate here
	if !i.opt.MonthEnd {
		if direction == DirectionForward && i.opt.AnchorDate.After(*i.opt.StartDate) {
			return nil, errors.New("anchor day cannot be after start day with FORWARD")
		}

		if direction == DirectionBackward && i.opt.AnchorDate.Before(*i.opt.StartDate) {
			return nil, errors.New("anchor day cannot be before start day with BACKWARD")
		}
	}

	if direction == DirectionForward && i.opt.AnchorDate.After(endDate) {
		return nil, errors.New("anchor day cannot be after end day with FORWARD")
	}

	if direction == DirectionBackward && i.opt.AnchorDate.Before(endDate) {
		return nil, errors.New("anchor day cannot be after end day with FORWARD")
	}

	maxAttemptCounter := defaultMaxAttempts(maxAttempt)
	currentAttemptCounter := maxAttemptCounter
	currentValue := i.opt.AnchorDate
	var nextValue time.Time

	return func(yield func(time.Time) bool) {
		if currentValue.Equal(*i.opt.StartDate) {
			if direction == DirectionForward {
				ok := yield(i.opt.AnchorDate)
				if !ok {
					return
				}
			}
		}

		for {
			nextValue = i.calculateNext(currentValue, direction)
			if nextValue.IsZero() {
				return
			}

			if direction == DirectionForward && nextValue.After(endDate) {
				return // ended iteration
			}

			if direction == DirectionBackward && nextValue.Before(endDate) {
				return // ended iteration
			}

			if direction == DirectionForward && nextValue.Before(*i.opt.StartDate) {
				currentAttemptCounter -= 1
				currentValue = nextValue // trying next
				continue
			}

			if direction == DirectionBackward && nextValue.After(*i.opt.StartDate) {
				currentAttemptCounter -= 1
				currentValue = nextValue // trying next
				continue
			}

			ok := yield(nextValue)
			if !ok {
				return
			}

			currentValue = nextValue
			currentAttemptCounter = maxAttemptCounter
		}
	}, nil
}

func (i *Intervaler) All(direction Direction, maxAttempt *int32) ([]time.Time, error) {
	if i.opt.EndDate == nil || i.opt.EndDate.IsZero() {
		return []time.Time{}, errors.New("cannot request all without end bound")
	}

	if i.opt.StartDate == nil || i.opt.StartDate.IsZero() {
		return []time.Time{}, errors.New("cannot request all without start bound")
	}

	estimatePeriods := i.estimateIntervalSize(*i.opt.StartDate, *i.opt.EndDate, i.opt.FrequencyUnit, int32(i.opt.IntervalValue))

	if estimatePeriods >= MAX_LIST_SIZE {
		return []time.Time{}, errors.New("maximum size is 25k - please select a smaller range")
	}

	list := make([]time.Time, 0, estimatePeriods*2)
	iterator, err := i.Iterate(direction, maxAttempt)
	if err != nil {
		return []time.Time{}, err
	}

	for value := range iterator {
		list = append(list, value)
	}

	return list, nil
}

func (i *Intervaler) Between(direction Direction, startTime time.Time, endTime time.Time, maxAttempt *int32) ([]time.Time, error) {
	estimatePeriods := i.estimateIntervalSize(startTime, endTime, i.opt.FrequencyUnit, int32(i.opt.IntervalValue))

	if estimatePeriods >= MAX_LIST_SIZE {
		return []time.Time{}, errors.New("maximum size is 25k - please select a smaller range")
	}

	list := make([]time.Time, 0, estimatePeriods*2)
	iterator, err := i.Iterate(direction, maxAttempt)
	if err != nil {
		return []time.Time{}, err
	}

	for value := range iterator {
		if (value.Before(endTime) && value.After(startTime)) || value.Equal(startTime) || value.Equal(endTime) {
			list = append(list, value)
		}
	}

	return list, nil
}
