package api

// This file holds question.answer: the owner answers an open question that
// blocks a change set. The owner's words are written as a Decision quoting
// them, linked answers to the question, in the change set the question
// blocks (--change or the current one). Once folded, the question is
// answered and drops out of ready and the inbox. With --run it answers a
// question a plan run stopped on at its gate instead: in the run's
// worktree, into the plan change set that creates the question.

import (
	"reflect"
	"slices"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
)

func init() {
	register(Operation{
		Name:     "question.answer",
		Summary:  "Answer an open question that blocks a change set: write the owner's words as a Decision quoting them, linked answers to the question, into --change",
		Write:    true,
		Request:  reflect.TypeFor[QuestionAnswerRequest](),
		Response: reflect.TypeFor[WriteResponse](),
		Run:      runQuestionAnswer,
	})
}

const questionAnswerNext = "rhiza describe question.answer"

// QuestionAnswerRequest names the question and the owner's words.
type QuestionAnswerRequest struct {
	Question string `json:"question" arg:"0" desc:"The question: an ID, a unique ID prefix of 4+ characters, or a live question's handle"`
	Text     string `json:"text" arg:"1" desc:"The owner's words, kept verbatim as the Decision's quote; must not be blank"`
	Run      string `json:"run,omitempty" flag:"run" desc:"A plan run stopped at its gate on this question: the Decision goes into the run's worktree, in the plan change set that creates the question (the run's change set when none does), and the question need block nothing; --change is ignored"`
}

// answerHandle is the handle of the Decision answering question handle q:
// answer-<q>, with a -2, -3 suffix when that is taken.
func answerHandle(m Model, q string) string {
	base := "answer-" + q
	if len(base) > maxHandle-4 {
		base = strings.TrimRight(base[:maxHandle-4], "-")
	}
	return uniqueHandle(base, func(h string) bool {
		_, taken := m.HandleTaken(h)
		return taken
	})
}

// runQuestionAnswer checks the request and writes the Decision and its
// answers link as one entry through the shared write pipeline.
func runQuestionAnswer(ctx *Ctx, req any) (any, error) {
	r := req.(QuestionAnswerRequest)
	ref := strings.TrimSpace(r.Question)
	if ref == "" {
		return nil, Usagef("question answer needs a question").WithNext(questionAnswerNext)
	}
	text := strings.TrimSpace(r.Text)
	if text == "" {
		return nil, Usagef("the answer is blank: say something").WithNext(questionAnswerNext)
	}
	gate := strings.TrimSpace(r.Run) != ""
	if gate {
		var err error
		if ctx, err = gateAnswerCtx(ctx, strings.TrimSpace(r.Run), ref); err != nil {
			return nil, err
		}
	}
	return runWriteOn(ctx, "question.answer", commandLine("question answer", ref, text), func(w *workspace, cs id.ID) (edit, error) {
		m := w.model(cs)
		q, err := m.Resolve(ref)
		if err != nil {
			if e := AsError(err); e != nil && e.Code == CodeNotFound {
				return edit{}, Usagef("no question %q", ref).WithNext("rhiza ls --kind question")
			}
			return edit{}, err
		}
		if q.Kind != op.KindQuestion {
			return edit{}, Usagef("%s is a %s, not a question", ref, q.Kind).WithNext("rhiza ls --kind question")
		}
		if !cs.IsZero() && !gate {
			blocks, err := questionBlocks(ctx, cs, q.ID)
			if err != nil {
				return edit{}, err
			}
			if !blocks {
				return edit{}, Usagef("%s is not an open question blocking this change set: it is answered, or blocks something else", ref).WithNext("rhiza inbox")
			}
		}
		handle := answerHandle(m, q.Handle)
		items := []applyItem{
			{op: "add", build: func(b writeBase) (edit, error) {
				return buildAdd(b, AddRequest{Kind: string(op.KindDecision), Handle: handle, Fields: []string{
					"statement=" + text,
					"quote=" + text,
					"quote_where=the owner, answering question " + q.Handle,
				}, lenient: true})
			}},
			{op: "link", build: func(b writeBase) (edit, error) {
				return buildLink(b, LinkRequest{From: handle, Type: "answers", To: q.Handle})
			}},
		}
		return buildApply(w, cs, ctx.IDs, items)
	})
}

// questionBlocks reports whether question q is open and blocks change set
// cs by the rule ready and the inbox use: a blocks link to cs, or to a node
// it touches or one of that node's ancestors.
func questionBlocks(ctx *Ctx, cs, q id.ID) (bool, error) {
	r, err := openReader(ctx, "question.answer", "", toBeView)
	if err != nil {
		return false, err
	}
	defer r.Close()
	in, err := planGather(r)
	if err != nil {
		return false, err
	}
	for _, p := range planCore(in).order {
		if p.info.ID != cs {
			continue
		}
		for _, b := range p.blockers {
			if b.question == q {
				return true, nil
			}
		}
	}
	return false, nil
}

// gateAnswerCtx is ctx moved into the worktree of plan run ref, which must
// be stopped at its gate for open questions, with the question q names as
// one of them, and with --change set to the plan change set whose ops
// create q (the run's change set when none does).
func gateAnswerCtx(ctx *Ctx, ref, q string) (*Ctx, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	j, err := readRun(run.NewStore(run.Dir(common)), ref)
	if err != nil {
		return nil, err
	}
	st := j.State()
	if st.Ended || !strings.HasPrefix(st.Stopped, run.OpenStop) {
		return nil, Conflictf("run %s is not stopped at its gate for open questions", j.ID).WithNext("rhiza run status " + j.ID)
	}
	if st.Worktree == "" {
		return nil, Errorf("run %s has no worktree to answer in", j.ID).WithNext("rhiza run status " + j.ID)
	}
	if ctx.Exec == nil {
		return nil, Errorf("cannot answer at a gate: this caller runs no programs")
	}
	c := *ctx
	dir := st.Worktree
	c.Root, c.IndexPath = dir, ""
	c.Git = func(args ...string) (string, error) {
		out, err := ctx.Exec("git", append([]string{"-C", dir}, args...), dir, nil)
		return string(out), err
	}
	w, err := openCtxWorkspace(&c)
	if err != nil {
		return nil, err
	}
	own, err := id.Parse(st.Change)
	if err != nil {
		return nil, Errorf("run %s: change set %q: %v", j.ID, st.Change, err)
	}
	sets := w.changeSets()
	plan, err := planChangeSets(sets, w.model(own), own)
	if err != nil {
		return nil, Errorf("run %s: %v", j.ID, err).WithNext("rhiza run status " + j.ID)
	}
	target, n, found := own, Node{}, false
	for _, cs := range plan {
		got, err := w.model(cs.ID).Resolve(q)
		if err != nil || got.Kind != op.KindQuestion {
			continue
		}
		n, found = got, true
		if slices.ContainsFunc(cs.Ops, func(o op.Op) bool { return o.Type == op.TypeNodeCreate && o.Node == got.ID }) {
			target = cs.ID
			break
		}
	}
	if !found || !slices.Contains(gateStopQuestions(st.Stopped), n.Handle) && !slices.Contains(gateStopQuestions(st.Stopped), n.ID.String()) {
		return nil, Usagef("%s is not a question run %s stopped on at its gate", q, j.ID).WithNext("rhiza run status " + j.ID)
	}
	c.ChangeRef = target.String()
	return &c, nil
}

// gateStopQuestions are the questions, by handle or ID, of a gate's stop
// detail "open questions: a, b; answer each ...".
func gateStopQuestions(stopped string) []string {
	rest, _ := strings.CutPrefix(stopped, run.OpenStop)
	rest, _, _ = strings.Cut(rest, ";")
	var qs []string
	for _, q := range strings.Split(rest, ",") {
		if q = strings.TrimSpace(q); q != "" {
			qs = append(qs, q)
		}
	}
	return qs
}
