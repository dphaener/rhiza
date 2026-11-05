package api

import (
	"fmt"
	"slices"
	"time"

	"github.com/dphaener/rhiza/internal/daemon"
	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// daemonNudgeText is what the daemon tells a stalled session.
const daemonNudgeText = "You have been idle for over 20 minutes without finishing. If you started a " +
	"background shell (a test run, a build), check its result and carry on. If your work is " +
	"finished, run `rhiza session done`: the run moves on only then."

// daemonNudgeTimeout bounds how long a nudge waits for the session's prompt,
// so one session cannot hold up a poll.
const daemonNudgeTimeout = "30s"

// daemonSend sends a session a message; tests replace it.
var daemonSend = func(ctx *Ctx, sid, text string) error {
	_, err := runSessionMessage(ctx, SessionMessageRequest{Session: sid, Text: text, Timeout: daemonNudgeTimeout})
	return err
}

// daemonNudge nudges each stalled session of the live runs, once per stall.
func daemonNudge(ctx *Ctx, rs *run.Store, state *daemon.State) error {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil
	}
	runs, err := rs.List()
	if err != nil {
		return Errorf("read runs: %v", err)
	}
	var live []run.State
	for _, j := range runs {
		if st := j.State(); !st.Ended && rs.Live(j.ID) {
			live = append(live, st)
		}
	}
	kept, err := daemonNudgeStalled(live, session.NewStore(session.Dir(common)), ctx.Now(), state.Nudged,
		func(sid, text string) error { return daemonSend(ctx, sid, text) })
	state.Nudged = kept
	return err
}

// daemonNudgeStalled sends, for each stalled session of an autonomous run
// among states (the run's own Autonomous, not the daemon's flag), one
// message per stall, a stall being keyed on the report time the session
// stalled at; the message is daemonNudgeText. It returns the keys of the stalls that still last. A send that
// fails is not recorded, so the next poll tries again; the first error is
// returned. Owner-attended runs are never nudged: the owner may be the one
// they wait on.
func daemonNudgeStalled(states []run.State, sessions *session.Store, now time.Time, nudged []string, send func(sid, text string) error) ([]string, error) {
	var kept []string
	var first error
	for _, st := range states {
		if !st.Autonomous {
			continue
		}
		for _, sl := range runStalledAll(st, true, sessions, now) {
			key := sl.Session + "@" + sl.Since
			if !slices.Contains(nudged, key) {
				if err := send(sl.Session, daemonNudgeText); err != nil {
					if first == nil {
						first = fmt.Errorf("nudge session %s: %w", sl.Session, err)
					}
					continue
				}
			}
			kept = append(kept, key)
		}
	}
	return kept, first
}
