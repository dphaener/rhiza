package api

// This file holds session.done, the command an agent runs as its very last
// action to say its work is finished. It writes only the record's done
// mark, which an agent step with until_when = "done" waits for (see
// run.Engine): idle at the prompt is not the same as finished.

import (
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "session.done",
		Summary:  "Mark a session's work as done; its agent runs this as its very last action",
		Request:  reflect.TypeFor[SessionDoneRequest](),
		Response: reflect.TypeFor[SessionDoneResponse](),
		Run:      runSessionDone,
	})
}

// envSession names the session a launched agent runs in.
const envSession = "RHIZA_SESSION"

// SessionDoneRequest names the session that is done.
type SessionDoneRequest struct {
	ID string `json:"id,omitempty" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters; default RHIZA_SESSION, which a launched agent has"`
}

// SessionDoneResponse is the session's done mark.
type SessionDoneResponse struct {
	Session string `json:"session" desc:"The session's full ID"`
	At      string `json:"at" desc:"When the session first said done, RFC 3339 UTC"`
	Already bool   `json:"already,omitempty" desc:"True if the session had already said done; the first mark stands"`
	DryRun  bool   `json:"dry_run,omitempty" desc:"True if nothing was written"`
}

// runSessionDone records that the session's agent says its work is done.
// The session is the argument or, when there is none, RHIZA_SESSION. A
// session that has exited cannot say anything more: that is a conflict.
// Saying done twice keeps the first mark.
func runSessionDone(ctx *Ctx, req any) (any, error) {
	r := req.(SessionDoneRequest)
	ref := strings.TrimSpace(r.ID)
	if ref == "" && ctx.Getenv != nil {
		ref = strings.TrimSpace(ctx.Getenv(envSession))
	}
	if ref == "" {
		return nil, Usagef("session done needs a session: name its ID, or run it where %s is set, as in a session Rhiza launched", envSession).
			WithNext(sessionLsNext)
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	if rec.Status.Exited() {
		return nil, exitedConflict(rec.ID)
	}
	if rec.Status.SaidDone() {
		return SessionDoneResponse{Session: rec.ID, At: rfc3339(rec.Status.Done.At), Already: true, DryRun: ctx.DryRun}, nil
	}
	mark := session.Report{State: session.StateDone, At: ctx.Now().UTC()}
	if ctx.DryRun {
		return SessionDoneResponse{Session: rec.ID, At: rfc3339(mark.At), DryRun: true}, nil
	}
	exited, already := false, false
	updated, err := st.Update(rec.ID, func(r *session.Record) {
		// Read again under the record's lock: the agent may have exited,
		// or said done from another process, since.
		switch {
		case r.Status.Exited():
			exited = true
		case r.Status.SaidDone():
			already = true
		default:
			r.Status.Done = &mark
		}
	})
	if err != nil {
		return nil, Errorf("record session %s done: %v", rec.ID, err)
	}
	if exited {
		return nil, exitedConflict(rec.ID)
	}
	return SessionDoneResponse{Session: rec.ID, At: rfc3339(updated.Status.Done.At), Already: already}, nil
}

// exitedConflict refuses a done mark on a session that has exited.
func exitedConflict(sid string) *Error {
	return Conflictf("session %s has exited: its exit is final, so it can no longer say done", sid).
		WithNext("rhiza session show " + sid)
}

// Human says the session is marked done.
func (r SessionDoneResponse) Human(w io.Writer) {
	switch {
	case r.DryRun && r.Already:
		fmt.Fprintf(w, "Session %s already said done at %s.\n\nThis was a dry run: nothing was written.\n", r.Session, r.At)
	case r.DryRun:
		fmt.Fprintf(w, "Would mark session %s done.\n\nThis was a dry run: nothing was written.\n", r.Session)
	case r.Already:
		fmt.Fprintf(w, "Session %s already said done at %s.\n", r.Session, r.At)
	default:
		fmt.Fprintf(w, "Session %s is done.\n", r.Session)
	}
}
