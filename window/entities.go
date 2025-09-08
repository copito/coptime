package window

import (
	"time"

	"github.com/copito/coptime/interval"
	rules "github.com/copito/coptime/rules"
)

type SubWindowResult struct {
	StartTime time.Time
	EndTime   time.Time
}

type WindowResult struct {
	StartWindow time.Time
	EndWindow   time.Time
	SubWindow   []SubWindowResult
}

type WindowOption struct {
	interval.IntervalOption

	// All rules applied on top of the windows
	Rules []rules.Rules
}
