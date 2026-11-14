package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/dphaener/rhiza/internal/id"
	"github.com/dphaener/rhiza/internal/op"
	"github.com/dphaener/rhiza/internal/run"
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
			// A question already answered is not answered twice: the run
			// waits for the others, or its continuation is queued.
			if done, err := gateAnswered(ctx(), run, question); err != nil {
				return "", refuse(err)
			} else if done {
				return "", refuse(Usagef("already answered"))
			}
			if _, err := runQuestionAnswer(ctx(), QuestionAnswerRequest{Question: question, Text: text, Run: run}); err != nil {
				return "", refuse(err)
			}
			// The run resumes once every question its gate stopped on is
			// answered, as the TUI Inbox app resumes it after answering.
			open, err := gateUnanswered(ctx(), run)
			if err != nil {
				return "", &web.ActionError{Status: http.StatusInternalServerError, Message: "Answered, but the run did not continue. The CLI has the details."}
			}
			if len(open) > 0 {
				return "Answered. The run continues once its other questions are answered.", nil
			}
			if _, err := runRunResume(ctx(), RunResumeRequest{Run: run}); err != nil {
				return "", &web.ActionError{Status: apiStatus(err), Message: "Answered, but the run did not continue. The CLI has the details."}
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

// gateAnswered reports whether question q, one plan run ref stopped on at
// its gate, has a live Decision answering it in the change set
// question.answer --run writes into.
func gateAnswered(ctx *Ctx, ref, q string) (bool, error) {
	c, err := gateAnswerCtx(ctx, ref, q)
	if err != nil {
		return false, err
	}
	c.IndexPath = ""
	w, err := openCtxWorkspace(c)
	if err != nil {
		return false, err
	}
	cs, err := id.Parse(c.ChangeRef)
	if err != nil {
		return false, err
	}
	m := w.model(cs)
	n, err := m.Resolve(q)
	if err != nil {
		return false, err
	}
	for _, l := range m.LinksTo(n.ID, op.LinkAnswers) {
		if d, ok := m.Node(l.From); ok && !d.Retired && d.Kind == op.KindDecision {
			return true, nil
		}
	}
	return false, nil
}

// gateUnanswered is the questions plan run ref stopped on at its gate that
// no Decision answers yet.
func gateUnanswered(ctx *Ctx, ref string) ([]string, error) {
	common, err := gitCommonDir(ctx)
	if err != nil {
		return nil, err
	}
	j, err := readRun(run.NewStore(run.Dir(common)), ref)
	if err != nil {
		return nil, err
	}
	var open []string
	for _, q := range gateStopQuestions(j.State().Stopped) {
		done, err := gateAnswered(ctx, j.ID, q)
		if err != nil {
			return nil, err
		}
		if !done {
			open = append(open, q)
		}
	}
	return open, nil
}
