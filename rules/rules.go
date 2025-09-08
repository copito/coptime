package rules

func FilterRules(rules []Rule, ruleType RuleType) []Rule {
	var out []Rule
	for _, r := range rules {
		if r.IntervalType == ruleType {
			out = append(out, r)
		}
	}
	return out
}
