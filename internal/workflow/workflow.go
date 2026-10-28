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
	// KindLeaves conducts the run change set's leaves: a builder session
	// and then a verifier session for each leaf.
	KindLeaves = "leaves"
	// KindChecks runs a review's checks: a fresh session per check, all
	// in parallel, then a further cycle of the checks that failed.
	KindChecks = "checks"
	// KindSignoff is a review's last step: it waits for the owner's
	// sign-off by stopping the run; rhiza run approve passes it and rhiza
	// run request-changes ends it with outcome changes, which goes back
	// to a checks step.
	KindSignoff = "signoff"
	// KindWorkflow runs another workflow's steps inside this run. It is
	// expanded when the file is parsed: the definition holds the inner
	// workflow's steps, named <workflow>/<step>, and no workflow step.
	KindWorkflow = "workflow"
)

// Leaves step roles: the sub-tables of a leaves step. A checks step has a
// verifier too.
const (
	RoleBuilder  = "builder"
	RoleVerifier = "verifier"
)

// Checks step roles: a check's session, and the sub-tables of a checks
// step beside its verifier.
const (
	RoleCheck  = "check"
	RoleCritic = "critic"
	RoleFixer  = "fixer"
)

// Split values: how a leaves step lays out its panes.
const (
	SplitAuto  = "auto"
	SplitRight = "right"
	SplitBelow = "below"
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
	// OutcomeChanges is a signoff step's outcome when the owner asked
	// for changes: its on.changes names the checks step that fixes them.
	OutcomeChanges = "changes"
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

// Placeholders a leaves step's builder or verifier prompt may use, and a
// checks step's verifier.
var leafPlaceholders = append(slices.Clone(placeholders), "leaf", "attempt")

// Placeholders a checks step's check or critic prompt may use.
var checkPlaceholders = append(slices.Clone(placeholders), "check", "cycle")

// Placeholders a checks step's fixer prompt may use; {leaf} is empty for
// the fixer of findings that name no leaf.
var fixerPlaceholders = append(slices.Clone(placeholders), "leaf", "cycle")

// kindOutcomes lists the outcomes each step kind can produce.
var kindOutcomes = map[string][]string{
	KindAgent:   {OutcomePass, OutcomeFail},
	KindWait:    {OutcomePass, OutcomeTimeout},
	KindCommand: {OutcomePass, OutcomeFail},
	KindArbiter: {OutcomePass, OutcomeFail},
	KindLeaves:  {OutcomePass, OutcomeFail},
	KindChecks:  {OutcomePass, OutcomeFail},
	KindSignoff: {OutcomePass, OutcomeChanges},
	// A workflow step passes when the inner workflow ends done and fails
	// when it ends failed.
	KindWorkflow: {OutcomePass, OutcomeFail},
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
	// Includes holds the source of every workflow a workflow step named,
	// directly or through another, by name: what a run must keep to parse
	// this definition again as it was.
	Includes map[string]string

	// declared describes each step the file declares, by its name, in
	// terms of the flat steps: see declStep.
	declared map[string]declStep
}

// declStep is a declared step as seen by a workflow that includes this
// one: where it is entered and where it goes when it passes ("done" when
// the workflow ends).
type declStep struct{ entry, pass string }

// Resolver gives the source of the workflow called name, for a workflow
// step.
type Resolver func(name string) ([]byte, error)

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
	// Models is an agent step's model list; empty means the agent's
	// default.
	Models []string
	// Builder and Verifier are a leaves step's roles; nil on other steps.
	// A checks step has a Verifier too, for a leaf a fixer fixed.
	Builder, Verifier *Role
	// Checks are a checks step's checks, in file order; nil on other
	// steps.
	Checks []Check
	// Critic and Fixer are a checks step's roles; nil on other steps.
	Critic, Fixer *Role
	// Split is how a leaves or checks step lays out its panes: SplitAuto
	// (the default), SplitRight or SplitBelow. Empty on other steps.
	Split string
	// MaxVisits caps how often a run may enter the step; 0 means 1.
	MaxVisits int
	// Workflow and Skip are a workflow step's included workflow and the
	// steps of it to take as passed. They are gone from a parsed
	// definition, whose workflow steps are expanded.
	Workflow string
	Skip     []string
	// On maps an outcome to the next step, or to done or failed.
	On map[string]string
}

// Role is one session kind of a leaves step: its builder or its verifier.
type Role struct {
	// Prompt is the role's first prompt to its agent, a template over
	// {change}, {slug}, {title}, {run}, {step}, {leaf} and {attempt}.
	Prompt string
	// Agent is the role's adapter name; empty means the run's.
	Agent string
	// PermissionMode is the role's permission mode; empty means the run's.
	PermissionMode string
	// Models is the role's model by attempt: attempt n takes the nth, and
	// the last stays for every later attempt.
	Models []string
}

// Check is one check of a checks step: its name, unique in the step, and
// its session's role. A check's Models are by cycle, as a leaves role's
// are by attempt.
type Check struct {
	Name string
	Role
}

// Model is the model for attempt (1-based): Models[min(attempt, len)-1],
// or "" when the role names no model.
func (r Role) Model(attempt int) string {
	if len(r.Models) == 0 {
		return ""
	}
	return r.Models[max(1, min(attempt, len(r.Models)))-1]
}

// Model is the model for visit (1-based) of an agent step: its Models
// laid out by visit as a Role's are by attempt, or "" when the step names
// no model.
func (s Step) Model(visit int) string { return Role{Models: s.Models}.Model(visit) }

// LeafVars are the values a leaves role's prompt placeholders take: a
// step's, plus the leaf's handle and the attempt.
type LeafVars struct {
	PromptVars
	Leaf    string
	Attempt int
}

// RenderPrompt fills the role's prompt placeholders from v. Validation has
// already refused an unknown placeholder.
func (r Role) RenderPrompt(v LeafVars) string {
	return renderPlaceholders(r.Prompt, map[string]string{
		"change": v.Change, "slug": v.Slug, "title": v.Title, "run": v.Run, "step": v.Step,
		"leaf": v.Leaf, "attempt": fmt.Sprint(v.Attempt),
	})
}

// CheckVars are the values a checks step's check or critic prompt
// placeholders take: a step's, plus the check's name and the cycle.
type CheckVars struct {
	PromptVars
	Check string
	Cycle int
}

// RenderCheckPrompt fills a check's or critic's prompt placeholders from
// v. Validation has already refused an unknown placeholder.
func (r Role) RenderCheckPrompt(v CheckVars) string {
	return renderPlaceholders(r.Prompt, map[string]string{
		"change": v.Change, "slug": v.Slug, "title": v.Title, "run": v.Run, "step": v.Step,
		"check": v.Check, "cycle": fmt.Sprint(v.Cycle),
	})
}

// FixerVars are the values a checks step's fixer prompt placeholders
// take: a step's, plus the leaf's handle ("" for findings naming no leaf)
// and the cycle.
type FixerVars struct {
	PromptVars
	Leaf  string
	Cycle int
}

// RenderFixerPrompt fills a fixer's prompt placeholders from v.
// Validation has already refused an unknown placeholder.
func (r Role) RenderFixerPrompt(v FixerVars) string {
	return renderPlaceholders(r.Prompt, map[string]string{
		"change": v.Change, "slug": v.Slug, "title": v.Title, "run": v.Run, "step": v.Step,
		"leaf": v.Leaf, "cycle": fmt.Sprint(v.Cycle),
	})
}

// Role is a leaves step's builder or verifier, or a checks step's critic,
// fixer or verifier, by name; nil for any other name or step. A check's
// role is its Check's.
func (s Step) Role(name string) *Role {
	switch name {
	case RoleBuilder:
		return s.Builder
	case RoleVerifier:
		return s.Verifier
	case RoleCritic:
		return s.Critic
	case RoleFixer:
		return s.Fixer
	}
	return nil
}

// Check is the checks step's check named name, if it has one.
func (s Step) Check(name string) (Check, bool) {
	for _, c := range s.Checks {
		if c.Name == name {
			return c, true
		}
	}
	return Check{}, false
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
	return renderPlaceholders(s.Prompt, map[string]string{
		"change": v.Change, "slug": v.Slug, "title": v.Title, "run": v.Run, "step": v.Step,
	})
}

// renderPlaceholders replaces each {name} in prompt that values has, and
// leaves the rest.
func renderPlaceholders(prompt string, values map[string]string) string {
	return placeholderPattern.ReplaceAllStringFunc(prompt, func(m string) string {
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
	Models         []string          `toml:"models"`
	Workflow       string            `toml:"workflow"`
	Skip           []string          `toml:"skip"`
	Split          string            `toml:"split"`
	Builder        *fileRole         `toml:"builder"`
	Verifier       *fileRole         `toml:"verifier"`
	Check          []fileCheck       `toml:"check"`
	Critic         *fileRole         `toml:"critic"`
	Fixer          *fileRole         `toml:"fixer"`
	On             map[string]string `toml:"on"`
}

// fileCheck is the TOML shape of one [[step.<name>.check]] of a checks
// step.
type fileCheck struct {
	Name           string   `toml:"name"`
	Prompt         string   `toml:"prompt"`
	Agent          string   `toml:"agent"`
	PermissionMode string   `toml:"permission_mode"`
	Models         []string `toml:"models"`
}

// fileRole is the TOML shape of a leaves step's builder or verifier.
type fileRole struct {
	Prompt         string   `toml:"prompt"`
	Agent          string   `toml:"agent"`
	PermissionMode string   `toml:"permission_mode"`
	Models         []string `toml:"models"`
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
	return ParseWith(name, data, nil)
}

// ParseWith parses as Parse does, resolving each workflow step's workflow
// with resolve; without one a workflow step is an error. The included
// workflows' steps replace the workflow step, named <workflow>/<step>.
func ParseWith(name string, data []byte, resolve Resolver) (*Definition, error) {
	includes := map[string]string{}
	d, err := parse(name, data, resolve, []string{name}, includes)
	if err != nil {
		return nil, err
	}
	d.Includes = includes
	return d, nil
}

func parse(name string, data []byte, resolve Resolver, stack []string, includes map[string]string) (*Definition, error) {
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

	start := d.Order[0]
	if f.Start != "" {
		if _, ok := d.Steps[f.Start]; !ok {
			return nil, fmt.Errorf("start names %q, which is not a step; want one of %s", f.Start, strings.Join(d.Order, ", "))
		}
		start = f.Start
	}
	if d.Steps[start].Kind == KindArbiter {
		return nil, fmt.Errorf("step %q is an arbiter step and cannot be the start step: the engine calls it on a question; start at another step", start)
	}
	if err := d.expand(start, resolve, stack, includes); err != nil {
		return nil, err
	}
	return d, nil
}

// unknownKey names a key the format does not know, by step when it is in one.
func unknownKey(k toml.Key) error {
	if len(k) >= 4 && k[0] == "step" && (k[2] == RoleBuilder || k[2] == RoleVerifier || k[2] == RoleCritic || k[2] == RoleFixer) {
		return fmt.Errorf("step %q: %s: unknown key %q; want prompt, agent, permission_mode or models", k[1], k[2], toml.Key(k[3:]).String())
	}
	if len(k) >= 4 && k[0] == "step" && k[2] == RoleCheck {
		return fmt.Errorf("step %q: check: unknown key %q; want name, prompt, agent, permission_mode or models", k[1], toml.Key(k[3:]).String())
	}
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
		Models:         fs.Models,
		On:             fs.On,
	}
	if s.On == nil {
		s.On = map[string]string{}
	}

	switch fs.Kind {
	case KindAgent, KindWait, KindCommand, KindArbiter, KindLeaves, KindChecks, KindSignoff, KindWorkflow:
	case "":
		return fail("no kind: want agent, wait, command, arbiter, leaves, checks, signoff or workflow")
	default:
		return fail("unknown kind %q: want agent, wait, command, arbiter, leaves, checks, signoff or workflow", fs.Kind)
	}
	if fs.Kind != KindWorkflow {
		if fs.Workflow != "" {
			return fail("workflow is only for workflow steps; %s step cannot have it", article(fs.Kind))
		}
		if fs.Skip != nil {
			return fail("skip is only for workflow steps; %s step cannot have it", article(fs.Kind))
		}
	}
	if fs.Models != nil && fs.Kind != KindAgent {
		if fs.Kind == KindChecks {
			return fail("a checks step cannot have models: each check and role has its own; set it under [[step.%s.check]] or the role's table", name)
		}
		if fs.Kind == KindSignoff {
			return fail("a signoff step cannot have models: it runs no agent, it waits for the owner; remove it")
		}
		if fs.Kind == KindLeaves {
			return fail("models is only for agent steps and a leaves step's builder and verifier; set it under [step.%s.builder] or [step.%s.verifier]", name, name)
		}
		return fail("models is only for agent steps and a leaves step's builder and verifier; %s step cannot have it", article(fs.Kind))
	}
	if fs.Split != "" && fs.Kind != KindLeaves && fs.Kind != KindChecks {
		return fail("split is only for leaves and checks steps; %s step cannot have it", article(fs.Kind))
	}
	if fs.Kind != KindLeaves && fs.Kind != KindChecks {
		for _, r := range []struct {
			role, kinds string
			set         bool
		}{{RoleBuilder, "leaves steps", fs.Builder != nil}, {RoleVerifier, "leaves and checks steps", fs.Verifier != nil}} {
			if r.set {
				return fail("%s is only for %s; %s step cannot have it", r.role, r.kinds, article(fs.Kind))
			}
		}
	}
	if fs.Kind != KindChecks {
		for _, r := range []struct {
			role string
			set  bool
		}{{RoleCheck, fs.Check != nil}, {RoleCritic, fs.Critic != nil}, {RoleFixer, fs.Fixer != nil}} {
			if r.set {
				return fail("%s is only for checks steps; %s step cannot have it", r.role, article(fs.Kind))
			}
		}
	}
	if err := checkModels(fs.Models); err != nil {
		return fail("models %v", err)
	}
	switch fs.Kind {
	case KindArbiter:
		return buildArbiter(s, fs, fail)
	case KindLeaves:
		return buildLeaves(s, fs, fail)
	case KindChecks:
		return buildChecks(s, fs, fail)
	case KindSignoff:
		return buildSignoff(s, fs, fail)
	case KindWorkflow:
		return buildWorkflow(s, fs, fail)
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

// buildLeaves finishes a leaves step: it needs a builder and a verifier,
// each with a prompt, takes split (SplitAuto when unset), interval
// (DefaultInterval when unset) and max_visits, and refuses the session keys
// that belong in its sub-tables.
func buildLeaves(s Step, fs fileStep, fail func(string, ...any) (Step, error)) (Step, error) {
	perRole := fmt.Sprintf("each role has its own; set it under [step.%s.builder] or [step.%s.verifier]", s.Name, s.Name)
	refused := []struct {
		key string
		set bool
		why string
	}{
		{"prompt", fs.Prompt != "", perRole},
		{"agent", fs.Agent != "", perRole},
		{"permission_mode", fs.PermissionMode != "", perRole},
		{"run", fs.Run != "", "it runs agents, not a command"},
		{"until", fs.Until != "", "it ends when every leaf has passed or failed"},
		{"until_when", fs.UntilWhen != "", "it ends when every leaf has passed or failed"},
	}
	for _, r := range refused {
		if r.set {
			return fail("a leaves step cannot have %s: %s; remove it", r.key, r.why)
		}
	}
	s.Split = SplitAuto
	if fs.Split != "" {
		if fs.Split != SplitAuto && fs.Split != SplitRight && fs.Split != SplitBelow {
			return fail("split %q: want %q (by the tab's shape), %q or %q", fs.Split, SplitAuto, SplitRight, SplitBelow)
		}
		s.Split = fs.Split
	}
	for _, r := range []struct {
		name string
		fr   *fileRole
		dst  **Role
	}{{RoleBuilder, fs.Builder, &s.Builder}, {RoleVerifier, fs.Verifier, &s.Verifier}} {
		if r.fr == nil {
			return fail("a leaves step needs a [step.%s.%s] table with prompt = \"<what the %s should do>\"", s.Name, r.name, r.name)
		}
		if strings.TrimSpace(r.fr.Prompt) == "" {
			return fail("%s needs prompt = \"<what the %s should do>\"; it may use {%s}", r.name, r.name, strings.Join(leafPlaceholders, "}, {"))
		}
		for _, m := range placeholderPattern.FindAllStringSubmatch(r.fr.Prompt, -1) {
			if !slices.Contains(leafPlaceholders, m[1]) {
				return fail("%s: prompt uses unknown placeholder {%s}; want one of {%s}", r.name, m[1], strings.Join(leafPlaceholders, "}, {"))
			}
		}
		if err := checkModels(r.fr.Models); err != nil {
			return fail("%s: models %v", r.name, err)
		}
		*r.dst = &Role{Prompt: r.fr.Prompt, Agent: r.fr.Agent, PermissionMode: r.fr.PermissionMode, Models: r.fr.Models}
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

// buildChecks finishes a checks step: it needs at least one check, each
// with a valid name of its own and a prompt, and a critic, a fixer and a
// verifier, each with a prompt; it takes split (SplitAuto when unset),
// interval (DefaultInterval when unset) and max_visits, and refuses the
// session keys that belong in its checks and roles, and a builder.
func buildChecks(s Step, fs fileStep, fail func(string, ...any) (Step, error)) (Step, error) {
	own := "each check and role has its own"
	ends := "it ends when its checks have passed or run out of cycles"
	refused := []struct {
		key string
		set bool
		why string
	}{
		{"prompt", fs.Prompt != "", own},
		{"agent", fs.Agent != "", own},
		{"permission_mode", fs.PermissionMode != "", own},
		{"run", fs.Run != "", "it runs agents, not a command"},
		{"until", fs.Until != "", ends},
		{"until_when", fs.UntilWhen != "", ends},
		{"timeout", fs.Timeout != "", ends},
		{"builder", fs.Builder != nil, "a fixer fixes what its checks find"},
	}
	for _, r := range refused {
		if r.set {
			return fail("a checks step cannot have %s: %s; remove it", r.key, r.why)
		}
	}
	if fs.MaxVisits < 0 {
		return fail("max_visits is %d; want 0 or more (0 means 1)", fs.MaxVisits)
	}
	s.Split = SplitAuto
	if fs.Split != "" {
		if fs.Split != SplitAuto && fs.Split != SplitRight && fs.Split != SplitBelow {
			return fail("split %q: want %q (by the tab's shape), %q or %q", fs.Split, SplitAuto, SplitRight, SplitBelow)
		}
		s.Split = fs.Split
	}
	if len(fs.Check) == 0 {
		return fail("a checks step needs at least one [[step.%s.check]] table with name = \"...\" and prompt = \"...\"", s.Name)
	}
	for i, fc := range fs.Check {
		if fc.Name == "" {
			return fail("check %d has no name: set name = \"...\"", i+1)
		}
		if err := checkName("check", fc.Name); err != nil {
			return fail("%v", err)
		}
		if _, dup := s.Check(fc.Name); dup {
			return fail("check %q is declared twice; give each check its own name", fc.Name)
		}
		r, err := buildRole(fmt.Sprintf("check %q", fc.Name), "check", fileRole{Prompt: fc.Prompt, Agent: fc.Agent, PermissionMode: fc.PermissionMode, Models: fc.Models}, checkPlaceholders)
		if err != nil {
			return fail("%v", err)
		}
		s.Checks = append(s.Checks, Check{Name: fc.Name, Role: r})
	}
	for _, r := range []struct {
		name  string
		fr    *fileRole
		dst   **Role
		names []string
	}{
		{RoleCritic, fs.Critic, &s.Critic, checkPlaceholders},
		{RoleFixer, fs.Fixer, &s.Fixer, fixerPlaceholders},
		{RoleVerifier, fs.Verifier, &s.Verifier, leafPlaceholders},
	} {
		if r.fr == nil {
			return fail("a checks step needs a [step.%s.%s] table with prompt = \"<what the %s should do>\"", s.Name, r.name, r.name)
		}
		role, err := buildRole(r.name, r.name, *r.fr, r.names)
		if err != nil {
			return fail("%v", err)
		}
		*r.dst = &role
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

// buildSignoff finishes a signoff step: it takes max_visits and on., and
// refuses every key of a session or a condition, as it runs nothing: it
// waits for the owner. checkTargets checks its on.changes.
func buildSignoff(s Step, fs fileStep, fail func(string, ...any) (Step, error)) (Step, error) {
	waits := "it runs nothing: it waits for the owner's rhiza run approve or rhiza run request-changes"
	refused := []struct {
		key string
		set bool
	}{
		{"prompt", fs.Prompt != ""},
		{"agent", fs.Agent != ""},
		{"permission_mode", fs.PermissionMode != ""},
		{"run", fs.Run != ""},
		{"until", fs.Until != ""},
		{"until_when", fs.UntilWhen != ""},
		{"timeout", fs.Timeout != ""},
		{"interval", fs.Interval != ""},
	}
	for _, r := range refused {
		if r.set {
			return fail("a signoff step cannot have %s: %s; remove it", r.key, waits)
		}
	}
	if fs.MaxVisits < 0 {
		return fail("max_visits is %d; want 0 or more (0 means 1)", fs.MaxVisits)
	}
	return s, nil
}

// buildRole checks one session role of a checks step, named what in its
// errors: a prompt using only names among its placeholders, and models.
// does is what the prompt is asked to say the role should do.
func buildRole(what, does string, fr fileRole, names []string) (Role, error) {
	if strings.TrimSpace(fr.Prompt) == "" {
		return Role{}, fmt.Errorf("%s needs prompt = \"<what the %s should do>\"; it may use {%s}", what, does, strings.Join(names, "}, {"))
	}
	for _, m := range placeholderPattern.FindAllStringSubmatch(fr.Prompt, -1) {
		if !slices.Contains(names, m[1]) {
			return Role{}, fmt.Errorf("%s: prompt uses unknown placeholder {%s}; want one of {%s}", what, m[1], strings.Join(names, "}, {"))
		}
	}
	if err := checkModels(fr.Models); err != nil {
		return Role{}, fmt.Errorf("%s: models %v", what, err)
	}
	return Role{Prompt: fr.Prompt, Agent: fr.Agent, PermissionMode: fr.PermissionMode, Models: fr.Models}, nil
}

// checkModels refuses a blank model name.
func checkModels(models []string) error {
	for i, m := range models {
		if strings.TrimSpace(m) == "" {
			return fmt.Errorf("entry %d is blank; want a model name, like \"opus\"", i+1)
		}
	}
	return nil
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
		if w, _, ok := strings.Cut(t, "/"); ok {
			// A step of an included workflow, checked once expanded.
			if ws, ok := steps[w]; ok && ws.Kind == KindWorkflow {
				continue
			}
		}
		ts, ok := steps[t]
		if !ok {
			return fmt.Errorf("step %q: on.%s names %q, which is not a step; want a step or done or failed", s.Name, o, t)
		}
		if ts.Kind == KindArbiter {
			return fmt.Errorf("step %q: on.%s names %q, an arbiter step; no step goes to it (the engine calls it on a question): want another step or done or failed", s.Name, o, t)
		}
	}
	if s.Kind == KindSignoff {
		t, ok := s.On[OutcomeChanges]
		if !ok {
			return fmt.Errorf("step %q: a signoff step needs on.changes = \"<checks step>\": the checks step that fixes what the owner asks for", s.Name)
		}
		if w, _, cut := strings.Cut(t, "/"); cut && steps[w].Kind == KindWorkflow {
			return nil
		}
		if ts, ok := steps[t]; !ok || ts.Kind != KindChecks {
			return fmt.Errorf("step %q: on.changes names %q, which is not a checks step; want the checks step that fixes what the owner asks for", s.Name, t)
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
