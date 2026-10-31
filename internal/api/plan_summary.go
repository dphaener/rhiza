package api

// This file holds planSummary: a plan run's journal as the markdown body
// of the plan's pull request. It names each review round, its confirmed
// and dropped findings (a dropped one stays, with the critic's reason),
// the questions the review left and the change sets the gate moved. It
// reads no clock, so a journal and a model always render the same.

import (
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/changeset"
	"github.com/dphaener/rhiza/internal/run"
)

// planSummary renders plan run st as a pull request body. sets are the
// change sets, by which the gate's change sets are named and their state
// told.
func planSummary(st run.State, sets []changeset.ChangeSet) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Plan of run %s\n\n", st.ID)
	fmt.Fprintf(&b, "Change set %s", st.Change)
	if st.Title != "" {
		fmt.Fprintf(&b, ": %s", oneLine(st.Title))
	}
	fmt.Fprintf(&b, ". Workflow %s; the run %s.\n", st.Workflow, summaryWhere(st))

	b.WriteString("\n## Review rounds\n")
	for _, rd := range st.Rounds {
		fmt.Fprintf(&b, "\n### Round %d\n\n", rd.Number)
		if len(rd.Findings) == 0 {
			b.WriteString("No findings.\n")
		}
		for _, f := range rd.Findings {
			verdict := "no verdict"
			reason := ""
			if v, ok := rd.Verdict(f.ID); ok {
				verdict, reason = v.Verdict, ": "+oneLine(v.Reason)
			}
			fmt.Fprintf(&b, "- %s %s on %s: %s — %s%s\n", verdict, f.ID, planFindingNode(f), oneLine(f.Reason), verdict, reason)
		}
	}

	b.WriteString("\n## Questions\n\n")
	qs := reviewQuestions(st)
	if len(qs) == 0 {
		b.WriteString("None.\n")
	}
	for _, q := range qs {
		fmt.Fprintf(&b, "- %s: %s\n", q.Handle, oneLine(q.Text))
	}

	b.WriteString("\n## Gate\n\n")
	if len(st.Gated) == 0 {
		b.WriteString("The gate has not run.\n")
	}
	for _, g := range st.Gated {
		title, state := g, "unknown"
		for _, c := range sets {
			if c.ID.String() == g {
				title, state = c.Open.DisplayTitle(), string(c.State)
			}
		}
		fmt.Fprintf(&b, "- %s (%s): %s\n", g, oneLine(title), state)
	}
	return b.String()
}

// planFindingNode names what a plan finding is about.
func planFindingNode(f run.Finding) string {
	for _, s := range []string{f.Handle, f.Node, f.Place} {
		if s != "" {
			return s
		}
	}
	return f.Kind
}
