package api

import (
	"net/http"
	"strconv"
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
	setPaused := func(key string, paused bool) error {
		_, err := runDaemonSet(ctx(), DaemonSetRequest{Key: key, Value: strconv.FormatBool(paused)})
		return err
	}
	return web.Actions{
		Pause: func(scope string, paused bool) (string, error) {
			key, what := "rhiza.daemon.paused", "The daemon"
			switch scope {
			case "daemon":
			case "project":
				key, what = "rhiza.project.paused", "This project"
			default:
				return "", &web.ActionError{Status: http.StatusBadRequest, Message: "Nothing to pause: say the daemon or the project."}
			}
			if err := setPaused(key, paused); err != nil {
				return "", &web.ActionError{Status: apiStatus(err), Message: "That did not work: " + strings.ToLower(what) + " was not changed. The CLI has the details."}
			}
			if paused {
				return what + " is paused.", nil
			}
			return what + " is resumed.", nil
		},
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
