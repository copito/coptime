package window

import (
	"github.com/copito/coptime/interval"
	rules "github.com/copito/coptime/rules"
)

const (
	DEFAULT_MAX_ATTEMPTS int32 = 10_000
)

type WindowOption struct {
	interval.IntervalOption

	Size float64

	// All rules applied on top of the windows
	Rules []rules.Rule
}
