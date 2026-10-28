package api

// This file holds run.summary: a review run's verdicts as markdown, from
// its journal alone. It names each check's outcome per cycle, its
// confirmed findings, the serious findings its critic dropped (they stay
// here, with the critic's reason), its minors and nits (reported, not
// fixed), the fix cycles, the code checks and the sign-offs. It reads no
// model and no clock, so a journal always renders the same.

import (
	"cmp"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/workflow"
)

func init() {
	register(Operation{
		Name:     "run.summary",
		Summary:  "Render a review run's verdicts as markdown from its journal: each check's outcome, its confirmed findings, the findings its critic dropped, the minors and nits, the fix cycles, the code checks and the sign-off",
		Request:  reflect.TypeFor[RunSummaryRequest](),
		Response: reflect.TypeFor[RunSummaryResponse](),
		Run:      runRunSummary,
	})
}

// RunSummaryRequest names the run to summarise.
type RunSummaryRequest struct {
	Run string `json:"run,omitempty" arg:"0" desc:"The run: a full ID or a unique prefix of 4+ characters; default the latest run of the current change set, else the latest run"`
}

// RunSummaryResponse is a review run's summary.
type RunSummaryResponse struct {
	Run      string `json:"run" desc:"The run's ID"`
	Change   string `json:"change" desc:"The run's change set"`
	Markdown string `json:"markdown" desc:"The summary as markdown, rendered from the run's journal alone"`
}

// Human prints the markdown as it is.
func (r RunSummaryResponse) Human(w io.Writer) {
	fmt.Fprint(w, r.Markdown)
}

// runRunSummary reads the run ref names, or the latest, finished or not,
// and renders its summary.
func runRunSummary(ctx *Ctx, req any) (any, error) {
	r := req.(RunSummaryRequest)
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	rs := run.NewStore(run.Dir(common))
	var j run.Journal
	if ref := strings.TrimSpace(r.Run); ref != "" {
		j, err = readRun(rs, ref)
	} else {
		j, err = latestRun(ctx, rs)
	}
	if err != nil {
		return nil, err
	}
	st := j.State()
	return RunSummaryResponse{Run: j.ID, Change: st.Change, Markdown: reviewSummary(st)}, nil
}

// reviewSummary renders st as a review's summary in markdown.
func reviewSummary(st run.State) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Review of run %s\n\n", st.ID)
	fmt.Fprintf(&b, "Change set %s", st.Change)
	if st.Title != "" {
		fmt.Fprintf(&b, ": %s", oneLine(st.Title))
	}
	fmt.Fprintf(&b, ". Workflow %s; the run %s.\n", st.Workflow, summaryWhere(st))

	b.WriteString("\n## Code check\n\n")
	if len(st.CodeChecks) == 0 {
		b.WriteString("None ran.\n")
	}
	for _, c := range st.CodeChecks {
		fmt.Fprintf(&b, "- step %s: %s", c.Step, c.Outcome)
		if c.Detail != "" {
			fmt.Fprintf(&b, ": %s", oneLine(c.Detail))
		}
		b.WriteString("\n")
	}

	steps, checks := summaryOrder(st)
	if len(steps) == 0 {
		b.WriteString("\n## Checks\n\nNone ran.\n")
	}
	for _, step := range steps {
		fmt.Fprintf(&b, "\n## Checks of step %s\n", step)
		for _, name := range checks[step] {
			if c, ok := st.Check(step, name); ok {
				summaryCheck(&b, st, c)
			}
		}
		summaryCycles(&b, st, step)
	}

	b.WriteString("\n## Sign-off\n\n")
	if len(st.Signoffs) == 0 && !signoffWaiting(st) {
		b.WriteString("None yet.\n")
	}
	for _, s := range st.Signoffs {
		so := runSignoff(s)
		fmt.Fprintf(&b, "- %s\n", so.line())
		for _, f := range so.Findings {
			fmt.Fprintf(&b, "  - %s %s\n", f.Check, findingLine(f))
		}
	}
	if signoffWaiting(st) {
		fmt.Fprintf(&b, "- %s (visit %d): waiting for the owner: rhiza run approve %s, or rhiza run request-changes <file> %s\n",
			st.Current, st.CurrentVisit, st.ID, st.ID)
	}
	return b.String()
}

// summaryWhere says where the run is, for the summary's first line.
func summaryWhere(st run.State) string {
	switch {
	case st.Ended:
		return st.Result
	case signoffWaiting(st):
		return "waits at step " + st.Current + " for the owner's sign-off"
	}
	return "is " + runWhere(st)
}

// summaryOrder is the checks steps, and each one's checks, in the
// workflow's order; a step or check the definition no longer names comes
// after, in journal order.
func summaryOrder(st run.State) ([]string, map[string][]string) {
	var steps []string
	checks := map[string][]string{}
	add := func(step, check string) {
		if !slices.Contains(steps, step) {
			steps = append(steps, step)
		}
		if check != "" && !slices.Contains(checks[step], check) {
			checks[step] = append(checks[step], check)
		}
	}
	var journalled []string
	for _, c := range st.Cycles {
		journalled = append(journalled, c.Step)
	}
	for _, c := range st.Checks {
		journalled = append(journalled, c.Step)
	}
	if def, err := workflow.Parse(st.Workflow, []byte(st.Definition)); err == nil {
		for _, name := range def.Order {
			if !slices.Contains(journalled, name) {
				continue
			}
			for _, c := range def.Steps[name].Checks {
				if _, ok := st.Check(name, c.Name); ok {
					add(name, c.Name)
				}
			}
			add(name, "")
		}
	}
	for _, c := range st.Checks {
		add(c.Step, c.Check)
	}
	for _, c := range st.Cycles {
		add(c.Step, "")
	}
	return steps, checks
}

// summaryCheck renders check c: its outcome, one line per cycle it ran
// with its detail (none for a plain pass), then its findings over every
// round: the serious ones standing, the serious ones its critic dropped
// and the minors and nits.
func summaryCheck(b *strings.Builder, st run.State, c run.CheckState) {
	outcome := "open"
	if o, ok := c.Latest(); ok {
		outcome = o.Outcome
	}
	fmt.Fprintf(b, "\n### %s: %s\n\n", c.Check, outcome)
	for _, cy := range st.Cycles {
		if cy.Step != c.Step || !slices.Contains(cy.Checks, c.Check) {
			continue
		}
		o, ok := c.Outcome(cy.Cycle)
		if !ok {
			fmt.Fprintf(b, "- cycle %d: open\n", cy.Cycle)
			continue
		}
		fmt.Fprintf(b, "- cycle %d: %s", cy.Cycle, o.Outcome)
		if d := strings.TrimPrefix(oneLine(o.Detail), fmt.Sprintf("cycle %d: ", cy.Cycle)); d != "" && d != "check passed" {
			fmt.Fprintf(b, ": %s", d)
		}
		b.WriteString("\n")
	}
	var confirmed, dropped, minor []string
	for _, rd := range c.Rounds {
		for _, f := range rd.Findings {
			line := fmt.Sprintf("- cycle %d: %s", rd.Cycle, findingLine(runFinding(f, rd)))
			v, judged := rd.Verdict(f.ID)
			switch {
			case !run.Serious(f.Severity):
				minor = append(minor, line)
			case judged && v.Verdict == run.Dropped:
				dropped = append(dropped, line)
			case !judged:
				confirmed = append(confirmed, line+" — no verdict, so it stands")
			default:
				confirmed = append(confirmed, line)
			}
		}
	}
	for _, part := range []struct {
		head  string
		lines []string
	}{
		{"Confirmed", confirmed},
		{"Dropped by the critic", dropped},
		{"Minors and nits, reported, not fixed", minor},
	} {
		if len(part.lines) > 0 {
			fmt.Fprintf(b, "\n%s:\n\n%s\n", part.head, strings.Join(part.lines, "\n"))
		}
	}
}

// summaryCycles renders checks step step's fix cycles: every cycle after
// the first, and any that raised questions, with what it re-ran, the
// owner's findings in an owner cycle, its leaf fixers and their
// verifiers, its change-set fixer and the questions raised once it ended
// with no fix cycle left.
func summaryCycles(b *strings.Builder, st run.State, step string) {
	var lines []string
	for _, cy := range st.Cycles {
		qs := st.ReviewQuestionsFor(step, cy.Cycle)
		if cy.Step != step || (cy.Cycle == 1 && cy.Kind == "" && len(qs) == 0) {
			continue
		}
		head := fmt.Sprintf("- cycle %d", cy.Cycle)
		switch {
		case cy.Owner():
			head += ", the owner's changes"
		case cy.Answered():
			head += ", answered by " + strings.Join(cy.Decisions, ", ")
		}
		lines = append(lines, head+": re-runs "+strings.Join(cy.Checks, ", "))
		for _, f := range cy.Findings {
			lines = append(lines, "  - owner's "+f.Check+" "+findingLine(runFinding(f, run.Round{})))
		}
		for _, l := range st.Leaves {
			if l.Step != step {
				continue
			}
			if _, ok := l.Session(workflow.RoleFixer, cy.Cycle); !ok {
				continue
			}
			line := "  - leaf " + l.Handle + ": fixer " + summaryJudged(l, workflow.RoleFixer, cy.Cycle)
			if _, ok := l.Session(workflow.RoleVerifier, cy.Cycle); ok {
				line += "; verifier " + summaryJudged(l, workflow.RoleVerifier, cy.Cycle)
			}
			lines = append(lines, line)
		}
		if f, ok := st.ReviewFix(step, cy.Cycle); ok {
			line := "  - change set: fixer "
			if !f.Judged {
				line += "open"
			} else {
				line += summaryOutcome(f.Outcome, f.Commit, f.Detail)
			}
			lines = append(lines, line)
		}
		for _, q := range qs {
			line := "  - question " + cmp.Or(q.QuestionHandle, q.Question) + " on check " + q.Check
			if q.Finding != "" {
				line += " " + q.Finding
			}
			if q.Handle != "" {
				line += ", blocks leaf " + q.Handle
			}
			lines = append(lines, line)
		}
	}
	if len(lines) > 0 {
		fmt.Fprintf(b, "\nFix cycles:\n\n%s\n", strings.Join(lines, "\n"))
	}
}

// summaryJudged is how role's session on cycle of leaf l was judged:
// open, or pass or fail with its commit and, for a fail, its detail.
func summaryJudged(l run.LeafState, role string, cycle int) string {
	a, ok := l.Judged(role, cycle)
	if !ok {
		return "open"
	}
	return summaryOutcome(a.Outcome, a.Commit, a.Detail)
}

// summaryOutcome is a fixer's or verifier's outcome, then its commit when
// it passed, or its detail when it failed.
func summaryOutcome(outcome, commit, detail string) string {
	if outcome == workflow.OutcomePass {
		if commit != "" {
			return outcome + ", commit " + shortRev(commit)
		}
		return outcome
	}
	if detail != "" {
		return outcome + " (" + oneLine(detail) + ")"
	}
	return outcome
}
