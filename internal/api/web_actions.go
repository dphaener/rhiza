package api

import (
	"net/http"
	"strings"

	"github.com/dphaener/rhiza/internal/web"
)

// WebActions are the actions the web UI offers on the project base works in:
// answering an open question a plan run stopped on (question.answer) and
// approving a sign-off (run.approve). Each runs the real operation on the
// project's log, as the owner at a terminal would, and answers in the
// owner's words. Nothing else is reachable from the web.
func WebActions(base *Ctx) web.Actions {
	ctx := func() *Ctx {
		c := *base
		c.Record, c.DryRun, c.Out, c.JSON = false, false, nil, true
		c.Stdin, c.Rex, c.Interrupts = nil, nil, nil
		return &c
	}
	// A refusal is told in the owner's words: the operation's own message
	// names runs, handles and commands (constraint owner-surfaces-no-ids).
	refuse := func(err error) error {
		status := apiStatus(err)
		if status == http.StatusInternalServerError {
			return &web.ActionError{Status: status, Message: "That did not work. The CLI has the details."}
		}
		return &web.ActionError{Status: status, Message: "Nothing to do: it is already handled, or no longer waits for you."}
	}
	return web.Actions{
		Answer: func(run, question, text string) (string, error) {
			if strings.TrimSpace(run) == "" || strings.TrimSpace(question) == "" {
				return "", &web.ActionError{Status: http.StatusBadRequest, Message: "Nothing to answer: the question is missing."}
			}
			if strings.TrimSpace(text) == "" {
				return "", &web.ActionError{Status: http.StatusBadRequest, Message: "The answer is blank: say something."}
			}
			if _, err := runQuestionAnswer(ctx(), QuestionAnswerRequest{Question: question, Text: text, Run: run}); err != nil {
				return "", refuse(err)
			}
			return "Answered. The run continues.", nil
		},
		Approve: func(run string) (string, error) {
			if strings.TrimSpace(run) == "" {
				return "", &web.ActionError{Status: http.StatusBadRequest, Message: "Nothing to approve: the run is missing."}
			}
			if _, err := runRunApprove(ctx(), RunApproveRequest{Run: run}); err != nil {
				return "", refuse(err)
			}
			return "Approved. The run continues.", nil
		},
	}
}
