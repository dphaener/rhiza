// Package check checks a repository's integrity: that its log reads and
// folds, that what has landed makes a sound model, that every link joins the
// kinds its type allows, that the rendered markdown is current and that the
// log has only grown since a base revision. It is what `genesis check` and
// CI run. Run is a pure function of its Input: it reads nothing but the
// fs.FS it is given, runs no git and never panics.
//
// The vocabulary:
//
//   - A finding is one thing check found. Its severity is fail or warn; any
//     fail makes the report not OK, warnings alone do not.
//   - A rule is the check that produced a finding; the rules are listed
//     below in report order.
//   - An endpoint's End is its node's kind, retired or not, or ChangeSet
//     for a change set; the endpoint table (endpoints.go) says which pairs
//     of Ends each link type allows. A live link to a retired node is not a
//     finding: only kinds matter.
//   - The in-flight problems are those of model.ToBe and of model.WhatIf for
//     each draft, proposed or accepted change set, minus those as-is already
//     has (same change set, op and error text), each distinct one once.
//
// The rules, and what each fails and warns on:
//
//   - log fails on each problem store.ReadLog reported.
//   - fold fails on each changeset.FoldLog problem but only warns on a
//     concurrent pair (ErrConcurrent). Fold problems are permanent, since
//     the log is append-only, and a concurrent pair is what an honest merge
//     of two branches leaves.
//   - model fails on each as-is problem. A later change set repairs one; for
//     a conflict, Next says how.
//   - in-flight warns on each in-flight problem.
//   - link fails on each live as-is link whose Ends the endpoint table does
//     not allow.
//   - render fails on each rendered file that differs from what
//     render.Render writes, is missing or is extra, and on a missing
//     rendered folder.
//   - append-only fails on each log file changed since the base other than
//     by being added, and warns once when the rule was skipped (no git, no
//     base).
//
// Findings are sorted by rule in the order above, then path, change set,
// entry, op, link, message, severity and next.
package check

import (
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io/fs"
	"slices"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
	"github.com/dphaener/genesis/internal/model"
	"github.com/dphaener/genesis/internal/render"
	"github.com/dphaener/genesis/internal/store"
)

// The two severities.
const (
	SeverityFail = "fail"
	SeverityWarn = "warn"
)

// The rules, in report order.
const (
	RuleLog        = "log"
	RuleFold       = "fold"
	RuleModel      = "model"
	RuleInFlight   = "in-flight"
	RuleLink       = "link"
	RuleRender     = "render"
	RuleAppendOnly = "append-only"
)

// Rules returns the rules in report order.
func Rules() []string {
	return []string{RuleLog, RuleFold, RuleModel, RuleInFlight, RuleLink, RuleRender, RuleAppendOnly}
}

// Paths in findings are relative to the project root.
const (
	genesisDir  = ".genesis/"
	renderedDir = ".genesis/rendered"
)

// Input is everything Run checks.
type Input struct {
	// Log is the log as store.ReadLog read it.
	Log store.Log
	// Rendered is rooted at .genesis/rendered; nil when the folder is
	// missing.
	Rendered fs.FS
	// Base is the base revision used, "" when the append-only rule was
	// skipped.
	Base string
	// Skipped is why the append-only rule was skipped, "" otherwise.
	Skipped string
	// Changes are the log files changed since Base.
	Changes []Change
}

// Change is one log file changed since the base: a git name-status letter
// ("A", "M", "D", ...) and a path relative to the project root, such as
// .genesis/log/<change set>/<entry>.json.
type Change struct {
	Status string
	Path   string
}

// Finding is one thing check found. Every field but Severity, Rule and
// Message may be "". Path is relative to the project root; ChangeSet,
// Entry, Op and Link are IDs; Next is a command that would help.
type Finding struct {
	Severity  string
	Rule      string
	Path      string
	ChangeSet string
	Entry     string
	Op        string
	Link      string
	Message   string
	Next      string
}

// Report is the result of a check. OK is true when no finding fails. Base
// is Input.Base. Findings is never nil.
type Report struct {
	OK       bool
	Base     string
	Findings []Finding
}

// Run checks in and reports every finding.
func Run(in Input) Report {
	c := checker{in: in}
	c.log()
	sets := c.fold()
	asIs := c.model(sets)
	c.links(asIs, sets)
	c.render(asIs, sets)
	c.appendOnly()

	rank := map[string]int{}
	for i, r := range Rules() {
		rank[r] = i
	}
	slices.SortFunc(c.findings, func(a, b Finding) int {
		return cmp.Or(cmp.Compare(rank[a.Rule], rank[b.Rule]), strings.Compare(a.Path, b.Path),
			strings.Compare(a.ChangeSet, b.ChangeSet), strings.Compare(a.Entry, b.Entry),
			strings.Compare(a.Op, b.Op), strings.Compare(a.Link, b.Link),
			strings.Compare(a.Message, b.Message), strings.Compare(a.Severity, b.Severity),
			strings.Compare(a.Next, b.Next))
	})
	ok := !slices.ContainsFunc(c.findings, func(f Finding) bool { return f.Severity == SeverityFail })
	return Report{OK: ok, Base: in.Base, Findings: c.findings}
}

type checker struct {
	in       Input
	findings []Finding
}

func (c *checker) add(f Finding) { c.findings = append(c.findings, f) }

// str is i's string form, "" for the zero ID.
func str(i id.ID) string {
	if i.IsZero() {
		return ""
	}
	return i.String()
}

func (c *checker) log() {
	for _, p := range c.in.Log.Problems {
		c.add(Finding{Severity: SeverityFail, Rule: RuleLog, Path: genesisDir + p.Path, Message: p.Err.Error()})
	}
}

func (c *checker) fold() []changeset.ChangeSet {
	sets, problems := changeset.FoldLog(c.in.Log)
	for _, p := range problems {
		f := Finding{Severity: SeverityFail, Rule: RuleFold, ChangeSet: str(p.ChangeSet), Entry: str(p.Entry),
			Message: p.Err.Error()}
		if errors.Is(p.Err, changeset.ErrConcurrent) {
			f.Severity = SeverityWarn
		}
		if !p.Entry.IsZero() {
			f.Path = fmt.Sprintf("%s%s/%s/%s.json", genesisDir, store.LogDir, p.ChangeSet, p.Entry)
		}
		c.add(f)
	}
	return sets
}

// problemKey identifies a model problem across views.
type problemKey struct {
	cs, op id.ID
	err    string
}

func keyOf(p model.Problem) problemKey { return problemKey{p.ChangeSet, p.Op, p.Err.Error()} }

func (c *checker) model(sets []changeset.ChangeSet) model.Model {
	asIs, problems := model.AsIs(sets)
	seen := map[problemKey]bool{}
	for _, p := range problems {
		seen[keyOf(p)] = true
		f := Finding{Severity: SeverityFail, Rule: RuleModel, ChangeSet: str(p.ChangeSet), Op: str(p.Op),
			Message: p.Err.Error()}
		if errors.Is(p.Err, model.ErrConflict) {
			f.Next = fmt.Sprintf("genesis change new %q", "Resolve the conflict on "+target(sets, p))
		}
		c.add(f)
	}

	_, toBe := model.ToBe(sets)
	views := [][]model.Problem{toBe}
	for _, s := range sets {
		switch s.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			_, ps := model.WhatIf(sets, s.ID)
			views = append(views, ps)
		}
	}
	for _, ps := range views {
		for _, p := range ps {
			k := keyOf(p)
			if seen[k] {
				continue
			}
			seen[k] = true
			c.add(Finding{Severity: SeverityWarn, Rule: RuleInFlight, ChangeSet: str(p.ChangeSet), Op: str(p.Op),
				Message: p.Err.Error()})
		}
	}
	return asIs
}

// target names what problem p's op acts on: its target's ID, or the op's
// own ID when the op is not found.
func target(sets []changeset.ChangeSet, p model.Problem) string {
	for _, s := range sets {
		if s.ID != p.ChangeSet {
			continue
		}
		for _, o := range s.Ops {
			if o.ID == p.Op && !o.Target().IsZero() {
				return o.Target().String()
			}
		}
	}
	return str(p.Op)
}

func (c *checker) links(m model.Model, sets []changeset.ChangeSet) {
	known := make(map[id.ID]bool, len(sets))
	for _, s := range sets {
		known[s.ID] = true
	}
	end := func(i id.ID) (End, string) {
		if n, ok := m.Node(i); ok {
			return NodeEnd(n.Kind), string(n.Kind)
		}
		if known[i] {
			return ChangeSet, string(ChangeSet)
		}
		return "", "unknown " + i.String()
	}
	for _, l := range m.Links() {
		if l.Retired {
			continue
		}
		from, fromName := end(l.From)
		to, toName := end(l.To)
		if Allowed(l.Type, from, to) {
			continue
		}
		f := Finding{Severity: SeverityFail, Rule: RuleLink, ChangeSet: str(l.Created), Link: l.ID.String(),
			Message: fmt.Sprintf("link type %q does not go from %s to %s", l.Type, fromName, toName)}
		if len(l.History) > 0 {
			f.Op = str(l.History[0].Op)
		}
		c.add(f)
	}
}

func (c *checker) render(m model.Model, sets []changeset.ChangeSet) {
	const next = "genesis render"
	fail := func(p, msg string) {
		c.add(Finding{Severity: SeverityFail, Rule: RuleRender, Path: p, Message: msg, Next: next})
	}
	if c.in.Rendered == nil {
		fail(renderedDir, "the rendered folder is missing")
		return
	}
	want := render.Render(m, sets)
	have := map[string]bool{}
	err := fs.WalkDir(c.in.Rendered, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			fail(renderedPath(p), "cannot read: "+err.Error())
			if d != nil && d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}
		if !d.IsDir() {
			have[p] = true
		}
		return nil
	})
	if err != nil {
		fail(renderedDir, "cannot read: "+err.Error())
	}
	for p := range have {
		if _, ok := want[p]; !ok {
			fail(renderedPath(p), "extra: render does not write this file")
		}
	}
	for p, b := range want {
		if !have[p] {
			fail(renderedPath(p), "missing: render writes this file")
			continue
		}
		got, err := fs.ReadFile(c.in.Rendered, p)
		switch {
		case err != nil:
			fail(renderedPath(p), "cannot read: "+err.Error())
		case !bytes.Equal(got, b):
			fail(renderedPath(p), "differs from what render writes")
		}
	}
}

func renderedPath(p string) string {
	if p == "." {
		return renderedDir
	}
	return renderedDir + "/" + p
}

func (c *checker) appendOnly() {
	if c.in.Skipped != "" {
		c.add(Finding{Severity: SeverityWarn, Rule: RuleAppendOnly,
			Message: "append-only check skipped: " + c.in.Skipped})
	}
	for _, ch := range c.in.Changes {
		if ch.Status == "A" {
			continue
		}
		f := Finding{Severity: SeverityFail, Rule: RuleAppendOnly, Path: ch.Path,
			Message: fmt.Sprintf("log file %s since the base (status %s); entries may only be added", verb(ch.Status), ch.Status)}
		if c.in.Base != "" {
			f.Next = fmt.Sprintf("git checkout %s -- %s", c.in.Base, ch.Path)
		}
		c.add(f)
	}
}

// verb is what git's name-status letter s means.
func verb(s string) string {
	switch {
	case strings.HasPrefix(s, "M"):
		return "modified"
	case strings.HasPrefix(s, "D"):
		return "deleted"
	case strings.HasPrefix(s, "T"):
		return "changed type"
	case strings.HasPrefix(s, "R"):
		return "renamed"
	case strings.HasPrefix(s, "C"):
		return "copied"
	}
	return "changed"
}
