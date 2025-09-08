# coptime

**coptime** is a Go library for working with time-based intervals, rules, and windows.  
It provides utilities for iterating over time periods, generating values between dates, and applying rule-based windows.

---

## Installation

Install the package using `go get`:

```bash
go get github.com/copito/coptime
```

```go
import (
    "github.com/copito/coptime/interval"
    "github.com/copito/coptime/window"
    rules "github.com/copito/coptime/rules"
)
```

## Features

- Intervals

1. Iterate forward or backward through time.

1. Generate values between two dates.

1. Collect all values within an interval.

1. Flexible frequency units (days, hours, etc.).

- Windows

1. Create windowed iterations using rules.

1. Combine interval logic with custom constraints.

## Usage Examples

See the examples in this repo for more reference. Below are some core use cases.

### 1. Interval Iterator

```go
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
    panic(err)
}

for value := range iterator {
    fmt.Printf("Value: %s\n", value)
}
```

### 2. Find Between

```go
values, err := iv.Between(
    interval.DirectionForward,
    time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
    time.Date(2025, 1, 25, 0, 0, 0, 0, time.UTC),
    nil,
)
if err != nil {
    panic(err)
}
fmt.Println(values)
```

### 3. Collect All Values

```go
valuesList, err := iv.All(interval.DirectionForward, nil)
if err != nil {
    panic(err)
}
fmt.Println(valuesList)
```

### 4. Hourly Intervals

```go
opts := interval.IntervalOption{
    AnchorDate:    time.Date(2024, 12, 31, 10, 23, 0, 0, time.UTC),
    StartDate:     &startTime,
    EndDate:       &endTime,
    FrequencyUnit: interval.FrequencyHour,
    IntervalValue: 3,
}
iv := interval.New(opts)
```

### 5. Window Iterator

```go
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
    panic(err)
}
```

## Running the Examples

To try out the full example program:

```bash
go run ./examples/main.go
```

This will run:

- Interval iteration

- Finding values between dates

- Collecting all values

- Hourly intervals

- Window iteration

## License

MIT [License](./LICENSE.md)

---
