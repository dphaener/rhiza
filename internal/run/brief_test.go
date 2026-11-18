package run

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func briefInput() BriefInput {
	return BriefInput{
		ChangeTitle: "Every agent session a run", ChangeID: "01ABC", Intent: "files",
		Leaves: []BriefLeaf{{Handle: "a", Statement: "does a", Acceptance: "go test a"},
			{Handle: "b", Statement: "does b", Acceptance: "go test b"}},
		Decisions:   []string{"d1"},
		Constraints: []string{"c1", "c2"},
		LocalGate:   "make gate",
	}
}

func TestBriefFiles(t *testing.T) {
	in := briefInput()
	got := RenderBrief(in)
	if got != RenderBrief(in) {
		t.Fatal("same input rendered different bytes")
	}
	last := -1
	for _, s := range []string{"## Change set", "Every agent session a run", "01ABC", "files", "### a", "go test a", "### b", "go test b",
		"## Decisions", "- d1", "## Constraints", "- c2", "## Check commands", "rhiza.ship.localGate: make gate", "## Contract"} {
		i := strings.Index(got, s)
		if i < 0 || i < last {
			t.Fatalf("%q missing or out of order in:\n%s", s, got)
		}
		last = i
	}
	if strings.Contains(got, "prePushCheck") {
		t.Error("unset prePushCheck listed")
	}
	if !strings.HasSuffix(got, "Full contract: rhiza guide contract\n") {
		t.Errorf("last line wrong:\n%s", got)
	}
	in.LocalGate = ""
	if none := RenderBrief(in); !strings.Contains(none, "## Check commands\n\nnone\n") {
		t.Errorf("no check commands should say none:\n%s", none)
	}
	in.PrePushCheck = "make self-check"
	if both := RenderBrief(in); !strings.Contains(both, "rhiza.ship.prePushCheck: make self-check") {
		t.Errorf("prePushCheck missing:\n%s", both)
	}
}

// A leaf lists a Files:, Must not: and Failure waiver: line under its
// acceptance for each contract field it has, and none for a field it lacks.
func TestBriefLeafContract(t *testing.T) {
	in := briefInput()
	in.Leaves[0].Files = "internal/a.go"
	in.Leaves[0].MustNot = "touch b"
	in.Leaves[0].FailureWaiver = "nothing can fail"
	got := RenderBrief(in)
	want := "### a\n\ndoes a\n\nAcceptance: go test a\nFiles: internal/a.go\nMust not: touch b\nFailure waiver: nothing can fail\n"
	if !strings.Contains(got, want) {
		t.Errorf("leaf a lacks its contract lines:\n%s", got)
	}
	_, b, _ := strings.Cut(got, "### b")
	b, _, _ = strings.Cut(b, "## Decisions")
	if b != "\n\ndoes b\n\nAcceptance: go test b\n\n" {
		t.Errorf("leaf b with no contract fields = %q", b)
	}
	for _, s := range []string{"Files:", "Must not:", "Failure waiver:"} {
		if strings.Count(got, s) != 1 {
			t.Errorf("%s listed %d times:\n%s", s, strings.Count(got, s), got)
		}
	}
}

func TestSessionDir(t *testing.T) {
	cd := "/repo/.git"
	dirs := map[string]string{
		"agent-1":         SessionDir(cd, "r1", "conduct", 1, AgentSessionName()),
		"agent-2":         SessionDir(cd, "r1", "conduct", 2, AgentSessionName()),
		"builder-a-1":     SessionDir(cd, "r1", "build", 1, LeafSessionName("builder", "a", 1)),
		"builder-a-2":     SessionDir(cd, "r1", "build", 1, LeafSessionName("builder", "a", 2)),
		"check":           SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "check", 1)),
		"critic-1":        SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "critic", 1)),
		"critic-2":        SessionDir(cd, "r1", "review", 1, CheckSessionName("quality", "critic", 1)),
		"check-cycle2":    SessionDir(cd, "r1", "review", 1, CheckSessionName("qa", "check", 2)),
		"leaffixer-1":     SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "a", 1, 1)),
		"leafverifier-1":  SessionDir(cd, "r1", "review", 1, LeafFixSessionName("verifier", "a", 1, 1)),
		"leaffixer-2":     SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "a", 1, 2)),
		"csfixer-1":       SessionDir(cd, "r1", "review", 1, ChangeSetFixSessionName(1, 1)),
		"csfixer-2":       SessionDir(cd, "r1", "review", 1, ChangeSetFixSessionName(1, 2)),
		"leaffixer-other": SessionDir(cd, "r1", "review", 1, LeafFixSessionName("fixer", "b", 1, 1)),
	}
	seen := map[string]string{}
	for name, d := range dirs {
		if other, ok := seen[d]; ok {
			t.Errorf("%s and %s share %s", name, other, d)
		}
		seen[d] = name
		if !strings.HasPrefix(d, filepath.Join(cd, "rhiza", "runs", "r1")) {
			t.Errorf("%s: %s outside the run's directory", name, d)
		}
	}
	if got, want := dirs["builder-a-1"], filepath.Join(cd, "rhiza/runs/r1/build-1/builder-a-1"); got != want {
		t.Errorf("builder dir = %s, want %s", got, want)
	}
}

func TestSessionDirWrite(t *testing.T) {
	root := t.TempDir()
	dir := SessionDir(root, "r1", "build", 1, LeafSessionName("builder", "a", 1))
	if err := WriteSessionFiles(dir, "brief", []byte(`{"x":1}`)); err != nil {
		t.Fatal(err)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "slice.json")); string(b) != `{"x":1}` {
		t.Errorf("slice.json = %q", b)
	}
	if b, _ := os.ReadFile(filepath.Join(dir, "brief.md")); string(b) != "brief" {
		t.Errorf("brief.md = %q", b)
	}
	file := filepath.Join(root, "file")
	os.WriteFile(file, nil, 0o644)
	bad := filepath.Join(file, "sub")
	err := WriteSessionFiles(bad, "b", nil)
	if err == nil || !strings.Contains(err.Error(), bad) {
		t.Errorf("error = %v, want one naming %s", err, bad)
	}
}

func TestHandbackFile(t *testing.T) {
	got, err := RenderHandback(HandbackCause{Step: "gate", Outcome: "failed",
		Commands: []FailedCommand{{"make gate", "/o/gate.log"}}})
	if err != nil || !strings.Contains(got, "gate") || !strings.Contains(got, "failed") ||
		!strings.Contains(got, "`make gate`") || !strings.Contains(got, "/o/gate.log") {
		t.Fatalf("gate: %v\n%s", err, got)
	}
	got, _ = RenderHandback(HandbackCause{Step: "ci", Outcome: "red", CIRuns: []FailedCI{{"run 42", "/l/42.log"}}})
	if !strings.Contains(got, "run 42") || !strings.Contains(got, "/l/42.log") {
		t.Errorf("ci:\n%s", got)
	}
	got, _ = RenderHandback(HandbackCause{Step: "merge", Outcome: "conflict", Conflicts: []string{"a.go", "b/c.go"}})
	if !strings.Contains(got, "- a.go") || !strings.Contains(got, "- b/c.go") {
		t.Errorf("conflict:\n%s", got)
	}
	if _, err := RenderHandback(HandbackCause{Outcome: "failed"}); err == nil {
		t.Error("cause with no failing step accepted")
	}
	dir := filepath.Join(t.TempDir(), "s")
	if err := WriteHandback(dir, HandbackCause{Step: "gate", Outcome: "failed"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "handback.md")); err != nil {
		t.Error(err)
	}
	if err := WriteHandback(dir, HandbackCause{}); err == nil {
		t.Error("empty cause written")
	}
}

// A spawned child's brief.md ends with the notes given, oldest first,
// after the caller's brief; with none it has no notes heading, and the
// contract section is the same either way.
func TestRenderSpawnBrief(t *testing.T) {
	notes := []Note{{Text: "first note"}, {Text: "second note"}}
	got := RenderSpawnBrief("Build leaf a.", notes)
	want := "# Brief\n\n## Contract\n\n" + contractSection + "\n## Brief\n\nBuild leaf a.\n\n" +
		"Notes from earlier sessions of this run:\n- first note\n- second note\n"
	if got != want {
		t.Errorf("brief with notes =\n%q\nwant\n%q", got, want)
	}
	bare := RenderSpawnBrief("Build leaf a.\n", nil)
	if bare != "# Brief\n\n## Contract\n\n"+contractSection+"\n## Brief\n\nBuild leaf a.\n" || strings.Contains(bare, NotesHeading) {
		t.Errorf("brief with no notes = %q", bare)
	}
	if NotesSection(nil) != "" {
		t.Errorf("NotesSection(nil) = %q", NotesSection(nil))
	}
}

// todayBrief is brief.md as RenderBrief wrote it before a brief could
// carry a previous attempt, for briefInput with leaf a's files and a
// pre-push check.
const todayBrief = "# Brief\n\n## Change set\n\n- Title: Every agent session a run\n- ID: 01ABC\n- Intent: files\n\n## Leaves\n\n### a\n\ndoes a\n\nAcceptance: go test a\nFiles: internal/a.go\n\n### b\n\ndoes b\n\nAcceptance: go test b\n\n## Decisions\n\n- d1\n\n## Constraints\n\n- c1\n- c2\n\n## Check commands\n\n- rhiza.ship.prePushCheck: make self-check\n- rhiza.ship.localGate: make gate\n\n## Contract\n\nBuild or check only what your role's prompt names, to the leaves' acceptance above.\nRun only the tests your change touches, record what you ran, never rewrite history,\nand end with `rhiza session done` as its own command.\nFull contract: rhiza guide contract\n"

// A brief with a previous attempt carries a Previous attempt section after
// the change set: the failed step, outcome, failing command, output tail
// and the round's findings with check, severity and text; (failure case)
// one without it renders byte for byte as before.
func TestBriefPreviousAttempt(t *testing.T) {
	in := briefInput()
	in.Leaves[0].Files = "internal/a.go"
	in.PrePushCheck = "make self-check"
	if got := RenderBrief(in); got != todayBrief {
		t.Fatalf("a brief with no previous attempt changed:\n%q", got)
	}
	in.Previous = &PreviousAttempt{
		Cause: HandbackCause{Step: "check", Outcome: "fail", Commands: []FailedCommand{{Command: "rhiza run review-check"}},
			Output: "refused: 1 major finding open"},
		Findings: []AttemptFinding{{Check: "plan review", Severity: "major", Text: "leaf x has no acceptance"},
			{Check: "code", Text: "unused import"}},
	}
	got := RenderBrief(in)
	want := "- Intent: files\n\n## Previous attempt\n\nThe run failed and was retried; this is how the attempt before this one ended.\n\n" +
		"- Failed step: check\n- Outcome: fail\n- Failing command: `rhiza run review-check`\n\n" +
		"### Output tail\n\n```\nrefused: 1 major finding open\n```\n\n" +
		"### Findings of the latest round\n\n- plan review, major: leaf x has no acceptance\n- code: unused import\n\n## Leaves\n"
	if !strings.Contains(got, want) {
		t.Errorf("brief lacks the previous attempt %q:\n%s", want, got)
	}
	if strings.Replace(got, want, "- Intent: files\n\n## Leaves\n", 1) != todayBrief {
		t.Errorf("the section changed the rest of the brief:\n%s", got)
	}
	in.Previous = &PreviousAttempt{Cause: HandbackCause{Step: "ship", Outcome: "failed", Cause: `step "ship" went over max_visits 2`}}
	got = RenderBrief(in)
	if want := "- Failed step: ship\n- Outcome: failed\n\n### Cause\n\nstep \"ship\" went over max_visits 2\n\n### Findings of the latest round\n\nnone\n"; !strings.Contains(got, want) {
		t.Errorf("brief lacks %q:\n%s", want, got)
	}
}

// retryBriefWF fails at check, a review-check-shaped command step, and
// retries from review.
const retryBriefWF = `
name = "wf"
[step.review]
kind = "command"
run = "review"
on.pass = "check"

[step.check]
kind = "command"
run = "rhiza run review-check"
retry_from = "review"
`

// A run retried after failing at a command step has a previous attempt:
// the step, its outcome, the command and its output tail, and the
// findings of the latest plan review round before the retry, not an
// earlier one's or one recorded after the retry; (failure case) a run
// never retried has none.
func TestJournalPreviousAttempt(t *testing.T) {
	e, r, _, _ := setup(t, retryBriefWF)
	r.codes["check"] = []int{1, 0}
	r.details["check"] = "review-check: refused\n1 major finding open"
	findings := func(reason string) {
		if err := e.Store.Append(runID, Entry{Type: TypeFindings, At: at(1), Findings: []Finding{{ID: "f1", Severity: "major", Reason: reason}}}); err != nil {
			t.Fatal(err)
		}
	}
	findings("an older round")
	r.onCommand = func(step string) {
		if step == "check" && r.commands["check"] == 1 {
			findings("leaf x has no acceptance")
		}
	}
	if _, err := e.Run(runID); err != nil {
		t.Fatal(err)
	}
	j, err := e.Store.Read(runID)
	if err != nil {
		t.Fatal(err)
	}
	if p := j.PreviousAttempt(); p != nil {
		t.Fatalf("a run never retried has a previous attempt: %+v", p)
	}
	r.onCommand = func(step string) {
		if step == "check" {
			findings("after the retry")
		}
	}
	if st, err := e.Retry(runID); err != nil || st.Result != Succeeded {
		t.Fatalf("retry: %v, %+v", err, st)
	}
	if j, err = e.Store.Read(runID); err != nil {
		t.Fatal(err)
	}
	want := &PreviousAttempt{
		Cause: HandbackCause{Step: "check", Outcome: "fail", Commands: []FailedCommand{{Command: "rhiza run review-check"}},
			Output: "review-check: refused\n1 major finding open"},
		Findings: []AttemptFinding{{Check: PlanReviewCheck, Severity: "major", Text: "leaf x has no acceptance"}},
	}
	if got := j.PreviousAttempt(); !reflect.DeepEqual(got, want) {
		t.Errorf("previous attempt = %+v\nwant %+v", got, want)
	}
}

// A review check's round brings the latest round of every check of its
// step in that cycle, each finding under its check's name; a leaf's
// findings are no round. A retried entry naming no failed step falls back
// to the finished entry's step, and a step that never ended (it went over
// its max_visits) gives "failed" and the finished entry's detail.
func TestJournalPreviousAttemptCheckRounds(t *testing.T) {
	def := "name = \"wf\"\n[step.review]\nkind = \"agent\"\nprompt = \"p\"\non.pass = \"ship\"\n[step.ship]\nkind = \"command\"\nrun = \"ship\"\n"
	fs := func(sev, reason string) []Finding { return []Finding{{ID: reason, Severity: sev, Reason: reason}} }
	j := Journal{ID: runID, Entries: []Entry{
		{Type: TypeStarted, At: at(0), Started: &Started{Workflow: "wf", Definition: def, Change: "CS1"}},
		{Type: TypeFindings, At: at(1), Step: "review", Check: "plan", Cycle: 1, Findings: fs("major", "cycle one")},
		{Type: TypeFindings, At: at(2), Step: "review", Check: "code", Cycle: 2, Findings: fs("minor", "code two")},
		{Type: TypeFindings, At: at(3), Step: "review", Check: "plan", Cycle: 2, Findings: fs("blocker", "plan two")},
		{Type: TypeFindings, At: at(4), Step: "review", Leaf: "L1", Attempt: 1, Findings: fs("major", "a leaf's")},
		{Type: TypeFinished, At: at(5), Step: "ship", End: Failed, Detail: `step "ship" went over max_visits 2`},
		{Type: TypeRetried, At: at(6), Step: "ship", Detail: `step "ship" went over max_visits 2`},
	}}
	want := &PreviousAttempt{
		Cause: HandbackCause{Step: "ship", Outcome: Failed, Commands: []FailedCommand{{Command: "ship"}},
			Cause: `step "ship" went over max_visits 2`},
		Findings: []AttemptFinding{{Check: "plan", Severity: "blocker", Text: "plan two"}, {Check: "code", Severity: "minor", Text: "code two"}},
	}
	if got := j.PreviousAttempt(); !reflect.DeepEqual(got, want) {
		t.Errorf("previous attempt = %+v\nwant %+v", got, want)
	}
}
