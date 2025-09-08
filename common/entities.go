package common

import "time"

type SubWindowResult struct {
	Start time.Time
	End   time.Time
}

type WindowResult struct {
	Start     time.Time
	End       time.Time
	SubWindow []SubWindowResult
}
