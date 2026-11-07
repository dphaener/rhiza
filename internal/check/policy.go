package check

import (
	"strings"

	"github.com/dphaener/rhiza/internal/policy"
)

// policy fails, under the policy rule, when the project's policy file does
// not read: it does not parse, or a value is out of range. The finding's
// Path names the file, so the message drops policy.Read's own path prefix.
func (c *checker) policy() {
	if c.in.Tree == nil {
		return
	}
	if _, err := policy.Read(c.in.Tree); err != nil {
		c.add(Finding{Severity: SeverityFail, Rule: RulePolicy, Path: policy.Path, Message: strings.TrimPrefix(err.Error(), policy.Path+": ")})
	}
}
