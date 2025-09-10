package interval

import "time"

func addMonth(t time.Time, months int) time.Time {
	year, month, day := t.Date()

	newMonth := month + time.Month(months)
	newYear := year

	// Normalize month and year
	if newMonth > 12 {
		newYear += (int(newMonth) - 1) / 12
		newMonth = time.Month((int(newMonth)-1)%12 + 1)
	} else if newMonth < 1 {
		newYear += int(newMonth-1)/12 - 1
		newMonth = time.Month(int(newMonth-1)%12+12)
	}

	// Check for invalid day
	lastDayOfMonth := time.Date(newYear, newMonth+1, 0, 0, 0, 0, 0, t.Location()).Day()
	if day > lastDayOfMonth {
		day = lastDayOfMonth
	}

	return time.Date(newYear, newMonth, day, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
}


func nextFromAnchorDays(anchor, prev time.Time, interval, step int) time.Time {
	// keep times in the same location
	loc := anchor.Location()
	a := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), loc)

	// number of whole intervals between anchor and prev
	delta := prev.Sub(a)
	intervalDur := time.Hour * 24 * time.Duration(interval)
	k := int(delta / intervalDur)

	occ := a.Add(time.Duration(k) * intervalDur)
	if step > 0 {
		if !occ.After(prev) {
			occ = occ.Add(intervalDur)
		}
	} else {
		if !occ.Before(prev) {
			occ = occ.Add(-intervalDur)
		}
	}
	return occ
}

func startOfWeek(t time.Time, wkst time.Weekday) time.Time {
	weekday := t.Weekday()
	offset := int(weekday - wkst)
	if offset < 0 {
		offset += 7
	}
	return t.AddDate(0, 0, -offset)
}

func nextFromAnchorWeeks(anchor, prev time.Time, interval, step int, byDay []int, wkst time.Weekday) time.Time {
	if len(byDay) == 0 {
		loc := anchor.Location()
		a := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), loc)

		weekDur := 7 * 24 * time.Hour
		intervalDur := time.Duration(interval) * weekDur

		delta := prev.Sub(a)
		k := int(delta / intervalDur)

		occ := a.Add(time.Duration(k) * intervalDur)
		if step > 0 {
			// forward: strictly after prev
			if !occ.After(prev) {
				occ = occ.Add(intervalDur)
			}
		} else {
			// backward: strictly before prev
			if !occ.Before(prev) {
				occ = occ.Add(-intervalDur)
			}
		}
		return occ
	}

	// BYDAY logic
	dayIterator := prev
	anchorWeekStart := startOfWeek(anchor, wkst)
	for {
		dayIterator = dayIterator.AddDate(0, 0, 1*step)

		isValidDay := false
		for _, wd := range byDay {
			if dayIterator.Weekday() == time.Weekday(wd) {
				isValidDay = true
				break
			}
		}
		if !isValidDay {
			continue
		}

		currentWeekStart := startOfWeek(dayIterator, wkst)
		weekDiff := int(currentWeekStart.Sub(anchorWeekStart).Hours() / 24 / 7)

		if weekDiff%interval == 0 {
			return dayIterator
		}
	}
}

func nextFromAnchorMonths(anchor, prev time.Time, interval, step int, monthEnd bool) time.Time {
	loc := anchor.Location()
	a := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), loc)

	// compute how many months between a and prev
	monthsBetween := func(from, to time.Time) int {
		y := to.Year() - from.Year()
		m := int(to.Month()) - int(from.Month())
		total := y*12 + m
		// adjust if the day-of-month/time of 'to' is before 'from'’s
		cand := addMonth(from, total)
		if to.Before(cand) {
			total--
		}
		return total
	}

	totalMonths := monthsBetween(a, prev)
	// align to our interval stride
	k := totalMonths / interval

	occ := addMonth(a, k*interval)
	if step > 0 {
		if !occ.After(prev) {
			occ = addMonth(occ, interval)
		}
	} else {
		if !occ.Before(prev) {
			occ = addMonth(occ, -interval)
		}
	}

	if monthEnd {
		// adjust occ to be the last day of its month
		year, month, _ := occ.Date()
		firstOfNextMonth := time.Date(year, month+1, 1, anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), anchor.Location())
		lastOfMonth := firstOfNextMonth.AddDate(0, 0, -1)
		occ = time.Date(lastOfMonth.Year(), lastOfMonth.Month(), lastOfMonth.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), anchor.Location())
	}

	return occ
}

func nextFromAnchorYears(anchor, prev time.Time, interval, step int) time.Time {
	loc := anchor.Location()
	a := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), loc)

	yearsBetween := func(from, to time.Time) int {
		y := to.Year() - from.Year()
		cand := time.Date(from.Year()+y, from.Month(), from.Day(), from.Hour(), from.Minute(), from.Second(), from.Nanosecond(), from.Location())
		if to.Before(cand) {
			y--
		}
		return y
	}

	totalYears := yearsBetween(a, prev)
	k := totalYears / interval

	occ := a.AddDate(k*interval, 0, 0)
	if step > 0 {
		if !occ.After(prev) {
			occ = occ.AddDate(interval, 0, 0)
		}
	} else {
		if !occ.Before(prev) {
			occ = occ.AddDate(-interval, 0, 0)
		}
	}
	return occ
}
