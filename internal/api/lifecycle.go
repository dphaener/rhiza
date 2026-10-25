package api

// This file holds the lifecycle operations: propose, accept, reject,
// withdraw and land. Each writes one entry of its type into --change or the
// current change set, after checking the move against the folded state with
// the same table changeset.Fold applies. A terminal move (reject, withdraw,
// land) on the current change set clears it.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/genesis/internal/changeset"
	"github.com/dphaener/genesis/internal/id"
)

func init() {
	for _, v := range []lifecycleVerb{
		{name: "propose", summary: "Propose the change set for review", done: "Proposed"},
		{name: "accept", summary: "Accept a proposed change set", done: "Accepted", decision: true},
		{name: "reject", summary: "Reject a proposed change set", done: "Rejected", decision: true},
		{name: "withdraw", summary: "Withdraw a change set that has not landed", done: "Withdrew"},
		{name: "land", summary: "Land the change set: an accepted intention, or a record", done: "Landed"},
	} {
		req := reflect.TypeFor[LifecycleRequest]()
		if v.decision {
			req = reflect.TypeFor[DecisionRequest]()
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

// LifecycleRequest is the request of propose, withdraw and land: they act
// on --change or the current change set and take nothing else.
type LifecycleRequest struct{}

// DecisionRequest is the request of accept and reject.
type DecisionRequest struct {
	Rule   string `json:"rule,omitempty" desc:"The rule the decision is made under; default owner"`
	Reason string `json:"reason,omitempty" desc:"Why"`
}

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
		return "genesis land"
	case c.State == changeset.StateDraft:
		return "genesis propose"
	}
	return "genesis accept"
}

// run checks the move and writes its entry.
func (v lifecycleVerb) run(ctx *Ctx, req any) (any, error) {
	var rule, reason string
	if r, ok := req.(DecisionRequest); ok {
		rule, reason = strings.TrimSpace(r.Rule), r.Reason
		if rule == "" {
			rule = defaultRule
		}
	}
	if ctx.Record {
		return nil, Usagef("%s does not take --record: it is for node and link writes", v.name).
			WithNext("genesis describe " + v.name)
	}
	body := v.body(rule, reason)
	if err := body.Validate(); err != nil {
		return nil, Usagef("%s", err.Error()).WithNext("genesis describe " + v.name)
	}
	w, err := openCtxWorkspace(ctx)
	if err != nil {
		return nil, err
	}
	current, err := w.current()
	if err != nil {
		return nil, err
	}
	cs := current
	if ref := strings.TrimSpace(ctx.ChangeRef); ref != "" {
		if cs, err = w.resolveChangeSet(ref); err != nil {
			return nil, err
		}
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
	actor, forID, err := w.resolveActor(cs, ctx.ActorRef, ctx.ForRef, ctx.GitConfig)
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
	return resp, nil
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
		fmt.Fprintln(w, "next: genesis accept")
	case changeset.StateAccepted:
		fmt.Fprintln(w, "next: genesis land")
	}
}
