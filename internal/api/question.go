package api

// This file holds question.answer: the owner answers an open question that
// blocks a change set. The owner's words are written as a Decision quoting
// them, linked answers to the question, in the change set the question
// blocks (--change or the current one). Once folded, the question is
// answered and drops out of ready and the inbox.

import (
	"reflect"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
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
		handle := answerHandle(m, q.Handle)
		items := []applyItem{
			{op: "add", build: func(b writeBase) (edit, error) {
				return buildAdd(b, AddRequest{Kind: string(op.KindDecision), Handle: handle, Fields: []string{
					"statement=" + text,
					"quote=" + text,
					"quote_where=the owner, answering question " + q.Handle,
				}})
			}},
			{op: "link", build: func(b writeBase) (edit, error) {
				return buildLink(b, LinkRequest{From: handle, Type: "answers", To: q.Handle})
			}},
		}
		return buildApply(w, cs, ctx.IDs, items)
	})
}
