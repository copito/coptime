package rules

func FilterRules(rules []Rule, ruleType IntervalRuleType) []Rule {
	var out []Rule
	for _, r := range rules {
		if r.IntervalType == ruleType {
			out = append(out, r)
		}
	}
	return out
}
