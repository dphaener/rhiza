package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/anchor"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
)

// anchors fails, under the anchor rule, once on each live as-is anchor that
// no longer resolves against Input.Tree: its path or its symbol is gone.
// It warns, under the same rule, when the anchor resolves but the fingerprint
// of what it names differs from the stored one. An anchor without a path has
// nothing to resolve.
func (c *checker) anchors(m model.Model) {
	if c.in.Tree == nil {
		return
	}
	implements := map[id.ID][]string{}
	for _, l := range m.Links() {
		if l.Retired || l.Type != op.LinkImplements {
			continue
		}
		if n, ok := m.Node(l.To); ok && !n.Retired {
			implements[l.From] = append(implements[l.From], name(n))
		}
	}
	for _, n := range m.Nodes() {
		if n.Kind != op.KindAnchor || n.Retired {
			continue
		}
		p, symbol := text(n, "path"), text(n, "symbol")
		if p == "" {
			continue
		}
		got, err := anchor.FingerprintFS(c.in.Tree, p, symbol)
		if err == nil {
			if want := text(n, "fingerprint"); want != "" && want != got {
				what := "file " + p
				if symbol != "" {
					what = "symbol " + symbol + " of " + p
				}
				c.add(Finding{Severity: SeverityWarn, Rule: RuleAnchor, Path: p, ChangeSet: str(n.Created),
					Message: fmt.Sprintf("anchor %s (%s): %s changed since it was fingerprinted", name(n), implementsOf(implements, n), what),
					Next:    "rhiza anchor refresh " + name(n)})
			}
			continue
		}
		var what string
		switch {
		case errors.Is(err, anchor.ErrNoSymbol), errors.Is(err, anchor.ErrBadSymbol):
			what = fmt.Sprintf("symbol %s of %s is gone: %v", symbol, p, err)
		default:
			what = fmt.Sprintf("path %s is gone: %v", p, err)
		}
		by := implementsOf(implements, n)
		c.add(Finding{Severity: SeverityFail, Rule: RuleAnchor, Path: p, ChangeSet: str(n.Created),
			Message: fmt.Sprintf("anchor %s (%s): %s", name(n), by, what),
			Next:    "rhiza set " + name(n) + " --field path=<path> --change <change set>"})
	}
}

// implementsOf says which capabilities n implements, sorted.
func implementsOf(implements map[id.ID][]string, n model.Node) string {
	caps := slices.Clone(implements[n.ID])
	slices.Sort(caps)
	if len(caps) == 0 {
		return "implements nothing"
	}
	return "implements " + strings.Join(caps, ", ")
}

func name(n model.Node) string {
	if n.Handle != "" {
		return n.Handle
	}
	return n.ID.String()
}

func text(n model.Node, field string) string {
	var s string
	if raw, ok := n.Fields[field]; ok && json.Unmarshal(raw, &s) == nil {
		return s
	}
	return ""
}
