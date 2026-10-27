package api

// This file holds session.run, the command a session's terminal tab runs,
// and session.hook, the command the agent's hooks run. Together they keep a
// session's status in its record: run starts the agent, copies the tab's
// Rex status reports while it runs and records how it exited; hook records
// what each hook says and, for a session in a run, journals in the run the
// messages the session sends. Each writes only its own source of the status.
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
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/adapter"
	"github.com/dphaener/rhiza/internal/run"
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
	// envAutonomous is set to 1 for an agent an autonomous run launched.
	envAutonomous = "RHIZA_AUTONOMOUS"
	// envLeaf and envAttempt are a leaf session's leaf handle and attempt.
	envLeaf    = run.EnvLeaf
	envAttempt = run.EnvAttempt
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
	Session   SessionInfo `json:"session" desc:"The session as recorded after the agent ended"`
	Argv      []string    `json:"argv" desc:"The agent's command"`
	Env       []string    `json:"env" desc:"The variables set over this terminal's environment for the agent"`
	Code      int         `json:"code" desc:"The agent's exit code, which is also this command's; 128 plus the signal number if a signal ended it"`
	Decisions []string    `json:"decisions,omitempty" desc:"Handles of the Decisions proposed from the questions the owner answered (see session.interview)"`
	Warnings  []string    `json:"warnings,omitempty" desc:"Status reports that could not be copied while the agent ran, token use that could not be recorded, and a failed capture of the owner's answers"`
	DryRun    bool        `json:"dry_run,omitempty" desc:"True if nothing was started or written"`
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
	Session   string `json:"session,omitempty" desc:"The session the payload names, if Rhiza launched it"`
	Recorded  bool   `json:"recorded" desc:"Whether the hook report was recorded"`
	State     string `json:"state,omitempty" enum:"starting,working,idle,blocked,done,error,exited" desc:"The state recorded"`
	Kind      string `json:"kind,omitempty" desc:"What a blocked session waits for"`
	Journaled bool   `json:"journaled,omitempty" desc:"Whether a message the session sent was journaled in its run"`
	// reply is what the hook prints on stdout for the agent to read: the
	// denial of a question in an autonomous session. Human writes it.
	reply []byte
}

// autonomousAskReason is what the agent of an autonomous session is told
// when its hook denies a question to the owner.
const autonomousAskReason = "This run is autonomous: no owner will answer. Your question is recorded and the arbiter will decide it; its answer arrives as a message. End your turn now and wait at your prompt."

// runSessionRun starts the session's agent in its worktree with this
// terminal, its environment stamped with the change set, actor, for-actor
// and run; copies this tab's Rex status reports while the agent runs; and
// records the exit; then proposes a Decision for each question the owner
// answered (see captureInterview). It returns the agent's exit code as its
// own: a capture that fails is only a warning.
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
		Prompt:         rec.Prompt,
		Model:          rec.Model,
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
	} else {
		rec = updated
	}
	if warn := recordTokenUse(ctx, st, rec); warn != "" {
		resp.Warnings = append(resp.Warnings, warn)
	} else if updated, err := st.Read(rec.ID); err == nil {
		rec = updated
	}
	resp.Session = sessionInfo(rec)
	resp.Decisions, resp.Warnings = captureAfterRun(ctx, rec, resp.Warnings)
	return resp, nil
}

// recordTokenUse reads the session's token use from its transcript,
// overwrites the record's tokens with it and, for a session in a run,
// journals a tokens entry: the later reading wins. It runs when the
// session says done, when its agent step ends and when its agent exits.
// It returns a warning, never an error: a missing or unreadable transcript
// records nothing. A caller that gives neither a home folder nor an
// environment has no transcript to find, so nothing is read.
func recordTokenUse(ctx *Ctx, st *session.Store, rec session.Record) string {
	if ctx.Home == "" && ctx.Getenv == nil {
		return ""
	}
	warn := func(msg string) string { return "record token use: " + msg }
	a, ok := adapter.Lookup(rec.Agent)
	if !ok {
		return warn(fmt.Sprintf("session %s names agent %q, which this rhiza does not know", rec.ID, rec.Agent))
	}
	transcript, err := readTranscript(ctx, a, rec.ID)
	if err != nil {
		return warn(AsError(err).Message)
	}
	use, err := a.TokenUse(transcript)
	if err != nil {
		return warn(fmt.Sprintf("read the token use of session %s: %v", rec.ID, err))
	}
	tokens := sessionTokenUse(use)
	if _, err := st.Update(rec.ID, func(r *session.Record) { r.Tokens = &tokens }); err != nil {
		return warn(fmt.Sprintf("record session %s: %v", rec.ID, err))
	}
	if rec.Run == "" {
		return ""
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return warn(AsError(err).Message)
	}
	ru := runTokenUse(tokens)
	e := run.Entry{Type: run.TypeTokens, At: ctx.Now(), Step: rec.Step, Session: rec.ID, Tokens: &ru}
	if err := run.NewStore(run.Dir(common)).Append(rec.Run, e); err != nil {
		return warn(fmt.Sprintf("journal session %s's token use in run %s: %v", rec.ID, rec.Run, err))
	}
	return ""
}

// sessionTokenUse is an adapter's reading as the record keeps it.
func sessionTokenUse(u adapter.TokenUse) session.TokenUse {
	c := func(t adapter.TokenCounts) session.TokenCounts {
		return session.TokenCounts{Input: t.Input, Output: t.Output, CacheCreation: t.CacheCreation, CacheRead: t.CacheRead}
	}
	return session.TokenUse{Session: c(u.Session), Helpers: c(u.Helpers)}
}

// runTokenUse is a record's reading as the journal keeps it.
func runTokenUse(u session.TokenUse) run.TokenUse {
	c := func(t session.TokenCounts) run.TokenCounts {
		return run.TokenCounts{Input: t.Input, Output: t.Output, CacheCreation: t.CacheCreation, CacheRead: t.CacheRead}
	}
	return run.TokenUse{Session: c(u.Session), Helpers: c(u.Helpers)}
}

// captureAfterRun proposes the session's answered questions as Decisions
// once its agent has exited, and returns their handles and warnings with a
// failure appended. A caller that gives neither a home folder nor an
// environment has no transcript to find, so nothing is captured.
func captureAfterRun(ctx *Ctx, rec session.Record, warnings []string) ([]string, []string) {
	if ctx.Home == "" && ctx.Getenv == nil {
		return nil, warnings
	}
	got, err := captureInterview(ctx, rec)
	if err != nil {
		e := AsError(err)
		msg := "propose decisions from the owner's answers: " + e.Message
		if e.Next != "" {
			msg += " (next: " + e.Next + ")"
		}
		return nil, append(warnings, msg)
	}
	var handles []string
	for _, p := range got.Proposed {
		handles = append(handles, p.Handle)
	}
	return handles, warnings
}

// launchEnv is what the agent's environment gets over the terminal's: the
// change set, the agent's actor, the actor it acts for, the session (so
// `rhiza session done` needs no argument) and, only when the launch belongs
// to one, the run, RHIZA_AUTONOMOUS=1 when that run is autonomous, and a
// leaf session's RHIZA_LEAF and RHIZA_ATTEMPT. A run, autonomous or leaf
// mark the terminal had is removed.
func launchEnv(rec session.Record) (set, unset []string) {
	change := rec.Change
	if rec.Writes != "" {
		change = rec.Writes
	}
	set = []string{
		envChange + "=" + change,
		envActor + "=" + rec.Actor,
		envFor + "=" + rec.For,
		envSession + "=" + rec.ID,
	}
	if rec.Run != "" {
		set = append(set, envRun+"="+rec.Run)
	} else {
		unset = append(unset, envRun)
	}
	if rec.Autonomous {
		set = append(set, envAutonomous+"=1")
	} else {
		unset = append(unset, envAutonomous)
	}
	if rec.Leaf != "" {
		set = append(set, envLeaf+"="+rec.Leaf, envAttempt+"="+strconv.Itoa(rec.Attempt))
	} else {
		unset = append(unset, envLeaf, envAttempt)
	}
	return set, unset
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

// runSessionHook records the status a hook reports and journals a run's
// messages (see sessionHook). It never fails: a hook
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
// sessions, records the report its adapter reads from the payload and, when
// the session belongs to a run, journals in the run any message the payload
// says it sent. The message is journaled whether or not the status was
// recorded, even after the session exited.
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
		serr := recordHook(ctx, st, rec, a, event, payload, rep, changed, &resp)
		journaled, jerr := journalHookMessage(ctx, a, rec, event, payload)
		resp.Journaled = journaled
		return resp, errors.Join(serr, jerr)
	}
	return resp, errors.Join(errs...)
}

// recordHook records the report a hook gave in the session's record and
// says so in resp. A report that changes nothing, or arrives after the
// session exited, is not recorded. In an autonomous session a question to
// the owner is denied (resp's reply) and kept as the asked mark: the
// arbiter answers it by message.
func recordHook(ctx *Ctx, st *session.Store, rec session.Record, a adapter.Adapter, event string, payload []byte, rep adapter.HookReport, changed bool, resp *SessionHookResponse) error {
	name := a.Name()
	if !changed {
		return nil
	}
	if !session.ValidState(rep.State) {
		return fmt.Errorf("%s reported unknown state %q", name, rep.State)
	}
	report := session.Report{State: rep.State, Kind: rep.Kind, Detail: rep.Detail, At: ctx.Now().UTC()}
	var reply []byte
	if rec.Autonomous {
		b, ask, err := a.AskReply(event, payload, autonomousAskReason)
		if err != nil {
			return fmt.Errorf("%s: %v", name, err)
		}
		if ask {
			reply = b
		}
	}
	recorded := false
	if _, err := st.Update(rec.ID, func(r *session.Record) {
		// A hook that arrives after the exit does not revive the session.
		if !r.Status.Exited() {
			r.Status.Hook, recorded = &report, true
			if reply != nil {
				r.Status.Asked = &session.Report{State: session.StateBlocked, Kind: session.KindQuestion, Detail: rep.Detail, At: report.At}
			}
		}
	}); err != nil {
		return fmt.Errorf("record session %s: %v", rec.ID, err)
	}
	if !recorded {
		return nil
	}
	resp.Recorded, resp.State, resp.Kind, resp.reply = true, report.State, report.Kind, reply
	return nil
}

// journalHookMessage appends the message the payload says the session sent
// to the session's run journal, and reports whether it did. A session in
// no run journals nothing.
func journalHookMessage(ctx *Ctx, a adapter.Adapter, rec session.Record, event string, payload []byte) (bool, error) {
	m, ok, err := a.HookMessage(event, payload)
	if err != nil {
		return false, fmt.Errorf("%s: %v", a.Name(), err)
	}
	if !ok || rec.Run == "" {
		return false, nil
	}
	common, err := gitCommonDir(ctx)
	if err != nil {
		return false, err
	}
	e := run.Entry{
		Type: run.TypeMessage, At: ctx.Now(), Step: rec.Step, Session: rec.ID, Actor: rec.Actor,
		To: m.To, Text: m.Text, Summary: m.Summary, Structured: m.Structured,
	}
	if err := run.NewStore(run.Dir(common)).Append(rec.Run, e); err != nil {
		return false, fmt.Errorf("journal the message session %s sent in run %s: %v", rec.ID, rec.Run, err)
	}
	return true, nil
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
	for _, h := range r.Decisions {
		fmt.Fprintf(w, "- Proposed decision %s\n", h)
	}
	for _, warn := range r.Warnings {
		fmt.Fprintf(w, "- Warning: %s\n", warn)
	}
	fmt.Fprintf(w, "\nnext: rhiza session show %s\n", s.ID)
}

// Human writes the reply for the agent, if there is one, and nothing else:
// the agent reads what a hook prints.
func (r SessionHookResponse) Human(w io.Writer) {
	if len(r.reply) > 0 {
		w.Write(r.reply)
		io.WriteString(w, "\n")
	}
}
