package check

import (
	"fmt"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// unlinkedConstraints warns, under the link rule, once on each live as-is
// constraint with no live constrains link from it. A constraint that
// constrains nothing binds nothing; links of other types (derived-from,
// relates-to) do not make it bind.
func (c *checker) unlinkedConstraints(m model.Model) {
	linked := map[id.ID]bool{}
	for _, l := range m.Links() {
		if !l.Retired && l.Type == op.LinkConstrains {
			linked[l.From] = true
		}
	}
	for _, n := range m.Nodes() {
		if n.Kind != op.KindConstraint || n.Retired || linked[n.ID] {
			continue
		}
		name := n.Handle
		if name == "" {
			name = n.ID.String()
		}
		c.add(Finding{Severity: SeverityWarn, Rule: RuleLink, ChangeSet: str(n.Created),
			Message: fmt.Sprintf("constraint %s constrains nothing: it has no constrains link", name),
			Next:    "rhiza link " + name + " constrains <capability-or-intent> --change <change set>"})
	}
}
