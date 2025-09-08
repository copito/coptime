package window

import (
	"time"

	"github.com/copito/coptime/interval"
	rules "github.com/copito/coptime/rules"
)

const (
	DEFAULT_MAX_ATTEMPTS int32 = 10_000
)

type SubWindowResult struct {
	Start time.Time
	End   time.Time
}

type WindowResult struct {
	Start     time.Time
	End       time.Time
	SubWindow []SubWindowResult
}

type WindowOption struct {
	interval.IntervalOption

	// All rules applied on top of the windows
	Rules []rules.Rule
}
