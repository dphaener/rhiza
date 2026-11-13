package run

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// BriefLeaf is one leaf of the change set as brief.md lists it.
type BriefLeaf struct {
	Handle     string
	Statement  string
	Acceptance string
}

// BriefInput is what brief.md is rendered from. The same input renders the
// same bytes.
type BriefInput struct {
	ChangeTitle  string
	ChangeID     string
	Intent       string
	Leaves       []BriefLeaf
	Decisions    []string // statements
	Constraints  []string // statements
	PrePushCheck string   // rhiza.ship.prePushCheck, "" when unset
	LocalGate    string   // rhiza.ship.localGate, "" when unset
}

// ContractLastLine ends brief.md's contract summary.
const ContractLastLine = "Full contract: rhiza guide contract"

// RenderBrief renders a session's brief.md: the change set, its leaves, the
// decisions, the constraints, the check commands and a short contract
// summary, in that order.
func RenderBrief(in BriefInput) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Brief\n\n## Change set\n\n- Title: %s\n- ID: %s\n- Intent: %s\n", in.ChangeTitle, in.ChangeID, in.Intent)

	b.WriteString("\n## Leaves\n")
	if len(in.Leaves) == 0 {
		b.WriteString("\nnone\n")
	}
	for _, l := range in.Leaves {
		fmt.Fprintf(&b, "\n### %s\n\n%s\n\nAcceptance: %s\n", l.Handle, l.Statement, l.Acceptance)
	}

	writeList := func(title string, items []string) {
		fmt.Fprintf(&b, "\n## %s\n\n", title)
		if len(items) == 0 {
			b.WriteString("none\n")
		}
		for _, s := range items {
			fmt.Fprintf(&b, "- %s\n", s)
		}
	}
	writeList("Decisions", in.Decisions)
	writeList("Constraints", in.Constraints)

	b.WriteString("\n## Check commands\n\n")
	if in.PrePushCheck == "" && in.LocalGate == "" {
		b.WriteString("none\n")
	} else {
		if in.PrePushCheck != "" {
			fmt.Fprintf(&b, "- rhiza.ship.prePushCheck: %s\n", in.PrePushCheck)
		}
		if in.LocalGate != "" {
			fmt.Fprintf(&b, "- rhiza.ship.localGate: %s\n", in.LocalGate)
		}
	}

	b.WriteString("\n## Contract\n\n")
	b.WriteString(contractSection)
	return b.String()
}

// contractSection is the body of a brief's Contract section, which every
// session's brief.md ends with or, for a spawned child's, begins with.
const contractSection = "Build or check only what your role's prompt names, to the leaves' acceptance above.\n" +
	"Run only the tests your change touches, record what you ran, never rewrite history,\n" +
	"and end with `rhiza session done` as its own command.\n" +
	ContractLastLine + "\n"

// RenderSpawnBrief renders a spawned child's brief.md: the contract section
// only, then the caller's brief file as the caller wrote it under its own
// heading.
func RenderSpawnBrief(callerBrief string) string {
	var b strings.Builder
	b.WriteString("# Brief\n\n## Contract\n\n")
	b.WriteString(contractSection)
	b.WriteString("\n## Brief\n\n")
	b.WriteString(callerBrief)
	if !strings.HasSuffix(callerBrief, "\n") {
		b.WriteString("\n")
	}
	return b.String()
}

// WriteBriefFile creates dir and writes brief.md alone there: a spawned
// child's directory holds no slice.json. A directory or file that cannot be
// created is an error naming its path.
func WriteBriefFile(dir, brief string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create session directory %s: %w", dir, err)
	}
	path := filepath.Join(dir, "brief.md")
	if err := os.WriteFile(path, []byte(brief), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}

// SessionDirs name a session's directory under a run's. Each is the
// <session> part of <step>-<visit>/<session>.

// AgentSessionName is an agent step's session.
func AgentSessionName() string { return "agent" }

// LeafSessionName is a leaves step's builder or verifier of leaf at attempt.
func LeafSessionName(role, leaf string, attempt int) string {
	return fmt.Sprintf("%s-%s-%d", role, leaf, attempt)
}

// CheckSessionName is a checks step's check or critic of cycle, named as
// OpenCheckCheckout names its checkout.
func CheckSessionName(check, role string, cycle int) string {
	return fmt.Sprintf("%s-%s-%d", check, role, cycle)
}

// LeafFixSessionName is a checks step's leaf fixer or verifier of leaf in cycle
// at attempt.
func LeafFixSessionName(role, leaf string, cycle, attempt int) string {
	return fmt.Sprintf("%s-%s-%d-%d", role, leaf, cycle, attempt)
}

// ChangeSetFixSessionName is a checks step's change set fixer in cycle at
// attempt.
func ChangeSetFixSessionName(cycle, attempt int) string {
	return fmt.Sprintf("fixer-%d-%d", cycle, attempt)
}

// SessionDir is <commondir>/rhiza/runs/<run>/<step>-<visit>/<session>. It
// lies in the git common dir, so inside no worktree, and is kept when the
// run ends.
func SessionDir(commonDir, runID, step string, visit int, session string) string {
	return filepath.Join(commonDir, "rhiza", "runs", safeName(runID), fmt.Sprintf("%s-%d", safeName(step), visit), safeName(session))
}

// KeyedSessionDir is <commondir>/rhiza/runs/<run>/<kind>/<session id>, the
// directory of an arbiter or advisor session: no step visit, so keyed on
// the session's own ID.
func KeyedSessionDir(commonDir, runID, kind, session string) string {
	return filepath.Join(commonDir, "rhiza", "runs", safeName(runID), safeName(kind), safeName(session))
}

// safeName keeps a name to one path element.
func safeName(s string) string {
	return strings.NewReplacer("/", "_", "\\", "_").Replace(s)
}

// WriteSessionFiles creates dir and writes brief.md and slice.json there,
// the slice's bytes as given. A directory or file that cannot be created is
// an error naming its path.
func WriteSessionFiles(dir, brief string, slice []byte) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create session directory %s: %w", dir, err)
	}
	for name, data := range map[string][]byte{"brief.md": []byte(brief), "slice.json": slice} {
		path := filepath.Join(dir, name)
		if err := os.WriteFile(path, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", path, err)
		}
	}
	return nil
}

// FailedCommand is a failing gate or pre-push command and the path of its
// kept output.
type FailedCommand struct {
	Command    string
	OutputPath string
}

// FailedCI is a failed CI run and its saved log file.
type FailedCI struct {
	Run     string
	LogFile string
}

// HandbackCause is why a step handed back: the failing step and its
// outcome, and whichever of the failing commands, failed CI runs and
// conflicted files the cause has.
type HandbackCause struct {
	Step      string
	Outcome   string
	Commands  []FailedCommand
	CIRuns    []FailedCI
	Conflicts []string
}

// RenderHandback renders handback.md for cause.
func RenderHandback(c HandbackCause) (string, error) {
	if strings.TrimSpace(c.Step) == "" {
		return "", errors.New("hand-back cause names no failing step")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Hand-back\n\n- Failing step: %s\n- Outcome: %s\n", c.Step, c.Outcome)
	if len(c.Commands) > 0 {
		b.WriteString("\n## Failing commands\n\n")
		for _, f := range c.Commands {
			fmt.Fprintf(&b, "- `%s`: output in %s\n", f.Command, f.OutputPath)
		}
	}
	if len(c.CIRuns) > 0 {
		b.WriteString("\n## Failed CI runs\n\n")
		for _, f := range c.CIRuns {
			fmt.Fprintf(&b, "- %s: log in %s\n", f.Run, f.LogFile)
		}
	}
	if len(c.Conflicts) > 0 {
		b.WriteString("\n## Conflicted files\n\n")
		for _, f := range c.Conflicts {
			fmt.Fprintf(&b, "- %s\n", f)
		}
	}
	return b.String(), nil
}

// WriteHandback writes handback.md for cause into dir (created if need be).
// run-handback calls it; nothing else does.
func WriteHandback(dir string, c HandbackCause) error {
	text, err := RenderHandback(c)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create session directory %s: %w", dir, err)
	}
	path := filepath.Join(dir, "handback.md")
	if err := os.WriteFile(path, []byte(text), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
