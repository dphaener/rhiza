package run

import (
	"fmt"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/workflow"
)

// capItem is one question a capped checks step raises: on check's
// surviving confirmed finding, or, when finding is zero, on the check's
// failure itself.
type capItem struct {
	check   string
	finding Finding
}

// name is how the step's detail names the item.
func (it capItem) name() string {
	if it.finding.ID == "" {
		return "check " + it.check
	}
	return "check " + it.check + " finding " + it.finding.ID
}

// capItems are the questions a checks step raises once cycle, with no fix
// cycle left, ended with checks failed: for each, in order, one per
// surviving confirmed finding, those the next cycle's fix phase would
// have fixed (confirmedFindings), or one on the check's failure when it
// has none (its session, critic or fix failed, or its round was missing).
func capItems(cycle int, failed []string, states map[string]CheckState) []capItem {
	var out []capItem
	for _, name := range failed {
		n := len(out)
		confirmedFindings(cycle+1, []string{name}, states, func(f Finding) {
			out = append(out, capItem{check: name, finding: f})
		})
		if len(out) == n {
			out = append(out, capItem{check: name})
		}
	}
	return out
}

// reviewQuestionText is the statement of the question it raises, check
// having run cycles cycles and failed the last for detail: the check, the
// finding's ID, severity, place, leaf and reason, and what to settle.
func reviewQuestionText(it capItem, cycles int, detail string) string {
	if it.finding.ID == "" {
		return fmt.Sprintf("Review check %s still fails after %d cycles: %s. How should it be fixed, or should the change set change?",
			it.check, cycles, clip(detail))
	}
	f := it.finding
	what := fmt.Sprintf("finding %s (%s", f.ID, f.Severity)
	if f.Place != "" {
		what += " at " + f.Place
	}
	what += ")"
	settle := "should the change set change?"
	if f.LeafHandle != "" {
		what += " on leaf " + f.LeafHandle
		settle = "should leaf " + f.LeafHandle + " change?"
	}
	return fmt.Sprintf("Review check %s still fails after %d cycles: %s: %s. Its last cycle: %s. How should it be fixed, or %s",
		it.check, cycles, what, clip(f.Reason), clip(detail), settle)
}

// capped settles a checks step that would end failed (a's detail) with
// no fix cycle left. First it raises a question per capItem of its latest
// cycle not raised yet: on a finding naming one of the change set's
// leaves with AddLeafQuestion, linked blocks to that leaf, otherwise with
// AddReviewQuestion, each journalled as a review-question entry. One that
// cannot be added is named in the step's detail, and a later visit tries
// it again. Once every question of the cycle is raised and answered by a
// Decision, it journals an answered check-cycle over the failed checks,
// one fresh cycle past the cap, and returns again true; otherwise the
// step ends failed, its detail naming the open questions.
//
// In an autonomous run whose step has had no answered cycle yet, the
// arbiter answers each question at once (arbitrate, or its open
// arbitration rejoined), and only a pass counts. So an autonomous review
// is arbitrated once: when its answered cycle fails too, the questions
// that raises are the owner's (LeafAnswer), never the arbiter's. A
// LeafAnswer failure ends the step with an error outcome. The error
// return is for journal failures.
func (r *running) capped(cr CheckRunner, step workflow.Step, a checkAction, leaves *[]Leaf) (result, bool, error) {
	fail := result{outcome: a.outcome, detail: a.detail}
	st := r.j.State()
	cur, ok := st.LatestCycle(step.Name)
	if !ok {
		return fail, false, nil
	}
	states := r.checkStates(step.Name)
	var failed []string
	for _, name := range cur.Checks {
		if o, ok := states[name].Outcome(cur.Cycle); ok && o.Outcome == CheckFailed {
			failed = append(failed, name)
		}
	}
	items := capItems(cur.Cycle, failed, states)
	if *leaves == nil && slices.ContainsFunc(items, func(it capItem) bool { return it.finding.Leaf != "" }) {
		l, err := cr.Leaves(st.Change)
		if err != nil {
			return errResult(fmt.Errorf("read the change set's leaves for the review's questions: %v", err)), false, nil
		}
		*leaves = l
	}
	var problems []string
	for _, it := range items {
		if slices.ContainsFunc(st.ReviewQuestionsFor(step.Name, cur.Cycle), func(q ReviewQuestion) bool {
			return q.Check == it.check && q.Finding == it.finding.ID
		}) {
			continue
		}
		cycles := 0
		for _, cy := range st.Cycles {
			if cy.Step == step.Name && slices.Contains(cy.Checks, it.check) {
				cycles++
			}
		}
		o, _ := states[it.check].Outcome(cur.Cycle)
		text := reviewQuestionText(it, cycles, o.Detail)
		en := Entry{Type: TypeReviewQuestion, Step: step.Name, Check: it.check, Cycle: cur.Cycle, Finding: it.finding.ID, Detail: text}
		var q Question
		var err error
		if i := slices.IndexFunc(*leaves, func(l Leaf) bool { return it.finding.Leaf != "" && l.ID == it.finding.Leaf }); i >= 0 {
			lf := (*leaves)[i]
			en.Leaf, en.Handle = lf.ID, lf.Handle
			q, err = cr.AddLeafQuestion(lf, text)
		} else {
			q, err = cr.AddReviewQuestion(text)
		}
		if err != nil {
			problems = append(problems, fmt.Sprintf("%s: its question could not be added: %v", it.name(), err))
			continue
		}
		en.Question, en.QuestionHandle = q.ID, q.Handle
		if err := r.append(en); err != nil {
			return result{}, false, err
		}
		r.log("step %s: %s: question %s (%s)", step.Name, it.name(), q.Handle, q.ID)
		st = r.j.State()
	}
	qs := st.ReviewQuestionsFor(step.Name, cur.Cycle)
	ending := func(open []string) (result, bool, error) {
		if len(open) > 0 {
			fail.detail += "; unanswered questions: " + strings.Join(open, ", ")
		}
		if len(problems) > 0 {
			fail.detail += "; " + strings.Join(problems, "; ")
		}
		return fail, false, nil
	}
	if len(problems) > 0 || len(qs) == 0 {
		var open []string
		for _, q := range qs {
			open = append(open, q.QuestionHandle)
		}
		return ending(open)
	}
	arbiter := st.Autonomous && !slices.ContainsFunc(st.Cycles, func(c CheckCycle) bool { return c.Step == step.Name && c.Answered() })
	var ids, decisions, open []string
	for _, q := range qs {
		ids = append(ids, q.Question)
		if arbiter {
			ar, ok := r.j.State().Arbitration(q.Question)
			if !ok || ar.Open() {
				if _, err := r.arbitrate(step, q.Question); err != nil {
					return result{}, false, err
				}
				ar, _ = r.j.State().Arbitration(q.Question)
			}
			if ar.Outcome == workflow.OutcomePass {
				decisions = append(decisions, ar.Decision)
			} else {
				open = append(open, q.QuestionHandle)
			}
			continue
		}
		d, _, ok, err := cr.LeafAnswer(q.Question)
		if err != nil {
			return errResult(fmt.Errorf("check %s: read the answer to question %s: %v", q.Check, q.QuestionHandle, err)), false, nil
		}
		if !ok {
			open = append(open, q.QuestionHandle)
			continue
		}
		decisions = append(decisions, d)
	}
	if len(open) > 0 {
		return ending(open)
	}
	if err := r.append(Entry{
		Type: TypeCheckCycle, Step: step.Name, Cycle: cur.Cycle + 1, Checks: failed,
		Kind: CycleAnswered, Questions: ids, Decisions: decisions,
	}); err != nil {
		return result{}, false, err
	}
	r.log("step %s: cycle %d, answered by %s: %s", step.Name, cur.Cycle+1, strings.Join(decisions, ", "), strings.Join(failed, ", "))
	return result{}, true, nil
}

// cycleAnswers are the answers the fixers of step's cycle are given: none
// unless it is an answered cycle; then, for each question it names, the
// Decision that answered it, its statement read anew (LeafAnswer), so the
// journal holds no model text. A statement is left "" when that Decision
// no longer answers the question.
func (r *running) cycleAnswers(lr LeafRunner, step workflow.Step, cycle int) ([]ReviewAnswer, error) {
	st := r.j.State()
	i := slices.IndexFunc(st.Cycles, func(c CheckCycle) bool { return c.Step == step.Name && c.Cycle == cycle })
	if i < 0 || !st.Cycles[i].Answered() {
		return nil, nil
	}
	cy := st.Cycles[i]
	var out []ReviewAnswer
	for n, q := range cy.Questions {
		a := ReviewAnswer{Question: q}
		if n < len(cy.Decisions) {
			a.Decision = cy.Decisions[n]
		}
		for _, rq := range st.ReviewQuestions {
			if rq.Step == step.Name && rq.Question == q {
				a.QuestionHandle = rq.QuestionHandle
			}
		}
		d, statement, ok, err := lr.LeafAnswer(q)
		if err != nil {
			return nil, fmt.Errorf("read the answer to question %s: %v", q, err)
		}
		if ok && d == a.Decision {
			a.Statement = statement
		}
		out = append(out, a)
	}
	return out, nil
}
