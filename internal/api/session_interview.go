package api

// This file holds session.interview: after an agent step, Rhiza reads the
// session's transcript and proposes a Decision in the session's change set
// for each question the agent asked and the owner answered. session.run
// calls the same capture once the agent exits.
//
// Who wrote each word is kept apart. A picked option is the agent's wording,
// accepted by the owner: it is recorded as picked, with a wording note, and
// never as the owner's quote. Only text the owner typed (an answer that is
// not an option, and notes) is quoted as the owner's.
//
// Nothing here names an agent: the session's adapter says where its
// transcript is and reads the answers out of it.

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "session.interview",
		Summary:  "Propose a Decision for each question the session's agent asked and the owner answered",
		Write:    true,
		Request:  reflect.TypeFor[SessionInterviewRequest](),
		Response: reflect.TypeFor[SessionInterviewResponse](),
		Run:      runSessionInterview,
	})
}

// answerHandlePrefix starts the handle of every Decision an answer proposes.
const answerHandlePrefix = "answer-"

// SessionInterviewRequest names the session whose answers to capture.
type SessionInterviewRequest struct {
	ID string `json:"id" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
}

// SessionInterviewResponse lists the Decisions proposed, or in a dry run
// those that would be.
type SessionInterviewResponse struct {
	Session  string             `json:"session" desc:"The session ID"`
	Change   string             `json:"change" desc:"The session's change set, which the Decisions go into while it is open"`
	Record   string             `json:"record,omitempty" desc:"The landed record change set the Decisions were written in, when the session's change set had landed"`
	Proposed []ProposedDecision `json:"proposed" desc:"The Decisions proposed, in transcript order"`
	Skipped  []string           `json:"skipped" desc:"Handles of answers already proposed in the change set, left alone"`
	DryRun   bool               `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// ProposedDecision is one answered question, proposed as a Decision.
type ProposedDecision struct {
	Handle   string   `json:"handle" desc:"The Decision's handle: answer-<first 8 characters of the session ID>-<position in the transcript>"`
	Question string   `json:"question" desc:"The question the agent asked, verbatim"`
	Picked   []string `json:"picked,omitempty" desc:"The agent's option labels the owner picked, verbatim"`
	Typed    string   `json:"typed,omitempty" desc:"The text the owner typed, verbatim"`
}

// answerWrite is the add one answer becomes.
type answerWrite struct {
	answer adapter.Answer
	add    AddRequest
}

// runSessionInterview captures the session's answers as Decisions.
func runSessionInterview(ctx *Ctx, req any) (any, error) {
	r := req.(SessionInterviewRequest)
	ref := strings.TrimSpace(r.ID)
	if ref == "" {
		return nil, Usagef("session interview needs a session ID").WithNext(sessionLsNext)
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	return captureInterview(ctx, rec)
}

// captureInterview reads rec's transcript and proposes one Decision per
// answered question into rec's change set, in the workspace of rec's
// worktree, as rec's agent actor for rec's for-actor, in one entry. Answers
// whose handle the change set already has are skipped, so a second capture
// of the same transcript writes nothing. Ctx.DryRun writes nothing.
//
// When rec's change set has landed, as a milestone planner's record does
// before the planner exits, the Decisions and their derived-from links go
// into one landed record change set over the landed view instead; answers
// whose handle the landed view has are skipped. Any other terminal state is
// a conflict.
func captureInterview(ctx *Ctx, rec session.Record) (SessionInterviewResponse, error) {
	resp := SessionInterviewResponse{Session: rec.ID, Change: rec.Change, Proposed: []ProposedDecision{},
		Skipped: []string{}, DryRun: ctx.DryRun}
	a, ok := adapter.Lookup(rec.Agent)
	if !ok {
		return resp, Errorf("session %s names agent %q, which this rhiza does not know: want one of %s",
			rec.ID, rec.Agent, strings.Join(adapter.Names(), ", "))
	}
	transcript, err := readTranscript(ctx, a, rec.ID)
	if err != nil {
		return resp, err
	}
	answers, err := a.Answers(transcript)
	if err != nil {
		return resp, Errorf("read the answers of session %s: %v", rec.ID, err)
	}
	writes, err := answerWrites(rec.ID, answers)
	if err != nil {
		return resp, err
	}
	if len(writes) == 0 {
		return resp, nil
	}

	ww, cs, err := sessionWorkspace(ctx, rec)
	if err != nil {
		return resp, err
	}
	c, _ := ww.changeSet(cs)
	// A landed change set takes no more ops: the Decisions go into a record
	// write over the landed view, where its Source is too. A milestone
	// planner lands its record before it exits.
	landed := c.State == changeset.StateLanded
	view := cs
	if landed {
		view = id.ID{}
	}
	m := ww.model(view)
	sources := sessionSources(m, c.Ops, rec.ID)
	if landed {
		sources = landedSessionSources(m, sources, rec.ID)
	}
	recorded := recordedQuestions(m, c.Ops, sources, landed)
	var items []applyItem
	for _, w := range writes {
		if answerProposed(m, c.Ops, w.add.Handle) || recorded[w.answer.Question] {
			resp.Skipped = append(resp.Skipped, w.add.Handle)
			continue
		}
		resp.Proposed = append(resp.Proposed, ProposedDecision{Handle: w.add.Handle, Question: w.answer.Question,
			Picked: w.answer.Picked, Typed: w.answer.Typed})
		add := w.add
		items = append(items, applyItem{op: "add", build: func(b writeBase) (edit, error) { return buildAdd(b, add) }})
		for _, src := range sources {
			link := LinkRequest{From: add.Handle, Type: string(op.LinkDerivedFrom), To: src.String()}
			items = append(items, applyItem{op: "link", build: func(b writeBase) (edit, error) { return buildLink(b, link) }})
		}
	}
	if len(items) == 0 {
		return resp, nil
	}
	if c.State.Terminal() && !landed {
		return resp, Conflictf("change set %s is %s; no Decision can be proposed in it", cs, c.State).
			WithNext("rhiza change show " + cs.String())
	}
	actor, err := liveActor(m, rec.Actor)
	if err != nil {
		return resp, Errorf("session %s's actor %s: %s", rec.ID, rec.Actor, AsError(err).Message).
			WithNext("rhiza session show " + rec.ID)
	}
	forID, err := sessionForActor(ctx, m, rec.For)
	if err != nil {
		return resp, Errorf("session %s's for-actor %s: %s", rec.ID, rec.For, AsError(err).Message).
			WithNext("rhiza session show " + rec.ID)
	}
	if landed {
		e, err := buildApply(ww, id.ID{}, ctx.IDs, items)
		if err == nil {
			err = e.check("add")
		}
		if err != nil {
			return resp, err
		}
		intent := fmt.Sprintf("Record the answers of session %s, whose change set %s has landed", rec.ID, cs)
		rc, _, err := ww.writeRecord(intent, e.ops, actor, forID)
		if err != nil {
			return resp, err
		}
		if !ctx.DryRun {
			resp.Record = rc.String()
		}
		return resp, nil
	}
	e, err := buildApply(ww, cs, ctx.IDs, items)
	if err == nil {
		err = e.check("add")
	}
	if err != nil {
		return resp, err
	}
	if _, err := ww.appendEntry(cs, changeset.Ops{Ops: e.ops, Drop: e.drop}, actor, forID); err != nil {
		return resp, err
	}
	return resp, nil
}

// readTranscript reads the first file that matches one of a's transcript
// globs for session sid, in the adapter's order. None is an error (exit 1),
// not not_found: the session exists, its agent's file does not.
func readTranscript(ctx *Ctx, a adapter.Adapter, sid string) ([]byte, error) {
	globs := a.TranscriptGlobs(sid, ctxGetenv(ctx))
	for _, g := range globs {
		matches, err := filepath.Glob(g)
		if err != nil {
			continue
		}
		slices.Sort(matches)
		for _, p := range matches {
			fi, err := os.Stat(p)
			if err != nil || !fi.Mode().IsRegular() {
				continue
			}
			data, err := os.ReadFile(p)
			if err != nil {
				return nil, Errorf("read the transcript of session %s: %v", sid, err)
			}
			return data, nil
		}
	}
	where := "nowhere: the home folder and environment are unknown"
	if len(globs) > 0 {
		where = strings.Join(globs, ", ")
	}
	return nil, Errorf("no transcript of session %s: looked in %s", sid, where).
		WithNext("rhiza session show " + sid)
}

// ctxGetenv is the environment an adapter may consult: Ctx.Getenv, with
// HOME from Ctx.Home when that is set.
func ctxGetenv(ctx *Ctx) func(string) string {
	return func(key string) string {
		if key == "HOME" && ctx.Home != "" {
			return ctx.Home
		}
		if ctx.Getenv == nil {
			return ""
		}
		return ctx.Getenv(key)
	}
}

// sessionWorkspace opens the Rhiza repository in rec's worktree (at the
// same place within the git repository as ctx.Root) and resolves rec's
// change set there.
func sessionWorkspace(ctx *Ctx, rec session.Record) (*workspace, id.ID, error) {
	prefix, err := ctx.Git("rev-parse", "--show-prefix")
	if err != nil {
		return nil, id.ID{}, Errorf("%v", err)
	}
	root := filepath.Join(rec.Worktree, strings.TrimSpace(prefix))
	ww, err := openWorkspace(root, ctx.Now, ctx.IDs, ctx.DryRun)
	if err != nil {
		if e := AsError(err); e != nil && e.Code == CodeNotFound {
			return nil, id.ID{}, NotFoundf("session %s's worktree has no Rhiza repository at %s", rec.ID, root).
				WithNext("rhiza session show " + rec.ID)
		}
		return nil, id.ID{}, err
	}
	cs, err := id.Parse(rec.Change)
	if err != nil {
		return nil, id.ID{}, Errorf("session %s names change set %q: %v", rec.ID, rec.Change, err)
	}
	if _, ok := ww.changeSet(cs); !ok {
		return nil, id.ID{}, NotFoundf("change set %s of session %s is not in the worktree's log at %s", cs, rec.ID, root).
			WithNext("rhiza session show " + rec.ID)
	}
	return ww, cs, nil
}

// answerProposed reports whether the change set's view already has a node
// with handle h: a live one in the model, or one the change set created,
// even if it was retired since, so an answer the owner turned down is not
// proposed again.
func answerProposed(m Model, ops []op.Op, h string) bool {
	if _, ok := m.HandleTaken(h); ok {
		return true
	}
	return slices.ContainsFunc(ops, func(o op.Op) bool { return o.Type == op.TypeNodeCreate && o.Handle == h })
}

// recordedQuestions is the set of questions already answered by a live
// Decision: one the change set's ops created, or, when the change set has
// landed, one derived-from a Source of the session, whose question field
// is the question verbatim. A planner's own Decision for an answer stops
// the capture from proposing a duplicate.
func recordedQuestions(m Model, ops []op.Op, sources []id.ID, landed bool) map[string]bool {
	out := map[string]bool{}
	add := func(n Node) {
		if n.Retired || n.Kind != op.KindDecision {
			return
		}
		if q := fieldText(n.Fields["question"]); q != "" {
			out[q] = true
		}
	}
	for _, o := range ops {
		if o.Type != op.TypeNodeCreate || o.Kind != op.KindDecision {
			continue
		}
		if n, ok := m.Node(o.Node); ok {
			add(n)
		}
	}
	if landed {
		for _, src := range sources {
			for _, l := range m.LinksTo(src, op.LinkDerivedFrom) {
				if n, ok := m.Node(l.From); ok {
					add(n)
				}
			}
		}
	}
	return out
}

// sessionSources lists the Sources the change set's ops create whose
// session field names session sid, in op order: those a source add
// --raw-session of that session wrote, still live in m.
func sessionSources(m Model, ops []op.Op, sid string) []id.ID {
	var out []id.ID
	for _, o := range ops {
		if o.Type != op.TypeNodeCreate || o.Kind != op.KindSource {
			continue
		}
		n, ok := m.Node(o.Node)
		if !ok || n.Retired || fieldText(n.Fields["session"]) != sid || slices.Contains(out, o.Node) {
			continue
		}
		out = append(out, o.Node)
	}
	return out
}

// landedSessionSources adds to found every other live landed Source whose
// session field names session sid, by ID: one added with --record, say.
func landedSessionSources(m Model, found []id.ID, sid string) []id.ID {
	out := slices.Clone(found)
	for _, n := range m.Children(op.KindSource, id.ID{}) {
		if !n.Retired && fieldText(n.Fields["session"]) == sid && !slices.Contains(out, n.ID) {
			out = append(out, n.ID)
		}
	}
	return out
}

// liveActor resolves ref to a live actor node.
func liveActor(m Model, ref string) (id.ID, error) {
	n, err := resolveActorRef(m, ref)
	if err != nil {
		return id.ID{}, err
	}
	return actorByID(m, n)
}

// sessionForActor resolves the session's for-actor, the launcher, to a live
// actor: in the worktree's model m, else, when m has no such node, in the
// model of ctx's own repository, where session start resolved it. A
// launcher recorded since the worktree's commit is only there, and session
// start already writes the agent's actor in the worktree for it.
func sessionForActor(ctx *Ctx, m Model, ref string) (id.ID, error) {
	n, err := liveActor(m, ref)
	if e := AsError(err); e == nil || e.Code != CodeNotFound {
		return n, err
	}
	w, werr := openCtxWorkspace(ctx)
	if werr != nil {
		return id.ID{}, err
	}
	return liveActor(w.model(id.ID{}), ref)
}

// answerWrites turns the answers of session sid, in transcript order, into
// one decision add each. The handle is answer-<sid8>-<n>, n the answer's
// 1-based position. Picked labels are the agent's wording the owner
// accepted, never a quote; typed text and notes are the owner's, quoted
// verbatim.
func answerWrites(sid string, answers []adapter.Answer) ([]answerWrite, error) {
	if len(sid) < 8 {
		return nil, Errorf("session ID %q is shorter than 8 characters", sid)
	}
	prefix := answerHandlePrefix + strings.ToLower(sid[:8]) + "-"
	out := make([]answerWrite, 0, len(answers))
	for i, a := range answers {
		parts := slices.Clone(a.Picked)
		if a.Typed != "" {
			parts = append(parts, a.Typed)
		}
		fields := []string{
			"statement=" + a.Question + " — " + strings.Join(parts, ", "),
			"question=" + a.Question,
			"session=" + sid,
		}
		if len(a.Picked) > 0 {
			fields = append(fields,
				"picked="+strings.Join(a.Picked, ", "),
				"wording=the agent's option, picked by the owner in session "+sid)
		}
		if a.Typed != "" {
			fields = append(fields,
				"quote="+a.Typed,
				"quote_where=session "+sid+`, answering "`+a.Question+`"`)
		}
		if a.Notes != "" {
			fields = append(fields, "notes_quote="+a.Notes)
		}
		var jsonFields []string
		if len(a.Options) > 0 {
			raw, err := json.Marshal(a.Options)
			if err != nil {
				return nil, Errorf("options of %q: %v", a.Question, err)
			}
			jsonFields = append(jsonFields, "options="+string(raw))
		}
		out = append(out, answerWrite{answer: a, add: AddRequest{
			Kind: string(op.KindDecision), Handle: prefix + strconv.Itoa(i+1), Fields: fields, JSONFields: jsonFields, lenient: true,
		}})
	}
	return out, nil
}

// Human lists the Decisions proposed and the answers skipped.
func (r SessionInterviewResponse) Human(w io.Writer) {
	verb := "Proposed"
	if r.DryRun {
		verb = "Would propose"
	}
	if len(r.Proposed) == 0 {
		fmt.Fprintf(w, "No new answered question in session %s: nothing to propose in change set %s.\n", r.Session, r.Change)
	} else if r.Record != "" {
		fmt.Fprintf(w, "Recorded %d decision(s) from session %s in record change set %s: change set %s has landed.\n\n",
			len(r.Proposed), r.Session, r.Record, r.Change)
	} else {
		fmt.Fprintf(w, "%s %d decision(s) in change set %s from session %s:\n\n", verb, len(r.Proposed), r.Change, r.Session)
	}
	for _, p := range r.Proposed {
		fmt.Fprintf(w, "- %s: %s\n", p.Handle, p.Question)
	}
	for _, h := range r.Skipped {
		fmt.Fprintf(w, "- Skipped %s: already recorded\n", h)
	}
	if r.DryRun {
		fmt.Fprintf(w, "\nThis was a dry run: nothing was written.\n")
		return
	}
	fmt.Fprintf(w, "\nnext: rhiza change show %s\n", r.Change)
}
