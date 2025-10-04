package window

import (
	"time"

	"github.com/copito/coptime/common"
	"github.com/copito/coptime/interval"
	rulespkg "github.com/copito/coptime/rules"
)

type ruleBinding struct {
	rule rulespkg.Rule
}

func buildRuleBindings(opt WindowOption, ruleType rulespkg.RuleType) []ruleBinding {
	bindings := make([]ruleBinding, 0, len(opt.Rules))
	for _, rule := range opt.Rules {
		if rule.IntervalType != ruleType {
			continue
		}

		bindings = append(bindings, ruleBinding{rule: rule})
	}
	return bindings
}

func generateSubWindowsForBindings(direction interval.Direction, unit interval.Frequency, start time.Time, end time.Time, bindings []ruleBinding, tz *time.Location) []common.SubWindowResult {
	if len(bindings) == 0 {
		return []common.SubWindowResult{}
	}

	ruleSet := make([]rulespkg.Rule, 0, len(bindings))
	for _, binding := range bindings {
		ruleSet = append(ruleSet, binding.rule)
	}

	return rulespkg.GenerateSubWindowsForRuleType(direction, unit, start, end, ruleSet, tz)
}
