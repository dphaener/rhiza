package api

// This file holds the review-run half of session.done's refusals: a check
// session may not say done before its findings round is journaled, nor a
// critic before every serious finding of its check's latest round has a
// verdict (see runSessionDone).

import (
	"fmt"
	"strings"

	"github.com/dphaener/rhiza/internal/run"
	"github.com/dphaener/rhiza/internal/session"
)

// isCheckSession reports whether rec is a review check's session.
func isCheckSession(rec session.Record) bool {
	return rec.Run != "" && rec.Check != "" && rec.Role == "check"
}

// isCriticSession reports whether rec is a review check's critic session.
func isCriticSession(rec session.Record) bool {
	return rec.Run != "" && rec.Check != "" && rec.Role == "critic"
}

// reviewRecordMissing says why the check or critic session rec has not
// recorded its result for this visit, or "" when it has, or when that cannot
// be told (a run that is gone or over, a step that moved on): then the run
// judges it as it always did.
func reviewRecordMissing(ctx *Ctx, rec session.Record) string {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return ""
	}
	j, err := run.NewStore(run.Dir(common)).Read(rec.Run)
	if err != nil {
		return ""
	}
	st := j.State()
	if st.Ended || st.Current != rec.Step {
		return ""
	}
	cs, _ := st.Check(rec.Step, rec.Check)
	if isCheckSession(rec) {
		if _, ok := cs.RoundFor(rec.Cycle); ok {
			return ""
		}
		return fmt.Sprintf("check %s has journaled no findings round for cycle %d: run `rhiza run findings --check %s -` (with {\"findings\":[]} for a clean check) before saying done", rec.Check, rec.Cycle, rec.Check)
	}
	rd, ok := cs.LatestRound()
	if !ok {
		return ""
	}
	var open []string
	for _, f := range rd.Findings {
		if !run.Serious(f.Severity) {
			continue
		}
		if _, judged := rd.Verdict(f.ID); !judged {
			open = append(open, f.ID)
		}
	}
	if len(open) == 0 {
		return ""
	}
	return fmt.Sprintf("check %s's round %d has blocker or major findings with no verdict (%s): run `rhiza run verdict --check %s -` with a verdict for each before saying done", rec.Check, rd.Number, strings.Join(open, ", "), rec.Check)
}
