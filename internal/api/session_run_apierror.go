package api

// This file holds how `rhiza session run` answers an agent whose last turn
// ended on an API error: after a backoff it sends the session a message to
// carry on, a few times, then leaves it to the inbox.

import (
	"fmt"
	"time"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/session"
)

// apiErrorBackoff is how long a session's API error stands before each
// re-prompt, from when the error is seen or the last message sent; its
// length is the most consecutive errors re-prompted. Tests shorten it.
var apiErrorBackoff = []time.Duration{time.Minute, 2 * time.Minute, 4 * time.Minute}

// sendAPIErrorMessage sends the session a message, as session.message does.
// Tests replace it.
var sendAPIErrorMessage = func(ctx *Ctx, sid, text string) error {
	_, err := runSessionMessage(ctx, SessionMessageRequest{Session: sid, Text: text})
	return err
}

// apiErrorWatch is the state of one session's API-error handling across
// the polls of its transcript.
type apiErrorWatch struct {
	ctx        *Ctx
	a          adapter.Adapter
	st         *session.Store
	sid        string
	autonomous bool

	sent     int       // messages sent for the current run of errors
	deadline time.Time // when the standing error is acted on; zero if not waiting
	recorded bool      // the record already says the session stopped
}

// poll acts on the transcript past this start's base: t. A transcript that
// does not end in an API error resets the count and cancels a pending send.
func (w *apiErrorWatch) poll(t []byte) {
	text, ok := w.a.APIError(t)
	if !ok {
		w.sent, w.deadline, w.recorded = 0, time.Time{}, false
		return
	}
	if w.recorded {
		return
	}
	if !w.autonomous {
		w.stop(text)
		return
	}
	if w.deadline.IsZero() {
		w.deadline = time.Now().Add(apiErrorBackoff[min(w.sent, len(apiErrorBackoff)-1)])
		return
	}
	if time.Now().Before(w.deadline) {
		return
	}
	if w.sent >= len(apiErrorBackoff) {
		w.stop(text)
		return
	}
	w.send(text)
}

// stop records the session for the inbox: idle, with the error's text.
func (w *apiErrorWatch) stop(text string) {
	w.recorded = true
	w.recordIdle(session.APIErrorDetail + text)
}

// send puts the session at its prompt in the record and sends the message.
func (w *apiErrorWatch) send(text string) {
	w.sent++
	w.deadline = time.Now().Add(apiErrorBackoff[min(w.sent, len(apiErrorBackoff)-1)])
	rec, err := w.st.Read(w.sid)
	if err != nil {
		return
	}
	if rep, _ := session.Resolve(rec); rep.State != session.StateIdle {
		w.recordIdle("")
	}
	msg := fmt.Sprintf("Your last turn ended on an API error: %s. Carry on where you left off.", text)
	if err := sendAPIErrorMessage(w.ctx, w.sid, msg); err != nil && w.ctx.Stderr != nil {
		fmt.Fprintf(w.ctx.Stderr, "rhiza: session %s: re-prompt after an API error: %v\n", w.sid, err)
	}
}

// recordIdle records an idle hook report with detail, unless the session
// has exited.
func (w *apiErrorWatch) recordIdle(detail string) {
	rep := session.Report{State: session.StateIdle, Detail: detail, At: w.ctx.Now().UTC()}
	_, _ = w.st.Update(w.sid, func(r *session.Record) {
		if !r.Status.Exited() {
			r.Status.Hook = &rep
		}
	})
}
