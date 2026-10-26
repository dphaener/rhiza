package api

// This file holds session.run, the command a session's terminal tab runs,
// and session.hook, the command the agent's hooks run. Together they keep a
// session's status in its record: run starts the agent, copies the tab's
// Rex status reports while it runs and records how it exited; hook records
// what each hook says. Each writes only its own source of the status.
//
// Nothing here names an agent: the record names its adapter.

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/session"
)

func init() {
	register(Operation{
		Name:     "session.run",
		Summary:  "Run a session's agent in this terminal and record its status and exit; its tab runs this",
		Request:  reflect.TypeFor[SessionRunRequest](),
		Response: reflect.TypeFor[SessionRunResponse](),
		Run:      runSessionRun,
	})
	register(Operation{
		Name:     "session.hook",
		Summary:  "Record the status an agent hook reports, read from stdin; the agent's hooks run this",
		Request:  reflect.TypeFor[SessionHookRequest](),
		Response: reflect.TypeFor[SessionHookResponse](),
		Run:      runSessionHook,
	})
}

// The environment variables a launched agent gets.
const (
	envChange = "RHIZA_CHANGE"
	envActor  = "RHIZA_ACTOR"
	envFor    = "RHIZA_FOR"
	envRun    = "RHIZA_RUN"
)

// The Rex event that carries a block's status report, and its envelopes.
const (
	rexSessionEvent       = "session_event"
	rexBlockEvent         = "block_event"
	rexProgramStatusEvent = "program_status_changed"
	rexStatusRemovedEvent = "program_status_removed"
	rexChildExited        = "child_exited"
)

// SessionRunRequest names the session to run.
type SessionRunRequest struct {
	ID string `json:"id" arg:"0" desc:"The session: its full ID or a unique prefix of 4+ characters"`
}

// SessionRunResponse is how the agent ended, or in a dry run what would
// start.
type SessionRunResponse struct {
	Session  SessionInfo `json:"session" desc:"The session as recorded after the agent ended"`
	Argv     []string    `json:"argv" desc:"The agent's command"`
	Env      []string    `json:"env" desc:"The variables set over this terminal's environment for the agent"`
	Code     int         `json:"code" desc:"The agent's exit code, which is also this command's; 128 plus the signal number if a signal ended it"`
	Warnings []string    `json:"warnings,omitempty" desc:"Status reports that could not be copied while the agent ran"`
	DryRun   bool        `json:"dry_run,omitempty" desc:"True if nothing was started or written"`
}

// ExitCode is the agent's exit code.
func (r SessionRunResponse) ExitCode() int { return r.Code }

// SessionHookRequest names the hook event; its payload is stdin.
type SessionHookRequest struct {
	Event string `json:"event,omitempty" arg:"0" desc:"The hook event, as the agent names it; the payload is read from stdin"`
}

// SessionHookResponse says what the hook recorded. Its markdown is empty:
// an agent may read a hook's output.
type SessionHookResponse struct {
	Session  string `json:"session,omitempty" desc:"The session the payload names, if Rhiza launched it"`
	Recorded bool   `json:"recorded" desc:"Whether the hook report was recorded"`
	State    string `json:"state,omitempty" enum:"starting,working,idle,blocked,done,error,exited" desc:"The state recorded"`
	Kind     string `json:"kind,omitempty" desc:"What a blocked session waits for"`
}

// runSessionRun starts the session's agent in its worktree with this
// terminal, its environment stamped with the change set, actor, for-actor
// and run; copies this tab's Rex status reports while the agent runs; and
// records the exit. It returns the agent's exit code as its own.
func runSessionRun(ctx *Ctx, req any) (any, error) {
	r := req.(SessionRunRequest)
	ref := strings.TrimSpace(r.ID)
	if ref == "" {
		return nil, Usagef("session run needs a session ID").WithNext(sessionLsNext)
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return nil, err
	}
	rec, err := readSession(st, ref)
	if err != nil {
		return nil, err
	}
	a, ok := adapter.Lookup(rec.Agent)
	if !ok {
		return nil, Errorf("session %s names agent %q, which this rhiza does not know: want one of %s",
			rec.ID, rec.Agent, strings.Join(adapter.Names(), ", "))
	}
	// Session start checked the mode; check again, since the record is a
	// file anyone can edit.
	if rec.PermissionMode == "" {
		return nil, Errorf("session %s has no permission mode: an agent never starts without an explicit one", rec.ID).
			WithNext(sessionStartNext)
	}
	if a.Bypass(rec.PermissionMode) {
		return nil, bypassRefused(rec.PermissionMode)
	}
	self := ctx.Self
	if self == "" {
		if !ctx.DryRun {
			return nil, Errorf("cannot run session %s: the path of the running rhiza, which its hooks call, is unknown", rec.ID)
		}
		self = "rhiza"
	}
	argv, err := a.Command(adapter.Launch{
		SessionID:      rec.ID,
		PermissionMode: rec.PermissionMode,
		HookCommand:    []string{self, "session", "hook"},
	})
	if err != nil {
		return nil, Errorf("session %s: %v", rec.ID, err).WithNext(sessionStartNext)
	}
	env, unset := launchEnv(rec)
	resp := SessionRunResponse{Session: sessionInfo(rec), Argv: argv, Env: env}
	if ctx.DryRun {
		resp.DryRun = true
		return resp, nil
	}
	if ctx.Start == nil {
		return nil, Errorf("cannot run session %s: this caller starts no programs", rec.ID)
	}

	follow := followRex(ctx, st, rec.ID)
	p, err := ctx.Start(Program{Name: argv[0], Args: argv[1:], Dir: rec.Worktree, Env: env, Unset: unset, Terminal: true})
	if err != nil {
		follow.stop()
		return nil, recordRunFailure(ctx, st, rec.ID, fmt.Sprintf("start %s: %v", a.Name(), err))
	}
	code, werr := p.Wait()
	resp.Warnings = follow.stop()
	if werr != nil {
		return nil, recordRunFailure(ctx, st, rec.ID, fmt.Sprintf("wait for %s: %v", a.Name(), werr))
	}
	resp.Code = code
	updated, err := st.Update(rec.ID, func(r *session.Record) {
		r.Status.Exit = &session.Report{State: session.StateExited, Code: &code, At: ctx.Now().UTC()}
	})
	if err != nil {
		resp.Warnings = append(resp.Warnings, fmt.Sprintf("record the exit: %v", err))
		return resp, nil
	}
	resp.Session = sessionInfo(updated)
	return resp, nil
}

// launchEnv is what the agent's environment gets over the terminal's: the
// change set, the agent's actor, the actor it acts for and, only when the
// launch belongs to one, the run. A run the terminal had is removed.
func launchEnv(rec session.Record) (set, unset []string) {
	set = []string{
		envChange + "=" + rec.Change,
		envActor + "=" + rec.Actor,
		envFor + "=" + rec.For,
	}
	if rec.Run != "" {
		return append(set, envRun+"="+rec.Run), nil
	}
	return set, []string{envRun}
}

// recordRunFailure records that the agent could not run as the exit
// source's error, and returns the error to report.
func recordRunFailure(ctx *Ctx, st *session.Store, sid, msg string) error {
	if _, err := st.Update(sid, func(r *session.Record) {
		r.Status.Exit = &session.Report{State: session.StateError, Detail: msg, At: ctx.Now().UTC()}
	}); err != nil {
		msg += "; and recording that failed: " + err.Error()
	}
	return Errorf("session %s: %s", sid, msg).WithNext("rhiza session show " + sid)
}

// rexFollower copies one block's Rex status reports into a session record
// until stopped.
type rexFollower struct {
	p        Process
	done     chan struct{}
	warnings []string
}

// followRex starts `rex events` for the caller's Rex session and copies
// every status report of the caller's block into the record's rex source.
// Outside Rex, or when Rex cannot be started, it follows nothing.
func followRex(ctx *Ctx, st *session.Store, sid string) *rexFollower {
	f := &rexFollower{}
	if ctx.RexBlock == "" || ctx.RexSession == "" || ctx.Start == nil {
		return f
	}
	p, err := ctx.Start(Program{Name: "rex", Args: []string{"events", "-s", ctx.RexSession, "--json"}})
	if err != nil {
		f.warnings = append(f.warnings, fmt.Sprintf("follow rex events: %v", err))
		return f
	}
	f.p, f.done = p, make(chan struct{})
	go func() {
		defer close(f.done)
		rd := bufio.NewReader(p.Output())
		for {
			line, err := rd.ReadBytes('\n')
			if rep, exited, ok := rexStatus(line, ctx.RexBlock); ok {
				rep.At = ctx.Now().UTC()
				if _, uerr := st.Update(sid, func(r *session.Record) { recordRex(r, rep, exited) }); uerr != nil {
					f.warnings = append(f.warnings, fmt.Sprintf("record rex status %s: %v", rep.State, uerr))
				}
			}
			if err != nil {
				if !errors.Is(err, io.EOF) && !errors.Is(err, fs.ErrClosed) {
					f.warnings = append(f.warnings, fmt.Sprintf("read rex events: %v", err))
				}
				return
			}
		}
	}()
	return f
}

// stop ends `rex events`, waits until every line it wrote is copied and
// returns what could not be.
func (f *rexFollower) stop() []string {
	if f.p == nil {
		return f.warnings
	}
	f.p.Kill()
	<-f.done
	f.p.Wait()
	return f.warnings
}

// recordRex copies a Rex report into the record. A report that the
// terminal's program exited becomes the exit source, unless the wrapper
// already recorded the exit with its code. Nothing is recorded once the
// session has exited: a late report does not revive it.
func recordRex(r *session.Record, rep session.Report, exited bool) {
	if r.Status.Exited() {
		return
	}
	if exited {
		r.Status.Exit = &rep
		return
	}
	r.Status.Rex = &rep
}

// rexEvent is the part of a `rex events --json` line that carries a block's
// status: a block_event inside a session_event, whose payload is a report
// (program_status_changed) or a removal (program_status_removed).
type rexEvent struct {
	Type    string `json:"type"`
	Payload struct {
		Event struct {
			Type    string `json:"type"`
			Payload struct {
				BlockID string `json:"block_id"`
				Name    string `json:"name"`
				Payload struct {
					Record struct {
						State string `json:"state"`
						Kind  string `json:"kind"`
						Msg   string `json:"msg"`
					} `json:"record"`
					Reason string `json:"reason"`
				} `json:"payload"`
			} `json:"payload"`
		} `json:"event"`
	} `json:"payload"`
}

// rexStatus reads one `rex events --json` line and returns the status it
// carries for block, if any. A status report gives state, kind and message;
// Rex asks that an unknown state be taken as idle and an unknown kind as
// none. A removal because the terminal's program exited gives an exited
// report and exited true, so Rex alone can end a session whose wrapper
// never saw the exit.
func rexStatus(line []byte, block string) (rep session.Report, exited, ok bool) {
	var ev rexEvent
	if err := json.Unmarshal(line, &ev); err != nil {
		return session.Report{}, false, false
	}
	be := ev.Payload.Event
	if ev.Type != rexSessionEvent || be.Type != rexBlockEvent || be.Payload.BlockID != block {
		return session.Report{}, false, false
	}
	switch be.Payload.Name {
	case rexProgramStatusEvent:
	case rexStatusRemovedEvent:
		if be.Payload.Payload.Reason != rexChildExited {
			return session.Report{}, false, false
		}
		return session.Report{State: session.StateExited, Detail: "rex: the terminal's program exited"}, true, true
	default:
		return session.Report{}, false, false
	}
	rec := be.Payload.Payload.Record
	rep = session.Report{State: rec.State, Detail: rec.Msg}
	switch rec.State {
	case session.StateIdle, session.StateWorking, session.StateDone, session.StateError:
	case session.StateBlocked:
		switch rec.Kind {
		case session.KindPermission, session.KindQuestion, session.KindAuth:
			rep.Kind = rec.Kind
		}
	default:
		rep.State = session.StateIdle
	}
	return rep, false, true
}

// runSessionHook records the status a hook reports. It never fails: a hook
// that failed would get in the agent's way, so problems go to Ctx.Stderr
// and the call succeeds. A payload for a session Rhiza did not launch is
// ignored without a word.
func runSessionHook(ctx *Ctx, req any) (any, error) {
	r := req.(SessionHookRequest)
	resp, err := sessionHook(ctx, strings.TrimSpace(r.Event))
	if err != nil && ctx.Stderr != nil {
		fmt.Fprintf(ctx.Stderr, "rhiza: session hook %s: %v\n", r.Event, err)
	}
	return resp, nil
}

// sessionHook finds the session the payload names, among every adapter's
// sessions, and records the report its adapter reads from the payload.
func sessionHook(ctx *Ctx, event string) (SessionHookResponse, error) {
	var resp SessionHookResponse
	if event == "" {
		return resp, errors.New("no hook event named")
	}
	var payload []byte
	if ctx.Stdin != nil {
		var err error
		if payload, err = io.ReadAll(ctx.Stdin); err != nil {
			return resp, fmt.Errorf("read the payload: %v", err)
		}
	}
	st, err := sessionStore(ctx)
	if err != nil {
		return resp, err
	}
	var errs []error
	for _, name := range adapter.Names() {
		a, _ := adapter.Lookup(name)
		rep, changed, err := a.HookStatus(event, payload)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %v", name, err))
			continue
		}
		rec, found, err := hookSession(st, rep.SessionID)
		if err != nil {
			return resp, err
		}
		if !found || rec.Agent != name {
			continue
		}
		resp.Session = rec.ID
		if !changed {
			return resp, nil
		}
		if !session.ValidState(rep.State) {
			return resp, fmt.Errorf("%s reported unknown state %q", name, rep.State)
		}
		report := session.Report{State: rep.State, Kind: rep.Kind, Detail: rep.Detail, At: ctx.Now().UTC()}
		recorded := false
		if _, err := st.Update(rec.ID, func(r *session.Record) {
			// A hook that arrives after the exit does not revive the session.
			if !r.Status.Exited() {
				r.Status.Hook, recorded = &report, true
			}
		}); err != nil {
			return resp, fmt.Errorf("record session %s: %v", rec.ID, err)
		}
		if !recorded {
			return resp, nil
		}
		resp.Recorded, resp.State, resp.Kind = true, report.State, report.Kind
		return resp, nil
	}
	return resp, errors.Join(errs...)
}

// hookSession reads the session with exactly the ID sid. A payload must
// name a session in full: no prefix matches.
func hookSession(st *session.Store, sid string) (session.Record, bool, error) {
	if sid == "" {
		return session.Record{}, false, nil
	}
	rec, err := st.Read(sid)
	var nf *session.NotFoundError
	switch {
	case errors.As(err, &nf):
		return rec, false, nil
	case err != nil:
		return rec, false, fmt.Errorf("read session %s: %v", sid, err)
	}
	return rec, rec.ID == sid, nil
}

// Human reports how the agent ended, or what a dry run would start.
func (r SessionRunResponse) Human(w io.Writer) {
	s := r.Session
	if r.DryRun {
		fmt.Fprintf(w, "Would run session %s for change set %s: %s\n\n", s.ID, s.Change, s.Title)
		fmt.Fprintf(w, "- Worktree: %s\n", s.Worktree)
		fmt.Fprintf(w, "- Command: %s\n", strings.Join(r.Argv, " "))
		fmt.Fprintf(w, "- Environment: %s\n", strings.Join(r.Env, " "))
		fmt.Fprintf(w, "\nThis was a dry run: nothing was started or written.\n")
		return
	}
	fmt.Fprintf(w, "Session %s ended: the agent exited with code %d.\n", s.ID, r.Code)
	for _, warn := range r.Warnings {
		fmt.Fprintf(w, "- Warning: %s\n", warn)
	}
	fmt.Fprintf(w, "\nnext: rhiza session show %s\n", s.ID)
}

// Human writes nothing: the agent may read what a hook prints.
func (SessionHookResponse) Human(io.Writer) {}
