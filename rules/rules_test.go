package rules

import (
	"testing"
)

func TestFilterRules(t *testing.T) {
	testTable := []struct {
		name     string
		rules    []Rule
		ruleType RuleType
		expected int
	}{
		{
			name:     "empty rules",
			rules:    []Rule{},
			ruleType: RuleTypeInclusion,
			expected: 0,
		},
		{
			name: "only inclusion rules",
			rules: []Rule{
				{IntervalType: RuleTypeInclusion},
				{IntervalType: RuleTypeInclusion},
			},
			ruleType: RuleTypeInclusion,
			expected: 2,
		},
		{
			name: "only exclusion rules",
			rules: []Rule{
				{IntervalType: RuleTypeExclusion},
				{IntervalType: RuleTypeExclusion},
			},
			ruleType: RuleTypeExclusion,
			expected: 2,
		},
		{
			name: "mixed rules",
			rules: []Rule{
				{IntervalType: RuleTypeInclusion},
				{IntervalType: RuleTypeExclusion},
				{IntervalType: RuleTypeInclusion},
				{IntervalType: RuleTypeExclusion},
			},
			ruleType: RuleTypeInclusion,
			expected: 2,
		},
		{
			name: "mixed rules, exclusion",
			rules: []Rule{
				{IntervalType: RuleTypeInclusion},
				{IntervalType: RuleTypeExclusion},
				{IntervalType: RuleTypeInclusion},
				{IntervalType: RuleTypeExclusion},
			},
			ruleType: RuleTypeExclusion,
			expected: 2,
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			result := FilterRules(tt.rules, tt.ruleType)
			if len(result) != tt.expected {
				t.Errorf("expected %d rules, got %d", tt.expected, len(result))
			}
		})
	}
}
