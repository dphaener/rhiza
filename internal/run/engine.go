package run

import (
	"errors"
	"fmt"
	"time"

	"github.com/dphaener/rhiza/internal/workflow"
)

// OutcomeError is the outcome an ended entry records when the Runner failed
// to run a step. It always ends the run failed, whatever the workflow's
// `on.` entries say.
const OutcomeError = "error"

var (
	// ErrLive means another process holds the run's lock.
	ErrLive = errors.New("run is running in another process")
	// ErrFinished means the run has already finished.
	ErrFinished = errors.New("run is finished")
)

// runError names the run and wraps ErrLive or ErrFinished.
type runError struct {
	id   string
	kind error
}

func (e *runError) Error() string {
	switch e.kind {
	case ErrLive:
		return fmt.Sprintf("run %s is running in another process", e.id)
	case ErrFinished:
		return fmt.Sprintf("run %s is finished", e.id)
	}
	return fmt.Sprintf("run %s: %v", e.id, e.kind)
}

func (e *runError) Unwrap() error { return e.kind }

// Clock is the time the engine polls by: wait conditions, agent sessions
// and wait timeouts. Journal timestamps never come from it.
type Clock interface {
	Now() time.Time
	Sleep(time.Duration)
}

// RealClock is the wall clock.
type RealClock struct{}

// Now returns time.Now.
func (RealClock) Now() time.Time { return time.Now() }

// Sleep calls time.Sleep.
func (RealClock) Sleep(d time.Duration) { time.Sleep(d) }

// Runner runs steps for the engine. The caller supplies it, so this package
// runs no program and names no agent. env holds the step's RHIZA_*
// variables, to be laid over the runner's own environment.
type Runner interface {
	// Command runs a command step's `run` once and returns its exit code.
	// A non-zero exit is not an error; failing to run it is.
	Command(step workflow.Step, env []string) (code int, detail string, err error)
	// Holds runs a wait step's `until` condition once.
	Holds(step workflow.Step, env []string) (bool, error)
	// StartAgent launches an agent step's session and returns its ID.
	StartAgent(step workflow.Step, env []string) (session string, err error)
	// AgentStatus reports where a session is: ended, and with which
	// outcome (pass or fail), or still running and perhaps blocked on the
	// owner.
	AgentStatus(session string) (AgentStatus, error)
}

// AgentStatus is one poll of an agent step's session.
type AgentStatus struct {
	// Done is true once the session's process has exited; exit is final.
	Done bool
	// Outcome (pass or fail) and Detail are set when Done.
	Outcome string
	Detail  string
	// Blocked is what the session waits on, nil when it is not blocked.
	// It is ignored once Done.
	Blocked *Block
}

// Notice tells the caller that a run's agent step became blocked on the
// owner (Blocked true) or that the block ended (Blocked false). Block is
// what the step waits, or waited, on.
type Notice struct {
	Run     string
	Change  string
	Title   string
	Step    string
	Session string
	Blocked bool
	Block   *Block
}

// Engine runs a run's steps in the foreground, journalling each one.
type Engine struct {
	Store  *Store
	Runner Runner
	// Clock drives polling and wait timeouts; nil means RealClock.
	Clock Clock
	// Now stamps journal entries; nil means time.Now.
	Now func() time.Time
	// NoWait stops the run at a wait step whose condition does not hold
	// yet, instead of polling it. The run stays resumable.
	NoWait bool
	// Log, when set, receives one progress line per event.
	Log func(string)
	// Notify, when set, is called once when an agent step becomes blocked
	// and once when that block ends. A resume of a step whose journal
	// already shows the block calls it only when the block ends.
	Notify func(Notice)
}

// Start creates a run's journal with its started entry. The run begins at
// the definition's start step.
func Start(store *Store, id string, s Started, now time.Time, def *workflow.Definition) error {
	start := def.Start
	if start == "" && len(def.Order) > 0 {
		start = def.Order[0]
	}
	if start == "" {
		return fmt.Errorf("run: workflow %q has no start step", def.Name)
	}
	return store.Create(id, Entry{Type: TypeStarted, At: now, Step: start, Started: &s})
}

// ActionKind is what the engine does next with a run.
type ActionKind int

const (
	// ActFinished means the run has finished; there is nothing to do.
	ActFinished ActionKind = iota
	// ActContinue means the current step was entered and not ended: run it
	// again without a new visit (the process was killed, or it stopped at
	// a wait).
	ActContinue
	// ActEnter means enter Step as visit Visit.
	ActEnter
	// ActOverCap means entering Step would go over its max_visits: the run
	// finishes failed with Detail.
	ActOverCap
	// ActFinish means the last step ended the run but no finished entry
	// was written: finish it with End and Detail.
	ActFinish
)

// Action is the engine's next move, decided from the folded journal alone.
type Action struct {
	Kind   ActionKind
	Step   string
	Visit  int
	End    workflow.End
	Detail string
}

// Decide is the pure transition: from a definition and a folded journal it
// says what to do next. An ended step is never run again; a step entered
// and not ended is continued without a new visit.
func Decide(def *workflow.Definition, st State) (Action, error) {
	if st.Ended {
		return Action{Kind: ActFinished}, nil
	}
	if st.Current == "" {
		if len(st.Steps) == 0 {
			return Action{}, fmt.Errorf("run %s: the journal has no step to start at", st.ID)
		}
		last := st.Steps[len(st.Steps)-1]
		end, detail := endOf(def, last.Step, last.Outcome, last.Detail)
		return Action{Kind: ActFinish, Step: last.Step, End: end, Detail: detail}, nil
	}
	step, ok := def.Steps[st.Current]
	if !ok {
		return Action{}, fmt.Errorf("run %s: step %q is not in workflow %q", st.ID, st.Current, def.Name)
	}
	if st.CurrentVisit > 0 {
		return Action{Kind: ActContinue, Step: st.Current, Visit: st.CurrentVisit}, nil
	}
	visit := st.Visits[st.Current] + 1
	if visit > step.Cap() {
		return Action{
			Kind: ActOverCap, Step: st.Current, Visit: visit, End: workflow.EndFailed,
			Detail: fmt.Sprintf("step %q went over max_visits %d", st.Current, step.Cap()),
		}, nil
	}
	return Action{Kind: ActEnter, Step: st.Current, Visit: visit}, nil
}

// endOf is how a run ends when step ended with outcome and the transition
// leaves the workflow, with the finished entry's detail. It is only
// meaningful when the transition has an End.
func endOf(def *workflow.Definition, step, outcome, detail string) (workflow.End, string) {
	if outcome == OutcomeError {
		return workflow.EndFailed, detail
	}
	_, end := def.Next(step, outcome)
	if end == workflow.EndFailed {
		return end, fmt.Sprintf("step %q ended with outcome %q", step, outcome)
	}
	return end, ""
}

// Run runs the run with full ID or unique prefix id until it finishes, or
// stops at a wait under NoWait. It returns the folded state either way. A
// run another process holds returns ErrLive; a finished one ErrFinished.
func (e *Engine) Run(id string) (State, error) {
	j, err := e.Store.Read(id)
	if err != nil {
		return State{}, err
	}
	unlock, ok, err := e.Store.TryLock(j.ID)
	if err != nil {
		return State{}, err
	}
	if !ok {
		return j.State(), &runError{id: j.ID, kind: ErrLive}
	}
	defer unlock()
	// Read again under the lock: the holder may have moved on.
	if j, err = e.Store.Read(j.ID); err != nil {
		return State{}, err
	}
	st := j.State()
	if st.Ended {
		return st, &runError{id: j.ID, kind: ErrFinished}
	}
	def, err := definition(st)
	if err != nil {
		return st, err
	}
	r := &running{e: e, j: j, def: def}
	return r.loop()
}

// definition parses the definition the run started with.
func definition(st State) (*workflow.Definition, error) {
	def, err := workflow.Parse(st.Workflow, []byte(st.Definition))
	if err != nil {
		return nil, fmt.Errorf("run %s: its workflow no longer parses: %w", st.ID, err)
	}
	return def, nil
}

// running is one Engine.Run call holding a run's lock.
type running struct {
	e   *Engine
	j   Journal
	def *workflow.Definition
}

func (r *running) clock() Clock {
	if r.e.Clock != nil {
		return r.e.Clock
	}
	return RealClock{}
}

func (r *running) now() time.Time {
	if r.e.Now != nil {
		return r.e.Now()
	}
	return time.Now()
}

func (r *running) log(format string, args ...any) {
	if r.e.Log != nil {
		r.e.Log(fmt.Sprintf(format, args...))
	}
}

// append journals one entry, stamped with Now.
func (r *running) append(en Entry) error {
	en.At = r.now()
	if err := r.e.Store.Append(r.j.ID, en); err != nil {
		return err
	}
	r.j.Entries = append(r.j.Entries, en)
	return nil
}

func (r *running) finish(step string, end workflow.End, detail string) (State, error) {
	if err := r.append(Entry{Type: TypeFinished, Step: step, End: end.String(), Detail: detail}); err != nil {
		return r.j.State(), err
	}
	if detail != "" {
		r.log("run %s %s: %s", r.j.ID, end, detail)
	} else {
		r.log("run %s %s", r.j.ID, end)
	}
	return r.j.State(), nil
}

func (r *running) loop() (State, error) {
	for {
		st := r.j.State()
		act, err := Decide(r.def, st)
		if err != nil {
			return st, err
		}
		switch act.Kind {
		case ActFinished:
			return st, nil
		case ActFinish, ActOverCap:
			return r.finish(act.Step, act.End, act.Detail)
		case ActEnter:
			if err := r.append(Entry{Type: TypeEntered, Step: act.Step, Visit: act.Visit}); err != nil {
				return r.j.State(), err
			}
			r.log("step %s: entered, visit %d", act.Step, act.Visit)
		case ActContinue:
			r.log("step %s: continuing visit %d", act.Step, act.Visit)
		}
		step := r.def.Steps[act.Step]
		res, err := r.runStep(step, st.Change)
		if err != nil {
			return r.j.State(), err
		}
		if res.stopped != "" {
			if err := r.append(Entry{Type: TypeStopped, Step: step.Name, Detail: res.stopped}); err != nil {
				return r.j.State(), err
			}
			r.log("step %s: stopped, %s", step.Name, res.stopped)
			return r.j.State(), nil
		}
		next, end := r.def.Next(step.Name, res.outcome)
		if res.outcome == OutcomeError {
			next, end = "", workflow.EndFailed
		}
		if err := r.append(Entry{
			Type: TypeEnded, Step: step.Name, Visit: act.Visit, Outcome: res.outcome,
			Next: next, End: end.String(), Detail: res.detail,
		}); err != nil {
			return r.j.State(), err
		}
		r.log("step %s: %s", step.Name, res.outcome)
		// The next turn of the loop enters next, or (via ActFinish) writes
		// the finished entry, the same path a resume takes.
	}
}

// result is how a step's run went: an outcome, or stopped with a reason.
type result struct {
	outcome string
	detail  string
	stopped string
}

func errResult(err error) result { return result{outcome: OutcomeError, detail: err.Error()} }

// runStep runs one step by kind. A Runner failure is an error outcome, not
// an error: it ends the run. The error return is for journal failures.
func (r *running) runStep(step workflow.Step, change string) (result, error) {
	env := []string{
		"RHIZA_CHANGE=" + change,
		"RHIZA_RUN=" + r.j.ID,
		"RHIZA_STEP=" + step.Name,
	}
	switch step.Kind {
	case workflow.KindCommand:
		return r.command(step, env), nil
	case workflow.KindWait:
		return r.wait(step, env), nil
	case workflow.KindAgent:
		return r.agent(step, env)
	}
	return errResult(fmt.Errorf("step %q: unknown kind %q", step.Name, step.Kind)), nil
}

func (r *running) command(step workflow.Step, env []string) result {
	code, detail, err := r.e.Runner.Command(step, env)
	if err != nil {
		return errResult(err)
	}
	d := fmt.Sprintf("exit %d", code)
	if detail != "" {
		d += ": " + detail
	}
	if code == 0 {
		return result{outcome: workflow.OutcomePass, detail: d}
	}
	return result{outcome: workflow.OutcomeFail, detail: d}
}

// wait checks the condition until it holds or the timeout passes, measured
// with the Clock from when this call began polling. Under NoWait a
// condition that does not hold stops the run instead.
func (r *running) wait(step workflow.Step, env []string) result {
	c := r.clock()
	interval := step.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	began := c.Now()
	for polls := 1; ; polls++ {
		held, err := r.e.Runner.Holds(step, env)
		if err != nil {
			return errResult(err)
		}
		if held {
			return result{outcome: workflow.OutcomePass, detail: fmt.Sprintf("held after %d checks", polls)}
		}
		if r.e.NoWait {
			return result{stopped: "waiting: " + step.Until}
		}
		elapsed := c.Now().Sub(began)
		if elapsed >= step.Timeout {
			return result{outcome: workflow.OutcomeTimeout, detail: fmt.Sprintf("did not hold within %s", step.Timeout)}
		}
		// Sleep at most to the deadline, so the last check falls on it.
		c.Sleep(min(interval, step.Timeout-elapsed))
	}
}

// agent launches the step's session, unless the journal already has one
// for this visit (a resume), and polls it until it ends.
func (r *running) agent(step workflow.Step, env []string) (result, error) {
	session := r.j.State().Session[step.Name]
	if session == "" {
		s, err := r.e.Runner.StartAgent(step, env)
		if err != nil {
			return errResult(err), nil
		}
		if err := r.append(Entry{Type: TypeSession, Step: step.Name, Session: s}); err != nil {
			return result{}, err
		}
		r.log("step %s: session %s", step.Name, s)
		session = s
	} else {
		r.log("step %s: rejoining session %s", step.Name, session)
	}
	c := r.clock()
	interval := step.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	for {
		status, err := r.e.Runner.AgentStatus(session)
		if err != nil {
			// The step ends with an error outcome: end its block first.
			if err := r.unblock(step.Name, session); err != nil {
				return result{}, err
			}
			return errResult(err), nil
		}
		if status.Done {
			if err := r.unblock(step.Name, session); err != nil {
				return result{}, err
			}
			outcome := status.Outcome
			if outcome == "" {
				outcome = workflow.OutcomeFail
			}
			return result{outcome: outcome, detail: status.Detail}, nil
		}
		if status.Blocked != nil {
			err = r.block(step.Name, session, status.Blocked)
		} else {
			err = r.unblock(step.Name, session)
		}
		if err != nil {
			return result{}, err
		}
		c.Sleep(interval)
	}
}

// blockedOn reports whether the journal shows session blocked.
func (r *running) blockedOn(session string) (State, bool) {
	st := r.j.State()
	return st, st.Blocked != nil && st.BlockedSession == session
}

// block journals that step's session became blocked on b and notifies,
// unless the journal already shows it blocked: one notice per block, a
// resume included.
func (r *running) block(step, session string, b *Block) error {
	st, already := r.blockedOn(session)
	if already {
		return nil
	}
	cp := *b
	if err := r.append(Entry{Type: TypeBlocked, Step: step, Session: session, Block: &cp}); err != nil {
		return err
	}
	on := cp.Kind
	if on == "" {
		on = "the owner"
	}
	r.log("step %s: blocked on %s", step, on)
	r.notify(st, step, session, true, &cp)
	return nil
}

// unblock journals the end of session's block and notifies, when the
// journal shows it blocked; otherwise it does nothing.
func (r *running) unblock(step, session string) error {
	st, blocked := r.blockedOn(session)
	if !blocked {
		return nil
	}
	if err := r.append(Entry{Type: TypeUnblocked, Step: step, Session: session}); err != nil {
		return err
	}
	r.log("step %s: unblocked", step)
	r.notify(st, step, session, false, st.Blocked)
	return nil
}

func (r *running) notify(st State, step, session string, blocked bool, b *Block) {
	if r.e.Notify == nil {
		return
	}
	cp := *b
	r.e.Notify(Notice{
		Run: r.j.ID, Change: st.Change, Title: st.Title, Step: step,
		Session: session, Blocked: blocked, Block: &cp,
	})
}
