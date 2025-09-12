package interval

import "time"

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

func nextFromAnchorWeeks(anchor, prev time.Time, interval, step int) time.Time {
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

func nextFromAnchorMonths(anchor, prev time.Time, interval, step int) time.Time {
	loc := anchor.Location()
	a := time.Date(anchor.Year(), anchor.Month(), anchor.Day(), anchor.Hour(), anchor.Minute(), anchor.Second(), anchor.Nanosecond(), loc)

	// compute how many months between a and prev
	monthsBetween := func(from, to time.Time) int {
		y := to.Year() - from.Year()
		m := int(to.Month()) - int(from.Month())
		total := y*12 + m
		// adjust if the day-of-month/time of 'to' is before 'from'’s
		cand := from.AddDate(0, total, 0)
		if to.Before(cand) {
			total--
		}
		return total
	}

	totalMonths := monthsBetween(a, prev)
	// align to our interval stride
	k := totalMonths / interval

	occ := a.AddDate(0, k*interval, 0)
	if step > 0 {
		if !occ.After(prev) {
			occ = occ.AddDate(0, interval, 0)
		}
	} else {
		if !occ.Before(prev) {
			occ = occ.AddDate(0, -interval, 0)
		}
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
