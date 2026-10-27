// Package workflow loads workflow definitions: TOML state machines of named
// steps, where each step ends with an outcome that picks the next step.
//
// The package is pure. It parses bytes and decides transitions; it runs no
// program, writes no file and keeps no journal. An agent step's `agent` is an
// opaque adapter name, resolved elsewhere.
//
// An arbiter step is not a state of the machine: no step goes to it and the
// run never starts at it. The engine calls it, in an autonomous run, to decide
// a question an agent step raised. A file declares at most one, only to set
// its agent, permission mode and interval; its prompt is fixed. Without one,
// Definition.Arbiter gives a built-in arbiter step.
package workflow

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

// Step kinds.
const (
	KindAgent   = "agent"
	KindWait    = "wait"
	KindCommand = "command"
	// KindArbiter is the step an autonomous run's engine calls to decide a
	// question an agent step raised. It is never a transition target.
	KindArbiter = "arbiter"
)

// ArbiterName is the name of the built-in arbiter step, used when a file
// declares none.
const ArbiterName = "arbiter"

// Reserved targets: an `on.<outcome>` naming one of these ends the run.
const (
	TargetDone   = "done"
	TargetFailed = "failed"
)

// Outcomes the step kinds produce.
const (
	OutcomePass    = "pass"
	OutcomeFail    = "fail"
	OutcomeTimeout = "timeout"
)

// When an agent step checks its `until` (the step's until_when).
const (
	// UntilWhenIdle checks it whenever the session sits idle at its
	// prompt. It is the default.
	UntilWhenIdle = "idle"
	// UntilWhenDone checks it only once the agent has said its work is
	// done (`rhiza session done`), or when its session exits 0; idle alone
	// is not enough.
	UntilWhenDone = "done"
)

// DefaultInterval is how often a wait step checks its condition, and an
// agent step its session, when the step sets no interval.
const DefaultInterval = 10 * time.Second

// AgentUntilInterval is an agent step's interval when it has an `until`
// and the file sets no interval.
const AgentUntilInterval = 30 * time.Second

// MaxNameLen is the longest step name allowed.
const MaxNameLen = 30

var (
	namePattern        = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)
	placeholderPattern = regexp.MustCompile(`\{([A-Za-z0-9_]+)\}`)
)

// Placeholders an agent step's prompt may use.
var placeholders = []string{"change", "slug", "title", "run", "step"}

// kindOutcomes lists the outcomes each step kind can produce.
var kindOutcomes = map[string][]string{
	KindAgent:   {OutcomePass, OutcomeFail},
	KindWait:    {OutcomePass, OutcomeTimeout},
	KindCommand: {OutcomePass, OutcomeFail},
	KindArbiter: {OutcomePass, OutcomeFail},
}

// Definition is a loaded, validated workflow.
type Definition struct {
	Name string
	// Start is the step the run enters first: the file's `start`, or the
	// first step in the file.
	Start string
	Steps map[string]Step
	// Order holds the step names in file order; Order[0] is the first step
	// in the file.
	Order []string
}

// Step is one named step of a workflow.
type Step struct {
	Name string
	Kind string
	// Run is a command step's shell command.
	Run string
	// Until is a wait step's condition, a shell command; exit 0 means it
	// holds. On an agent step it is optional: the step also passes when
	// its session is idle at its prompt and the condition holds, leaving
	// the session open.
	Until string
	// Timeout is how long a wait step waits for its condition.
	Timeout time.Duration
	// UntilWhen is when an agent step with an `until` checks it:
	// UntilWhenIdle (also when empty) or UntilWhenDone. Empty on other
	// steps.
	UntilWhen string
	// Interval is how often a wait step checks its condition, or an agent
	// step its session. It is DefaultInterval when the file sets none
	// (AgentUntilInterval on an agent step with until), and zero on a
	// command step.
	Interval time.Duration
	// Agent is an agent or arbiter step's adapter name; empty means the
	// run's, else the default.
	Agent string
	// PermissionMode is an agent or arbiter step's permission mode; empty
	// means the run's.
	PermissionMode string
	// Prompt is an agent step's first prompt to its agent, a template over
	// {change}, {slug}, {title}, {run} and {step}.
	Prompt string
	// MaxVisits caps how often a run may enter the step; 0 means 1.
	MaxVisits int
	// On maps an outcome to the next step, or to done or failed.
	On map[string]string
}

// Cap is how many times a run may enter the step.
func (s Step) Cap() int {
	if s.MaxVisits > 0 {
		return s.MaxVisits
	}
	return 1
}

// Arbiter is the step that decides an autonomous run's questions: the
// file's arbiter step, or a built-in one (the run's agent and mode, checked
// every DefaultInterval) when the file declares none.
func (d *Definition) Arbiter() Step {
	for _, n := range d.Order {
		if s := d.Steps[n]; s.Kind == KindArbiter {
			return s
		}
	}
	return Step{Name: ArbiterName, Kind: KindArbiter, Interval: DefaultInterval, On: map[string]string{}}
}

// WaitsForDone reports whether the step checks its `until` only once its
// agent has said done (or exited 0), never on idle alone.
func (s Step) WaitsForDone() bool { return s.UntilWhen == UntilWhenDone }

// PromptVars are the values an agent step's prompt placeholders take.
type PromptVars struct {
	Change, Slug, Title, Run, Step string
}

// RenderPrompt fills the step's prompt placeholders from v. Validation has
// already refused an unknown placeholder.
func (s Step) RenderPrompt(v PromptVars) string {
	values := map[string]string{
		"change": v.Change, "slug": v.Slug, "title": v.Title, "run": v.Run, "step": v.Step,
	}
	return placeholderPattern.ReplaceAllStringFunc(s.Prompt, func(m string) string {
		if val, ok := values[m[1:len(m)-1]]; ok {
			return val
		}
		return m
	})
}

// End says whether a transition ends the run, and how.
type End int

const (
	// EndNone means the run goes on to the next step.
	EndNone End = iota
	// EndSucceeded means the run ends and succeeds.
	EndSucceeded
	// EndFailed means the run ends and fails.
	EndFailed
)

func (e End) String() string {
	switch e {
	case EndSucceeded:
		return "succeeded"
	case EndFailed:
		return "failed"
	default:
		return ""
	}
}

// Next decides where a run goes after step ends with outcome. An
// `on.<outcome>` entry wins: a step name continues there, done ends the run
// succeeded, failed ends it failed. With no entry, `pass` ends the run
// succeeded and any other outcome ends it failed. A step the definition does
// not have ends the run failed.
func (d *Definition) Next(step, outcome string) (next string, end End) {
	s, ok := d.Steps[step]
	if !ok {
		return "", EndFailed
	}
	target, ok := s.On[outcome]
	if !ok {
		if outcome == OutcomePass {
			return "", EndSucceeded
		}
		return "", EndFailed
	}
	switch target {
	case TargetDone:
		return "", EndSucceeded
	case TargetFailed:
		return "", EndFailed
	}
	return target, EndNone
}

// fileDef is the TOML shape of a workflow file.
type fileDef struct {
	Name  string              `toml:"name"`
	Start string              `toml:"start"`
	Step  map[string]fileStep `toml:"step"`
}

type fileStep struct {
	Kind           string            `toml:"kind"`
	Run            string            `toml:"run"`
	Until          string            `toml:"until"`
	UntilWhen      string            `toml:"until_when"`
	Timeout        string            `toml:"timeout"`
	Interval       string            `toml:"interval"`
	Agent          string            `toml:"agent"`
	PermissionMode string            `toml:"permission_mode"`
	Prompt         string            `toml:"prompt"`
	MaxVisits      int               `toml:"max_visits"`
	On             map[string]string `toml:"on"`
}

// Load reads and parses the workflow file at path. A file without a `name`
// takes its base name, less the extension.
func Load(path string) (*Definition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	d, err := Parse(name, data)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return d, nil
}

// Parse decodes and validates a workflow. name is the workflow's name when
// the file sets none. Any key the format does not know is an error.
func Parse(name string, data []byte) (*Definition, error) {
	var f fileDef
	md, err := toml.Decode(string(data), &f)
	if err != nil {
		return nil, err
	}
	if und := md.Undecoded(); len(und) > 0 {
		return nil, unknownKey(und[0])
	}
	if f.Name != "" {
		name = f.Name
	}
	if name == "" {
		return nil, errors.New("workflow has no name: set name = \"...\"")
	}
	if len(f.Step) == 0 {
		return nil, errors.New("workflow has no steps: add a [step.<name>] table")
	}

	d := &Definition{Name: name, Steps: make(map[string]Step, len(f.Step))}
	d.Order = stepOrder(md.Keys(), f.Step)
	arbiter := ""
	for _, n := range d.Order {
		s, err := buildStep(n, f.Step[n])
		if err != nil {
			return nil, err
		}
		if s.Kind == KindArbiter {
			if arbiter != "" {
				return nil, fmt.Errorf("steps %q and %q are both arbiter steps; a workflow has at most one", arbiter, n)
			}
			arbiter = n
		}
		d.Steps[n] = s
	}
	for _, n := range d.Order {
		if err := checkTargets(d.Steps[n], d.Steps); err != nil {
			return nil, err
		}
	}

	d.Start = d.Order[0]
	if f.Start != "" {
		if _, ok := d.Steps[f.Start]; !ok {
			return nil, fmt.Errorf("start names %q, which is not a step; want one of %s", f.Start, strings.Join(d.Order, ", "))
		}
		d.Start = f.Start
	}
	if d.Steps[d.Start].Kind == KindArbiter {
		return nil, fmt.Errorf("step %q is an arbiter step and cannot be the start step: the engine calls it on a question; start at another step", d.Start)
	}
	return d, nil
}

// unknownKey names a key the format does not know, by step when it is in one.
func unknownKey(k toml.Key) error {
	if len(k) >= 3 && k[0] == "step" {
		return fmt.Errorf("step %q: unknown key %q", k[1], toml.Key(k[2:]).String())
	}
	return fmt.Errorf("unknown key %q", k.String())
}

// stepOrder recovers the steps' file order from the decoder's key order:
// each step at its first `step.<name>` key. Go maps keep no order.
func stepOrder(keys []toml.Key, steps map[string]fileStep) []string {
	seen := make(map[string]bool, len(steps))
	order := make([]string, 0, len(steps))
	for _, k := range keys {
		if len(k) < 2 || k[0] != "step" || seen[k[1]] {
			continue
		}
		if _, ok := steps[k[1]]; !ok {
			continue
		}
		seen[k[1]] = true
		order = append(order, k[1])
	}
	// The decoder reports every key it decodes, so this only guards
	// against a decoder that does not.
	var rest []string
	for n := range steps {
		if !seen[n] {
			rest = append(rest, n)
		}
	}
	sort.Strings(rest)
	return append(order, rest...)
}

// checkName refuses a step or outcome name that does not match the pattern.
func checkName(what, n string) error {
	if !namePattern.MatchString(n) {
		return fmt.Errorf("%s %q: want lowercase letters and digits in words joined by single hyphens, like \"run-tests\"", what, n)
	}
	return nil
}

func buildStep(name string, fs fileStep) (Step, error) {
	fail := func(format string, args ...any) (Step, error) {
		return Step{}, fmt.Errorf("step %q: %s", name, fmt.Sprintf(format, args...))
	}
	if err := checkName("step", name); err != nil {
		return Step{}, err
	}
	if len(name) > MaxNameLen {
		return Step{}, fmt.Errorf("step %q: name is %d characters; want at most %d", name, len(name), MaxNameLen)
	}
	if name == TargetDone || name == TargetFailed {
		return fail("%q is a reserved target; rename the step", name)
	}

	s := Step{
		Name:           name,
		Kind:           fs.Kind,
		Run:            fs.Run,
		Until:          fs.Until,
		UntilWhen:      fs.UntilWhen,
		Agent:          fs.Agent,
		PermissionMode: fs.PermissionMode,
		Prompt:         fs.Prompt,
		MaxVisits:      fs.MaxVisits,
		On:             fs.On,
	}
	if s.On == nil {
		s.On = map[string]string{}
	}

	switch fs.Kind {
	case KindAgent, KindWait, KindCommand:
	case KindArbiter:
		return buildArbiter(s, fs, fail)
	case "":
		return fail("no kind: want agent, wait, command or arbiter")
	default:
		return fail("unknown kind %q: want agent, wait, command or arbiter", fs.Kind)
	}

	if fs.Run != "" && fs.Kind != KindCommand {
		return fail("run is only for command steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Until != "" && fs.Kind == KindCommand {
		return fail("until is only for wait and agent steps; a command step cannot have it")
	}
	if fs.UntilWhen != "" {
		if fs.Kind != KindAgent {
			return fail("until_when is only for agent steps; a %s step cannot have it", fs.Kind)
		}
		if fs.Until == "" {
			return fail("until_when needs until = \"<shell command>\" on the same step: it says when that condition is checked")
		}
		if fs.UntilWhen != UntilWhenIdle && fs.UntilWhen != UntilWhenDone {
			return fail("until_when %q: want %q (check until whenever the session is idle) or %q (only once the agent says it is done)", fs.UntilWhen, UntilWhenIdle, UntilWhenDone)
		}
	}
	if fs.Timeout != "" && fs.Kind != KindWait {
		return fail("timeout is only for wait steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Interval != "" && fs.Kind == KindCommand {
		return fail("interval is only for wait and agent steps; a command step cannot have it")
	}
	if fs.Agent != "" && fs.Kind != KindAgent {
		return fail("agent is only for agent and arbiter steps; a %s step cannot have it", fs.Kind)
	}
	if fs.PermissionMode != "" && fs.Kind != KindAgent {
		return fail("permission_mode is only for agent and arbiter steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Prompt != "" && fs.Kind != KindAgent {
		return fail("prompt is only for agent steps; a %s step cannot have it", fs.Kind)
	}
	if fs.MaxVisits < 0 {
		return fail("max_visits is %d; want 0 or more (0 means 1)", fs.MaxVisits)
	}

	switch fs.Kind {
	case KindAgent:
		if strings.TrimSpace(fs.Prompt) == "" {
			return fail("an agent step needs prompt = \"<what the agent should do>\"; it may use {%s}", strings.Join(placeholders, "}, {"))
		}
		for _, m := range placeholderPattern.FindAllStringSubmatch(fs.Prompt, -1) {
			if !slices.Contains(placeholders, m[1]) {
				return fail("prompt uses unknown placeholder {%s}; want one of {%s}", m[1], strings.Join(placeholders, "}, {"))
			}
		}
	case KindCommand:
		if strings.TrimSpace(fs.Run) == "" {
			return fail("a command step needs run = \"<shell command>\"")
		}
	case KindWait:
		if strings.TrimSpace(fs.Until) == "" {
			return fail("a wait step needs until = \"<shell command>\" (exit 0 means the condition holds)")
		}
		if fs.Timeout == "" {
			return fail("a wait step needs timeout = \"<duration>\", like \"30m\"")
		}
		d, err := positiveDuration(fs.Timeout)
		if err != nil {
			return fail("timeout %v", err)
		}
		s.Timeout = d
	}
	if fs.Kind == KindWait || fs.Kind == KindAgent {
		s.Interval = DefaultInterval
		if fs.Kind == KindAgent && fs.Until != "" {
			s.Interval = AgentUntilInterval
		}
		if fs.Interval != "" {
			d, err := positiveDuration(fs.Interval)
			if err != nil {
				return fail("interval %v", err)
			}
			s.Interval = d
		}
	}
	return s, nil
}

// buildArbiter finishes an arbiter step: it takes agent, permission_mode
// and interval (DefaultInterval when unset) and refuses every other key,
// its prompt above all, which is fixed so that the arbiter is handed only
// the question.
func buildArbiter(s Step, fs fileStep, fail func(string, ...any) (Step, error)) (Step, error) {
	refused := []struct {
		key string
		set bool
		why string
	}{
		{"prompt", fs.Prompt != "", "its prompt is fixed: it names only the question and rhiza guide arbiter"},
		{"run", fs.Run != "", "it runs an agent, not a command"},
		{"until", fs.Until != "", "it ends when its agent says done or exits"},
		{"until_when", fs.UntilWhen != "", "it ends when its agent says done or exits"},
		{"timeout", fs.Timeout != "", "it ends when its agent says done or exits"},
		{"max_visits", fs.MaxVisits != 0, "no step goes to it, so it has no visits"},
	}
	for _, r := range refused {
		if r.set {
			return fail("an arbiter step cannot have %s: %s; remove it", r.key, r.why)
		}
	}
	if len(fs.On) > 0 {
		outcomes := make([]string, 0, len(fs.On))
		for o := range fs.On {
			outcomes = append(outcomes, o)
		}
		sort.Strings(outcomes)
		return fail("on.%s: an arbiter step cannot have on.: its outcome goes back to the step that asked; remove it", outcomes[0])
	}
	s.Interval = DefaultInterval
	if fs.Interval != "" {
		d, err := positiveDuration(fs.Interval)
		if err != nil {
			return fail("interval %v", err)
		}
		s.Interval = d
	}
	return s, nil
}

func positiveDuration(v string) (time.Duration, error) {
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%q is not a duration; want one like \"30s\" or \"10m\"", v)
	}
	if d <= 0 {
		return 0, fmt.Errorf("%q is not positive; want a duration above zero", v)
	}
	return d, nil
}

// checkTargets refuses an outcome name that does not match the pattern or
// that the step's kind cannot produce, a target that is neither a step nor
// done or failed, and a target that is the arbiter step.
func checkTargets(s Step, steps map[string]Step) error {
	outcomes := make([]string, 0, len(s.On))
	for o := range s.On {
		outcomes = append(outcomes, o)
	}
	sort.Strings(outcomes)
	for _, o := range outcomes {
		if err := checkName("outcome", o); err != nil {
			return fmt.Errorf("step %q: on.%s: %w", s.Name, o, err)
		}
		if want := kindOutcomes[s.Kind]; !slices.Contains(want, o) {
			return fmt.Errorf("step %q: on.%s: %s step has no outcome %q; want on.%s", s.Name, o, article(s.Kind), o, strings.Join(want, " or on."))
		}
		t := s.On[o]
		if t == TargetDone || t == TargetFailed {
			continue
		}
		ts, ok := steps[t]
		if !ok {
			return fmt.Errorf("step %q: on.%s names %q, which is not a step; want a step or done or failed", s.Name, o, t)
		}
		if ts.Kind == KindArbiter {
			return fmt.Errorf("step %q: on.%s names %q, an arbiter step; no step goes to it (the engine calls it on a question): want another step or done or failed", s.Name, o, t)
		}
	}
	return nil
}

// article prefixes a step kind with "a" or "an".
func article(kind string) string {
	if kind == KindAgent || kind == KindArbiter {
		return "an " + kind
	}
	return "a " + kind
}
