package rules

import (
	"testing"
)

func TestFilterRules(t *testing.T) {
	testTable := []struct {
		name     string
		rules    []Rule
		ruleType IntervalRuleType
		expected int
	}{
		{
			name:     "empty rules",
			rules:    []Rule{},
			ruleType: IntervalRuleTypeInclusion,
			expected: 0,
		},
		{
			name: "only inclusion rules",
			rules: []Rule{
				{IntervalType: IntervalRuleTypeInclusion},
				{IntervalType: IntervalRuleTypeInclusion},
			},
			ruleType: IntervalRuleTypeInclusion,
			expected: 2,
		},
		{
			name: "only exclusion rules",
			rules: []Rule{
				{IntervalType: IntervalRuleTypeExclusion},
				{IntervalType: IntervalRuleTypeExclusion},
			},
			ruleType: IntervalRuleTypeExclusion,
			expected: 2,
		},
		{
			name: "mixed rules",
			rules: []Rule{
				{IntervalType: IntervalRuleTypeInclusion},
				{IntervalType: IntervalRuleTypeExclusion},
				{IntervalType: IntervalRuleTypeInclusion},
				{IntervalType: IntervalRuleTypeExclusion},
			},
			ruleType: IntervalRuleTypeInclusion,
			expected: 2,
		},
		{
			name: "mixed rules, exclusion",
			rules: []Rule{
				{IntervalType: IntervalRuleTypeInclusion},
				{IntervalType: IntervalRuleTypeExclusion},
				{IntervalType: IntervalRuleTypeInclusion},
				{IntervalType: IntervalRuleTypeExclusion},
			},
			ruleType: IntervalRuleTypeExclusion,
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
