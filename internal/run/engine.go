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
	// ErrNotFailed means a retry was asked of a run that did not fail:
	// only a failed run can be retried.
	ErrNotFailed = errors.New("run did not fail")
)

// runError names the run and wraps ErrLive, ErrFinished or ErrNotFailed.
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
	case ErrNotFailed:
		return fmt.Sprintf("run %s did not fail: only a failed run can be retried", e.id)
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
	// Holds runs a wait or agent step's `until` condition once.
	Holds(step workflow.Step, env []string) (bool, error)
	// StartAgent launches an agent step's session on model ("" for the
	// agent's default) and returns its ID.
	StartAgent(step workflow.Step, model string, env []string) (session string, err error)
	// AgentStatus reports where a session is: ended, and with which
	// outcome (pass or fail), or still running and perhaps blocked on the
	// owner.
	AgentStatus(session string) (AgentStatus, error)
	// StartArbiter launches the arbiter step's session to decide question
	// node question, on the run's change set, and returns its ID. Its
	// prompt names only the question and how to decide it.
	StartArbiter(step workflow.Step, question string, env []string) (session string, err error)
	// Answer reports whether the model, as the change set arbiter session
	// session writes into sees it (the run's, or the one the question lives
	// in), holds a valid answer to question node question, and if not, why.
	Answer(question, session string) (Answer, error)
	// AddQuestion adds what session asked the owner, text, as one question
	// node in the run's change set, written as the session's actor, and
	// returns the node. It refuses when the change set is finished.
	AddQuestion(session, text string) (Question, error)
	// Deliver gives session the arbiter's answer to question node
	// question: Decision decision, whose statement is statement, as a
	// message once the session waits at its prompt; then the session is no
	// longer blocked on the question.
	Deliver(session, question, decision, statement string) error
	// AgentEnded reports the end of an agent step to its session: it is
	// called once each time an agent step ends with a session, whatever its
	// outcome (pass, fail, timeout or error), the session perhaps left
	// open; never for a step that stopped, nor when StartAgent failed. An
	// error is logged and changes nothing.
	AgentEnded(step workflow.Step, session string) error
}

// Question is a question node the Runner added: its ID and handle.
type Question struct {
	ID     string
	Handle string
}

// Answer is the Runner's verdict on an arbiter's work: OK when a live
// Decision answers the question as an arbiter must (assumed wording and a
// reason). Decision and Statement are that Decision's handle and
// statement; Why says what was missing when it is not OK.
type Answer struct {
	OK        bool
	Decision  string
	Statement string
	Why       string
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
	// Idle is true while the live session sits at its prompt: neither
	// working nor blocked. It is ignored once Done.
	Idle bool
	// SaidDone is true once the agent has said its work is done (its
	// session's done mark). It is ignored once Done.
	SaidDone bool
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
	// already shows the block calls it only when the block ends. An
	// autonomous run's question is the arbiter's, not the owner's: it is
	// never notified.
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
// and not ended is continued without a new visit. A step's max_visits
// counts only its visits since the run's last retried entry (see
// State.CapVisits), so a retry can enter a step whose cap was reached.
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
	if st.CapVisits(st.Current)+1 > step.Cap() {
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
	return e.run(id, false)
}

// Retry reopens a failed run at the step it failed at (State.FailedStep)
// by appending a retried entry, then runs it as Run does. The entry is
// appended under the run's lock, after reading the journal again, so it
// is never written beside another process's run. A run that is not
// finished is run as Run would; one that succeeded returns ErrNotFailed
// and writes nothing.
func (e *Engine) Retry(id string) (State, error) {
	return e.run(id, true)
}

func (e *Engine) run(id string, retry bool) (State, error) {
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
	if st.Ended && !retry {
		return st, &runError{id: j.ID, kind: ErrFinished}
	}
	def, err := definition(st)
	if err != nil {
		return st, err
	}
	r := &running{e: e, j: j, def: def}
	if st.Ended {
		step, err := st.FailedStep()
		if err != nil {
			return st, &runError{id: j.ID, kind: ErrNotFailed}
		}
		if _, ok := def.Steps[step]; !ok {
			return st, fmt.Errorf("run %s: step %q is not in workflow %q", j.ID, step, def.Name)
		}
		if err := r.append(Entry{Type: TypeRetried, Step: step, Detail: st.Detail}); err != nil {
			return r.j.State(), err
		}
		r.log("run %s: retrying step %s", j.ID, step)
	}
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
	env := r.stepEnv(step.Name, change)
	switch step.Kind {
	case workflow.KindCommand:
		res := r.command(step, env)
		if res.outcome == OutcomeError {
			return res, nil
		}
		return r.arbitrateAsked(step, res)
	case workflow.KindWait:
		return r.wait(step, env), nil
	case workflow.KindAgent:
		return r.agent(step, env)
	case workflow.KindLeaves:
		return r.leaves(step, env)
	}
	return errResult(fmt.Errorf("step %q: unknown kind %q", step.Name, step.Kind)), nil
}

// stepEnv is the RHIZA_* variables of step on the run's change set.
func (r *running) stepEnv(step, change string) []string {
	return []string{
		"RHIZA_CHANGE=" + change,
		"RHIZA_RUN=" + r.j.ID,
		"RHIZA_STEP=" + step,
	}
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

// arbitrateAsked has the arbiter decide, in an autonomous run, each
// question a command step's process journalled (an asked entry with no
// session, such as `rhiza run review-check` writes) that it has not
// decided yet, once the step's command has run with result res. The
// process appended to the journal without the run's lock, so the journal
// is read again first. No session asked, so nothing is delivered: the
// answering Decision is in the model. A failed verdict ends the step with
// an error outcome naming the question, so the run fails whatever the
// step's on. entries say; otherwise res stands. A normal run leaves the
// questions open for the owner. The error return is for journal failures.
func (r *running) arbitrateAsked(step workflow.Step, res result) (result, error) {
	if !r.j.State().Autonomous {
		return res, nil
	}
	j, err := r.e.Store.Read(r.j.ID)
	if err != nil {
		return result{}, err
	}
	r.j = j
	for _, asked := range undecidedAsks(r.j.Entries) {
		v, err := r.arbitrate(step, asked.Question)
		if err != nil {
			return result{}, err
		}
		if v.outcome != workflow.OutcomePass {
			return result{
				outcome: OutcomeError,
				detail:  fmt.Sprintf("arbiter failed on question %s (%s): %s", asked.Handle, asked.Question, v.detail),
			}, nil
		}
	}
	return res, nil
}

// undecidedAsks are the session-less asked entries, in journal order, whose
// question has no verdict yet (an open arbitration is rejoined). A verdict
// other than pass counts only until the next retried entry: a retry has the
// arbiter decide the question anew.
func undecidedAsks(entries []Entry) []Asked {
	var asks []Asked
	verdict := map[string]string{}
	for _, e := range entries {
		switch e.Type {
		case TypeAsked:
			if e.Session == "" {
				asks = append(asks, Asked{Step: e.Step, Question: e.Question, Handle: e.Handle, Text: e.Detail, At: e.At})
			}
		case TypeArbitrated:
			verdict[e.Question] = e.Outcome
		case TypeRetried:
			for q, o := range verdict {
				if o != workflow.OutcomePass {
					delete(verdict, q)
				}
			}
		}
	}
	var out []Asked
	seen := map[string]bool{}
	for _, a := range asks {
		if _, ok := verdict[a.Question]; ok || seen[a.Question] {
			continue
		}
		seen[a.Question] = true
		out = append(out, a)
	}
	return out
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
// for this visit (a resume), and polls it until it ends. A step with an
// `until` also passes when a poll finds the session ready and the
// condition holds; the session is left open. Ready is idle at its prompt,
// or, for a step that waits for done (until_when = "done"), that the agent
// has said its work is done: idle alone is then not enough. The condition
// is never checked while the session is blocked on the owner. A step that
// waits for done also checks it when the session exits 0, and passes only
// if it holds; any other exit fails, as on every agent step. Once the step
// ends with an outcome, the Runner hears of it (Runner.AgentEnded).
func (r *running) agent(step workflow.Step, env []string) (result, error) {
	session := r.j.State().Session[step.Name]
	if session == "" {
		s, err := r.e.Runner.StartAgent(step, step.Model(r.j.State().CurrentVisit), env)
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
	res, err := r.pollAgent(step, env, session)
	if err != nil || res.stopped != "" {
		return res, err
	}
	if err := r.e.Runner.AgentEnded(step, session); err != nil {
		r.log("step %s: %v", step.Name, err)
	}
	return res, nil
}

// pollAgent polls the step's session until the step ends (see agent).
func (r *running) pollAgent(step workflow.Step, env []string, session string) (result, error) {
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
			if outcome == workflow.OutcomePass && step.Until != "" && step.WaitsForDone() {
				return r.exitedUntil(step, env, status.Detail), nil
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
		if r.arbitrates(status.Blocked) {
			res, ended, err := r.ask(step, session, status.Blocked)
			if err != nil || ended {
				return res, err
			}
		}
		if ready, why := readyForUntil(step, status); ready {
			held, err := r.e.Runner.Holds(step, env)
			if err != nil {
				return errResult(err), nil
			}
			if held {
				return result{outcome: workflow.OutcomePass, detail: why + " and until held; session left open"}, nil
			}
		}
		c.Sleep(interval)
	}
}

// arbitrates reports whether the engine, not the owner, answers block b:
// only a question, and only in an autonomous run.
func (r *running) arbitrates(b *Block) bool {
	return b != nil && b.Kind == "question" && r.j.State().Autonomous
}

// ask has the arbiter answer the question step's session is blocked on,
// in an autonomous run. Once per block (the journal's asked entry for it,
// State.BlockAsked), it adds the question as a node in the run's change
// set and journals an asked entry; then it arbitrates the question, or
// rejoins its open arbitration (a resume), or reads its verdict. A failed
// verdict ends the step (ended true) with an error outcome naming the
// question, its block ended first, so the run finishes failed whatever
// the step's on. entries say. On a pass the step goes on polling its
// session: the arbiter's answer is the run's, and the session's own
// outcome decides the step. The error return is for journal failures.
func (r *running) ask(step workflow.Step, session string, b *Block) (result, bool, error) {
	asked, ok := r.j.State().BlockAsked()
	if !ok {
		q, err := r.e.Runner.AddQuestion(session, b.Question)
		if err != nil {
			return r.askFailed(step, session, errResult(fmt.Errorf("add the question session %s asked: %v", session, err)))
		}
		asked = Asked{Step: step.Name, Session: session, Question: q.ID, Handle: q.Handle, Text: b.Question}
		if err := r.append(Entry{
			Type: TypeAsked, Step: step.Name, Session: session, Question: q.ID, Handle: q.Handle, Detail: b.Question,
		}); err != nil {
			return result{}, false, err
		}
		r.log("step %s: session %s asked question %s (%s); the arbiter answers it", step.Name, session, q.Handle, q.ID)
	}
	var res result
	if a, ok := r.j.State().Arbitration(asked.Question); ok && !a.Open() {
		res = result{outcome: a.Outcome, detail: a.Detail}
	} else {
		var err error
		if res, err = r.arbitrate(step, asked.Question); err != nil {
			return result{}, false, err
		}
	}
	if res.outcome == workflow.OutcomePass {
		return r.answer(step, session, asked)
	}
	return r.askFailed(step, session, result{
		outcome: OutcomeError,
		detail:  fmt.Sprintf("arbiter failed on question %s (%s): %s", asked.Handle, asked.Question, res.detail),
	})
}

// answer delivers the arbiter's passing answer to asked's question to the
// asking session, journals an answered entry and ends the block; then the
// step goes on polling the session, whose own outcome decides the step.
// When the question already has its answered entry (a resume after the
// delivery), it only ends the block. The statement is read anew through
// Runner.Answer, so the journal holds no model text. A failed read or
// delivery ends the step with an error outcome naming the question.
func (r *running) answer(step workflow.Step, session string, asked Asked) (result, bool, error) {
	if asked.Answered == "" {
		a, _ := r.j.State().Arbitration(asked.Question)
		ans, err := r.e.Runner.Answer(asked.Question, a.Session)
		if err == nil && !ans.OK {
			err = errors.New(ans.Why)
		}
		if err == nil {
			err = r.e.Runner.Deliver(session, asked.Question, a.Decision, ans.Statement)
		}
		if err != nil {
			return r.askFailed(step, session, result{
				outcome: OutcomeError,
				detail:  fmt.Sprintf("could not deliver the answer to question %s (%s): %v", asked.Handle, asked.Question, err),
			})
		}
		if err := r.append(Entry{
			Type: TypeAnswered, Step: step.Name, Session: session, Question: asked.Question, Decision: a.Decision,
		}); err != nil {
			return result{}, false, err
		}
		r.log("step %s: session %s has the answer to question %s: decision %s", step.Name, session, asked.Handle, a.Decision)
	}
	if err := r.unblock(step.Name, session); err != nil {
		return result{}, false, err
	}
	return result{}, false, nil
}

// askFailed ends session's block and returns res as the step's outcome.
func (r *running) askFailed(step workflow.Step, session string, res result) (result, bool, error) {
	if err := r.unblock(step.Name, session); err != nil {
		return result{}, false, err
	}
	return res, true, nil
}

// readyForUntil reports whether a live session's poll lets the step check
// its `until`, and why: never without an `until` or while blocked; for a
// step that waits for done, only once the agent said done; otherwise
// whenever the session is idle.
func readyForUntil(step workflow.Step, status AgentStatus) (bool, string) {
	switch {
	case step.Until == "" || status.Blocked != nil:
		return false, ""
	case step.WaitsForDone():
		return status.SaidDone, "said done"
	default:
		return status.Idle, "idle"
	}
}

// exitedUntil is the outcome of a step that waits for done when its
// session exited 0 (detail): pass if its `until` holds, else fail.
func (r *running) exitedUntil(step workflow.Step, env []string, detail string) result {
	held, err := r.e.Runner.Holds(step, env)
	if err != nil {
		return errResult(err)
	}
	if held {
		return result{outcome: workflow.OutcomePass, detail: detail + " and until held"}
	}
	return result{outcome: workflow.OutcomeFail, detail: detail + " but until did not hold"}
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
	if !r.arbitrates(&cp) {
		r.notify(st, step, session, true, &cp)
	}
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
	if !r.arbitrates(st.Blocked) {
		r.notify(st, step, session, false, st.Blocked)
	}
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

// arbitrate has the run's arbiter (Definition.Arbiter) decide question
// node question, raised by agent step asking, and journals its verdict. It
// rejoins the question's open arbitration (a resume), else launches the
// arbiter's session and journals an arbiter entry. It polls the session at
// the arbiter's interval until the agent says done or the session exits,
// then asks the Runner for the answer: a valid one passes, anything else
// fails with why. A session that blocks, on a question above all, fails at
// once: an arbiter never asks (it would recurse). The session is left open
// either way. Runner failures are an error outcome; the error return is
// for journal failures. The arbitrated entry records the verdict.
func (r *running) arbitrate(asking workflow.Step, question string) (result, error) {
	arb := r.def.Arbiter()
	st := r.j.State()
	session := ""
	if a, ok := st.OpenArbitration(question); ok && a.Session != "" {
		session = a.Session
		r.log("step %s: rejoining arbiter session %s on question %s", asking.Name, session, question)
	} else {
		s, err := r.e.Runner.StartArbiter(arb, question, r.stepEnv(arb.Name, st.Change))
		if err != nil {
			return r.arbitrated(asking, question, errResult(err), "")
		}
		if err := r.append(Entry{Type: TypeArbiter, Step: asking.Name, Session: s, Question: question}); err != nil {
			return result{}, err
		}
		r.log("step %s: arbiter session %s on question %s", asking.Name, s, question)
		session = s
	}
	interval := arb.Interval
	if interval <= 0 {
		interval = workflow.DefaultInterval
	}
	for {
		status, err := r.e.Runner.AgentStatus(session)
		if err != nil {
			return r.arbitrated(asking, question, errResult(err), "")
		}
		if !status.Done && status.Blocked != nil {
			return r.arbitrated(asking, question, result{outcome: workflow.OutcomeFail, detail: arbiterBlocked(status.Blocked)}, "")
		}
		if status.Done || status.SaidDone {
			break
		}
		r.clock().Sleep(interval)
	}
	ans, err := r.e.Runner.Answer(question, session)
	if err != nil {
		return r.arbitrated(asking, question, errResult(err), "")
	}
	if !ans.OK {
		return r.arbitrated(asking, question, result{outcome: workflow.OutcomeFail, detail: ans.Why}, "")
	}
	return r.arbitrated(asking, question, result{outcome: workflow.OutcomePass, detail: "answered by " + ans.Decision}, ans.Decision)
}

// arbiterBlocked is why an arbiter whose session blocked on b fails.
func arbiterBlocked(b *Block) string {
	switch {
	case b.Kind != "" && b.Kind != "question":
		return "the arbiter blocked on " + b.Kind
	case b.Question != "":
		return "the arbiter asked a question: " + b.Question
	}
	return "the arbiter asked a question"
}

// arbitrated journals the arbiter's verdict res on question, asked from
// step asking, with the answering Decision's handle on a pass, and
// returns res.
func (r *running) arbitrated(asking workflow.Step, question string, res result, decision string) (result, error) {
	if err := r.append(Entry{
		Type: TypeArbitrated, Step: asking.Name, Question: question,
		Outcome: res.outcome, Detail: res.detail, Decision: decision,
	}); err != nil {
		return result{}, err
	}
	r.log("step %s: arbiter %s on question %s: %s", asking.Name, res.outcome, question, res.detail)
	return res, nil
}
