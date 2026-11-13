package web

import (
	"context"
	"net/http"
	"net/url"
	"strings"
)

// Actions runs what the owner does from the web UI. Each call answers with
// what to tell the owner; a *ActionError carries the status to send with a
// refusal. A nil func leaves the action off (404).
type Actions struct {
	// Answer answers the open question (its handle) that the plan run
	// stopped on, with the owner's words.
	Answer func(run, question, text string) (string, error)
	// Approve signs off the run waiting at a sign-off step.
	Approve func(run string) (string, error)
	// Pause pauses (paused true) or resumes the daemon (scope "daemon")
	// or the project the page is of (scope "project").
	Pause func(scope string, paused bool) (string, error)
}

// ActionError is a refusal an action reports: Status is the HTTP status
// (400 when zero), Message what the owner reads.
type ActionError struct {
	Status  int
	Message string
}

func (e *ActionError) Error() string { return e.Message }

type actionsFuncKey struct{}

// WithActions serves h with the actions a is the runner of. Only a request
// the token gate passed (WithToken) may reach one.
func WithActions(a Actions, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), actionsFuncKey{}, a)))
	})
}

// Paths the action forms post to, under the project's prefix.
const (
	AnswerPath  = "/actions/answer"
	ApprovePath = "/actions/approve"
	PausePath   = "/actions/pause"
)

// formSpec is how a page builds its rows' action forms: none without the
// CSRF value (the request carried no token).
type formSpec struct {
	CSRF string
	// Base is the project prefix the actions post under, empty when the
	// page is itself served in the project (its links are rewritten).
	Base string
}

func formsOf(r *http.Request, base string) formSpec {
	return formSpec{CSRF: CSRF(r), Base: base}
}

// of is the forms of an inbox row: an answer form on an open question that
// names its handle, an approve form on a sign-off.
func (f formSpec) of(it InboxItemView) []RowForm {
	if f.CSRF == "" || it.Run == "" {
		return nil
	}
	switch {
	case it.Kind == "open" && it.Question != "":
		return []RowForm{{Kind: "answer", Action: f.Base + AnswerPath, Run: it.Run, Question: it.Question, CSRF: f.CSRF}}
	case it.Kind == "signoff":
		return []RowForm{{Kind: "approve", Action: f.Base + ApprovePath, Run: it.Run, CSRF: f.CSRF}}
	}
	return nil
}

// ActionResult is what the confirmation page shows.
type ActionResult struct {
	Message string
	Failed  bool
	// Back is where the link under the message goes.
	Back string
}

func serveAction(w http.ResponseWriter, r *http.Request, run func(a Actions) (string, error)) {
	a, _ := r.Context().Value(actionsFuncKey{}).(Actions)
	if !ActionsOn(r) {
		http.Error(w, "actions are off", http.StatusForbidden)
		return
	}
	msg, err := run(a)
	back := "/"
	if sc := scopeOf(r); sc != nil {
		back = sc.Prefix + "/"
	} else if ref := r.Referer(); ref != "" {
		if u, perr := url.Parse(ref); perr == nil && u.Host == r.Host && strings.HasPrefix(u.Path, "/") && !strings.HasPrefix(u.Path, "//") {
			back = u.Path
		}
	}
	res := ActionResult{Message: msg, Back: back}
	status := http.StatusOK
	if err != nil {
		res.Failed, res.Message, status = true, err.Error(), http.StatusBadRequest
		if ae, ok := err.(*ActionError); ok && ae.Status != 0 {
			status = ae.Status
		}
	}
	renderStatus(w, r, status, "action", page{Title: "Action", Nav: "overview", NoDaemon: true, Data: res})
}

func serveAnswer(w http.ResponseWriter, r *http.Request) {
	serveAction(w, r, func(a Actions) (string, error) {
		if a.Answer == nil {
			return "", &ActionError{Status: http.StatusNotFound, Message: "answering is not available here"}
		}
		return a.Answer(r.PostFormValue("run"), r.PostFormValue("question"), r.PostFormValue("text"))
	})
}

func serveApprove(w http.ResponseWriter, r *http.Request) {
	serveAction(w, r, func(a Actions) (string, error) {
		if a.Approve == nil {
			return "", &ActionError{Status: http.StatusNotFound, Message: "approving is not available here"}
		}
		return a.Approve(r.PostFormValue("run"))
	})
}

func servePause(w http.ResponseWriter, r *http.Request) {
	serveAction(w, r, func(a Actions) (string, error) {
		if a.Pause == nil {
			return "", &ActionError{Status: http.StatusNotFound, Message: "pausing is not available here"}
		}
		scope := r.PostFormValue("scope")
		if scope != "daemon" && scope != "project" {
			return "", &ActionError{Status: http.StatusBadRequest, Message: "Nothing to pause: say the daemon or the project."}
		}
		var paused bool
		switch r.PostFormValue("paused") {
		case "true":
			paused = true
		case "false":
		default:
			return "", &ActionError{Status: http.StatusBadRequest, Message: "Say whether to pause or resume."}
		}
		return a.Pause(scope, paused)
	})
}
