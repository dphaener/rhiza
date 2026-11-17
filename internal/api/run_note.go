package api

// This file holds run.note: a session journals a one-line learning on its
// run (a test that needs a scoped command, a fixture's trap), which run
// status lists oldest first. A note is work-layer state in the run's
// journal: it never touches the model or its log, and it is not rhiza note
// add, which notes a model node.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "run.note",
		Summary:  "Journal a one-line learning (at most 200 characters) on the run RHIZA_RUN names, or --run, with the noting session, role and leaf from RHIZA_SESSION, RHIZA_ROLE and RHIZA_LEAF when set; rhiza run status lists the run's notes oldest first",
		Write:    true,
		Request:  reflect.TypeFor[RunNoteRequest](),
		Response: reflect.TypeFor[RunNoteResponse](),
		Run:      runRunNote,
	})
}

const runNoteNext = "rhiza describe run.note"

// RunNoteRequest is the note and, outside a run, the run it goes on.
type RunNoteRequest struct {
	Text string `json:"text" arg:"0" desc:"The note: one line of at most 200 characters"`
	Run  string `json:"run,omitempty" flag:"run" desc:"The run: a full ID or a unique prefix of 4+ characters; defaults to RHIZA_RUN's"`
}

// RunNoteResponse is the journaled note, or in a dry run the one that
// would be.
type RunNoteResponse struct {
	Run    string  `json:"run" desc:"The run's full ID"`
	Note   RunNote `json:"note" desc:"The note as journaled"`
	DryRun bool    `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// RunNote is one note on a run.
type RunNote struct {
	At      string `json:"at" desc:"When it was journaled, RFC 3339 UTC"`
	Step    string `json:"step,omitempty" desc:"The step the run was at; absent when it was at none"`
	Session string `json:"session,omitempty" desc:"The noting session (RHIZA_SESSION); absent when not set"`
	Role    string `json:"role,omitempty" desc:"The noting session's role (RHIZA_ROLE), such as builder or verifier; absent when not set"`
	Leaf    string `json:"leaf,omitempty" desc:"The leaf it worked on (RHIZA_LEAF); absent when not set"`
	Text    string `json:"text" desc:"The note"`
}

// Human prints the journaled note.
func (r RunNoteResponse) Human(w io.Writer) {
	verb := "Noted"
	if r.DryRun {
		verb = "Would note"
	}
	fmt.Fprintf(w, "%s on run %s%s: %s\n", verb, r.Run, noteBy(r.Note), r.Note.Text)
}

// noteBy is " (builder on leaf-a, session s1)" for a note's role, leaf
// and session, leaving out what is not set, or "" with none of them.
func noteBy(n RunNote) string {
	var parts []string
	who := n.Role
	if n.Leaf != "" {
		who = strings.TrimSpace(who + " on " + n.Leaf)
	}
	if who != "" {
		parts = append(parts, who)
	}
	if n.Session != "" {
		parts = append(parts, "session "+n.Session)
	}
	if len(parts) == 0 {
		return ""
	}
	return " (" + strings.Join(parts, ", ") + ")"
}

// runNotes is the response form of a journal's notes, oldest first.
func runNotes(in []run.Note) []RunNote {
	var out []RunNote
	for _, n := range in {
		out = append(out, RunNote{At: rfc3339(n.At), Step: n.Step, Session: n.Session, Role: n.Role, Leaf: n.Leaf, Text: n.Text})
	}
	return out
}

func runRunNote(ctx *Ctx, req any) (any, error) {
	const cmd = "run note"
	r := req.(RunNoteRequest)
	if p := run.NoteProblem(r.Text); p != "" {
		return nil, Usagef("%s: %s", cmd, p).WithNext(runNoteNext)
	}
	getenv := ctxGetenv(ctx)
	ref := strings.TrimSpace(r.Run)
	if ref == "" {
		ref = strings.TrimSpace(getenv(envRun))
	}
	if ref == "" {
		return nil, Usagef("%s works inside a run, and %s is not set: run it from a session of the run, or name the run with --run", cmd, envRun).WithNext(runNoteNext)
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := runStore(ctx, common)
	j, err := readRun(rs, ref)
	if err != nil {
		return nil, err
	}
	e := run.Entry{
		Type: run.TypeNote, At: ctx.Now(), Step: j.State().Current, Text: r.Text,
		Session: strings.TrimSpace(getenv(envSession)),
		Role:    strings.TrimSpace(getenv(run.EnvRole)),
		Leaf:    strings.TrimSpace(getenv(run.EnvLeaf)),
	}
	resp := RunNoteResponse{Run: j.ID, DryRun: ctx.DryRun, Note: runNotes([]run.Note{{
		Step: e.Step, Session: e.Session, Role: e.Role, Leaf: e.Leaf, Text: e.Text, At: e.At,
	}})[0]}
	if ctx.DryRun {
		return resp, nil
	}
	if err := rs.Append(j.ID, e); err != nil {
		return nil, Errorf("journal the note in run %s: %v", j.ID, err).WithNext("rhiza run status " + j.ID)
	}
	return resp, nil
}
