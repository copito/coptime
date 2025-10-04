package rules

import (
	"sort"

	"github.com/copito/coptime/common"
)

// NewOccurrenceFilter returns a filter that selects the nth occurrences from the
// provided slice of sub-windows based on the supplied indexes. Positive indexes
// select from the start (1-based) and negative indexes select from the end.
func NewOccurrenceFilter(indexes []int) SubWindowFilter {
	if len(indexes) == 0 {
		return nil
	}

	// Copy indexes to avoid mutating the caller's slice when the closure runs.
	positions := append([]int(nil), indexes...)

	return func(windows []common.SubWindowResult) []common.SubWindowResult {
		if len(positions) == 0 || len(windows) == 0 {
			return windows
		}

		sorted := append([]common.SubWindowResult(nil), windows...)
		sort.Slice(sorted, func(i, j int) bool {
			return sorted[i].Start.Before(sorted[j].Start)
		})

		selected := make([]common.SubWindowResult, 0, len(positions))
		used := make(map[int]struct{}, len(positions))

		for _, pos := range positions {
			if pos == 0 {
				continue
			}

			idx := 0
			if pos > 0 {
				idx = pos - 1
			} else {
				idx = len(sorted) + pos
			}

			if idx < 0 || idx >= len(sorted) {
				continue
			}

			if _, exists := used[idx]; exists {
				continue
			}

			selected = append(selected, sorted[idx])
			used[idx] = struct{}{}
		}

		sort.Slice(selected, func(i, j int) bool {
			return selected[i].Start.Before(selected[j].Start)
		})

		if len(selected) == 0 {
			return []common.SubWindowResult{}
		}

		return selected
	}
}
