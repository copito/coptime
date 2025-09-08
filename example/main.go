package main

import (
	"fmt"
	"time"

	rules "github.com/copito/coptime/rules"

	"github.com/copito/coptime/interval"
	"github.com/copito/coptime/window"
)

func exampleIterator() {
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	// ---------------------
	// Iterator
	// ---------------------
	opts := interval.IntervalOption{
		AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		StartDate:     &startTime,
		EndDate:       &endTime,
		FrequencyUnit: interval.FrequencyDay,
		IntervalValue: 1,
	}
	iv := interval.New(opts)

	iterator, err := iv.Iterate(interval.DirectionForward, nil)
	if err != nil {
		panic(fmt.Errorf("error: %v", err.Error()))
	}

	i := 1
	for value := range iterator {
		if i >= 10 {
			break
		}

		fmt.Printf("Values for: %s\n", value)
		i++
	}

	fmt.Println("Done!")
}

func exampleBetween() {
	// ---------------------
	// Find the between
	// ---------------------
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	// ---------------------
	// Iterator
	// ---------------------
	opts := interval.IntervalOption{
		AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		StartDate:     &startTime,
		EndDate:       &endTime,
		FrequencyUnit: interval.FrequencyDay,
		IntervalValue: 1,
	}
	iv := interval.New(opts)

	values, err := iv.Between(
		interval.DirectionForward,
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
		nil,
	)
	if err != nil {
		panic(fmt.Errorf("error: %v", err.Error()))
	}

	fmt.Println(values)
}

func exampleAll() {
	// ---------------------
	// Find the ALL
	// ---------------------
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	// ---------------------
	// Iterator
	// ---------------------
	opts := interval.IntervalOption{
		AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
		StartDate:     &startTime,
		EndDate:       &endTime,
		FrequencyUnit: interval.FrequencyDay,
		IntervalValue: 1,
	}
	iv := interval.New(opts)
	valuesList, err := iv.All(interval.DirectionForward, nil)
	if err != nil {
		panic(fmt.Errorf("error: %v", err.Error()))
	}

	fmt.Println(valuesList)
}

func exampleHourly() {
	// ---------------------
	// Iterator Hourly
	// ---------------------
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)
	opts := interval.IntervalOption{
		AnchorDate:    time.Date(2024, 12, 31, 10, 23, 0, 0, time.UTC),
		StartDate:     &startTime,
		EndDate:       &endTime,
		FrequencyUnit: interval.FrequencyHour,
		IntervalValue: 3,
	}
	iv := interval.New(opts)

	iterator, err := iv.Iterate(interval.DirectionForward, nil)
	if err != nil {
		panic(fmt.Errorf("error: %v", err.Error()))
	}

	j := 1
	for value := range iterator {
		if j >= 10 {
			break
		}

		fmt.Printf("Values for: %s\n", value)
		j++
	}

	fmt.Println("Done!")
}

func exampleWindowIterator() {
	startTime := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endTime := time.Date(2025, 2, 2, 0, 0, 0, 0, time.UTC)

	// ---------------------
	// Iterator
	// ---------------------
	opts := window.WindowOption{
		IntervalOption: interval.IntervalOption{
			AnchorDate:    time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			StartDate:     &startTime,
			EndDate:       &endTime,
			FrequencyUnit: interval.FrequencyDay,
			IntervalValue: 1,
		},
		Rules: []rules.Rules{},
	}
	w := window.New(opts)

	iterator, err := w.Iterate(interval.DirectionForward, nil)
	if err != nil {
		panic(fmt.Errorf("error: %v", err.Error()))
	}

	i := 1
	for value := range iterator {
		if i >= 10 {
			break
		}

		fmt.Printf("Values for: %s\n", value)
		i++
	}

	fmt.Println("Done!")
}

func main() {
	// ---------------------
	// Interval - Iterator
	// ---------------------
	exampleIterator()

	// ---------------------
	// Interval - Find the between
	// ---------------------
	exampleBetween()

	// ---------------------
	// Interval - Find the ALL
	// ---------------------
	exampleAll()

	// ---------------------
	// Interval - Iterator Hourly
	// ---------------------
	exampleHourly()

	// ---------------------
	// Window - Iterator
	// ---------------------
	exampleWindowIterator()
}
