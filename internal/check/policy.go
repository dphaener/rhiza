package check

import "github.com/dphaener/rhiza/internal/policy"

// policy fails, under the policy rule, when the project's policy file does
// not read: it does not parse, or a value is out of range.
func (c *checker) policy() {
	if c.in.Tree == nil {
		return
	}
	if _, err := policy.Read(c.in.Tree); err != nil {
		c.add(Finding{Severity: SeverityFail, Rule: RulePolicy, Path: policy.Path, Message: err.Error()})
	}
}
