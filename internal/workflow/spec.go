package workflow

import "slices"

// This file describes the file format for the generated reference: what each
// step kind is, which keys it takes and which outcomes it produces. The
// outcomes come from kindOutcomes, the code's own table; a test keeps the
// kinds here and the parser's in step.

// Kinds lists every step kind, in the order the reference gives them.
var Kinds = []string{KindAgent, KindCommand, KindWait, KindLeaves, KindChecks, KindSignoff, KindWorkflow, KindArbiter}

// KeyDoc is one key of a step.
type KeyDoc struct {
	Key      string
	Desc     string
	Required bool
}

// KindDoc describes one step kind.
type KindDoc struct {
	Kind    string
	Summary string
	// Keys are the keys the kind takes, beside the shared ones.
	Keys []KeyDoc
	// Tables are the sub-tables the kind takes, as written in the file.
	Tables []KeyDoc
	// Outcomes are the outcomes the kind produces: the keys of its on table.
	Outcomes []string
}

// SharedKeys are the keys any step may set, where its kind allows.
var SharedKeys = []KeyDoc{
	{Key: "kind", Desc: "The step kind: agent, command, wait, leaves, checks, signoff, workflow or arbiter", Required: true},
	{Key: "on.<outcome>", Desc: "The step to go to for an outcome, or done or failed to end the run"},
	{Key: "max_visits", Desc: "How often a run may enter the step; 0 or unset means once"},
	{Key: "stall", Desc: "How long the watchdog lets the step go without progress before it acts, such as \"20m\"; unset takes the watchdog's default for the kind"},
}

// FileKeys are the top-level keys of a workflow file.
var FileKeys = []KeyDoc{
	{Key: "name", Desc: "The workflow's name; unset takes the file's base name"},
	{Key: "start", Desc: "The step the run enters first; unset takes the first step in the file"},
	{Key: "autonomous", Desc: "true or false: set every autonomy stage to auto or owner at once; not with [autonomy]"},
	{Key: "autonomy", Desc: "A table giving a mode, auto or owner, to any of the stages plan, accept, ship and signoff; a stage left out follows rhiza.autonomy.<stage>"},
	{Key: "step.<name>", Desc: "A step, by its name: lowercase words joined by hyphens, at most 30 characters, never done or failed", Required: true},
}

var (
	kAgent     = KeyDoc{Key: "agent", Desc: "The agent adapter's name; unset takes the run's, else the default"}
	kPerm      = KeyDoc{Key: "permission_mode", Desc: "The agent's permission mode; unset takes the run's"}
	kModels    = KeyDoc{Key: "models", Desc: "Models in order of the attempt, visit or cycle; the last stays for every later one; unset takes the agent's default"}
	kInterval  = KeyDoc{Key: "interval", Desc: "How often the step checks, such as \"10s\""}
	kSplit     = KeyDoc{Key: "split", Desc: "How the panes lay out: auto (by the tab's shape, the default), right or below"}
	kRolePrmpt = KeyDoc{Key: "prompt", Desc: "The role's first prompt to its agent", Required: true}
)

// KindDocs describes every step kind, in Kinds order.
func KindDocs() []KindDoc {
	var out []KindDoc
	for _, k := range Kinds {
		d := kindDocs[k]
		d.Kind = k
		d.Outcomes = slices.Clone(kindOutcomes[k])
		out = append(out, d)
	}
	return out
}

// RoleKeys are the keys of a role or check sub-table.
func RoleKeys() []KeyDoc { return []KeyDoc{kRolePrmpt, kAgent, kPerm, kModels} }

// PromptPlaceholder names the placeholders a kind of prompt may use.
type PromptPlaceholder struct {
	Prompt string
	Names  []string
}

// PromptPlaceholders lists the placeholders for each kind of prompt.
func PromptPlaceholders() []PromptPlaceholder {
	return []PromptPlaceholder{
		{"agent step prompt", placeholders},
		{"leaves builder and verifier, checks verifier", leafPlaceholders},
		{"checks check and critic", checkPlaceholders},
		{"checks fixer", fixerPlaceholders},
	}
}

var kindDocs = map[string]KindDoc{
	KindAgent: {
		Summary: "Runs one agent session with a prompt. It passes when the session says it is done (`rhiza session done`) or exits 0, and fails when it exits otherwise.",
		Keys: []KeyDoc{
			{Key: "prompt", Desc: "The agent's first prompt, a template over {change}, {slug}, {title}, {run} and {step}", Required: true},
			kAgent, kPerm, kModels,
			{Key: "until", Desc: "A shell command: the step also passes when the session is idle and it exits 0, leaving the session open"},
			{Key: "until_when", Desc: "With until: idle (check whenever the session sits idle, the default) or done (only once the agent says it is done)"},
			{Key: "interval", Desc: "How often the session is checked; 10s, or 30s with until"},
		},
	},
	KindCommand: {
		Summary: "Runs a shell command. It passes on exit 0 and fails otherwise.",
		Keys:    []KeyDoc{{Key: "run", Desc: "The shell command", Required: true}},
	},
	KindWait: {
		Summary: "Waits for a condition. It passes when the condition holds and times out otherwise.",
		Keys: []KeyDoc{
			{Key: "until", Desc: "A shell command; exit 0 means the condition holds", Required: true},
			{Key: "timeout", Desc: "How long to wait, such as \"30m\"", Required: true},
			{Key: "interval", Desc: "How often the condition is checked; unset is 10s"},
		},
	},
	KindLeaves: {
		Summary: "Conducts the run change set's leaves: a builder session and then a verifier session for each leaf, retrying a failed leaf. It passes when every leaf passes.",
		Keys:    []KeyDoc{kSplit, kInterval},
		Tables: []KeyDoc{
			{Key: "step.<name>.builder", Desc: "The builder role: prompt, agent, permission_mode, models by attempt", Required: true},
			{Key: "step.<name>.verifier", Desc: "The verifier role, with the same keys", Required: true},
		},
	},
	KindChecks: {
		Summary: "Runs a review's checks, each in a fresh session, all in parallel, then a further cycle of the checks that failed, with a critic judging findings and fixers fixing them. It passes when the checks pass.",
		Keys:    []KeyDoc{kSplit, kInterval},
		Tables: []KeyDoc{
			{Key: "step.<name>.check", Desc: "An array of tables, one per check: name (unique in the step) plus the role keys; at least one", Required: true},
			{Key: "step.<name>.critic", Desc: "The critic role", Required: true},
			{Key: "step.<name>.fixer", Desc: "The fixer role", Required: true},
			{Key: "step.<name>.verifier", Desc: "The verifier role, for a leaf a fixer fixed", Required: true},
		},
	},
	KindSignoff: {
		Summary: "A review's last step. It runs nothing: it stops the run until the owner signs off. `rhiza run approve` passes it; `rhiza run request-changes` gives the outcome changes.",
		Keys:    []KeyDoc{{Key: "on.changes", Desc: "The checks step that fixes what the owner asked for", Required: true}},
	},
	KindWorkflow: {
		Summary: "Runs another workflow's steps inside this run. The file is expanded when parsed: the steps come in named `<workflow>/<step>`. It passes when the inner workflow ends done and fails when it ends failed.",
		Keys: []KeyDoc{
			{Key: "workflow", Desc: "The name of the workflow to include", Required: true},
			{Key: "skip", Desc: "Steps of it to take as passed"},
		},
	},
	KindArbiter: {
		Summary: "Not a state of the machine: no step goes to it. In an autonomous run the engine calls it to decide a question an agent step raised, and its outcome goes back to the step that asked. A file declares at most one, only to set its agent, permission mode and interval; its prompt is fixed. Without one a built-in arbiter is used.",
		Keys:    []KeyDoc{kAgent, kPerm, kInterval},
	},
}
