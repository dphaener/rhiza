// Package check checks a repository's integrity: that its log reads and
// folds, that what has landed makes a sound model, that every link joins the
// kinds its type allows, that the rendered markdown is current and that the
// log has only grown since a base revision. It is what `rhiza check` and
// CI run. Run is a pure function of its Input: it reads nothing but the
// fs.FS it is given, runs no git and never panics. What git knows, the base
// entries and the commits' trailers, its caller gathers into the Input.
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
//   - in-flight warns on each in-flight problem, Next on a conflict saying
//     where to look, and on each op of an in-flight change set that rests on
//     a dropped one (changeset.RestsOn), on each pair of proposed or accepted
//     change sets writing overlapping aspects of one target
//     (model.Overlaps; Next is rhiza slice --change <A> --view what-if), and on each live link that a draft,
//     proposed or accepted change set wrote itself, in its what-if, whose
//     Ends the endpoint table does not allow, each once.
//   - link fails on each live as-is link whose Ends the endpoint table does
//     not allow.
//   - link also warns once on each live as-is constraint with no live
//     constrains link from it; derived-from or relates-to links alone do
//     not count (unlinked.go).
//   - depends-on fails once on each depends-on loop among change sets
//     (changeset.Cycles), whatever their states, and once on each loop of
//     live as-is depends-on links among capabilities.
//   - evidence-actor warns once on each live Evidence with a live as-is
//     verifies link to a capability when the Evidence's actor (the actor of
//     its create, model.CreatedBy) is named, by handle or ID, in an Actor:
//     trailer of a commit that changes a file outside .rhiza/, whose
//     Change: trailers name the change set that created the capability and
//     whose Leaf: trailers, if it has any, name the capability by handle or
//     ID: whoever wrote a leaf's code should not be the one to verify it. A
//     commit of model writes alone is not a leaf's code, so its caller lists
//     only commits that change a file outside .rhiza/ (Input.Commits). The
//     finding names the oldest such commit. Without git (Input.Commits nil)
//     the rule is skipped silently.
//   - unclaimed warns once on each commit of Input.Unclaimed (the non-merge
//     commits of <cutoff>..HEAD that change a file outside .rhiza/) none of
//     whose Claims names a change set in the log. The message says whether
//     the Change: trailer is missing or names an unknown change set. A Change:
//     counts from any wholly-trailer paragraph after the subject, so a
//     GitHub squash body is claimed by its branch commits' paragraphs; Actor:,
//     Leaf: and Review: trailers alone claim nothing. Without a cutoff
//     (Input.UnclaimedSince "") the rule is off.
//   - anchor fails once on each live as-is anchor whose path is gone from
//     Input.Tree or whose symbol is gone from its file (a path that is not
//     a clean relative one, a symbol on a non-Go file and a Go file that
//     does not parse count as gone). Each finding names the anchor and the
//     capabilities it implements by live implements links. It warns, with
//     Next rhiza anchor refresh <anchor>, on each anchor that resolves but
//     whose stored fingerprint differs from that of its declaration (or
//     file). Without a Tree the rule is skipped silently.
//   - render fails on each rendered file that differs from what
//     render.Render writes, is missing or is extra, and on a missing
//     rendered folder.
//   - append-only compares log entries by identity, not by path: an entry
//     is its entry ID (its file name minus .json) and its bytes. It fails on
//     each entry of the base whose ID is in the current log nowhere (deleted)
//     or only with other bytes (modified), and warns once when the rule was
//     skipped (no git, no base). An entry found with the same bytes passes
//     wherever it now sits, so a moved log passes; entries new since the
//     base pass.
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
	"maps"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/render"
	"github.com/dphaener/rhiza/internal/store"
)

// The two severities.
const (
	SeverityFail = "fail"
	SeverityWarn = "warn"
)

// The rules, in report order.
const (
	RuleLog           = "log"
	RuleFold          = "fold"
	RuleModel         = "model"
	RuleInFlight      = "in-flight"
	RuleLink          = "link"
	RuleDependsOn     = "depends-on"
	RuleEvidenceActor = "evidence-actor"
	RuleUnclaimed     = "unclaimed"
	RuleAnchor        = "anchor"
	RuleRender        = "render"
	RuleAppendOnly    = "append-only"
)

// Rules returns the rules in report order.
func Rules() []string {
	return []string{RuleLog, RuleFold, RuleModel, RuleInFlight, RuleLink, RuleDependsOn, RuleEvidenceActor, RuleUnclaimed, RuleAnchor,
		RuleRender, RuleAppendOnly}
}

// Paths in findings are relative to the project root.
const (
	rhizaDir    = ".rhiza/"
	renderedDir = ".rhiza/rendered"
)

// Input is everything Run checks.
type Input struct {
	// Log is the log as store.ReadLog read it.
	Log store.Log
	// Rendered is rooted at .rhiza/rendered; nil when the folder is
	// missing.
	Rendered fs.FS
	// Rev is the base revision used, "" when the append-only rule was
	// skipped.
	Rev string
	// Skipped is why the append-only rule was skipped, "" otherwise.
	Skipped string
	// Base are the log entries at Rev, Current those in the work tree.
	Base, Current []Entry
	// Tree is the working tree rooted at the project; nil means the anchor
	// rule is skipped.
	Tree fs.FS
	// Commits are the commits of HEAD's history that change a file outside
	// .rhiza/, newest first as git log lists them; a commit of model writes
	// alone is not listed. Nil means no git, and the evidence-actor rule is
	// skipped.
	Commits []Commit
	// UnclaimedSince is the full hash of the unclaimed-changes cutoff from
	// .rhiza/check.json, "" when the rule is off.
	UnclaimedSince string
	// Unclaimed are the non-merge commits of UnclaimedSince..HEAD that change
	// a file outside .rhiza/, newest first.
	Unclaimed []Commit
}

// Commit is one git commit and the values of its trailers: Changes of its
// Change: trailers (change set IDs), Actors of its Actor: trailers (actor
// handles or IDs) and Leaves of its Leaf: trailers (capability handles or
// IDs), and Reviews of its Review: trailers (the run IDs of the reviews
// whose change-set fixer made it). Hash is the full commit hash.
type Commit struct {
	Hash                    string
	Changes, Actors, Leaves []string
	Reviews                 []string
	// Subject is the first line of the message. Claims are the values of the
	// Change: trailers of every wholly-trailer paragraph after the subject.
	Subject string
	Claims  []string
}

// Entry is one log entry file: its entry ID (the file name minus .json),
// its path relative to the project root, such as
// .rhiza/log/<change set>/<entry>.json, and the git blob ID of its bytes.
// Two entries with the same Blob have the same bytes.
type Entry struct {
	ID   string
	Path string
	Blob string
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
// is Input.Rev. Findings is never nil.
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
	c.unlinkedConstraints(asIs)
	c.inFlightLinks(sets)
	c.dependsOn(sets)
	c.capabilityLoops(asIs)
	c.evidenceActor(asIs, sets)
	c.unclaimed(sets)
	c.anchors(asIs)
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
	return Report{OK: ok, Base: in.Rev, Findings: c.findings}
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
		c.add(Finding{Severity: SeverityFail, Rule: RuleLog, Path: rhizaDir + p.Path, Message: p.Err.Error()})
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
			f.Path = fmt.Sprintf("%s%s/%s/%s.json", rhizaDir, store.LogDir, p.ChangeSet, p.Entry)
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
			f.Next = fmt.Sprintf("rhiza change new %q", "Resolve the conflict on "+target(sets, p))
		}
		c.add(f)
	}

	_, toBe := model.ToBe(sets)
	views := [][]model.Problem{toBe}
	var inFlight []changeset.ChangeSet
	for _, s := range sets {
		switch s.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
			_, ps := model.WhatIf(sets, s.ID)
			views = append(views, ps)
			inFlight = append(inFlight, s)
		}
	}
	for _, ps := range views {
		for _, p := range ps {
			k := keyOf(p)
			if seen[k] {
				continue
			}
			seen[k] = true
			f := Finding{Severity: SeverityWarn, Rule: RuleInFlight, ChangeSet: str(p.ChangeSet), Op: str(p.Op),
				Message: p.Err.Error()}
			if errors.Is(p.Err, model.ErrConflict) && !p.ChangeSet.IsZero() {
				f.Next = fmt.Sprintf("rhiza slice --change %s --view what-if", p.ChangeSet)
			}
			c.add(f)
		}
	}
	for _, s := range inFlight {
		for _, x := range changeset.RestsOn(sets, s.ID) {
			for _, o := range x.Ops {
				c.add(Finding{Severity: SeverityWarn, Rule: RuleInFlight, ChangeSet: s.ID.String(), Op: o.String(),
					Message: fmt.Sprintf("op rests on change set %s (%s)", x.ChangeSet, x.State),
					Next:    "rhiza change show " + s.ID.String()})
			}
		}
	}
	for _, ov := range model.Overlaps(sets) {
		c.add(Finding{Severity: SeverityWarn, Rule: RuleInFlight, ChangeSet: ov.A.String(),
			Message: fmt.Sprintf("change sets %s and %s both write %s of %s", ov.A, ov.B,
				strings.Join(ov.Aspects, ", "), ov.Target),
			Next: fmt.Sprintf("rhiza slice --change %s --view what-if", ov.A)})
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
	for _, l := range badLinks(m, sets) {
		c.add(linkFinding(SeverityFail, RuleLink, l))
	}
}

// inFlightLinks warns on each live link, in the what-if of a draft,
// proposed or accepted change set, that the set itself wrote and whose Ends
// the endpoint table does not allow, each once. As-is links are the link
// rule's.
func (c *checker) inFlightLinks(sets []changeset.ChangeSet) {
	seen := map[id.ID]bool{}
	for _, s := range sets {
		switch s.State {
		case changeset.StateDraft, changeset.StateProposed, changeset.StateAccepted:
		default:
			continue
		}
		whatIf, _ := model.WhatIf(sets, s.ID)
		for _, l := range badLinks(whatIf, sets) {
			if l.Created != s.ID || seen[l.ID] {
				continue
			}
			seen[l.ID] = true
			c.add(linkFinding(SeverityWarn, RuleInFlight, l))
		}
	}
}

// badLink is a live link of a model whose Ends its type does not allow,
// with the names of its Ends.
type badLink struct {
	model.Link
	fromName, toName string
}

func badLinks(m model.Model, sets []changeset.ChangeSet) []badLink {
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
	var bad []badLink
	for _, l := range m.Links() {
		if l.Retired {
			continue
		}
		from, fromName := end(l.From)
		to, toName := end(l.To)
		if !Allowed(l.Type, from, to) {
			bad = append(bad, badLink{l, fromName, toName})
		}
	}
	return bad
}

func linkFinding(severity, rule string, l badLink) Finding {
	f := Finding{Severity: severity, Rule: rule, ChangeSet: str(l.Created), Link: l.ID.String(),
		Message: fmt.Sprintf("link type %q does not go from %s to %s", l.Type, l.fromName, l.toName)}
	if len(l.History) > 0 {
		f.Op = str(l.History[0].Op)
	}
	return f
}

func (c *checker) dependsOn(sets []changeset.ChangeSet) {
	for _, loop := range changeset.Cycles(sets) {
		msg := fmt.Sprintf("change set %s depends on itself", loop[0])
		if len(loop) > 1 {
			names := make([]string, len(loop))
			for i, cs := range loop {
				names[i] = cs.String()
			}
			msg = "depends-on loop among change sets " + strings.Join(names, ", ")
		}
		c.add(Finding{Severity: SeverityFail, Rule: RuleDependsOn, ChangeSet: loop[0].String(), Message: msg,
			Next: "rhiza change show " + loop[0].String()})
	}
}

// capabilityLoops fails once on each depends-on loop among capabilities:
// each group of two or more that reach one another over live as-is links,
// and each that depends on itself, named by handle in handle order.
func (c *checker) capabilityLoops(m model.Model) {
	g := map[id.ID][]id.ID{}
	for _, l := range m.Links() {
		if l.Retired || l.Type != op.LinkDependsOn {
			continue
		}
		from, okFrom := m.Node(l.From)
		to, okTo := m.Node(l.To)
		if okFrom && okTo && from.Kind == op.KindCapability && to.Kind == op.KindCapability {
			g[l.From] = append(g[l.From], l.To)
		}
	}
	for _, loop := range loops(g) {
		names := make([]string, len(loop))
		for i, n := range loop {
			names[i] = n.String()
			if v, _ := m.Node(n); v.Handle != "" {
				names[i] = v.Handle
			}
		}
		slices.Sort(names)
		msg := "capability " + names[0] + " depends on itself"
		if len(names) > 1 {
			msg = "depends-on loop among capabilities " + strings.Join(names, ", ")
		}
		c.add(Finding{Severity: SeverityFail, Rule: RuleDependsOn, Message: msg, Next: "rhiza show " + names[0]})
	}
}

// evidenceActor warns on each Evidence whose actor wrote the code of a
// capability it verifies (see the package comment), once per Evidence and
// capability, naming the oldest matching commit by its first 12 hex digits.
func (c *checker) evidenceActor(m model.Model, sets []changeset.ChangeSet) {
	if len(c.in.Commits) == 0 {
		return
	}
	name := func(n model.Node) string {
		if n.Handle != "" {
			return n.Handle
		}
		return n.ID.String()
	}
	type pair struct{ ev, capab id.ID }
	seen := map[pair]bool{}
	for _, l := range m.Links() {
		if l.Retired || l.Type != op.LinkVerifies || seen[pair{l.From, l.To}] {
			continue
		}
		ev, okEv := m.Node(l.From)
		capab, okCap := m.Node(l.To)
		if !okEv || !okCap || ev.Retired || ev.Kind != op.KindEvidence || capab.Kind != op.KindCapability {
			continue
		}
		who, ok := model.CreatedBy(ev, sets)
		if !ok {
			continue
		}
		actorNames := []string{who.String()}
		actorName := who.String()
		if a, ok := m.Node(who); ok && a.Handle != "" {
			actorNames = append(actorNames, a.Handle)
			actorName = a.Handle
		}
		var hit *Commit
		for i := range c.in.Commits {
			cm := &c.in.Commits[i]
			if slices.Contains(cm.Changes, capab.Created.String()) &&
				(len(cm.Leaves) == 0 || (capab.Handle != "" && slices.Contains(cm.Leaves, capab.Handle)) ||
					slices.Contains(cm.Leaves, capab.ID.String())) &&
				slices.ContainsFunc(cm.Actors, func(a string) bool { return slices.Contains(actorNames, a) }) {
				hit = cm
			}
		}
		if hit == nil {
			continue
		}
		seen[pair{l.From, l.To}] = true
		hash := hit.Hash
		if len(hash) > 12 {
			hash = hash[:12]
		}
		c.add(Finding{Severity: SeverityWarn, Rule: RuleEvidenceActor, ChangeSet: str(capab.Created), Link: l.ID.String(),
			Message: fmt.Sprintf("evidence %s verifies capability %s, but its actor %s wrote that capability in commit %s",
				name(ev), name(capab), actorName, hash),
			Next: "rhiza show " + name(ev)})
	}
}

// unclaimed warns on each commit of Input.Unclaimed that no known change set
// claims (see the package comment).
func (c *checker) unclaimed(sets []changeset.ChangeSet) {
	if c.in.UnclaimedSince == "" {
		return
	}
	known := map[string]bool{}
	for _, s := range sets {
		known[s.ID.String()] = true
	}
	for _, cm := range c.in.Unclaimed {
		if slices.ContainsFunc(cm.Claims, func(v string) bool { return known[v] }) {
			continue
		}
		hash := cm.Hash
		if len(hash) > 12 {
			hash = hash[:12]
		}
		why := "has no Change: trailer (missing)"
		if len(cm.Claims) > 0 {
			why = fmt.Sprintf("names unknown change set %s in its Change: trailer", strings.Join(cm.Claims, ", "))
		}
		c.add(Finding{Severity: SeverityWarn, Rule: RuleUnclaimed,
			Message: fmt.Sprintf("commit %s %q changes code outside .rhiza/ but %s", hash, cm.Subject, why)})
	}
}

// loops returns the loops of graph g: each strongly connected component of
// two or more nodes, and each node with an edge to itself. It visits nodes
// and edges in ID order, so the result is the same for the same g.
func loops(g map[id.ID][]id.ID) [][]id.ID {
	nodes := slices.SortedFunc(maps.Keys(g), id.Compare)
	for _, n := range nodes {
		slices.SortFunc(g[n], id.Compare)
	}
	// Tarjan's strongly connected components.
	index := map[id.ID]int{}
	low := map[id.ID]int{}
	onStack := map[id.ID]bool{}
	var stack []id.ID
	var out [][]id.ID
	var visit func(v id.ID)
	visit = func(v id.ID) {
		index[v] = len(index)
		low[v] = index[v]
		stack = append(stack, v)
		onStack[v] = true
		for _, w := range g[v] {
			if _, ok := index[w]; !ok {
				visit(w)
				low[v] = min(low[v], low[w])
			} else if onStack[w] {
				low[v] = min(low[v], index[w])
			}
		}
		if low[v] != index[v] {
			return
		}
		var scc []id.ID
		for {
			w := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			onStack[w] = false
			scc = append(scc, w)
			if w == v {
				break
			}
		}
		if len(scc) > 1 || slices.Contains(g[v], v) {
			out = append(out, scc)
		}
	}
	for _, n := range nodes {
		if _, ok := index[n]; !ok {
			visit(n)
		}
	}
	return out
}

func (c *checker) render(m model.Model, sets []changeset.ChangeSet) {
	const next = "rhiza render"
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

// appendOnly judges each base entry by its ID and bytes (see the package
// comment). Of several current entries with a base entry's ID, one with the
// same bytes passes it; else the one at the base path, else the first by
// path, is the modified one. Next restores the base bytes and, when the base
// path lies outside the store folder (history from before the rename), moves
// them where the entry belongs now.
func (c *checker) appendOnly() {
	if c.in.Skipped != "" {
		c.add(Finding{Severity: SeverityWarn, Rule: RuleAppendOnly,
			Message: "append-only check skipped: " + c.in.Skipped})
	}
	byID := map[string][]Entry{}
	for _, e := range c.in.Current {
		byID[e.ID] = append(byID[e.ID], e)
	}
	for _, cur := range byID {
		slices.SortFunc(cur, func(a, b Entry) int { return strings.Compare(a.Path, b.Path) })
	}
	seen := map[Entry]bool{}
	for _, b := range c.in.Base {
		if seen[b] {
			continue
		}
		seen[b] = true
		cur := byID[b.ID]
		if slices.ContainsFunc(cur, func(e Entry) bool { return e.Blob == b.Blob }) {
			continue
		}
		moved := !strings.HasPrefix(b.Path, rhizaDir)
		f := Finding{Severity: SeverityFail, Rule: RuleAppendOnly, Entry: b.ID}
		if len(cur) == 0 {
			f.Path = b.Path
			f.Message = fmt.Sprintf("log entry %s deleted since the base; entries may only be added", b.ID)
			if c.in.Rev != "" {
				f.Next = fmt.Sprintf("git checkout %s -- %s", c.in.Rev, b.Path)
				if moved {
					f.Next += fmt.Sprintf(" && git mv %s %s", b.Path, underStore(b.Path))
				}
			}
			c.add(f)
			continue
		}
		at := cur[0]
		if i := slices.IndexFunc(cur, func(e Entry) bool { return e.Path == b.Path }); i >= 0 {
			at = cur[i]
		}
		f.Path = at.Path
		f.Message = fmt.Sprintf("log entry %s modified since the base; entries may only be added", b.ID)
		if at.Path != b.Path {
			f.Message = fmt.Sprintf("log entry %s modified since the base, where it was %s; entries may only be added",
				b.ID, b.Path)
		}
		if c.in.Rev != "" {
			f.Next = fmt.Sprintf("git checkout %s -- %s", c.in.Rev, b.Path)
			if moved {
				f.Next += fmt.Sprintf(" && git mv -f %s %s", b.Path, at.Path)
			}
		}
		c.add(f)
	}
}

// underStore is path p with its top folder replaced by the store folder.
func underStore(p string) string {
	if _, rest, ok := strings.Cut(p, "/"); ok {
		return rhizaDir + rest
	}
	return rhizaDir + p
}
