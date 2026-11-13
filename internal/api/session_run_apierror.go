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

// apiErrorBackoff is how long a session's API error stands, from when it
// is seen, before each re-prompt; its length is the most consecutive errors
// re-prompted. Tests shorten it.
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
	sentAt   int       // transcript length at the last send: a new error lies past it
	deadline time.Time // when the standing error is acted on; zero if not waiting
	recorded bool      // the record already says the session stopped
}

// poll acts on the transcript past this start's base: t. A transcript that
// does not end in an API error resets the count and cancels a pending send.
// Each wait starts when an error is seen: after a send, only a new error
// line, past what the transcript held then, starts the next one, and a new
// error once every re-prompt is spent goes to the inbox at once.
func (w *apiErrorWatch) poll(t []byte) {
	text, ok := w.a.APIError(t)
	if !ok {
		w.sent, w.sentAt, w.deadline, w.recorded = 0, 0, time.Time{}, false
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
		if w.sent > 0 {
			if len(t) < w.sentAt {
				return
			}
			if _, fresh := w.a.APIError(t[w.sentAt:]); !fresh {
				return
			}
		}
		if w.sent >= len(apiErrorBackoff) {
			w.stop(text)
			return
		}
		w.deadline = time.Now().Add(apiErrorBackoff[w.sent])
		return
	}
	if time.Now().Before(w.deadline) {
		return
	}
	w.send(text, len(t))
}

// stop records the session for the inbox: idle, with the error's text. A
// record that fails is reported and tried again on the next poll.
func (w *apiErrorWatch) stop(text string) {
	w.recorded = w.recordIdle(session.APIErrorDetail+text, "record its API error")
}

// send puts the session at its prompt in the record and sends the message;
// size is the transcript's length now, before the message lands in it. A
// record it cannot read is tried again on the next poll; a message that
// fails counts, and the standing error starts the next wait.
func (w *apiErrorWatch) send(text string, size int) {
	rec, err := w.st.Read(w.sid)
	if err != nil {
		return
	}
	w.sent, w.sentAt, w.deadline = w.sent+1, size, time.Time{}
	if rep, _ := session.Resolve(rec); rep.State != session.StateIdle {
		w.recordIdle("", "record it idle before a re-prompt")
	}
	msg := fmt.Sprintf("Your last turn ended on an API error: %s. Carry on where you left off.", text)
	if err := sendAPIErrorMessage(w.ctx, w.sid, msg); err != nil {
		w.sentAt = 0
		if w.ctx.Stderr != nil {
			fmt.Fprintf(w.ctx.Stderr, "rhiza: session %s: re-prompt after an API error: %v\n", w.sid, err)
		}
	}
}

// recordIdle records an idle hook report with detail, unless the session
// has exited. A failed write goes to Stderr, saying what it was for; it
// reports whether the write succeeded.
func (w *apiErrorWatch) recordIdle(detail, what string) bool {
	rep := session.Report{State: session.StateIdle, Detail: detail, At: w.ctx.Now().UTC()}
	if _, err := w.st.Update(w.sid, func(r *session.Record) {
		if !r.Status.Exited() {
			r.Status.Hook = &rep
		}
	}); err != nil {
		if w.ctx.Stderr != nil {
			fmt.Fprintf(w.ctx.Stderr, "rhiza: session %s: %s: %v\n", w.sid, what, err)
		}
		return false
	}
	return true
}
