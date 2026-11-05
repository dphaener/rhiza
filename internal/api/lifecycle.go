package api

// This file holds the lifecycle operations: propose, accept, reject,
// withdraw, land and reopen. Each writes one entry of its type into --change
// or the current change set, after checking the move against the folded
// state with the same table changeset.Fold applies. A terminal move (reject,
// withdraw, land) on the current change set clears it; a reopen makes its
// change set current.
//
// Land checks two things more, before anything is written, so a dry run
// refuses the same: every change set it depends on has landed, and its
// what-if view has no conflict on its own ops. Either refusal exits 3.
//
// Accept and reject under rule autonomous are allowed only inside an
// autonomous run: RHIZA_RUN must name a run that is unfinished and was
// started --autonomous, or the move is refused as usage before anything is
// written. With no --reason the reason names that run. Once the log entry
// is written, the run's journal gets an accepted or rejected entry naming
// the change set and the log entry, so the run lists its own gate moves.
//
// Reopen moves a landed change set back to accepted, only while its land is
// on its branch alone: none of its land entries is in the log at check's
// base. Landed work on main is final. With no base to look at it refuses
// too. Both refusals exit 3, before anything is written.

import (
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/entry"
	"github.com/dphaener/rhiza/internal/id"
	gmodel "github.com/dphaener/rhiza/internal/model"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	for _, v := range []lifecycleVerb{
		{name: "propose", summary: "Propose the change set for review", done: "Proposed"},
		{name: "accept", summary: "Accept a proposed change set", done: "Accepted", decision: true},
		{name: "reject", summary: "Reject a proposed change set", done: "Rejected", decision: true},
		{name: "withdraw", summary: "Withdraw a change set that has not landed", done: "Withdrew"},
		{name: "land", summary: "Land the change set: an accepted intention, or a record", done: "Landed"},
		{name: "reopen", summary: "Reopen a change set landed only on its branch, so it takes writes and lands again", done: "Reopened"},
	} {
		req := reflect.TypeFor[LifecycleRequest]()
		switch {
		case v.decision:
			req = reflect.TypeFor[DecisionRequest]()
		case v.name == "reopen":
			req = reflect.TypeFor[ReopenRequest]()
		}
		register(Operation{
			Name:     v.name,
			Summary:  v.summary,
			Write:    true,
			Request:  req,
			Response: reflect.TypeFor[LifecycleResponse](),
			Run:      v.run,
		})
	}
}

// defaultRule is the rule an accept or reject is made under when none is
// given.
const defaultRule = "owner"

// autonomousRule is the rule of a gate an autonomous run moved past without
// the owner. It is allowed only while RHIZA_RUN names an unfinished
// autonomous run, so it always means such a run made the move.
const autonomousRule = "autonomous"

// autonomyRule is the rule of an accept the gate of an owner-run plan made
// because rhiza.autonomy.accept is auto. It is allowed only while RHIZA_RUN
// names an unfinished run, autonomous or not.
const autonomyRule = "autonomy"

// autonomousRuleWhy is why rule autonomous is refused outside an
// autonomous run, for its messages.
const autonomousRuleWhy = "it is only for a move made inside a run started with rhiza run start --autonomous"

// LifecycleRequest is the request of propose, withdraw and land: they act
// on --change or the current change set and take nothing else.
type LifecycleRequest struct{}

// DecisionRequest is the request of accept and reject.
type DecisionRequest struct {
	Rule   string `json:"rule,omitempty" desc:"The rule the decision is made under; default owner"`
	Reason string `json:"reason,omitempty" desc:"Why"`
}

// ReopenRequest is the request of reopen: the base its land must not be in.
type ReopenRequest struct {
	Base string `json:"base,omitempty" desc:"The git revision the change set's land must not be in; default the merge-base of HEAD with origin/HEAD, origin/main or main"`
}

// reopenNext is the next command when reopen is refused: the work goes on
// in a new change set.
const reopenNext = `rhiza change new "<intent>" --json`

// LifecycleResponse reports one lifecycle move.
type LifecycleResponse struct {
	Verb           string          `json:"-"`
	ChangeSet      id.ID           `json:"change_set" desc:"The change set moved"`
	Entry          id.ID           `json:"entry,omitzero" desc:"The entry written; absent in a dry run"`
	From           changeset.State `json:"from" desc:"The state before the move"`
	To             changeset.State `json:"to" desc:"The state after the move"`
	Rule           string          `json:"rule,omitempty" desc:"The rule of an accept or reject"`
	Reason         string          `json:"reason,omitempty" desc:"The reason of an accept or reject"`
	ClearedCurrent bool            `json:"cleared_current,omitempty" desc:"True if the move ended the current change set and cleared it"`
	DryRun         bool            `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// lifecycleVerb is one lifecycle operation.
type lifecycleVerb struct {
	name, summary string
	// done is the past tense for human output.
	done string
	// decision is true for accept and reject.
	decision bool
}

// body builds the entry body for the verb.
func (v lifecycleVerb) body(rule, reason string) changeset.Body {
	switch v.name {
	case "propose":
		return changeset.Propose{}
	case "accept":
		return changeset.Accept{Rule: rule, Reason: reason}
	case "reject":
		return changeset.Reject{Rule: rule, Reason: reason}
	case "withdraw":
		return changeset.Withdraw{}
	case "reopen":
		return changeset.Reopen{}
	}
	return changeset.Land{}
}

// to is the state the verb moves c to, or false if the move is illegal.
// It mirrors the lifecycle table of changeset.Fold.
func (v lifecycleVerb) to(c changeset.ChangeSet) (changeset.State, bool) {
	inFlight := c.State == changeset.StateDraft || c.State == changeset.StateProposed ||
		c.State == changeset.StateAccepted
	switch v.name {
	case "propose":
		return changeset.StateProposed, c.State == changeset.StateDraft
	case "accept":
		return changeset.StateAccepted, c.State == changeset.StateProposed
	case "reject":
		return changeset.StateRejected, c.State == changeset.StateProposed
	case "withdraw":
		return changeset.StateWithdrawn, inFlight
	case "reopen":
		return changeset.StateAccepted, c.State == changeset.StateLanded
	}
	if c.Open.Kind == changeset.KindRecord {
		return changeset.StateLanded, inFlight
	}
	return changeset.StateLanded, c.State == changeset.StateAccepted
}

// legalNext is the move that can be made on c next.
func legalNext(c changeset.ChangeSet) string {
	switch {
	case c.State.Terminal() || c.State == "":
		return newEntryNext
	case c.Open.Kind == changeset.KindRecord || c.State == changeset.StateAccepted:
		return "rhiza land"
	case c.State == changeset.StateDraft:
		return "rhiza propose"
	}
	return "rhiza accept"
}

// run checks the move and writes its entry.
func (v lifecycleVerb) run(ctx *Ctx, req any) (any, error) {
	var rule, reason, base string
	switch r := req.(type) {
	case DecisionRequest:
		rule, reason = strings.TrimSpace(r.Rule), r.Reason
		if rule == "" {
			rule = defaultRule
		}
	case ReopenRequest:
		base = strings.TrimSpace(r.Base)
	}
	var auto *autonomousGate
	if v.decision && (rule == autonomousRule || rule == autonomyRule) {
		g, err := checkAutonomousRule(ctx, v.name, rule)
		if err != nil {
			return nil, err
		}
		auto = &g
		if strings.TrimSpace(reason) == "" {
			reason = "autonomous run " + g.run
			if rule == autonomyRule {
				reason = "rhiza.autonomy.accept is auto, run " + g.run
			}
		}
	}
	body := v.body(rule, reason)
	if err := body.Validate(); err != nil {
		return nil, Usagef("%s", err.Error()).WithNext("rhiza describe " + v.name)
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	cs, err := w.effectiveChangeSet(ctx)
	if err != nil {
		return nil, err
	}
	if cs.IsZero() {
		return nil, Usagef("no change set: pass --change or make one current with change use").
			WithNext(changeLsNext)
	}
	c, ok := w.changeSet(cs)
	if !ok {
		return nil, NotFoundf("no change set %s", cs).WithNext(changeLsNext)
	}
	to, legal := v.to(c)
	if !legal {
		return nil, Conflictf("cannot %s change set %s: it is %s", v.name, cs, c.State).
			WithNext(legalNext(c))
	}
	switch v.name {
	case "land":
		if err := checkLand(w.changeSets(), cs); err != nil {
			return nil, err
		}
	case "reopen":
		if err := checkReopen(ctx, c, base); err != nil {
			return nil, err
		}
	}
	actor, forID, err := w.resolveActor(cs, ctx)
	if err != nil {
		return nil, err
	}
	e, err := w.appendEntry(cs, body, actor, forID)
	if err != nil {
		return nil, err
	}
	resp := LifecycleResponse{Verb: v.done, ChangeSet: cs, Entry: e.ID, From: c.State, To: to,
		Rule: rule, Reason: reason, DryRun: ctx.DryRun}
	if ctx.DryRun {
		resp.Entry = id.ID{}
	}
	if to.Terminal() && cs == current && !ctx.DryRun {
		if err := w.clearCurrent(); err != nil {
			return nil, err
		}
		resp.ClearedCurrent = true
	}
	if v.name == "reopen" && !ctx.DryRun {
		if err := w.setCurrent(cs); err != nil {
			return nil, err
		}
	}
	if auto != nil && !ctx.DryRun {
		if err := auto.journal(ctx, v.name, cs, e.ID); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// autonomousGate is the run an accept or reject under rule autonomous is
// made in, and the store its journal is in.
type autonomousGate struct {
	store *run.Store
	run   string
}

// checkAutonomousRule finds the run RHIZA_RUN names and refuses, as
// usage, unless it is an unfinished autonomous run.
func checkAutonomousRule(ctx *Ctx, verb, rule string) (autonomousGate, error) {
	next := "rhiza describe " + verb
	why := autonomousRuleWhy
	if rule == autonomyRule {
		why = "it is only for a move made inside a run"
	}
	ref := strings.TrimSpace(ctxGetenv(ctx)(envRun))
	if ref == "" {
		return autonomousGate{}, Usagef("cannot %s under rule %s: %s, and %s is not set",
			verb, rule, why, envRun).WithNext(next)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return autonomousGate{}, err
	}
	rs := run.NewStore(run.Dir(common))
	j, err := rs.Read(ref)
	var nf *run.NotFoundError
	var amb *run.AmbiguousError
	switch {
	case errors.As(err, &nf), errors.As(err, &amb):
		return autonomousGate{}, Usagef("cannot %s under rule %s: %s, and %s=%s names no run: %v", verb, rule, why, envRun, ref, err).
			WithNext(next)
	case err != nil:
		return autonomousGate{}, Errorf("read run %s: %v", ref, err)
	}
	st := j.State()
	if !st.Autonomous && rule == autonomousRule {
		return autonomousGate{}, Usagef("cannot %s under rule %s: %s, and run %s is not autonomous", verb, rule, why, j.ID).
			WithNext(next)
	}
	if st.Ended {
		return autonomousGate{}, Usagef("cannot %s under rule %s: %s, and run %s is finished: it %s", verb, rule, why, j.ID, st.Result).
			WithNext(next)
	}
	return autonomousGate{store: rs, run: j.ID}, nil
}

// journal appends the move to the run's journal: accepted or rejected,
// with the change set and the log entry written.
func (g autonomousGate) journal(ctx *Ctx, verb string, cs, logEntry id.ID) error {
	typ := run.TypeAccepted
	if verb == "reject" {
		typ = run.TypeRejected
	}
	e := run.Entry{Type: typ, At: ctx.Now(), Change: cs.String(), LogEntry: logEntry.String()}
	if err := g.store.Append(g.run, e); err != nil {
		return Errorf("journal the %s of change set %s in run %s (log entry %s is written): %v",
			verb, cs, g.run, logEntry, err).WithNext("rhiza run status " + g.run)
	}
	return nil
}

// checkReopen refuses to reopen c while any of its land entries is in the
// log at the base (base, else check's default), or when there is no base to
// look at, since then nothing shows the land is only on this branch.
func checkReopen(ctx *Ctx, c changeset.ChangeSet, base string) error {
	rev, _, skipped, err := checkBase(ctx.Git, base)
	if err != nil {
		return err
	}
	if rev == "" {
		return Conflictf("cannot reopen change set %s: no base to show its land is only on this branch (%s)",
			c.ID, skipped).WithNext(reopenNext)
	}
	es, err := logAt(ctx, rev)
	if err != nil {
		return err
	}
	inBase := make(map[string]bool, len(es))
	for _, e := range es {
		inBase[e.ID] = true
	}
	for _, m := range c.Lifecycle {
		if m.Type != entry.TypeLand || !inBase[m.Entry.String()] {
			continue
		}
		short := shortID(c.ID)
		return Conflictf("cannot reopen change set %s: it landed on the base (%s), and landed work on main is final; "+
			"open a new change set and link it: rhiza link <new> discovered-from %s --change <new>",
			c.ID, shortRev(rev), short).WithNext(reopenNext)
	}
	return nil
}

// shortRev shortens a commit ID for messages; a ref name is kept whole.
func shortRev(rev string) string {
	if len(rev) == 40 && strings.Trim(rev, "0123456789abcdef") == "" {
		return rev[:12]
	}
	return rev
}

// checkLand refuses to land cs while a change set it depends on has not
// landed, naming the first by ID, or while cs's what-if view has a problem
// on cs itself or one of its own ops: a conflict (next rhiza check) or
// any other, such as an op resting on a node that is gone (next the what-if
// slice), so nothing lands that check would then fail.
func checkLand(sets []changeset.ChangeSet, cs id.ID) error {
	states := make(map[id.ID]changeset.State, len(sets))
	for _, s := range sets {
		states[s.ID] = s.State
	}
	for _, dep := range changeset.DependsOn(sets, cs) {
		state, ok := states[dep]
		if ok && state == changeset.StateLanded {
			continue
		}
		what := string(state)
		if !ok {
			what = "not in the log"
		}
		next := "rhiza change show " + shortID(dep)
		if state == changeset.StateAccepted {
			next = "rhiza land --change " + shortID(dep)
		}
		return Conflictf("cannot land change set %s: it depends on %s (%s)", cs, dep, what).WithNext(next)
	}
	_, problems := gmodel.WhatIf(sets, cs)
	for _, p := range problems {
		if p.ChangeSet != cs {
			continue
		}
		next := "rhiza slice --change " + shortID(cs) + " --view what-if"
		if errors.Is(p.Err, gmodel.ErrConflict) {
			next = "rhiza check"
		}
		return Conflictf("cannot land change set %s: %v", cs, p.Err).WithNext(next)
	}
	return nil
}

// Human reports the move.
func (r LifecycleResponse) Human(w io.Writer) {
	if r.DryRun {
		fmt.Fprintf(w, "Would move change set %s (%s → %s).\n\n", r.ChangeSet, r.From, r.To)
		fmt.Fprintf(w, "This was a dry run: nothing was written and the current change set is unchanged.\n")
		return
	}
	verb := r.Verb
	if verb == "" {
		verb = "Moved"
	}
	fmt.Fprintf(w, "%s change set %s (%s → %s).\n", verb, r.ChangeSet, r.From, r.To)
	if r.Rule != "" {
		fmt.Fprintf(w, "Rule: %s\n", r.Rule)
	}
	if r.Reason != "" {
		fmt.Fprintf(w, "Reason: %s\n", r.Reason)
	}
	if r.ClearedCurrent {
		fmt.Fprintln(w, "No current change set now.")
	}
	switch r.To {
	case changeset.StateProposed:
		fmt.Fprintln(w, "next: rhiza accept")
	case changeset.StateAccepted:
		fmt.Fprintln(w, "next: rhiza land")
	}
}
