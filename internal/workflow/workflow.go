// Package workflow loads workflow definitions: TOML state machines of named
// steps, where each step ends with an outcome that picks the next step.
//
// The package is pure. It parses bytes and decides transitions; it runs no
// program, writes no file and keeps no journal. An agent step's `agent` is an
// opaque adapter name, resolved elsewhere.
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
)

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

// DefaultInterval is how often a wait step checks its condition, and an
// agent step its session, when the step sets no interval.
const DefaultInterval = 10 * time.Second

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
	// holds.
	Until string
	// Timeout is how long a wait step waits for its condition.
	Timeout time.Duration
	// Interval is how often a wait step checks its condition, or an agent
	// step its session. It is DefaultInterval when the file sets none, and
	// zero on a command step.
	Interval time.Duration
	// Agent is an agent step's adapter name; empty means the default.
	Agent string
	// PermissionMode is an agent step's permission mode; empty means the
	// run's.
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
	for _, n := range d.Order {
		s, err := buildStep(n, f.Step[n])
		if err != nil {
			return nil, err
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
	case "":
		return fail("no kind: want agent, wait or command")
	default:
		return fail("unknown kind %q: want agent, wait or command", fs.Kind)
	}

	if fs.Run != "" && fs.Kind != KindCommand {
		return fail("run is only for command steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Until != "" && fs.Kind != KindWait {
		return fail("until is only for wait steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Timeout != "" && fs.Kind != KindWait {
		return fail("timeout is only for wait steps; a %s step cannot have it", fs.Kind)
	}
	if fs.Interval != "" && fs.Kind == KindCommand {
		return fail("interval is only for wait and agent steps; a command step cannot have it")
	}
	if fs.Agent != "" && fs.Kind != KindAgent {
		return fail("agent is only for agent steps; a %s step cannot have it", fs.Kind)
	}
	if fs.PermissionMode != "" && fs.Kind != KindAgent {
		return fail("permission_mode is only for agent steps; a %s step cannot have it", fs.Kind)
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
// that the step's kind cannot produce, and a target that is neither a step
// nor done or failed.
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
		if _, ok := steps[t]; !ok {
			return fmt.Errorf("step %q: on.%s names %q, which is not a step; want a step or done or failed", s.Name, o, t)
		}
	}
	return nil
}

// article prefixes a step kind with "a" or "an".
func article(kind string) string {
	if kind == KindAgent {
		return "an " + kind
	}
	return "a " + kind
}
